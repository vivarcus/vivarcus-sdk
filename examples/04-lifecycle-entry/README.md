# 04-lifecycle-entry

**系统自动调用的 Record Action**：entry_action、event_action、workflow action step、workflow cancel。

与 [02-update-field](02-update-field) 的差别：

1. Go 里 `Meta.Usages` 用 `LifecycleEntryAction` / `EventAction` / `WorkflowStep` / `WorkflowCancel`（不是 `UserAction`）
2. 管理员在 lifecycle / workflow 规则里挂 **`Recordaction.<FQN>`**（本示例）或 **`Objectaction.<obj>.<act>`**（若还要状态页按钮）

> `OnPreExecute` / `OnPostExecute` **不会**在这些路径执行；只有用户点按钮（`UserAction` / `LifecycleUserAction`）才有确认框。

## 什么时候用本示例

| 场景 | 用哪个示例 | 本 README 步骤 |
|------|------------|----------------|
| 记录详情 **All Actions** 按钮 | [02-update-field](02-update-field) | — |
| 进某状态时自动改字段 | **本示例** | Step 3–4 |
| 创建/更新记录事件自动跑 | **本示例** | Step 6 |
| Workflow Action 步骤里跑 SDK | **本示例** | Step 5 |
| Workflow 取消时跑 SDK | **本示例** | Step 7 |
| 状态页按钮 + 系统自动跑 | 02 的代码 + `LifecycleUserAction`，Objectaction | Step 3 方式 B |

## Quick start（目标 Vault，一键跑通）

1. 配置 CLI（指向你的 Vivarcus Vault）：

```bash
vivarcus auth login
vivarcus config set default_vault <vault-uuid>
```

2. 阅读 **MDL 组件**（`mdl/` 目录，与 Step 3–7 一一对应）：

```bash
cat mdl/README.md
```

3. 进入本示例目录，一键预置 + 验证：

```bash
cd examples/04-lifecycle-entry
./scripts/dev-demo.sh           # lifecycle MDL + combo VPK（对象+gosdk）+ workflow MDL + 验证
./scripts/dev-demo.sh --setup-only
./scripts/dev-demo.sh --verify-only
```

脚本**不包含**组件定义：对象 / lifecycle / workflow 写在 [`mdl/*.mdl`](mdl/) 中。Recordaction FQN 在 `vivarcus-sdk build` 之后由 `vivarcus-sdk describe action.wasm` 填入。

> **验证 event 须走 CAP-OM 创建**：`POST /api/v1/objects/{obj}/records`。`vivarcus object create` 不会触发 `create_record` event_action。

> 对象上已有记录时，`03-bind-object-lifecycle.mdl` 可能失败（脚本会跳过并继续）。首次跑建议在空对象上。

下文 Step 1–7 为手工分步说明（Admin UI / MDL 逐项配置）。

## Prerequisites

| 项 | 值 |
|----|-----|
| 对象 | `demo_request__c`，含字段 `title__c`（[01-object.mdl](../../templates/mdl/01-object.mdl)） |
| 生命周期 | `demo_request_lc__c`，含 `draft__c`、`in_review__c`（[02-lifecycle.mdl](../../templates/mdl/02-lifecycle.mdl)） |
| 工具 | Go、TinyGo、`vivarcus-sdk`、`vivarcus` CLI |

## Step 1 — 构建 wasm

```bash
vivarcus-sdk build . -o action.wasm
vivarcus-sdk describe action.wasm
# component_name 即为 Recordaction 名，由 go.mod + 类型名派生
```

若 TinyGo 报 `requires go version 1.19 through 1.26`，在命令前加 `GOTOOLCHAIN=go1.26.2`。

FQN 来自 `Meta.Name`（若设置）否则 `github.com/vivarcus/vivarcus-sdk` → `com.example.StampOnEnter`。系统路径直连 **不需要** Objectaction；`Meta.Usages` 应含 `LifecycleEntryAction` / `EventAction` / `WorkflowStep` / `WorkflowCancel`。

## Step 2 — 部署 VPK

使用 [`_shared/scripts/package-vpk.sh`](../_shared/scripts/package-vpk.sh)（对象 MDL 可选 `--component`）→ `import` → `validate` → `deploy --confirm`。gosdk 步会创建 **active** 的 Recordaction（`Meta.ObjectAction` 非空时还有 Objectaction）。

## Step 3 — entry_action（进状态时自动跑）

