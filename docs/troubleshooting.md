# 排错

## ov-sdk build 失败

| 现象 | 处理 |
|------|------|
| `requires go version 1.19 through 1.26, got go1.27` | 本机 Go 版本过新，TinyGo 暂不支持。降级 Go 或仅用 `--skip-compile` 验证扫描/codegen |
| `tinygo not found` | 安装 TinyGo 并加入 `PATH` |
| `no Record Action entry type found` | 确保有且仅有一个类型实现 `Meta`/`IsExecutable`/`Execute` |
| `multiple entry types` | 一个模块只能有一个入口类型 |
| `imports forbidden package` | 移除 `os`/`net` 等禁止包 |
| `module exceeds 2MB` | 精简代码或依赖 |

## ov package validate 失败

| issue | 处理 |
|-------|------|
| `gosdk_invalid` + checksum | 重新计算 `sha256sum action.wasm`，更新 manifest |
| `gosdk_invalid` + api_version | 使用 `"api_version": "1"` |
| `gosdk_invalid` + import / whitelist | wasm 含非法 import；勿手写 wasm，用 `ov-sdk build` |

```bash
sha256sum gosdk/action.wasm
# 将 hex 填入 sdk_manifest.json 的 sha256 字段
```

## 部署成功但按钮不出现

1. 确认 `ALTER Recordaction ... active(true)` 已执行
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
- 示例对照：[examples/99-full-stack](../examples/99-full-stack)
