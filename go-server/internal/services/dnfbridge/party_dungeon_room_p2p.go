package dnfbridge

import (
	"encoding/binary"
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// 副本内房间 P2P 表（A21 `DungeonRoomP2pBuilder.cs`，SC Noti 0x0184）。
//
// A21/86JP 的注释写明：副本里"队友互相直连、填血量、清连接中"用的是 **0x184 房间表
// （18B/peer）**，与城镇组队的 **0x0B（22B/成员）是两张不同的表**，"别混用"。
// 我们此前进副本只发 0x000b（城镇表，见 sendDungeonEntryPeerEndpointRemotes），
// 从没发过 0x184 —— 现场表现正好是"副本里看不到队友，重选角色后才有"。
//
// 布局（逐字节照 A21，含其"待真机验证"字段）：
//   HEADER 8B : u8 f1=0, u8 f2=0, u32 session/room=partyId, u8 globalFlag=0, u8 mode=1
//   BODY      : u16 roomId=partyId, u8 count1=0, u8 count2=N,
//               N × 18B { u8 slot, u16 uid(LE), u8[4] inner_ip, u8[4] outer_ip,
//                         u16 port(BE), u8 nat=0, u16 extra1=0, u16 extra2=0 }
//   TRAILER   : u8 0
//
// 开关：DNFBRIDGE_PARTY_DUNGEON_ROOM_P2P
//
//	默认 on（实验档）；off / 0 / false / no / legacy = 一个字节都不发。
const (
	dungeonRoomP2PMsgID uint16 = 0x0184
	dungeonRoomP2PEnv          = "DNFBRIDGE_PARTY_DUNGEON_ROOM_P2P"
)

func dungeonRoomP2PEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(dungeonRoomP2PEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// sendDungeonRoomP2PTable 给单个会话发一条副本房间 P2P 表。
func (s *Service) sendDungeonRoomP2PTable(session *gameSession, state alignedcmd.PartyState, source string) error {
	if s == nil || session == nil || session.selectedCharacterID == 0 {
		return nil
	}
	if !dungeonRoomP2PEnabled() {
		return nil
	}
	members := runtimePartyMembers(state)
	if len(members) < 2 {
		return nil
	}
	body := s.buildDungeonRoomP2PBody(session, state, members)
	s.logGameEvent(session, "game-party-dungeon-room-p2p-send",
		"source", source,
		"char_id", session.selectedCharacterID,
		"party_id", state.PartyID,
		"room_id", state.PartyID,
		"peer_count", len(members),
		"body_len", len(body),
		"expected_body_len", 13+18*len(members),
		"msg_id", dungeonRoomP2PMsgID,
		"classification", 0,
		"env", dungeonRoomP2PEnv,
		"mode", "a21_dungeon_room_p2p")
	return s.sendGameUpperRawClassCodec(session, dungeonRoomP2PMsgID, body, 0, true)
}

func (s *Service) buildDungeonRoomP2PBody(
	session *gameSession,
	state alignedcmd.PartyState,
	members []alignedcmd.PartyMemberState,
) []byte {
	body := make([]byte, 0, 13+18*len(members))
	// HEADER（8B）
	body = append(body, 0, 0)
	body = binary.LittleEndian.AppendUint32(body, uint32(state.PartyID))
	body = append(body, 0, 1) // globalFlag, mode=1（A21：1 触发 BODY 解析）
	// BODY
	body = binary.LittleEndian.AppendUint16(body, uint16(state.PartyID)) // roomId
	body = append(body, 0)                                               // count1 = 0
	body = append(body, byte(len(members)))                              // count2
	relay := s.partyRelayBackend()
	for _, member := range members {
		ip := [4]byte{127, 0, 0, 1}
		port := uint16(defaultPartyPeerPort)
		if relay != nil {
			if address, resolved, ok := relay.Endpoint(uint16(state.PartyID), session.selectedCharacterID, member.UserID); ok {
				ip = address
				port = resolved
			}
		}
		body = append(body, byte(member.Slot))
		body = binary.LittleEndian.AppendUint16(body, member.UserID)
		body = append(body, ip[:]...)  // inner
		body = append(body, ip[:]...)  // outer（同局域网时与 inner 相同）
		body = append(body, byte(port>>8), byte(port&0xFF)) // port 按网络序（大端）
		body = append(body, 0)                              // nat_type = 0
		body = binary.LittleEndian.AppendUint16(body, 0)    // extra1
		body = binary.LittleEndian.AppendUint16(body, 0)    // extra2
	}
	body = append(body, 0) // trailer
	return body
}
