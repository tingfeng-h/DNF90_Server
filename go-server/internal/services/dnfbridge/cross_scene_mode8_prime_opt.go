package dnfbridge

import (
	"context"
	"strings"
	"sync"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const (
	// crossSceneMode8PrimeOptEnv 控制跨房 mode8 预热的按需触发与会话级去重优化。
	// 默认开启：
	//   1) 禁止新角色进频道时的被动在线广播（presence/online/login）主动推送 mode8，
	//      彻底消除"未同屏突然 loading"及人多时城镇连续卡顿；
	//   2) 同一 viewer 对同一 (peer, channel, appearanceRevision) 仅推送 1 次，
	//      避免同公会 + 冒险团好友 + 角色好友多入口重复推包。
	// 设为 off/false/0/no/legacy 时逐字节回退至旧版行为。
	crossSceneMode8PrimeOptEnv = "DNFBRIDGE_CROSS_SCENE_MODE8_PRIME_OPT"

	// crossSceneMode8LightweightEnv 控制跨房 mode8 内嵌 mode0 是否使用轻量外壳
	// （保留职业/名称/公会/冒险团/基础头像，剥离宠物/名片/14槽重型时装特效与多余查库）。
	// 设为 off/false/0/no/legacy 时回退使用完整 mode0。
	crossSceneMode8LightweightEnv = "DNFBRIDGE_CROSS_SCENE_MODE8_LIGHTWEIGHT"

	crossSceneMode8PrimeCacheMaxEntries = 256
)

type crossSceneMode8PrimeKey struct {
	peerCharID uint16
	channel    int
}

type crossSceneMode8ViewerCache struct {
	mu         sync.Mutex
	character  uint16
	generation uint64
	revisions  map[crossSceneMode8PrimeKey]uint64
}

// 按 *gameSession 指针隔离，不改动 gameSession/Service 结构体字段。
var crossSceneMode8ViewerCaches sync.Map // key: *gameSession -> *crossSceneMode8ViewerCache

func crossSceneMode8PrimeOptEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(envString(crossSceneMode8PrimeOptEnv, ""))) {
	case "off", "false", "0", "no", "legacy":
		return false
	default:
		return true
	}
}

func crossSceneMode8LightweightEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(envString(crossSceneMode8LightweightEnv, ""))) {
	case "off", "false", "0", "no", "legacy":
		return false
	default:
		return true
	}
}

// crossSceneMode8PrimeAllowedForSource 判定当前刷新来源是否允许触发跨房 mode8 预热。
// 仅在玩家主动打开面板/执行变更（如 client、mutation、guild_chat）时按需触发；
// 被动在线广播（presence、online、offline、character_switched、登录尾自动刷新等）
// 只刷新在线状态与频道号，不再推送重型 mode8 外观包。
func crossSceneMode8PrimeAllowedForSource(source string) bool {
	if !crossSceneMode8PrimeOptEnabled() {
		return true
	}
	normalized := strings.ToLower(strings.TrimSpace(source))
	switch normalized {
	case "client", "mutation", "guild_chat", "party_invite":
		return true
	default:
		return false
	}
}

func crossSceneMode8CacheForViewer(session *gameSession) *crossSceneMode8ViewerCache {
	if session == nil {
		return nil
	}
	if existing, ok := crossSceneMode8ViewerCaches.Load(session); ok {
		return existing.(*crossSceneMode8ViewerCache)
	}
	created := &crossSceneMode8ViewerCache{
		revisions: make(map[crossSceneMode8PrimeKey]uint64),
	}
	actual, _ := crossSceneMode8ViewerCaches.LoadOrStore(session, created)
	return actual.(*crossSceneMode8ViewerCache)
}

func crossSceneMode8AlreadyPrimed(viewer boundGameSessionCharacter, peerCharID uint16, channel int, revision uint64) bool {
	if !crossSceneMode8PrimeOptEnabled() || viewer.session == nil || peerCharID == 0 {
		return false
	}
	cache := crossSceneMode8CacheForViewer(viewer.session)
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.character != viewer.character || cache.generation != viewer.generation {
		return false
	}
	prev, ok := cache.revisions[crossSceneMode8PrimeKey{peerCharID: peerCharID, channel: channel}]
	return ok && prev == revision
}

func markCrossSceneMode8Primed(viewer boundGameSessionCharacter, peerCharID uint16, channel int, revision uint64) {
	if !crossSceneMode8PrimeOptEnabled() || viewer.session == nil || peerCharID == 0 {
		return
	}
	cache := crossSceneMode8CacheForViewer(viewer.session)
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.character != viewer.character || cache.generation != viewer.generation {
		cache.character = viewer.character
		cache.generation = viewer.generation
		cache.revisions = make(map[crossSceneMode8PrimeKey]uint64)
	}
	if len(cache.revisions) >= crossSceneMode8PrimeCacheMaxEntries {
		cache.revisions = make(map[crossSceneMode8PrimeKey]uint64)
	}
	cache.revisions[crossSceneMode8PrimeKey{peerCharID: peerCharID, channel: channel}] = revision
}

// buildCurrentLightweightCrossSceneMode0 为跨房 mode8 缓存构造轻量 mode0：
// 保留职业、转职、等级、角色名、冒险团名、公会名/图标及名册基础头像装备，
// 跳过宠物实体、名片装饰、技能树及 14 槽全套重型时装/光环/特效查库，
// 大幅降低客户端主线程同步加载 NPK 贴图造成的 loading 卡顿。
func (s *Service) buildCurrentLightweightCrossSceneMode0(
	ctx context.Context,
	peerSession *gameSession,
	peerCharID uint16,
	peerName string,
	character dnfrepo.CharacterRecord,
	ownerChannel byte,
) []byte {
	if peerName == "" {
		peerName = csharpSelectCharacterName(peerCharID)
	}
	objectKey := currentSceneActorObjectKey(peerCharID)
	adventureName := s.currentSceneObjectAdventureOverheadName(ctx, peerSession, peerCharID, character, true)
	guild, err := s.currentSceneGuild(ctx, peerCharID)
	if err != nil {
		guild = currentSceneGuildProjection{}
	}
	// 仅保留角色名册基础外观，不附加宠物实体（currentEquippedCreatureSnapshot{}）。
	return buildCurrentSceneObjectListBodyWithCreatureAndAdventureNameInContext(
		objectKey,
		character,
		true,
		peerName,
		currentEquippedCreatureSnapshot{},
		adventureName,
		ownerChannel,
		s.currentSceneActorRoute(),
		guild,
	)
}

