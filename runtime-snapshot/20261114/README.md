# 运行时快照 2026-11-14

对应基线：

| 项 | 值 |
|---|---|
| commit（代码） | `9ff371f8`（tag `prod-20261114-upstreammerge`） |
| 镜像（游戏服务） | `servers5a3-upstream:merge-20261114` = `2cffc0ab1270` |
| 二进制 sha256 | `a7e06e9c87ac90c1695fb8d800724c2a062608587c7c256266f85140c6baf940` |
| 数据库镜像 | `dnf90-db-data:20261112` = `f731f63e2924`（MariaDB 11.8，307 表） |
| GM 镜像 | `dnf90_gm:20261021` = `57923d548ddb` |
| 上游基线 | `Cec1c/ServerS5A3` `2d99f016` |
| 客户端版本 | `official-24.9` |
| PVF sha256 | `29b66dc9274915cb1c2eabc6373dd9cf003557a7afc12e6d9bc9985cd92d5ae4` |

## 目录

```
runtime-snapshot/20261114/
├── compose.json                 # 生产 compose（项目名 servers5a3-local）
├── server.env                   # 游戏服务 env_file
├── database.env                 # 数据库 env_file
├── gm.json                      # GM 工具配置
├── schema.sql                   # dnf_local 全量表结构（307 表，无数据）
├── data/
│   ├── channel_info.etc         # 频道表（38 频道）
│   ├── channel.override.etc     # 频道覆盖表
│   └── manifest.json            # 频道表清单与校验
└── server-configs/
    ├── dnfbridge.toml           # 服务配置（admin token 已脱敏）
    ├── dnf/logic.toml           # 逻辑层配置（DB DSN 已脱敏）
    └── servergroup/plan.json    # 分片/分组计划
```

## 脱敏

| 占位符 | 位置 |
|---|---|
| `<DB_PASSWORD>` | `server.env`、`database.env`、`gm.json`、`server-configs/dnf/logic.toml` |
| `<DB_ROOT_PASSWORD>` | `database.env` |
| `<ADMIN_SECRET>` | `server.env`（`DNFBRIDGE_PLAYER_ADMIN_SECRET`）、`gm.json`（`gameToken`） |
| `<ADMIN_TOKEN>` | `server-configs/dnfbridge.toml`（`[admin] token`） |
| `<GM_PASSWORD>` | `gm.json` |

网络地址保留原值（`10.253.90.0/24` 容器内网、`10.0.4.9` 宿主机内网、`124.223.217.102` 公网），换机部署时改这三处。

## 未纳入仓库的大文件

| 文件 | 大小 | 原因 | 校验 |
|---|---|---|---|
| `data/Script.pvf` | 121 MB | 超过 GitHub 单文件 100 MB 限制 | sha256 `29b66dc9…5ae4` |
| `gm/state/item-index.json` | 48 MB | GM 首次启动按 PVF 自动重建 | 由 `pvfSha256` 校验 |

PVF 必须与上面 sha256 完全一致，否则物品 ID、任务、掉落都会错位。

## 复现步骤

```sh
# 1. 目录
mkdir -p runtime/docker/{data,server/configs/dnf,server/configs/servergroup,gm/state}
cp runtime-snapshot/20261114/{compose.json,server.env,database.env} runtime/docker/
cp runtime-snapshot/20261114/gm.json runtime/docker/gm/gm.json
cp runtime-snapshot/20261114/server-configs/dnfbridge.toml runtime/docker/server/configs/
cp runtime-snapshot/20261114/server-configs/dnf/logic.toml runtime/docker/server/configs/dnf/
cp runtime-snapshot/20261114/server-configs/servergroup/plan.json runtime/docker/server/configs/servergroup/
cp runtime-snapshot/20261114/data/* runtime/docker/data/

# 2. 填回占位符（见上表）

# 3. 放 PVF（sha256 必须匹配）
cp Script.pvf runtime/docker/data/Script.pvf

# 4. 初始化数据库
docker run -d --name tmp-db -e MARIADB_ROOT_PASSWORD=<DB_ROOT_PASSWORD> -e MARIADB_DATABASE=dnf_local \
  -e MARIADB_USER=dnf -e MARIADB_PASSWORD=<DB_PASSWORD> mariadb:11.8
docker exec -i tmp-db mariadb -uroot -p<DB_ROOT_PASSWORD> dnf_local < runtime-snapshot/20261114/schema.sql

# 5. 起服务
cd runtime/docker && docker compose -p servers5a3-local -f compose.json up -d
```

`compose.json` 里 db 服务用的是预置数据目录镜像 `dnf90-db-data:20261112`；若改用 `mariadb:11.8` + 卷，需把第 4 步的数据库写在同名卷 `servers5a3-local_database` 上，并把 `compose.json` 的 `db.image` 改掉。

## 备注

- `DNFBRIDGE_BUFF_GROUP_SYNC` 未设置（双组 Buff 默认关闭）；V212+ 客户端需要时设为 `on`。
- GM 的 `gameApi` 已修正为 `http://10.0.4.9:17002/api/game/admin`（旧副本误写 `10.0.0.9`）。
- 已知问题：选图界面直接回城时，红眼《血之狂暴》双刀仍会偶发（客户端行为）。
