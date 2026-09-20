# 05-stamp-trigger

**Record Trigger** 示例：`BEFORE_INSERT` 在 `name__v` 后追加 `-trig`。

端到端（VPK deploy 投影 Recordtrigger + CREATE 时自动执行）见 [multi-component](../multi-component)。

```bash
cd examples/05-stamp-trigger
vivarcus-sdk build . -o action.wasm
vivarcus-sdk describe action.wasm
```

组件与对象定义见 [main.go](main.go) 与 [mdl/](mdl/)。
