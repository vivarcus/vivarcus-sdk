# 02-update-field

记录页 **All Actions** 按钮：通过 `SetValue` + `platform.Update` 写入 `title__c`。

## Quick start

```bash
cd examples/02-update-field
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
```

对象 / 按钮名写在 [main.go](main.go) 的 `Meta()` 与 [`mdl/`](mdl/) 里。

## Build only

```bash
cd examples/02-update-field
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
```

## 手工跟做

| 步骤 | 位置 |
|------|------|
| 创建对象 | [mdl/01-object.mdl](mdl/01-object.mdl) |
| 构建 + VPK + deploy | [_shared/scripts/package-vpk.sh](../_shared/scripts/package-vpk.sh) |

端到端示例见 [multi-component](../multi-component)。
