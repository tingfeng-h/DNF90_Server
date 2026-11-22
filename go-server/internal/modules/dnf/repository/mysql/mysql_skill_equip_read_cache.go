package mysql

import (
	"sync"
	"time"
)

const mysqlSkillEquipReadCacheTTL = time.Second

type mysqlRecordReadCacheEntry[T any] struct {
	record    T
	expiresAt time.Time
}

type mysqlRecordReadCacheFlight[T any] struct {
	done   chan struct{}
	record T
	found  bool
	err    error
}

// mysqlRecordReadCache 复用 P0 角色只读缓存的相同安全契约（1s TTL + singleflight + 写入口 0ms 失效），
// 仅在非事务只读路径生效，所有事务内读改写（FOR UPDATE / REPEATABLE READ / *sqlTransactionDB）
// 天然绕过，覆盖 op37 技能点与进图/换房装备高频只读查询。
type mysqlRecordReadCache[T any] struct {
	mu          sync.Mutex
	clone       func(T) T
	entries     map[string]mysqlRecordReadCacheEntry[T]
	flights     map[string]*mysqlRecordReadCacheFlight[T]
	generations map[string]uint64
}

func newMySQLRecordReadCache[T any](clone func(T) T) *mysqlRecordReadCache[T] {
	return &mysqlRecordReadCache[T]{
		clone:       clone,
		entries:     make(map[string]mysqlRecordReadCacheEntry[T], 64),
		flights:     make(map[string]*mysqlRecordReadCacheFlight[T], 16),
		generations: make(map[string]uint64, 64),
	}
}

func (c *mysqlRecordReadCache[T]) getOrLoad(
	key string,
	now time.Time,
	loader func() (T, bool, error),
) (T, bool, error) {
	var zero T
	if c == nil || key == "" {
		return loader()
	}
	c.mu.Lock()
	if entry, ok := c.entries[key]; ok && now.Before(entry.expiresAt) {
		c.mu.Unlock()
		return c.clone(entry.record), true, nil
	}
	if flight, inflight := c.flights[key]; inflight {
		c.mu.Unlock()
		<-flight.done
		if flight.err != nil || !flight.found {
			return zero, flight.found, flight.err
		}
		return c.clone(flight.record), true, nil
	}
	startGen := c.generations[key]
	flight := &mysqlRecordReadCacheFlight[T]{done: make(chan struct{})}
	c.flights[key] = flight
	c.mu.Unlock()

	record, found, err := loader()
	flight.record = record
	flight.found = found
	flight.err = err

	c.mu.Lock()
	if err == nil && found && c.generations[key] == startGen {
		if len(c.entries) > 512 {
			pruneNow := time.Now()
			for k, v := range c.entries {
				if !pruneNow.Before(v.expiresAt) {
					delete(c.entries, k)
				}
			}
		}
		c.entries[key] = mysqlRecordReadCacheEntry[T]{
			record:    c.clone(record),
			expiresAt: time.Now().Add(mysqlSkillEquipReadCacheTTL),
		}
	}
	delete(c.flights, key)
	close(flight.done)
	c.mu.Unlock()

	if err != nil || !found {
		return zero, found, err
	}
	return c.clone(record), true, nil
}

func (c *mysqlRecordReadCache[T]) invalidate(key string) {
	if c == nil || key == "" {
		return
	}
	c.mu.Lock()
	c.generations[key]++
	delete(c.entries, key)
	c.mu.Unlock()
}
