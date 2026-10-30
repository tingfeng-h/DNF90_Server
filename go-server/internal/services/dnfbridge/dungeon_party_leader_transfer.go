package dnfbridge

import (
	"os"
	"sort"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/party"
)

// \u961f\u957f\u9000\u56fe\uff08\u56de\u57ce\uff09\u65f6\u6309 DNF \u5b98\u65b9\u8bed\u4e49\u5904\u7406\uff1a**\u961f\u4f0d\u4e0d\u89e3\u6563**\uff0c\u961f\u957f\u4f4d\u79fb\u4ea4\u7ed9\u5269\u4f59\u6210\u5458\u4e2d
// \u69fd\u4f4d\u6700\u9760\u524d\u7684\u90a3\u4e00\u4f4d\uff0c\u7136\u540e\u65e7\u961f\u957f\u4ee5\u666e\u901a\u6210\u5458\u8eab\u4efd\u79bb\u961f\u3002\u961f\u4f0d ID\u3001P2P \u7aef\u70b9\u4e0d\u53d8\u3002
//
// \u590d\u7528\u624b\u52a8\u59d4\u4efb\u961f\u957f\u90a3\u6761\u5df2\u9a8c\u6536\u7684\u7ebf\u683c\u5f0f\uff08manager.TransferLeader + publish +
// sendChangeHostRuntimePartySnapshots\uff09\uff0c\u907f\u514d\u53e6\u5199\u4e00\u5957\u540d\u518c\u6295\u5f71\u3002
//
// \u4efb\u4f55\u4e00\u6b65\u5931\u8d25\u90fd\u8fd4\u56de ok=false\uff0c\u7531\u8c03\u7528\u65b9\u56de\u9000\u5230\u201c\u89e3\u6563\u201d\u6863\uff0c\u4e0d\u4f1a\u628a\u961f\u4f0d\u5361\u5728\u4e2d\u95f4\u6001\u3002
func (s *Service) transferRuntimePartyOnDungeonLeaderExit(
	session *gameSession,
	instance *dungeonPartyInstance,
	leavingCharacterID uint16,
	source string,
) (int, bool) {
	if s == nil || session == nil || leavingCharacterID == 0 {
		return 0, false
	}
	state := s.runtimePartyStateForSession(session)
	members := runtimePartyMembers(state)
	if state.PartyID <= 0 || len(members) < 2 || state.UserID != leavingCharacterID {
		return 0, false
	}
	// 只有**队长走了之后还剩 ≥2 人**才转让。
	// 2 人局队长退图只剩 1 人，而客户端根本画不出“1 人队伍”：
	// 名册帧会被丢弃、界面按旧名册显示（旧队长那行残留），
	// 直到幸存者自己出图重建场景。这种情形直接走解散（旧行为，已验证干净）。
	if len(members) < 3 {
		return 0, false
	}
	manager := s.runtimePartyManagerForService()
	if manager == nil {
		return 0, false
	}
	// \u65b0\u961f\u957f\uff1a\u6309\u69fd\u4f4d\u987a\u5e8f\uff08runtimePartyMembers \u5df2\u6392\u5e8f\uff09\u53d6\u7b2c\u4e00\u4e2a\u4e0d\u662f\u81ea\u5df1\u7684\u6210\u5458\u3002
	var next alignedcmd.PartyMemberState
	for _, member := range members {
		if member.UserID == 0 || member.UserID == leavingCharacterID {
			continue
		}
		next = member
		break
	}
	if next.UserID == 0 {
		return 0, false
	}
	// \u65b0\u961f\u957f\u7684\u4f1a\u8bdd\u4ee3\u9645\u53ea\u80fd\u4ece\u6743\u5a01\u5feb\u7167\u91cc\u53d6\uff08PartyMemberState \u4e0d\u5e26\u5b57\u6bb5\uff09\u3002
	newLeaderGeneration := uint64(0)
	if snapshot, found := manager.SnapshotByUser(next.UserID, 0); found {
		for _, member := range snapshot.Members {
			if member.UserID == next.UserID {
				newLeaderGeneration = member.SessionGeneration
				break
			}
		}
	}
	if newLeaderGeneration == 0 {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-transfer-skipped",
			"source", source, "party_id", state.PartyID,
			"leaving_char_id", leavingCharacterID, "next_leader_id", next.UserID,
			"reason", "next_leader_has_no_authoritative_session_generation")
		return 0, false
	}
	identity, bound := s.boundGameSessionCharacterSnapshot(session)
	if !bound || !s.boundGameSessionCharacterCurrent(identity) {
		return 0, false
	}
	if result := manager.TransferLeader(leavingCharacterID, identity.generation, next.UserID, newLeaderGeneration); !result.OK {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-transfer-failed",
			"source", source, "party_id", state.PartyID,
			"leaving_char_id", leavingCharacterID, "next_leader_id", next.UserID,
			"reason", "manager_transfer_failed")
		return 0, false
	}
	// \u8f6c\u79fb\u5b8c\u6210\u540e\u65e7\u961f\u957f\u5df2\u662f\u666e\u901a\u6210\u5458\uff0c\u8fd9\u65f6\u624d\u80fd\u5e72\u51c0\u5730\u79bb\u961f\uff08\u4e0d\u4f1a\u89e6\u53d1\u201c\u961f\u957f\u79bb\u961f\u2192\u91cd\u5efa\u4e00\u4ee3\u201d\uff09\u3002
	var nextSnapshot party.RuntimePartySnapshot
	if leaveResult := manager.LeaveExpected(leavingCharacterID, identity.generation, uint16(state.PartyID)); leaveResult.OK {
		nextSnapshot = leaveResult.Party
	} else {
		// \u961f\u957f\u4f4d\u5df2\u7ecf\u8f6c\u51fa\u53bb\u4e86\uff0c\u8fd9\u91cc\u4e0d\u80fd\u518d\u8d70\u89e3\u6563\uff1a\u53ea\u8bb0\u4e8b\u4ef6\u5e76\u6309\u5f53\u524d\u6743\u5a01\u5feb\u7167\u6295\u5f71\u3002
		s.logGameEvent(session, "game-party-dungeon-leader-exit-leave-failed",
			"source", source, "party_id", state.PartyID,
			"leaving_char_id", leavingCharacterID, "reason", leaveResult.Reason)
		snapshot, found := manager.SnapshotByID(uint16(state.PartyID))
		if !found {
			return 0, true
		}
		nextSnapshot = snapshot
	}
	nextState := s.publishRuntimePartySnapshot(nextSnapshot)
	// 客户端的队伍窗口是"从 slot 0 推队长"：转让后必须把新队长挤到 slot 0，
	// 其余成员依次 1..n-1。否则新队长还挂在原来的 slot 1，客户端标记不换人。
	normalized := normalizeRuntimePartySlotsForProjection(nextState)
	s.syncRuntimePartyUDPRelay(nextState)
	// 两段式：客户端在副本里只会“更新已有行”，不会因为名册里
	// 槽位变了就把队长标记换人。先发一次空名册（让它把本地队伍表
	// 拆掉），紧接着发新队长在 slot 0 的完整名册，相当于“重建”。
	if leaderTransferTwoPhaseEnabled() {
		for _, member := range runtimePartyMembers(normalized) {
			memberIdentity, online := s.onlineGameSessionCharacterSnapshot(member.UserID)
			if !online || memberIdentity.session == nil {
				continue
			}
			if _, err := s.sendCurrentPartyActorFrameProjectionInContext(
				memberIdentity.session,
				memberIdentity.session,
				alignedcmd.PartyState{},
				source+"_transfer_phase1_clear",
				nil,
				currentPartyFrameOwnerSelection{Mode: currentPartyFrameOwnerDungeonContext},
			); err != nil {
				s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-phase1-failed",
					"source", source, "char_id", member.UserID, "error", err)
			}
		}
		s.sendDungeonPartyFrameRoster(normalized, source+"_transfer_phase2_rebuild")
		s.logGameEvent(session, "game-party-dungeon-leader-exit-two-phase",
			"source", source, "party_id", normalized.PartyID,
			"member_count", len(runtimePartyMembers(normalized)),
			"new_leader_id", next.UserID,
			"switch", leaderTransferTwoPhaseEnv,
			"reason", "client_only_updates_existing_dungeon_rows_so_the_table_is_cleared_then_rebuilt")
	} else if err := s.sendChangeHostRuntimePartySnapshots(nextState); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-roster-failed",
			"source", source, "party_id", nextState.PartyID, "error", err)
	}
	// 再按归一化后的槽位补一轮副本名册（slot 0 = 新队长）。
	if sent := s.sendDungeonPartyFrameRoster(normalized, "managed_runtime_party_change_host_dungeon_roster_normalized"); sent > 0 {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-normalized-roster",
			"source", source, "party_id", normalized.PartyID, "sent", sent,
			"leader_slot", "zero", "member_count", len(runtimePartyMembers(normalized)))
	}
	// 副本内的幸存者：即使只剩 1 人，也必须用**副本上下文**重发一次自己的名册。
	// sendDungeonPartyFrameRoster 在成员数 < 2 时直接返回，于是图里的客户端还按旧名册画：
	// 旧队长标记不换人、旧队长那一行也残留，直到幸存者自己出图重建场景才继承队长。
	for _, member := range runtimePartyMembers(nextState) {
		memberIdentity, online := s.onlineGameSessionCharacterSnapshot(member.UserID)
		if !online || memberIdentity.session == nil {
			continue
		}
		// 顺序契约（见 sendDungeonPartyFrameRoster 注释）：副本名册必须
		// **op153（血条/槽位实时表）在前，op9（槽位记录）在后**。
		// 反过来发的话客户端会按旧实时表画，表现就是"出图才继承队长"。
		if err := s.sendGameUpperRawClassCodec(memberIdentity.session, 0x0099, runtimePartyRealtimeInfoBody(nextState), 0, true); err != nil {
			s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-dungeon-realtime-failed",
				"source", source, "char_id", member.UserID, "error", err)
		}
		if _, err := s.sendCurrentPartyActorFrameProjectionInContext(
			memberIdentity.session,
			memberIdentity.session,
			normalized,
			"leader_exit_transfer_dungeon_self",
			nil,
			currentPartyFrameOwnerSelection{Mode: currentPartyFrameOwnerDungeonContext},
		); err != nil {
			s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-dungeon-roster-failed",
				"source", source, "char_id", member.UserID, "error", err)
		}
	}

	// \u5e78\u5b58\u8005\uff1a\u65e7\u961f\u957f\u7684 actor \u4e0e\u4e24\u5957\u4e0a\u4e0b\u6587\u7684\u961f\u4f0d\u8bb0\u5f55\u90fd\u8981\u6e05\u6389\uff0c
	// \u5426\u5219\u56fe\u91cc\u4f1a\u7559\u4e0b\u4ed6\u90a3\u4e00\u884c\u8840\u6761\uff1b\u65b0\u961f\u957f\u6807\u8bb0\u7531\u4e0a\u9762\u7684\u540d\u518c\u5e27\u5e26\u8fc7\u53bb\u3002
	retired := []alignedcmd.PartyMemberState{{UserID: leavingCharacterID}}
	cleared := 0
	for _, member := range runtimePartyMembers(nextState) {
		memberIdentity, online := s.onlineGameSessionCharacterSnapshot(member.UserID)
		if !online || memberIdentity.session == nil {
			continue
		}
		if err := s.clearDungeonSurvivorTownPartyRecords(memberIdentity.session, retired, source+"_retired_leader_town"); err != nil {
			s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-retired-town-clear-failed",
				"source", source, "char_id", member.UserID, "error", err)
		}
		if err := s.clearDisbandedDungeonPartyRecords(memberIdentity.session, retired, leavingCharacterID, false, source+"_retired_leader_dungeon"); err != nil {
			s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-retired-dungeon-clear-failed",
				"source", source, "char_id", member.UserID, "error", err)
		}
		if err := s.sendDungeonDepartedActorRemoval(memberIdentity.session, leavingCharacterID, source+"_retired_leader_actor"); err != nil {
			s.logGameEvent(memberIdentity.session, "game-party-dungeon-leader-exit-departed-actor-failed",
				"source", source, "char_id", member.UserID, "error", err)
		}
		cleared++
	}

	// \u65e7\u961f\u957f\u672c\u4eba\uff1a\u57ce\u9547\u4e0a\u4e0b\u6587\u7684\u9500\u6bc1\u5e27 + \u9000\u961f ACK\uff0c\u754c\u9762\u624d\u4f1a\u7acb\u5373\u62c6\u6389\u3002
	if err := s.sendDungeonExitSelfFrameClear(session, source+"_self"); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-self-clear-failed", "source", source, "error", err)
	}
	if err := s.sendDungeonExitPartyLeaveAck(session, source+"_self"); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-self-ack-failed", "source", source, "error", err)
	}
	s.clearRuntimePartyProjection(leavingCharacterID, identity.generation)

	s.logGameEvent(session, "game-party-dungeon-leader-exit-transferred",
		"source", source,
		"party_id", nextState.PartyID,
		"old_party_id", state.PartyID,
		"leaving_char_id", leavingCharacterID,
		"new_leader_id", next.UserID,
		"member_count", len(runtimePartyMembers(nextState)),
		"survivors_cleared", cleared,
		"instance_id", instanceIDOfPartyDungeonInstance(instance),
		"reason", "official_party_keeps_leader_transferred_to_the_first_remaining_slot")
	return cleared, true
}

