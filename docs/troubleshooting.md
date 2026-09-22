# 排错

## vivarcus 登录 HTTP 429（too many login attempts）

密码登录限流：**同一 IP + 用户名，1 分钟最多 4 次**（`POST /ui/auth/login`、Vault REST `POST /api/{version}/auth`）。

| 现象 | 处理 |
|------|------|
| `too many login attempts` / HTTP 429 | 按响应头 `Retry-After` 等待；勿连打 login |
| Agent 每条命令前 login | 改为注入 `VIVARCUS_TOKEN`，整段任务复用（最长 48h） |
| 多脚本同一用户连登 | 各脚本共用同一 token，或错开登录 |

详见 [01-prerequisites](01-prerequisites.md) 与 [vivarcus CLI 认证章节](https://github.com/vivarcus/vivarcus-cli/blob/main/docs/cli.md#agent--自动化避免登录限流)。

## 平台编译 / validate 失败（默认看这里）

VPK 只含源码时，这些错误来自 **Vault 编译**。编号示例改 `.go` 后重新 `sdk put`；multi-component 重新打包 deploy。

| issue | 处理 |
|-------|------|
| `gosdk_invalid` + 2mb / size | wasm ≤ 2 MB；精简客户代码或依赖 |
| `gosdk_invalid` + api_version | `__sdk_describe` 须返回 `api_version`=`1`（入口类型实现不全时常出现） |
| `gosdk_invalid` + import / whitelist | 含非法 import（`os`/`net` 等）；勿手写 wasm 塞进 VPK |
| `gosdk_invalid` + go.mod / go.sum | VPK `gosdk/` 只放 `.go`；`module` 写在 `vaultpackage.xml` |
| `gosdk_invalid` + local module harness | 根 `main.go` 只有空导入；留在本地，不要打进 VPK 或 `sdk put`。入口类型写在该文件里时可以上传 |
| `gosdk_invalid` + must not declare func main | 删掉客户源码里的 `func main`。平台 reactor 会生成唯一的 `func main` |
| `gosdk_invalid` + module path / `com.example` | `<gosdk><module>` 必填且不能是静默示例前缀 |
| `no Record Action entry type found` | 至少一个类型实现 `Meta`/`IsExecutable`/`Execute`（可在 `actions/`、`entries/` 或模块根） |
| `subdirectory entries must use a named package` | 子目录入口不能 `package main`，改为 `package actions` 等 |
| `must not live in shared/` | `shared/` 只放 helper；把 `Meta()` 入口移到 `actions/` / `triggers/` / `entries/` 或模块根 |
| `duplicate ... type name` | 同一 module 内类型名须唯一；或用 `Meta.Name` 钉 FQN |
| `imports forbidden package` | 移除 `os`/`net` 等禁止包 |

VPK `gosdk/` 只放业务 `.go`，不要放 `.wasm`、`go.mod`、薄 `main.go`，也不要在源码里写 `func main`。

## 部署成功但按钮不出现

1. `vivarcus sdk enable <FQN> --json`（或 `ALTER Recordaction ... (active(true))`）
2. 确认对应 `Objectaction` 也已 `active(true)`
3. 确认 `Meta.Object` 与当前记录对象一致
4. 确认 `IsExecutable` 返回 `true`
5. 确认用户有 object action 执行权限
6. `vivarcus sdk get <FQN> -o /tmp/src.go --json` 能拉到客户 Go，说明组件已投影。编号示例用 `vivarcus sdk put -f <子目录>/<类型名>.go --json` 部署或覆盖；多文件树才 re-import VPK

## 点击按钮报错

| 错误 | 处理 |
|------|------|
| `record_action_inactive` | `vivarcus sdk enable <FQN> --json` |
| `record_action_timeout` | 优化逻辑，避免长循环 |
| `record_action_wasm_trap` | 检查空指针、除零等；拉 Runtime Log（见下） |
| `record_action_host_call_failed` | `platform.Update` 字段名/类型错误，或记录不存在 |

## Runtime Log（`vivarcus sdk logs`）

guest `platform.LogInfo` / 执行异常落在 Vault 的 SDK Runtime Log。Agent 用用户 CLI 按 UTC 日下载 ZIP：

```bash
DATE=$(date -u +%Y-%m-%d)
./bin/vivarcus sdk logs --date "$DATE" --format csv -o /tmp/SdkLog-$DATE.zip --json
unzip -p /tmp/SdkLog-$DATE.zip
```

- `--date` 必填；最多 30 天；需要 Logs 权限（`debug_log` 或 `all_audit`）。
- `--json` 只打印 `{"path":...}`，ZIP 在 `-o` 文件里。不要 `-o -` 配 `--json`。
- 完整调用约定见仓库 `vivarcus-cli` SKILL。

## Vault 未启用客户 wasm

托管环境须已启用客户 Record Action（wasm）运行时。若部署与 validate 均成功但执行报「不支持」类错误，请联系 Vault 管理员确认环境已开通该能力。

## 获取帮助

- Agent：读 [AGENTS.md](../AGENTS.md) 检查清单
- 示例对照：[examples/multi-component](../examples/multi-component)
