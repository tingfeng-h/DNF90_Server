package dnfbridge

import "errors"

// 移植自「轻装上阵」分支 refreshCurrentCharacterFriendPresence 的语义：
// 某个角色的在线状态/频道/角色切换发生变化时，通知所有把他当好友的观察者
// 重新投影角色好友列表。上游只在登录时投影一次，观察者侧因此会残留旧状态
// （切频道后好友"消失"就是这种残留）。
//
// 这里不搬支线的行级投影函数，直接复用既有 sendCurrentCharacterFriendList，
// 语义等价且不依赖支线内部结构。
func (s *Service) refreshCurrentCharacterFriendPresence(characterID uint16, source string) {
	if s == nil || s.options.currentClientVersion != "official-24.9" || characterID == 0 {
		return
	}
	var peers []boundGameSessionCharacter
	s.mu.Lock()
	for id, session := range s.gameSessions {
		if session != nil && id != characterID {
			peers = append(peers, boundGameSessionCharacter{session: session, character: id, generation: session.characterGeneration})
		}
	}
	s.mu.Unlock()
	if len(peers) == 0 {
		return
	}
	var submitErrors []error
	for _, peer := range peers {
		peer := peer
		err := s.postGameSessionCharacterEvent(peer.session, "character_friend_"+source, peer.character, peer.generation, func() error {
			if !s.boundGameSessionCharacterCurrent(peer) || !peer.session.loginPacketSeen || !peer.session.accountFromToken {
				return nil
			}
			if err := s.sendCurrentCharacterFriendList(peer.session, source); err != nil {
				s.logGameEvent(peer.session, "game-character-friend-presence-failed",
					"source", source, "target_char_id", characterID, "error", err)
				return err
			}
			s.logGameEvent(peer.session, "game-character-friend-presence-refreshed",
				"source", source, "target_char_id", characterID, "observer_char_id", peer.character)
			return nil
		})
		if err != nil {
			submitErrors = append(submitErrors, err)
		}
	}
	_ = errors.Join(submitErrors...)
}
