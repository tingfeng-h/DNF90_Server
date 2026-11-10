package dnfbridge

import "testing"

// 成员重入同一 party id 时，服务端必须把队伍搬到新 id：成员、队长、代次不变，
// 旧 id 消失，客户端因此拿到全新的队伍身份和中继端点。
func TestRebuildRuntimePartyOnRejoinKeepsMembersAndAllocatesNewID(t *testing.T) {
	f := newLeaderExitFixture(t)
	manager := f.service.runtimePartyManagerForService()
	snapshot, found := manager.SnapshotByUser(f.leader.selectedCharacterID, 0)
	if !found {
		t.Fatal("fixture party missing")
	}
	oldID := snapshot.ID
	state, ok := f.service.rebuildRuntimePartyOnRejoin(int(oldID), f.leader.selectedCharacterID)
	if !ok {
		t.Fatal("rejoin rebuild failed")
	}
	if state.PartyID == int(oldID) || state.PartyID == 0 {
		t.Fatalf("rejoin did not allocate a new party id: old=%d new=%d", oldID, state.PartyID)
	}
	if len(state.Members) != len(snapshot.Members) {
		t.Fatalf("rebuilt members=%d, want %d", len(state.Members), len(snapshot.Members))
	}
	if state.UserID != snapshot.Leader {
		t.Fatalf("rebuilt leader=%d, want %d", state.UserID, snapshot.Leader)
	}
	if _, found := manager.SnapshotByID(oldID); found {
		t.Fatalf("old party id %d survived the rebuild", oldID)
	}
	rebuilt, found := manager.SnapshotByID(uint16(state.PartyID))
	if !found || len(rebuilt.Members) != len(snapshot.Members) {
		t.Fatalf("new party %d snapshot=%+v found=%t", state.PartyID, rebuilt, found)
	}
}

func TestRecentPartyDepartureOnlyMatchesTheSameParty(t *testing.T) {
	service := &Service{}
	service.notePartyDeparture(7, 12)
	if !service.departedFromPartyRecently(7, 12) {
		t.Fatal("recent departure for the same party was not remembered")
	}
	if service.departedFromPartyRecently(7, 13) {
		t.Fatal("departure from party 12 must not match party 13")
	}
	if service.departedFromPartyRecently(8, 12) {
		t.Fatal("departure of char 7 must not match char 8")
	}
	service.clearPartyDeparture(7)
	if service.departedFromPartyRecently(7, 12) {
		t.Fatal("cleared departure still matched")
	}
}
