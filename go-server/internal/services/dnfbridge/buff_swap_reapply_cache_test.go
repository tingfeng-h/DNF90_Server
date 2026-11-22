package dnfbridge

import (
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/buffswap"
)

func TestBuffSwapReapplyCacheHitExpireInvalidate(t *testing.T) {
	s := &Service{}
	key := buffSwapReapplyCacheKey{character: 1001, generation: 3, hasGroup: true, group: 0}
	now := time.Unix(1700000000, 0)

	if _, ok := s.lookupBuffSwapReapplyCache(key, now); ok {
		t.Fatalf("expected initial cache miss")
	}
	s.storeBuffSwapReapplyCache(key, buffSwapReapplyCacheEntry{
		found: true,
		state: buffswap.State{WireProfile: buffswap.CurrentWireProfile, Skill: 14},
	}, now)

	got, ok := s.lookupBuffSwapReapplyCache(key, now.Add(10*time.Second))
	if !ok || !got.found || got.state.Skill != 14 {
		t.Fatalf("expected cache hit within 30s TTL, got ok=%v entry=%+v", ok, got)
	}
	if _, ok := s.lookupBuffSwapReapplyCache(key, now.Add(buffSwapReapplyCacheTTL+time.Millisecond)); ok {
		t.Fatalf("expected cache miss after TTL expiry")
	}

	s.storeBuffSwapReapplyCache(key, buffSwapReapplyCacheEntry{found: true, state: buffswap.State{Skill: 14}}, now)
	s.invalidateBuffSwapReapplyCache(1001)
	if _, ok := s.lookupBuffSwapReapplyCache(key, now.Add(time.Second)); ok {
		t.Fatalf("expected cache miss after invalidation")
	}
}
