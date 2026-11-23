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

// primeCurrentSameChannelGuildRosterForBound silently primes category-4
// cross-scene actor caches (mode8 cacheOnly) for online guild members who are
// currently on the SAME channel as the viewer, before the op140 guild roster is
// sent. This ensures same-channel cross-room guild members already have full
// avatar appearance and adventure club metadata when right-clicked in the guild
// panel, without creating visible town actors or leaking cross-channel members
// into local menus.
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
		// local category-4 cache so their guild menu remains cross-channel.
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
