package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/dnfenum"
	dnfproto "longheng.io/server/internal/modules/dnf/protocol"
)

// 退图时服务端把玩家从队伍里摘掉（队员退图自己退队 / 队长退图解散），
// 但**从来没发过退队 ACK**。客户端的队伍界面是由这条 ACK 驱动的：
// 手动退队（op13）路径在 manager.Leave 之后发
//   sendGameUpperRawClassCodec(session, CmdPacketLeaveParty, []byte{1}, DefaultChannelClassification, true)
// 所以手动退队不留残留；而退图路径只发了名册销毁帧 + 空 op153，
// 客户端始终认为自己还在队伍里 —— 就是“谁先退图谁残留”。
//
// DNFBRIDGE_PARTY_DUNGEON_EXIT_LEAVE_ACK：默认 on（补齐 ACK）；显式 off/0/false/no 回退。

const partyDungeonExitLeaveAckEnv = "DNFBRIDGE_PARTY_DUNGEON_EXIT_LEAVE_ACK"

func partyDungeonExitLeaveAckEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyDungeonExitLeaveAckEnv))) {
	case "off", "0", "false", "no":
		return false
	default:
		return true
	}
}

func (s *Service) sendDungeonExitPartyLeaveAck(session *gameSession, source string) error {
	if s == nil || session == nil || !partyDungeonExitLeaveAckEnabled() {
		return nil
	}
	if err := s.sendGameUpperRawClassCodec(
		session,
		uint16(dnfenum.CmdPacketLeaveParty),
		[]byte{1},
		dnfproto.DefaultChannelClassification,
		true,
	); err != nil {
		s.logGameEvent(session, "game-party-dungeon-exit-leave-ack-failed",
			"source", source, "char_id", selectedCharacterID(session), "error", err)
		return err
	}
	s.logGameEvent(session, "game-party-dungeon-exit-leave-ack",
		"source", source,
		"char_id", selectedCharacterID(session),
		"msg_id", uint16(dnfenum.CmdPacketLeaveParty),
		"body", "01",
		"switch", partyDungeonExitLeaveAckEnv,
		"reason", "server_removed_the_member_on_dungeon_exit_so_the_client_needs_the_native_leave_ack")
	return nil
}
