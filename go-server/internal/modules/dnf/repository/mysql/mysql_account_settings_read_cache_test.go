package mysql

import (
	"sync/atomic"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

func TestMySQLAccountAndSettingsReadCacheHitAndInvalidate(t *testing.T) {
	now := time.Now()
	accCache := newMySQLRecordReadCache(repository.CloneAccount)
	var accCalls atomic.Int32
	acc, found, err := accCache.getOrLoad("acc-1", now, func() (repository.AccountRecord, bool, error) {
		accCalls.Add(1)
		return repository.AccountRecord{AccountID: "acc-1", Metadata: map[string]string{"k": "v1"}}, true, nil
	})
	if err != nil || !found || acc.Metadata["k"] != "v1" {
		t.Fatalf("unexpected account load: err=%v found=%v acc=%+v", err, found, acc)
	}
	acc.Metadata["k"] = "mutated"
	cachedAcc, found, err := accCache.getOrLoad("acc-1", now.Add(100*time.Millisecond), func() (repository.AccountRecord, bool, error) {
		t.Fatalf("expected account cache hit")
		return repository.AccountRecord{}, false, nil
	})
	if err != nil || !found || cachedAcc.Metadata["k"] != "v1" || accCalls.Load() != 1 {
		t.Fatalf("account cache polluted: err=%v found=%v acc=%+v", err, found, cachedAcc)
	}
	accCache.invalidate("acc-1")
	_, found, err = accCache.getOrLoad("acc-1", now.Add(200*time.Millisecond), func() (repository.AccountRecord, bool, error) {
		accCalls.Add(1)
		return repository.AccountRecord{AccountID: "acc-1"}, true, nil
	})
	if err != nil || !found || accCalls.Load() != 2 {
		t.Fatalf("expected account reload after invalidate: calls=%d", accCalls.Load())
	}

	setCache := newMySQLRecordReadCache(repository.CloneSettings)
	var setCalls atomic.Int32
	setRec, found, err := setCache.getOrLoad("scope-1", now, func() (repository.SettingsRecord, bool, error) {
		setCalls.Add(1)
		return repository.SettingsRecord{Scope: "scope-1", Values: map[string]string{"a": "1"}}, true, nil
	})
	if err != nil || !found || setRec.Values["a"] != "1" {
		t.Fatalf("unexpected settings load: err=%v found=%v rec=%+v", err, found, setRec)
	}
	setCache.invalidate("scope-1")
	_, found, err = setCache.getOrLoad("scope-1", now.Add(100*time.Millisecond), func() (repository.SettingsRecord, bool, error) {
		setCalls.Add(1)
		return repository.SettingsRecord{Scope: "scope-1"}, true, nil
	})
	if err != nil || !found || setCalls.Load() != 2 {
		t.Fatalf("expected settings reload after invalidate: calls=%d", setCalls.Load())
	}
}
