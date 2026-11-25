package dnfbridge

import (
	"testing"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/channelcatalog"
	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
)

func TestTownPartyKeepSingletonInPlaceRetainsSamePartyAndFixedSlot(t *testing.T) {

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

	// 1) Town 2 -> 1 kick: same PartyID retained in place, leader stays in slot 0.
	kicked := manager.KickInContext(1001, leaderID.generation, 1002, memberID.generation, false)
	if !kicked.OK || kicked.Disbanded || kicked.Party.ID != created.Party.ID || len(kicked.Party.Members) != 1 || kicked.Party.Members[0].Slot != 0 {
		t.Fatalf("expected town 2->1 kick to retain same party %d in slot 0, got %+v", created.Party.ID, kicked)
	}

	// Re-join 1002 into slot 1, then leader leaves in town: 2P (1002) becomes leader and KEEPS fixed slot 1.
	rejoined := manager.Join(created.Party.ID, party.RuntimePartyMember{
		UserID:            1002,
		SessionGeneration: memberID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1002, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	if !rejoined.OK {
		t.Fatalf("rejoin failed: %+v", rejoined)
	}
	service.publishRuntimePartySnapshot(rejoined.Party)

	handled, err := service.handleAlignedGameCommand(leader, byte(dnfenum.GameCmdCommand), uint16(dnfenum.CmdPacketLeaveParty), nil)
	if err != nil || !handled {
		t.Fatalf("leader leave handled=%t err=%v", handled, err)
	}
	retained, found := manager.SnapshotByUser(1002, 0)
	if !found || len(retained.Members) != 1 || retained.Leader != 1002 || retained.Members[0].Slot != 1 {
		t.Fatalf("expected 1002 retained as leader in fixed slot 1, got found=%t snapshot=%+v", found, retained)
	}
}

func TestTownPartyKeepSingletonSkipsDungeonSurvivor(t *testing.T) {

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

	// 1) Dungeon 2 -> 1 LeaveDungeonExpected must still disband.
	created := manager.Create(party.RuntimePartyMember{
		UserID:            1001,
		SessionGeneration: leaderID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1001, UserState: 1, HPPercent: 100, MPPercent: 100},
	}, alignedcmd.PartyState{NameBytes: []byte("dungeon-leave"), MaxMembers: 4, UserState: 1})
	joined := manager.Join(created.Party.ID, party.RuntimePartyMember{
		UserID:            1002,
		SessionGeneration: memberID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1002, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	if !joined.OK {
		t.Fatalf("join failed: %+v", joined)
	}
	dungeonLeave := manager.LeaveDungeonExpected(1001, leaderID.generation, joined.Party.ID)
	if !dungeonLeave.OK || !dungeonLeave.Disbanded {
		t.Fatalf("dungeon 2->1 LeaveDungeonExpected must disband, got %+v", dungeonLeave)
	}

	// 2) Dungeon 2 -> 1 KickInContext(..., true) must still disband.
	created2 := manager.Create(party.RuntimePartyMember{
		UserID:            1001,
		SessionGeneration: leaderID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1001, UserState: 1, HPPercent: 100, MPPercent: 100},
	}, alignedcmd.PartyState{NameBytes: []byte("dungeon-kick"), MaxMembers: 4, UserState: 1})
	joined2 := manager.Join(created2.Party.ID, party.RuntimePartyMember{
		UserID:            1002,
		SessionGeneration: memberID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1002, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	if !joined2.OK {
		t.Fatalf("join2 failed: %+v", joined2)
	}
	dungeonKick := manager.KickInContext(1001, leaderID.generation, 1002, memberID.generation, true)
	if !dungeonKick.OK || !dungeonKick.Disbanded {
		t.Fatalf("dungeon 2->1 KickInContext must disband, got %+v", dungeonKick)
	}
}
