package dnfbridge

import (
	"os"
	"strings"
)

// 现场（2026-10-26）：接不了也放不掉的深渊派对 Ep3 任务（quest 4539）返回
// `quest give-up requires an inventory transaction` —— owner 要求放弃时把任务物品
// 一起回收，而当前实现只做"无资产"放弃，于是整条拒绝。
//
// 本开关跳到"只移除任务状态、任务物品留在背包"的放弃路径（与客户端表现一致：
// 任务消失，背包里的邀请函等道具还能自己处理）。
//
// 开关：DNFBRIDGE_QUEST_GIVEUP_ALLOW_ASSETS
//
//	默认 on；off / 0 / false / no / legacy = 回到上游严格语义（持有任务物品的任务不可放弃）。
const questGiveUpAllowAssetsEnv = "DNFBRIDGE_QUEST_GIVEUP_ALLOW_ASSETS"

func questGiveUpAllowAssets() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(questGiveUpAllowAssetsEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}
