package dnfbridge

import (
	"os"
	"strings"
)

// 影片阶段的"全员制"：默认让影片自然放完；只有**全部实际成员**都按了取消
// （C2S 659 body=0，厂商 018FBF20 "skip input while the movie is active"）
// 才回关闭提示，让所有人的影片一起结束。
//
// 另一半：阶段推进不再硬等一个影片时长。原来 raid_reward_flow.go 里是
// `SquadShownAt + movie`（movie = PVF [MOVIE END TIME]），所以"影片放完还要卡一会"
// 是服务端自己造成的。开启本档后，只要**全部实际成员**都发过 659 body=1
// （厂商 018FBED0 "close completion"，只发一次）就立刻 FinishedAt → 推牌。
//
// "全员"= 攻坚队实际成员（`state.Members`，随加入动态 append），
// 不是团本容量：卢克 8 人团只有 2 人参与时，2 人都完成即可。
//
// 开关：DNFBRIDGE_RAID_MOVIE_ALL_CANCEL
//
//	默认（空/未设置/非法值）= on（本次需求的行为）；
//	off / 0 / false / no / legacy = 逐字节回到改动前（卢克二阶段单发提示、
//	阶段推进等满一个影片时长）。
const raidMovieAllCancelEnv = "DNFBRIDGE_RAID_MOVIE_ALL_CANCEL"

func raidMovieAllCancelEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(raidMovieAllCancelEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// raidMovieVotingMembersLocked 返回本次攻坚里需要投票的实际成员。
// 调用方已持 s.raidMu。这里**不做在线过滤**：掉线成员挡住了提前终止，
// 结果只是影片自然放完（正是默认行为），不会卡死；阶段推进另有超时兜底。
func raidMovieVotingMembersLocked(live *runtimeRaidState) []uint16 {
	if live == nil {
		return nil
	}
	members := make([]uint16, 0, len(live.Members))
	for _, member := range live.Members {
		if member.CharID == 0 {
			continue
		}
		members = append(members, member.CharID)
	}
	return members
}

// raidMovieAllMembersCompletedLocked 报告本次阶段里**全部实际成员**是否都已完成影片
// （即都发过 659 body=1）。无人发过时返回 false，于是回到原来的时长兜底。
// 调用方已持 s.raidMu。
func raidMovieAllMembersCompletedLocked(live *runtimeRaidState, phase byte) bool {
	if live == nil || live.Battle == nil {
		return false
	}
	index := int(phase)
	if index < 0 || index >= len(live.Battle.RewardCardReady) {
		return false
	}
	ready := live.Battle.RewardCardReady[index]
	if len(ready) == 0 {
		return false
	}
	count := 0
	for _, member := range live.Members {
		if member.CharID == 0 {
			continue
		}
		count++
		if !ready[member.CharID] {
			return false
		}
	}
	return count > 0
}

// recordRaidMovieCancelVote 处理 C2S 659 body=0：记一票。
// **没凑齐全员时一个字节都不回**（客户端影片继续播）；凑齐后才逐个回关闭提示。
func (s *Service) recordRaidMovieCancelVote(session *gameSession, captured runtimeRaidState, op uint16) error {
	s.raidMu.Lock()
	live := s.runtimeRaidByMemberLocked(session.selectedCharacterID)
	if live == nil || live.Battle == nil || captured.Battle == nil ||
		live.Battle.InstanceID != captured.Battle.InstanceID || live.Battle.Phase != captured.Battle.Phase {
		s.raidMu.Unlock()
		return s.sendRaidCommandFailure(session, op)
	}
	phase := int(live.Battle.Phase)
	if phase < 0 || phase >= len(live.Battle.MovieCancelVote) {
		s.raidMu.Unlock()
		return s.sendRaidCommandFailure(session, op)
	}
	if live.Battle.MovieCancelVote[phase] == nil {
		live.Battle.MovieCancelVote[phase] = make(map[uint16]bool, len(live.Members))
	}
	live.Battle.MovieCancelVote[phase][session.selectedCharacterID] = true
	votes := len(live.Battle.MovieCancelVote[phase])
	members := raidMovieVotingMembersLocked(live)
	allVoted := len(members) > 0
	for _, id := range members {
		if !live.Battle.MovieCancelVote[phase][id] {
			allVoted = false
			break
		}
	}
	instanceID, battlePhase := live.Battle.InstanceID, live.Battle.Phase
	s.raidMu.Unlock()

	s.logGameEvent(session, "game-raid-movie-cancel-vote",
		"char_id", session.selectedCharacterID,
		"instance_id", instanceID,
		"phase", battlePhase,
		"votes", votes,
		"members", len(members),
		"all_voted", allVoted,
		"env", raidMovieAllCancelEnv,
		"mode", "all_cancel")
	if !allVoted {
		return nil
	}
	for _, id := range members {
		peer, online := s.onlineGameSession(id)
		if !online || peer == nil {
			continue
		}
		peerState, err := s.currentAntonRewardState(peer)
		if err != nil {
			continue
		}
		if err := s.sendCurrentRaidMovieSkipPrompt(peer, peerState, op); err != nil {
			s.logGameEvent(session, "game-raid-movie-cancel-release-failed",
				"char_id", id, "error", err)
		}
	}
	return nil
}

// raidMovieAllReadyLocked 报告本次副本/阶段里是否所有实际成员都已确认影片结束（659/1）。
func (s *Service) raidMovieAllReadyLocked(session *gameSession, captured runtimeRaidState) bool {
	s.raidMu.Lock()
	defer s.raidMu.Unlock()
	live := s.runtimeRaidByMemberLocked(session.selectedCharacterID)
	if live == nil || live.Battle == nil || captured.Battle == nil ||
		live.Battle.InstanceID != captured.Battle.InstanceID || live.Battle.Phase != captured.Battle.Phase {
		return false
	}
	return raidMovieAllMembersCompletedLocked(live, live.Battle.Phase)
}

// claimRaidMovieCompletionDeliveryLocked 保证"全员就绪后的统一投递"每个阶段只发生一次。
func (s *Service) claimRaidMovieCompletionDeliveryLocked(session *gameSession, captured runtimeRaidState) bool {
	s.raidMu.Lock()
	defer s.raidMu.Unlock()
	live := s.runtimeRaidByMemberLocked(session.selectedCharacterID)
	if live == nil || live.Battle == nil || captured.Battle == nil ||
		live.Battle.InstanceID != captured.Battle.InstanceID || live.Battle.Phase != captured.Battle.Phase {
		return false
	}
	phase := int(live.Battle.Phase)
	if phase < 0 || phase >= len(live.Battle.MovieDeliveryClaimed) {
		return false
	}
	if live.Battle.MovieDeliveryClaimed[phase] {
		return false
	}
	live.Battle.MovieDeliveryClaimed[phase] = true
	s.logGameEvent(session, "game-raid-movie-all-ready",
		"char_id", session.selectedCharacterID,
		"instance_id", live.Battle.InstanceID,
		"phase", live.Battle.Phase,
		"members", len(live.Members))
	return true
}

// deliverRaidMovieCompletionToAllMembers 在全员确认后，给每个在线成员各投递一次翻牌。
// 每个成员用**自己的**状态投影，避免把触发者的状态套到别人身上。
func (s *Service) deliverRaidMovieCompletionToAllMembers(trigger *gameSession, captured runtimeRaidState, op uint16) error {
	s.raidMu.Lock()
	live := s.runtimeRaidByMemberLocked(trigger.selectedCharacterID)
	var members []uint16
	if live != nil {
		members = raidMovieVotingMembersLocked(live)
	}
	s.raidMu.Unlock()
	var firstErr error
	for _, id := range members {
		peer, online := s.onlineGameSession(id)
		if !online || peer == nil {
			continue
		}
		peerState, err := s.currentAntonRewardState(peer)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := s.deliverCurrentRaidMovieCompletion(peer, peerState, op); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			s.logGameEvent(peer, "game-raid-movie-delivery-failed", "char_id", id, "error", err)
		}
	}
	return firstErr
}
