# 排错

## vivarcus 登录 HTTP 429（too many login attempts）

密码登录限流：**同一 IP + 用户名，1 分钟最多 4 次**（`POST /ui/auth/login`、Vault REST `POST /api/{version}/auth`）。

| 现象 | 处理 |
|------|------|
| `too many login attempts` / HTTP 429 | 按响应头 `Retry-After` 等待；勿连打 login |
| Agent 每条命令前 login | 改为注入 `VIVARCUS_TOKEN`，整段任务复用（最长 48h） |
| 多脚本同一用户连登 | 各脚本共用同一 token，或错开登录 |

详见 [01-prerequisites](01-prerequisites.md) 与 [vivarcus CLI 认证章节](https://github.com/vivarcus/vivarcus-cli/blob/main/docs/cli.md#agent--自动化避免登录限流)。

## vivarcus-sdk build 失败

| 现象 | 处理 |
|------|------|
| `requires go version 1.19 through 1.26, got go1.27` | 本机 Go 版本过新，TinyGo 暂不支持。用 `GOTOOLCHAIN=go1.26.2 vivarcus-sdk build ...`，或降级 Go；仅验证扫描/codegen 时加 `--skip-compile` |
| `tinygo not found` | 安装 TinyGo 并加入 `PATH` |
| `no Record Action entry type found` | 至少一个类型实现 `Meta`/`IsExecutable`/`Execute`（可在 `actions/`、`entries/` 或模块根） |
| `subdirectory entries must use a named package` | 子目录入口不能 `package main`，改为 `package actions` 等 |
| `must not live in shared/` | `shared/` 只放 helper；把 `Meta()` 入口移到 `actions/` / `triggers/` / `entries/` 或模块根 |
| `duplicate ... type name` | 同一 module 内类型名须唯一；或用 `Meta.Name` 钉 FQN |
| `imports forbidden package` | 移除 `os`/`net` 等禁止包 |
| `module exceeds 2MB` | 精简代码或依赖 |

## vivarcus package validate 失败

| issue | 处理 |
|-------|------|
| `gosdk_invalid` + 2mb / size | wasm ≤ 2 MB；重新 `vivarcus-sdk build` |
| `gosdk_invalid` + api_version | `__sdk_describe` 须返回 `api_version`=`1` |
| `gosdk_invalid` + import / whitelist | wasm 含非法 import；勿手写 wasm，用 `vivarcus-sdk build` |
| `gosdk_invalid` + go.mod / go.sum | VPK `gosdk/` 只放 `.go`；`module` 写在 `vaultpackage.xml` |
| `gosdk_invalid` + module path / `com.example` | `<gosdk><module>` 必填且不能是静默示例前缀 |

```bash
本地可用 `vivarcus-sdk build` 验证；VPK `gosdk/` 只放 `.go`，不要放 `.wasm` 或 `go.mod`
```

## 部署成功但按钮不出现

1. 确认 `ALTER Recordaction ... (active(true))` 已执行
2. 确认对应 `Objectaction` 也已 `active(true)`
3. 确认 `Meta.Object` 与当前记录对象一致
4. 确认 `IsExecutable` 返回 `true`
5. 确认用户有 object action 执行权限

## 点击按钮报错

| 错误 | 处理 |
|------|------|
| `record_action_inactive` | 重新激活组件 |
| `record_action_timeout` | 优化逻辑，避免长循环 |
| `record_action_wasm_trap` | 检查空指针、除零等；看 Vault SDK 日志 |
| `record_action_host_call_failed` | `platform.Update` 字段名/类型错误，或记录不存在 |

## Vault 未启用客户 wasm

托管环境须已启用客户 Record Action（wasm）运行时。若部署与 validate 均成功但执行报「不支持」类错误，请联系 Vault 管理员确认环境已开通该能力。

## 获取帮助

- Agent：读 [AGENTS.md](../AGENTS.md) 检查清单
- 示例对照：[examples/multi-component](../examples/multi-component)
