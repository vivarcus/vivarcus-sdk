# 04-lifecycle-entry

生命周期 **entry_action**、event_action、workflow step / cancel（见 `mdl/`）。

## Quick start

```bash
cd examples/04-lifecycle-entry
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
# 先 apply lifecycle MDL，再 package-vpk + deploy（见 mdl/README.md）
```

## Build only

```bash
cd examples/04-lifecycle-entry
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
```

Combo VPK 与 [multi-component](../multi-component) 相同模式。
