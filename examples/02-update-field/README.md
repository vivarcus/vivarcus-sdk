# 02-update-field

记录页 **All Actions** 按钮：通过 `SetValue` + `platform.Update` 写入 `title__c`。

## Quick start

```bash
vivarcus auth login && vivarcus config set default_vault <uuid>
./scripts/dev-demo.sh
```

组件定义在 [`mdl/`](mdl/) 与 [main.go](main.go) 的 `Meta()`。脚本：build → combo VPK → API 验证。

## Build only

```bash
cd examples/02-update-field
vivarcus-sdk build . -o action.wasm
```

## 手工跟做

| 步骤 | 位置 |
|------|------|
| 创建对象 | VPK `components/00010/`（[mdl/01-object.mdl](mdl/01-object.mdl)） |
| 打 VPK + deploy | [_shared/scripts/package-vpk.sh](../_shared/scripts/package-vpk.sh) |

部署细节：[docs/05-deploy.md](../../docs/05-deploy.md)。端到端索引：[multi-component](../multi-component)。系统路径见 [04-lifecycle-entry](../04-lifecycle-entry)。
