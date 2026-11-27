package dnfbridge

import (
	"bytes"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/raid"
)

func TestRaidStartCountdownOpeningStateAndRollbackSwitch(t *testing.T) {
	now := time.Date(2026, 11, 27, 14, 0, 0, 0, time.UTC)
	battle := &runtimeRaidBattle{
		Configuration: raid.BattleConfiguration{StartDelaySeconds: 3},
		EnterAfter:    now.Add(3 * time.Second),
	}

	// 倒计时窗口内：先发 state=1（触发客户端 0x01F0A840 3秒倒计时 UI）。
	state, delay, countdown := currentRaidStartOpeningState(battle, now)
	if state != 1 || delay != 3*time.Second || !countdown {
		t.Fatalf("within countdown window got state=%d delay=%v countdown=%v", state, delay, countdown)
	}

	// 已过 EnterAfter（例如重连）：直接发 state=2，不重播倒计时。
	state, delay, countdown = currentRaidStartOpeningState(battle, now.Add(3*time.Second))
	if state != 2 || delay != 0 || countdown {
		t.Fatalf("after EnterAfter got state=%d delay=%v countdown=%v", state, delay, countdown)
	}

	// StartDelaySeconds == 0：直接发 state=2。
	noDelay := &runtimeRaidBattle{EnterAfter: now.Add(3 * time.Second)}
	state, delay, countdown = currentRaidStartOpeningState(noDelay, now)
	if state != 2 || delay != 0 || countdown {
		t.Fatalf("zero delay got state=%d delay=%v countdown=%v", state, delay, countdown)
	}

	// 环境变量回滚开关 off：直接发 state=2。
	t.Setenv(raidStartCountdownEnv, "off")
	state, delay, countdown = currentRaidStartOpeningState(battle, now)
	if state != 2 || delay != 0 || countdown {
		t.Fatalf("rollback switch got state=%d delay=%v countdown=%v", state, delay, countdown)
	}
}

func TestRaidStartCountdownCompletionSendsState2AfterDelay(t *testing.T) {
	now := time.Date(2026, 11, 27, 14, 0, 0, 0, time.UTC)
	queue := newFakeCurrentDungeonDeathTimerQueue()
	queue.now = now

	battle := &runtimeRaidBattle{
		InstanceID:    "inst-countdown-1",
		Phase:         0,
		Configuration: raid.BattleConfiguration{StartDelaySeconds: 3},
		EnterAfter:    now.Add(3 * time.Second),
	}
	conn := &bufferConn{}
	session := &gameSession{
		conn:                conn,
		connID:              "raid-countdown-conn",
		selectedCharacterID: 19,
		characterGeneration: 1,
	}
	s := &Service{
		gameplayTimers: queue,
		raids: map[uint32]*runtimeRaidState{
			77: {
				RaidKey: 77,
				Started: true,
				Battle:  battle,
			},
		},
	}

	s.scheduleCurrentRaidStartCountdownCompletion(session, battle, 77, 0, 3*time.Second, false)
	if conn.write.Len() != 0 {
		t.Fatalf("state2 sent prematurely before countdown elapsed: %x", conn.write.Bytes())
	}

	// 推进 3 秒：触发定时回调，补发 op587 = {2, 0}。
	queue.fire(queue.task(t, 0), false)
	p, rest := splitGameServerUpperPacket(t, conn.write.Bytes())
	if len(rest) != 0 || p.Header.MsgID != 587 || p.Header.Classification != 0 || !bytes.Equal(p.Body, []byte{2, 0}) {
		t.Fatalf("unexpected completion packet header=%+v body=%x rest=%x", p.Header, p.Body, rest)
	}

	// 若团本已结束，即使倒计时触发也不应再发送 state2。
	conn.write.Reset()
	battle.PhaseCompleted = true
	s.scheduleCurrentRaidStartCountdownCompletion(session, battle, 77, 0, 3*time.Second, false)
	queue.fire(queue.task(t, 1), false)
	if conn.write.Len() != 0 {
		t.Fatalf("ended raid emitted late state2: %x", conn.write.Bytes())
	}
}
