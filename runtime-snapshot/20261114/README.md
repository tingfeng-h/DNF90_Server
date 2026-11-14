# 运行时快照 2026-11-14

对应基线：

| 项 | 值 |
|---|---|
| commit | `9ff371f8`（tag `prod-20261114-upstreammerge`） |
| 镜像 | `servers5a3-upstream:merge-20261114` = `2cffc0ab1270` |
| 二进制 sha256 | `a7e06e9c87ac90c1695fb8d800724c2a062608587c7c256266f85140c6baf940` |
| 上游基线 | `Cec1c/ServerS5A3` `2d99f016` |
| 客户端版本 | `official-24.9` |
| PVF sha256 | `29b66dc9274915cb1c2eabc6373dd9cf003557a7afc12e6d9bc9985cd92d5ae4` |

## 文件

| 文件 | 说明 |
|---|---|
| `server.env` | 游戏服务容器环境变量（`compose.json` 的 `server.env_file`） |
| `database.env` | 数据库容器环境变量（`compose.json` 的 `db.env_file`） |
| `compose.json` | 生产 compose 定义（`servers5a3-local` 项目），不含明文密钥 |
| `gm.json` | GM 工具配置，挂载为容器内 `/config/gm.json` |

## 脱敏说明

以下值已替换为占位符，恢复时填回真实值：

| 占位符 | 出现位置 | 含义 |
|---|---|---|
| `<DB_PASSWORD>` | `server.env`、`database.env`、`gm.json` | `dnf` 数据库用户密码 |
| `<DB_ROOT_PASSWORD>` | `database.env` | MariaDB root 密码 |
| `<ADMIN_SECRET>` | `server.env`（`DNFBRIDGE_PLAYER_ADMIN_SECRET`）、`gm.json`（`gameToken`） | 游戏管理 API 与 GM 工具共享的令牌 |
| `<GM_PASSWORD>` | `gm.json` | GM Web 界面口令哈希 |

网络地址保留原值（`10.253.90.0/24` 为容器内网、`10.0.4.9` 为宿主机内网卡、`124.223.217.102` 为公网 IP），换机器部署时按实际环境调整。

`gm.json` 的 `gameApi` 已修正为 `http://10.0.4.9:17002/api/game/admin`（旧副本里误写成 `10.0.0.9`）。

## 目录结构

```
runtime/docker/
├── compose.json
├── server.env
├── database.env
├── gm/gm.json
├── data/            # Script.pvf、channel_info.etc 等只读资源
└── server/          # 容器内 /runtime，包含 configs/、data/、logs/
```

## 复现

```sh
cd <runtime>/docker
docker compose -p servers5a3-local -f compose.json up -d
```

数据库卷 `servers5a3-local_database` 必须保留（307 张表），只重建游戏服务时用：

```sh
docker compose -p servers5a3-local -f compose.json up -d --no-deps --force-recreate server
```

## 备注

- `DNFBRIDGE_BUFF_GROUP_SYNC` 未设置（双组 Buff 默认关闭）；客户端为 V212+ 且需要双组时显式设为 `on`。
- GM 的物品索引 `gm/state/item-index.json` 由 GM 首次启动按 `pvfSha256` 生成；PVF 更换后应删除该缓存重建。
- 已知问题：选图界面直接回城时，红眼《血之狂暴》双刀仍会偶发（客户端行为，服务端未消除）。
