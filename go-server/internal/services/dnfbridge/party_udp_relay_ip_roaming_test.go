package dnfbridge

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestCurrentPartyUDPRelayIPRoamingAllowsSplitEgressAndHandshakeSwitch(t *testing.T) {
	t.Setenv(partyUDPRelayIPRoamingEnv, "roam")
	t.Setenv(partyUDPRelayRoamGraceMsEnv, "5000")
	t.Setenv(partyUDPRelayRoamIdleMsEnv, "500")

	portStart := testCurrentPartyUDPRelayPortStart(t)
	relay, err := newCurrentPartyUDPRelay(true, "127.0.0.1", portStart, 12, nil)
	if err != nil {
		t.Fatalf("new UDP relay: %v", err)
	}
	defer relay.Close()

	// Both TCP sessions observed 127.0.0.1, while A's UDP egress comes from 127.0.0.2
	// and then switches immediately to 127.0.0.3 during the handshake grace window.
	bindings := []currentPartyUDPRelayBinding{
		{characterID: 101, generation: 1, address: [4]byte{127, 0, 0, 1}},
		{characterID: 202, generation: 1, address: [4]byte{127, 0, 0, 1}},
	}
	if !relay.Sync(88, bindings) {
		t.Fatal("relay did not create room")
	}
	addrAB, portAB, ok := relay.Endpoint(88, 101, 202)
	if !ok {
		t.Fatal("missing AB endpoint")
	}
	addrBA, portBA, ok := relay.Endpoint(88, 202, 101)
	if !ok {
		t.Fatal("missing BA endpoint")
	}

	clientA1 := testCurrentPartyUDPClient(t, "127.0.0.2")
	defer clientA1.Close()
	clientA2 := testCurrentPartyUDPClient(t, "127.0.0.3")
	defer clientA2.Close()
	clientB := testCurrentPartyUDPClient(t, "127.0.0.1")
	defer clientB.Close()

	epAB := &net.UDPAddr{IP: net.IPv4(addrAB[0], addrAB[1], addrAB[2], addrAB[3]), Port: int(portAB)}
	epBA := &net.UDPAddr{IP: net.IPv4(addrBA[0], addrBA[1], addrBA[2], addrBA[3]), Port: int(portBA)}

	buf := make([]byte, 128)
	deadline := time.Now().Add(4 * time.Second)
	for {
		_, _ = clientB.WriteToUDP([]byte("b-hi"), epBA)
		_, _ = clientA1.WriteToUDP([]byte("from-A1"), epAB)
		_ = clientB.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, _, err := clientB.ReadFromUDP(buf)
		if err == nil && bytes.Equal(buf[:n], []byte("from-A1")) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("did not receive split-egress payload from A1: %v", err)
		}
	}

	// Immediate A1 -> A2 switch inside handshake grace window (no 2s stall).
	deadline = time.Now().Add(4 * time.Second)
	for {
		_, _ = clientA2.WriteToUDP([]byte("from-A2"), epAB)
		_ = clientB.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, _, err := clientB.ReadFromUDP(buf)
		if err == nil && bytes.Equal(buf[:n], []byte("from-A2")) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("did not receive switched payload from A2 during grace window: %v", err)
		}
	}

	// Verify memo recorded actual UDP egress IP (127.0.0.3), not TCP IP (127.0.0.1).
	memo := relay.lookupEndpointMemo(101)
	if memo == nil || !memo.valid {
		t.Fatalf("expected valid endpoint memo for char 101")
	}
	if memo.address != ([4]byte{127, 0, 0, 3}) {
		t.Fatalf("memo.address = %v, want [127 0 0 3]", memo.address)
	}
}

func TestCurrentPartyUDPRelayIPRoamingSteadyStateAntiHijackGuard(t *testing.T) {
	t.Setenv(partyUDPRelayIPRoamingEnv, "roam")
	t.Setenv(partyUDPRelayRoamGraceMsEnv, "0")
	t.Setenv(partyUDPRelayRoamIdleMsEnv, "500")

	now := time.Now()
	leg := &currentPartyUDPRelayLeg{
		binding:         currentPartyUDPRelayBinding{characterID: 101, address: [4]byte{10, 0, 0, 1}},
		createdAt:       now.Add(-10 * time.Second),
		learnedEndpoint: &net.UDPAddr{IP: net.IPv4(10, 0, 0, 2), Port: 40000},
		lastReceivedAt:  now.Add(-100 * time.Millisecond),
	}
	leg.fastTarget.Store(leg.learnedEndpoint)

	// Active learned IP (10.0.0.2) is allowed.
	if !currentPartyUDPRelaySourceAllowed(leg, &net.UDPAddr{IP: net.IPv4(10, 0, 0, 2), Port: 40001}, now) {
		t.Fatal("active learned IP must be allowed")
	}
	// Third-party IP (10.0.0.9) is rejected while 10.0.0.2 was active 100ms ago (< 500ms).
	if currentPartyUDPRelaySourceAllowed(leg, &net.UDPAddr{IP: net.IPv4(10, 0, 0, 9), Port: 40002}, now) {
		t.Fatal("third-party IP must be rejected during active steady-state traffic")
	}
	// Once 10.0.0.2 is silent for >= 500ms, roaming to 10.0.0.9 is allowed.
	if !currentPartyUDPRelaySourceAllowed(leg, &net.UDPAddr{IP: net.IPv4(10, 0, 0, 9), Port: 40002}, now.Add(600*time.Millisecond)) {
		t.Fatal("roaming should be allowed after old IP is silent >= idle threshold")
	}
}
