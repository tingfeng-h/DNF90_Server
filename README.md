# DNF90 服务端（上游版本 + 本地修复）

本仓库是上游 [ServerS5A3](https://gitgud.io/16ww/servers5a3) 的二次分发：
**服务端代码整体跟随上游**，只在其上叠加少量本地修复，不另起分支重写。

| 项 | 值 |
|---|---|
| 上游基线 | `ed09f60`（2026-10-27） |
| 本分支 | `guild-merge` |
| 部署方式 | Docker Compose |
| 许可 | GPL v3（跟随上游） |

## 与上游的差异

### 玩法 / 修复

| 项目 | 效果 | 开关 |
|---|---|---|
| 深渊史诗掉率 | 深渊史诗固定掉落，件数 1–4 件按权重随机 | 代码内常量 |
| 任务放弃 | 持有任务物品的任务也能放弃（如深渊派对 Ep3），物品留在背包 | `DNFBRIDGE_QUEST_GIVEUP_ALLOW_ASSETS`（默认开） |
| 安图恩跳过影片 | 按 ESC 先弹"结束视频"确认框（与卢克一致） | `DNFBRIDGE_RAID_MOVIE_PROMPT_ANTON`（默认开） |
| 换频道团本状态 | 换频道/重连落点主动重推团本状态，修复"没有邀请攻坚队入口、看不到申请" | `DNFBRIDGE_RAID_CHANNEL_RESYNC`（默认开） |

### 性能

| 项目 | 效果 | 位置 |
|---|---|---|
| 分解 PVF 配置缓存 | 分解不再每次重读 PVF 配置 | 代码 |
| 分解背包写入范围 | 分解只写受影响的槽位，不整容器重写 | 代码 |
| 批量写入 | 一条 INSERT 最多 200 行，接任务/分解明显变快 | `DNFBRIDGE_MYSQL_INSERT_BATCH=200`（上游新版本已默认开启，保留作双保险） |

### 部署

| 项目 | 说明 |
|---|---|
| `deploy/docker/Dockerfile` | 构建阶段固定 `GOPROXY=https://goproxy.cn,direct`（国内网络构建必需） |
| `deploy/docker/server.env.example` | 完整环境变量示例（含性能与组队关键项），复制为 `server.env` 后填自己的值 |

## 部署

```bash
# 1. 拉代码
git clone -b guild-merge <本仓库地址> && cd ServerS5A3

# 2. 准备环境变量
cp deploy/docker/server.env.example server.env
#    填 4 个 ***：公网 IP、数据库密码、登录令牌、管理密钥
#    按自己的网络改 *_LISTEN / *_HOST 的 IP

# 3. 构建镜像并启动（compose 文件按自己的目录结构调整）
docker compose -p <项目名> -f <compose.json> up -d --build

# 4. 验证
docker ps                                  # 容器 healthy
docker logs --since 2m <server容器> | grep -c level=ERROR   # 应为 0
```

上游的完整构建/频道表/活动开关说明见 [Docs/BUILD_AND_DEPLOY.md](Docs/BUILD_AND_DEPLOY.md) 与 [Docs/CHANNEL_TABLE_DEPLOYMENT.md](Docs/CHANNEL_TABLE_DEPLOYMENT.md)。

## 部署后必看

| 环境变量 | 作用 | 建议值 |
|---|---|---|
| `DNFBRIDGE_MYSQL_INSERT_BATCH` | 批量写入行数，性能关键 | `200` |
| `DNFBRIDGE_PARTY_DUNGEON_PEER_POSITION` | 组队进图后队友"一移动就消失" | `off` |
| `DNFBRIDGE_PARTY_DUNGEON_PEER_ACTOR_PRE_REMOVE` | 预删除会连带清掉血条 | `off` |
| `DNFBRIDGE_PARTY_UDP_RELAY_DATAGRAM_TRACE` | UDP 抓包日志（量大） | `off` |

## 已知未修

| 问题 | 状态 |
|---|---|
| 装备耐久（客户端 op48） | 服务端未实现落库；上游未修，本仓库不做本地补丁 |
| 冒险团好友"在线状态" | 客户端不认行内状态字段；上游标注 P57"原生 UI 待验" |
| 全部修理后快捷栏装备不实时刷新 | 只修了穿戴槽，快捷栏未处理 |

## 同步上游

```bash
git fetch origin
git merge origin/master     # 上游历史偶有重写，以 merge 为准
```

冲突通常集中在 `initial_town_*.go` / `town_reconnect_seria_origin.go`：上游已实现同名的赛丽亚回城逻辑，**取上游版本**即可。

## 致谢

- 上游项目：[16ww/servers5a3](https://gitgud.io/16ww/servers5a3)（服务端、客户端补丁、文档）
- 参考实现：[rewio/ServerS4A21](https://gitgud.io/rewio/ServerS4A21)（86/A21 协议对照）
