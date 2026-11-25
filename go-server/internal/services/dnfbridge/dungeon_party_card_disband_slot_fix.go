package dnfbridge

// Graduated production policy: a surviving dungeon member whose party disbanded
// mid-dungeon (e.g. 2 -> 1 leader exit) always retains their dungeon-entry party
// slot (runtime.partyMemberIndex) for op35/op71 card reward projection so it
// stays aligned with V212 localPlayer->vtable[0xC80]().

// currentDungeonCardEntryPartySlotFallback returns the session's frozen
// dungeon-entry party slot when the active party has already disbanded
// mid-dungeon. Callers hold session.dungeon.mu.
func currentDungeonCardEntryPartySlotFallback(session *gameSession) int {
	if session == nil {
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
