package dnfbridge

// Graduated production policy: guild member character-info and adventure-group
// inspections via C2S819 require viewer and target to be on the same channel.
func currentGuildMemberProfileSameChannelOnlyEnabled() bool {
	return true
}

func currentGuildMemberProfileChannelAllowed(viewerChannel int, targetChannel byte) bool {
	return viewerChannel > 0 && viewerChannel == int(targetChannel)
}
