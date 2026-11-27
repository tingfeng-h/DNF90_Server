package dnfbridge

import (
	"context"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	dnfpremium "longheng.io/server/internal/modules/dnf/premium"
	dnfpvf "longheng.io/server/internal/modules/dnf/pvf"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	"longheng.io/server/internal/modules/dnf/worldmap"
)

const (
	// dailySequentialDungeonEnv 控制安图恩每日与卢克每日的顺序副本解锁与每日 06:00 重置逻辑。
	// 默认开启；设为 off/false/0/no/legacy 时逐字节回退至旧版仅卢克终身记录行为。
	dailySequentialDungeonEnv = "DNFBRIDGE_DAILY_SEQUENTIAL_DUNGEON"

	dailySequentialInfoPath          = "etc/sequential_dungeon_info.etc"
	dailySequentialClearScopePrefix  = "daily_sequential_clear:"
	dailySequentialClearDayKey       = "day"
	dailySequentialClearDungeonsKey  = "cleared"
	dailySequentialNativeMaxDungeons = 8
)

type dailySequentialDefinition struct {
	SequenceID uint32
	Dungeons   []uint32
	AutoEnter  bool
}

type dailySequentialCatalog struct {
	BySequence map[uint32]dailySequentialDefinition
	ByDungeon  map[uint32][]dailySequentialDefinition
}

var dailySequentialClearLocks sync.Map // key: string(scope) -> *sync.Mutex

func dailySequentialDungeonEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(envString(dailySequentialDungeonEnv, ""))) {
	case "off", "false", "0", "no", "legacy":
		return false
	default:
		return true
	}
}

func dailySequentialClearScope(characterID string) string {
	return dailySequentialClearScopePrefix + strings.TrimSpace(characterID)
}

func dailySequentialClearMutex(scope string) *sync.Mutex {
	if existing, ok := dailySequentialClearLocks.Load(scope); ok {
		return existing.(*sync.Mutex)
	}
	created := &sync.Mutex{}
	actual, _ := dailySequentialClearLocks.LoadOrStore(scope, created)
	return actual.(*sync.Mutex)
}

