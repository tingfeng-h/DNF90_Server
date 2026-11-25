package dnfbridge

// Graduated production policy: when right-clicking a same-party member in town
// (op10 mode=15), skip op9 projection to prevent [localPlayer+0x4C8] hijacking
// and still emit 0x0007 BuildRequestPeerSelectionNotice so the context menu opens.
func samePartyTownSelectionNoticeEnabled() bool {
	return true
}
