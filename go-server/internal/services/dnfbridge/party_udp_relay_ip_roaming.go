package dnfbridge

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	partyUDPRelayIPRoamingEnv      = "DNFBRIDGE_PARTY_UDP_RELAY_IP_ROAMING"
	partyUDPRelayRoamIdleMsEnv     = "DNFBRIDGE_PARTY_UDP_RELAY_ROAM_IDLE_MS"
	partyUDPRelayRoamGraceMsEnv    = "DNFBRIDGE_PARTY_UDP_RELAY_ROAM_GRACE_MS"
	defaultPartyUDPRelayRoamIdle   = 500 * time.Millisecond
	defaultPartyUDPRelayRoamGrace  = 5 * time.Second
)

type partyUDPRelayIPRoamingMode uint8

const (
	partyUDPRelayIPRoamingStrict partyUDPRelayIPRoamingMode = iota
	partyUDPRelayIPRoamingRoam
	partyUDPRelayIPRoamingAny
)

// currentPartyUDPRelayIPRoamingMode parses DNFBRIDGE_PARTY_UDP_RELAY_IP_ROAMING:
//   - strict / off / 0 / false / "" (default): legacy behavior (source.IP == TCP IPv4),
//     preserving all existing source_mismatch unit tests unchanged.
//   - roam / on / 1 / true / yes / allow: Scheme 1 — allow TCP IP, current learned UDP IP,
//     seed IP, first-packet split IP, immediate A->B rebind during the initial handshake grace
//     window (default 5s), and steady-state roaming after old-IP silence >= idle threshold
//     (default 500ms, below the client's 1000ms red-computer timeout).
//   - any / last-seen / lsw: allow any valid IPv4 UDP source (Last-Seen-Wins, 0ms idle).
func currentPartyUDPRelayIPRoamingMode() partyUDPRelayIPRoamingMode {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyUDPRelayIPRoamingEnv))) {
	case "1", "on", "true", "yes", "roam", "allow":
		return partyUDPRelayIPRoamingRoam
	case "any", "last-seen", "lsw":
		return partyUDPRelayIPRoamingAny
	default:
		return partyUDPRelayIPRoamingStrict
	}
}

func partyUDPRelayIPRoamingEnabled() bool {
	return currentPartyUDPRelayIPRoamingMode() != partyUDPRelayIPRoamingStrict
}

func partyUDPRelayRoamIdleThreshold() time.Duration {
	raw := strings.TrimSpace(os.Getenv(partyUDPRelayRoamIdleMsEnv))
	if raw == "" {
		return defaultPartyUDPRelayRoamIdle
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return defaultPartyUDPRelayRoamIdle
	}
	return time.Duration(ms) * time.Millisecond
}

func partyUDPRelayRoamGraceWindow() time.Duration {
	raw := strings.TrimSpace(os.Getenv(partyUDPRelayRoamGraceMsEnv))
	if raw == "" {
		return defaultPartyUDPRelayRoamGrace
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return defaultPartyUDPRelayRoamGrace
	}
	return time.Duration(ms) * time.Millisecond
}

// currentPartyUDPRelaySourceAllowed decides whether an inbound UDP datagram
// from `source` should be accepted on `leg`.
func currentPartyUDPRelaySourceAllowed(leg *currentPartyUDPRelayLeg, source *net.UDPAddr, now time.Time) bool {
	if leg == nil || source == nil {
		return false
	}
	if currentPartyUDPRelaySourceMatches(source, leg.binding.address) {
		return true
	}
	mode := currentPartyUDPRelayIPRoamingMode()
	if mode == partyUDPRelayIPRoamingStrict {
		return false
	}
	ip4 := source.IP.To4()
	if ip4 == nil || source.Port <= 0 {
		return false
	}
	if mode == partyUDPRelayIPRoamingAny {
		return true
	}
	// Fast lock-free check against the currently active learned UDP IP.
	if fast := leg.fastTarget.Load(); fast != nil && fast.IP != nil && fast.IP.Equal(ip4) {
		return true
	}
	leg.mu.Lock()
	defer leg.mu.Unlock()
	if leg.learnedEndpoint != nil && leg.learnedEndpoint.IP != nil && leg.learnedEndpoint.IP.Equal(ip4) {
		return true
	}
	if leg.seedEndpoint != nil && leg.seedEndpoint.IP != nil && leg.seedEndpoint.IP.Equal(ip4) {
		return true
	}
	// First-packet split egress (e.g. CGNAT pool / accelerator TCP-UDP split):
	// no UDP endpoint has been learned on this dedicated leg yet.
	if leg.learnedEndpoint == nil {
		return true
	}
	// Handshake grace window: allow immediate A -> B switch right after leg
	// creation (e.g. initial direct probe followed by accelerator takeover).
	if grace := partyUDPRelayRoamGraceWindow(); grace > 0 && !leg.createdAt.IsZero() && now.Sub(leg.createdAt) < grace {
		return true
	}
	// Steady-state anti-hijack guard: while the previously learned IP is still
	// actively sending packets, reject third-party IPs. Allow smooth roaming
	// once the old IP has been silent for at least the idle threshold.
	lastRecv := leg.lastReceivedAt
	if fastNs := leg.fastLastRecvUnixNano.Load(); fastNs > 0 {
		if fastTime := time.Unix(0, fastNs); fastTime.After(lastRecv) {
			lastRecv = fastTime
		}
	}
	if lastRecv.IsZero() {
		return true
	}
	return now.Sub(lastRecv) >= partyUDPRelayRoamIdleThreshold()
}

// currentPartyUDPRelayMemoAddressForSource returns the IPv4 to persist in
// endpointMemo: in strict mode it stays leg.binding.address (exact legacy
// behavior); in roam/any mode it records the actual learned UDP source IPv4.
func currentPartyUDPRelayMemoAddressForSource(leg *currentPartyUDPRelayLeg, source *net.UDPAddr) [4]byte {
	if leg == nil {
		return [4]byte{}
	}
	if partyUDPRelayIPRoamingEnabled() && source != nil && source.IP != nil {
		if ip4 := source.IP.To4(); len(ip4) == net.IPv4len {
			return [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]}
		}
	}
	return leg.binding.address
}

// currentPartyUDPRelayMemoBootstrapAddr resolves a fresh endpointMemo target.
// In strict mode it requires memo.address == leg.binding.address; in roam/any
// mode it reuses the actual learned UDP IPv4:port even when it differs from TCP.
func currentPartyUDPRelayMemoBootstrapAddr(leg *currentPartyUDPRelayLeg, memo *currentPartyUDPRelayEndpointMemo) (*net.UDPAddr, bool) {
	if leg == nil || memo == nil || !memo.valid || memo.port == 0 {
		return nil, false
	}
	if !partyUDPRelayIPRoamingEnabled() {
		if port, ok := memo.portForTTL(leg.binding.address, leg.memoTTL); ok {
			return &net.UDPAddr{
				IP:   net.IPv4(leg.binding.address[0], leg.binding.address[1], leg.binding.address[2], leg.binding.address[3]),
				Port: int(port),
			}, true
		}
		return nil, false
	}
	ttl := leg.memoTTL
	if ttl <= 0 || ttl > memoTTL {
		ttl = memoTTL
	}
	if time.Since(memo.seenAt) > ttl {
		return nil, false
	}
	if memo.address == ([4]byte{}) {
		return nil, false
	}
	return &net.UDPAddr{
		IP:   net.IPv4(memo.address[0], memo.address[1], memo.address[2], memo.address[3]),
		Port: int(memo.port),
	}, true
}
