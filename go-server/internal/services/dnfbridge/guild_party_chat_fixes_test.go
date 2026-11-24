package dnfbridge

import (
	"bytes"
	"context"
	"testing"

	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

func TestNamedChatSenderServerGroupUsesCharacterServerGroupNotChannel(t *testing.T) {
	s, peers, _ := guildHallFixture(t)
	// Sender is on channel 10 (or whatever townActorOwnerChannel is), server_group_id = 0 (90CN).
	peers[29].townActorOwnerChannel = 10
	peers[29].channel.ID = 10
	peers[30].townPositionSnapshot.TownID = 38

	if err := s.handleCurrentPeerChat(peers[29], 17, guildChatRequest(29, "同服测试")); err != nil {
		t.Fatal(err)
	}
	var want bytes.Buffer
	w := packetWriter{}
	w.writeBytes([]byte{currentGuildChatType, 0, 0, 0})
	w.writeRawDstr(rosterNameBytes("角色29"))
	w.writeByte(0) // server_group_id=0 (90CN), NOT channel=10 (华北1)
	w.writeRawDstr(rosterNameBytes("同服测试"))
	w.writeRawDstr(nil)
	want.Write(w.bytes())

	wire := peers[30].conn.(*bufferConn).write.Bytes()
	packet, rest := splitCurrentGameServerUpperPacketAuto(t, wire)
	if len(rest) != 0 || packet.Header.MsgID != currentOfficialWhisperNotifyOpcode || !bytes.Equal(packet.Body, want.Bytes()) {
		t.Fatalf("guild chat notice=%x want=%x", packet.Body, want.Bytes())
	}

	// Verify explicit rollback switch restores legacy channel byte.
	t.Setenv(namedChatServerGroupEnv, "legacy")
	got, ok := currentNamedChatSenderServerGroup(peers[29], dnfrepo.CharacterRecord{Stats: map[string]int64{"server_group_id": 0}})
	if !ok || got != 10 {
		t.Fatalf("legacy mode got=%d ok=%v want=10", got, ok)
	}
}

func TestGuildMemberProfileChannelGuardBlocksCrossChannelByDefault(t *testing.T) {
	f, _ := guildMemberProfileFixture(t)
	// Viewer on channel 91, target actor on channel 11 (different channels).
	f.target.channel.ID, f.target.connectionTownActorOwnerChannel = 91, 91
	f.actor.channel.ID, f.actor.connectionTownActorOwnerChannel = 11, 11

	request := []byte{11, 17, 0, 0, 1}
	if err := f.service.handleGameCommand(f.target, byte(dnfenum.GameCmdCommand), currentGuildMemberProfileOpcode, request); err != nil {
		t.Fatal(err)
	}
	if f.targetConn.write.Len() != 0 {
		t.Fatalf("cross-channel guild profile should be blocked, wrote %d bytes", f.targetConn.write.Len())
	}

	// Same channel must succeed.
	f.target.channel.ID, f.target.connectionTownActorOwnerChannel = 11, 11
	if err := f.service.handleGameCommand(f.target, byte(dnfenum.GameCmdCommand), currentGuildMemberProfileOpcode, request); err != nil {
		t.Fatal(err)
	}
	if f.targetConn.write.Len() == 0 {
		t.Fatal("same-channel guild profile should succeed")
	}
}

func TestPartySameChannelCrossRoomInviteAllowsAndPrimesMode8(t *testing.T) {
	f := newTownPartyActorHPFixture(t, false)
	s := f.service
	s.options.currentClientVersion = "official-24.9"
	// Same channel (10), different rooms (town 1 vs town 2).
	f.target.channel.ID, f.target.connectionTownActorOwnerChannel = 10, 10
	f.actor.channel.ID, f.actor.connectionTownActorOwnerChannel = 10, 10
	if !s.allowSameChannelCrossRoomPartyRequest(f.target, f.target.selectedCharacterID, f.actor.selectedCharacterID, 0) {
		t.Fatal("expected same-channel cross-room party invite (mode=0) to be allowed")
	}
	if s.allowSameChannelCrossRoomPartyRequest(f.target, f.target.selectedCharacterID, f.actor.selectedCharacterID, 1) {
		t.Fatal("trade (mode=1) must NOT bypass PeerInSameArea")
	}
	if s.allowSameChannelCrossRoomPartyRequest(f.target, f.target.selectedCharacterID, f.actor.selectedCharacterID, 15) {
		t.Fatal("scene click (mode=15) must NOT bypass PeerInSameArea")
	}
	// Different channels must be rejected.
	f.actor.channel.ID, f.actor.connectionTownActorOwnerChannel = 11, 11
	if s.allowSameChannelCrossRoomPartyRequest(f.target, f.target.selectedCharacterID, f.actor.selectedCharacterID, 0) {
		t.Fatal("different-channel party invite must be rejected")
	}
	_ = context.Background()
	_ = party.RequestPeerRequest{}
}

func TestPartyCrossRoomTownProjectionSwitches(t *testing.T) {
	if !partyPostSnapshotTownRefreshSkipEnabled() {
		t.Fatal("post-snapshot town refresh skip should be enabled by default")
	}
	if !partySeriaActorNoOp23Enabled() {
		t.Fatal("seria actor no-op23 projection should be enabled by default")
	}
	t.Setenv(partyPostSnapshotTownRefreshSkipEnv, "legacy")
	t.Setenv(partySeriaActorNoOp23Env, "legacy")
	if partyPostSnapshotTownRefreshSkipEnabled() {
		t.Fatal("legacy env should disable post-snapshot town refresh skip")
	}
	if partySeriaActorNoOp23Enabled() {
		t.Fatal("legacy env should disable seria actor no-op23 projection")
	}
}

func TestCharacterFriendRealtimePresenceSwitchAndWatchers(t *testing.T) {
	if !characterFriendRealtimePresenceEnabled() {
		t.Fatal("character friend realtime presence should be enabled by default")
	}
	t.Setenv(characterFriendRealtimePresenceEnv, "legacy")
	if characterFriendRealtimePresenceEnabled() {
		t.Fatal("legacy env should disable character friend realtime presence")
	}
	t.Setenv(characterFriendRealtimePresenceEnv, "")

	s := &Service{
		options:      options{currentClientVersion: "official-24.9"},
		gameSessions: map[uint16]*gameSession{},
	}
	viewer := &gameSession{selectedCharacterID: 101, characterGeneration: 1}
	s.gameSessions[101] = viewer
	identity := boundGameSessionCharacter{session: viewer, character: 101, generation: 1}
	s.recordCurrentCharacterFriendWatchersAndPrime(context.Background(), viewer, identity, []dnfrepo.CharacterRecord{
		{CharacterID: "202", AccountID: "acc2", Name: "PeerB"},
	}, "test")
	hub := s.characterFriendWatcherHub(false)
	if hub == nil {
		t.Fatal("expected watcher hub to be initialized")
	}
	hub.mu.Lock()
	_, watching := hub.byTarget[202][viewer]
	hub.mu.Unlock()
	if !watching {
		t.Fatal("expected viewer to be registered as watcher of character 202")
	}
	s.unregisterCurrentCharacterFriendWatcher(viewer)
	hub.mu.Lock()
	remaining := len(hub.byTarget[202])
	hub.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("expected watcher cleanup after unregister, got %d", remaining)
	}
}

