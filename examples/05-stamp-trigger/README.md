# 05-stamp-trigger

**Record Trigger** 示例：`BEFORE_INSERT` 给 `name__v` 追加 `-trig`。

## Build

```bash
cd examples/05-stamp-trigger
GOTOOLCHAIN=go1.22.12 vivarcus-sdk build . -o action.wasm
vivarcus-sdk describe action.wasm
```

Deploy 见 [multi-component](../multi-component)（Trigger 与 Action 可同 VPK）。
