package dnfbridge

import (
	"testing"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/channelcatalog"
	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
)

func TestTownPartyKeepSingletonInPlaceRetainsSamePartyAndFixedSlot(t *testing.T) {
	t.Setenv("DNFBRIDGE_TOWN_PARTY_KEEP_SINGLETON", "1")

	service := &Service{options: options{
		currentClientVersion: "official-24.9",
		gameUpperBodyCodec:   gameUpperBodyCodecPlain,
	}}
	channel := channelcatalog.Channel{ID: 10, Type: 11, Name: "ch.10", Port: 10010}
	leader := &gameSession{conn: &bufferConn{}, channel: channel, selectedCharacterID: 1001}
	member := &gameSession{conn: &bufferConn{}, channel: channel, selectedCharacterID: 1002}
	service.bindGameSessionCharacter(leader, 1001)
	service.bindGameSessionCharacter(member, 1002)

	leaderID, _ := service.boundGameSessionCharacterSnapshot(leader)
	memberID, _ := service.boundGameSessionCharacterSnapshot(member)
	manager := service.runtimePartyManagerForService()
	created := manager.Create(party.RuntimePartyMember{
		UserID:            1001,
		SessionGeneration: leaderID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1001, UserState: 1, HPPercent: 100, MPPercent: 100},
	}, alignedcmd.PartyState{NameBytes: []byte("inplace"), MaxMembers: 4, UserState: 1})
	joined := manager.Join(created.Party.ID, party.RuntimePartyMember{
		UserID:            1002,
		SessionGeneration: memberID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1002, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	service.publishRuntimePartySnapshot(joined.Party)

	// 1) Leader leaves in town: 2P (1002) becomes leader and KEEPS its fixed slot (slot 1)!
	handled, err := service.handleAlignedGameCommand(leader, byte(dnfenum.GameCmdCommand), uint16(dnfenum.CmdPacketLeaveParty), nil)
	if err != nil || !handled {
		t.Fatalf("leader leave handled=%t err=%v", handled, err)
	}
	retained, found := manager.SnapshotByUser(1002, 0)
	if !found || len(retained.Members) != 1 || retained.Leader != 1002 || retained.Members[0].Slot != 1 {
		t.Fatalf("expected 1002 retained as leader in fixed slot 1, got found=%t snapshot=%+v", found, retained)
	}

	// 2) Dungeon 2 -> 1 LeaveDungeonExpected must STILL disband.
	rejoined := manager.Join(retained.ID, party.RuntimePartyMember{
		UserID:            1001,
		SessionGeneration: leaderID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1001, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	if !rejoined.OK {
		t.Fatalf("rejoin failed: %+v", rejoined)
	}
	dungeonLeave := manager.LeaveDungeonExpected(1001, leaderID.generation, rejoined.Party.ID)
	if !dungeonLeave.OK || !dungeonLeave.Disbanded {
		t.Fatalf("dungeon 2->1 must still disband, got %+v", dungeonLeave)
	}
}
