# 超时空（FiendWar）直登：频道号装进 u8 owner

2026-10-30。代码提交 `def100f1`（`feat(channel): 超时空(type42)直登，落点真实城镇 98`）只覆盖服务端逻辑；
真正让客户端不崩的另一半是**运行时数据**，不在本仓库里，本文记录完整可复现步骤。

## 现象

- 客户端点超时空（改号前 261）：`game-connect channel_id=261 channel_type=42` 成功，握手、选角、`town_id=98 area_id=0`（`GContents/180621_Fiendwar/Town/Gate_Fiend.map`）全部完成，客户端还回了 type1345 场景就绪。
- 随后 1~7 秒内进程直接消失，启动器记录 `退出代码 3221225477`（`0xC0000005`），且**没有** op3 EXIT（正常退频道会先发 op3）。
- 中间态还有一次 `game-connection-bootstrap-failed error_len=80`，错误全文为
  `dnf resident channel field exceeds current EXE range: connected=261 resident=261`。

## 根因

- 场景 actor 的 owner 字段是 **u8**（`channel_owner_range.go`），客户端用它和自己的场景频道号比较；
  不一致时走 `mode0-owner-compare ... branch=remote-actor`，把自己的角色当跨频道远端 actor，再访问违例。
  这与当年 508（死亡之塔）崩溃是同一机制。
- 261..263 装不进 u8，且**没有任何** u8 值能等于 261，所以「公告 261 + owner 261」不可能同时成立。
- 上游现行口径（`raid_fiendwar_channel_availability.go`）也是：频道号装不进 u8 就判定不可用、从目录隐藏。

## 改动

### 1. 运行时频道表（`runtime/docker/data/channel_info.etc`）

```
-   261   `<chn_channel_info_048>`   42   `[fiendwar]`  ...
-   262   `<chn_channel_info_048>`   42   `[fiendwar]`  ...
-   263   `<chn_channel_info_048>`   42   `[fiendwar]`  ...
+   61    `<chn_channel_info_048>`   42   `[fiendwar]`  ...
+   62    `<chn_channel_info_048>`   42   `[fiendwar]`  ...
+   63    `<chn_channel_info_048>`   42   `[fiendwar]`  ...
```

type 42 与 `[fiendwar]` 组不变；端口由 `port = 10000 + id` 自动变成 `10061-10063`。

### 2. 运行时端口发布（`runtime/docker/compose.json`）

```diff
-        "10.0.0.9:10261:10261/tcp",
-        "10.0.0.9:10262:10262/tcp",
-        "10.0.0.9:10263:10263/tcp",
+        "10.0.0.9:10061:10061/tcp",
+        "10.0.0.9:10062:10062/tcp",
+        "10.0.0.9:10063:10063/tcp",
```

现网发布端口是 123 个，正好被 118 个频道占满，所以新增的三个 id 必须同时补端口映射。

### 3. 代码（本仓库，`def100f1`）

- `game.go`：`DNFBRIDGE_FIENDWAR_DIRECT_CHANNEL=on` 时 type42 不再被 `catalog.ResidentFor` 换成普通落点。
- `channel_resident_notice.go`：owner 回退 `context 0` 改为**仅当频道号装不下 u8 时**生效；
  61..63 装得下，场景 owner 因此等于真实频道号（这一步是崩溃的最后一处）。
- `game_bootstrap.go` / `initial_town_reconnect.go`：装不下时的兜底与换线重连放行。

## 验证（2026-10-30 现网）

| 项 | 结果 |
| --- | --- |
| 连接 | `game-connect channel_id=61 channel_type=42 game_port=10061` |
| owner | `owner_channel=61`（场景对象列表、town 入口、route endpoint 全部一致，无 0） |
| 落点 | `town_id=98 area_id=0`，随后 `town-set-user-area-committed town_id=98 area_id=1` |
| 团本 UI | `game-current-raid-profiles-sent`、`game-raid-entry-cost-sent item_id=490007844 required=1`、`超时空频道开放状态下发 available=true` |
| 稳定性 | 无 `game-read-close`/`game-close`，启动器无 `3221225477` |

## 回退

- 运行时：还原 `channel_info.etc.bak-renumber-fiendwar-*` 与 `compose.json.bak-renumber-fiendwar-*`，再 `--force-recreate server`。
- 代码：`git revert def100f1`。

## 注意

- 频道号对玩家可见：列表里显示 61/62/63，不再是 261/262/263。
- 换机器/重装必须同时带走这两个运行时文件，只部署代码不会生效。
- 若将来 PVF 或频道表重新出现 >255 的专用频道，同样的问题会复现；`DNFBRIDGE_CHANNEL_DIRECTORY_OWNER_RANGE=on` 可把这类频道从目录隐藏（不崩溃、也进不去）。
