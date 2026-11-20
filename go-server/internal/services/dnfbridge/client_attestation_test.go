package dnfbridge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeAttestFile(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "client-attest.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("写入允许清单失败: %v", err)
	}
	t.Setenv("DNFBRIDGE_CLIENT_ATTEST_FILE", path)
	t.Setenv("DNFBRIDGE_ATTEST_SECRET", "")
	return path
}

const attestTestHash = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestClientAttestMissingFileDisablesCheck(t *testing.T) {
	t.Setenv("DNFBRIDGE_CLIENT_ATTEST_FILE", filepath.Join(t.TempDir(), "absent.json"))
	store := newClientAttestStore()
	outcome := store.verify(accountAPIRequest{Account: "alice"}, "127.0.0.1")
	if !outcome.Allow || outcome.Mode != clientAttestModeOff {
		t.Fatalf("缺文件应视为关闭: %+v", outcome)
	}
}

func TestClientAttestWarnAndEnforceShareMatching(t *testing.T) {
	writeAttestFile(t, `{"mode":"warn","versions":[{"label":"v211","exe":"`+attestTestHash+`","pvf":"*","launcher":"*"}]}`)
	store := newClientAttestStore()
	matching := accountAPIRequest{Account: "alice", ClientExeSHA256: strings.ToUpper(attestTestHash), ClientPvfSHA256: attestTestHash, LauncherSHA256: attestTestHash}
	if outcome := store.verify(matching, "10.0.0.5"); !outcome.Allow || outcome.Label != "v211" {
		t.Fatalf("大写哈希应命中: %+v", outcome)
	}
	other := accountAPIRequest{Account: "alice", ClientExeSHA256: strings.Repeat("b", 64), ClientPvfSHA256: attestTestHash, LauncherSHA256: attestTestHash}
	outcome := store.verify(other, "10.0.0.5")
	if outcome.Allow || outcome.Mode != clientAttestModeWarn || outcome.Reason == "" {
		t.Fatalf("不匹配应被记录并拒绝: %+v", outcome)
	}
}

func TestClientAttestRejectsMalformedEntry(t *testing.T) {
	writeAttestFile(t, `{"mode":"enforce","versions":[{"label":"broken","exe":"zz","pvf":"*","launcher":"*"}]}`)
	store := newClientAttestStore()
	outcome := store.verify(accountAPIRequest{Account: "alice", ClientExeSHA256: attestTestHash, ClientPvfSHA256: attestTestHash, LauncherSHA256: attestTestHash}, "10.0.0.5")
	if outcome.Allow {
		t.Fatalf("坏条目不应放行任何客户端: %+v", outcome)
	}
}

func TestClientAttestChallengeIsSingleUseAndSigned(t *testing.T) {
	writeAttestFile(t, `{"mode":"enforce","versions":[{"label":"v211","exe":"`+attestTestHash+`","pvf":"*","launcher":"*"}]}`)
	t.Setenv("DNFBRIDGE_ATTEST_SECRET", "unit-test-secret")
	store := newClientAttestStore()

	nonce, _, hmacRequired, err := store.challenge("10.0.0.5")
	if err != nil || !hmacRequired || nonce == "" {
		t.Fatalf("挑战生成失败: %q %v %v", nonce, hmacRequired, err)
	}
	request := accountAPIRequest{Account: "alice", AttestNonce: nonce, ClientExeSHA256: attestTestHash, ClientPvfSHA256: attestTestHash, LauncherSHA256: attestTestHash}
	request.AttestHmac = clientAttestSignature("unit-test-secret", request)
	if outcome := store.verify(request, "10.0.0.5"); !outcome.Allow {
		t.Fatalf("正确签名应放行: %+v", outcome)
	}
	if outcome := store.verify(request, "10.0.0.5"); outcome.Allow {
		t.Fatalf("nonce 应一次性: %+v", outcome)
	}

	nonce, _, _, _ = store.challenge("10.0.0.5")
	bad := accountAPIRequest{Account: "alice", AttestNonce: nonce, AttestHmac: strings.Repeat("0", 64), ClientExeSHA256: attestTestHash, ClientPvfSHA256: attestTestHash, LauncherSHA256: attestTestHash}
	if outcome := store.verify(bad, "10.0.0.5"); outcome.Allow {
		t.Fatalf("错误签名不应放行: %+v", outcome)
	}
}

func TestClientAttestChallengeExpires(t *testing.T) {
	store := newClientAttestStore()
	nonce, _, _, err := store.challenge("10.0.0.5")
	if err != nil {
		t.Fatalf("挑战生成失败: %v", err)
	}
	store.mu.Lock()
	store.challenges[nonce] = clientAttestChallenge{expiresAt: time.Now().Add(-time.Second), ip: "10.0.0.5"}
	store.mu.Unlock()
	if store.consumeChallenge(nonce, "10.0.0.5") {
		t.Fatal("过期挑战不应可用")
	}
}

func TestClientAttestModeFallsBackToWarn(t *testing.T) {
	if got := normalizeClientAttestMode("EnFoRcE"); got != clientAttestModeEnforce {
		t.Fatalf("大小写应归一: %s", got)
	}
	if got := normalizeClientAttestMode("typo"); got != clientAttestModeWarn {
		t.Fatalf("未知取值应按 warn: %s", got)
	}
}