func TestPartyChatNoti12BubbleSwitchAndPayload(t *testing.T) {
	if !partyChatNoti12BubbleEnabled() {
		t.Fatal("party chat NOTI12 bubble should be enabled by default")
	}
	op, body := currentPartyChatNoticePayload(0, 29, []byte("Peer29"), []byte("hi"))
	if op != currentAreaChatNotifyOpcode {
		t.Fatalf("expected NOTI 12, got %d", op)
	}
	if len(body) < 4 || body[0] != currentPartyChatType || body[1] != 29 || body[2] != 0 || body[3] != 0 {
		t.Fatalf("unexpected NOTI 12 party chat header: %x", body)
	}
	t.Setenv(partyChatNoti12BubbleEnv, "legacy")
	if partyChatNoti12BubbleEnabled() {
		t.Fatal("legacy env should disable party chat NOTI12 bubble")
	}
	opLegacy, _ := currentPartyChatNoticePayload(0, 29, []byte("Peer29"), []byte("hi"))
	if opLegacy != currentOfficialWhisperNotifyOpcode {
		t.Fatalf("expected legacy NOTI 65, got %d", opLegacy)
	}
}

func TestMainWhisperZeroTargetIDAndMegaphoneHyperlink(t *testing.T) {
	// 1. Main whisper (type=1) with TargetID=0 and trailing server_group=0 must decode.
	var w packetWriter
	w.writeByte(1)
	w.writeUint16(0)
	w.writeUint32(0)
	w.writeRawDstr([]byte("hello"))
	w.writeRawDstr([]byte("PeerB"))
	w.writeByte(0)
	req, err := decodeCurrentPeerChatRequest(w.bytes())
	if err != nil {
		t.Fatalf("expected type=1 TargetID=0 main whisper to decode, got %v", err)
	}
	if req.MessageType != 1 || req.TargetID != 0 || !req.HasTrailingValue || req.TrailingValue != 0 {
		t.Fatalf("unexpected decoded whisper: %+v", req)
	}

	// 2. Megaphone with item hyperlink (op424 mode=11, slot=7, 1 item) -> NOTI 372 layout.
	var hw packetWriter
	hw.writeByte(11)
	hw.writeUint16(0)
	hw.writeUint32(7)
	hw.writeRawDstr([]byte("2[item]"))
	hw.writeByte(1)
	item := make([]byte, currentItemHyperlinkItemSize)
	item[0] = 0xAA
	hw.writeBytes(item)
	hreq, err := decodeCurrentMegaphoneItemHyperlinkRequest(hw.bytes())
	if err != nil {
		t.Fatalf("expected megaphone hyperlink to decode, got %v", err)
	}
	if hreq.MessageType != 11 || hreq.Reserved != 7 || len(hreq.Items) != 1 {
		t.Fatalf("unexpected decoded megaphone hyperlink: %+v", hreq)
	}
	noti := buildCurrentMegaphoneItemHyperlinkBody(11, 0, 141, 29, []byte("Peer29"), hreq.Message, hreq.Items)
	if len(noti) < 6 || noti[0] != 11 || noti[3] != 141 {
		t.Fatalf("unexpected NOTI 372 header: %x", noti[:6])
	}
}


