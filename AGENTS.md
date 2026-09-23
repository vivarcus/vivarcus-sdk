# Agent 开发手册 — Vivarcus SDK

本文件供编码 Agent（Cursor、Copilot 等）读取。人类开发者请从 [README.md](README.md) 开始。

## 仓库定位

- **目标**：帮助开发者在 Vivarcus Vault 上实现 **Record Action**、**Record Trigger**、**Job Processor**、**Custom Web API** 与 **Record Workflow Action**。
- **语言**：Go（**不是 Java**）。平台把客户 `.go` 编成 wasm。
- **部署路径**：
  - **编号示例 01–09**：`vivarcus sdk put -f <子目录>/<类型名>.go`（有对象则先 `component apply-mdl`）。入口在命名子目录，例如 `actions/set_title.go`（`package actions`）。
  - **仅 [multi-component](examples/multi-component)**：`package-vpk.sh` → `vivarcus package import/validate/deploy`。
- **默认状态**：客户入口部署后 **active**，记录页按钮完成后即可见。

## 默认循环（先做这个）

1. 对照模板实现接口；`Meta` 与对象 / 按钮 api_name 对齐。
2. `go test ./...`：直接调 `Execute` / `Process`，mock `platform.*` 的 `*Func`（见 [examples/02-update-field/actions/set_title_test.go](examples/02-update-field/actions/set_title_test.go)）。`*_test.go` **不**部署。
3. 编号示例：对象尚不存在时 `vivarcus component apply-mdl --confirm -f mdl/01-object.mdl`，然后 `vivarcus sdk put -f actions/set_title.go --json`（路径按示例子目录）。`sdk_demo__c` 已存在则不要再 apply 这条 `RECREATE`，只 `sdk put`。
4. 多文件树才打 VPK：只做 [examples/multi-component](examples/multi-component)。编译失败时读 `gosdk_invalid` 改源码再 `put` 或 re-deploy。
5. 在 Vault 验证按钮 / Trigger / Job / Web API。

不要在本机编 wasm。不要给 01–09 打 VPK。模块布局见 [docs/03-build.md](docs/03-build.md)。

## 任务路由

| 用户意图 | 先读 | 再执行 |
|----------|------|--------|
| 新建 Record Action | `templates/action/action.go.tpl` | 实现接口 → `go test ./...` → `sdk put -f actions/<类型名>.go` |
| 新建 Job Processor | `templates/job/job.go.tpl` | 实现 Meta/Init/Process → `go test` → `sdk put` |
| 新建 Record Workflow Action | `templates/workflowaction/workflow_action.go.tpl` | 实现 Meta/Execute → `go test` → `sdk put` |
| 按钮改字段 | `examples/02-update-field/` | `apply-mdl` + `sdk put`；对齐 `Meta.Object` |
| 执行前确认框 | `examples/03-confirm-dialog/` | `apply-mdl` + `sdk put`；实现 `OnPreExecute` |
| 记录页按钮（多文件端到端） | `examples/multi-component/` | **仅此示例打 VPK** |
| **系统自动执行（entry/event/workflow/cancel）** | **`examples/04-lifecycle-entry/`** | `sdk put` 后再绑 lifecycle MDL |
| 生命周期概念与两种挂法 | `docs/06-lifecycle.md` | 对照表 + 链到 04 示例 |
| 创建对象 | `templates/mdl/01-object.mdl` | `vivarcus component apply-mdl`（需 Vault Owner） |
| 打包 VPK | `docs/04-package-vpk.md` | **仅 multi-component**；`gosdk/` + 可选 `components/` |
| 部署到 Vault | `docs/05-deploy.md` | 01–09：`sdk put`；multi-component：`package import/validate/deploy` |
| 下载已部署客户源码 / 单文件部署 / 启停入口 | `docs/05-deploy.md`（`vivarcus sdk get\|put\|enable\|disable`） | 编号示例命令见各自 README；不要 `curl /code` |
| 自定义代码执行失败 / `platform.LogInfo` 对不上 | `vivarcus sdk logs`（见 `vivarcus-cli` SKILL） | UTC 日 ZIP；见 troubleshooting Runtime Log |

## 端到端检查清单

完成一个 Record Action 时，按顺序确认：

1. **对象存在**：`templates/mdl/01-object.mdl` 中 `{{OBJECT}}` 已创建，字段 api_name 与 Go 代码一致。
2. **Action 代码**：实现 `Meta()`、`IsExecutable()`、`Execute()`；同一 module 可有多个类型；可选 `OnPreExecute`/`OnPostExecute`。
3. **Meta 对齐**：`Meta.Object` / `Meta.ObjectAction` / `Meta.Label` 与 Vault 对象和按钮 api_name 一致；FQN 由 module+类型名派生（或 `Meta.Name`）。
4. **本地**：`go test` 覆盖主路径；不要 import `os`/`net`/`unsafe` 等（见 `docs/07-limits.md`）。
5. **部署（编号示例）**：`vivarcus sdk put -f <子目录>/<类型名>.go --json`。有对象先 `component apply-mdl`。不要 `package-vpk.sh`。
6. **部署（仅 multi-component）**：`package-vpk.sh` → `validate` → `deploy --confirm`。
7. **验证**：
   - 按钮：记录详情 **All Actions**
   - 系统路径：见 [04-lifecycle-entry](examples/04-lifecycle-entry)（entry / event / workflow / cancel 分步）
   - 投影：`vivarcus sdk get <FQN> -o /tmp/<Type>.go --json`（FQN 写在各示例 README）
   - 止血：`vivarcus sdk disable <FQN> --json`，验证完再 `enable`。不要对 `Sdkcode` helper 调 enable/disable
   - 运行时日志：`vivarcus sdk logs --date $(date -u +%Y-%m-%d) -o /tmp/SdkLog.zip --json`，再 `unzip -p` 看 `platform.Log*` / EXCEPTION

