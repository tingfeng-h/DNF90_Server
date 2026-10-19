# v141 可玩基线

## 源码
- 基线来源：servers5a3 upstream master
- 导入提交：dfd6587cf5bfbe1a424622b56e58ba55f8811c20
- 开发分支：local/v141-playable-baseline

## 客户端
- DNF.exe SHA256：
  26abf551753c2745c2e5b349fb24336422a189187ff847c995e0df9a68871079
- dstr.dat SHA256：
  ee8904c8540ef63807dc984604a3ea0f5cc4040228e5ac2e95d490c03d21ec08
- 客户端 Script.pvf SHA256：
  9dcd31a3e6001fd5563d36dc65cbaa96bed8b782700cba30625197e27cf1e43c
- 发行清单要求的 PVF SHA256：
  a56cacf96087ba34ca5299e04ab65f8268f694b214b97dec9879fdebf46f5290
- PVF 哈希不同；当前可玩测试不代表资源版本完全匹配。

## 已验证的运行配置
- 使用独立 Docker Compose 项目和独立 MariaDB volume。
- channel_advertise_server=1。
- 使用上游参考频道表：
  go-server/internal/services/dnfbridge/testdata/channel_info.etc
- 参考频道表 SHA256：
  b1741336a9f91646ece865c799d187d6c36460bc826bf9efbc1b850b39a71a36
- 频道引导 TCP 端口：7001。
- 账号 API TCP 端口：17002。
- 管理端口：18111，仅 loopback。
- 游戏频道端口和 UDP 中继以运行时 catalog 为准。

## 实机结果
- v141 客户端通过自制登录器启动。
- 成功连接频道、进入游戏并刷图。
- 团本、组队和旧数据迁移尚未完成全面验收。

## 安全注意
- runtime/docker/instance.json 和数据库卷不进入 Git。
- 不提交密码、Token、私钥、运行日志或数据库转储。
- 临时 HTTP 账号 API 仅用于隔离测试；长期运行前必须加固认证入口。
- 排查结束后关闭 packet_log，避免记录敏感认证数据。
