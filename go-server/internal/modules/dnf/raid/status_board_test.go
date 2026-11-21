package raid

import (
	"bytes"
	"testing"
)

func TestBuildStatusBoardBody(t *testing.T) {
	body, err := BuildStatusBoardBody([]StatusBoardSquad{
		{
			GroupIndex:      1,
			MemberCharIDs:   []uint16{101, 102},
			UsedCoins:       3,
			DungeonProgress: 0,
		},
		{
			GroupIndex:      2,
			MemberCharIDs:   []uint16{201},
			UsedCoins:       0,
			DungeonProgress: 1,
			ExtraValues:     []uint32{50},
		},
	})
	if err != nil {
		t.Fatalf("BuildStatusBoardBody returned error: %v", err)
	}
	want := []byte{
		2,          // squad count
		1, 0,       // groupIndex = 1
		2,          // memberCount = 2
		101, 0,     // char 101
		102, 0,     // char 102
		3, 0, 0, 0, // usedCoins = 3
		0, 0, 0, 0, // dungeonProgress = 0
		0,          // extraCount = 0
		2, 0,       // groupIndex = 2
		1,          // memberCount = 1
		201, 0,     // char 201
		0, 0, 0, 0, // usedCoins = 0
		1, 0, 0, 0, // dungeonProgress = 1
		1,          // extraCount = 1
		50, 0, 0, 0,
	}
	if !bytes.Equal(body, want) {
		t.Fatalf("body mismatch:\ngot:  %x\nwant: %x", body, want)
	}
}
