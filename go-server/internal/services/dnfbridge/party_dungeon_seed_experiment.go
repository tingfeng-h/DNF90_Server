package dnfbridge

import (
	"os"
	"strings"
)

// dungeonSeedIgnoreEnv 控制「普通副本退队时是否忽略房间 Seed」。
//
// 默认 off = 现行为：只有 RaidKey != 0 且 RaidInstanceID / RaidRunID 都相等的
// 跑团，才把 Seed 从比较里剔除。普通副本（RaidKey == 0）仍拿 Seed 比较。
//
// 问题在于导航推进房间会替换 Seed：noteDungeonPartyRoomLanding 把
// session.dungeonPartyRun.Seed 更新为新房间种子，而 instance.key 作为
// registry.byKey 的 Map 键，始终停在进图首房的 Seed。于是普通副本**过图之后**
// 任何队员退队都会在 withCurrentDungeonPartyRetirement 里
// `current != retired` 直接返回 false：
//   - RuntimePartyManager 里根本没删掉退队的人；
//   - 图里剩下的人收不到 op9 减员包和 NOTI6，退队者的角色/血条留在房间里，
//     过门检查卡死；
//   - 最后回城时同样被拦，两人在城镇里仍然是一个队伍。
//
// 打开后：Seed 一律不参与比较。跑团的实例身份仍由 RaidKey / RaidInstanceID /
// RaidRunID 三个字段在结构体比较里兜底，不会因为去掉 Seed 而把两次不同的
// 跑团混为一谈。
const dungeonSeedIgnoreEnv = "DNFBRIDGE_PARTY_DUNGEON_SEED_IGNORE"

func dungeonSeedIgnore() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(dungeonSeedIgnoreEnv)), "on")
}
