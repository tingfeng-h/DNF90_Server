# 运行时配置快照（2026-11-22，已脱敏）

对应基线：`baseline-20261122-perf`（`2fdc48c8`）。

| 文件 | 说明 |
| :--- | :--- |
| `compose.json` | 容器编排（端口、卷、环境变量引用） |
| `server.env` | 服务端环境变量（密钥已置 `CHANGE_ME`） |
| `database.env` | 数据库账号（密码已置 `CHANGE_ME`） |
| `server/client-attest.json` | 客户端校验（哈希置 `REDACTED_SHA256`，`mode` 保持原值） |
| `server/configs/` | dnfbridge / dnf logic / servergroup 配置 |
| `gm/gm.json` | GM 配置 |

## 脱敏规则

| 类型 | 替换为 |
| :--- | :--- |
| 本服公网 IP | `203.0.113.x`（RFC 5737 文档地址） |
| 内网网段与单机地址（形如 `10.x.x.x`） | `10.0.0.x` |
| 服务器绝对路径 | `/opt/dnf90/...` |
| 本机 Windows 客户端路径 | `C:\Users\example` / `D:\dnf90` |
| 个人账号标识 | `example-user` / `example` |
| env/json 的 password/secret/token/HMAC | `CHANGE_ME` |
| 64 位 SHA-256（client-attest） | `REDACTED_SHA256` |
| 邮箱 | `user@example.com` |

## 使用前

把 `CHANGE_ME` 换成自己的值；`client-attest.json` 的哈希换成自己客户端真实哈希，或先把 `mode` 设为 `off`。
