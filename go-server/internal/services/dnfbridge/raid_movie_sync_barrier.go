package dnfbridge

import (
	"fmt"

	"longheng.io/server/internal/modules/dnf/raid"
)

func currentRaidMoviePhaseSlot(phase byte) int {
	if phase > 1 {
		return 1
	}
	return int(phase)
}

func (s *Service) currentRaidOnlineMovieSessions(state runtimeRaidState) []*gameSession {
	if s == nil || state.Battle == nil {
		return nil
	}
	b := state.Battle
	candidates := make([]*gameSession, 0, len(state.Members))
	s.raidMu.Lock()
	s.mu.Lock()
	for _, member := range state.Members {
		if member.UserState == 0 {
			continue
		}
		if _, disconnected := s.currentRaidDisconnects[member.CharID]; disconnected {
			continue
		}
		if peer := s.gameSessions[member.CharID]; peer != nil {
			candidates = append(candidates, peer)
		}
	}
	s.mu.Unlock()
	s.raidMu.Unlock()

	sessions := make([]*gameSession, 0, len(candidates))
	for _, peer := range candidates {
		if !s.currentRaidAuthenticated(peer) {
			continue
		}
		if peer.channel.ID != b.ChannelID || peer.channel.Type != b.ChannelType {
			continue
		}
		if b.MemberAccounts[peer.selectedCharacterID] != s.accountIDForSession(peer) {
			continue
		}
		sessions = append(sessions, peer)
	}
	return sessions
}

func (s *Service) handleCurrentRaidSuccessMovieBarrier(session *gameSession, state runtimeRaidState, op uint16) error {
	phaseSlot := currentRaidMoviePhaseSlot(state.Battle.Phase)
	charID := session.selectedCharacterID

	s.raidMu.Lock()
	live := s.runtimeRaidByMemberLocked(charID)
	valid := live != nil && live.Battle != nil && live.Battle.InstanceID == state.Battle.InstanceID && live.Battle.Phase == state.Battle.Phase
	if valid {
		if live.Battle.MovieSyncConfirmed[phaseSlot] == nil {
			live.Battle.MovieSyncConfirmed[phaseSlot] = map[uint16]bool{}
		}
		live.Battle.MovieSyncConfirmed[phaseSlot][charID] = true
		if live.Battle.MovieSyncReleased[phaseSlot] {
			if live.Battle.RewardCardReady[phaseSlot] == nil {
				live.Battle.RewardCardReady[phaseSlot] = map[uint16]bool{}
			}
			live.Battle.RewardCardReady[phaseSlot][charID] = true
		}
	}
	alreadyReleased := valid && live.Battle.MovieSyncReleased[phaseSlot]
	var snapshot runtimeRaidState
	if valid {
		snapshot = cloneRuntimeRaidState(*live)
	}
	s.raidMu.Unlock()
	if !valid {
		return s.sendRaidCommandFailure(session, op)
	}

	o := &session.raidProjection
	o.mu.Lock()
	o.movieIdentity = currentRaidProjectionIdentity{
		InstanceID: state.Battle.InstanceID,
		Phase:      state.Battle.Phase,
		CharID:     charID,
		Generation: session.characterGeneration,
	}
	o.movieCompleted = true
	o.mu.Unlock()
	s.logInfo("dnfbridge raid movie stage", "kind", "game-raid-movie-stage", "char_id", charID, "instance_id", state.Battle.InstanceID, "mode", state.Battle.Mode, "phase", state.Battle.Phase, "stage", 1)

	if err := s.sendCurrentAntonRewardPackets(session, snapshot, []currentAntonRewardPacket{{opcode: op, class: 1, body: []byte{1}}}, false); err != nil {
		return err
	}
	if alreadyReleased {
		return s.sendCurrentAntonSuccessRewardPanel(session, snapshot)
	}
	return s.tryReleaseCurrentRaidMovieBarrier(snapshot)
}

func (s *Service) handleCurrentRaidFailureMovieBarrier(session *gameSession, state runtimeRaidState, op uint16) error {
	phaseSlot := currentRaidMoviePhaseSlot(state.Battle.Phase)
	charID := session.selectedCharacterID

	s.raidMu.Lock()
	live := s.runtimeRaidByMemberLocked(charID)
	valid := live != nil && live.Battle != nil && live.Battle.TimedOut && !live.Battle.PhaseCompleted &&
		live.Battle.InstanceID == state.Battle.InstanceID && live.Battle.Phase == state.Battle.Phase
	if valid {
		if live.Battle.MovieSyncConfirmed[phaseSlot] == nil {
			live.Battle.MovieSyncConfirmed[phaseSlot] = map[uint16]bool{}
		}
		live.Battle.MovieSyncConfirmed[phaseSlot][charID] = true
	}
	alreadyReleased := valid && live.Battle.MovieSyncReleased[phaseSlot]
	var snapshot runtimeRaidState
	if valid {
		snapshot = cloneRuntimeRaidState(*live)
	}
	s.raidMu.Unlock()
	if !valid {
		return s.sendRaidCommandFailure(session, op)
	}

	if err := s.sendCurrentAntonRewardPackets(session, snapshot, []currentAntonRewardPacket{{opcode: op, class: 1, body: []byte{1}}}, false); err != nil {
		return err
	}
	if alreadyReleased {
		return s.sendCurrentRaidFailureMovieClose(session, snapshot)
	}
	return s.tryReleaseCurrentRaidMovieBarrier(snapshot)
}

