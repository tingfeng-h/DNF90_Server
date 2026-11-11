package dnfbridge

import (
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// bug1RetireDungeonPartyAliasesEnv 控制「回城重建城镇组队投影之前，先把副本作用域
// 的别名 op9 记录用 kind3 退役掉」。
//
// 默认（未设置 / on）= 退役；legacy / off / 0 / false / no = 回到改动前，只发
// 城镇 R_self、不退役别名。
//
// 对应 dnf90_party_dungeon_bugfix_prompt.md 修改点 4 第 1 条。
const bug1RetireDungeonPartyAliasesEnv = "DNFBRIDGE_PARTY_TOWN_RETURN_RETIRE_ALIASES"

func bug1RetireDungeonPartyAliasesEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(bug1RetireDungeonPartyAliasesEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// retireDungeonPartyAliasRecords 在回城后重建城镇组队投影之前，把副本作用域
// （owner = 0）的别名 op9 记录逐条用 kind3 退役：**先队友，最后自己**。
//
// 为什么必须做：客户端 NOTI9 只有在 kind3 的 (objectKey, owner) 与
// [local_actor+0x4C8] 当前指向的 CPartyRecord 完全一致时，才会走退队清场
// （清空 0x56A3768 与 16 个血条槽位）。图内 3->2 时 survivor 收到了队友的
// 别名 op9，[local_actor+0x4C8] 已经被改写成了**队友**的记录；回城后只发
// 城镇 R_self 的话，这个绑定对不上，客户端会误入"首次入队"分支，`.gld127`
// 调 ResetAllPeerConnections 把原本正常的两人 UDP 槽位全部抹成 -1，而随后
// 的重连循环又因为队友已在旧槽位快照里而跳过，于是再次拉人时卡"连接中"并
// 最终爆防火墙错误。
//
// 先退役队友别名、再退役自己，最后才发城镇 R_self，让 [local_actor+0x4C8]
// 严格绑定回自己的记录。
func (s *Service) retireDungeonPartyAliasRecords(receiver *gameSession, members []alignedcmd.PartyMemberState, source string) error {
	if s == nil || receiver == nil || receiver.selectedCharacterID == 0 {
		return nil
	}
	route := s.currentSceneActorRoute()
	retire := func(characterID uint16) error {
		objectKey := currentSceneActorObjectKey(characterID)
		if objectKey == 0 {
			return nil
		}
		body := buildCurrentSceneOp9ActorRemovalBodyInContext(
			objectKey,
			currentSceneObjectContext,
			route,
		)
		return s.sendGameUpperRawClass(
			receiver,
			uint16(dnfenum.CmdPacketRecoverStamina),
			body,
			0,
		)
	}
	retiredPeers := 0
	for _, member := range members {
		if member.UserID == receiver.selectedCharacterID {
			continue
		}
		if err := retire(member.UserID); err != nil {
			return err
		}
		retiredPeers++
	}
	// Self last: the immediately following town roster must be the record the
	// client ends up bound to.
	if err := retire(receiver.selectedCharacterID); err != nil {
		return err
	}
	s.logGameEvent(receiver, "game-party-town-return-dungeon-alias-retired",
		"source", source,
		"char_id", receiver.selectedCharacterID,
		"retired_peers", retiredPeers,
		"owner_channel", currentSceneObjectContext,
		"switch", bug1RetireDungeonPartyAliasesEnv,
		"reason", "kind3_owner0_peers_then_self_before_town_self_record")
	return nil
}
