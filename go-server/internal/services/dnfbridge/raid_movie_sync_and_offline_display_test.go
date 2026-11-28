package dnfbridge

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	"longheng.io/server/internal/modules/dnf/raid"
)

func TestCurrentRuntimeRaidMembersOfflineWireFormat(t *testing.T) {
	disconnectedAt := time.Unix(1735000123, 0)
	state := runtimeRaidState{
		RaidKey:       0x01020001,
		LeaderID:      101,
		CurrentLayout: true,
		Started:       true,
		Members: []runtimeRaidMemberState{
			{
				CharID:         101,
				Name:           "Leader",
				GroupIndex:     1,
				SlotOrder:      0,
				UserState:      0,
				Level:          90,
				DisconnectedAt: disconnectedAt,
			},
			{
				CharID:            102,
				Name:              "OfflinePeer",
				GroupIndex:        2,
				SlotOrder:         1,
				UserState:         0,
				Level:             90,
				OrdinaryPartySolo: true,
				DisconnectedAt:    disconnectedAt.Add(10 * time.Second),
			},
			{
				CharID:     103,
				Name:       "OnlinePeer",
				GroupIndex: 2,
				SlotOrder:  0,
				UserState:  1,
				Level:      90,
			},
		},
	}

	records := currentRuntimeRaidMembers(state)
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
	for _, rec := range records {
		switch rec.Core.Name {
		case "Leader":
			if rec.Core.CharKey != 0 || rec.Core.Field4 != 0 || rec.Core.Field48 != uint32(disconnectedAt.Unix()) || rec.Core.GroupIndex != 1 {
				t.Fatalf("unexpected offline leader wire record: %+v", rec.Core)
			}
		case "OfflinePeer":
			if rec.Core.CharKey != 0 || rec.Core.Field4 != 0 || rec.Core.Field48 != uint32(disconnectedAt.Add(10*time.Second).Unix()) || rec.Core.GroupIndex != 2 || rec.Core.SlotOrder != 1 {
				t.Fatalf("unexpected offline peer wire record: %+v", rec.Core)
			}
		case "OnlinePeer":
			if rec.Core.CharKey != 103 || rec.Core.Field4 != 1 || rec.Core.Field48 != 90 || rec.Core.GroupIndex != 2 {
				t.Fatalf("unexpected online peer wire record: %+v", rec.Core)
			}
		}
	}

	if _, err := raid.BuildCurrentMemberRefreshMode3Body(state.RaidKey, records); err != nil {
		t.Fatalf("BuildCurrentMemberRefreshMode3Body rejected multiple offline members: %v", err)
	}
	meta := currentRuntimeRaidMetadata(state)
	if meta.Leader.Core.Name != "Leader" || meta.Leader.Core.CharKey != 0 {
		t.Fatalf("offline leader not preserved in metadata: %+v", meta.Leader.Core)
	}
	if _, err := raid.BuildCurrentMetadataBody(meta, 0); err != nil {
		t.Fatalf("BuildCurrentMetadataBody rejected offline leader: %v", err)
	}
	if _, err := raid.BuildCurrentInfoReplyBody(meta, 2); err != nil {
		t.Fatalf("BuildCurrentInfoReplyBody rejected offline leader: %v", err)
	}
	if _, err := raid.BuildCurrentDirectoryBody([]raid.CurrentMetadata{meta}); err != nil {
		t.Fatalf("BuildCurrentDirectoryBody rejected offline leader: %v", err)
	}
}

func TestCurrentRaidMovieSyncBarrierWaitsAndReleasesOnDisconnect(t *testing.T) {
	now := time.Unix(1735000200, 0)
	s := &Service{
		options:                options{currentClientVersion: "official-24.9"},
		raids:                  map[uint32]*runtimeRaidState{},
		gameSessions:           map[uint16]*gameSession{},
		currentRaidDisconnects: map[uint16]currentRaidDisconnectSeat{},
		repositoryProvider: func() (dnfrepo.Group, bool) {
			return dnfrepo.Group{Character: movieBarrierCharacterRepo{}}, true
		},
	}
	sA := &gameSession{conn: &bufferConn{}, accountID: "accA", selectedCharacterID: 11, characterGeneration: 1, loginPacketSeen: true, accountFromToken: true}
	sB := &gameSession{conn: &bufferConn{}, accountID: "accB", selectedCharacterID: 12, characterGeneration: 1, loginPacketSeen: true, accountFromToken: true}
	sA.channel.ID, sA.channel.Type = 5, 1
	sB.channel.ID, sB.channel.Type = 5, 1
	s.gameSessions[11] = sA
	s.gameSessions[12] = sB

	idA := currentRaidProjectionIdentity{InstanceID: "inst-1", Phase: 0, CharID: 11, Generation: 1}
	idB := currentRaidProjectionIdentity{InstanceID: "inst-1", Phase: 0, CharID: 12, Generation: 1}
	sA.raidProjection.resultSent, sA.raidProjection.resultIdentity, sA.raidProjection.failureResultAt = true, idA, now
	sB.raidProjection.resultSent, sB.raidProjection.resultIdentity, sB.raidProjection.failureResultAt = true, idB, now

	state := &runtimeRaidState{
		RaidKey:       0x01050001,
		LeaderID:      11,
		CurrentLayout: true,
		Started:       true,
		Members: []runtimeRaidMemberState{
			{CharID: 11, Name: "A", GroupIndex: 1, SlotOrder: 0, UserState: 1},
			{CharID: 12, Name: "B", GroupIndex: 1, SlotOrder: 1, UserState: 1},
			{CharID: 13, Name: "OfflineC", GroupIndex: 1, SlotOrder: 2, UserState: 0, DisconnectedAt: now},
		},
		Battle: &runtimeRaidBattle{
			InstanceID:     "inst-1",
			Mode:           1,
			ChannelID:      5,
			ChannelType:    1,
			Phase:          0,
			TimedOut:       true,
			MemberAccounts: map[uint16]string{11: "accA", 12: "accB", 13: "accC"},
		},
	}
	s.raids[state.RaidKey] = state

	if err := s.handleCurrentRaidFailureMovieBarrier(sA, cloneRuntimeRaidState(*state), 659); err != nil {
		t.Fatalf("player A confirm failed: %v", err)
	}
	if state.Battle.MovieSyncReleased[0] {
		t.Fatal("barrier released before online teammate B confirmed")
	}
	if sA.raidProjection.movieCompleted {
		t.Fatal("player A failure movie should stay open while waiting for B")
	}

	// Teammate B disconnects mid-movie; barrier must release immediately for A.
	state.Members[1].UserState = 0
	state.Members[1].DisconnectedAt = now.Add(2 * time.Second)
	if err := s.tryReleaseCurrentRaidMovieBarrier(cloneRuntimeRaidState(*state)); err != nil {
		t.Fatalf("barrier release after disconnect failed: %v", err)
	}
	if !state.Battle.MovieSyncReleased[0] {
		t.Fatal("barrier did not release after remaining unconfirmed peer disconnected")
	}
	if !sA.raidProjection.movieCompleted {
		t.Fatal("player A failure movie not marked completed after barrier release")
	}
}

