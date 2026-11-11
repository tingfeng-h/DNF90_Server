package dnfbridge

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// bug1RetireDungeonPartyAliasesEnv 控制「回城重建城镇组队投影之前，先把副本与城镇作用域
// 的旧别名 op9 记录用 kind3 退役掉」。
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

// consumeDungeonPartyAliasRetirePending 原子检查并消费"本角色有一笔来自组队副本的
// 别名记录待退役"。同一场副本回城只返回一次 true；选图页 ESC 回城、普通城镇移动、
// 以及未进过副本的会话恒返回 false。
func (s *Service) consumeDungeonPartyAliasRetirePending(session *gameSession) bool {
	if s == nil || session == nil || session.selectedCharacterID == 0 {
		return false
	}
	if s.sessionHasActiveDungeonRuntime(session) {
		return false
	}
	session.dungeonPartyRunMu.Lock()
	defer session.dungeonPartyRunMu.Unlock()
	if session.dungeonPartyAliasRetireCharID == 0 {
		return false
	}
	if session.dungeonPartyAliasRetireCharID != session.selectedCharacterID {
		session.dungeonPartyAliasRetireCharID = 0
		return false
	}
	session.dungeonPartyAliasRetireCharID = 0
	return true
}

// retireDungeonPartyAliasRecordsIfPending 仅在 receiver 确实刚从组队副本回城、
// 且本轮副本的别名尚未退役时执行一次 retireDungeonPartyAliasRecords。
func (s *Service) retireDungeonPartyAliasRecordsIfPending(
	receiver *gameSession,
	members []alignedcmd.PartyMemberState,
	source string,
) (bool, error) {
	if !bug1RetireDungeonPartyAliasesEnabled() {
		return false, nil
	}
	if !s.consumeDungeonPartyAliasRetirePending(receiver) {
		return false, nil
	}
	if err := s.retireDungeonPartyAliasRecords(receiver, members, source); err != nil {
		return false, err
	}
	return true, nil
}

// retireDungeonPartyAliasRecords 在真正从副本回城后、重建城镇组队投影之前：
//  1. 对副本作用域（owner = 0）：用 kind3 退役同屏非本队单人角色、本队队友、以及自己（self）；
//  2. 对城镇作用域（owner = channel）：只用 kind3 退役同屏非本队单人角色（如图内已主动退队者
//     的 retainedRecordID 记录）与本队队友的别名记录，**绝不退役 (self, owner = channel)**。
//
// 为什么绝不能退役 (self, owner = channel)：
// 图内 3->2 时 sendMemberGiveupSurvivorRosterInContext 已经按「先队友、最后 self」发过
// owner=channel，因此 [local_actor+0x4C8] 此时正绑定在 record(self, channel) 上。
// 如果对 (self, channel) 发 kind3，客户端 NOTI9 会命中 [ebp-0x356]==1 && in_party==0，
// 走 0x201786B 清场块把 [local_actor+0x4C8] 置 0；紧接着收到 sendRuntimePartyRosterLocal
// 的 op9 kind0(self, channel) 时就会因 before==0 && after==1 误入 0x2017D7C 首次入队分支，
// 调用 .gld127 的 ResetAllPeerConnections 抹掉全部 UDP 槽位，导致队员血槽出现"连接中"。
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

	// 1. 副本作用域（owner = 0）：非本队单人 -> 队友 -> 自己（self）全部退役。
	for _, charID := range ordered {
		if err := retire(charID, currentSceneObjectContext); err != nil {
			return err
		}
	}

	// 2. 城镇作用域（owner = channel）：只退役非本队单人（如图内退队者）与队友别名，
	// 保留 (self, owner = channel) 不被拆毁，避免随后的 kind0(self, channel) 触发首次入队 UDP 重置。
	if channel := currentClientSelectedChannel(receiver); channel != currentSceneObjectContext {
		for _, charID := range ordered {
			if charID == receiver.selectedCharacterID {
				continue
			}
			if err := retire(charID, channel); err != nil {
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
		"reason", "kind3_owner0_all_and_channel_peers_only_before_town_self_record")
	return nil
}

func buildPartyTownReturnRebuildSignature(state alignedcmd.PartyState, sessions []*gameSession) string {
	type memberSig struct {
		charID       uint16
		charGen      uint64
		placementEp  uint64
		postTransGen uint64
		townSceneGen uint64
		channelID    int
		townID       byte
		areaID       byte
	}
	rows := make([]memberSig, 0, len(sessions))
	for _, m := range sessions {
		if m == nil {
			continue
		}
		m.townPostTransition.mu.Lock()
		postGen := m.townPostTransition.generation
		m.townPostTransition.mu.Unlock()

		m.townMu.Lock()
		sceneGen := m.townSceneGeneration
		townID := m.townPositionSnapshot.TownID
		areaID := m.townPositionSnapshot.AreaID
		m.townMu.Unlock()

		rows = append(rows, memberSig{
			charID:       m.selectedCharacterID,
			charGen:      m.characterGeneration,
			placementEp:  m.dungeonPlacementEpoch.Load(),
			postTransGen: postGen,
			townSceneGen: sceneGen,
			channelID:    currentTownPresenceChannelID(m),
			townID:       townID,
			areaID:       areaID,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].charID < rows[j].charID
	})
	var b strings.Builder
	fmt.Fprintf(&b, "p=%d/l=%d/n=%d", state.PartyID, state.UserID, len(rows))
	for _, r := range rows {
		fmt.Fprintf(&b, ";c=%d:%d:%d:%d:%d:%d:%d:%d",
			r.charID, r.charGen, r.placementEp, r.postTransGen,
			r.townSceneGen, r.channelID, r.townID, r.areaID)
	}
	return b.String()
}

func claimPartyTownReturnRebuild(state alignedcmd.PartyState, sessions []*gameSession) (string, bool) {
	if len(sessions) == 0 {
		return "", false
	}
	sig := buildPartyTownReturnRebuildSignature(state, sessions)
	if sig == "" {
		return "", false
	}
	allMatched := true
	for _, m := range sessions {
		if m == nil {
			return "", false
		}
		m.partyTownReturnRebuildMu.Lock()
		matched := m.lastPartyTownReturnRebuildSig == sig
		m.partyTownReturnRebuildMu.Unlock()
		if !matched {
			allMatched = false
		}
	}
	if allMatched {
		return sig, false
	}
	for _, m := range sessions {
		m.partyTownReturnRebuildMu.Lock()
		m.lastPartyTownReturnRebuildSig = sig
		m.partyTownReturnRebuildMu.Unlock()
	}
	return sig, true
}

func releasePartyTownReturnRebuild(sessions []*gameSession, sig string) {
	if sig == "" {
		return
	}
	for _, m := range sessions {
		if m == nil {
			continue
		}
		m.partyTownReturnRebuildMu.Lock()
		if m.lastPartyTownReturnRebuildSig == sig {
			m.lastPartyTownReturnRebuildSig = ""
		}
		m.partyTownReturnRebuildMu.Unlock()
	}
}
