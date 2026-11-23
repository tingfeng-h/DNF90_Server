package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

const (
	// DNFBRIDGE_PARTY_POST_SNAPSHOT_TOWN_REFRESH_SKIP controls whether
	// finishManagedPartyJoin skips the second (post-op11) town actor projection.
	// Line 70 already sends mode0+mode1(+op23) before op11/op153/op9; sending a
	// second town mode0+mode1+op23 after op11 starts P2P handshake resets the
	// remote peer actor mid-handshake and causes a red-computer failure on
	// cross-room public-town joins. Default: on; set off/0/false/no/legacy to
	// restore the old double refresh.
	partyPostSnapshotTownRefreshSkipEnv = "DNFBRIDGE_PARTY_POST_SNAPSHOT_TOWN_REFRESH_SKIP"

	// DNFBRIDGE_PARTY_SERIA_ACTOR_NO_OP23 controls whether two same-channel
	// party members who are both in their own Seria rooms receive mode0+mode1
	// (registering the peer into the active scene table 0x056A3758 so op153 HP,
	// op9 CParty* binding, party HUD right-click menu, and P2P work) while
	// omitting the trailing op23 (NOTI 0x17) so the client never renders the
	// peer inside the local Seria room. Default: on; set off/0/false/no/legacy
	// to restore the old full block.
	partySeriaActorNoOp23Env = "DNFBRIDGE_PARTY_SERIA_ACTOR_NO_OP23"
)

func partyPostSnapshotTownRefreshSkipEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyPostSnapshotTownRefreshSkipEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func partySeriaActorNoOp23Enabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partySeriaActorNoOp23Env))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// townPrivateRoomPeerProjectionPolicy decides how a same-physical-Seria-map
// party peer should be projected:
//   - blockAll=true: legacy behavior, skip both mode0+mode1 and op23.
//   - skipOp23=true: send mode0+mode1 so the peer enters the active scene table
//     (0x056A3758) for HP/party/P2P, but omit op23 so the peer is not rendered
//     inside the receiver's private Seria room.
func (s *Service) townPrivateRoomPeerProjectionPolicy(receiver *gameSession, actor onlinePlayerInfo) (blockAll bool, skipOp23 bool) {
	if !s.townPrivateRoomPeerProjectionBlocked(receiver, actor) {
		return false, false
	}
	if partySeriaActorNoOp23Enabled() {
		return false, true
	}
	return true, false
}

// refreshPostSnapshotPartyActorsAfterJoin runs at the end of
// finishManagedPartyJoin (after op11/op153/op9). When the skip switch is active,
// town receivers are NOT re-projected a second time (they were already projected
// before op11 in step 3), while active-dungeon receivers keep their idempotent
// post-snapshot projection path untouched.
func (s *Service) refreshPostSnapshotPartyActorsAfterJoin(state alignedcmd.PartyState, source string) {
	if s == nil || s.onlinePlayers == nil || state.PartyID <= 0 || len(runtimePartyMembers(state)) < 2 {
		return
	}
	if !partyPostSnapshotTownRefreshSkipEnabled() {
		s.refreshTownPartyActorsAfterJoin(state, source)
		return
	}
	dungeonReceivers := make([]*gameSession, 0, len(runtimePartyMembers(state)))
	for _, member := range runtimePartyMembers(state) {
		receiver, ok := s.onlineGameSession(member.UserID)
		if !ok || receiver == nil {
			continue
		}
		if s.sessionHasActiveDungeonRuntime(receiver) {
			dungeonReceivers = append(dungeonReceivers, receiver)
		}
	}
	if len(dungeonReceivers) == 0 {
		return
	}
	for _, receiver := range dungeonReceivers {
		if !s.recordCurrentDungeonPartyRunKey(receiver) {
			s.logGameEvent(receiver, "game-dungeon-party-run-key-after-join-failed",
				"source", source,
				"party_id", state.PartyID,
				"char_id", receiver.selectedCharacterID,
				"reason", "active_dungeon_run_or_scene_unavailable")
		}
	}
	for _, receiver := range dungeonReceivers {
		if err := s.projectDungeonPartyPeerActors(receiver, source+"_dungeon_peer_actor"); err != nil {
			s.logGameEvent(receiver, "game-dungeon-party-actor-after-join-failed",
				"source", source,
				"party_id", state.PartyID,
				"char_id", receiver.selectedCharacterID,
				"error", err)
		}
	}
}
