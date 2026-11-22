package dnfbridge

import (
	"sync"
	"time"

	"longheng.io/server/internal/modules/dnf/buffswap"
)

// buffSwapReapplyCacheTTL 覆盖副本连续换房窗口（30s），避免每次 op37 换房
// 都为 group 0/1 各开一次 REPEATABLE READ MySQL 事务（40ms/人）。
// 任何登记/穿脱装备路径（diagnostics == true）会先显式失效本缓存与队友投影缓存，
// 保证换装与改登记 0ms 立即生效。
const buffSwapReapplyCacheTTL = 30 * time.Second

type buffSwapReapplyCacheKey struct {
	character  uint16
	generation uint64
	hasGroup   bool
	group      byte
}

type buffSwapReapplyCacheEntry struct {
	found     bool
	state     buffswap.State
	omitted   int
	expiresAt time.Time
}

type buffSwapReapplyCache struct {
	mu      sync.Mutex
	entries map[buffSwapReapplyCacheKey]buffSwapReapplyCacheEntry
}

func (s *Service) lookupBuffSwapReapplyCache(key buffSwapReapplyCacheKey, now time.Time) (buffSwapReapplyCacheEntry, bool) {
	if s == nil || key.character == 0 || key.generation == 0 {
		return buffSwapReapplyCacheEntry{}, false
	}
	s.buffSwapReapplyMemo.mu.Lock()
	defer s.buffSwapReapplyMemo.mu.Unlock()
	if s.buffSwapReapplyMemo.entries == nil {
		return buffSwapReapplyCacheEntry{}, false
	}
	entry, ok := s.buffSwapReapplyMemo.entries[key]
	if !ok || !now.Before(entry.expiresAt) {
		if ok {
			delete(s.buffSwapReapplyMemo.entries, key)
		}
		return buffSwapReapplyCacheEntry{}, false
	}
	return entry, true
}

func (s *Service) storeBuffSwapReapplyCache(key buffSwapReapplyCacheKey, entry buffSwapReapplyCacheEntry, now time.Time) {
	if s == nil || key.character == 0 || key.generation == 0 {
		return
	}
	entry.expiresAt = now.Add(buffSwapReapplyCacheTTL)
	s.buffSwapReapplyMemo.mu.Lock()
	defer s.buffSwapReapplyMemo.mu.Unlock()
	if s.buffSwapReapplyMemo.entries == nil {
		s.buffSwapReapplyMemo.entries = make(map[buffSwapReapplyCacheKey]buffSwapReapplyCacheEntry, 32)
	}
	if len(s.buffSwapReapplyMemo.entries) > 256 {
		for k, v := range s.buffSwapReapplyMemo.entries {
			if !now.Before(v.expiresAt) {
				delete(s.buffSwapReapplyMemo.entries, k)
			}
		}
	}
	s.buffSwapReapplyMemo.entries[key] = entry
}

func (s *Service) invalidateBuffSwapReapplyCache(character uint16) {
	if s == nil || character == 0 {
		return
	}
	s.buffSwapReapplyMemo.mu.Lock()
	defer s.buffSwapReapplyMemo.mu.Unlock()
	for k := range s.buffSwapReapplyMemo.entries {
		if k.character == character {
			delete(s.buffSwapReapplyMemo.entries, k)
		}
	}
}
