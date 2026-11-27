package dnfbridge

import (
	"context"
	"strconv"
	"sync"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const (
	characterFriendSameChannelPrimeLimit = 64
)

// Graduated production policy: character friend realtime presence notifications are always enabled.
func characterFriendRealtimePresenceEnabled() bool {
	return true
}

type characterFriendWatcherHub struct {
	mu        sync.Mutex
	byTarget  map[uint16]map[*gameSession]struct{}
	bySession map[*gameSession]map[uint16]struct{}
}

var (
	characterFriendWatcherHubsMu sync.Mutex
	characterFriendWatcherHubs   = make(map[*Service]*characterFriendWatcherHub)
)

func (s *Service) characterFriendWatcherHub(create bool) *characterFriendWatcherHub {
	if s == nil {
		return nil
	}
	characterFriendWatcherHubsMu.Lock()
	defer characterFriendWatcherHubsMu.Unlock()
	hub := characterFriendWatcherHubs[s]
	if hub == nil && create {
		hub = &characterFriendWatcherHub{
			byTarget:  make(map[uint16]map[*gameSession]struct{}),
			bySession: make(map[*gameSession]map[uint16]struct{}),
		}
		characterFriendWatcherHubs[s] = hub
	}
	return hub
}

func (s *Service) unregisterCurrentCharacterFriendWatcher(session *gameSession) {
	if s == nil || session == nil {
		return
	}
	hub := s.characterFriendWatcherHub(false)
	if hub == nil {
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	prev := hub.bySession[session]
	if len(prev) == 0 {
		delete(hub.bySession, session)
		return
	}
	for targetID := range prev {
		if watchers := hub.byTarget[targetID]; watchers != nil {
			delete(watchers, session)
			if len(watchers) == 0 {
				delete(hub.byTarget, targetID)
			}
		}
	}
	delete(hub.bySession, session)
}

// recordCurrentCharacterFriendWatchersAndPrime updates the in-memory reverse
// watcher index (targetCharID -> watching online sessions) for `session` and
// silently primes same-channel cross-room mode8 actor caches for online
// character friends right before op273 is sent.
func (s *Service) recordCurrentCharacterFriendWatchersAndPrime(
	ctx context.Context,
	session *gameSession,
	identity boundGameSessionCharacter,
	rows []dnfrepo.CharacterRecord,
	source string,
) {
	if s == nil || session == nil || !characterFriendRealtimePresenceEnabled() ||
		s.options.currentClientVersion != "official-24.9" ||
		!s.boundGameSessionCharacterCurrent(identity) {
		return
	}
	nextTargets := make(map[uint16]struct{}, len(rows))
	for _, row := range rows {
		parsed, err := strconv.ParseUint(row.CharacterID, 10, 16)
		if err != nil || parsed == 0 {
			continue
		}
		nextTargets[uint16(parsed)] = struct{}{}
	}
	hub := s.characterFriendWatcherHub(true)
	hub.mu.Lock()
	prev := hub.bySession[session]
	for oldTarget := range prev {
		if _, keep := nextTargets[oldTarget]; !keep {
			if watchers := hub.byTarget[oldTarget]; watchers != nil {
				delete(watchers, session)
				if len(watchers) == 0 {
					delete(hub.byTarget, oldTarget)
				}
			}
		}
	}
	if len(nextTargets) == 0 {
		delete(hub.bySession, session)
	} else {
		hub.bySession[session] = nextTargets
		for targetID := range nextTargets {
			watchers := hub.byTarget[targetID]
			if watchers == nil {
				watchers = make(map[*gameSession]struct{})
				hub.byTarget[targetID] = watchers
			}
			watchers[session] = struct{}{}
		}
	}
	hub.mu.Unlock()

	if !guildRosterSameChannelPrimeEnabled() || !crossSceneMode8PrimeAllowedForSource(source) || len(rows) == 0 {
		return
	}
	session.dungeon.mu.Lock()
	inDungeon := session.dungeon.runtime != nil
	session.dungeon.mu.Unlock()
	if inDungeon {
		return
	}
	viewerChannel := currentTownPresenceChannelID(session)
	if viewerChannel <= 0 || viewerChannel >= 255 {
		return
	}
	primed := 0
	for _, row := range rows {
		if primed >= characterFriendSameChannelPrimeLimit {
			break
		}
		parsed, err := strconv.ParseUint(row.CharacterID, 10, 16)
		if err != nil || parsed == 0 || uint16(parsed) == identity.character {
			continue
		}
		targetID := uint16(parsed)
		peer, online := s.onlineGameSessionCharacterSnapshot(targetID)
		if !online || peer.session == nil || peer.session == session ||
			s.accountIDForSession(peer.session) != row.AccountID {
			continue
		}
		if currentTownPresenceChannelID(peer.session) != viewerChannel {
			continue
		}
		if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(identity.character, targetID) {
			continue
		}
		delivered, err := s.primeSingleSameChannelGuildMode8(ctx, identity, peer, viewerChannel)
		if err == nil && delivered {
			primed++
		}
	}
	if primed > 0 {
		s.logGameEvent(session, "game-character-friend-same-channel-primed",
			"source", source,
			"viewer_char_id", identity.character,
			"viewer_channel", viewerChannel,
			"primed_friends", primed)
	}
}

func (s *Service) refreshCurrentCharacterFriendPresence(session *gameSession, source string) {
	if s == nil || session == nil {
		return
	}
	charID := selectedCharacterID(session)
	if charID == 0 {
		return
	}
	s.refreshCurrentCharacterFriendPresenceForCharacterID(session, charID, source)
}

// refreshCurrentCharacterFriendPresenceForCharacterID pushes a fresh op273
// character-friend snapshot to every online session currently watching
// `targetCharID`, so character login, logout, character switch, and channel
// switch update other players' Character Friend tab in real time.
func (s *Service) refreshCurrentCharacterFriendPresenceForCharacterID(exclude *gameSession, targetCharID uint16, source string) {
	if s == nil || targetCharID == 0 || !characterFriendRealtimePresenceEnabled() ||
		s.options.currentClientVersion != "official-24.9" {
		return
	}
	hub := s.characterFriendWatcherHub(false)
	if hub == nil {
		return
	}
	hub.mu.Lock()
	watchersMap := hub.byTarget[targetCharID]
	if len(watchersMap) == 0 {
		hub.mu.Unlock()
		return
	}
	candidates := make([]*gameSession, 0, len(watchersMap))
	for watcher := range watchersMap {
		if watcher != nil && watcher != exclude {
			candidates = append(candidates, watcher)
		}
	}
	hub.mu.Unlock()
	if len(candidates) == 0 {
		return
	}

	var peers []boundGameSessionCharacter
	var stale []*gameSession
	s.mu.Lock()
	for _, watcher := range candidates {
		charID := watcher.selectedCharacterID
		if charID == 0 || s.gameSessions[charID] != watcher {
			stale = append(stale, watcher)
			continue
		}
		peers = append(peers, boundGameSessionCharacter{
			session:    watcher,
			character:  charID,
			generation: watcher.characterGeneration,
		})
	}
	s.mu.Unlock()
	for _, dead := range stale {
		s.unregisterCurrentCharacterFriendWatcher(dead)
	}
	for _, peer := range peers {
		peer := peer
		_ = s.postGameSessionCharacterEvent(peer.session, "character_friend_"+source, peer.character, peer.generation, func() error {
			if !s.boundGameSessionCharacterCurrent(peer) {
				return nil
			}
			return s.sendCurrentCharacterFriendList(peer.session, "presence_"+source)
		})
	}
}
