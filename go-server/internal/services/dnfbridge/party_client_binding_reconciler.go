package dnfbridge

// 本文件是"客户端组队绑定协调器"：只负责 DNF.exe 侧 CPartyRecord 绑定([localPlayer+0x4C8])
// 的跟踪与修复，不承担任何业务判断（谁在队里、谁是队长由 RuntimePartyManager 决定）。
//
// 依据（客户端逆向，见 Docs 之外的分析记录）：
//   - [actor+0x4C8] 全程序只有 7 处写入；其中 record(self,ch) 想重新绑定自己，
//     必须满足 CPartyRecord::SetMemberUser(0x029733B0) 里 oldUser != newUser；
//   - 若某条非本机的城镇记录 record(D,ch) 先把 [self+0x4C8] 抢绑走，
//     之后所有 record(self,ch) 的 kind0/kind2 都会因 oldUser==newUser 变成空操作。
// 因此需要在"被抢绑过"的会话上先发一条空 kind2 把 record(self,ch) 的槽位复位成 0xFFFF，
// 再发权威 kind0，才能把绑定夺回；随后才能安全用 kind3 清掉退图者的残留记录。

import (
	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// partyClientBindingState 是单个会话的客户端绑定投影账本（受 session.party.mu 保护）。
type partyClientBindingState struct {
	// foreignChannelRecordBound 表示 [localPlayer+0x4C8] 曾被"非本机、owner==频道号"的记录抢绑。
	foreignChannelRecordBound bool
	// retainedAliasKeys 是当前组队周期内向该客户端投影过的非本机角色 ID（如中途退图者）。
	// 只在下发 kind3 退役或整队离场时消费。
	retainedAliasKeys map[uint16]struct{}
}

// recordClientPartyProjectionObserved 在向 receiver 实际发出非空 NOTI 9(kind0/kind2) 后登记。
func (s *Service) recordClientPartyProjectionObserved(
	receiver *gameSession,
	recordCharacterID uint16,
	ownerChannel byte,
	localChannelID int,
) {
	if s == nil || receiver == nil || recordCharacterID == 0 {
		return
	}
	selfID := receiver.selectedCharacterID
	if selfID == 0 || recordCharacterID == selfID {
		return
	}
	receiver.party.mu.Lock()
	if receiver.party.clientBinding.retainedAliasKeys == nil {
		receiver.party.clientBinding.retainedAliasKeys = make(map[uint16]struct{}, 4)
	}
	receiver.party.clientBinding.retainedAliasKeys[recordCharacterID] = struct{}{}
	if ownerChannel != currentSceneObjectContext && int(ownerChannel) == localChannelID {
		receiver.party.clientBinding.foreignChannelRecordBound = true
	}
	receiver.party.mu.Unlock()
}

// recordClientPartyAliasKeys 批量登记该客户端接触过的非本机角色记录键。
func (s *Service) recordClientPartyAliasKeys(receiver *gameSession, characterIDs []uint16, markForeignChannelHijacked bool) {
	if s == nil || receiver == nil || len(characterIDs) == 0 {
		return
	}
	selfID := receiver.selectedCharacterID
	receiver.party.mu.Lock()
	if receiver.party.clientBinding.retainedAliasKeys == nil {
		receiver.party.clientBinding.retainedAliasKeys = make(map[uint16]struct{}, len(characterIDs))
	}
	for _, id := range characterIDs {
		if id != 0 && id != selfID {
			receiver.party.clientBinding.retainedAliasKeys[id] = struct{}{}
		}
	}
	if markForeignChannelHijacked {
		receiver.party.clientBinding.foreignChannelRecordBound = true
	}
	receiver.party.mu.Unlock()
}

// takeForeignChannelRecordHijacked 原子检查并消费"曾被抢绑"标志。
func takeForeignChannelRecordHijacked(receiver *gameSession) bool {
	if receiver == nil {
		return false
	}
	receiver.party.mu.Lock()
	hijacked := receiver.party.clientBinding.foreignChannelRecordBound
	receiver.party.clientBinding.foreignChannelRecordBound = false
	receiver.party.mu.Unlock()
	return hijacked
}

// consumeNonMemberClientPartyAliasKeys 取出并消费所有已不在当前队伍里的历史别名键。
func consumeNonMemberClientPartyAliasKeys(receiver *gameSession, activeMembers []alignedcmd.PartyMemberState) []uint16 {
	if receiver == nil {
		return nil
	}
	selfID := receiver.selectedCharacterID
	memberSet := make(map[uint16]bool, len(activeMembers)+1)
	if selfID != 0 {
		memberSet[selfID] = true
	}
	for _, m := range activeMembers {
		if m.UserID != 0 {
			memberSet[m.UserID] = true
		}
	}
	receiver.party.mu.Lock()
	defer receiver.party.mu.Unlock()
	if len(receiver.party.clientBinding.retainedAliasKeys) == 0 {
		return nil
	}
	var retired []uint16
	for id := range receiver.party.clientBinding.retainedAliasKeys {
		if id == 0 || id == selfID {
			delete(receiver.party.clientBinding.retainedAliasKeys, id)
			continue
		}
		if !memberSet[id] {
			retired = append(retired, id)
			delete(receiver.party.clientBinding.retainedAliasKeys, id)
		}
	}
	return retired
}

// consumeAllClientPartyAliasKeysOnDeparture 在本机彻底离队/被踢时取出全部历史别名键并复位。
func consumeAllClientPartyAliasKeysOnDeparture(receiver *gameSession) []uint16 {
	if receiver == nil {
		return nil
	}
	selfID := receiver.selectedCharacterID
	receiver.party.mu.Lock()
	defer receiver.party.mu.Unlock()
	receiver.party.clientBinding.foreignChannelRecordBound = false
	if len(receiver.party.clientBinding.retainedAliasKeys) == 0 {
		return nil
	}
	out := make([]uint16, 0, len(receiver.party.clientBinding.retainedAliasKeys))
	for id := range receiver.party.clientBinding.retainedAliasKeys {
		if id != 0 && id != selfID {
			out = append(out, id)
		}
	}
	receiver.party.clientBinding.retainedAliasKeys = nil
	return out
}

// resetUnboundSelfChannelRecordIfHijacked 在绑定被抢时，先给 record(self,ch) 发空 kind2 复位槽位。
// 此时 [self+0x4C8] != record(self,ch)，客户端走 wasBound==0 && isNowBound==0 静默分支：
// 不关 P2P、不清血条、不弹提示，只把 memberUsers[0..7] 复位成 0xFFFF。
func (s *Service) resetUnboundSelfChannelRecordIfHijacked(
	receiver *gameSession,
	state alignedcmd.PartyState,
	source string,
) error {
	if s == nil || receiver == nil || receiver.selectedCharacterID == 0 || state.PartyID <= 0 {
		return nil
	}
	if s.options.currentClientVersion != "official-24.9" {
		return nil
	}
	if !takeForeignChannelRecordHijacked(receiver) {
		return nil
	}
	channel := currentClientSelectedChannel(receiver)
	if channel == currentSceneObjectContext {
		return nil
	}
	body := s.buildCurrentPartyDepartureBody(receiver.selectedCharacterID, channel)
	if _, err := s.sendRetiredPartyPacket(
		receiver,
		uint16(dnfenum.CmdPacketRecoverStamina),
		body,
		&state,
		nil,
		false,
	); err != nil {
		return err
	}
	s.logGameEvent(receiver, "game-party-client-self-record-pre-rebind-reset",
		"source", source,
		"char_id", receiver.selectedCharacterID,
		"owner_channel", channel,
		"party_id", state.PartyID,
		"reason", "reset_unbound_self_record_slots_so_setmemberuser_rebinds_local_actor_from_departed_alias")
	return nil
}

// retireDepartedPartyAliasRecordsAfterSelfRebind 在绑定已回到 record(self,ch) 之后，
// 用 kind3 静默退役已不在队伍中的历史记录（如中途退图者 D）的两套 owner。
func (s *Service) retireDepartedPartyAliasRecordsAfterSelfRebind(
	receiver *gameSession,
	state alignedcmd.PartyState,
	source string,
) error {
	if s == nil || receiver == nil || receiver.selectedCharacterID == 0 {
		return nil
	}
	if s.options.currentClientVersion != "official-24.9" {
		return nil
	}
	members := runtimePartyMembers(state)
	retiredIDs := consumeNonMemberClientPartyAliasKeys(receiver, members)
	if len(retiredIDs) == 0 {
		return nil
	}
	route := s.currentSceneActorRoute()
	owners := partyRecordOwnerChain(receiver, false)
	for _, owner := range owners {
		for _, id := range retiredIDs {
			objectKey := currentSceneActorObjectKey(id)
			if objectKey == 0 {
				continue
			}
			body := buildCurrentSceneOp9ActorRemovalBodyInContext(objectKey, owner, route)
			if _, err := s.sendRetiredPartyPacket(
				receiver,
				uint16(dnfenum.CmdPacketRecoverStamina),
				body,
				&state,
				nil,
				false,
			); err != nil {
				return err
			}
		}
	}
	s.logGameEvent(receiver, "game-party-client-departed-alias-records-retired",
		"source", source,
		"char_id", receiver.selectedCharacterID,
		"retired_alias_ids", retiredIDs,
		"party_id", state.PartyID,
		"reason", "retire_mid_dungeon_departed_member_records_after_self_channel_record_is_rebound")
	return nil
}
