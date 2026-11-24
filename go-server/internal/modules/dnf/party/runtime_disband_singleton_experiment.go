package party

import (
	"os"
	"strings"
)

const disbandOnSingletonEnv = "DNFBRIDGE_PARTY_DISBAND_ON_SINGLETON"
const townPartyKeepSingletonEnv = "DNFBRIDGE_TOWN_PARTY_KEEP_SINGLETON"

func disbandOnSingleton() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(disbandOnSingletonEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// KeepTownSingletonEnabled controls whether town 2 -> 1 departures (Leave /
// town Kick) retain the surviving member in place (keeping fixed slots) instead
// of disbanding. Dungeon 2 -> 1 (LeaveDungeonExpected and dungeon Kick) always
// follows disbandOnSingleton() and disbands.
func KeepTownSingletonEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(townPartyKeepSingletonEnv))) {
	case "1", "on", "true", "yes":
		return true
	default:
		return false
	}
}

func disbandOnSingletonTown() bool {
	if KeepTownSingletonEnabled() {
		return false
	}
	return disbandOnSingleton()
}