func loadDailySequentialCatalog(source initialEquipmentTextSource) (*dailySequentialCatalog, error) {
	if source == nil {
		return nil, fmt.Errorf("sequential dungeon source unavailable")
	}
	text, err := source.ReadText(dailySequentialInfoPath)
	if err != nil {
		return nil, err
	}
	doc, err := dnfpvf.Parse(dailySequentialInfoPath, text)
	if err != nil {
		return nil, err
	}
	catalog := &dailySequentialCatalog{
		BySequence: make(map[uint32]dailySequentialDefinition),
		ByDungeon:  make(map[uint32][]dailySequentialDefinition),
	}
	var current *dailySequentialDefinition
	var currentDup bool
	flush := func() error {
		if current == nil {
			return nil
		}
		def := *current
		current = nil
		if currentDup {
			currentDup = false
			return fmt.Errorf("sequential dungeon %d duplicated", def.SequenceID)
		}
		if len(def.Dungeons) == 0 {
			return nil
		}
		if _, exists := catalog.BySequence[def.SequenceID]; exists {
			return fmt.Errorf("sequential dungeon %d duplicated", def.SequenceID)
		}
		catalog.BySequence[def.SequenceID] = def
		for _, dungeonID := range def.Dungeons {
			catalog.ByDungeon[dungeonID] = append(catalog.ByDungeon[dungeonID], def)
		}
		return nil
	}

	for _, section := range doc.Sections {
		name := strings.ToLower(strings.TrimSpace(section.Name))
		values := doc.Tokens[section.Start:section.End]
		switch name {
		case "sequential dungeon":
			if err := flush(); err != nil {
				return nil, err
			}
			if len(values) != 1 || values[0].Kind != dnfpvf.TokenInt || values[0].Int <= 0 || values[0].Int > int64(^uint32(0)) {
				continue
			}
			seqID := uint32(values[0].Int)
			if _, exists := catalog.BySequence[seqID]; exists {
				current = &dailySequentialDefinition{SequenceID: seqID}
				currentDup = true
				continue
			}
			current = &dailySequentialDefinition{SequenceID: seqID}
		case "/sequential dungeon":
			if err := flush(); err != nil {
				return nil, err
			}
		default:
			if current == nil {
				continue
			}
			switch name {
			case "dungeon index check":
				if len(current.Dungeons) != 0 {
					return nil, fmt.Errorf("sequential dungeon %d list duplicated", current.SequenceID)
				}
				seen := make(map[uint32]bool, len(values))
				for _, token := range values {
					if token.Kind != dnfpvf.TokenInt || token.Int <= 0 || token.Int > int64(^uint32(0)) {
						return nil, fmt.Errorf("sequential dungeon %d contains invalid dungeon id", current.SequenceID)
					}
					id := uint32(token.Int)
					if seen[id] {
						return nil, fmt.Errorf("sequential dungeon %d contains duplicate dungeon %d", current.SequenceID, id)
					}
					seen[id] = true
					current.Dungeons = append(current.Dungeons, id)
				}
			case "next dungeon auto enter":
				if len(values) != 0 {
					return nil, fmt.Errorf("sequential dungeon %d auto-enter marker invalid", current.SequenceID)
				}
				current.AutoEnter = true
			}
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	for dungeonID := range catalog.ByDungeon {
		sort.Slice(catalog.ByDungeon[dungeonID], func(i, j int) bool {
			return catalog.ByDungeon[dungeonID][i].SequenceID < catalog.ByDungeon[dungeonID][j].SequenceID
		})
	}
	return catalog, nil
}

func (s *Service) currentDailySequentialCatalog() (*dailySequentialCatalog, error) {
	source, err := s.currentDungeonRuleSource()
	if err != nil {
		return nil, err
	}
	return loadDailySequentialCatalog(source)
}

func (s *Service) currentDailySequentialServiceDay() int64 {
	return dnfpremium.ServiceDayID(s.gameplayNow())
}

func parseDailySequentialClearedValues(values map[string]string, serviceDay int64) map[uint32]bool {
	cleared := make(map[uint32]bool)
	if len(values) == 0 {
		return cleared
	}
	storedDay, err := strconv.ParseInt(strings.TrimSpace(values[dailySequentialClearDayKey]), 10, 64)
	if err != nil || storedDay != serviceDay {
		return cleared
	}
	raw := strings.TrimSpace(values[dailySequentialClearDungeonsKey])
	if raw == "" {
		return cleared
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, parseErr := strconv.ParseUint(part, 10, 32)
		if parseErr != nil || id == 0 {
			continue
		}
		cleared[uint32(id)] = true
	}
	return cleared
}

func formatDailySequentialClearedValues(serviceDay int64, cleared map[uint32]bool) map[string]string {
	ids := make([]uint32, 0, len(cleared))
	for id, ok := range cleared {
		if ok && id != 0 {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatUint(uint64(id), 10)
	}
	return map[string]string{
		dailySequentialClearDayKey:      strconv.FormatInt(serviceDay, 10),
		dailySequentialClearDungeonsKey: strings.Join(parts, ","),
	}
}

func (s *Service) loadDailySequentialClearedSet(ctx context.Context, repos dnfrepo.Group, characterID string, serviceDay int64) (map[uint32]bool, error) {
	if repos.Settings == nil {
		return nil, fmt.Errorf("daily sequential settings repository unavailable")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	record, found, err := repos.Settings.Load(ctx, dailySequentialClearScope(characterID))
	if err != nil {
		return nil, err
	}
	if !found {
		return make(map[uint32]bool), nil
	}
	return parseDailySequentialClearedValues(record.Values, serviceDay), nil
}

func (s *Service) recordDailySequentialClear(ctx context.Context, repos dnfrepo.Group, characterID string, dungeonID uint32, serviceDay int64) (map[uint32]bool, error) {
	if repos.Settings == nil {
		return nil, fmt.Errorf("daily sequential settings repository unavailable")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	scope := dailySequentialClearScope(characterID)
	mu := dailySequentialClearMutex(scope)
	mu.Lock()
	defer mu.Unlock()

	record, found, err := repos.Settings.Load(ctx, scope)
	if err != nil {
		return nil, err
	}
	var cleared map[uint32]bool
	if found {
		cleared = parseDailySequentialClearedValues(record.Values, serviceDay)
	} else {
		cleared = make(map[uint32]bool)
	}
	cleared[dungeonID] = true
	nextRecord := dnfrepo.SettingsRecord{
		Scope:  scope,
		Values: formatDailySequentialClearedValues(serviceDay, cleared),
	}
	if err := repos.Settings.Save(ctx, nextRecord); err != nil {
		return nil, err
	}
	return cleared, nil
}

// buildDailySequentialSelectorProgress 构建客户端原生 class0/op602 的 9 字节正文：
// u32 sequence_id + u8 enterable_cursor + u32 individual_clear_mask。
// 客户端 0x018C0310 据此自动计算：
//   - visible   = ids[:min(cursor+2, len)]（当前可进 + 下一张带锁图可见，更后续图自动隐藏）
//   - enterable = ids[:cursor+1]          （第 cursor+1 张图自动显示锁图标且不可进入）
//   - completed = clear_mask 置位集合
func buildDailySequentialSelectorProgress(sequence dailySequentialDefinition, cleared map[uint32]bool) []byte {
	var mask uint32
	var cursor byte
	prefix := true
	for i, id := range sequence.Dungeons {
		if cleared[id] {
			mask |= uint32(1) << uint(i)
			if prefix {
				cursor = byte(i + 1)
			}
		} else {
			prefix = false
		}
	}
	if len(sequence.Dungeons) > 0 && int(cursor) >= len(sequence.Dungeons) {
		cursor = byte(len(sequence.Dungeons) - 1)
	}
	body := make([]byte, 9)
	binary.LittleEndian.PutUint32(body, sequence.SequenceID)
	body[4] = cursor
	binary.LittleEndian.PutUint32(body[5:], mask)
	return body
}

func (s *Service) handleDailySequentialSelectorProgress(session *gameSession, body []byte) error {
	if len(body) != 4 {
		s.logGameEvent(session, "game-luke-selector-progress-rejected", "reason", "unsupported_sequence_or_shape", "body_len", len(body))
		return nil
	}
	seqID := binary.LittleEndian.Uint32(body)
	ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
	defer cancel()
	progress, err := s.currentDailySequentialSelectorProgress(ctx, session, seqID)
	if err != nil {
		s.logGameEvent(session, "game-luke-selector-progress-rejected", "reason", err, "sequence_id", seqID)
		return nil
	}
	if err := s.sendGameUpperRawClass(session, currentLukeSelectorProgressMsgID, progress, 0); err != nil {
		return err
	}
	session.lukeSelectorProgressCharacterID = session.selectedCharacterID
	s.logGameEvent(session, "game-luke-selector-progress-sent",
		"character_id", session.selectedCharacterID,
		"sequence_id", seqID,
		"cursor", progress[4],
		"clear_mask", binary.LittleEndian.Uint32(progress[5:]),
		"service_day", s.currentDailySequentialServiceDay(),
		"source", "character_daily_sequential_clears")
	return nil
}

func (s *Service) currentDailySequentialSelectorProgress(ctx context.Context, session *gameSession, sequenceID uint32) ([]byte, error) {
	if session == nil || session.selectedCharacterID == 0 {
		return nil, fmt.Errorf("sequential progress character required")
	}
	repos, ok := s.repositoryGroupFor(ctx)
	if !ok || repos.Character == nil || repos.Settings == nil {
		return nil, fmt.Errorf("sequential progress repository unavailable")
	}
	id := strconv.Itoa(int(session.selectedCharacterID))
	character, found, err := repos.Character.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	account := strings.TrimSpace(s.accountIDForSession(session))
	if !found || character.CharacterID != id || account == "" || strings.TrimSpace(character.AccountID) != account {
		return nil, fmt.Errorf("sequential progress character owner mismatch")
	}
	catalog, err := s.currentDailySequentialCatalog()
	if err != nil {
		return nil, err
	}
	sequence, exists := catalog.BySequence[sequenceID]
	if !exists {
		return nil, fmt.Errorf("sequential dungeon %d unavailable", sequenceID)
	}
	if len(sequence.Dungeons) == 0 || len(sequence.Dungeons) > dailySequentialNativeMaxDungeons {
		return nil, fmt.Errorf("sequential dungeon %d exceeds native mask", sequenceID)
	}
	cleared, err := s.loadDailySequentialClearedSet(ctx, repos, id, s.currentDailySequentialServiceDay())
	if err != nil {
		return nil, err
	}
	return buildDailySequentialSelectorProgress(sequence, cleared), nil
}

// validateDailySequentialDungeonEntry 校验安图恩每日/卢克每日的当日顺序入场资格：
// 序列第 0 张图恒放行；第 i>0 张图必须先在当前服务日（06:00 起）通关全部前置图 0..i-1。
func (s *Service) validateDailySequentialDungeonEntry(ctx context.Context, session *gameSession, dungeon worldmap.Dungeon) error {
	if dungeon.ID <= 0 || dungeon.ID > int64(^uint32(0)) {
		return nil
	}
	catalog, err := s.currentDailySequentialCatalog()
	if err != nil {
		// 非顺序副本或不具备 PVF 规则源时，仅对卢克王座 5106 保持硬校验，避免影响普通白图测试。
		if dungeon.ID == 5106 && strings.EqualFold(strings.ReplaceAll(dungeon.Path, "\\", "/"), "dungeon/luke_normal/solium.dgn") {
			return err
		}
		return nil
	}
	sequences := catalog.ByDungeon[uint32(dungeon.ID)]
	if len(sequences) == 0 {
		return nil
	}
	// 检查是否存在需要前置图的序列位置。
	var prerequisites []uint32
	for _, seq := range sequences {
		for idx, id := range seq.Dungeons {
			if id == uint32(dungeon.ID) && idx > 0 {
				prerequisites = append(prerequisites, seq.Dungeons[:idx]...)
				break
			}
		}
	}
	if len(prerequisites) == 0 {
		return nil
	}
	if session == nil || session.selectedCharacterID == 0 {
		return errDungeonCharacterRequired
	}
	repos, ok := s.repositoryGroupFor(ctx)
	if !ok || repos.Settings == nil {
		return fmt.Errorf("daily sequential clear repository unavailable")
	}
	cleared, err := s.loadDailySequentialClearedSet(ctx, repos, strconv.Itoa(int(session.selectedCharacterID)), s.currentDailySequentialServiceDay())
	if err != nil {
		return err
	}
	for _, reqID := range prerequisites {
		if !cleared[reqID] {
			return fmt.Errorf("daily sequential dungeon %d prerequisite %d incomplete", dungeon.ID, reqID)
		}
	}
	return nil
}

// commitAndBuildDailySequentialProgress 在非团本顺序副本结算后记录今日通关，
// 若当前会话已查询过选图进度则返回最新的 9 字节 op602 包供结算原子推送。
func (s *Service) commitAndBuildDailySequentialProgress(ctx context.Context, session *gameSession, runtime *runtimeDungeonState) ([]byte, error) {
	if session == nil || session.selectedCharacterID == 0 || runtime == nil || runtime.raidEntry != nil {
		return nil, nil
	}
	if runtime.Dungeon.ID <= 0 || runtime.Dungeon.ID > int64(^uint32(0)) {
		return nil, nil
	}
	catalog, err := s.currentDailySequentialCatalog()
	if err != nil {
		return nil, nil
	}
	dungeonID := uint32(runtime.Dungeon.ID)
	sequences := catalog.ByDungeon[dungeonID]
	if len(sequences) == 0 {
		return nil, nil
	}
	repos, ok := s.repositoryGroupFor(ctx)
	if !ok || repos.Settings == nil {
		return nil, fmt.Errorf("daily sequential settings repository unavailable")
	}
	charID := strconv.Itoa(int(session.selectedCharacterID))
	cleared, err := s.recordDailySequentialClear(ctx, repos, charID, dungeonID, s.currentDailySequentialServiceDay())
	if err != nil {
		return nil, err
	}
	if session.lukeSelectorProgressCharacterID != session.selectedCharacterID {
		return nil, nil
	}
	seq := sequences[0]
	if len(seq.Dungeons) == 0 || len(seq.Dungeons) > dailySequentialNativeMaxDungeons {
		return nil, nil
	}
	return buildDailySequentialSelectorProgress(seq, cleared), nil
}
