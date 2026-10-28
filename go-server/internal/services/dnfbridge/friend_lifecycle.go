package dnfbridge

import "errors"

// refreshCurrentFriendLifecycle reprojects durable friend identities after a
// committed rename or deletion. Work belongs to each observer's event queue;
// a character switch or replacement connection invalidates the captured view.
// The projection owners retain their last successful view on repository errors.
func (s *Service) refreshCurrentFriendLifecycle(characterID, source string) error {
	if s == nil || s.options.currentClientVersion != "official-24.9" || characterID == "" {
		return nil
	}
	var peers []boundGameSessionCharacter
	s.mu.Lock()
	for id, session := range s.gameSessions {
		if session != nil {
			peers = append(peers, boundGameSessionCharacter{session: session, character: id, generation: session.characterGeneration})
		}
	}
	s.mu.Unlock()
	var submitErrors []error
	for _, peer := range peers {
		peer := peer
		err := s.postGameSessionCharacterEvent(peer.session, "friend_lifecycle_"+source, peer.character, peer.generation, func() error {
			if !s.boundGameSessionCharacterCurrent(peer) || !peer.session.loginPacketSeen || !peer.session.accountFromToken {
				return nil
			}
			characterErr := s.sendCurrentCharacterFriendList(peer.session, source)
			if !s.boundGameSessionCharacterCurrent(peer) {
				return characterErr
			}
			accountErr := s.sendCurrentAccountFriendLists(peer.session, source)
			if err := errors.Join(characterErr, accountErr); err != nil {
				s.logGameEvent(peer.session, "game-friend-lifecycle-refresh-failed", "source", source, "target_char_id", characterID, "error", err)
				return err
			}
			s.logGameEvent(peer.session, "game-friend-lifecycle-refreshed", "source", source, "target_char_id", characterID, "observer_char_id", peer.character)
			return nil
		})
		if err != nil {
			submitErrors = append(submitErrors, err)
		}
	}
	return errors.Join(submitErrors...)
}
