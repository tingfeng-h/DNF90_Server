package dnfbridge

import (
	"context"
	"testing"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

func TestCrossSceneMode8PrimeSourceFilterAndRollbackSwitch(t *testing.T) {
	// 默认开启：仅允许按需来源，屏蔽新角色上线/切频道等被动广播来源。
	for _, allowed := range []string{"client", "mutation", "guild_chat", "party_invite", " CLIENT "} {
		if !crossSceneMode8PrimeAllowedForSource(allowed) {
			t.Fatalf("expected source %q to be allowed", allowed)
		}
	}
	for _, blocked := range []string{
		"presence",
		"online",
		"offline",
		"character_switched",
		"channel",
		"login_after_adventure_actor_before_scene_commit",
		"finish_loading_after_player_state",
		"",
	} {
		if crossSceneMode8PrimeAllowedForSource(blocked) {
			t.Fatalf("expected passive source %q to be suppressed", blocked)
		}
	}

	// 回滚开关 off：全部放行，逐字节退回旧行为。
	t.Setenv(crossSceneMode8PrimeOptEnv, "off")
	for _, src := range []string{"presence", "online", "login_after_adventure_actor_before_scene_commit"} {
		if !crossSceneMode8PrimeAllowedForSource(src) {
			t.Fatalf("rollback switch did not allow source %q", src)
		}
	}
}

func TestCrossSceneMode8SessionDeduplicationAndInvalidation(t *testing.T) {
	viewerSession := &gameSession{connID: "viewer-dedup"}
	viewer := boundGameSessionCharacter{session: viewerSession, character: 101, generation: 1}

	if crossSceneMode8AlreadyPrimed(viewer, 202, 11, 7) {
		t.Fatal("fresh viewer unexpectedly marked primed")
	}
	markCrossSceneMode8Primed(viewer, 202, 11, 7)
	if !crossSceneMode8AlreadyPrimed(viewer, 202, 11, 7) {
		t.Fatal("same (viewer, peer, channel, revision) was not deduplicated")
	}

	// peer 换装（revision 变化）、换频道、不同 peer、或 viewer 切角色（generation 变化）均自动失效。
	if crossSceneMode8AlreadyPrimed(viewer, 202, 11, 8) {
		t.Fatal("stale revision remained deduplicated")
	}
	if crossSceneMode8AlreadyPrimed(viewer, 202, 12, 7) {
		t.Fatal("channel change remained deduplicated")
	}
	if crossSceneMode8AlreadyPrimed(viewer, 203, 11, 7) {
		t.Fatal("different peer wrongly deduplicated")
	}
	switchedViewer := boundGameSessionCharacter{session: viewerSession, character: 102, generation: 2}
	if crossSceneMode8AlreadyPrimed(switchedViewer, 202, 11, 7) {
		t.Fatal("switched viewer character inherited old dedup cache")
	}

	// 回滚开关 off：不去重。
	t.Setenv(crossSceneMode8PrimeOptEnv, "off")
	if crossSceneMode8AlreadyPrimed(viewer, 202, 11, 7) {
		t.Fatal("rollback switch did not disable deduplication")
	}
}

func TestCrossSceneMode8LightweightMode0BuildsValidMode8Frame(t *testing.T) {
	s := &Service{}
	peerSession := &gameSession{connID: "peer-light"}
	character := dnfrepo.CharacterRecord{
		CharacterID: "202",
		AccountID:   "acc-202",
		Name:        "TestPeer",
		Level:       90,
		Job:         "0",
	}
	mode0 := s.buildCurrentLightweightCrossSceneMode0(context.Background(), peerSession, 202, "TestPeer", character, 11)
	if len(mode0) <= currentCrossSceneMode0KeyOffset+2 || mode0[0] != currentCrossSceneMode0Subtype {
		t.Fatalf("invalid lightweight mode0 header: len=%d first=%v", len(mode0), mode0)
	}
	mode3 := s.buildCurrentSelectedUserInfoMode3BodyInContext(context.Background(), peerSession, nil, character, true, 202, 11)
	frame, err := buildCurrentCrossSceneUserInfoBody(mode0, mode3)
	if err != nil {
		t.Fatalf("buildCurrentCrossSceneUserInfoBody failed with lightweight mode0: %v", err)
	}
	if len(frame) == 0 || frame[0] != currentCrossSceneMode8Subtype {
		t.Fatalf("unexpected mode8 frame: len=%d", len(frame))
	}

	// 验证轻量开关回退。
	if !crossSceneMode8LightweightEnabled() {
		t.Fatal("lightweight mode8 should be enabled by default")
	}
	t.Setenv(crossSceneMode8LightweightEnv, "off")
	if crossSceneMode8LightweightEnabled() {
		t.Fatal("lightweight rollback switch failed")
	}
}
