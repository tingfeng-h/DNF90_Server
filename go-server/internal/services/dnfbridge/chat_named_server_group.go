package dnfbridge

import (
	"os"
	"strings"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

const namedChatServerGroupEnv = "DNFBRIDGE_NAMED_CHAT_SERVER_GROUP"

// currentNamedChatServerGroupLegacyEnabled returns true only when the operator
// explicitly opts back into the historical channel-byte wire behavior.
func currentNamedChatServerGroupLegacyEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(namedChatServerGroupEnv))) {
	case "legacy", "channel", "off", "0", "false":
		return true
	default:
		return false
	}
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

// namedChatServerGroupEnabled 是"原生字段模式是否启用"的可读别名。
func namedChatServerGroupEnabled() bool { return !currentNamedChatServerGroupLegacyEnabled() }
