package mysql

import (
	"sync/atomic"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

func TestMySQLInventoryAndQuestReadCacheHitAndInvalidate(t *testing.T) {
	now := time.Now()
	invCache := newMySQLRecordReadCache(repository.CloneInventory)
	var invCalls atomic.Int32
	rec, found, err := invCache.getOrLoad("1001", now, func() (repository.InventoryRecord, bool, error) {
		invCalls.Add(1)
		return repository.InventoryRecord{
			CharacterID: "1001",
			Slots:       map[string]repository.ItemStack{"0:1": {ItemID: 3037, Count: 10}},
		}, true, nil
	})
	if err != nil || !found || rec.Slots["0:1"].ItemID != 3037 {
		t.Fatalf("unexpected inventory load: err=%v found=%v rec=%+v", err, found, rec)
	}
	rec.Slots["0:1"] = repository.ItemStack{ItemID: 9999}
	cached, found, err := invCache.getOrLoad("1001", now.Add(100*time.Millisecond), func() (repository.InventoryRecord, bool, error) {
		t.Fatalf("expected inventory cache hit")
		return repository.InventoryRecord{}, false, nil
	})
	if err != nil || !found || cached.Slots["0:1"].ItemID != 3037 || invCalls.Load() != 1 {
		t.Fatalf("inventory cache polluted or missed: err=%v found=%v rec=%+v", err, found, cached)
	}
	invCache.invalidate("1001")
	_, found, err = invCache.getOrLoad("1001", now.Add(200*time.Millisecond), func() (repository.InventoryRecord, bool, error) {
		invCalls.Add(1)
		return repository.InventoryRecord{CharacterID: "1001"}, true, nil
	})
	if err != nil || !found || invCalls.Load() != 2 {
		t.Fatalf("expected inventory reload after invalidate: calls=%d", invCalls.Load())
	}

	questCache := newMySQLRecordReadCache(repository.CloneQuest)
	var qCalls atomic.Int32
	q, found, err := questCache.getOrLoad("1001", now, func() (repository.QuestRecord, bool, error) {
		qCalls.Add(1)
		return repository.QuestRecord{CharacterID: "1001"}, true, nil
	})
	if err != nil || !found || q.CharacterID != "1001" {
		t.Fatalf("unexpected quest load: err=%v found=%v q=%+v", err, found, q)
	}
	questCache.invalidate("1001")
	_, found, err = questCache.getOrLoad("1001", now.Add(100*time.Millisecond), func() (repository.QuestRecord, bool, error) {
		qCalls.Add(1)
		return repository.QuestRecord{CharacterID: "1001"}, true, nil
	})
	if err != nil || !found || qCalls.Load() != 2 {
		t.Fatalf("expected quest reload after invalidate: calls=%d", qCalls.Load())
	}
}
