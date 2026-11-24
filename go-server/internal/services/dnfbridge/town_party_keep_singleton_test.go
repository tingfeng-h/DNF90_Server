package dnfbridge

import (
	"longheng.io/server/internal/modules/dnf/worldmap"
	"testing"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/channelcatalog"
	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
)

func TestTownPartyKeepSingletonRetainsTownSurvivorOnTwoToOneLeave(t *testing.T) {
	t.Setenv(townPartyKeepSingletonEnv, "1")

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
	}, alignedcmd.PartyState{NameBytes: []byte("solo-keep"), MaxMembers: 4, UserState: 1})
	if !created.OK {
		t.Fatalf("create party failed: %+v", created)
	}
	joined := manager.Join(created.Party.ID, party.RuntimePartyMember{
		UserID:            1002,
		SessionGeneration: memberID.generation,
		State:             alignedcmd.PartyMemberState{UserID: 1002, UserState: 1, HPPercent: 100, MPPercent: 100},
	})
	if !joined.OK {
		t.Fatalf("join party failed: %+v", joined)
	}
	service.publishRuntimePartySnapshot(joined.Party)

	handled, err := service.handleAlignedGameCommand(member, byte(dnfenum.GameCmdCommand), uint16(dnfenum.CmdPacketLeaveParty), nil)
	if err != nil || !handled {
		t.Fatalf("leave party handled=%t err=%v", handled, err)
	}

	// Departed member is solo; town survivor is retained in a fresh 1-player party.
	if _, found := manager.SnapshotByUser(1002, 0); found {
		t.Fatal("departed member should not remain in any party")
	}
	retained, found := manager.SnapshotByUser(1001, 0)
	if !found || len(retained.Members) != 1 || retained.Leader != 1001 {
		t.Fatalf("expected town survivor retained in 1-player party, got found=%t snapshot=%+v", found, retained)
	}
}

func TestTownPartyKeepSingletonSkipsDungeonSurvivor(t *testing.T) {
	t.Setenv(townPartyKeepSingletonEnv, "1")

	service := &Service{options: options{
		currentClientVersion: "official-24.9",
		gameUpperBodyCodec:   gameUpperBodyCodecPlain,
	}}
	channel := channelcatalog.Channel{ID: 10, Type: 11, Name: "ch.10", Port: 10010}
	survivor := &gameSession{conn: &bufferConn{}, channel: channel, selectedCharacterID: 1001}
	survivor.dungeon.runtime = &dungeonRuntime{Session: &worldmap.DungeonSession{}}
	service.bindGameSessionCharacter(survivor, 1001)

	prev := party.RuntimePartySnapshot{
		ID:     5,
		Leader: 1001,
		Members: []party.RuntimePartyMember{
			{UserID: 1001},
			{UserID: 1002},
		},
	}
	if err := service.retainTownSingletonPartyAfterDisband(survivor, prev, "dungeon_guard"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, found := service.runtimePartyManagerForService().SnapshotByUser(1001, 0); found {
		t.Fatal("dungeon survivor must not be recreated into a 1-player party")
	}
}
