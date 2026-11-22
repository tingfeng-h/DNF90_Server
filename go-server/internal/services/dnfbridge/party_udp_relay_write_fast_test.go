package dnfbridge

import (
	"net"
	"testing"
	"time"
)

func TestPartyUDPRelayWriteFastPathAfterSlowPathDeadlineExpiry(t *testing.T) {
	receiverConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen receiver: %v", err)
	}
	defer receiverConn.Close()

	legConnA, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen leg A: %v", err)
	}
	defer legConnA.Close()

	legConnB, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen leg B: %v", err)
	}
	defer legConnB.Close()

	keyA := currentPartyUDPRelayLegKey{partyID: 1, ownerID: 13, peerID: 1}
	keyB := currentPartyUDPRelayLegKey{partyID: 1, ownerID: 1, peerID: 13}
	legA := &currentPartyUDPRelayLeg{key: keyA, conn: legConnA}
	legB := &currentPartyUDPRelayLeg{key: keyB, conn: legConnB}

	relay := &currentPartyUDPRelay{
		legs: map[currentPartyUDPRelayLegKey]*currentPartyUDPRelayLeg{
			keyA: legA,
			keyB: legB,
		},
	}
	target := receiverConn.LocalAddr().(*net.UDPAddr)

	// 1. 模拟首包走慢路径（probe/fanout），慢路径会设置 50ms 写截止并在返回前清除。
	if _, err := relay.writeCurrentPairDatagram(legA, legB, legB.conn, []byte("probe"), target, true); err != nil {
		t.Fatalf("slow path write failed: %v", err)
	}

	// 2. 等待超过 50ms 写截止窗口（复现 v15-c20 的 177ms 后全员红电脑场景）。
	time.Sleep(currentPartyRelayWriteTimeout + 25*time.Millisecond)

	// 3. 常规快路径必须直接成功，绝不出现残留的 i/o timeout。
	if n, err := relay.writeCurrentPairDatagram(legA, legB, legB.conn, []byte("steady"), target); err != nil || n != 6 {
		t.Fatalf("fast path write after slow path deadline expiry failed: n=%d err=%v", n, err)
	}

	// 4. 退役/配对不匹配/共享 socket（conn != opposite.conn）均必须返回 handled=false 回退慢路径。
	if _, _, handled := relay.writeCurrentDedicatedLegFastPath(legA, legB, legConnA, []byte("x"), target); handled {
		t.Fatalf("shared/mismatched conn must fall back to slow path")
	}
	legB.retired.Store(true)
	if _, _, handled := relay.writeCurrentDedicatedLegFastPath(legA, legB, legB.conn, []byte("x"), target); handled {
		t.Fatalf("retired leg must fall back to slow path")
	}
	legB.retired.Store(false)
	wrongLeg := &currentPartyUDPRelayLeg{key: currentPartyUDPRelayLegKey{partyID: 1, ownerID: 99, peerID: 1}, conn: legConnB}
	if _, _, handled := relay.writeCurrentDedicatedLegFastPath(legA, wrongLeg, wrongLeg.conn, []byte("x"), target); handled {
		t.Fatalf("mismatched pair key must fall back to slow path")
	}
}
