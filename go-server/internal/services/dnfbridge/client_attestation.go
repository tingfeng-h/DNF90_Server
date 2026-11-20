package dnfbridge

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// 客户端文件指纹校验（登录器登录时提交，服务端按允许清单比对）。
//
// 设计要点：
//   - 允许清单是运行时可热更新的 JSON（默认 /runtime/client-attest.json），按 mtime 重载；
//   - 挑战 nonce 一次性、两分钟有效、绑定发起 IP，防止重放；
//   - DNFBRIDGE_ATTEST_SECRET 非空时要求登录器对 (账号|nonce|三个哈希) 做 HMAC-SHA256；
//   - 清单里某项写 "*" 表示不校验该项（便于先只锁定 exe / pvf，登录器稍后再上线）。
//
// 边界：指纹由登录器本地计算，篡改登录器可以伪造；这是"统一版本、提高门槛"，不是
// 密码学意义上的强证明。
const (
	defaultClientAttestFile = "/runtime/client-attest.json"
	clientAttestNonceTTL    = 2 * time.Minute
	clientAttestMaxPending  = 4096
)

const (
	clientAttestModeOff     = "off"
	clientAttestModeWarn    = "warn"
	clientAttestModeEnforce = "enforce"
)

// clientAttestVersion 是一条被允许的客户端组合。
type clientAttestVersion struct {
	Label    string `json:"label"`
	Exe      string `json:"exe"`
	Pvf      string `json:"pvf"`
	Launcher string `json:"launcher"`
	Note     string `json:"note,omitempty"`
}

// clientAttestFile 是允许清单文件内容。
type clientAttestFile struct {
	Mode     string                `json:"mode"`
	Versions []clientAttestVersion `json:"versions"`
}

type clientAttestChallenge struct {
	expiresAt time.Time
	ip        string
}

// clientAttestStore 缓存允许清单与未过期的挑战。
type clientAttestStore struct {
	mu         sync.Mutex
	path       string
	cfg        clientAttestFile
	cfgLoaded  bool
	cfgModTime time.Time
	cfgSize    int64
	cfgError   string
	challenges map[string]clientAttestChallenge
}

// clientAttestOutcome 是一次校验的结论。
type clientAttestOutcome struct {
	Mode   string
	Label  string
	Reason string
	Allow  bool
	// ObserveOnly 表示允许清单还是空的：放行并把完整指纹写日志，便于首次收集。
	ObserveOnly bool
}

func newClientAttestStore() *clientAttestStore {
	path := strings.TrimSpace(os.Getenv("DNFBRIDGE_CLIENT_ATTEST_FILE"))
	if path == "" {
		path = defaultClientAttestFile
	}
	return &clientAttestStore{path: path, challenges: make(map[string]clientAttestChallenge)}
}

func clientAttestHMACSecret() string {
	return strings.TrimSpace(os.Getenv("DNFBRIDGE_ATTEST_SECRET"))
}

func normalizeClientAttestMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case clientAttestModeOff:
		return clientAttestModeOff
	case clientAttestModeEnforce:
		return clientAttestModeEnforce
	case "", clientAttestModeWarn:
		return clientAttestModeWarn
	}
	// 未知取值按 warn 处理：宁可只记录，也不要因为写错一个字母就全服登不上。
	return clientAttestModeWarn
}

func normalizeClientAttestHash(value string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(value))
	switch {
	case v == "":
		return "", false
	case v == "*":
		return "*", true
	case len(v) != 64:
		return "", false
	}
	if _, err := hex.DecodeString(v); err != nil {
		return "", false
	}
	return v, true
}

func normalizeClientAttestVersions(in []clientAttestVersion) []clientAttestVersion {
	out := make([]clientAttestVersion, 0, len(in))
	for _, version := range in {
		exe, okExe := normalizeClientAttestHash(version.Exe)
		pvf, okPvf := normalizeClientAttestHash(version.Pvf)
		launcher, okLauncher := normalizeClientAttestHash(version.Launcher)
		if !okExe || !okPvf || !okLauncher {
			// 写坏的条目直接丢弃，避免"空值等于放行"。
			continue
		}
		out = append(out, clientAttestVersion{
			Label:    strings.TrimSpace(version.Label),
			Exe:      exe,
			Pvf:      pvf,
			Launcher: launcher,
			Note:     version.Note,
		})
	}
	return out
}

func clientAttestHashMatches(expected, actual string) bool {
	if expected == "" || expected == "*" {
		return true
	}
	return expected == actual
}

