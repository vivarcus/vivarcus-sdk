# SDK examples

| 目录 | 场景 | 本地 build |
|------|------|------------|
| **[multi-component](multi-component)** | **推荐**：多 Action/Trigger + `actions/`/`triggers/`/`shared/`，一个 VPK | `cd multi-component && vivarcus-sdk build .` |
| [01-hello-action](01-hello-action) | 工具链 | `cd 01-hello-action && vivarcus-sdk build . --skip-compile` |
| [02-update-field](02-update-field) | 记录页按钮 | `cd 02-update-field && vivarcus-sdk build .` |
| [03-confirm-dialog](03-confirm-dialog) | 确认框 + 横幅 | `cd 03-confirm-dialog && vivarcus-sdk build .` |
| [04-lifecycle-entry](04-lifecycle-entry) | entry / event / workflow / cancel | `cd 04-lifecycle-entry && vivarcus-sdk build .` |
| [05-stamp-trigger](05-stamp-trigger) | Record Trigger | `cd 05-stamp-trigger && vivarcus-sdk build .` |
| [06-query-field](06-query-field) | 记录页按钮 + 只读 VQL | `cd 06-query-field && vivarcus-sdk build .` |
| [07-job-processor](07-job-processor) | Job Processor | `cd 07-job-processor && vivarcus-sdk build .` |
| [08-hello-webapi](08-hello-webapi) | Custom Web API | `cd 08-hello-webapi && vivarcus-sdk build .` |
| [09-record-workflow-action](09-record-workflow-action) | Record Workflow Action | `cd 09-record-workflow-action && vivarcus-sdk build .` |

打包 VPK、import / validate / deploy 见 [docs/04-package-vpk.md](../docs/04-package-vpk.md) 与 [docs/05-deploy.md](../docs/05-deploy.md)。

脚本：[`_shared/scripts/package-vpk.sh`](_shared/scripts/package-vpk.sh)
