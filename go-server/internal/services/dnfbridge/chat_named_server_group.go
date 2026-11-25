package dnfbridge

import (
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

// Graduated production policy: named chat packets (NOTI 65 / NOTI 371) always
// emit the sender's server_group_id rather than the connection channel byte.
func currentNamedChatServerGroupLegacyEnabled() bool {
	return false
}

// NOTI65 (VA 0x0201B510) and NOTI371 (VA 0x02A04520) store the single u8 field
// after the sender DSTR into chat record +0x44 as the sender's server group ID
// (e.g. 0 = 90CN), NOT the connection's channel ID. Writing channel ID (e.g. 10)
// causes the chat click handler (0x021E6F96 -> 0x0296D810) to treat a same-server
// guild/party speaker as a foreign-server character (e.g. 华北1服务器) and drop
// party-invite / character-info menu items.
func currentNamedChatSenderServerGroup(session *gameSession, character dnfrepo.CharacterRecord) (byte, bool) {
	if currentNamedChatServerGroupLegacyEnabled() {
		return currentSessionOwnerChannelField(session)
	}
	return currentChatCharacterServerGroup(character)
}

func namedChatServerGroupEnabled() bool {
	return !currentNamedChatServerGroupLegacyEnabled()
}
