# 09-record-workflow-action

**Record Workflow Action** 示例：在工作流 **Start** 步骤的 `GET_PARTICIPANTS` 事件里，把当前参与人组默认填成发起人。

在工作流步骤事件（如 `GET_PARTICIPANTS`）里运行客户逻辑。这不是记录页按钮（`Recordaction`），也不是对象工作流 XML 里的系统 Workflow Action 步骤（见 [04-lifecycle-entry](../04-lifecycle-entry)）。

[`mdl/04-workflow.mdl`](mdl/04-workflow.mdl) 在 Start 步 participant 控件上引用本示例的 `Recordworkflowaction`（见 [`mdl/README.md`](mdl/README.md)）。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/09-record-workflow-action`（monorepo：`cd sdk/examples/09-record-workflow-action`）。

对象与生命周期先就绪，**put 工作流 Action 后再** apply 引用 FQN 的 `04-workflow.mdl`。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus component apply-mdl --confirm -f mdl/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f mdl/03-bind-object-lifecycle.mdl
vivarcus sdk put -f workflowactions/capture_participants.go --json
vivarcus component apply-mdl --confirm -f mdl/04-workflow.mdl
```

FQN：`acme.corp.sdkdemo.CaptureParticipants`。对象是 `capture_demo__c`（不共用 `sdk_demo__c`）。在记录上执行生命周期动作 **Start Capture Participants WF**，打开工作流启动对话框验证参与人预填。

当前平台会投影 `Recordworkflowaction` 组件。Start 步骤上已挂的 FQN 会在 `GET_PARTICIPANTS` 时执行客户 wasm，并把返回的参与人写回启动对话框。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
