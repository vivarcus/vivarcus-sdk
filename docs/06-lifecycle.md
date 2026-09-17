# 生命周期与 Workflow 中的 Action

SDK Record Action 可挂在 **记录页按钮**、**进状态（entry）**、**事件（event）**、**workflow 步骤**、**workflow 取消** 等场景。

**完整跟做示例**：[examples/04-lifecycle-entry](../examples/04-lifecycle-entry)（Step 1–7：构建 → deploy → 四种绑法 → 验证）。

## 快速对照

| 场景 | `Meta.Usages` | Vault 里怎么挂 | 示例步骤 |
|------|---------------|----------------|----------|
| 记录详情 All Actions | `UserAction` | Objectaction | [02-update-field](../examples/02-update-field) |
| 生命周期状态页按钮 | `LifecycleUserAction` | user_action：`Objectaction.*` | [06-lifecycle §按钮](06-lifecycle.md) |
| **进状态时自动跑** | `LifecycleEntryAction` | entry_action rule | [04 Step 4–5](../examples/04-lifecycle-entry) |
| **create_record 等事件** | `EventAction` | event_action rule | [04 Step 7](../examples/04-lifecycle-entry) |
| **Workflow Action 步骤** | `WorkflowStep` | step_detail rule | [04 Step 6](../examples/04-lifecycle-entry) |
| **Workflow 取消** | `WorkflowCancel` | cancelation rule | [04 Step 8](../examples/04-lifecycle-entry) |

同一 wasm 可在 `Usages` 里声明多个系统场景（`UserBulkAction` 除外）。示例代码见 [04-lifecycle-entry/main.go](../examples/04-lifecycle-entry/main.go)。

## Go 代码（系统路径）

```go
Usages: []action.Usage{
    action.UsageLifecycleEntryAction,
    action.UsageEventAction,
    action.UsageWorkflowStep,
    action.UsageWorkflowCancel,
},
```

## 两种挂法：`Recordaction` vs `Objectaction`

### 方式 A — `Recordaction.<FQN>`（entry / event / workflow 常用）

```xml
<action type="Recordaction.com.example.StampOnEnter"/>
```

| 模板 | 用途 |
|------|------|
| [04-entry-action-rule.mdl](../templates/mdl/04-entry-action-rule.mdl) | entry_action |
| [05-event-action-rule.mdl](../templates/mdl/05-event-action-rule.mdl) | event_action |
| [06-workflow-action-step.mdl](../templates/mdl/06-workflow-action-step.mdl) | workflow Action 步 |
| [07-workflow-cancel-rule.mdl](../templates/mdl/07-workflow-cancel-rule.mdl) | workflow 取消 |

### 方式 B — `Objectaction.<obj>.<act>`

deploy 时会创建 active 的 Recordaction；`Meta.ObjectAction` 非空时还有 active Objectaction。**user_action 状态页按钮只能用此方式**。

## UI 钩子

| Usage | OnPreExecute / OnPostExecute |
|-------|------------------------------|
| `UserAction` / `LifecycleUserAction` | ✅ |
| `LifecycleEntryAction` / `EventAction` / `WorkflowStep` / `WorkflowCancel` | ❌ |

## 推荐学习路径

1. [02-update-field](../examples/02-update-field) + [99-full-stack](../examples/99-full-stack) — 按钮
2. [04-lifecycle-entry](../examples/04-lifecycle-entry) — entry → event → workflow → cancel

Record Action API：[02-record-action](02-record-action.md)；部署：[05-deploy](05-deploy.md)。
