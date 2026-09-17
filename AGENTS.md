# Agent 开发手册 — Vivarcus SDK

本文件供编码 Agent（Cursor、Copilot 等）读取。人类开发者请从 [README.md](README.md) 开始。

## 仓库定位

- **目标**：帮助开发者在 Vivarcus Vault 上实现 **Record Action**（记录页自定义按钮）。
- **语言**：Go → TinyGo → WebAssembly（**不是 Java**）。
- **部署路径**：`gosdk/` Inbound VPK → `vivarcus package import/validate/deploy`。
- **默认状态**：客户 Action 部署后 **active**，记录页按钮 deploy 完成即可见。

## 任务路由

| 用户意图 | 先读 | 再执行 |
|----------|------|--------|
| 新建 Record Action | `templates/action/main.go.tpl` | 实现接口 → `vivarcus-sdk build` |
| 按钮改字段 | `examples/02-update-field/` | 对齐 `Meta.Object` 与 MDL 对象 api_name |
| 执行前确认框 | `examples/03-confirm-dialog/` | 实现 `OnPreExecute` |
| 记录页按钮（端到端） | `examples/99-full-stack/` | MDL + VPK + deploy |
| **系统自动执行（entry/event/workflow/cancel）** | **`examples/04-lifecycle-entry/`** | 同一 wasm；Step 3–7 分路径绑 rule + 验证 |
| 生命周期概念与两种挂法 | `docs/06-lifecycle.md` | 对照表 + 链到 04 示例 |
| 创建对象 | `templates/mdl/01-object.mdl` | `vivarcus mdl run`（需 Vault Owner） |
| 打包 VPK | `docs/04-package-vpk.md` | `gosdk/*.wasm` + 可选 `components/` |
| 部署到 Vault | `docs/05-deploy.md` | `vivarcus package import/validate/deploy --confirm` |
| validate 失败 | `docs/troubleshooting.md` | 查 wasm import / describe `api_version` / 体积 |

## 端到端检查清单

完成一个 Record Action 时，按顺序确认：

1. **对象存在**：`templates/mdl/01-object.mdl` 中 `{{OBJECT}}` 已创建，字段 api_name 与 Go 代码一致。
2. **Action 代码**：实现 `Meta()`、`IsExecutable()`、`Execute()`；同一 module 可有多个类型；可选 `OnPreExecute`/`OnPostExecute`。
3. **Meta 对齐**：`Meta.Object` / `Meta.ObjectAction` / `Meta.Label` 与 Vault 对象和按钮 api_name 一致；FQN 由 module+类型名派生（或 `Meta.Name`）。
4. **构建**：`vivarcus-sdk build <dir> -o action.wasm` 成功。
5. **VPK 结构**：
   ```
   vaultpackage.xml
   gosdk/action.wasm
   ```
6. **部署**：`vivarcus package validate` 通过（`deployment_status` 非 `not_verified__v`），`deploy --confirm` 创建 active Recordaction / Objectaction。
7. **验证**：
   - 按钮：记录详情 **All Actions**
   - 系统路径：见 [04-lifecycle-entry](examples/04-lifecycle-entry)（entry / event / workflow / cancel 分步）

## 占位符约定

模板文件使用 `{{NAME}}` 占位符，Agent 替换时须全局一致：

| 占位符 | 含义 | 示例 |
|--------|------|------|
| `{{OBJECT}}` | 自定义对象 api_name | `demo_request__c` |
| `{{ACTION_FQN}}` | Recordaction 名 | `com.example.SetTitle`（`vivarcus-sdk describe`） |
| `{{OBJECT_ACTION}}` | Objectaction api_name | `demo_request__c.approve__c` |
| `{{LABEL}}` | 按钮显示名 | `Approve Request` |
| `{{SHA256}}` | wasm 文件 SHA-256 hex | 由 `sha256sum action.wasm` 得到 |

`component_name`（FQN）默认由 module 路径 + 类型名派生；需要固定名字时在 `Meta.Name` 显式写出。

## 禁止事项

- **不要用 Java** 或 `javasdk/` VPK（平台返回 `not_supported__v`）。
- **不要** import `os`、`net`、`database/sql`、`unsafe` 等（编译期拒绝）。
- **不要** 跳过 `validate` 直接 `deploy`。
- **不要** 假设部署后按钮自动可见（须 `active(true)`）。
- Phase 1 **无** VQL、通知、HTTP、生命周期切换 host API（见 `docs/07-limits.md`）。

## API 速查

```go
// 必选
func (T) Meta() action.Meta
func (T) IsExecutable(ctx action.RecordActionContext) bool
func (T) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error)

// 可选（实现即导出 wasm 钩子）
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

## 示例索引

| 目录 | 场景 |
|------|------|
| `examples/01-hello-action` | 工具链冒烟 |
| `examples/02-update-field` | 读写记录字段（UserAction 按钮） |
| `examples/03-confirm-dialog` | 确认框 + 结果横幅 |
| `examples/04-lifecycle-entry` | entry / event / workflow step / cancel（Step 4–8 + MDL 模板 04–07） |
| `examples/99-full-stack` | MDL + VPK 脚本一条龙（按钮） |

## 反馈

问题与功能请求请通过 Vivarcus 客户支持渠道联系，不要在本仓库提交 PR。
