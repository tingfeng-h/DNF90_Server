package mysql

import (
	"sync"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

const mysqlCharReadCacheTTL = time.Second

type mysqlCharReadCacheEntry struct {
	record    repository.CharacterRecord
	expiresAt time.Time
}

type mysqlCharReadCacheFlight struct {
	done   chan struct{}
	record repository.CharacterRecord
	found  bool
	err    error
}

// mysqlCharReadCache 是挂载在 mysqlRouter 上的非事务只读角色缓存（1s TTL + singleflight）。
//   - 仅对纯非事务只读（!lockReads && !consistentRead && 非 *sqlTransactionDB）生效，
//     自动覆盖 fillPrivateCombatSnapshot / partyFrameCharacterBatch / currentWeddingCharacter /
//     characterfriend / selectedCharacterForEnter 等所有绕过上层的热点读；
//   - 所有事务内读改写（结算 / FOR UPDATE / REPEATABLE READ）天然绕过缓存，保证严格 ACID；
//   - 任何角色写入口（含事务内 SaveFields / CreateCharacter / SaveProgressionCharacter /
//     SaveExpertJobExperience / AdvanceStoryDigest / swapCharacterSlots）均会立即失效该角色键，
//     保证换装、升级、结算后零延迟读到最新值。
type mysqlCharReadCache struct {
	mu          sync.Mutex
	entries     map[string]mysqlCharReadCacheEntry
	flights     map[string]*mysqlCharReadCacheFlight
	generations map[string]uint64
}

func newMySQLCharReadCache() *mysqlCharReadCache {
	return &mysqlCharReadCache{
		entries:     make(map[string]mysqlCharReadCacheEntry, 64),
		flights:     make(map[string]*mysqlCharReadCacheFlight, 16),
		generations: make(map[string]uint64, 64),
	}
}

func (r mysqlRouter) canUseCharReadCache() bool {
	if r.charReadCache == nil || r.lockReads || r.consistentRead {
		return false
	}
	_, inTx := r.db.(*sqlTransactionDB)
	return !inTx
}

func (c *mysqlCharReadCache) getOrLoad(
	characterID string,
	now time.Time,
	loader func() (repository.CharacterRecord, bool, error),
) (repository.CharacterRecord, bool, error) {
	if c == nil || characterID == "" {
		return loader()
	}
	c.mu.Lock()
	if entry, ok := c.entries[characterID]; ok && now.Before(entry.expiresAt) {
		c.mu.Unlock()
		return repository.CloneCharacter(entry.record), true, nil
	}
	if flight, inflight := c.flights[characterID]; inflight {
		c.mu.Unlock()
		<-flight.done
		if flight.err != nil || !flight.found {
			return repository.CharacterRecord{}, flight.found, flight.err
		}
		return repository.CloneCharacter(flight.record), true, nil
	}
	startGen := c.generations[characterID]
	flight := &mysqlCharReadCacheFlight{done: make(chan struct{})}
	c.flights[characterID] = flight
	c.mu.Unlock()

	record, found, err := loader()
	flight.record = record
	flight.found = found
	flight.err = err

	c.mu.Lock()
	// 仅在加载期间没有发生任何写入失效（generation 未变）时才写入缓存，防止旧读覆盖新写。
	if err == nil && found && c.generations[characterID] == startGen {
		if len(c.entries) > 512 {
			pruneNow := time.Now()
			for k, v := range c.entries {
				if !pruneNow.Before(v.expiresAt) {
					delete(c.entries, k)
				}
			}
		}
		c.entries[characterID] = mysqlCharReadCacheEntry{
			record:    repository.CloneCharacter(record),
			expiresAt: time.Now().Add(mysqlCharReadCacheTTL),
		}
	}
	delete(c.flights, characterID)
	close(flight.done)
	c.mu.Unlock()

	if err != nil || !found {
		return repository.CharacterRecord{}, found, err
	}
	return repository.CloneCharacter(record), true, nil
}

func (c *mysqlCharReadCache) invalidate(characterID string) {
	if c == nil || characterID == "" {
		return
	}
	c.mu.Lock()
	c.generations[characterID]++
	delete(c.entries, characterID)
	c.mu.Unlock()
}

func (c *mysqlCharReadCache) invalidateAll() {
	if c == nil {
		return
	}
	c.mu.Lock()
	for k := range c.entries {
		c.generations[k]++
		delete(c.entries, k)
	}
	c.mu.Unlock()
}
