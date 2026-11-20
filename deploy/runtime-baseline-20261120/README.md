# 运行时配置（脱敏基线 2026-11-20）

本目录是线上运行配置的**脱敏副本**，仅用于记录"这版基线跑的是什么参数"。

## 脱敏规则

| 内容 | 处理 |
|---|---|
| `*PASSWORD*` / `*SECRET*` / `*TOKEN*` / `*DSN*` / `*_KEY*` 的值 | 替换为 `***REDACTED***` |
| 公网 IP | 替换为 `<SERVER_PUBLIC_IP>` |
| 内网 IPv4（10.x） | 替换为 `<LAN_IP>` |
| `127.0.0.1` | 保留（本机回环，无敏感信息） |
| GM 工具配置（`gm/gm.json`） | **不收录**（独立系统，见 DNF90GM 仓库） |

## 本基线对应

| 项 | 值 |
|---|---|
| Git 提交 | `22957087`（分支 `baseline-20261120` 后续会更新到此） |
| 镜像 | `servers5a3-upstream:v15-c7` = `6a0b4a1e5ea3` |
| 二进制 sha256 | `d9672c52c92be96dfbfba232eb95eed130392a92b76fd7acb8a783bbb7d2293a` |
| cgroup 内存 | 5.5 GiB（`5905580032`） |
| Go 运行时 | `GOGC=100` `GOMEMLIMIT=5GiB` `GOMAXPROCS=2` |
| 中继 | embedded，原生探针端口 2311/2312，转发端口 30000+ |

## 还原方式

把本目录的 `compose.json` / `server.env` / `database.env` 复制回部署机
`runtime/docker/`，再把 `***REDACTED***` 与占位 IP 换成本机真实值即可。

> 注意：`manifest.json` 是资源清单（无密钥），`compose.json` 里的端口映射/内存上限
> 与生产一致；`server.env` 里的 `DNFBRIDGE_*` 开关是本基线的实际配置。
