# 05-stamp-trigger

**Record Trigger** 示例：`Meta()` + `Execute()` + `SetValue`，`BEFORE_INSERT` 给 `name__v` 追加 `-trig`。

端到端（VPK deploy 投影 Recordtrigger + CREATE 时自动执行）见 [multi-component](../multi-component)。本目录用于精读单一 Trigger 入口。

## Quick start

```bash
make -C sdk/examples build-05-stamp-trigger
vivarcus-sdk describe sdk/examples/05-stamp-trigger/action.wasm
```

`describe` 输出应含 `triggers[]`（FQN、object、events、event_segment、order）。

## Build only

```bash
vivarcus-sdk build ./examples/05-stamp-trigger -o action.wasm
```

Record Action 入门见 [01-hello-action](../01-hello-action)；按钮改字段见 [02-update-field](../02-update-field)。
