package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/raid"
)

// 换频道/重连后客户端会丢掉团本 UI 上下文：没有"邀请攻坚队"入口、别人的申请
// 也显示不出来。现场（2026-10-26 01:31）：连续换 3 次频道（11→10→23，队伍已解散），
// 客户端只在点开界面时才发 raid-info 请求，服务端不主动推 → UI 一直空着，
// 直到 01:33:23 某次完整状态被推下去才恢复。
//
// 这里在登录落点序列（initial town 的 raid channel ready 之后）主动补一次完整
// 团本状态：profiles → info ACK → roster → 入团消耗，只对"确实是该团成员"的角色发。
//
// 开关：DNFBRIDGE_RAID_CHANNEL_RESYNC
//
//	默认（空/未设置/非法值）= on；
//	off / 0 / false / no / legacy = 逐字节回到改动前（只在客户端请求时回应）。
const raidChannelResyncEnv = "DNFBRIDGE_RAID_CHANNEL_RESYNC"

func raidChannelResyncEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(raidChannelResyncEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// pushCurrentRaidMembershipAfterChannelResync 在登录/重连落点主动重推当前团本状态。
// 只对当前团成员、且该团状态对本会话可见时发送；否则一个字节都不发。
func (s *Service) pushCurrentRaidMembershipAfterChannelResync(session *gameSession, source string) error {
	if s == nil || session == nil || session.selectedCharacterID == 0 {
		return nil
	}
	if !raidChannelResyncEnabled() {
		return nil
	}
	if s.options.currentClientVersion != "official-24.9" {
		return nil
	}
	charID := session.selectedCharacterID
	s.raidMu.Lock()
	live := s.runtimeRaidByMemberLocked(charID)
	var state runtimeRaidState
	if live != nil {
		state = cloneRuntimeRaidState(*live)
	}
	s.raidMu.Unlock()
	if state.RaidKey == 0 || runtimeRaidMemberIndex(state, charID) < 0 {
		return nil
	}
	if !s.currentRaidInfoVisible(session, state) {
		return nil
	}
	if state.Battle != nil && state.Battle.MemberAccounts[charID] != s.accountIDForSession(session) {
		return nil
	}
	reply, err := raid.BuildCurrentInfoReplyBody(currentRuntimeRaidMetadata(state), 2)
	if err != nil {
		return err
	}
	if err := s.sendCurrentRuntimeRaidProfiles(session, state); err != nil {
		return err
	}
	if err := s.sendGameUpperRawClassCodec(session, uint16(dnfenum.CmdPacketRequestRaidInfo), reply, 1, true); err != nil {
		return err
	}
	roster, err := s.buildRuntimeRaidRosterBody(state)
	if err != nil {
		return err
	}
	if err := s.sendGameUpperRawClassCodec(session, raid.RaidMemberRefreshMsgID, roster, 0, true); err != nil {
		return err
	}
	if err := s.sendCurrentRuntimeRaidEntryCost(session, state); err != nil {
		return err
	}
	s.logGameEvent(session, "game-raid-channel-resync-push",
		"char_id", charID,
		"raid_key", state.RaidKey,
		"source", source,
		"env", raidChannelResyncEnv,
		"mode", "push")
	return nil
}
