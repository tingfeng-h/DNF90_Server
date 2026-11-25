package dnfbridge

import "testing"

func TestSamePartyTownSelectionNoticeAlwaysEnabled(t *testing.T) {
	if !samePartyTownSelectionNoticeEnabled() {
		t.Fatal("same-party town selection notice must stay enabled by default")
	}
}
