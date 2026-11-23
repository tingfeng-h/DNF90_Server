package dnfbridge

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	dnfguild "longheng.io/server/internal/modules/dnf/guild"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const (
	guildRosterSameChannelPrimeEnv   = "DNFBRIDGE_GUILD_ROSTER_SAME_CHANNEL_PRIME"
	guildRosterSameChannelPrimeLimit = 64
)

func guildRosterSameChannelPrimeEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(guildRosterSameChannelPrimeEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// primeSingleSameChannelGuildMode8 builds and sends one silent mode8 cache
// frame of `peer` to `viewer` using sendCurrentCrossSceneProfileRefreshFrame.
// Crucially, it does NOT call reserveCurrentGuildParticipants(viewer, peer):
// callers such as runCurrentGuildPresenceRefresh already hold the non-reentrant
// participant reservation on `viewer`, so re-entering it would self-deadlock
// until ctx timeout and silently drop the mode8 update when a guildmate enters
// the viewer's channel.
func (s *Service) primeSingleSameChannelGuildMode8(ctx context.Context, viewer boundGameSessionCharacter, peer boundGameSessionCharacter, viewerChannel int) (bool, error) {
	if s == nil || viewer.session == nil || peer.session == nil || viewer.session == peer.session ||
		viewerChannel <= 0 || viewerChannel >= 255 {
		return false, nil
	}
	if !s.boundGameSessionCharacterCurrent(viewer) || !s.boundGameSessionCharacterCurrent(peer) {
		return false, nil
	}
	if currentTownPresenceChannelID(viewer.session) != viewerChannel ||
		currentTownPresenceChannelID(peer.session) != viewerChannel {
		return false, nil
	}
	if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(viewer.character, peer.character) {
		return false, nil
	}
	owner := byte(viewerChannel)
	revision := currentCrossSceneAppearanceRevision(peer.session)
	route := s.currentSceneActorRoute()
	loadedID, name, character, found := s.characterForEnter(ctx, peer.session, peer.character)
	if !found || loadedID != peer.character || character.CharacterID != strconv.Itoa(int(peer.character)) || name == "" {
		return false, nil
	}
	mode0, err := s.buildCurrentSceneObjectListBodyForSessionInContextStrict(ctx, peer.session, peer.character, name, character, true, owner)
	if err != nil {
		return false, err
	}
	var legacyRepo dnfrepo.LegacyUserInfoRepository
	if repositories, ok := s.repositoryGroupFor(ctx); ok {
		legacyRepo = repositories.LegacyUserInfo
	}
	mode3Ctx := context.WithValue(ctx, nativePeerProfileSkillsKey{}, strconv.Itoa(int(peer.character)))
	if buffSwapReapplyEnabled() {
		mode3Ctx = context.WithValue(mode3Ctx, nativePeerBuffInspectKey{}, strconv.Itoa(int(peer.character)))
	}
	mode3 := s.buildCurrentSelectedUserInfoMode3BodyInContext(mode3Ctx, peer.session, legacyRepo, character, true, peer.character, owner)
	priming, err := buildCurrentCrossSceneUserInfoBody(mode0, mode3)
	if err != nil {
		return false, err
	}
	return s.sendCurrentCrossSceneProfileRefreshFrame(ctx, viewer.session, peer.session, priming, func() bool {
		return revision == peer.session.crossSceneAppearanceRevision &&
			route == s.currentSceneActorRoute() &&
			s.boundGameSessionCharacterCurrent(viewer) &&
			s.boundGameSessionCharacterCurrent(peer) &&
			currentTownPresenceChannelID(viewer.session) == viewerChannel &&
			currentTownPresenceChannelID(peer.session) == viewerChannel
	}, false)
}

// primeCurrentSameChannelGuildRosterForBound silently primes category-1
// cross-scene actor caches (mode8 cacheOnly) for online guild members who are
// currently on the SAME channel as the viewer, before the op140 guild roster or
// presence refresh is sent.
func (s *Service) primeCurrentSameChannelGuildRosterForBound(ctx context.Context, bound boundGameSessionCharacter, g *dnfguild.Guild) error {
	if s == nil || g == nil || !guildRosterSameChannelPrimeEnabled() ||
		s.options.currentClientVersion != "official-24.9" ||
		!s.boundGameSessionCharacterCurrent(bound) {
		return nil
	}
	viewer := bound.session
	viewerChannel := currentTownPresenceChannelID(viewer)
	if viewerChannel <= 0 || viewerChannel >= 255 {
		return nil
	}
	viewerID := strconv.Itoa(int(bound.character))
	var (
		published error
		primed    int
	)
	for _, member := range dnfguild.SortedMembers(g) {
		if primed >= guildRosterSameChannelPrimeLimit {
			break
		}
		if member.CharacterID == viewerID {
			continue
		}
		parsedID, err := strconv.ParseUint(member.CharacterID, 10, 16)
		if err != nil || parsedID == 0 || uint16(parsedID) == bound.character {
			continue
		}
		targetID := uint16(parsedID)
		peer, online := s.onlineGameSessionCharacterSnapshot(targetID)
		if !online || peer.session == nil || peer.session == viewer ||
			s.accountIDForSession(peer.session) != member.AccountID {
			continue
		}
		if currentTownPresenceChannelID(peer.session) != viewerChannel {
			continue
		}
		if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(bound.character, targetID) {
			continue
		}
		delivered, err := s.primeSingleSameChannelGuildMode8(ctx, bound, peer, viewerChannel)
		if err != nil {
			published = errors.Join(published, err)
			continue
		}
		if delivered {
			primed++
		}
	}
	if primed > 0 {
		s.logGameEvent(viewer, "game-guild-roster-same-channel-primed",
			"guild_id", g.ID,
			"viewer_char_id", bound.character,
			"viewer_channel", viewerChannel,
			"primed_members", primed)
	}
	return published
}

// queuePrimeGuildChatSenderSameChannel asynchronously primes a single mode8
// frame of the guild chat sender onto a same-channel cross-room recipient's
// event loop.
func (s *Service) queuePrimeGuildChatSenderSameChannel(sender boundGameSessionCharacter, recipient boundGameSessionCharacter, guildID uint32) {
	if s == nil || !guildRosterSameChannelPrimeEnabled() ||
		s.options.currentClientVersion != "official-24.9" ||
		sender.session == nil || recipient.session == nil ||
		sender.session == recipient.session ||
		sender.character == 0 || recipient.character == 0 || sender.character == recipient.character {
		return
	}
	viewerChannel := currentTownPresenceChannelID(recipient.session)
	if viewerChannel <= 0 || viewerChannel >= 255 || currentTownPresenceChannelID(sender.session) != viewerChannel {
		return
	}
	if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(recipient.character, sender.character) {
		return
	}
	run := func() error {
		ctx, cancel := guildChatContext(recipient.session)
		defer cancel()
		if _, err := s.primeSingleSameChannelGuildMode8(ctx, recipient, sender, viewerChannel); err != nil {
			s.logGameEvent(recipient.session, "game-guild-chat-same-channel-prime-deferred",
				"guild_id", guildID, "sender_char_id", sender.character, "viewer_char_id", recipient.character, "error", err)
		}
		return nil
	}
	if recipient.session.events == nil || recipient.session.events.loop == nil {
		_ = run()
		return
	}
	_ = s.tryPostGameSessionCharacterEventWithDiscard(
		recipient.session,
		"guild-chat-same-channel-prime",
		recipient.character,
		recipient.generation,
		run,
		nil,
	)
}
