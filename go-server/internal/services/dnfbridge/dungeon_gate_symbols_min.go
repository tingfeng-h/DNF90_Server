package dnfbridge

// 独立提供 action-object 运行时需要的"条件通过门"符号。
// 90CN 把完整实现放在 dungeon_boss_entrance.go（依赖 worldmap 的迷宫门解析），
// 那部分暂不移植；这里只提供运行时用到的 opcode 与固定 body。
//
// 当前客户端把 class-0/op0x0138 单独注册为 COMPLETE_CONDITION_PASS_GATE，
// body 固定 5 字节 u32(0)+u8(0)。
const currentDungeonCompleteConditionPassGateMsgID uint16 = 0x0138

var currentDungeonCompleteConditionPassGateBody = [...]byte{0, 0, 0, 0, 0}

func buildCurrentDungeonCompleteConditionPassGateBody() []byte {
	return append([]byte(nil), currentDungeonCompleteConditionPassGateBody[:]...)
}
