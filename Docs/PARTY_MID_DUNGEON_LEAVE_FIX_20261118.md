# 组队「中途退图 → 回城踢人/退队」修复记录（2026-11-17 ~ 2026-11-18）

本文件记录本服针对 DNF90 组队系统的一组修复，作为 **2026-11-18 基线**。所有结论均来自线上 `packet_log` 抓包、
客户端逆向结论（`.prsafe` / `.xscene` / `.gld127` 补丁段）与逐版实机回归。

## 1. 客户端事实（逆向结论，修复依据）

| 事实 | 含义 |
| :-- | :-- |
| `[CUser + 0x4C8]` 是客户端"我当前属于哪条队伍记录"的唯一指针 | 只有 `record == [localPlayer+0x4C8]` 的 `NOTI 9` 才走 `wasBound=1` 分支 |
| `CPartyRecord::SetMemberUser`（`0x029733B0`）在 `oldUser == newUser` 时提前返回 | 一旦绑定被别的记录抢走，再发 `record(self,ch)` 的 `kind2/kind0` 都是空操作 |
| `.prsafe 0x06EC536B`：`[localPlayer+0x4C8] == 0` 时即使 `oldUser==newUser` 也会绑定 | 回城时 `op24` 重建 `localPlayer`，一条 `kind0` 即可自动绑回——**不需要空 `kind2`** |
| `.prsafe 0x06EC50BE`：空 `kind2`/`kind3` 在 `[actor+0x4C8] == this` 时会把绑定清零 | 对已绑定的记录发空 `kind2` ⇒ `wasBound 1→0` ⇒ "已退出队伍" + `P2P->closeAll()` |
| `.prsafe 0x06EC5227`：空 `kind2` 的拦截**只在城镇场景（scene==1）**生效 | 副本内发空 `kind2` 会直接清零绑定 |
| `.prsafe 0x06EC5361`：副本（scene==3）收到 `owner==0` 记录仅当 `[actor+0x4C8]==0` 时写入 | 副本内不会被 `owner=0` 记录抢绑；**城镇 `owner=ch` 记录会** |
| `NOTI 9` 的四个分支 | `0→0` 静默；`0→1` 建队/重组（可能 `P2P->closeAll()`）；`1→0` 退队（`P2P->closeAll()`）；`1→1` 增量更新（只 `CloseSlot`） |
| `op10`（`NOTI 10`）按 `CParty::GetSlotActor(权威固定槽位)` 索引 | 不能用客户端命令里的紧凑行号 |

## 2. 逐版过程

| 版本 | 镜像 | 改动 | 实机结果 |
| :-- | :-- | :-- | :--- |
| `v15-recon` | `d56af45d66d8` | 从当前树回退到 v15 行为（7 个函数级差异里退掉 4 个） | 回城行为接近 v15 |
| `v15-bindfix` | `4aeab4829470` | 引入绑定协调器：退图者只投 `owner=0`、回城 `kind3` 退役别名、`op10` 用权威槽位 | 回城踢人/退队恢复；但回城后队长看队员红电脑 |
| `v15-redpcfix` | `3c4d1ecd858f` | 退图后中继 `rebuild`、副本内空 `kind2` 复位、回城重建中继、队长优先编排 | **变差**：副本内退图即"已退出队伍→重连→防火墙" |
| `v15-norebuild` | `e46bdf84af39` | 回退上述四项（中继回 `sync`、删除空 `kind2` 复位、回城时序还原） | 副本内正常、回城正常；城镇反复踢人仍有缺槽/连接中 |
| **`v15-r4`（基线）** | **`be77ee836c64`** | `mode 15` 同队跳过 + 投影兜底守卫（永不投影"同队非本机城镇记录"）、入队先建 `CUser*` 再发 `op11`、中继 `portOwners` 清理、`resetUnboundSelfChannelRecordIfHijacked` 变空操作 | 作为基线；待长期观察 |

## 3. 与"真 v15"镜像的函数级差异（已知，未对齐）

| 函数 | 差异 |
| :-- | :-- |
| `syncRuntimePartyUDPRelay` | 比 v15 大 4397B（v18+ 中继硬化仍在） |
| `retireDungeonPartyAliasRecords` | 比 v15 小 302B（已修其中 `memberSet` 漏赋值的空 `if` 体） |
| `refreshDungeonScopedPartyMemberTables` | 只刷 `owner=0`（v15 行为），比 v15 小 129B |

> 这几处与组队记录/踢人无关，属于中继与副本记录刷新路径，留作后续。

## 4. 回滚清单

| tag | 镜像 ID | 说明 |
| :-- | :-- | :--- |
| `rollback-before-v15r4` | `e46bdf84af39` | 上一版（v15-norebuild） |
| `rollback-before-v15norebuild` | `3c4d1ecd858f` | v15-redpcfix（有副本内重连问题） |
| `rollback-before-v15redpcfix` | `4aeab4829470` | v15-bindfix |
| `merge-20261117-v15` | `7880221c924f` | 真 v15（无源码） |
| `v5 基线` | `8708f2ada37b` | v5 |

## 5. 验证方法

- 抓包：`/root/servers5a3-upstream/runtime/docker/server/logs/packet_log.txt`，按 `SEND kind=game-upper conn_id=... msg_id=9/10/11/153` 解码
- 关注事件：`game-party-client-departed-alias-records-retired`、`game-party-actor-projection-party-frame-skipped`、`game-party-walkout-invalid-slot`
- 反例特征：`owner=频道号` 且 `obj=退图者` 的 `kind2`、空 `kind2`（`count=0`）、`partse membership_changed` 后紧跟 `legs closed reason=close_party`
