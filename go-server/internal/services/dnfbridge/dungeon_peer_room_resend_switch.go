package dnfbridge

import (
	"os"
	"strings"
)

// dungeonPeerRoomResendEnv 控制同一场地下城内“换房间 op37”是否重复补投队友全量 actor：
//   - 默认（空 / off / 0 / false）：按整场地下城去重（归零 RoomX/RoomY/MapID，保留
//     PartyID/RaidKey/RaidInstanceID/RaidRunID/DungeonID/MazeIndex/Seed）。
//     首房进图仍完整执行一次补投，同图过门换房不再重复向客户端发全量 mode0/mode1/subtype0/op9；
//   - on / 1 / true / legacy：恢复逐房间补投，便于随时 A/B 回退对比。
const dungeonPeerRoomResendEnv = "DNFBRIDGE_DUNGEON_PEER_ROOM_RESEND"

func dungeonPeerRoomResendEnabled() bool {
	return parseDungeonPeerRoomResendEnabled(os.Getenv(dungeonPeerRoomResendEnv))
}

func parseDungeonPeerRoomResendEnabled(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "on", "1", "true", "yes", "legacy":
		return true
	default:
		return false
	}
}

func normalizeDungeonPeerActorResendKey(key dungeonPartyRunKey) dungeonPartyRunKey {
	if dungeonPeerRoomResendEnabled() {
		return key
	}
	key.RoomX = 0
	key.RoomY = 0
	key.MapID = 0
	return key
}