在目标 lifecycle **状态** 的 **Entry Actions** 加 rule：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rule>
  <actions>
    <action type="Recordaction.com.example.StampOnEnter"/>
  </actions>
</rule>
```

- 绑定到 **进入 `in_review__c` 时** 触发的 entry action
- 模板：[04-entry-action-rule.mdl](../../templates/mdl/04-entry-action-rule.mdl)

**方式 B（还要状态页按钮）**：rule 改用 `Objectaction.demo_request__c.stamp_on_enter__c`（deploy 时已创建 active Objectaction）。

## Step 4 — 验证 entry

1. 创建 `demo_request__c` 记录（`draft__c`）
2. 用 user action 切到 `in_review__c`（或你们环境中等价操作）
3. 进入 `in_review__c` 后，`title__c` → `stamped-on-enter`

---

## Step 5 — workflow action step

### 5a. 建最小 workflow

在 Vault 配置 **Objectworkflow**（绑定 `demo_request_lc__c`）：

| 项 | 建议值 |
|----|--------|
| 名称 | `stamp_on_enter_wf__c`（示例） |
| cardinality | One |
| auto_start | false（用 user action 手动启动便于验证） |
| 步骤 | `start` → `stamp__c`（type=**action**）→ `end` |

### 5b. Action 步骤 step_detail

在 `stamp__c` 步骤粘贴 Perform actions rule（完整 XML 见 [06-workflow-action-step.mdl](../../templates/mdl/06-workflow-action-step.mdl)）：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<vwf:stepdetails xmlns:vwf="VeevaVault">
  <vwf:rules>
    <vwf:rule>
      <actions>
        <action type="Recordaction.com.example.StampOnEnter"/>
      </actions>
    </vwf:rule>
  </vwf:rules>
</vwf:stepdetails>
```

另建 lifecycle **user_action**：`LIFECYCLE_RUN_WORKFLOW_ACTION` 指向该 workflow（仅用于启动验证）。

### 5c. 验证 workflow

1. 创建 `demo_request__c` 记录
2. 点「Start Workflow」类 user action 启动 `stamp_on_enter_wf__c`
3. workflow 跑过 action 步后，`title__c` → `stamped-on-enter`

---

## Step 6 — event_action（创建记录等事件）

在 lifecycle **Event Actions** 增加一条：

| 字段 | 值 |
|------|-----|
| event | `create_record` |
| rule | 见下方 XML |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rule>
  <actions>
    <action type="Recordaction.com.example.StampOnEnter"/>
  </actions>
</rule>
```

模板：[05-event-action-rule.mdl](../../templates/mdl/05-event-action-rule.mdl)

### 验证 event

1. **新建**一条 `demo_request__c`（触发 `create_record`）
2. 创建完成后 `title__c` 应已为 `stamped-on-enter`（无需切状态、无需点按钮）

> 若同时配置了 Step 3 entry 与 Step 6 event，同一次创建可能触发多次；验收时建议 **只开一种** 绑法。

---

## Step 7 — workflow cancel

在 workflow **Cancelation Actions** 粘贴 rule（[07-workflow-cancel-rule.mdl](../../templates/mdl/07-workflow-cancel-rule.mdl)）：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rule>
  <actions>
    <action type="Recordaction.com.example.StampOnEnter"/>
  </actions>
</rule>
```

### 验证 cancel

1. 启动 Step 5 的 workflow，在结束前 **取消** 该实例
2. 取消完成后检查 `title__c` 是否被改写

---

## 两种挂法对照

| | `Recordaction.<FQN>` | `Objectaction.<obj>.<act>` |
|---|---|---|
| entry / event / workflow / cancel | ✅ 常用 | ✅ 也可 |
| user_action（状态页按钮） | ❌ 不支持 | ✅ 必须 |
| 需要 Objectaction 组件 | **否** | entry 可选；按钮 **是** |
| deploy 后组件状态 | active（Recordaction） | active（Recordaction + Objectaction） |

## 路径速查

| 路径 | `Meta.Usage` | rule 写在哪 |
|------|--------------|-------------|
| entry | `LifecycleEntryAction` | lifecycle 状态 Entry Actions |
| event | `EventAction` | lifecycle Event Actions |
| workflow step | `WorkflowStep` | workflow Action 步 step_detail |
| cancel | `WorkflowCancel` | workflow Cancelation Actions |

更多概念说明：[docs/06-lifecycle.md](../../docs/06-lifecycle.md)。
