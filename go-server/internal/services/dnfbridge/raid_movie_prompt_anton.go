package dnfbridge

import (
	"os"
	"strings"
)

// 把卢克（mode 2）的"按 ESC 先弹确认框"语义扩展到安图恩（mode 1）。
//
// 原有行为：只有卢克二阶段会在收到 C2S 659 body=0（ESC/跳过输入）时回一帧
// NOTI 599（影片控制器提示），客户端据此弹"是否结束视频"；安图恩（mode 1）
// 不回任何包，客户端本地直接关片，没有确认框。
//
// 开关：DNFBRIDGE_RAID_MOVIE_PROMPT_ANTON
//
//	默认（空/未设置/非法值）= on：安图恩也弹确认框；
//	off / 0 / false / no / legacy = 逐字节回到改动前（安图恩 ESC 直接关片）。
const raidMoviePromptAntonEnv = "DNFBRIDGE_RAID_MOVIE_PROMPT_ANTON"

func raidMoviePromptAntonEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(raidMoviePromptAntonEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}
