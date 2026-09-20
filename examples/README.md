# 示例

| 目录 | 场景 | 一键 demo |
|------|------|-----------|
| **[multi-component](multi-component)** | **推荐**：多 Action/Trigger + go.mod + shared/ | `cd multi-component && ./scripts/dev-demo.sh` |
| [01-hello-action](01-hello-action) | 工具链 | `vivarcus-sdk build . --skip-compile` |
| [02-update-field](02-update-field) | 记录页按钮 | `cd 02-update-field && ./scripts/dev-demo.sh` |
| [03-confirm-dialog](03-confirm-dialog) | 确认框 + 横幅 | `cd 03-confirm-dialog && ./scripts/dev-demo.sh` |
| [04-lifecycle-entry](04-lifecycle-entry) | entry / event / workflow / cancel | `cd 04-lifecycle-entry && ./scripts/dev-demo.sh` |
| [05-stamp-trigger](05-stamp-trigger) | Record Trigger 工具链 | `cd 05-stamp-trigger && vivarcus-sdk build .` |

先 `vivarcus auth login` 并 `vivarcus config set default_vault <uuid>`，再进入对应目录运行。

**Agent / CI**：注入 `VIVARCUS_TOKEN` 等环境变量复用 session，勿每条命令前密码 login（限 4 次/分钟/IP+用户）。见 [01-prerequisites](../docs/01-prerequisites.md)。

文档见 [README.md](../README.md)；Agent 见 [AGENTS.md](../AGENTS.md)。
