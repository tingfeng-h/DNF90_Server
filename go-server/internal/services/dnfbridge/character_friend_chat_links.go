package dnfbridge

import (
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	"time"
)

// Native chat menus copy the displayed channel into op290's server-group
// byte. Retain only links actually delivered by this server; never treat an
// arbitrary nonlocal group as a local name lookup. The wire chat format stays
// unchanged because native whisper replies need its channel byte.
type characterFriendChatLinkKey struct {
	name    string
	channel byte
}
type characterFriendChatLink struct {
	viewer             boundGameSessionCharacter
	character, account string
	group              byte
	delivered          time.Time
}

func (s *Service) rememberCharacterFriendChatLink(viewer boundGameSessionCharacter, row dnfrepo.CharacterRecord, channel byte) {
	if !s.boundGameSessionCharacterCurrent(viewer) || row.CharacterID == "" || row.AccountID == "" {
		return
	}
	viewer.session.friendChatLinkMu.Lock()
	defer viewer.session.friendChatLinkMu.Unlock()
	links := viewer.session.friendChatLinks
	if links == nil {
		links = make(map[characterFriendChatLinkKey]characterFriendChatLink)
		viewer.session.friendChatLinks = links
	}
	for key, link := range links {
		if link.viewer != viewer {
			delete(links, key)
		}
	}
	key := characterFriendChatLinkKey{row.Name, channel}
	if _, exists := links[key]; !exists && len(links) >= 128 {
		var oldest characterFriendChatLinkKey
		var when time.Time
		for k, link := range links {
			if when.IsZero() || link.delivered.Before(when) {
				oldest, when = k, link.delivered
			}
		}
		delete(links, oldest)
	}
	links[key] = characterFriendChatLink{viewer, row.CharacterID, row.AccountID, currentAdventureCharacterServerGroup(row), time.Now()}
}

func (s *Service) characterFriendChatLinkMatches(viewer boundGameSessionCharacter, request currentCharacterFriendMutation, target dnfrepo.CharacterRecord) bool {
	if request.TargetID != 0 && request.TargetID != 0xffff {
		return false
	}
	viewer.session.friendChatLinkMu.Lock()
	link, found := viewer.session.friendChatLinks[characterFriendChatLinkKey{request.Name, request.ServerGroup}]
	viewer.session.friendChatLinkMu.Unlock()
	return found && link.viewer == viewer && link.character == target.CharacterID && link.account == target.AccountID &&
		link.group == currentAdventureCharacterServerGroup(target) && s.boundGameSessionCharacterCurrent(viewer)
}
