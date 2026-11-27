package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/party"
)

const raidStartGroupSyncEnv = "DNFBRIDGE_RAID_START_GROUP_SYNC"

// raidStartGroupSyncEnabled 默认开启：
//   1. 开团时在 applyRuntimeRaidParties 真正创建普通队伍后，补发一次 NOTI591 攻坚队名册，
//      避免开团前导 NOTI591 在建队前运行导致 1队/2队 全员被投影为 GroupIndex=0（单人）；
//   2. 已开团且成员持有活着的 1 人普通队伍、开团前被团长显式编入有效队伍且尚未退队时，
//      保留其 1队/2队 编号显示，不误投影为单人。
//
// 设 DNFBRIDGE_RAID_START_GROUP_SYNC=0 / false / off 可一键回退到改动前行为。
func raidStartGroupSyncEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(raidStartGroupSyncEnv)))
	return v != "0" && v != "false" && v != "off"
}

// syncCurrentRaidRosterAfterStartParties 在 applyRuntimeRaidParties 建队并发布快照之后调用，
// 补发一次基于实时普通队伍状态的 NOTI591，确保客户端攻坚队列表正确显示 1队/2队 编队。
func (s *Service) syncCurrentRaidRosterAfterStartParties(state runtimeRaidState) error {
	if s == nil || !raidStartGroupSyncEnabled() || s.options.currentClientVersion != "official-24.9" {
		return nil
	}
	if !state.CurrentLayout || !state.Started || state.RouteOrRaidType > 1 {
		return nil
	}
	if state.Battle != nil && (state.Battle.Single || state.Battle.FiendWar != nil) {
		return nil
	}
	return s.broadcastRuntimeRaidRoster(state)
}

// currentRaidMemberRetainsExplicitStartedGroup 判断已开团成员在仅持有 1 人普通队伍时，
// 是否应保留其显式编队号（例如双开 1队1人+2队1人开团，或尚未触发退队标 solo 的 1 人编队）。
// 一旦成员经历退队被 refreshCurrentRaidSoloForPartyMembers 置位 OrdinaryPartySolo=true，
// 或不再属于任何普通队伍（!found / len==0），均回归单人显示，不影响既有退队语义。
func currentRaidMemberRetainsExplicitStartedGroup(
	state runtimeRaidState,
	member runtimeRaidMemberState,
	liveParty party.RuntimePartySnapshot,
	found bool,
) bool {
	if !raidStartGroupSyncEnabled() {
		return false
	}
	if !state.Started || !found || len(liveParty.Members) != 1 {
		return false
	}
	if !member.RecruitingPlacementExplicit || member.OrdinaryPartySolo {
		return false
	}
	return member.GroupIndex > 0 && member.GroupIndex <= currentRaidGroupMaximum(state)
}
