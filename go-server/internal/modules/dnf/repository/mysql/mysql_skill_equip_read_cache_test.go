package mysql

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

func TestMySQLSkillEquipReadCacheSingleflightAndInvalidate(t *testing.T) {
	skillCache := newMySQLRecordReadCache(repository.CloneSkill)
	now := time.Now()
	var calls atomic.Int32
	gate := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec, found, err := skillCache.getOrLoad("1001", now, func() (repository.SkillRecord, bool, error) {
				calls.Add(1)
				<-gate
				return repository.SkillRecord{CharacterID: "1001", Levels: map[string]int64{"14": 10}}, true, nil
			})
			if err != nil || !found || rec.Levels["14"] != 10 {
				t.Errorf("unexpected skill load: err=%v found=%v rec=%+v", err, found, rec)
			}
			rec.Levels["14"] = 0
		}()
	}
	close(gate)
	wg.Wait()

	if calls.Load() != 1 {
		t.Fatalf("expected 1 SQL load, got %d", calls.Load())
	}
	cached, found, err := skillCache.getOrLoad("1001", now.Add(200*time.Millisecond), func() (repository.SkillRecord, bool, error) {
		t.Fatalf("expected cache hit")
		return repository.SkillRecord{}, false, nil
	})
	if err != nil || !found || cached.Levels["14"] != 10 {
		t.Fatalf("unexpected cached skill: err=%v found=%v rec=%+v", err, found, cached)
	}
	skillCache.invalidate("1001")
	reloaded, found, err := skillCache.getOrLoad("1001", now.Add(300*time.Millisecond), func() (repository.SkillRecord, bool, error) {
		calls.Add(1)
		return repository.SkillRecord{CharacterID: "1001", Levels: map[string]int64{"14": 11}}, true, nil
	})
	if err != nil || !found || reloaded.Levels["14"] != 11 || calls.Load() != 2 {
		t.Fatalf("expected reload after invalidate: err=%v found=%v rec=%+v calls=%d", err, found, reloaded, calls.Load())
	}
}
