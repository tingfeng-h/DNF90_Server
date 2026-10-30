package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// Docs/CLIENT_PARTY_CLEAR_CONTEXT_V140_20261019.md：客户端 NOTI9 的删除键是
// **包里的 route/owner**，不是本地连接的 owner。
//
// 退图的那一瞬间，退出者的会话还在副本上下文（owner=0），而他的队伍记录是在
// 城镇上下文（owner=10）建的。两边键不一致 ⇒ kind3 找不到记录 ⇒ 队伍界面残留，
// 直到玩家手动退队（那时 owner 已是城镇的 10）才清掉。
//
// 这里对"服务端主动把人摘掉"的退图路径额外补一条 **用收包者当前所选频道
// 作为 owner** 的空名册（kind3 销毁）帧，让客户端能按城镇键找到那条记录。
//
// DNFBRIDGE_PARTY_DUNGEON_EXIT_CLIENT_CHANNEL_CLEAR：默认 on；显式 off/0/false/no 回退。

const partyDungeonExitClientChannelClearEnv = "DNFBRIDGE_PARTY_DUNGEON_EXIT_CLIENT_CHANNEL_CLEAR"

func partyDungeonExitClientChannelClearEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyDungeonExitClientChannelClearEnv))) {
	case "off", "0", "false", "no":
		return false
	default:
		return true
	}
}

func (s *Service) sendDungeonExitSelfFrameClear(session *gameSession, source string) error {
	if s == nil || session == nil || !partyDungeonExitClientChannelClearEnabled() {
		return nil
	}
	sent, err := s.sendCurrentPartyActorFrameProjectionInContext(
		session,
		session,
		alignedcmd.PartyState{},
		source+"_client_channel_owner",
		nil,
		currentPartyFrameOwnerSelection{Mode: currentPartyFrameOwnerClientChannel},
	)
	if err != nil {
		s.logGameEvent(session, "game-party-dungeon-exit-frame-clear-failed",
			"source", source, "char_id", selectedCharacterID(session), "error", err)
		return err
	}
	s.logGameEvent(session, "game-party-dungeon-exit-frame-clear",
		"source", source,
		"char_id", selectedCharacterID(session),
		"owner_mode", "client_channel",
		"owner_channel", currentClientSelectedChannel(session),
		"sent", sent,
		"switch", partyDungeonExitClientChannelClearEnv,
		"reason", "client_noti9_delete_key_is_the_packet_owner_so_the_town_context_record_needs_a_town_owner")
	return nil
}
