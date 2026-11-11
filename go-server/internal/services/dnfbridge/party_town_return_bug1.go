package dnfbridge

import (
	"os"
	"strings"
	"sync/atomic"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// bug1RetireDungeonPartyAliasesEnv 控制「回城重建城镇组队投影之前，先把副本与城镇作用域
// 的旧 op9 记录用 kind3 退役掉」。
//
// 默认（未设置 / on）= 退役；legacy / off / 0 / false / no = 回到改动前。
const bug1RetireDungeonPartyAliasesEnv = "DNFBRIDGE_PARTY_TOWN_RETURN_RETIRE_ALIASES"

func bug1RetireDungeonPartyAliasesEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(bug1RetireDungeonPartyAliasesEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func armDungeonReturnAliasRetirePending(session *gameSession) {
	if session == nil {
		return
	}
	atomic.StoreUint32(&session.dungeonReturnAliasRetirePending, 1)
}

func consumeDungeonReturnAliasRetirePending(session *gameSession) bool {
	if session == nil {
		return false
	}
	return atomic.CompareAndSwapUint32(&session.dungeonReturnAliasRetirePending, 1, 0)
}

// allSameChannelPartyMembersOutsideDungeon 判断当前队伍中所有在线且同频道的成员是否都已经脱离副本运行时。
// 只有当同频道队友全部回到城镇（彼此的 town actor 已通过 broadcastTownPlayerEnter 创建完毕）时，
// 才在 [town_local_actor+0x4C8] == 0 的状态下执行一次 retireDungeonPartyAliasRecords + sendRuntimePartyRosterLocal，
// 避免先回城的成员提前绑定 [town_local_actor+0x4C8] 后再被二次 kind3 触发 0x201786B 队伍解散。
func (s *Service) allSameChannelPartyMembersOutsideDungeon(receiver *gameSession, state alignedcmd.PartyState) bool {
	if s == nil || receiver == nil {
		return false
	}
	members := runtimePartyMembers(state)
	if len(members) <= 1 {
		return true
	}
	receiverChannelID := currentTownPresenceChannelID(receiver)
	for _, member := range members {
		if member.UserID == 0 || member.UserID == receiver.selectedCharacterID {
			continue
		}
		memberSession, online := s.onlineGameSession(member.UserID)
		if !online || memberSession == nil {
			continue
		}
		if currentTownPresenceChannelID(memberSession) != receiverChannelID {
			continue
		}
		memberSession.dungeon.mu.Lock()
		inDungeon := memberSession.dungeon.runtime != nil || memberSession.confirmedDungeonReturnStatePending
		memberSession.dungeon.mu.Unlock()
		if inDungeon {
			return false
		}
	}
	return true
}

// retireDungeonPartyAliasRecords 在回城后首次重建城镇组队投影（sendRuntimePartyRosterLocal）之前，
// 趁 [town_local_actor+0x4C8] 仍为 0（before==0 && after==0，不会触发 0x201786B 队伍解散），
// 把副本作用域（owner = 0）与城镇作用域（owner = channel）的旧 op9 记录逐条用 kind3 静默清槽：
// **先同屏非本队单人角色（如图内已退队者）、再本队队友、最后自己**。
func (s *Service) retireDungeonPartyAliasRecords(receiver *gameSession, members []alignedcmd.PartyMemberState, source string) error {
	if s == nil || receiver == nil || receiver.selectedCharacterID == 0 {
		return nil
	}
	route := s.currentSceneActorRoute()
	memberSet := make(map[uint16]bool, len(members))
	for _, member := range members {
		if member.UserID != 0 {
			memberSet[member.UserID] = true
		}
	}

	ordered := make([]uint16, 0, len(members)+4)
	seen := make(map[uint16]bool, len(members)+4)
	if s.onlinePlayers != nil {
		for _, peer := range s.onlinePlayers.PeersInSameArea(receiver.selectedCharacterID) {
			if peer.CharacterID == 0 || memberSet[peer.CharacterID] || seen[peer.CharacterID] {
				continue
			}
			if peer.Session != nil {
				peerState := s.townCopresenceActorPartyState(peer.Session)
				if peerState.PartyID > 0 {
					continue
				}
			}
			seen[peer.CharacterID] = true
			ordered = append(ordered, peer.CharacterID)
		}
	}
	for _, member := range members {
		if member.UserID == 0 || member.UserID == receiver.selectedCharacterID || seen[member.UserID] {
			continue
		}
		seen[member.UserID] = true
		ordered = append(ordered, member.UserID)
	}
	if !seen[receiver.selectedCharacterID] {
		seen[receiver.selectedCharacterID] = true
		ordered = append(ordered, receiver.selectedCharacterID)
	}

	owners := []byte{currentSceneObjectContext}
	if channel := currentClientSelectedChannel(receiver); channel != currentSceneObjectContext {
		owners = append(owners, channel)
	}

	retire := func(characterID uint16, owner byte) error {
		objectKey := currentSceneActorObjectKey(characterID)
		if objectKey == 0 {
			return nil
		}
		body := buildCurrentSceneOp9ActorRemovalBodyInContext(
			objectKey,
			owner,
			route,
		)
		return s.sendGameUpperRawClass(
			receiver,
			uint16(dnfenum.CmdPacketRecoverStamina),
			body,
			0,
		)
	}

	for _, owner := range owners {
		for _, charID := range ordered {
			if err := retire(charID, owner); err != nil {
				return err
			}
		}
	}

	s.logGameEvent(receiver, "game-party-town-return-dungeon-alias-retired",
		"source", source,
		"char_id", receiver.selectedCharacterID,
		"retired_peers", len(ordered)-1,
		"owner_channel", currentSceneObjectContext,
		"switch", bug1RetireDungeonPartyAliasesEnv,
		"reason", "kind3_owner0_and_channel_peers_then_self_while_town_actor_unlinked")
	return nil
}
