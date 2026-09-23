# Vivarcus SDK

在 Vivarcus Vault 上开发 **Record Action**（记录页自定义按钮）、**Record Trigger**、**Job Processor**、**Custom Web API** 与 **Record Workflow Action** 的 **Go** 开发套件。编号示例用 **`vivarcus sdk put`** 部署单个 `.go`；多文件树用 **Inbound VPK**（仅 [multi-component](examples/multi-component)）。**Vault 编译源码**。

本仓库是 guest API、模板与示例。扫描 / codegen / tinygo 在 Vault 镜像里，不随本模块发布。

## 5 分钟 Quickstart

### 1. 准备工具

| 工具 | 用途 |
|------|------|
| [Go](https://go.dev/dl/)（见 `sdk/go.mod`） | 编写 Action / Trigger，本地 `go test` |
| [vivarcus CLI](docs/01-prerequisites.md) | MDL、`sdk put`、把 VPK 部署到 Vault |

```bash
git clone https://github.com/vivarcus/vivarcus-sdk.git
cd vivarcus-sdk
```

编码 Agent 请先读 [AGENTS.md](AGENTS.md)。

### 2. 最小部署：单文件 `sdk put`

从 **[examples/01-hello-action](examples/01-hello-action)** 开始。先连接 Vault（地址、Vault ID；账号密码仅用于 `vivarcus auth login`，见 [examples/README — 连接 Vault](examples/README.md)）：

```bash
vivarcus auth login --endpoint https://<你的租户>.vivarcus.com
vivarcus config set default_vault <vault_id>
cd examples/01-hello-action
./deploy.sh
```

有对象的编号示例：先 `vivarcus component apply-mdl --confirm -f mdl/01-object.mdl`，再 `sdk put`。详见 [examples/README.md — 对象 MDL](examples/README.md#对象-mdl)。

### 3. 多文件：combo VPK（仅 multi-component）

**[examples/multi-component](examples/multi-component)** 才打 VPK——一个 module、多个 Action + Trigger：

```bash
cd examples/multi-component
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:multi_demo__c:mdl/01-object.mdl
# vivarcus package import → validate → deploy（见 docs/05-deploy.md）
```

涵盖本地 `go.mod` + `actions/` / `triggers/` / `shared/` 布局、对象 MDL、VPK 打包（只上传 `.go`）。平台编译 `gosdk/`，不要往包里放 wasm。不要用这套命令部署 01–09。

### 4. 写 Action

从模板开始：

```bash
mkdir -p my-action/actions
cp templates/action/action.go.tpl my-action/actions/set_title.go
# 编辑 my-action/actions/set_title.go，实现 Meta / IsExecutable / Execute。包名是 package actions。
cd my-action && go test ./...
```

本地反馈优先 `go test`（mock `platform.*`），对照 [examples/02-update-field](examples/02-update-field)。

分场景精读示例：

| 示例 | 路径 | 说明 |
|------|------|------|
| Hello Action | [examples/01-hello-action](examples/01-hello-action) | 最小空操作（`sdk put`） |
| 更新字段 | [examples/02-update-field](examples/02-update-field) | 记录页按钮（`sdk put`） |
| 确认对话框 | [examples/03-confirm-dialog](examples/03-confirm-dialog) | `OnPreExecute` / `OnPostExecute`（`sdk put`） |
| **系统自动执行** | [examples/04-lifecycle-entry](examples/04-lifecycle-entry) | entry / event / workflow / cancel（`sdk put` + MDL） |
| Record Trigger | [examples/05-stamp-trigger](examples/05-stamp-trigger) | 单一 Trigger（`sdk put`） |
| Job Processor | [examples/07-job-processor](examples/07-job-processor) | Init / Process（`sdk put` + 对象 MDL） |
| Custom Web API | [examples/08-hello-webapi](examples/08-hello-webapi) | Custom Web API（`sdk put` + Webapigroup MDL） |
| Record Workflow Action | [examples/09-record-workflow-action](examples/09-record-workflow-action) | GET_PARTICIPANTS（`sdk put` + MDL） |

### 5. 准备 Vault 元数据（MDL）

部署前，目标对象须已存在。复制并修改模板：

```bash
# 见 templates/mdl/README.md
vivarcus component apply-mdl --confirm -f templates/mdl/01-object.mdl
```

多文件才把对象 MDL 放进 combo VPK（见 multi-component）。

### 6. 部署

- **一个 `.go`**：`vivarcus sdk put -f actions/noop_action.go --json`（见 [05-deploy](docs/05-deploy.md) 与 [01-hello-action](examples/01-hello-action)）。
- **多文件树**：仅 [multi-component](examples/multi-component) 打 VPK。`gosdk/` 只放 `.go`，编译在 Vault 上完成。

### 7. 验证

- 按钮：记录详情 **All Actions**
- Trigger：创建记录时自动执行（multi-component demo 会 API 验收）
- Job Processor：Admin > Operations 调度 SDK Job
- Agent：`vivarcus sdk get <FQN> -o /tmp/src.go --json`；编号示例部署用 `vivarcus sdk put -f <子目录>/<类型名>.go --json`；启停见 [05-deploy](docs/05-deploy.md)

## 文档

| 文档 | 内容 |
|------|------|
| [01-prerequisites](docs/01-prerequisites.md) | 环境、权限、CLI 配置 |
| [02-record-action](docs/02-record-action.md) | API：Meta、Execute、可选钩子 |
| [03-build](docs/03-build.md) | 模块布局与平台如何编译 |
| [04-package-vpk](docs/04-package-vpk.md) | VPK 目录结构与 `vaultpackage.xml` |
| [05-deploy](docs/05-deploy.md) | 01–09：`sdk put`；multi-component：import → validate → deploy |
| [06-lifecycle](docs/06-lifecycle.md) | entry / workflow / 按钮：Usages 与两种 rule 挂法 |
| [07-limits](docs/07-limits.md) | 标准库边界、配额、Phase 1 能力 |
| [08-job-processor](docs/08-job-processor.md) | Job Processor：Meta / Init / Process |
| [troubleshooting](docs/troubleshooting.md) | 常见错误 |

## Go API 包

| 包 | 说明 |
|----|------|
| `action` | Record Action 接口与上下文 |
| `trigger` | Record Trigger 接口与上下文 |
| `job` | Job Processor 接口与上下文 |
| `webapi` | Custom Web API 接口与上下文 |
| `workflowaction` | Record Workflow Action 接口与上下文 |
| `platform` | 宿主能力：`Get` / `Update` / `Log` |
| `wire` | ABI 编解码（一般由平台 codegen 使用） |

```go
import (
    "github.com/vivarcus/vivarcus-sdk/action"
    "github.com/vivarcus/vivarcus-sdk/platform"
)
```

## 编码 Agent

使用 Cursor / Copilot 等 Agent 时，请先读根目录 [AGENTS.md](AGENTS.md)。

## License

Apache License 2.0 — 见 [LICENSE](LICENSE)。
