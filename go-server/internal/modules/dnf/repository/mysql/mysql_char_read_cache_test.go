package mysql

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

func TestMySQLCharReadCacheSingleflightAndWriteInvalidation(t *testing.T) {
	cache := newMySQLCharReadCache()
	now := time.Now()
	var calls atomic.Int32
	gate := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec, found, err := cache.getOrLoad("1001", now, func() (repository.CharacterRecord, bool, error) {
				calls.Add(1)
				<-gate
				return repository.CharacterRecord{CharacterID: "1001", Name: "Slayer", Level: 90, Stats: map[string]int64{"exp": 123}}, true, nil
			})
			if err != nil || !found || rec.Stats["exp"] != 123 {
				t.Errorf("unexpected load: err=%v found=%v rec=%+v", err, found, rec)
			}
			rec.Stats["exp"] = 999
		}()
	}
	close(gate)
	wg.Wait()

	if calls.Load() != 1 {
		t.Fatalf("expected 1 SQL load under singleflight, got %d", calls.Load())
	}

	// TTL 内命中且隔离调用方修改
	cached, found, err := cache.getOrLoad("1001", now.Add(200*time.Millisecond), func() (repository.CharacterRecord, bool, error) {
		t.Fatalf("expected cache hit within TTL")
		return repository.CharacterRecord{}, false, nil
	})
	if err != nil || !found || cached.Stats["exp"] != 123 {
		t.Fatalf("unexpected cached record: err=%v found=%v rec=%+v", err, found, cached)
	}

	// 写入失效后立刻重新查库
	cache.invalidate("1001")
	reloaded, found, err := cache.getOrLoad("1001", now.Add(300*time.Millisecond), func() (repository.CharacterRecord, bool, error) {
		calls.Add(1)
		return repository.CharacterRecord{CharacterID: "1001", Name: "Slayer", Level: 91}, true, nil
	})
	if err != nil || !found || reloaded.Level != 91 || calls.Load() != 2 {
		t.Fatalf("expected reload after invalidate: err=%v found=%v level=%d calls=%d", err, found, reloaded.Level, calls.Load())
	}
}

func TestMySQLRouterBypassesCharReadCacheInTransactions(t *testing.T) {
	base := mysqlRouter{charReadCache: newMySQLCharReadCache()}
	if !base.canUseCharReadCache() {
		t.Fatalf("expected root non-tx router to use charReadCache")
	}
	locked := base
	locked.lockReads = true
	if locked.canUseCharReadCache() {
		t.Fatalf("lockReads router must bypass charReadCache")
	}
	consistent := base
	consistent.consistentRead = true
	if consistent.canUseCharReadCache() {
		t.Fatalf("consistentRead router must bypass charReadCache")
	}
	inTx := base
	inTx.db = &sqlTransactionDB{}
	if inTx.canUseCharReadCache() {
		t.Fatalf("sqlTransactionDB router must bypass charReadCache")
	}
}
