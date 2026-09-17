# 01-hello-action

最小 Record Action — 验证 `vivarcus-sdk` scan、codegen、wasm 编译。

## Build

```bash
vivarcus-sdk build ./examples/01-hello-action --skip-compile
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build ./examples/01-hello-action -o action.wasm
```

## What it does

`Execute` 为空操作。确认工具链后做 [02-update-field](../02-update-field) 或 [04-lifecycle-entry](../04-lifecycle-entry)。

Vault 完整端到端见 [99-full-stack](../99-full-stack)（自包含 `main.go` + `mdl/`）。
