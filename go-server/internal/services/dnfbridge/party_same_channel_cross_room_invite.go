package dnfbridge

import (
	"context"
	"os"
	"strconv"
	"strings"

	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const partySameChannelCrossRoomInviteEnv = "DNFBRIDGE_PARTY_SAME_CHANNEL_CROSS_ROOM_INVITE"

func partySameChannelCrossRoomInviteEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partySameChannelCrossRoomInviteEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func isPartyCrossRoomInviteMode(mode byte) bool {
	return mode == 0 || mode == 12
}

// allowSameChannelCrossRoomPartyRequest decides whether an op10 party invite
// (mode=0) or join application (mode=12) may bypass the same-area (same-room)
// gate when both players are online in the exact same town channel.
// Ordinary trade (mode=1), town PvP (mode=4/13), and scene click (mode=15)
// strictly keep their existing PeerInSameArea requirement.
func (s *Service) allowSameChannelCrossRoomPartyRequest(session *gameSession, sourceID uint16, targetID uint16, mode byte) bool {
	if s == nil || session == nil || !partySameChannelCrossRoomInviteEnabled() ||
		s.options.currentClientVersion != "official-24.9" ||
		!isPartyCrossRoomInviteMode(mode) ||
		sourceID == 0 || targetID == 0 || sourceID == targetID {
		return false
	}
	if isPVPChannel(session.channel) || s.sessionHasActiveDungeonRuntime(session) {
		return false
	}
	sourceIdentity, sourceBound := s.boundGameSessionCharacterSnapshot(session)
	if !sourceBound || sourceIdentity.character != sourceID {
		return false
	}
	targetIdentity, targetBound := s.onlineGameSessionCharacterSnapshot(targetID)
	if !targetBound || targetIdentity.session == nil || targetIdentity.session == session {
		return false
	}
	if isPVPChannel(targetIdentity.session.channel) || s.sessionHasActiveDungeonRuntime(targetIdentity.session) {
		return false
	}
	if s.partyChannelMismatch(session, targetIdentity.session) {
		return false
	}
	return currentTownPresenceChannelID(session) > 0
}

// primeCrossRoomPartyInviteSender silently sends a mode8 cache frame of the
// inviter to the invitee right before forwarding NOTI7 when the two players are
// in different rooms of the same channel. With the companion client fallback on
// NOTI7 (0x02035F9D), the invitee's client resolves the inviter from category-4
// cache and displays the native party invitation dialog without requiring both
// players to stand in the same room.
func (s *Service) primeCrossRoomPartyInviteSender(sender, recipient *gameSession, senderID uint16) {
	if s == nil || sender == nil || recipient == nil || sender == recipient || senderID == 0 ||
		!partySameChannelCrossRoomInviteEnabled() || s.options.currentClientVersion != "official-24.9" {
		return
	}
	recipientID := selectedCharacterID(recipient)
	if recipientID == 0 {
		return
	}
	if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(senderID, recipientID) {
		return
	}
	owner := currentConnectionTownActorOwnerContext(recipient)
	if owner == 0 || s.partyChannelMismatch(sender, recipient) {
		return
	}
	senderBound, senderOK := s.boundGameSessionCharacterSnapshot(sender)
	recipientBound, recipientOK := s.boundGameSessionCharacterSnapshot(recipient)
	if !senderOK || !recipientOK || senderBound.character != senderID {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
	defer cancel()
	revision := currentCrossSceneAppearanceRevision(sender)
	route := s.currentSceneActorRoute()
	loadedID, name, character, found := s.characterForEnter(ctx, sender, senderID)
	if !found || loadedID != senderID || character.CharacterID != strconv.Itoa(int(senderID)) || name == "" {
		return
	}
	mode0, err := s.buildCurrentSceneObjectListBodyForSessionInContextStrict(ctx, sender, senderID, name, character, true, owner)
	if err != nil {
		s.logGameEvent(recipient, "game-party-cross-room-invite-prime-deferred",
			"inviter_char_id", senderID, "invitee_char_id", recipientID, "stage", "mode0", "error", err)
		return
	}
	var legacyRepo dnfrepo.LegacyUserInfoRepository
	if repositories, ok := s.repositoryGroupFor(ctx); ok {
		legacyRepo = repositories.LegacyUserInfo
	}
	mode3Ctx := context.WithValue(ctx, nativePeerProfileSkillsKey{}, strconv.Itoa(int(senderID)))
	if s.buffSwapReapplyEnabled() {
		mode3Ctx = context.WithValue(mode3Ctx, nativePeerBuffInspectKey{}, strconv.Itoa(int(senderID)))
	}
	mode3 := s.buildCurrentSelectedUserInfoMode3BodyInContext(mode3Ctx, sender, legacyRepo, character, true, senderID, owner)
	priming, err := buildCurrentCrossSceneUserInfoBody(mode0, mode3)
	if err != nil {
		s.logGameEvent(recipient, "game-party-cross-room-invite-prime-deferred",
			"inviter_char_id", senderID, "invitee_char_id", recipientID, "stage", "mode8", "error", err)
		return
	}
	delivered, err := s.sendCurrentCrossSceneProfileFrames(ctx, recipient, sender, [][]byte{priming}, func() bool {
		return revision == sender.crossSceneAppearanceRevision && route == s.currentSceneActorRoute() &&
			s.boundGameSessionCharacterCurrent(senderBound) && s.boundGameSessionCharacterCurrent(recipientBound) &&
			owner == currentConnectionTownActorOwnerContext(recipient) &&
			!s.partyChannelMismatch(sender, recipient)
	})
	if err != nil {
		s.logGameEvent(recipient, "game-party-cross-room-invite-prime-deferred",
			"inviter_char_id", senderID, "invitee_char_id", recipientID, "stage", "send", "error", err)
		return
	}
	if delivered {
		s.logGameEvent(recipient, "game-party-cross-room-invite-primed",
			"inviter_char_id", senderID,
			"invitee_char_id", recipientID,
			"msg_id", uint16(dnfenum.CmdPacketSetUDPIPPort),
			"mode", currentCrossSceneMode8Subtype,
			"owner_context", owner,
			"body_len", len(priming))
	}
}

// shouldPrimeCrossRoomPartyInvite reports whether a forwarded op10 request
// needs a silent mode8 pre-prime before sending NOTI7.
func (s *Service) shouldPrimeCrossRoomPartyInvite(sourceID uint16, targetID uint16, parsed party.RequestPeerRequest) bool {
	if s == nil || !partySameChannelCrossRoomInviteEnabled() || !isPartyCrossRoomInviteMode(parsed.Mode) {
		return false
	}
	return s.onlinePlayers == nil || !s.onlinePlayers.PeerInSameArea(sourceID, targetID)
}
