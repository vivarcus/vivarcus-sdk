# 01-hello-action

Minimal Record Action — validates `vivarcus-sdk` scan, codegen, and wasm compile.

## Build

```bash
# scan + codegen only (no tinygo)
vivarcus-sdk build ./examples/01-hello-action --skip-compile

# full wasm (requires tinygo)
vivarcus-sdk build ./examples/01-hello-action -o action.wasm
```

## What it does

`Execute` is a no-op. Use this example to verify your toolchain before deploying real logic.

See [99-full-stack](../99-full-stack) for MDL + VPK + deploy.
