package dnfbridge

import (
	"testing"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/channelcatalog"
	"longheng.io/server/internal/modules/dnf/party"
)

func TestRaidStartGroupSyncProjectsMultiMemberGroupsAfterPartyApply(t *testing.T) {
	service := &Service{options: options{
		currentClientVersion: "official-24.9",
		gameUpperBodyCodec:   gameUpperBodyCodecPlain,
	}}
	ch := channelcatalog.Channel{ID: 23, Type: 23, Name: "anton.23", Port: 10023}
	makeSession := func(id uint16) (*gameSession, *bufferConn) {
		conn := &bufferConn{}
		sess := &gameSession{conn: conn, channel: ch, selectedCharacterID: id}
		service.bindGameSessionCharacter(sess, id)
		return sess, conn
	}
	s1, c1 := makeSession(101)
	s2, _ := makeSession(102)
	s3, _ := makeSession(201)
	s4, _ := makeSession(202)
	_ = s1
	_ = s2
	_ = s3
	_ = s4

	state := runtimeRaidState{
		RaidKey:         23<<16 | 101,
		LeaderID:        101,
		RouteOrRaidType: 0,
		NameBytes:       []byte("anton-two-groups"),
		CurrentLayout:   true,
		Started:         true,
		Members: []runtimeRaidMemberState{
			{CharID: 101, Name: "A1", GroupIndex: 1, SlotOrder: 0, RecruitingPlacementExplicit: true},
			{CharID: 102, Name: "A2", GroupIndex: 1, SlotOrder: 1, RecruitingPlacementExplicit: true},
			{CharID: 201, Name: "B1", GroupIndex: 2, SlotOrder: 0, RecruitingPlacementExplicit: true},
			{CharID: 202, Name: "B2", GroupIndex: 2, SlotOrder: 1, RecruitingPlacementExplicit: true},
		},
	}

	// 建队前：无普通队伍，投影会显示为单人（这是旧代码在 applyRuntimeRaidParties 前发 NOTI591 的根因）。
	before := service.currentRaidLivePartyProjection(state)
	for _, m := range before.Members {
		if !m.OrdinaryPartySolo {
			t.Fatalf("before applyRuntimeRaidParties expected solo=true for char %d", m.CharID)
		}
	}

	// 建队后：1队(101,102) 和 2队(201,202) 在管理器中生效，补发 NOTI591 必须恢复 1队/2队 编队号。
	if err := service.applyRuntimeRaidParties(state); err != nil {
		t.Fatalf("applyRuntimeRaidParties failed: %v", err)
	}
	c1.write.Reset()
	if err := service.syncCurrentRaidRosterAfterStartParties(state); err != nil {
		t.Fatalf("syncCurrentRaidRosterAfterStartParties failed: %v", err)
	}
	after := service.currentRaidLivePartyProjection(state)
	for _, m := range after.Members {
		if m.OrdinaryPartySolo {
			t.Fatalf("after applyRuntimeRaidParties expected solo=false for char %d (group=%d)", m.CharID, m.GroupIndex)
		}
	}
	rows := parseR195CurrentRaidRoster(t, mustFirstPacketBody(t, c1.write.Bytes(), 591), state.RaidKey)
	wantGroups := map[uint16]byte{101: 1, 102: 1, 201: 2, 202: 2}
	for id, wantGroup := range wantGroups {
		row := r195SoloPrefixRosterMember(t, rows, id)
		if row.GroupIndex != wantGroup {
			t.Fatalf("char %d NOTI591 GroupIndex=%d want %d", id, row.GroupIndex, wantGroup)
		}
	}
}

func TestRaidStartGroupSyncRetainsExplicitSingletonGroupsWhenStarted(t *testing.T) {
	service := &Service{options: options{
		currentClientVersion: "official-24.9",
		gameUpperBodyCodec:   gameUpperBodyCodecPlain,
	}}
	ch := channelcatalog.Channel{ID: 23, Type: 23, Name: "anton.23", Port: 10023}
	s1 := &gameSession{conn: &bufferConn{}, channel: ch, selectedCharacterID: 301}
	s2 := &gameSession{conn: &bufferConn{}, channel: ch, selectedCharacterID: 302}
	service.bindGameSessionCharacter(s1, 301)
	service.bindGameSessionCharacter(s2, 302)

	id1, _ := service.boundGameSessionCharacterSnapshot(s1)
	id2, _ := service.boundGameSessionCharacterSnapshot(s2)
	mgr := service.runtimePartyManagerForService()
	mgr.Create(party.RuntimePartyMember{UserID: 301, SessionGeneration: id1.generation}, alignedcmd.PartyState{MaxMembers: 4})
	mgr.Create(party.RuntimePartyMember{UserID: 302, SessionGeneration: id2.generation}, alignedcmd.PartyState{MaxMembers: 4})

	state := runtimeRaidState{
		RaidKey:         23<<16 | 301,
		LeaderID:        301,
		RouteOrRaidType: 0,
		CurrentLayout:   true,
		Started:         true,
		Members: []runtimeRaidMemberState{
			{CharID: 301, GroupIndex: 1, SlotOrder: 0, RecruitingPlacementExplicit: true, OrdinaryPartySolo: false},
			{CharID: 302, GroupIndex: 2, SlotOrder: 0, RecruitingPlacementExplicit: true, OrdinaryPartySolo: false},
		},
	}

	// 默认开启：1队1人 + 2队1人开团时，攻坚队列表保持 1队/2队，不误标为单人。
	projected := service.currentRaidLivePartyProjection(state)
	for _, m := range projected.Members {
		if m.OrdinaryPartySolo {
			t.Fatalf("explicit started singleton group char %d unexpectedly marked solo", m.CharID)
		}
	}

	// 模拟 302 显式退队（被置位 OrdinaryPartySolo=true）后，应回归单人显示。
	state.Members[1].OrdinaryPartySolo = true
	projectedAfterLeave := service.currentRaidLivePartyProjection(state)
	if projectedAfterLeave.Members[0].OrdinaryPartySolo {
		t.Fatal("char 301 should remain in group 1")
	}
	if !projectedAfterLeave.Members[1].OrdinaryPartySolo {
		t.Fatal("char 302 marked solo after leave should stay solo")
	}

	// 回滚开关：DNFBRIDGE_RAID_START_GROUP_SYNC=0 恢复旧行为（1人队按单人投影）。
	t.Setenv(raidStartGroupSyncEnv, "0")
	rolledBack := service.currentRaidLivePartyProjection(state)
	if !rolledBack.Members[0].OrdinaryPartySolo {
		t.Fatal("rollback switch should restore legacy len<2 solo projection")
	}
}

func mustFirstPacketBody(t *testing.T, raw []byte, wantMsgID uint16) []byte {
	t.Helper()
	rest := raw
	for len(rest) > 0 {
		pkt, next := splitGameServerUpperPacket(t, rest)
		if pkt.Header.MsgID == wantMsgID {
			return pkt.Body
		}
		rest = next
	}
	t.Fatalf("packet msgID=%d not found in %d bytes", wantMsgID, len(raw))
	return nil
}
