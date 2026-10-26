package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// 退队时只删了队伍框（op9 / NOTI9），没有删 3D actor（op6 / NOTI6），
// 客户端 actor 表会保留同 object_key（= 角色 ID）的旧记录。下一次再组队进同一个地图
// （map_id 不变，客户端认为"场景没换"、不重建）时，新的创建帧被旧记录吞掉 ——
// 现场表现：不重选角色就看不到队友，重选角色后必好（重选会整体重建场景/actor 表）。
//
// 这里在退队确认后，把双方互相的 3D actor 各清一次：
//   幸存者 ← 删除离队者；离队者 ← 删除幸存者
// 只在退队这一刻发一次，不动副本内的投影路径（之前"每次投影前都删"会把血条一起删掉）。
//
// 开关：DNFBRIDGE_PARTY_LEAVE_ACTOR_REMOVAL
//
//	默认（空/未设置/非法值）= on；
//	off / 0 / false / no / legacy = 逐字节回到改动前（退队只删队伍框）。
const partyLeaveActorRemovalEnv = "DNFBRIDGE_PARTY_LEAVE_ACTOR_REMOVAL"

func partyLeaveActorRemovalEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyLeaveActorRemovalEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// cleanupPartyLeaveActors 在退队确认后清理双方的 3D actor 投影。
// previous 是离队前的队伍快照；sourceID 是离队角色。
func (s *Service) cleanupPartyLeaveActors(session *gameSession, previous alignedcmd.PartyState, sourceID uint16) {
	if s == nil || session == nil || sourceID == 0 {
		return
	}
	if !partyLeaveActorRemovalEnabled() {
		return
	}
	if s.options.currentClientVersion != "official-24.9" {
		return
	}
	for _, member := range runtimePartyMembers(previous) {
		if member.UserID == 0 || member.UserID == sourceID {
			continue
		}
		if peer, online := s.onlineGameSession(member.UserID); online && peer != nil {
			if err := s.sendDungeonDepartedActorRemoval(peer, sourceID, "client_leave_party_op13_survivor"); err != nil {
				s.logGameEvent(peer, "game-party-leave-actor-remove-failed",
					"departing_char_id", sourceID, "receiver_char_id", member.UserID, "error", err)
			}
		}
		if err := s.sendDungeonDepartedActorRemoval(session, member.UserID, "client_leave_party_op13_retired_peer"); err != nil {
			s.logGameEvent(session, "game-party-leave-actor-remove-failed",
				"departing_char_id", member.UserID, "receiver_char_id", sourceID, "error", err)
		}
	}
	s.logGameEvent(session, "game-party-leave-actors-cleared",
		"source_char_id", sourceID,
		"member_count", len(runtimePartyMembers(previous)),
		"env", partyLeaveActorRemovalEnv,
		"mode", "bidirectional")
}
