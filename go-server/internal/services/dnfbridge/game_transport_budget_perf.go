package dnfbridge

import (
	"errors"
	"io"
	"time"
)

// writeGameSessionFastObserved 是常规单包 TCP 写出的零额外分配快路径（调用方已持 wireMu）。
// 当不存在自定义父 wireWriteContext 时，直接复用连接自身的 transportCtx 校验与单次
// SetWriteDeadline，省去每包分配 context.WithTimeout、chan struct{} 与 context.AfterFunc
// 闭包带来的 GC 压力与全局定时器堆争用。连接被主动关闭时 closeGameTransport() 会直接
// Close 底层 conn 立即唤醒阻塞写，保持完全一致的 30s 超时与即时取消语义。
func (s *Service) writeGameSessionFastObserved(p *gameSession, data []byte, stage int) (int, error) {
	writeProbeStarted := time.Now()
	defer s.probeGameWriteDuration(p, stage, len(data), writeProbeStarted)
	lifetime := p.gameTransportContext()
	if err := lifetime.Err(); err != nil {
		s.logFirstGameTransportFailure(p, "before_socket_write", stage, len(data), 0, 0, lifetime)
		p.closeGameTransport()
		return 0, err
	}
	writeStarted := time.Now()
	deadline := writeStarted.Add(gameTransportWriteTimeout)
	if err := p.conn.SetWriteDeadline(deadline); err != nil {
		s.logFirstGameTransportFailure(p, "set_write_deadline", stage, len(data), 0, 0, lifetime)
		p.closeGameTransport()
		return 0, err
	}
	n, err := s.writeGameConnectionObserved(p.conn, data, stage)
	if clearErr := p.conn.SetWriteDeadline(time.Time{}); clearErr != nil {
		s.logFirstGameTransportFailure(p, "clear_write_deadline", stage, len(data), 0, 0, lifetime)
		p.closeGameTransport()
	}
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		s.logFirstGameTransportFailure(p, "socket_write", stage, len(data), n, time.Since(writeStarted), lifetime)
		if p.dproto != nil || n > 0 || err == io.ErrShortWrite || lifetime.Err() != nil || !time.Now().Before(deadline) {
			p.closeGameTransport()
		}
		if lifetime.Err() != nil {
			return n, errors.Join(err, lifetime.Err())
		}
		return n, err
	}
	return n, nil
}
