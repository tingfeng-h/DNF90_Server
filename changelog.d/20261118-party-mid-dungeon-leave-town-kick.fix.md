---
category: 功能修复
---

# 组队：中途退图后回城踢人/退队残留血条、回城红电脑

- 退图者的队伍记录只在 `owner=0`（副本作用域）投影，不再在城镇作用域新建 `record(D,ch)` 抢绑幸存者的
  `[localPlayer+0x4C8]`（客户端 `.prsafe` 的 `SetMemberUser` 在 `oldUser==newUser` 时不重绑，一旦被抢绑就再也绑不回来）。
- 彻底移除空 `kind2` 复位：城镇下会命中 `.prsafe 0x06EC50BE` 把绑定清零（`1→0` 分支 → `P2P->closeAll()`）；
  回城依赖 `.prsafe 0x06EC536B`（`[localPlayer+0x4C8]==0` 时自动绑回）。
- 退图后中继只 `sync`，不再 `CloseParty` 重建，避免掐断幸存者已建立的 UDP。
- 城镇点击队友（`mode 15`）与记录投影增加"同队非本机记录"双重守卫；`op10` 改用权威槽位 `target.Slot`。
- 入队顺序改为先建城镇 `CUser*` 再发 `op11`；中继删除 leg 时同步清 `portOwners`。
