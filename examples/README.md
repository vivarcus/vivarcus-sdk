# 示例

| 目录 | 场景 | 一键 demo |
|------|------|-----------|
| **[multi-component](multi-component)** | **推荐**：多 Action/Trigger + go.mod + shared/ | `make demo-multi-component` |
| [01-hello-action](01-hello-action) | 工具链 | `make build-01-hello-action` |
| [02-update-field](02-update-field) | 记录页按钮 | `make demo-02-update-field` |
| [03-confirm-dialog](03-confirm-dialog) | 确认框 + 横幅 | `make demo-03-confirm-dialog` |
| [04-lifecycle-entry](04-lifecycle-entry) | entry / event / workflow / cancel | `make demo-04-lifecycle-entry` |
| [05-stamp-trigger](05-stamp-trigger) | Record Trigger 工具链 | `make build-05-stamp-trigger` |

在 `examples/` 目录：`vivarcus auth login` → `make demo-multi-component`（或各示例 `./scripts/dev-demo.sh`）

文档见 [README.md](../README.md)；Agent 见 [AGENTS.md](../AGENTS.md)。
