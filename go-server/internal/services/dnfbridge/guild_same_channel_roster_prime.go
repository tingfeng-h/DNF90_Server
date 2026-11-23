package dnfbridge

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	dnfguild "longheng.io/server/internal/modules/dnf/guild"
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

// primeCurrentSameChannelGuildRosterForBound silently primes category-1
// cross-scene actor caches (mode8 cacheOnly) for online guild members who are
// currently on the SAME channel as the viewer, before the op140 guild roster or
// presence refresh is sent. This ensures same-channel cross-room guild members
// already have full avatar appearance and adventure club metadata when
// right-clicked in the guild panel or chat window, without creating visible
// town actors or leaking cross-channel members into local menus.
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
	route := s.currentSceneActorRoute()
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
		// Strictly same-channel only: cross-channel members must stay out of
		// local cache so their guild/chat menu remains cross-channel.
		if currentTownPresenceChannelID(peer.session) != viewerChannel {
			continue
		}
		// Peers already in the exact same town area already receive mode0.
		if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(bound.character, targetID) {
			continue
		}
		access, err := s.loadCurrentGuildMemberProfileAccess(ctx, viewer, currentGuildMemberProfileRequest{
			channel:   byte(viewerChannel),
			character: targetID,
			kind:      0,
			directory: route,
		})
		if err != nil {
			continue
		}
		if access.guildID != g.ID || !s.boundGameSessionCharacterCurrent(bound) {
			continue
		}
		if err := s.sendCurrentGuildMemberProfileWithPolicy(ctx, access, true); err != nil {
			published = errors.Join(published, err)
			continue
		}
		primed++
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
// event loop, so clicking the sender's name in chat immediately resolves to
// the same-channel menu with avatar & adventure name even if the recipient has
// never opened the guild panel after the sender entered this channel.
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
	route := s.currentSceneActorRoute()
	run := func() error {
		if !s.boundGameSessionCharacterCurrent(recipient) || !s.boundGameSessionCharacterCurrent(sender) {
			return nil
		}
		if currentTownPresenceChannelID(recipient.session) != viewerChannel || currentTownPresenceChannelID(sender.session) != viewerChannel {
			return nil
		}
		if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(recipient.character, sender.character) {
			return nil
		}
		ctx, cancel := guildChatContext(recipient.session)
		defer cancel()
		access, err := s.loadCurrentGuildMemberProfileAccess(ctx, recipient.session, currentGuildMemberProfileRequest{
			channel:   byte(viewerChannel),
			character: sender.character,
			kind:      0,
			directory: route,
		})
		if err != nil || access.guildID != guildID || !s.boundGameSessionCharacterCurrent(recipient) {
			return nil
		}
		if err := s.sendCurrentGuildMemberProfileWithPolicy(ctx, access, true); err != nil {
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
