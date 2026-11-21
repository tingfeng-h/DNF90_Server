package dnfbridge

import (
	"net"
	"time"
)

const (
	// partyUDPRelaySocketBufferBytes 为每条组队 UDP 中继腿分配 1 MiB 内核读写缓冲，
	// 避免 4 人全互联（12 条有向腿）在副本高帧率技能爆发时触发内核 UDP 丢包与卡顿。
	partyUDPRelaySocketBufferBytes = 1 << 20
)

// configurePartyUDPRelaySocket 扩大中继 UDP 套接字的内核收发缓冲；失败时静默降级不阻断建腿。
func configurePartyUDPRelaySocket(conn *net.UDPConn) {
	if conn == nil {
		return
	}
	_ = conn.SetReadBuffer(partyUDPRelaySocketBufferBytes)
	_ = conn.SetWriteBuffer(partyUDPRelaySocketBufferBytes)
}

// configureGameTCPConn 为游戏 TCP 会话显式关闭 Nagle 算法并适当扩大缓冲，降低小包同步抖动。
func configureGameTCPConn(conn net.Conn) {
	tc, ok := conn.(*net.TCPConn)
	if !ok || tc == nil {
		return
	}
	_ = tc.SetNoDelay(true)
}

// partyLANSamePrivateSubnet24 判断两个客户端自报的私网 IPv4 是否位于同一个 RFC1918 /24 局域网段。
// 仅在服务端已核验两端公网观测地址（observed outer）完全相同时才作为同局域网双机直连条件使用，
// 避免不同家庭 NAT 后面的同名私网地址误配。
func partyLANSamePrivateSubnet24(a, b [4]byte) bool {
	if !partyLANIsRFC1918(a) || !partyLANIsRFC1918(b) {
		return false
	}
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

func partyLANIsRFC1918(ip [4]byte) bool {
	if ip[0] == 10 {
		return true
	}
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}
	return ip[0] == 192 && ip[1] == 168
}

// observeInboundFastPath 在热路径端点未变化时无锁判定是否可直接转发，
// 避免每收到一个 UDP 数据报都争抢全局 r.mu 互斥锁。
func (leg *currentPartyUDPRelayLeg) observeInboundFastPath(source *net.UDPAddr, nowUnix int64) (hit bool, refreshMemo bool) {
	if leg == nil || source == nil {
		return false, false
	}
	current := leg.fastTarget.Load()
	if current == nil || !currentPartyUDPRelayAddrEqual(current, source) {
		return false, false
	}
	leg.fastReceived.Add(1)
	leg.fastLastRecvUnixNano.Store(time.Now().UnixNano())
	prev := leg.lastMemoUnixSec.Load()
	if nowUnix > prev && leg.lastMemoUnixSec.CompareAndSwap(prev, nowUnix) {
		refreshMemo = true
	}
	return true, refreshMemo
}

// linkOppositeLegsLocked 在持有 r.mu 时把同一对有向腿的 opposite 原子指针互相绑定，
// 使 readLoop 热路径查找对向腿完全无锁。
func (r *currentPartyUDPRelay) linkOppositeLegsLocked(keys []currentPartyUDPRelayLegKey) {
	for _, key := range keys {
		leg := r.legs[key]
		if leg == nil {
			continue
		}
		oppKey := currentPartyUDPRelayLegKey{
			pvp:     key.pvp,
			partyID: key.partyID,
			ownerID: key.peerID,
			peerID:  key.ownerID,
		}
		opp := r.legs[oppKey]
		leg.opposite.Store(opp)
		if opp != nil {
			opp.opposite.Store(leg)
		}
	}
}

// retireLegsLocked 在持有 r.mu 时标记待关闭腿已退役并解绑双向 opposite 指针。
func (r *currentPartyUDPRelay) retireLegsLocked(legs []*currentPartyUDPRelayLeg) {
	for _, leg := range legs {
		if leg == nil {
			continue
		}
		leg.retired.Store(true)
		if opp := leg.opposite.Swap(nil); opp != nil {
			opp.opposite.CompareAndSwap(leg, nil)
		}
	}
}

// totalReceivedLocked 返回包含无锁快路径计数在内的累计收包数（调用方持 leg.mu）。
func (leg *currentPartyUDPRelayLeg) totalReceivedLocked() uint64 {
	if leg == nil {
		return 0
	}
	return leg.received + leg.fastReceived.Load()
}

// lastReceivedTimeLocked 返回包含无锁快路径时间戳在内的最近收包时刻（调用方持 leg.mu）。
func (leg *currentPartyUDPRelayLeg) lastReceivedTimeLocked() time.Time {
	if leg == nil {
		return time.Time{}
	}
	t := leg.lastReceivedAt
	if ns := leg.fastLastRecvUnixNano.Load(); ns > 0 {
		ft := time.Unix(0, ns)
		if ft.After(t) {
			return ft
		}
	}
	return t
}
