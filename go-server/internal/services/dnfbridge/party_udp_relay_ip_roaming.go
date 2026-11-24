package dnfbridge

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	partyUDPRelayIPRoamingEnv     = "DNFBRIDGE_PARTY_UDP_RELAY_IP_ROAMING"
	partyUDPRelayRoamIdleMsEnv    = "DNFBRIDGE_PARTY_UDP_RELAY_ROAM_IDLE_MS"
	defaultPartyUDPRelayRoamIdle  = 2 * time.Second
)

// partyUDPRelayIPRoamingEnabled controls whether party UDP relay legs allow
// UDP egress IP roaming (CGNAT pool split, accelerator TCP/UDP split, multi-WAN,
// or mid-session network switch).
//   - off / strict (default): exact legacy behavior (source.IP == TCP IPv4),
//     preserving all existing source_mismatch unit tests unchanged.
//   - on / roam / 1 / true: allow TCP IP, current learned UDP IP, probe-seeded
//     UDP IP, first-packet split IP (when no endpoint learned yet), and smooth
//     rebind to a new IP after the old IP has been silent for >= roamIdle.
func partyUDPRelayIPRoamingEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyUDPRelayIPRoamingEnv))) {
	case "1", "on", "true", "yes", "roam", "allow":
		return true
	default:
		return false
	}
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

// currentPartyUDPRelaySourceAllowed decides whether an inbound UDP datagram
// from `source` should be accepted on `leg`.
func currentPartyUDPRelaySourceAllowed(leg *currentPartyUDPRelayLeg, source *net.UDPAddr, now time.Time) bool {
	if leg == nil || source == nil {
		return false
	}
	if currentPartyUDPRelaySourceMatches(source, leg.binding.address) {
		return true
	}
	if !partyUDPRelayIPRoamingEnabled() {
		return false
	}
	ip4 := source.IP.To4()
	if ip4 == nil || source.Port <= 0 {
		return false
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
	// Anti-hijack guard: while the previously learned IP is still actively
	// sending packets, reject any third-party IP. Only allow roaming to a new
	// IP once the old IP has gone silent for at least the idle threshold.
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