func TestCurrentRaidSuccessMovieBarrierHoldsRewardCardReadyAcrossTwoTeams(t *testing.T) {
	s := &Service{
		options:                options{currentClientVersion: "official-24.9"},
		raids:                  map[uint32]*runtimeRaidState{},
		gameSessions:           map[uint16]*gameSession{},
		currentRaidDisconnects: map[uint16]currentRaidDisconnectSeat{},
		repositoryProvider: func() (dnfrepo.Group, bool) {
			return dnfrepo.Group{Character: movieBarrierCharacterRepo{}}, true
		},
	}
	sTeam1 := &gameSession{conn: &bufferConn{}, accountID: "acc1", selectedCharacterID: 21, characterGeneration: 1, loginPacketSeen: true, accountFromToken: true}
	sTeam2 := &gameSession{conn: &bufferConn{}, accountID: "acc2", selectedCharacterID: 22, characterGeneration: 1, loginPacketSeen: true, accountFromToken: true}
	sTeam1.channel.ID, sTeam1.channel.Type = 5, 1
	sTeam2.channel.ID, sTeam2.channel.Type = 5, 1
	s.gameSessions[21] = sTeam1
	s.gameSessions[22] = sTeam2

	id1 := currentRaidProjectionIdentity{InstanceID: "inst-2", Phase: 0, CharID: 21, Generation: 1}
	sTeam1.raidProjection.resultSent, sTeam1.raidProjection.resultIdentity = true, id1
	// Team 2 is online in the raid but has not finished loading resultSent yet.

	state := &runtimeRaidState{
		RaidKey:       0x01050002,
		LeaderID:      21,
		CurrentLayout: true,
		Started:       true,
		Members: []runtimeRaidMemberState{
			{CharID: 21, Name: "Team1", GroupIndex: 1, SlotOrder: 0, UserState: 1},
			{CharID: 22, Name: "Team2", GroupIndex: 2, SlotOrder: 0, UserState: 1},
		},
		Battle: &runtimeRaidBattle{
			InstanceID:     "inst-2",
			Mode:           1,
			ChannelID:      5,
			ChannelType:    1,
			Phase:          0,
			PhaseCompleted: true,
			MemberAccounts: map[uint16]string{21: "acc1", 22: "acc2"},
		},
	}
	s.raids[state.RaidKey] = state

	if err := s.handleCurrentRaidSuccessMovieBarrier(sTeam1, cloneRuntimeRaidState(*state), 659); err != nil {
		t.Fatalf("team1 confirm failed: %v", err)
	}
	if state.Battle.MovieSyncReleased[0] {
		t.Fatal("barrier released while online team2 has not confirmed")
	}
	if state.Battle.RewardCardReady[0][21] {
		t.Fatal("RewardCardReady set before barrier release (would cause 1s projection tick to open cards early)")
	}
}

// movieBarrierCharacterRepo 只实现 currentRaidAuthenticated 需要的 Character.Load，
// 其余方法由嵌入的 nil 接口承担（测试不会走到）。
type movieBarrierCharacterRepo struct {
	dnfrepo.CharacterRepository
}

var movieBarrierTestAccounts = map[uint16]string{
	11: "accA", 12: "accB", 13: "accC", 21: "acc1", 22: "acc2",
}

func (movieBarrierCharacterRepo) Load(_ context.Context, id string) (dnfrepo.CharacterRecord, bool, error) {
	n, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		return dnfrepo.CharacterRecord{}, false, nil
	}
	account, ok := movieBarrierTestAccounts[uint16(n)]
	if !ok {
		return dnfrepo.CharacterRecord{}, false, nil
	}
	return dnfrepo.CharacterRecord{CharacterID: id, AccountID: account}, true, nil
}
