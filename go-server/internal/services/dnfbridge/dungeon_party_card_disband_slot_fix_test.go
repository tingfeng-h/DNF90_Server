package dnfbridge

import "testing"

func TestDungeonCardRetainsEntryPartySlotAfterMidDungeonDisband(t *testing.T) {
	plan, err := newDungeonCardRewardPlan(
		dungeonCardPlanIdentity{CharacterID: "1002", DungeonID: 100, MazeIndex: 0, RunSeed: 42},
		"test",
		dungeonCardRewardBundle{Gold: 12924},
		dungeonCardRewardBundle{Gold: 13180, Cost: 13180},
		dungeonCardRewardBundle{},
	)
	if err != nil {
		t.Fatalf("newDungeonCardRewardPlan: %v", err)
	}
	cardState, err := newDungeonCardState(plan)
	if err != nil {
		t.Fatalf("newDungeonCardState: %v", err)
	}
	runtime := &runtimeDungeonState{
		partyMemberIndex:          1,
		partyMemberIndexed:        true,
		settlementCardRewardState: cardState,
	}
	// Party already disbanded (session.party.state is empty), but 2P entered dungeon in slot 1.
	session := &gameSession{selectedCharacterID: 1002}
	session.dungeon.runtime = runtime

	// Default (enabled): flip slot stays 1, op35 populates CardSlots[1] (free) and CardSlots[5] (paid).
	if got := currentDungeonCardFlipRowStateValue(session); got != 1 {
		t.Fatalf("expected entry party slot 1 after mid-dungeon disband, got %d", got)
	}
	slots := currentDungeonClearRewardCardSlotsFromRuntime(runtime, currentDungeonCardFlipRowStateValue(session))
	if len(slots[1]) == 0 || slots[1][0].Value != 12924 {
		t.Fatalf("expected op35 CardSlots[1] to hold 2P free gold 12924, got %+v", slots[1])
	}
	if len(slots[0]) != 0 {
		t.Fatalf("expected op35 CardSlots[0] to remain empty for 2P, got %+v", slots[0])
	}

	// Rollback switch: off restores legacy 0 fallback.
	t.Setenv(dungeonPartyCardDisbandEntrySlotEnv, "off")
	if got := currentDungeonCardFlipRowStateValue(session); got != 0 {
		t.Fatalf("expected rollback fallback 0 when disabled, got %d", got)
	}
}
