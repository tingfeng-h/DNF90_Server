package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// 移植自「轻装上阵」分支 party_dungeon_actor_leave.go 的 host-transfer 段：
// 队长离开/解散时，先用原生 op26 + 槽位告诉留图者"这个槽位由你接管"，
// 再走已有的 op6 活动 actor 移除。op9 名册清除不能替代这两步——
// 它只清队伍表，副本场景投影出来的 actor 仍然可见。
//
//	DNFBRIDGE_PARTY_ACTOR_LEAVE_OP26=off   # 默认开启；置 off 回退
const partyActorLeaveOp26Env = "DNFBRIDGE_PARTY_ACTOR_LEAVE_OP26"

func partyActorLeaveOp26Enabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyActorLeaveOp26Env))) {
	case "0", "off", "false", "no":
		return false
	default:
		return true
	}
}

// sendDungeonPartyHostTransferToSurvivors 给留图者各发一条 op26（自身槽位）。
func (s *Service) sendDungeonPartyHostTransferToSurvivors(
	instance *dungeonPartyInstance,
	leavingCharacterID uint16,
	source string,
) int {
	if s == nil || instance == nil || leavingCharacterID == 0 || !partyActorLeaveOp26Enabled() {
		return 0
	}
	if s.options.currentClientVersion != "official-24.9" {
		return 0
	}
	peers := dungeonPartyInstanceMemberSessions(instance, leavingCharacterID)
	if len(peers) == 0 {
		return 0
	}
	var oldMembers []alignedcmd.PartyMemberState
	if leaving, online := s.onlineGameSession(leavingCharacterID); online {
		state := s.runtimePartyStateForSession(leaving)
		if state.UserID == leavingCharacterID {
			oldMembers = runtimePartyMembers(state)
		}
	}
	if len(oldMembers) == 0 {
		return 0
	}
	sent := 0
	for _, peer := range peers {
		if peer.Session == nil || peer.CharacterID == 0 {
			continue
		}
		for _, member := range oldMembers {
			if member.UserID != peer.CharacterID {
				continue
			}
			if err := s.sendGameUpperRawClass(peer.Session, 26, []byte{member.Slot}, 0); err != nil {
				s.logGameEvent(peer.Session, "game-party-dungeon-host-transfer-failed",
					"source", source, "char_id", peer.CharacterID,
					"leaving_char_id", leavingCharacterID, "host_slot", member.Slot, "error", err)
				break
			}
			sent++
			s.logGameEvent(peer.Session, "game-party-dungeon-host-transfer",
				"source", source, "char_id", peer.CharacterID,
				"leaving_char_id", leavingCharacterID, "host_slot", member.Slot,
				"body_source", "fork_native_op26_slot_host_transfer")
			break
		}
	}
	return sent
}