func (s *Service) tryReleaseCurrentRaidMovieBarrier(state runtimeRaidState) error {
	if s == nil || state.Battle == nil {
		return nil
	}
	b := state.Battle
	if !b.PhaseCompleted && !b.TimedOut {
		return nil
	}
	phaseSlot := currentRaidMoviePhaseSlot(b.Phase)
	online := s.currentRaidOnlineMovieSessions(state)
	if len(online) == 0 {
		return nil
	}

	s.raidMu.Lock()
	live := s.raids[state.RaidKey]
	if live == nil || live.Battle == nil || live.Battle.InstanceID != b.InstanceID || live.Battle.Phase != b.Phase {
		s.raidMu.Unlock()
		return nil
	}
	if live.Battle.MovieSyncReleased[phaseSlot] {
		s.raidMu.Unlock()
		return nil
	}
	confirmed := live.Battle.MovieSyncConfirmed[phaseSlot]
	if len(confirmed) == 0 {
		s.raidMu.Unlock()
		return nil
	}
	for _, peer := range online {
		if !confirmed[peer.selectedCharacterID] {
			s.raidMu.Unlock()
			return nil
		}
	}
	live.Battle.MovieSyncReleased[phaseSlot] = true
	if !live.Battle.TimedOut {
		if live.Battle.RewardCardReady[phaseSlot] == nil {
			live.Battle.RewardCardReady[phaseSlot] = map[uint16]bool{}
		}
		for _, peer := range online {
			live.Battle.RewardCardReady[phaseSlot][peer.selectedCharacterID] = true
		}
	}
	snapshot := cloneRuntimeRaidState(*live)
	s.raidMu.Unlock()

	var firstErr error
	for _, peer := range online {
		if !currentAntonResultDelivered(peer, snapshot) {
			continue
		}
		var err error
		if snapshot.Battle.TimedOut {
			err = s.sendCurrentRaidFailureMovieClose(peer, snapshot)
		} else {
			err = s.sendCurrentAntonSuccessRewardPanel(peer, snapshot)
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Service) sendCurrentRaidFailureMovieClose(session *gameSession, state runtimeRaidState) error {
	b := state.Battle
	if session == nil || b == nil {
		return nil
	}
	identity := currentRaidProjectionIdentity{
		InstanceID: b.InstanceID,
		Phase:      b.Phase,
		CharID:     session.selectedCharacterID,
		Generation: session.characterGeneration,
	}
	// NOTI587 state 6 closes the client's movie controller (0x02001238 ->
	// 0x18fc7d0) without opening the reward card panel.
	packets := []currentAntonRewardPacket{{opcode: 587, class: 0, body: []byte{6, 0}}}
	if err := s.sendCurrentAntonRewardPackets(session, state, packets, false); err != nil {
		return err
	}
	o := &session.raidProjection
	o.mu.Lock()
	o.movieIdentity = identity
	o.movieCompleted = true
	o.mu.Unlock()
	return nil
}

func (s *Service) sendCurrentAntonSuccessRewardPanel(session *gameSession, state runtimeRaidState) error {
	if flow := state.Battle.RewardFlow[state.Battle.Phase]; currentRaidOwnerRewardsCommitted(flow, session.selectedCharacterID) && flow.squadRequired() {
		return s.sendCurrentAntonRewardFlowProjection(session, state)
	}
	// Successful reward state4 has stateArgument0 in either phase. Argument1
	// is reserved for failure and must not be derived from the phase index.
	packets := []currentAntonRewardPacket{{opcode: 587, class: 0, body: []byte{4, 0}}}
	if currentAntonPhaseClearCount(state.Battle, session.selectedCharacterID) < state.Battle.AntonRewards.Phases[state.Battle.Phase].MinimumClears {
		return s.sendCurrentAntonRewardPackets(session, state, packets, true)
	}
	memberIdx := runtimeRaidMemberIndex(state, session.selectedCharacterID)
	if memberIdx < 0 {
		return nil
	}
	member := state.Members[memberIdx]
	peers := state.Members
	if flow := state.Battle.RewardFlow[state.Battle.Phase]; flow != nil {
		peers = nil
		found := false
		for _, owner := range flow.Recipients {
			if owner.CharacterID == member.CharID {
				member.GroupIndex, found = owner.Group, true
			}
			peers = append(peers, runtimeRaidMemberState{CharID: owner.CharacterID, GroupIndex: owner.Group})
		}
		if !found {
			return nil
		}
	}
	var cards []raid.RewardDisplayCard
	visible := s.currentRaidPartyRewardDisplayMembers(session, state)
	for _, peer := range peers {
		if visible[peer.CharID] && currentAntonPhaseClearCount(state.Battle, peer.CharID) >= state.Battle.AntonRewards.Phases[state.Battle.Phase].MinimumClears {
			card := raid.RewardDisplayCard{CharacterID: peer.CharID, CardType: byte(len(cards) + 1)}
			if selected := state.Battle.RewardCards[state.Battle.Phase][peer.CharID]; selected != nil && selected.Committed {
				card.CardType = selected.Card
				card.Items = currentAntonPartyDisplayItems(selected.Plan)
			}
			cards = append(cards, card)
		}
	}
	if len(cards) == 0 || len(cards) > 4 {
		return fmt.Errorf("Anton reward party exceeds native four-card layout")
	}
	for _, kind := range []byte{1, 0} {
		projection := cards
		if kind == 0 {
			projection = append([]raid.RewardDisplayCard(nil), cards...)
			for i := range projection {
				projection[i].Items = nil
			}
		}
		body, err := raid.BuildRewardListBody(kind, projection)
		if err != nil {
			return err
		}
		packets = append(packets, currentAntonRewardPacket{opcode: 600, body: body})
	}
	return s.sendCurrentAntonRewardPackets(session, state, packets, true)
}
