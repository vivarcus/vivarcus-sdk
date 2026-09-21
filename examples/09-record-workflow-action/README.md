# 09-record-workflow-action

**Record Workflow Action** 示例：在工作流 **Start** 步骤的 `GET_PARTICIPANTS` 事件里，把当前参与人组默认填成发起人。

对标 Veeva Java SDK `com.veeva.vault.sdk.api.workflow.RecordWorkflowAction`（`@RecordWorkflowActionInfo` + `execute`）。这不是记录页按钮（`Recordaction`），也不是对象工作流 XML 里的系统 Workflow Action 步骤（见 [04-lifecycle-entry](../04-lifecycle-entry)）。

部署后在对象工作流 Start 步骤的 Participant Control 上引用：

```xml
<vwf:recordWorkflowAction reference="Recordworkflowaction.acme.corp.recordworkflowaction.CaptureParticipants" />
```

## Build

```bash
cd examples/09-record-workflow-action
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
vivarcus-sdk describe action.wasm
```

Deploy 见 [multi-component](../multi-component)。对象工作流须已存在，并在 Start 步骤挂上上面的 `Recordworkflowaction` 引用。

当前平台会投影 `Recordworkflowaction` 组件；Start 步骤上已挂的 FQN 会在 `GET_PARTICIPANTS` 时被调用。客户 wasm 执行路径正在接入，若运行时报 `customer wasm not supported yet`，那是运行时还未跑 guest 模块，与本示例代码无关。
