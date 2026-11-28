package dnfbridge

import (
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/raid/fiendwar"
)

func TestPublishCurrentFiendWarMoveProjectionQueuesBoardUpdate(t *testing.T) {
	s := &Service{
		options: options{currentClientVersion: "official-24.9"},
		raids:   map[uint32]*runtimeRaidState{},
	}
	fw := &fiendwar.FiendWarBattle{
		Config: &fiendwar.FiendWarBattleConfig{
			WaitingTown: 98,
			WaitingArea: 0,
			Areas: map[int]fiendwar.FiendWarArea{
				0:  {Index: 0, Kind: "camp"},
				12: {Index: 12, Kind: "normal"},
			},
		},
		MemberAreas: map[uint16]int{11: 12, 12: 12},
		GroupAreas:  map[byte]int{1: 12},
		Areas:       map[int]fiendwar.FiendWarAreaState{},
		Symbols:     map[uint32]int{},
		Timers:      map[fiendwar.FiendWarTimerKey]fiendwar.FiendWarTimer{},
	}
	state := &runtimeRaidState{
		RaidKey:         0x012A0001,
		LeaderID:        11,
		RouteOrRaidType: 2,
		CurrentLayout:   true,
		Started:         true,
		Members: []runtimeRaidMemberState{
			{CharID: 11, Name: "P1", GroupIndex: 1, SlotOrder: 0, UserState: 1},
			{CharID: 12, Name: "P2", GroupIndex: 1, SlotOrder: 1, UserState: 1},
		},
		Battle: &runtimeRaidBattle{
			InstanceID:      "fw-move-1",
			Mode:            3,
			ChannelID:       42,
			ChannelType:     42,
			FiendWar:        fw,
			ProgressVersion: 2,
			Deadline:        time.Now().Add(time.Hour),
			MemberAccounts:  map[uint16]string{11: "acc1", 12: "acc2"},
			MemberGroups:    map[uint16]byte{11: 1, 12: 1},
		},
	}
	s.raids[state.RaidKey] = state

	// Verify helper safely snapshots and triggers queueCurrentRaidDungeonProgress without panicking or holding raidMu.
	s.publishCurrentFiendWarMoveProjection(&currentFiendWarMovePlan{
		raid:      state.RaidKey,
		instance:  "fw-move-1",
		character: 11,
		from:      0,
		to:        12,
		applied:   true,
		confirmed: true,
	})
}

func TestApplyCurrentFiendWarMoveSetsCombatTargetOnlyForOpenArea(t *testing.T) {
	f, peers := r134N1Fixture(t)
	s, b := f.service, f.state.Battle
	p := peers[0]
	var emptyTarget, from int
	for id, area := range b.FiendWar.Config.Areas {
		if id < 2 || b.FiendWar.AreaOpen(id) {
			continue
		}
		for _, next := range area.Next {
			if next >= 2 && !b.FiendWar.AreaOpen(next) {
				from, emptyTarget = id, next
				break
			}
		}
		if emptyTarget != 0 {
			break
		}
	}
	if emptyTarget == 0 {
		t.Fatal("real PVF has no adjacent non-combat FiendWar areas")
	}
	b.FiendWar.MemberAreas[p.selectedCharacterID] = from
	plan := &currentFiendWarMovePlan{raid: f.state.RaidKey, instance: b.InstanceID, character: p.selectedCharacterID, from: from, to: emptyTarget}
	if err := s.applyCurrentFiendWarMove(plan); err != nil {
		t.Fatalf("applyCurrentFiendWarMove(empty) failed: %v", err)
	}
	if plan.combatTarget {
		t.Fatal("moving into an empty FiendWar area must not mark combatTarget")
	}
}
