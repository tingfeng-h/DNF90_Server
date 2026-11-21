package dnfbridge

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

// startDebugPprof 在 DNFBRIDGE_PPROF_LISTEN 非空时启动一个只读调试端点。
// 只用于排查期（127.0.0.1 监听），生产环境不设置该变量即完全关闭。
func (s *Service) startDebugPprof() {
	addr := strings.TrimSpace(os.Getenv("DNFBRIDGE_PPROF_LISTEN"))
	if addr == "" {
		return
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil || (host != "127.0.0.1" && host != "localhost" && host != "::1") {
		s.logWarn("dnfbridge debug pprof rejected non-loopback address", "addr", addr)
		return
	}
	go func() {
		s.logWarn("dnfbridge debug pprof listening", "addr", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			s.logWarn("dnfbridge debug pprof stopped", "addr", addr, "error", err)
		}
	}()
}
