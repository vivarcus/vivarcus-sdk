# 对象工作流绑定

`Recordworkflowaction` 组件由 gosdk VPK deploy 从 `__sdk_describe` 的 `workflow_actions[]` 投影，不必手写进 `components/`。

在对象工作流 **Start** 步骤的 Participant Control 上引用 FQN（`go.mod` module + 类型名，本例如下）：

```xml
<vwf:recordWorkflowAction reference="Recordworkflowaction.acme.corp.recordworkflowaction.CaptureParticipants" />
```

端到端打包见 [multi-component](../../multi-component)。
