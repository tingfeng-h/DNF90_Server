package dnfbridge

import "net"

// writeCurrentDedicatedLegFastPath 是稳态常规 P2P 转发的免全局锁快路径：
// 当写入目标是该对端腿专属的独立 socket（conn == opposite.conn，非共享探针口）
// 且无需跨成员观测端点校验时，直接以原子标志校验生命周期后执行 WriteToUDP：
//   - 消除所有 12 条腿/全服所有队伍在每包写 socket 时争抢同一把全局 r.mu.Lock()；
//   - 消除每包调用 SetWriteDeadline 引发的 netpoller 全局定时器堆争用；
//   - 腿淘汰时 closeLeg 会置 retired=true 并 Close 专属 socket，不会跨队伍误投；
//   - 共享探针口及 observed 探测路径仍走原加锁慢路径，零改动包体与顺序。
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
		return 0, net.ErrClosed, true
	}
	if leg.key.pvp != opposite.key.pvp || leg.key.partyID != opposite.key.partyID ||
		leg.key.ownerID != opposite.key.peerID || leg.key.peerID != opposite.key.ownerID {
		return 0, net.ErrClosed, true
	}
	n, err := conn.WriteToUDP(payload, target)
	return n, err, true
}
