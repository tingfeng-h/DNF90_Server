package party

import (
	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"testing"
)

func TestDungeonLeaveTransfersLeaderWithoutReplacingInstanceParty(t *testing.T) {
	m := NewRuntimePartyManager()
	created := m.Create(testRuntimePartyMember(101, 1), alignedcmd.PartyState{MaxMembers: 4})
	m.Join(created.Party.ID, testRuntimePartyMember(202, 2))
	m.Join(created.Party.ID, testRuntimePartyMember(303, 3))
	before, _ := m.SnapshotByID(created.Party.ID)
	result := m.LeaveDungeonExpected(101, 1, created.Party.ID)
	if !result.OK || !result.LeaderMoved || result.Retired != nil || result.Disbanded || result.Party.ID != created.Party.ID || result.Party.Leader != 202 || result.Previous == nil {
		t.Fatalf("bad departure: %+v", result)
	}
	for _, member := range result.Party.Members {
		for _, old := range before.Members {
			if old.UserID == member.UserID && old.Slot != member.Slot {
				t.Fatal("survivor slots changed")
			}
		}
	}
	if result := m.LeaveDungeonExpected(202, 999, created.Party.ID); result.OK || result.Reason != "stale_session" {
		t.Fatal("stale session changed membership")
	}
	if result := m.LeaveDungeonExpected(202, 2, created.Party.ID+1); result.OK || result.Reason != "stale_party" {
		t.Fatal("stale party changed membership")
	}
	// Dungeon 2 -> 1 must disband (graduated production policy: DNF90 has no 1-player
	// dungeon party; retaining a singleton causes HP bar residue and door-lock).
	result = m.LeaveDungeonExpected(202, 2, created.Party.ID)
	if !result.OK || !result.Disbanded || result.Retired == nil {
		t.Fatalf("dungeon 2->1 must disband: %+v", result)
	}
	if result.Party.ID != 0 || len(result.Party.Members) != 0 {
		t.Fatalf("disbanded dungeon party must be empty: %+v", result.Party)
	}
	if result.Retired.ID != created.Party.ID || result.Retired.Leader != 202 || len(result.Retired.Members) != 2 {
		t.Fatalf("retired pre-disband snapshot = %+v, want id=%d leader=202 members=2", result.Retired, created.Party.ID)
	}
	if _, found := m.SnapshotByID(created.Party.ID); found {
		t.Fatal("disbanded dungeon party remained addressable")
	}
	if _, found := m.SnapshotByUser(303, 3); found {
		t.Fatal("survivor 303 remained bound after dungeon 2->1 disband")
	}
	if result := m.LeaveDungeonExpected(303, 3, created.Party.ID); result.OK || result.Reason != "not_in_party" {
		t.Fatalf("already-disbanded survivor leave = %+v, want not_in_party", result)
	}
}

// Companion contrast: unlike dungeon 2->1 (which disbands), town 2->1 retains
// the surviving member in place with the same PartyID and fixed slot.
func TestTownLeaveRetainsSingletonInPlace(t *testing.T) {
	m := NewRuntimePartyManager()
	created := m.Create(testRuntimePartyMember(101, 1), alignedcmd.PartyState{MaxMembers: 4})
	m.Join(created.Party.ID, testRuntimePartyMember(202, 2))
	result := m.Leave(202, 2)
	if !result.OK || result.Disbanded || result.Retired != nil {
		t.Fatalf("town 2->1 must retain singleton: %+v", result)
	}
	if result.Party.ID != created.Party.ID || result.Party.Leader != 101 || len(result.Party.Members) != 1 || result.Party.Members[0].Slot != 0 {
		t.Fatalf("retained town singleton = %+v, want id=%d leader=101 slot=0", result.Party, created.Party.ID)
	}
}
