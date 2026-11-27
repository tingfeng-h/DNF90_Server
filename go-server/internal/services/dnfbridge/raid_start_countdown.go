package dnfbridge

import (
	"fmt"
	"strings"
	"time"
)

const (
	// raidStartCountdownEnv 控制团本开团/换阶段时的原生 3 秒屏幕中央倒计时：
	//   - 默认开启：当 StartDelaySeconds > 0 且尚未到达 EnterAfter 时，
	//     先发 class0/op587 = {1, phase}（触发客户端 0x01F0A840: push 3 倒计时 UI），
	//     在 EnterAfter 到点后补发 class0/op587 = {2, phase}（自动关闭倒计时并播放开始横幅）；
	//   - 设为 off/false/0/no/legacy 时逐字节回退至旧版直接发 {2, phase}。
	raidStartCountdownEnv = "DNFBRIDGE_RAID_START_COUNTDOWN"
)

func raidStartCountdownEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(envString(raidStartCountdownEnv, ""))) {
	case "off", "false", "0", "no", "legacy":
		return false
	default:
		return true
	}
}

func currentRaidStartOpeningState(b *runtimeRaidBattle, now time.Time) (state byte, delay time.Duration, countdown bool) {
	if !raidStartCountdownEnabled() || b == nil || b.Configuration.StartDelaySeconds == 0 || b.EnterAfter.IsZero() {
		return 2, 0, false
	}
	remaining := b.EnterAfter.Sub(now)
	if remaining <= 0 {
		return 2, 0, false
	}
	return 1, remaining, true
}

func (s *Service) scheduleCurrentRaidStartCountdownCompletion(
	session *gameSession,
	b *runtimeRaidBattle,
	raidKey uint32,
	phase byte,
	delay time.Duration,
	includeFiendWarState bool,
) {
	if s == nil || session == nil || b == nil || delay <= 0 {
		return
	}
	bound, ok := s.boundGameSessionCharacterSnapshot(session)
	if !ok {
		bound = boundGameSessionCharacter{
			session:    session,
			character:  session.selectedCharacterID,
			generation: session.characterGeneration,
		}
	}
	if bound.character == 0 {
		return
	}
	instanceID := b.InstanceID
	timerName := fmt.Sprintf("raid-start-countdown:%s:%d:%d", instanceID, phase, bound.character)
	deliver := func() error {
		if session.conn == nil || session.selectedCharacterID != bound.character || session.characterGeneration != bound.generation {
			return nil
		}
		s.raidMu.Lock()
		state := s.raids[raidKey]
		active := state != nil && state.Started && state.Battle != nil &&
			state.Battle.InstanceID == instanceID &&
			state.Battle.Phase == phase &&
			!state.Battle.TimedOut &&
			!state.Battle.PhaseCompleted
		s.raidMu.Unlock()
		if !active {
			return nil
		}
		if includeFiendWarState {
			if err := s.sendGameUpperRawClass(session, 587, []byte{2, phase}, 0); err != nil {
				return err
			}
			if err := s.sendGameUpperRawClass(session, currentFiendWarStateNotice, []byte{2, phase}, 0); err != nil {
				return err
			}
		} else {
			if err := s.sendGameUpperRawClassCodec(session, 587, []byte{2, phase}, 0, true); err != nil {
				return err
			}
		}
		s.logGameEvent(session, "game-raid-start-countdown-completed",
			"raid_key", raidKey,
			"instance_id", instanceID,
			"phase", phase,
			"char_id", bound.character,
			"fiendwar_state", includeFiendWarState)
		return nil
	}
	callback := func(time.Time) {
		if session.events == nil || session.events.loop == nil {
			_ = deliver()
			return
		}
		_ = s.tryPostGameSessionCharacterEventWithDiscard(
			session,
			"raid-start-countdown-complete",
			bound.character,
			bound.generation,
			deliver,
			nil,
		)
	}
	if queue := s.ensureGameplayTimeQueue(); queue != nil {
		if err := queue.ScheduleAfter(timerName, delay, callback); err != nil {
			s.logWarn("dnfbridge raid start countdown schedule failed",
				"raid_key", raidKey,
				"char_id", bound.character,
				"phase", phase,
				"error", err)
		}
	}
}
