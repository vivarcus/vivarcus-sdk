# 04-lifecycle-entry

生命周期 **entry_action**、event_action、workflow step / cancel（见 `mdl/`）。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/04-lifecycle-entry`（monorepo：`cd sdk/examples/04-lifecycle-entry`）。

**顺序很重要**：先对象 MDL → `sdk put` 投影 Action → 再 apply 引用 `Recordaction.<FQN>` 的生命周期 MDL。`./deploy.sh` 已按这个顺序做。module 与其它示例相同，类型名 `StampOnEnter` 不与 `multi-component` 的 `StampDemoOnEnter` 冲突。

**一键**（推荐）：`./deploy.sh`（含 workflow 相关 MDL，与集成测试一致）

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/stamp_on_enter.go --json

export ACTION_FQN=acme.corp.sdkdemo.StampOnEnter
python3 ../_shared/scripts/render-mdl.py mdl/02-lifecycle.mdl /tmp/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f /tmp/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f mdl/03-bind-object-lifecycle.mdl
```

FQN：`acme.corp.sdkdemo.StampOnEnter`。对象是 `demo_request__c`（不是 `sdk_demo__c`）。workflow 绑法见 [`mdl/README.md`](mdl/README.md) 的 `04-workflow-action.mdl` / `05-workflow-cancel.mdl`。

不要上传 `zz_generated_reactor.go`（`sdk put` 只传 `actions/stamp_on_enter.go`）。

模块布局：[docs/03-build.md](../../docs/03-build.md)。系统路径概念：[docs/06-lifecycle.md](../../docs/06-lifecycle.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
