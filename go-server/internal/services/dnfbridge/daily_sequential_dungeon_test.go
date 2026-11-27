package dnfbridge

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
	"time"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	"longheng.io/server/internal/modules/dnf/worldmap"
)

func dailySequentialTestPVFText(sequences map[uint32][]uint32) string {
	var b strings.Builder
	b.WriteString("#PVF_File\n")
	for seqID, dungeons := range sequences {
		b.WriteString(fmt.Sprintf("[sequential dungeon]\n%d\n[dungeon index check]\n", seqID))
		for _, id := range dungeons {
			b.WriteString(fmt.Sprintf("%d\n", id))
		}
		b.WriteString("[/sequential dungeon]\n")
	}
	return b.String()
}

func TestDailySequentialAntonAndLukeUnlockOrderAndSixAMReset(t *testing.T) {
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

	querySeq := func(seqID uint32) (cursor byte, mask uint32) {
		t.Helper()
		f.conn.write.Reset()
		body := make([]byte, 4)
		binary.LittleEndian.PutUint32(body, seqID)
		handled, err := currentGameplayModules.dispatchUpperRequest(f.service, f.session, gameplayRequest{
			Opcode:         869,
			Classification: 1,
			Body:           body,
		})
		if !handled || err != nil {
			t.Fatalf("seq=%d handled=%v err=%v", seqID, handled, err)
		}
		p, rest := splitGameServerUpperPacket(t, f.conn.write.Bytes())
		if len(rest) != 0 || p.Header.MsgID != 602 || p.Header.Classification != 0 || len(p.Body) != 9 {
			t.Fatalf("seq=%d unexpected packet header=%+v body=%x rest=%x", seqID, p.Header, p.Body, rest)
		}
		if gotSeq := binary.LittleEndian.Uint32(p.Body[0:4]); gotSeq != seqID {
			t.Fatalf("seq=%d got=%d", seqID, gotSeq)
		}
		return p.Body[4], binary.LittleEndian.Uint32(p.Body[5:9])
	}

	// 初始状态（0 图通关）：cursor=0, mask=0 -> 客户端显示第1图可进、第2图加锁、第3~5图隐藏。
	if cursor, mask := querySeq(28); cursor != 0 || mask != 0 {
		t.Fatalf("anton initial cursor=%d mask=%d want 0/0", cursor, mask)
	}
	if cursor, mask := querySeq(42); cursor != 0 || mask != 0 {
		t.Fatalf("luke initial cursor=%d mask=%d want 0/0", cursor, mask)
	}

	// 逐图通关安图恩每日 5 图。
	ctx := context.Background()
	wantAnton := []struct {
		cursor byte
		mask   uint32
	}{{1, 1}, {2, 3}, {3, 7}, {4, 15}, {4, 31}}
	for i, id := range antonSeq {
		if _, err := f.service.recordDailySequentialClear(ctx, f.repos, "99", id, f.service.currentDailySequentialServiceDay()); err != nil {
			t.Fatal(err)
		}
		cursor, mask := querySeq(28)
		if cursor != wantAnton[i].cursor || mask != wantAnton[i].mask {
			t.Fatalf("anton step %d cursor=%d mask=%d want %d/%d", i+1, cursor, mask, wantAnton[i].cursor, wantAnton[i].mask)
		}
	}
	// 卢克每日未受安图恩影响，仍为 0/0。
	if cursor, mask := querySeq(42); cursor != 0 || mask != 0 {
		t.Fatalf("luke polluted cursor=%d mask=%d", cursor, mask)
	}

	// 通关卢克每日第 1 图。
	if _, err := f.service.recordDailySequentialClear(ctx, f.repos, "99", 5101, f.service.currentDailySequentialServiceDay()); err != nil {
		t.Fatal(err)
	}
	if cursor, mask := querySeq(42); cursor != 1 || mask != 1 {
		t.Fatalf("luke after first clear cursor=%d mask=%d want 1/1", cursor, mask)
	}

	// 跨越北京时间 06:00 服务日边界：安图恩每日与卢克每日全部自动重置回 0/0。
	now = time.Date(2026, 11, 27, 6, 0, 0, 0, tz)
	if cursor, mask := querySeq(28); cursor != 0 || mask != 0 {
		t.Fatalf("anton after 06:00 reset cursor=%d mask=%d want 0/0", cursor, mask)
	}
	if cursor, mask := querySeq(42); cursor != 0 || mask != 0 {
		t.Fatalf("luke after 06:00 reset cursor=%d mask=%d want 0/0", cursor, mask)
	}
}