func TestWhisperAndPartyItemHyperlinkDecodeAndEncode(t *testing.T) {
	item := make([]byte, currentItemHyperlinkItemSize)
	item[0] = 0x55

	// 1. Main whisper (type=1) with TargetID=0, TargetName, TrailingGroup, and 1 item at the end.
	var ww packetWriter
	ww.writeByte(1)
	ww.writeUint16(0)
	ww.writeUint32(0)
	ww.writeRawDstr([]byte("check [item]"))
	ww.writeRawDstr([]byte("PeerB"))
	ww.writeByte(0)
	ww.writeByte(1)
	ww.writeBytes(item)
	wreq, witems, err := decodeCurrentWhisperItemHyperlinkRequest(ww.bytes())
	if err != nil {
		t.Fatalf("expected whisper hyperlink to decode, got %v", err)
	}
	if wreq.MessageType != 1 || wreq.TargetID != 0 || len(witems) != 1 || string(wreq.TargetName) != "PeerB" {
		t.Fatalf("unexpected whisper hyperlink decode: %+v items=%d", wreq, len(witems))
	}
	noti371 := buildCurrentWhisperItemHyperlinkNoti371Body(1, 0, "PeerA", wreq.Message, witems)
	if len(noti371) < 2 || noti371[0] != 1 || noti371[1] != 0 {
		t.Fatalf("unexpected NOTI 371 whisper header: %x", noti371[:4])
	}

	// 2. Party chat (type=2) with 1 item.
	var pw packetWriter
	pw.writeByte(currentPartyChatType)
	pw.writeUint16(0)
	pw.writeUint32(0)
	pw.writeRawDstr([]byte("party [item]"))
	pw.writeByte(1)
	pw.writeBytes(item)
	preq, err := decodeCurrentPartyItemHyperlinkRequest(pw.bytes())
	if err != nil {
		t.Fatalf("expected party hyperlink to decode, got %v", err)
	}
	if preq.MessageType != currentPartyChatType || len(preq.Items) != 1 {
		t.Fatalf("unexpected party hyperlink decode: %+v", preq)
	}
	noti370 := buildCurrentPartyItemHyperlinkNoti370Body(0, 29, preq.Message, preq.Items)
	if len(noti370) < 4 || noti370[0] != currentPartyChatType || noti370[1] != 29 {
		t.Fatalf("unexpected NOTI 370 party header: %x", noti370[:4])
	}
}
