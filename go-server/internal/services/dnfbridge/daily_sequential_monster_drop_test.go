package dnfbridge

import (
	"testing"
	"time"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	"longheng.io/server/internal/modules/dnf/worldmap"
)

func TestDailySequentialMonsterDropOncePerDayAndSixAMReset(t *testing.T) {
	tz := time.FixedZone("CN", 8*3600)
	now := time.Date(2026, 11, 27, 5, 59, 59, 0, tz)
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}
	lukeSeq := []uint32{5101, 5102, 5103, 5104, 5105, 5106}

	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(lukeSeq, lukeSeq, true))
	f.service.gameplayTimers = &processGameplayTimeQueue{now: func() time.Time { return now }}
	f.source[dailySequentialInfoPath] = dailySequentialTestPVFText(map[uint32][]uint32{
		28: antonSeq,
		42: lukeSeq,
	})

	antonSmoke := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 1001, Path: "dungeon/anton/black_smoke.dgn"}}
	lukeBirth := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 5101, Path: "dungeon/luke_normal/birth.dgn"}}

	// 安图恩每日吞噬魔：第 1 次击杀允许掉落，第 2 次拦截；同图另一怪物（内尔贝）互不干扰。
	if !f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("first kill of devourer blocked")
	}
	if f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("second kill of devourer in same service day was not suppressed")
	}
	if !f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61002) {
		t.Fatal("different monster in same dungeon blocked")
	}

	// 卢克每日同理：第 1 次允许，第 2 次拦截。
	if !f.service.claimDailySequentialMonsterDrop(99, lukeBirth, 56001) {
		t.Fatal("first kill in luke daily blocked")
	}
	if f.service.claimDailySequentialMonsterDrop(99, lukeBirth, 56001) {
		t.Fatal("second kill in luke daily was not suppressed")
	}

	// 跨北京时间 06:00 刷新：重新放行首杀，随后再次拦截重复击杀。
	now = time.Date(2026, 11, 27, 6, 0, 0, 0, tz)
	if !f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("devourer drop did not reset at 06:00")
	}
	if f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("second kill after 06:00 reset was not suppressed")
	}
	if !f.service.claimDailySequentialMonsterDrop(99, lukeBirth, 56001) {
		t.Fatal("luke daily drop did not reset at 06:00")
	}
}

func TestDailySequentialMonsterDropPersistsAcrossServiceRestart(t *testing.T) {
	tz := time.FixedZone("CN", 8*3600)
	now := time.Date(2026, 11, 27, 10, 0, 0, 0, tz)
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}
	source := bridgePVFSource{
		dailySequentialInfoPath: dailySequentialTestPVFText(map[uint32][]uint32{28: antonSeq}),
	}
	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(antonSeq, antonSeq, false))
	f.service.gameplayTimers = &processGameplayTimeQueue{now: func() time.Time { return now }}
	f.service.dungeonRuleSource = source

	antonSmoke := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 1001, Path: "dungeon/anton/black_smoke.dgn"}}
	if !f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("first kill blocked")
	}

	// 模拟服务端重启：新 Service 实例无内存缓存，但共享同一个持久化 Settings 仓储。
	restarted := &Service{
		dungeonRuleSource:  source,
		gameplayTimers:     &processGameplayTimeQueue{now: func() time.Time { return now }},
		repositoryProvider: func() (dnfrepo.Group, bool) { return f.repos, true },
	}
	if restarted.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("daily kill record lost after service restart")
	}
}

func TestDailySequentialMonsterDropPartyIsolationAndScopeBoundaries(t *testing.T) {
	tz := time.FixedZone("CN", 8*3600)
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, tz)
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}
	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(antonSeq, antonSeq, false))
	f.service.gameplayTimers = &processGameplayTimeQueue{now: func() time.Time { return now }}
	f.source[dailySequentialInfoPath] = dailySequentialTestPVFText(map[uint32][]uint32{28: antonSeq})

	antonSmoke := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 1001, Path: "dungeon/anton/black_smoke.dgn"}}
	// 组队角色独立：角色 99 已拿过掉落，队友 100 首次击杀仍正常掉落。
	if !f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("char 99 first kill blocked")
	}
	if f.service.claimDailySequentialMonsterDrop(99, antonSmoke, 61001) {
		t.Fatal("char 99 duplicate kill allowed")
	}
	if !f.service.claimDailySequentialMonsterDrop(100, antonSmoke, 61001) {
		t.Fatal("teammate 100 first kill wrongly blocked by char 99")
	}

	// 团本（raidEntry != nil）绝不受拦截。
	raidRuntime := &runtimeDungeonState{
		Dungeon:   worldmap.Dungeon{ID: 210, Path: "dungeon/anton/raid_smoke.dgn"},
		raidEntry: &currentRaidDungeonEntry{ChannelType: 23},
	}
	for i := 0; i < 3; i++ {
		if !f.service.claimDailySequentialMonsterDrop(99, raidRuntime, 61001) {
			t.Fatalf("raid kill %d wrongly suppressed", i)
		}
	}

	// 普通白图绝不受拦截。
	ordinaryRuntime := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 3, Path: "dungeon/lorien/lorien.dgn"}}
	for i := 0; i < 3; i++ {
		if !f.service.claimDailySequentialMonsterDrop(99, ordinaryRuntime, 10001) {
			t.Fatalf("ordinary dungeon kill %d wrongly suppressed", i)
		}
	}
}

func TestDailySequentialMonsterDropRollbackSwitchAllowsRepeat(t *testing.T) {
	t.Setenv(dailyMonsterDropOnceEnv, "off")
	lukeSeq := []uint32{5101, 5102, 5103, 5104, 5105, 5106}
	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(lukeSeq, lukeSeq, true))
	lukeBirth := &runtimeDungeonState{Dungeon: worldmap.Dungeon{ID: 5101, Path: "dungeon/luke_normal/birth.dgn"}}
	if !f.service.claimDailySequentialMonsterDrop(99, lukeBirth, 56001) || !f.service.claimDailySequentialMonsterDrop(99, lukeBirth, 56001) {
		t.Fatal("rollback switch did not allow repeat monster drops")
	}
}
