package dnfbridge

import "testing"

func TestNormalizeDungeonPeerActorResendKey(t *testing.T) {
	roomA := dungeonPartyRunKey{PartyID: 1, DungeonID: 100, MazeIndex: 2, Seed: 42, RoomX: 0, RoomY: 0, MapID: 1001}
	roomB := dungeonPartyRunKey{PartyID: 1, DungeonID: 100, MazeIndex: 2, Seed: 42, RoomX: 1, RoomY: 0, MapID: 1002}
	nextRun := dungeonPartyRunKey{PartyID: 1, DungeonID: 100, MazeIndex: 2, Seed: 43, RoomX: 0, RoomY: 0, MapID: 1001}

	t.Setenv(dungeonPeerRoomResendEnv, "")
	if normalizeDungeonPeerActorResendKey(roomA) != normalizeDungeonPeerActorResendKey(roomB) {
		t.Fatalf("same dungeon rooms should share resend dedupe key by default")
	}
	if normalizeDungeonPeerActorResendKey(roomA) == normalizeDungeonPeerActorResendKey(nextRun) {
		t.Fatalf("different dungeon run seed must not share resend dedupe key")
	}

	t.Setenv(dungeonPeerRoomResendEnv, "on")
	if normalizeDungeonPeerActorResendKey(roomA) == normalizeDungeonPeerActorResendKey(roomB) {
		t.Fatalf("legacy on mode must keep per-room resend key")
	}
}
