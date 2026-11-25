package dnfbridge

import (
	"errors"
	"fmt"
	"time"
)

// deferPartyLoadingSocial moves non-critical social/DB notifications (friend
// lists, presence fan-out, wedding replay) off the synchronous op37 finish-
// loading hot path onto the gameplay time queue with a small per-character
// stagger (800ms + (charID%4)*150ms). This prevents 4 party members returning
// to town or entering a dungeon simultaneously from saturating a 2-core CPU
// and stalling their TCP read loops for 230-780ms.
func (s *Service) deferPartyLoadingSocial(session *gameSession, source string) {
	identity, ok := s.boundGameSessionCharacterSnapshot(session)
	if !ok {
		return
	}
	owner := snapshotPartyDungeonLoadingOwner(session)
	name := fmt.Sprintf("party-loading-social:%p:%d", session, identity.generation)
	delay := 800*time.Millisecond + time.Duration(identity.character%4)*150*time.Millisecond
	err := s.ensureGameplayTimeQueue().ScheduleAfter(name, delay, func(time.Time) {
		err := s.tryPostGameSessionCharacterEvent(session, "party-loading-social", identity.character, identity.generation, func() error {
			if !s.boundGameSessionCharacterCurrent(identity) || snapshotPartyDungeonLoadingOwner(session) != owner {
				return nil
			}
			return s.sendLoadingSocial(session, source)
		})
		if err != nil {
			s.logGameEvent(session, "game-party-loading-social-deferred", "error", err)
		}
	})
	if err != nil {
		s.logGameEvent(session, "game-party-loading-social-schedule-failed", "error", err)
	}
}

func (s *Service) shouldDeferLoadingSocial(session *gameSession) bool {
	if session == nil || session.events == nil || session.events.loop == nil {
		return false
	}
	if currentDungeonSceneActive(session) {
		return true
	}
	return runtimePartyStateSnapshot(session).PartyID > 0
}

func (s *Service) sendLoadingSocial(session *gameSession, source string) error {
	if err := s.sendCurrentCharacterFriendList(session, source); err != nil {
		return err
	}
	if err := s.sendCurrentAccountFriendLists(session, source); err != nil {
		return err
	}
	s.refreshCurrentAccountFriendPresence(session, source)
	s.refreshCurrentCharacterFriendPresence(session, source)
	return errors.Join(s.queueCurrentWeddingBenefitsReplay(session, source), s.queueCurrentWeddingCeremonyRecovery(session))
}
