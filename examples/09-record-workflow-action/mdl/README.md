| 文件 | 说明 |
|------|------|
| [01-object.mdl](01-object.mdl) | `RECREATE` `sdk_demo__c`（与 02/07 等同对象） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | `RECREATE` 生命周期 + 用户动作启动工作流 |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | 对象绑定 `sdk_demo_lc__c` |
| [04-workflow.mdl](04-workflow.mdl) | `RECREATE` 对象工作流；Start 步 participant 控件挂 `Recordworkflowaction`（须在 `sdk put` **之后** apply） |

`Recordworkflowaction` 由 `sdk put` 投影，不要手写进 VPK `components/`。Start 步 XML 中的 reference 须与 FQN 一致：`Recordworkflowaction.acme.corp.recordworkflowaction.CaptureParticipants`（与 [workflowactions/capture_participants.go](../workflowactions/capture_participants.go) 的 module + 类型名对齐）。

端到端多组件打包见 [multi-component](../../multi-component)。
