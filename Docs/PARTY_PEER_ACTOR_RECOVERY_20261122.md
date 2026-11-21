# 组队战斗卡顿：队友 actor 延迟重投链（2026-11-22）

## 现象

- 两台电脑 4 开、**4 人组队刷图**时，战斗过程中"一卡一卡"（整屏顿，非只有队友），
  **单刷或 4 开但不组队完全流畅**。
- 2 人组队"比 4 人好一点但仍卡"。
- 客户端侧：CPU 总占用 6.5%（无任何逻辑核心打满）、GPU 11.5%。
- 服务端侧：CPU 7~10%、GC 每 41 秒一次（标记 323ms）、战斗中无 300ms 以上慢事件、
  socket 写超过 20ms 全程仅 1 次（27ms）。
- 链路：4 条游戏连接 RTT 65~71ms、抖动 4~11ms、`app_limited`，每客户端 0.4~1.0 Mbps。

## 根因

`scheduleDungeonPeerActorRecovery`（`dungeon_peer_actor_recovery.go`）会在队友 actor 投影后
启动一条**延迟补投链**：

```go
dungeonPeerRecoveryDelays = {750ms, 1250ms, 3s, 3s, 3s, 3s}   // 最多 6 轮
dungeonPeerRecoveryBudget = 15s
```

每一轮都会**重发该队友完整的 mode0 / mode1 / subtype0 / op9 角色状态**。停止条件是客户端
回应 `0x7f14 (ARQ1)` 原生注册查询并把 `nativeRegistered` 置位。

但本服客户端**自 2026-11-15 起再没有回应过 ARQ1**，于是每次投影都会把 6 轮重投跑满，
客户端每收到一次重复的全量 actor 状态就重建一次该角色，表现为**每隔 1~3 秒卡一下**。

### 日志证据（`packet_log.txt`）

| 事件 | 次数 |
| :--- | ---: |
| `game-dungeon-party-peer-actor-send`（正常投递） | 50,576 |
| `game-dungeon-party-peer-actor-delivered` | 50,528 |
| **`game-dungeon-party-peer-actor-recovery`（重投）** | **37,754** |
| `game-dungeon-party-peer-actor-resend` | 3,491 |
| `game-dungeon-party-peer-native-registration`（客户端确认） | 6,706 |

- 最近 3000 行日志：14 次 send / 14 次 delivered / **12 次 recovery** / 2 次 resend；
- 重投轮次实测跑到 `round=1..5`；
- 重投记录里 `send_completed=true` 全部成功、`client_native_result=unknown` 全部未知；
- **最后一次** `native-registration`：`2026-11-15T17:04:12Z`（`native_status=3`），之后再无记录。

### CPU 档案佐证（20 秒 pprof，4 人战斗）

```
sendDungeonPartyPeerActorStateObserved        36~47%
  └─ armDungeonPeerActorRecovery.func1.1      20~38%
       ├─ buildCurrentUserInfoMode1WithEquipmentInContext
       └─ mysqlCharStore.Load / loadCharacterSupplement
```

## 修复

新增开关 `DNFBRIDGE_DUNGEON_PEER_RECOVERY`（`dungeon_peer_recovery_switch.go`）：

| 取值 | 行为 |
| :--- | :--- |
| 空、`off`、`0`、`false`、`no`、其它 | **关闭延迟补投链（默认）** |
| `on`、`1`、`true`、`yes`、`legacy` | 恢复旧的 6 轮重投链（A/B 对比用） |

`scheduleDungeonPeerActorRecovery` 入口直接按开关短路。

**只关"延迟补投链"**，不影响：

- 首次队友投影（`projectDungeonPartyPeerActors`）；
- 进图落位、组队状态机与同屏逻辑；
- 帧内容与帧内顺序。

## 验证

现场验证（2026-11-22）：关闭后 4 人组队战斗卡顿**消失**。

回滚方法：

```bash
# runtime/docker/server.env 追加
DNFBRIDGE_DUNGEON_PEER_RECOVERY=on
# 然后重启 server 容器
```

## 后续

- 若将来确认客户端会正确回应 `0x7f14 (ARQ1)`，或把服务端改为"TCP 写成功即视为投递完成"，
  才考虑重新开启补投链；
- 在此之前保持默认关闭。

## 相关文件

- `go-server/internal/services/dnfbridge/dungeon_peer_actor_recovery.go`（补投链本体）
- `go-server/internal/services/dnfbridge/dungeon_peer_native_registration.go`（`0x7f14` ARQ1 查询/回包）
- `go-server/internal/services/dnfbridge/dungeon_peer_recovery_switch.go`（新增开关）

## 附：观测方法

```bash
# 重投链计数
grep -c 'game-dungeon-party-peer-actor-recovery' packet_log.txt
grep 'game-dungeon-party-peer-actor-recovery' packet_log.txt | tail -5
grep 'game-dungeon-party-peer-native-registration' packet_log.txt | tail -3

# 慢事件（300ms 以上事件循环阻塞）
DNFBRIDGE_EVENT_SLOW_THRESHOLD_MS=300

# CPU 档案
docker exec servers5a3-local-server-1 sh -c \
  'wget -q -O /tmp/cpu.pb.gz "http://127.0.0.1:6060/debug/pprof/profile?seconds=20"'
go tool pprof -top -cum DNF90Server cpu.pb.gz
```