## 占位符约定

模板文件使用 `{{NAME}}` 占位符，Agent 替换时须全局一致：

| 占位符 | 含义 | 示例 |
|--------|------|------|
| `{{OBJECT}}` | 自定义对象 api_name | `demo_request__c` |
| `{{ACTION_FQN}}` | Recordaction 名 | `com.example.SetTitle`（module 路径 + 类型名；需要固定时写 `Meta.Name`） |
| `{{OBJECT_ACTION}}` | Objectaction api_name | `demo_request__c.approve__c` |
| `{{LABEL}}` | 按钮显示名 | `Approve Request` |
| `{{SHA256}}` | （已废弃）客户不再上传 wasm | 平台编译产出 wasm |

`component_name`（FQN）默认由 module 路径 + 类型名派生；需要固定名字时在 `Meta.Name` 显式写出。

## vivarcus CLI 认证（避免登录限流）

部署与 MDL 命令依赖 `vivarcus`。**不要**每条命令前密码 login：服务端对 `POST /ui/auth/login` 限 **10 次/分钟/IP+用户**。整段任务注入并复用 `VIVARCUS_TOKEN`（+ `VIVARCUS_ENDPOINT`、`VIVARCUS_VAULT`），详见 [01-prerequisites](docs/01-prerequisites.md)。

## 禁止事项

- **不要用 Java** 或 `javasdk/` VPK（平台返回 `not_supported__v`）。
- **不要**在本机编 wasm，也不要把 `.wasm` 打进 VPK。
- **不要**在自动化循环中反复 `curl /ui/auth/login` 或每条 `vivarcus` 前重新 login（会 429）。
- **不要** import `os`、`net`、`database/sql`、`unsafe` 等（平台编译拒绝）。`testing` 只放 `*_test.go`。
- **不要** 跳过 `validate` 直接 `deploy`（仅 VPK / multi-component）。
- **不要** 给 01–09 打 VPK；那些示例用 `sdk put`。
- **不要** 假设部署后按钮自动可见（须 `active(true)`）。
- Phase 1 支持 **只读 VQL**（`platform.Query`）；无通知、HTTP、生命周期切换 host API（见 `docs/07-limits.md`）。

## API 速查

```go
// 必选
func (T) Meta() action.Meta
func (T) IsExecutable(ctx action.RecordActionContext) bool
func (T) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error)

// Job Processor
func (T) Meta() job.Meta
func (T) Init(ctx job.InitContext) (job.Input, error)
func (T) Process(ctx job.ProcessContext) (job.ProcessResult, error)

// Record Workflow Action
func (T) Meta() workflowaction.Meta
func (T) Execute(ctx workflowaction.RecordWorkflowActionContext) error

// 可选（实现即由平台编译导出 wasm 钩子）
func (T) OnPreExecute(ctx action.RecordActionContext) (action.PreExecuteResult, error)
func (T) OnPostExecute(ctx action.RecordActionContext) (action.PostExecuteResult, error)
```

```go
// 宿主能力（Phase 1）
platform.Get(object, recordID)
platform.Update(object, recordID, map[string]any{...})
platform.LogInfo(msg)
rec.SetValue(field, value) // Execute 内暂存，配合 platform.Update 持久化
```

本地 `go test` 给 `platform.UpdateRecordFunc` / `GetRecordFunc` / `LogFunc` 赋值即可 mock，无需 wasm。

## 示例索引

| 目录 | 场景 |
|------|------|
| `examples/01-hello-action` | 最小空操作（`sdk put`） |
| `examples/02-update-field` | 读写记录字段（UserAction 按钮）+ `go test`（`sdk put`） |
| `examples/03-confirm-dialog` | 确认框 + 结果横幅（`sdk put`） |
| `examples/04-lifecycle-entry` | entry / event / workflow step / cancel（`sdk put` + MDL） |
| `examples/05-stamp-trigger` | Record Trigger（`sdk put`） |
| `examples/06-query-field` | 记录页按钮 + 只读 VQL（`sdk put`） |
| `examples/07-job-processor` | Job Processor（`sdk put` + 对象 MDL） |
| `examples/08-hello-webapi` | Custom Web API（`sdk put` + Webapigroup MDL） |
| `examples/09-record-workflow-action` | Record Workflow Action（`sdk put` + 生命周期/工作流 MDL） |
| `examples/multi-component` | **唯一 VPK 示例**：多 Action/Trigger + MDL + combo 包 |

## 反馈

问题与功能请求请通过 Vivarcus 客户支持渠道联系，不要在本仓库提交 PR。
