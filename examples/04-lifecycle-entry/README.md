# 04-lifecycle-entry

生命周期 **entry_action**、event_action、workflow step / cancel（见 `mdl/`）。

## Quick start

```bash
cd examples/04-lifecycle-entry
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/stamp_on_enter.go --json

# 生命周期 / workflow 规则引用 Recordaction.<FQN>，必须在 put 之后
export ACTION_FQN=acme.corp.lifecycle.StampOnEnter
python3 ../_shared/scripts/render-mdl.py mdl/02-lifecycle.mdl /tmp/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f /tmp/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f mdl/03-bind-object-lifecycle.mdl
```

FQN：`acme.corp.lifecycle.StampOnEnter`。对象是 `demo_request__c`（不是 `sdk_demo__c`）。workflow 绑法见 [`mdl/README.md`](mdl/README.md) 的 `04-workflow-action.mdl` / `05-workflow-cancel.mdl`。

不要上传 `zz_generated_reactor.go`（`sdk put` 只传 `actions/stamp_on_enter.go`）。

模块布局：[docs/03-build.md](../../docs/03-build.md)。系统路径概念：[docs/06-lifecycle.md](../../docs/06-lifecycle.md)。
