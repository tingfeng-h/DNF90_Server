package dnfbridge

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	dnfpremium "longheng.io/server/internal/modules/dnf/premium"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const (
	// dailyMonsterDropOnceEnv 控制安图恩每日与卢克每日每角色每日 06:00 每怪首杀掉落去重。
	// 默认开启；设为 off/false/0/no/legacy 时逐字节回退至旧版每次击杀均掉落行为。
	dailyMonsterDropOnceEnv = "DNFBRIDGE_DAILY_MONSTER_DROP_ONCE"

	dailyMonsterDropScopePrefix = "daily_monster_drop:"
	dailyMonsterDropDayKey      = "day"
	dailyMonsterDropKillsKey    = "kills"
)

type dailyMonsterDropCharBucket struct {
	mu     sync.Mutex
	day    int64
	loaded bool
	kills  map[string]struct{}
}

type dailyMonsterDropLedger struct {
	mu      sync.Mutex
	buckets map[uint16]*dailyMonsterDropCharBucket
}

// 按 *Service 指针隔离进程内热缓存，既不改动 Service 结构体字段避免并发合入冲突，
// 也能确保并行单元测试中的不同 Service 实例完全互不干扰。
var dailyMonsterDropLedgers sync.Map // key: *Service -> *dailyMonsterDropLedger

func dailyMonsterDropOnceEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(envString(dailyMonsterDropOnceEnv, ""))) {
	case "off", "false", "0", "no", "legacy":
		return false
	default:
		return true
	}
}

func dailyMonsterDropScope(characterID uint16) string {
	return dailyMonsterDropScopePrefix + strconv.FormatUint(uint64(characterID), 10)
}

func dailyMonsterDropKillKey(dungeonID, monsterID int64) string {
	return strconv.FormatInt(dungeonID, 10) + ":" + strconv.FormatInt(monsterID, 10)
}

func (s *Service) dailyMonsterDropCharBucket(characterID uint16) *dailyMonsterDropCharBucket {
	var ledger *dailyMonsterDropLedger
	if existing, ok := dailyMonsterDropLedgers.Load(s); ok {
		ledger = existing.(*dailyMonsterDropLedger)
	} else {
		created := &dailyMonsterDropLedger{buckets: make(map[uint16]*dailyMonsterDropCharBucket)}
		actual, _ := dailyMonsterDropLedgers.LoadOrStore(s, created)
		ledger = actual.(*dailyMonsterDropLedger)
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	bucket, ok := ledger.buckets[characterID]
	if !ok {
		bucket = &dailyMonsterDropCharBucket{kills: make(map[string]struct{})}
		ledger.buckets[characterID] = bucket
	}
	return bucket
}

// isDailySequentialMonsterDropDungeon 判定当前副本是否属于安图恩每日或卢克每日：
// 1) 严格排除一切团本（runtime.raidEntry != nil）；
// 2) 优先匹配 PVF etc/sequential_dungeon_info.etc 中的顺序副本；
// 3) 兜底匹配普通卢克（5101..5106 / dungeon/luke_normal/）及路径含 anton 的非团本副本。
func (s *Service) isDailySequentialMonsterDropDungeon(runtime *runtimeDungeonState) bool {
	if runtime == nil || runtime.raidEntry != nil || runtime.Dungeon.ID <= 0 {
		return false
	}
	if currentLukeNormalDungeon(runtime) {
		return true
	}
	if runtime.Dungeon.ID <= int64(^uint32(0)) {
		if catalog, err := s.currentDailySequentialCatalog(); err == nil {
			if len(catalog.ByDungeon[uint32(runtime.Dungeon.ID)]) > 0 {
				return true
			}
		}
	}
	path := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(runtime.Dungeon.Path, "\\", "/")))
	if strings.Contains(path, "anton") && !strings.Contains(path, "raid") {
		return true
	}
	return false
}

func parseDailyMonsterDropKills(values map[string]string, serviceDay int64) map[string]struct{} {
	kills := make(map[string]struct{})
	if len(values) == 0 {
		return kills
	}
	storedDay, err := strconv.ParseInt(strings.TrimSpace(values[dailyMonsterDropDayKey]), 10, 64)
	if err != nil || storedDay != serviceDay {
		return kills
	}
	raw := strings.TrimSpace(values[dailyMonsterDropKillsKey])
	if raw == "" {
		return kills
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kills[part] = struct{}{}
	}
	return kills
}

func formatDailyMonsterDropKills(serviceDay int64, kills map[string]struct{}) map[string]string {
	keys := make([]string, 0, len(kills))
	for k := range kills {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return map[string]string{
		dailyMonsterDropDayKey:   strconv.FormatInt(serviceDay, 10),
		dailyMonsterDropKillsKey: strings.Join(keys, ","),
	}
}

// claimDailySequentialMonsterDrop 判定某角色在当前服务日（北京时间 06:00 刷新）
// 是否首次击杀安图恩每日/卢克每日的指定怪物：
//   - 返回 true：允许本次产生掉落奖励（并登记今日已掉落，持久化到 Settings 仓储）；
//   - 返回 false：该角色今日已从该怪物获取过掉落奖励，本次应跳过掉落。
//
// 仅作用于普通怪物地面掉落规划路径，不影响每日任务物品掉落（quest drop）。
func (s *Service) claimDailySequentialMonsterDrop(characterID uint16, runtime *runtimeDungeonState, monsterID int64) bool {
	if s == nil || !dailyMonsterDropOnceEnabled() || characterID == 0 || monsterID <= 0 {
		return true
	}
	if !s.isDailySequentialMonsterDropDungeon(runtime) {
		return true
	}
	serviceDay := dnfpremium.ServiceDayID(s.gameplayNow())
	killKey := dailyMonsterDropKillKey(runtime.Dungeon.ID, monsterID)

	bucket := s.dailyMonsterDropCharBucket(characterID)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if bucket.day != serviceDay {
		bucket.day = serviceDay
		bucket.loaded = false
		bucket.kills = make(map[string]struct{})
	}
	repos, hasRepos := s.repositoryGroup()
	if !bucket.loaded {
		if hasRepos && repos.Settings != nil {
			ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
			record, found, err := repos.Settings.Load(ctx, dailyMonsterDropScope(characterID))
			cancel()
			if err == nil && found {
				bucket.kills = parseDailyMonsterDropKills(record.Values, serviceDay)
			}
		}
		bucket.loaded = true
	}
	if _, alreadyDropped := bucket.kills[killKey]; alreadyDropped {
		return false
	}
	bucket.kills[killKey] = struct{}{}
	if hasRepos && repos.Settings != nil {
		ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
		err := repos.Settings.Save(ctx, dnfrepo.SettingsRecord{
			Scope:  dailyMonsterDropScope(characterID),
			Values: formatDailyMonsterDropKills(serviceDay, bucket.kills),
		})
		cancel()
		if err != nil {
			s.logWarn("dnfbridge daily sequential monster drop ledger save failed",
				"char_id", characterID,
				"dungeon_id", runtime.Dungeon.ID,
				"monster_id", monsterID,
				"service_day", serviceDay,
				"error", err)
		}
	}
	return true
}
