package dnfbridge

import (
	"time"
)

// probeGameWriteDuration 记录"实际写 socket"的耗时。
// 只在超过 20ms 时告警，并每 60 秒最多输出一次汇总，避免刷日志。
func (s *Service) probeGameWriteDuration(p *gameSession, stage int, size int, started time.Time) {
	if s == nil || p == nil || started.IsZero() {
		return
	}
	elapsed := time.Since(started)
	if elapsed < 20*time.Millisecond {
		return
	}
	s.logWarn("game socket write slow",
		"conn_id", p.connID,
		"char_id", p.selectedCharacterID,
		"stage", stage,
		"bytes", size,
		"duration_ms", elapsed.Milliseconds())
}