// config 返回当前允许清单；文件变化时重新加载，解析失败保留上一次成功结果。
func (st *clientAttestStore) config() (clientAttestFile, string) {
	if st == nil {
		return clientAttestFile{Mode: clientAttestModeOff}, ""
	}
	info, err := os.Stat(st.path)
	st.mu.Lock()
	defer st.mu.Unlock()
	if err != nil {
		if st.cfgLoaded {
			return st.cfg, st.cfgError
		}
		st.cfgError = fmt.Sprintf("客户端允许清单不可读: %v", err)
		return clientAttestFile{Mode: clientAttestModeOff}, st.cfgError
	}
	if st.cfgLoaded && info.ModTime().Equal(st.cfgModTime) && info.Size() == st.cfgSize {
		return st.cfg, st.cfgError
	}
	raw, readErr := os.ReadFile(st.path)
	if readErr != nil {
		st.cfgError = fmt.Sprintf("客户端允许清单读取失败: %v", readErr)
		if st.cfgLoaded {
			return st.cfg, st.cfgError
		}
		return clientAttestFile{Mode: clientAttestModeOff}, st.cfgError
	}
	var cfg clientAttestFile
	if err := json.Unmarshal(raw, &cfg); err != nil {
		st.cfgError = fmt.Sprintf("客户端允许清单解析失败: %v", err)
		if st.cfgLoaded {
			return st.cfg, st.cfgError
		}
		return clientAttestFile{Mode: clientAttestModeOff}, st.cfgError
	}
	cfg.Mode = normalizeClientAttestMode(cfg.Mode)
	cfg.Versions = normalizeClientAttestVersions(cfg.Versions)
	st.cfg, st.cfgLoaded = cfg, true
	st.cfgModTime, st.cfgSize = info.ModTime(), info.Size()
	st.cfgError = ""
	return st.cfg, ""
}

// challenge 发一个一次性 nonce。
func (st *clientAttestStore) challenge(ip string) (string, time.Time, bool, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", time.Time{}, false, err
	}
	nonce := hex.EncodeToString(raw[:])
	expiresAt := time.Now().Add(clientAttestNonceTTL)
	st.mu.Lock()
	if len(st.challenges) >= clientAttestMaxPending {
		st.challenges = make(map[string]clientAttestChallenge)
	}
	now := time.Now()
	for key, item := range st.challenges {
		if now.After(item.expiresAt) {
			delete(st.challenges, key)
		}
	}
	st.challenges[nonce] = clientAttestChallenge{expiresAt: expiresAt, ip: ip}
	st.mu.Unlock()
	return nonce, expiresAt, clientAttestHMACSecret() != "", nil
}

// consumeChallenge 取用并作废 nonce。
func (st *clientAttestStore) consumeChallenge(nonce, ip string) bool {
	nonce = strings.ToLower(strings.TrimSpace(nonce))
	if nonce == "" {
		return false
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	item, found := st.challenges[nonce]
	if !found {
		return false
	}
	delete(st.challenges, nonce)
	if time.Now().After(item.expiresAt) {
		return false
	}
	if item.ip != "" && ip != "" && item.ip != ip {
		return false
	}
	return true
}

func clientAttestSignature(secret string, req accountAPIRequest) string {
	payload := strings.Join([]string{
		strings.TrimSpace(req.Account),
		strings.ToLower(strings.TrimSpace(req.AttestNonce)),
		strings.ToLower(strings.TrimSpace(req.ClientExeSHA256)),
		strings.ToLower(strings.TrimSpace(req.ClientPvfSHA256)),
		strings.ToLower(strings.TrimSpace(req.LauncherSHA256)),
	}, "|")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// verify 校验一次登录请求的客户端指纹。
func (st *clientAttestStore) verify(req accountAPIRequest, ip string) clientAttestOutcome {
	cfg, cfgErr := st.config()
	mode := normalizeClientAttestMode(cfg.Mode)
	outcome := clientAttestOutcome{Mode: mode}
	if mode == clientAttestModeOff {
		outcome.Allow = true
		return outcome
	}
	if cfgErr != "" {
		outcome.Reason = cfgErr
		return outcome
	}
	if secret := clientAttestHMACSecret(); secret != "" {
		if !st.consumeChallenge(req.AttestNonce, ip) {
			outcome.Reason = "登录挑战已过期，请重新登录"
			return outcome
		}
		expected := clientAttestSignature(secret, req)
		got := strings.ToLower(strings.TrimSpace(req.AttestHmac))
		if got == "" || !hmac.Equal([]byte(expected), []byte(got)) {
			outcome.Reason = "登录签名无效"
			return outcome
		}
	}
	exe, okExe := normalizeClientAttestHash(req.ClientExeSHA256)
	pvf, okPvf := normalizeClientAttestHash(req.ClientPvfSHA256)
	launcher, okLauncher := normalizeClientAttestHash(req.LauncherSHA256)
	if !okExe || !okPvf || !okLauncher {
		outcome.Reason = "缺少客户端文件指纹"
		return outcome
	}
	if len(cfg.Versions) == 0 {
		// 还没配置允许清单：放行，并把完整指纹交给调用方记日志。
		outcome.Allow, outcome.ObserveOnly = true, true
		return outcome
	}
	for _, version := range cfg.Versions {
		if clientAttestHashMatches(version.Exe, exe) &&
			clientAttestHashMatches(version.Pvf, pvf) &&
			clientAttestHashMatches(version.Launcher, launcher) {
			outcome.Allow, outcome.Label = true, version.Label
			return outcome
		}
	}
	outcome.Reason = "客户端文件与允许清单不一致"
	return outcome
}

func shortClientAttestHash(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	if len(v) > 12 {
		return v[:12]
	}
	return v
}