// normalizeRuntimePartySlotsForProjection 把队伍成员按"队长在 slot 0、其余依次 1..n-1"重排。
// 客户端队伍窗口用 slot 0 推队长，原槽位不连续（如 1/2）时标记不会跟着转移。
func normalizeRuntimePartySlotsForProjection(state alignedcmd.PartyState) alignedcmd.PartyState {
	members := runtimePartyMembers(state)
	if len(members) == 0 {
		return state
	}
	leaderID := state.UserID
	ordered := make([]alignedcmd.PartyMemberState, 0, len(members))
	rest := make([]alignedcmd.PartyMemberState, 0, len(members))
	for _, member := range members {
		if member.UserID == leaderID {
			ordered = append(ordered, member)
			continue
		}
		rest = append(rest, member)
	}
	sort.Slice(rest, func(i, j int) bool { return rest[i].Slot < rest[j].Slot })
	ordered = append(ordered, rest...)
	for i := range ordered {
		ordered[i].Slot = byte(i)
	}
	state.Members = ordered
	return state
}

const leaderTransferTwoPhaseEnv = "DNFBRIDGE_PARTY_LEADER_TRANSFER_TWOPHASE"

// leaderTransferTwoPhaseEnabled 默认 on；显式 off/0/false/no 回退到单段名册。
func leaderTransferTwoPhaseEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(leaderTransferTwoPhaseEnv))) {
	case "off", "0", "false", "no":
		return false
	default:
		return true
	}
}
