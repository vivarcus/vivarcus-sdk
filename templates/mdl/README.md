# MDL 模板

按顺序执行（替换 `{{...}}` 占位符后）：

| 顺序 | 文件 | 说明 |
|------|------|------|
| 1 | [01-object.mdl](01-object.mdl) | 创建自定义对象与业务字段 |
| 2 | [02-lifecycle.mdl](02-lifecycle.mdl) | 可选：生命周期与状态 |
| 3 | — | 部署 gosdk VPK（见 [docs/05-deploy.md](../../docs/05-deploy.md)） |
| 4 | [04-entry-action-rule.mdl](04-entry-action-rule.mdl) | entry_action rule（`Recordaction.<FQN>`） |
| 5 | [05-event-action-rule.mdl](05-event-action-rule.mdl) | event_action rule（如 `create_record`） |
| 6 | [06-workflow-action-step.mdl](06-workflow-action-step.mdl) | workflow Action 步 step_detail |
| 7 | [07-workflow-cancel-rule.mdl](07-workflow-cancel-rule.mdl) | workflow Cancelation Actions |

```bash
vivarcus mdl run templates/mdl/01-object.mdl
# deploy VPK ...
vivarcus mdl run templates/mdl/04-entry-action-rule.mdl
```

占位符须与 Go `Meta()`（及派生或显式的 Recordaction FQN）一致。
