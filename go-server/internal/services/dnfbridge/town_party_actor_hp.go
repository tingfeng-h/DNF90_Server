package dnfbridge

import (
	"encoding/binary"
	"fmt"
)

// 城镇队友投影此前只发 mode0/mode1 + op9 队伍帧，**不发 op153**。
// 当前 EXE 读 mode0 时会先把 HP 描述符清零（descriptor+31C），之后要靠 op153 恢复；
// 城镇同屏这条链少了这一步，于是队伍框里那名队友的血条是 0/旧值。
//
// 这里只做两件事，且只在 official-24.9 的**普通城镇**队伍里生效：
//  1. 把同一条 op153 里的 HP 百分比写进 mode0 的对应字节（初始化）；
//  2. 在 mode0/mode1 之后补发 op153（只更新 HP，不重建 op9、不重启 op11 链路）。
//
// 副本内、跨频道、超时空频道一律不碰：那些路径有自己的生命/血量通道。

func (s *Service) sendCurrentTownPartyHPAfterActor(target, actor *gameSession) error {
	body := s.currentTownPartyRealtimeBody(target, actor)
	if len(body) == 0 {
		return nil
	}
	return s.sendGameUpperRawClassCodec(target, 0x0099, body, 0, true)
}

func (s *Service) currentTownPartyRealtimeBody(target, actor *gameSession) []byte {
	if !s.currentTownPartyHPProjectionEligible(target, actor) ||
		sessionHasActiveDungeonRuntime(target) || sessionHasActiveDungeonRuntime(actor) {
		return nil
	}
	return s.currentTownPartyRealtimeBodyOutsideDungeon(target, actor)
}

func (s *Service) currentTownPartyHPProjectionEligible(target, actor *gameSession) bool {
	return s != nil && s.options.currentClientVersion == "official-24.9" &&
		target != nil && actor != nil && target != actor &&
		target.selectedCharacterID != actor.selectedCharacterID &&
		!isCurrentFiendWarChannel(target) && !isCurrentFiendWarChannel(actor)
}

// 只取"两侧都在同一城镇频道、且权威队伍名单里确实有这名队友"的 op153 行。
// 缓存里的旧名单会在退队后残留，因此以中央管理器为准。
func (s *Service) currentTownPartyRealtimeBodyOutsideDungeon(target, actor *gameSession) []byte {
	channel := currentTownPresenceChannelID(target)
	if channel == 0 || channel != currentTownPresenceChannelID(actor) {
		return nil
	}
	targetBound, ok := s.boundGameSessionCharacterSnapshot(target)
	if !ok {
		return nil
	}
	actorBound, ok := s.boundGameSessionCharacterSnapshot(actor)
	if !ok || actorBound.character == targetBound.character {
		return nil
	}
	state := s.runtimePartyStateForSession(target)
	members := runtimePartyMembers(state)
	if state.PartyID <= 0 || len(members) <= 1 ||
		!containsRuntimePartyMember(members, actorBound.character) {
		return nil
	}
	if !s.boundGameSessionCharacterCurrent(targetBound) || !s.boundGameSessionCharacterCurrent(actorBound) {
		return nil
	}
	return runtimePartyRealtimeInfoBody(state)
}

// 只改这名 actor 的 mode0 HP 字节，其余（名字、装备、克隆链接、尾块）逐字节保持。
// 只接受与 op153 完全一致的身份：count 在 2..4、行宽 5、actor 恰好出现一次、百分比 ≤100。
func projectCurrentTownPartyHPIntoMode0(mode0, realtime []byte, actorID uint16) error {
	if len(realtime) == 0 {
		return nil
	}
	if actorID == 0 || actorID == ^uint16(0) || realtime[0] < 2 || realtime[0] > 4 ||
		len(realtime) != 1+int(realtime[0])*5 {
		return fmt.Errorf("invalid town realtime HP projection")
	}
	hp, matches := byte(0), 0
	for offset := 1; offset < len(realtime); offset += 5 {
		if binary.LittleEndian.Uint16(realtime[offset:offset+2]) == actorID {
			hp = realtime[offset+2]
			matches++
		}
	}
	if matches != 1 || hp > 100 {
		return fmt.Errorf("town HP projection actor identity or percentage invalid")
	}
	const keyOffset = 5 + 0x47
	if len(mode0) < keyOffset+2 || mode0[0] != 0 ||
		binary.LittleEndian.Uint16(mode0[1:3]) != 1 ||
		binary.LittleEndian.Uint16(mode0[keyOffset:keyOffset+2]) != currentSceneActorObjectKey(actorID) {
		return fmt.Errorf("town mode0 HP projection identity differs")
	}
	_, tailStart, ok := currentSceneObjectLevelForLog(mode0)
	if !ok || tailStart >= len(mode0) {
		return fmt.Errorf("town mode0 HP projection tail missing")
	}
	tail := mode0[tailStart:]
	equipEnd, ok := currentSceneObjectEquipSummaryEnd(tail, 6)
	if !ok || equipEnd+4 >= len(tail) {
		return fmt.Errorf("town mode0 HP projection equipment boundary invalid")
	}
	tail[equipEnd+4] = hp
	return nil
}