func TestDailySequentialEntryBlocksSkippingAndResetsAtSixAM(t *testing.T) {
	tz := time.FixedZone("CN", 8*3600)
	now := time.Date(2026, 11, 27, 5, 59, 0, 0, tz)
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}

	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(antonSeq, antonSeq, false))
	f.service.gameplayTimers = &processGameplayTimeQueue{now: func() time.Time { return now }}
	f.source[dailySequentialInfoPath] = dailySequentialTestPVFText(map[uint32][]uint32{28: antonSeq})

	ctx := context.Background()
	// 第 1 图永远允许进入，第 2 图未通关第 1 图时被拦截。
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1001}); err != nil {
		t.Fatalf("first dungeon blocked: %v", err)
	}
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1002}); err == nil {
		t.Fatal("second dungeon admitted before clearing first")
	}

	// 通关第 1 图后第 2 图放行，第 3 图仍被拦截。
	if _, err := f.service.recordDailySequentialClear(ctx, f.repos, "99", 1001, f.service.currentDailySequentialServiceDay()); err != nil {
		t.Fatal(err)
	}
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1002}); err != nil {
		t.Fatalf("second dungeon blocked after clearing first: %v", err)
	}
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1003}); err == nil {
		t.Fatal("third dungeon admitted before clearing second")
	}

	// 角色隔离：角色 100 未通关第 1 图仍不能进第 2 图。
	peer := &gameSession{selectedCharacterID: 100, accountID: "account-100"}
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, peer, worldmap.Dungeon{ID: 1002}); err == nil {
		t.Fatal("peer inherited another character's daily clear")
	}

	// 非顺序副本不受影响。
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 3}); err != nil {
		t.Fatalf("ordinary dungeon blocked: %v", err)
	}

	// 06:00 刷新后，第 2 图重新锁定，第 1 图照常可进。
	now = time.Date(2026, 11, 27, 6, 0, 0, 0, tz)
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1002}); err == nil {
		t.Fatal("second dungeon stayed unlocked across 06:00 reset")
	}
	if err := f.service.validateCurrentLukeNormalThroneEntry(ctx, f.session, worldmap.Dungeon{ID: 1001}); err != nil {
		t.Fatalf("first dungeon blocked after 06:00 reset: %v", err)
	}
}

func TestDailySequentialSettlementCommitsAndPushesAntonProgress(t *testing.T) {
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}
	s, rt, session, _, _ := prepareCompletedSettlementRuntime(t)
	repos, ok := s.repositoryGroup()
	if !ok {
		t.Fatal("repository missing")
	}
	character, found, err := repos.Character.Load(context.Background(), "99")
	if err != nil || !found {
		t.Fatal(err)
	}
	character.Level = 1
	character.Stats["exp"] = 0
	if err := repos.Character.Save(context.Background(), character); err != nil {
		t.Fatal(err)
	}
	if err := repos.Skill.Save(context.Background(), dnfrepo.SkillRecord{CharacterID: "99", Points: dnfrepo.SkillPointState{TotalSP: 10, RemainingSP: 10, SyncedLevel: 1}}); err != nil {
		t.Fatal(err)
	}
	rt.Character = character
	session.accountID = character.AccountID
	session.lukeSelectorProgressCharacterID = session.selectedCharacterID
	rt.Dungeon.ID = 1001
	rt.Dungeon.Path = "dungeon/anton/black_smoke.dgn"
	rt.Dungeon.Metadata.BasisLevel.Set = true
	rt.Dungeon.Metadata.BasisLevel.Value = 1
	rt.Dungeon.Metadata.ExperienceIncreasingPoint.Set = true
	rt.Dungeon.Metadata.ExperienceIncreasingPoint.Value = 1
	rt.Dungeon.Metadata.DifficultyLevels = []int64{1, 1, 1, 1, 1}
	rt.Request.Difficulty = 0
	s.dungeonRuleSource = bridgePVFSource{
		dailySequentialInfoPath: dailySequentialTestPVFText(map[uint32][]uint32{28: antonSeq}),
	}
	rt.settlementPlayResultReceived = true
	rt.settlementPlayResultBody = make([]byte, currentDungeonPlayResultBaseSize)
	rt.settlementPlayResultBody[currentDungeonPlayResultClientRankPointOffset] = 73
	resources := currentDungeonSettlementTestResources(t)
	if err := s.produceCurrentDungeonSettlementPlanWithResourcesLocked(context.Background(), session, rt, resources); err != nil {
		t.Fatal(err)
	}
	want := []byte{28, 0, 0, 0, 1, 1, 0, 0, 0}
	if !bytes.Equal(rt.settlementResultPlan.LukeSelectorProgressBody, want) {
		t.Fatalf("anton settlement op602=%x want %x", rt.settlementResultPlan.LukeSelectorProgressBody, want)
	}
}

func TestDailySequentialRollbackSwitchRestoresLegacyLukeOnlyBehavior(t *testing.T) {
	t.Setenv(dailySequentialDungeonEnv, "off")
	antonSeq := []uint32{1001, 1002, 1003, 1004, 1005}
	f := newSeekerLukeNextFixture(t, seekerLukeNextSource(antonSeq, antonSeq, false))
	f.source[dailySequentialInfoPath] = dailySequentialTestPVFText(map[uint32][]uint32{28: antonSeq})

	body := make([]byte, 4)
	binary.LittleEndian.PutUint32(body, 28)
	handled, err := currentGameplayModules.dispatchUpperRequest(f.service, f.session, gameplayRequest{
		Opcode:         869,
		Classification: 1,
		Body:           body,
	})
	if !handled || err != nil || f.conn.write.Len() != 0 {
		t.Fatalf("legacy switch did not reject non-42 sequence: handled=%v err=%v sent=%x", handled, err, f.conn.write.Bytes())
	}
}
