# 01-hello-action

最小 Record Action — 验证 `vivarcus-sdk` scan、codegen、wasm 编译。

## Build

```bash
# scan + codegen only (no tinygo)
vivarcus-sdk build ./examples/01-hello-action --skip-compile

# full wasm (requires tinygo)
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build ./examples/01-hello-action -o action.wasm
```

## What it does

`Execute` 为空操作。用于确认工具链后再做 [02-update-field](../02-update-field) 或 [04-lifecycle-entry](../04-lifecycle-entry)。

部署到 Vault 的完整流程见 [multi-component](../multi-component)（不依赖本目录）。
