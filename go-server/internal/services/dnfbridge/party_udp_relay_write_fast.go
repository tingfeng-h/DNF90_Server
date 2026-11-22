package dnfbridge

import "net"

// writeCurrentDedicatedLegFastPath 是稳态常规 P2P 转发的免全局锁快路径：
//   - 仅在专属腿 socket（conn == opposite.conn）、双方未退役、配对键完全匹配时启用；
//   - 任何前置条件不满足或 WriteToUDP 返回错误时一律返回 handled=false，回退到原加锁慢路径
//     （慢路径会重新校验在册状态并刷新 50ms SetWriteDeadline，不会吞掉异常或误判 ErrClosed）；
//   - 慢路径写完后已统一清除写截止（SetWriteDeadline(time.Time{})），彻底杜绝首包 probe 残留的
//     50ms 过期写截止导致后续免锁写持续报 i/o timeout（全员红电脑）。
func (r *currentPartyUDPRelay) writeCurrentDedicatedLegFastPath(
	leg, opposite *currentPartyUDPRelayLeg,
	conn *net.UDPConn,
	payload []byte,
	target *net.UDPAddr,
) (int, error, bool) {
	if r == nil || leg == nil || opposite == nil || conn == nil || target == nil || conn != opposite.conn {
		return 0, nil, false
	}
	if r.closedFlag.Load() || leg.retired.Load() || opposite.retired.Load() {
		return 0, nil, false
	}
	if leg.key.pvp != opposite.key.pvp || leg.key.partyID != opposite.key.partyID ||
		leg.key.ownerID != opposite.key.peerID || leg.key.peerID != opposite.key.ownerID {
		return 0, nil, false
	}
	n, err := conn.WriteToUDP(payload, target)
	if err != nil {
		return 0, nil, false
	}
	return n, nil, true
}
