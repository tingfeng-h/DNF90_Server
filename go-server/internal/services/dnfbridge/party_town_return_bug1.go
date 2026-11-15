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
		if inDungeon || !currentPartyTownReturnSceneReady(memberSession) {
			return false
		}
		if s.onlinePlayers != nil {
			peerInfo, found := s.onlinePlayers.PlayerForCharacter(member.UserID)
			if !found || peerInfo.Session != memberSession {
				return false
			}
		}
	}
	return true
}

// sendRuntimePartyTownReturnRebind 在回城同屏 actor 重建完成后，按顺序执行：
//  1. retireDungeonPartyAliasRecords（若本轮回城待决）：清理 owner=0 与非本机别名的旧 CPartyRecord；
//  2. sendRuntimePartyPeerEndpointInfo（op11）：将对端 IP 与端口写进新创建的城镇 CUser(+0x4B4..+0x4BC)；
//  3. sendRuntimePartyRosterLocal（op153 + op9）：原位绑定 [town_actor+0x4C8] 并触发有效端点 P2P 重连。
func (s *Service) sendRuntimePartyTownReturnRebind(
	session *gameSession,
	state alignedcmd.PartyState,
	source string,
) error {
	if s == nil || session == nil || session.selectedCharacterID == 0 {
		return nil
	}
	receiverState := stateForRuntimePartyReceiver(state, session.selectedCharacterID)
	aliasRetireConsumed := false
	if bug1RetireDungeonPartyAliasesEnabled() && consumeDungeonReturnAliasRetirePending(session) {
		aliasRetireConsumed = true
		if err := s.retireDungeonPartyAliasRecords(session, runtimePartyMembers(receiverState), source+"_retire_aliases"); err != nil {
			return err
		}
	}
	if aliasRetireConsumed && receiverState.PartyID > 0 && len(runtimePartyMembers(receiverState)) >= 2 {
		if _, err := s.sendRuntimePartyPeerEndpointInfo(session, receiverState, s.partyFormationOrder()); err != nil {
			return err
		}
	}
	return s.sendRuntimePartyRosterLocal(session, receiverState, source)
}

// retireDungeonPartyAliasRecords 在回城后首次重建城镇组队投影（sendRuntimePartyRosterLocal）之前，
// 趁 [town_local_actor+0x4C8] 仍为 0（before==0 && after==0，不会触发 0x201786B 队伍解散），
// 把副本作用域（owner = 0）与城镇作用域（owner = channel）的非本机旧 op9 记录逐条用 kind3 静默清槽。
func (s *Service) retireDungeonPartyAliasRecords(receiver *gameSession, members []alignedcmd.PartyMemberState, source string) error {
	if s == nil || receiver == nil || receiver.selectedCharacterID == 0 {
		return nil
	}
	route := s.currentSceneActorRoute()
	memberSet := make(map[uint16]bool, len(members)+1)
	memberSet[receiver.selectedCharacterID] = true
	for _, member := range members {
		if member.UserID != 0 {
			memberSet[member.UserID] = true
		}
	}

	nonPartyIDs := make([]uint16, 0, 4)
	peerIDs := make([]uint16, 0, len(members))
	ordered := make([]uint16, 0, len(members)+4)
	seen := make(map[uint16]bool, len(members)+4)
	if s.onlinePlayers != nil {
		for _, peer := range s.onlinePlayers.PeersInSameArea(receiver.selectedCharacterID) {
			if peer.CharacterID == 0 || memberSet[peer.CharacterID] || seen[peer.CharacterID] {
				continue
			}
			seen[peer.CharacterID] = true
			nonPartyIDs = append(nonPartyIDs, peer.CharacterID)
			ordered = append(ordered, peer.CharacterID)
		}
		receiverChannelID := currentTownPresenceChannelID(receiver)
		s.onlinePlayers.mu.RLock()
		var sameChannelCandidates []onlinePlayerInfo
		for area, bucket := range s.onlinePlayers.byArea {
			if area.ChannelID != receiverChannelID {
				continue
			}
			for _, info := range bucket {
				if info != nil && info.CharacterID != 0 && !memberSet[info.CharacterID] && !seen[info.CharacterID] {
					sameChannelCandidates = append(sameChannelCandidates, *info)
				}
			}
		}
		s.onlinePlayers.mu.RUnlock()
		for _, cand := range sameChannelCandidates {
			if seen[cand.CharacterID] {
				continue
			}
			seen[cand.CharacterID] = true
			nonPartyIDs = append(nonPartyIDs, cand.CharacterID)
			ordered = append(ordered, cand.CharacterID)
		}
	}
	for _, member := range members {
		if member.UserID == 0 || member.UserID == receiver.selectedCharacterID || seen[member.UserID] {
			continue
		}
		seen[member.UserID] = true
		peerIDs = append(peerIDs, member.UserID)
		ordered = append(ordered, member.UserID)
	}
	if !seen[receiver.selectedCharacterID] {
		seen[receiver.selectedCharacterID] = true
		ordered = append(ordered, receiver.selectedCharacterID)
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

	// 1. 副本作用域（owner = 0）：仅退役已退队散人（nonPartyIDs）与本机（self），
	// 不对当前仍在队内的队友（peerIDs）发送 kind3(peerID, 0)，避免穿透误伤同区队友模型。
	for _, id := range nonPartyIDs {
		if err := retire(id, currentSceneObjectContext); err != nil {
			return err
		}
	}
	if err := retire(receiver.selectedCharacterID, currentSceneObjectContext); err != nil {
		return err
	}

	// 2. 城镇作用域（owner = channel）：
	// 若当前仍在多人队伍中，绝不退役 record(self, channel) 且绝不使用 handoffPeer 借槽（避免触发
	// 0x0201786B 静默解散与清零 [actor+0x4C8]），仅清理 nonPartyIDs 与 peerIDs 的别名记录；
	// 随后由 sendRuntimePartyRosterLocal 的 kind0 原位更新 record(self, channel) 并通过
	// .prsafe (0x06EC5334..0x06EC536B) 把本机与所有队友的新 town actor 绑回 record(self, channel)。
	if channel := currentClientSelectedChannel(receiver); channel != currentSceneObjectContext {
		state := stateForRuntimePartyReceiver(s.townCopresenceActorPartyState(receiver), receiver.selectedCharacterID)
		for _, id := range nonPartyIDs {
			if err := retire(id, channel); err != nil {
				return err
			}
		}
		if state.PartyID > 0 && len(runtimePartyMembers(state)) >= 2 {
			for _, id := range peerIDs {
				if err := retire(id, channel); err != nil {
					return err
				}
			}
		} else {
			if err := retire(receiver.selectedCharacterID, channel); err != nil {
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
		"reason", "kind3_owner0_and_channel_peers_preserving_self_channel_record")
	return nil
}
