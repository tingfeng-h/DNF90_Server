package dnfbridge

import (
	"os"
	"strings"
)

const guildProfileSameChannelOnlyEnv = "DNFBRIDGE_GUILD_PROFILE_SAME_CHANNEL_ONLY"

// currentGuildMemberProfileSameChannelOnlyEnabled enforces that guild member
// character-info and adventure-group inspections via C2S819 require the viewer
// and target to be on the same channel (matching the native cross-channel menu
// restriction). Set DNFBRIDGE_GUILD_PROFILE_SAME_CHANNEL_ONLY=off to restore
// cross-channel guild inspection.
func currentGuildMemberProfileSameChannelOnlyEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(guildProfileSameChannelOnlyEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func currentGuildMemberProfileChannelAllowed(viewerChannel int, targetChannel byte) bool {
	if !currentGuildMemberProfileSameChannelOnlyEnabled() {
		return true
	}
	return viewerChannel > 0 && viewerChannel == int(targetChannel)
}
