package dnfbridge

import (
	"os"
	"strings"
)

// DNFBRIDGE_PARTY_CARD_DISBAND_ENTRY_SLOT controls whether a surviving dungeon
// member whose party disbanded mid-dungeon (e.g. 2 -> 1 leader exit) retains
// their dungeon-entry party slot (runtime.partyMemberIndex) for op35/op71 card
// reward projection.
//
// In V212 DNF.exe, SetFreeCardFlip (0x0359CC90) and SetPaidCardFlip (0x0359CDA0)
// compare op71's StateA/StateB against localPlayer->vtable[0xC80]() (the fixed
// slot assigned when entering the dungeon). After a mid-dungeon 2 -> 1 disband,
// 2P's client keeps vtable[0xC80]() == 1 while session.party.state is empty.
// Falling back to 0 causes 2P's card to miss both the local nameplate and reward
// lookup; returning runtime.partyMemberIndex keeps op35/op71 aligned with the client.
//
// Default: enabled (on). Set DNFBRIDGE_PARTY_CARD_DISBAND_ENTRY_SLOT=off to roll back.
const dungeonPartyCardDisbandEntrySlotEnv = "DNFBRIDGE_PARTY_CARD_DISBAND_ENTRY_SLOT"

func dungeonPartyCardDisbandEntrySlotEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(dungeonPartyCardDisbandEntrySlotEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// currentDungeonCardEntryPartySlotFallback returns the session's frozen
// dungeon-entry party slot when the active party has already disbanded
// mid-dungeon. Callers hold session.dungeon.mu.
func currentDungeonCardEntryPartySlotFallback(session *gameSession) int {
	if !dungeonPartyCardDisbandEntrySlotEnabled() || session == nil {
		return -1
	}
	runtime := session.dungeon.runtime
	if runtime == nil || !runtime.partyMemberIndexed {
		return -1
	}
	if int(runtime.partyMemberIndex) >= dungeonCardSlotsPerSide {
		return -1
	}
	return int(runtime.partyMemberIndex)
}
