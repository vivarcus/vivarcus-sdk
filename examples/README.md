# SDK examples

两种部署，不要混用：

| 示例 | 怎么部署 |
|------|----------|
| **[multi-component](multi-component)** | **唯一走 VPK**：`package-vpk.sh` → `import` → `validate` → `deploy` |
| **[01](01-hello-action)–[09](09-record-workflow-action)** | **单文件**：对象/配置用 `vivarcus component apply-mdl`，Go 用 `vivarcus sdk put -f <子目录>/<类型名>.go` |

编号示例把入口放在命名子目录（`actions/`、`triggers/`、`jobs/`、`webapis/`、`workflowactions/`），根目录不放客户 `.go`。`sdk put` 把该文件 merge 进 vault 树并整树重编译。本地 `go test ./...` 见 [02-update-field](02-update-field)。模块布局见 [docs/03-build.md](../docs/03-build.md)。

| 目录 | 场景 |
|------|------|
| **[multi-component](multi-component)** | 多 Action/Trigger + `actions/`/`triggers/`/`shared/`，**combo VPK** |
| [01-hello-action](01-hello-action) | 最小空操作（`sdk put`） |
| [02-update-field](02-update-field) | 记录页按钮 + `go test`（`sdk put`） |
| [03-confirm-dialog](03-confirm-dialog) | 确认框 + 横幅（`sdk put`） |
| [04-lifecycle-entry](04-lifecycle-entry) | entry / event / workflow / cancel（`sdk put` + 生命周期 MDL） |
| [05-stamp-trigger](05-stamp-trigger) | Record Trigger（`sdk put`） |
| [06-query-field](06-query-field) | 记录页按钮 + 只读 VQL（`sdk put`） |
| [07-job-processor](07-job-processor) | Job Processor（`sdk put`） |
| [08-hello-webapi](08-hello-webapi) | Custom Web API（`sdk put`） |
| [09-record-workflow-action](09-record-workflow-action) | Record Workflow Action（`sdk put`） |

VPK 细节：[04-package-vpk](../docs/04-package-vpk.md)、[05-deploy](../docs/05-deploy.md)。脚本仅 multi-component 需要：[`_shared/scripts/package-vpk.sh`](_shared/scripts/package-vpk.sh)
