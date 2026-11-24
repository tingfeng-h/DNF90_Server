package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/party"
)

const townPartyKeepSingletonEnv = "DNFBRIDGE_TOWN_PARTY_KEEP_SINGLETON"

// townPartyKeepSingletonEnabled controls whether a 2 -> 1 town party departure
// (leave, kick, or disconnect) retains the surviving town player in a fresh
// 1-player party instead of leaving them solo. Dungeon 2 -> 1 departures always
// stay disbanded to avoid V212 dungeon ghost-HUD / door-lock bugs.
//   - off / 0 / false / "" (default): legacy disband behavior, preserving all
//     existing 2 -> 1 disband unit tests unchanged.
//   - on / 1 / true / yes: after sending the clean kind3/empty-kind2 departure
//     to clear the 2P HUD, seamlessly create a fresh 1-player party for the
//     surviving town player.
func townPartyKeepSingletonEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(townPartyKeepSingletonEnv))) {
	case "1", "on", "true", "yes":
		return true
	default:
		return false
	}
}

func (s *Service) retainTownSingletonPartyAfterDisband(
	survivor *gameSession,
	previous party.RuntimePartySnapshot,
	reason string,
) error {
	if s == nil || !townPartyKeepSingletonEnabled() || s.options.currentClientVersion != "official-24.9" {
		return nil
	}
	if survivor == nil || survivor.loginRevoked.Load() || len(previous.Members) != 2 {
		return nil
	}
	// Strictly town-only: never recreate a 1-player party for a player inside a dungeon.
	if s.sessionHasActiveDungeonRuntime(survivor) {
		return nil
	}
	identity, bound := s.boundGameSessionCharacterSnapshot(survivor)
	if !bound || !s.boundGameSessionCharacterCurrent(identity) {
		return nil
	}
	if _, inRaid := s.runtimeRaidForMember(identity.character); inRaid {
		return nil
	}
	manager := s.runtimePartyManagerForService()
	if manager == nil {
		return nil
	}
	if _, alreadyInParty := manager.SnapshotByUser(identity.character, 0); alreadyInParty {
		return nil
	}
	member, ok := s.runtimePartyMemberForIdentity(identity)
	if !ok {
		return nil
	}
	settings := previous.StateFor(identity.character)
	settings.PartyID = 0
	settings.UserID = identity.character
	settings.IsLeader = true
	settings.Members = nil
	if settings.MaxMembers <= 0 {
		settings.MaxMembers = 4
	}
	created := manager.Create(member, settings)
	if !created.OK {
		return nil
	}
	next := s.publishRuntimePartySnapshot(created.Party)
	s.logGameEvent(survivor, "game-town-party-singleton-retained",
		"reason", reason,
		"char_id", identity.character,
		"previous_party_id", previous.ID,
		"new_party_id", next.PartyID)
	return s.sendRuntimePartySnapshot(survivor, next)
}

// dungeonRuntime 是 runtimeDungeonState 的可读别名（测试与外部引用用，零运行时影响）。
type dungeonRuntime = runtimeDungeonState
