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
