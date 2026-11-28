# 运行时配置快照（脱敏）— baseline-20261128

本目录是 `baseline-20261128`（v15-c60）上线版本的运行时配置快照，已**脱敏**：所有密码、令牌、DSN 口令、
管理密钥与真实公网/内网地址都替换为占位符，可直接作为部署模板使用。

| 文件 | 说明 |
|---|---|
| `compose.json` | Docker Compose（db + server，端口映射、健康检查、内存上限） |
| `server.env` | 游戏服环境变量（DNFBRIDGE_* 开关） |
| `database.env` | MariaDB 账号环境变量（口令已脱敏） |
| `server/configs/dnfbridge.toml` | 服务端主配置（admin token 已脱敏） |
| `server/configs/dnf/logic.toml` | 逻辑层配置（MySQL DSN 口令已脱敏） |
| `server/configs/servergroup/plan.json` | 区服分组计划 |
| `server/client-attest.json` | 客户端校验允许清单（file hash，非敏感） |
| `gm/gm.json` | GM 工具配置（password / gameToken / DSN 口令 / PVF hash 已脱敏） |

占位符：

- `<SERVER_PUBLIC_IP>`：玩家可达的公网 IP
- `<INTERNAL_IP>`：Docker 内网 IP（db/server 通信）
- `<HOST_LAN_IP>`：宿主机局域网 IP（端口发布地址）
- `CHANGE_ME`：密码 / 令牌 / DSN 口令
- `<PVF_SHA256>`：GM 工具记录的 PVF 哈希

部署前把上述占位符替换为你的实际值即可。真实 PVF、数据库文件与日志不包含在本快照内。
