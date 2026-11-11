package dnfbridge

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// bug1RetireDungeonPartyAliasesEnv 控制「回城重建城镇组队投影之前，先把副本作用域
// 的别名 op9 记录用 kind3 退役掉」。
//
// 默认（未设置 / on）= 退役；legacy / off / 0 / false / no = 回到改动前，只发
// 城镇 R_self、不退役别名。
const bug1RetireDungeonPartyAliasesEnv = "DNFBRIDGE_PARTY_TOWN_RETURN_RETIRE_ALIASES"

func bug1RetireDungeonPartyAliasesEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(bug1RetireDungeonPartyAliasesEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// armDungeonPartyAliasRetirePending 在会话进入组队副本或收到副本作用域（owner=0）
// 队伍帧时置位，表示该角色离开副本回城时需要退役一次 owner=0 别名记录。
func armDungeonPartyAliasRetirePending(session *gameSession) {
	if session == nil || session.selectedCharacterID == 0 {
		return
	}
	session.dungeonPartyRunMu.Lock()
	session.dungeonPartyAliasRetireCharID = session.selectedCharacterID
	session.dungeonPartyRunMu.Unlock()
}

// consumeDungeonPartyAliasRetirePending 原子检查并消费"本角色有一笔来自副本的
// owner=0 别名记录待退役"。同一场副本回城只返回一次 true；普通城镇移动、选图页
// 进出、以及未进过副本的会话恒返回 false。
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
// 且本轮副本的 owner=0 别名尚未退役时执行一次 retireDungeonPartyAliasRecords。
// 返回 (retired, err)：retired==true 表示本次发出了 kind3(owner=0) 清场帧。
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

// retireDungeonPartyAliasRecords 把副本作用域（owner = 0）的别名 op9 记录逐条用
// kind3 退役：**先队友，最后自己**。
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

// buildPartyTownReturnRebuildSignature 生成全队当前城镇落位世代的唯一幂等签名。
// 任何成员再次进出副本/选图页（dungeonPlacementEpoch 自增）、触发新的城镇过场世代
// （townPostTransition.generation 自增）、跨城镇区域（townSceneGeneration / TownID / AreaID 变化）、
// 换频道、换角（characterGeneration 变化）或队伍成员/队长变化，签名都会改变。
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

// claimPartyTownReturnRebuild 判断全队是否已经对该城镇落位世代执行过同屏重建。
// 若所有成员都已记录相同签名，返回 false（幂等跳过）；否则更新全队成员的签名并返回 true。
// 按会话指针排序依次单独加锁（或逐个更新），不嵌套任何其它锁。
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

// releasePartyTownReturnRebuild 在重建中途写 socket 失败时回滚签名，允许下一次确认重试。
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
