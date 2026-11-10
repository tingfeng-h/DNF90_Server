package dnfbridge

import (
	"net"
	"testing"
	"time"
)

// 成员重入/新腿建立时，优先用同角色在别的腿上已经学到的**真实 P2P 端点**，
// 而不是持续保活的探针口。现场：探针源端口 6969、游戏 P2P 源端口 8278；
// 旧顺序让新腿把种子写到探针口，包进不了游戏 socket，老成员的新腿永远不开口。
func TestSeedEndpointPrefersFreshLearnedP2PEndpointOverProbe(t *testing.T) {
	recorder := &relayLogRecorder{}
	relay, err := newCurrentPartyUDPRelayWithSettings(currentPartyUDPRelaySettings{
		enabled:     true,
		advertiseIP: "127.0.0.1",
		portStart:   testCurrentPartyUDPRelayPortStart(t),
		portCount:   12,
		probeReplay: true,
		seedForward: true,
		memoTTL:     30 * time.Second,
		logWarn:     recorder.log,
	})
	if err != nil {
		t.Fatalf("new UDP relay: %v", err)
	}
	t.Cleanup(func() { _ = relay.Close() })

	address := [4]byte{127, 0, 0, 1}
	relay.rememberEndpoint(&currentPartyUDPRelayLeg{
		key:     currentPartyUDPRelayLegKey{partyID: 1, ownerID: 101},
		binding: currentPartyUDPRelayBinding{characterID: 101, address: address},
	}, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 23456})
	relay.mu.Lock()
	relay.probes[101] = currentPartyUDPRelayProbeEndpoint{
		characterID: 101,
		address:     address,
		port:        34567,
		hitPort:     relay.portStart,
		seenAt:      time.Now(),
	}
	relay.mu.Unlock()

	seed := relay.seedEndpointFor(101, address)
	if seed == nil || seed.Port != 23456 {
		t.Fatalf("seed endpoint = %v, want the fresh learned P2P port 23456 (probe port must not win)", seed)
	}
}

// 备忘过期时仍回落到探针口：探针是“客户端现在还在发”的唯一证据。
func TestSeedEndpointFallsBackToProbeWhenLearnedEndpointIsStale(t *testing.T) {
	recorder := &relayLogRecorder{}
	relay, err := newCurrentPartyUDPRelayWithSettings(currentPartyUDPRelaySettings{
		enabled:     true,
		advertiseIP: "127.0.0.1",
		portStart:   testCurrentPartyUDPRelayPortStart(t),
		portCount:   12,
		probeReplay: true,
		seedForward: true,
		memoTTL:     30 * time.Second,
		logWarn:     recorder.log,
	})
	if err != nil {
		t.Fatalf("new UDP relay: %v", err)
	}
	t.Cleanup(func() { _ = relay.Close() })

	address := [4]byte{127, 0, 0, 1}
	relay.rememberEndpoint(&currentPartyUDPRelayLeg{
		key:     currentPartyUDPRelayLegKey{partyID: 1, ownerID: 102},
		binding: currentPartyUDPRelayBinding{characterID: 102, address: address},
	}, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 23457})
	relay.memoMu.Lock()
	memo := relay.endpointMemo[102]
	memo.seenAt = time.Now().Add(-90 * time.Second)
	relay.endpointMemo[102] = memo
	relay.memoMu.Unlock()
	relay.mu.Lock()
	relay.probes[102] = currentPartyUDPRelayProbeEndpoint{
		characterID: 102,
		address:     address,
		port:        34568,
		hitPort:     relay.portStart,
		seenAt:      time.Now(),
	}
	relay.mu.Unlock()

	seed := relay.seedEndpointFor(102, address)
	if seed == nil || seed.Port != 34568 {
		t.Fatalf("seed endpoint = %v, want fresh probe port 34568 when the learned memo is stale", seed)
	}
}
