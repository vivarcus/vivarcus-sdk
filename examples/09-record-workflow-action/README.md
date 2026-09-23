# 09-record-workflow-action

**Record Workflow Action** 示例：在工作流 **Start** 步骤的 `GET_PARTICIPANTS` 事件里，把当前参与人组默认填成发起人。

对标 Veeva Java SDK `com.veeva.vault.sdk.api.workflow.RecordWorkflowAction`（`@RecordWorkflowActionInfo` + `execute`）。这不是记录页按钮（`Recordaction`），也不是对象工作流 XML 里的系统 Workflow Action 步骤（见 [04-lifecycle-entry](../04-lifecycle-entry)）。

[`mdl/04-workflow.mdl`](mdl/04-workflow.mdl) 在 Start 步 participant 控件上引用本示例的 `Recordworkflowaction`（见 [`mdl/README.md`](mdl/README.md)）。

## Quick start

```bash
cd examples/09-record-workflow-action
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus component apply-mdl --confirm -f mdl/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f mdl/03-bind-object-lifecycle.mdl
vivarcus sdk put -f workflowactions/capture_participants.go --json
vivarcus component apply-mdl --confirm -f mdl/04-workflow.mdl
```

FQN：`acme.corp.recordworkflowaction.CaptureParticipants`。在 `sdk_demo__c` 记录上执行生命周期动作 **Start Capture Participants WF**，打开工作流启动对话框验证参与人预填。

当前平台会投影 `Recordworkflowaction` 组件；Start 步骤上已挂的 FQN 会在 `GET_PARTICIPANTS` 时被调用。客户 wasm 执行路径正在接入，若运行时报 `customer wasm not supported yet`，那是运行时还未跑 guest 模块，与本示例代码无关。
