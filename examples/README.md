# SDK examples

示例演示如何把 Go 源码部署到 Vivarcus Vault。每个子目录的 **README** 里有该场景的完整命令顺序；本文是总览。

## 部署入门

### 1. 连接 Vault（地址 + 身份）

部署**需要**知道 **Vault 地址**、**Vault ID**，以及已登录的 **session**（`deploy.sh` **不会**在脚本里写账号密码）。

| 项 | 说明 |
|------|------|
| **Vault 地址** | `https://<租户>.vivarcus.com`（或你的专用 endpoint） |
| **Vault ID** | 数字 ID，Admin 或 `vivarcus config` 可见 |
| **身份** | 账号密码只在 **首次** `vivarcus auth login` 时使用；之后 CLI 保存 session，或你在 `deploy.env` 里放 `VIVARCUS_TOKEN` |

**本机首次（推荐）**：

```bash
vivarcus auth login --endpoint https://<你的租户>.vivarcus.com
vivarcus config set default_vault <vault_id>
vivarcus auth status
```

**CI / 不想交互登录**：复制 [`deploy.env.example`](deploy.env.example) 为同目录下的 `deploy.env`（已 gitignore），填写 `VIVARCUS_ENDPOINT`、`VIVARCUS_VAULT`、`VIVARCUS_TOKEN`。Token 来自一次 login 或浏览器 session，**不是密码**。

另需 [vivarcus CLI](../docs/01-prerequisites.md)（支持 `sdk put`）与 Go（仅本地 `go test` 时需要）。权限见 [01-prerequisites — Vault 权限](../docs/01-prerequisites.md#vault-权限)。

### 2. 进入示例目录

在 [vivarcus-sdk](https://github.com/vivarcus/vivarcus-sdk) clone 根目录：

```bash
cd examples/<示例目录>    # 例如 examples/01-hello-action
```

在 Vivarcus monorepo 内开发时：

```bash
cd sdk/examples/<示例目录>
```

### 3. 一个 Vault 只有一棵 Go 源码树

全部示例的 `go.mod` 都是 **`github.com/acme.corp.sdkdemo`**。`sdk put` 把文件合并进这棵树。入口的相对路径和类型名互不重复（`multi-component` 用 `SetTitleShared` / `StampDemoOnEnter` / `StampDemoName`，不占用 02 / 04 / 05 的名字），所以可以按顺序部署多个示例，下一次编译仍然通过。FQN 都是 `acme.corp.sdkdemo.<类型名>`。

| 方式 | 适用 | 做什么 |
|------|------|--------|
| **A. 单文件 `sdk put`** | **01–09** | 按需 `component apply-mdl`，再 `vivarcus sdk put -f <路径>.go` |
| **B. Inbound VPK** | **[multi-component](multi-component)** | `package-vpk.sh` → `package import` → `validate` → `deploy --confirm`（增量，同一 module） |

**推荐路径**：从 [01-hello-action](01-hello-action) 的 `./deploy.sh` 开始；要按钮就跑 [02-update-field](02-update-field)；多入口再看 [multi-component](multi-component)。

### 4. 一键部署

在示例目录执行（须已完成上一节连接 Vault）：

```bash
cd examples/<示例目录>    # monorepo: cd sdk/examples/<示例目录>
./deploy.sh               # 开头会打印当前 endpoint / vault
```

脚本步骤与 `integration_test.go` 一致。

### 5. 各示例部署速查

下表是 `./deploy.sh` 里的关键步骤。源码一步都是 **replace** 本目录，不是增量 `sdk put`。

| 目录 | 先 apply-mdl | 再装上的源码 | 备注 |
|------|----------------|------------|------|
| [01-hello-action](01-hello-action) | — | `actions/noop_action.go` | 无对象，只验证编译 |
| [02-update-field](02-update-field) | `mdl/01-object.mdl` | `actions/set_title.go` | |
| [03-confirm-dialog](03-confirm-dialog) | `mdl/01-object.mdl` | `actions/confirm_dialog.go` | |
| [04-lifecycle-entry](04-lifecycle-entry) | 对象 → **源码** → 生命周期 MDL（见 README） | `actions/stamp_on_enter.go` | 规则引用 FQN 须在源码部署 **之后** |
| [05-stamp-trigger](05-stamp-trigger) | `mdl/01-object.mdl` | `triggers/stamp_name.go` | |
| [06-query-field](06-query-field) | `mdl/01-object.mdl` | `actions/query_and_stamp.go` | |
| [07-job-processor](07-job-processor) | `mdl/01-object.mdl` | `jobs/stamp_records.go` | 调度在 Admin Job Definitions |
| [08-hello-webapi](08-hello-webapi) | `mdl/01-webapigroup.mdl` | `webapis/hello_api.go` | 非 Owner 须 Permission Set |
| [09-record-workflow-action](09-record-workflow-action) | 多条 `mdl/*.mdl`（见 README 顺序） | `workflowactions/capture_participants.go` | 对象 `capture_demo__c`；workflow MDL 在源码部署之后 |
| [multi-component](multi-component) | 生命周期 MDL，对象打进 VPK | 目录内全部入口 | `package-vpk.sh --replace` |

### 6. 验证部署成功

- `./deploy.sh` 结束且无报错；`package deploy` 的 `deployment_status` 为 `deployed__v`
- README 中的 **FQN** 可 `vivarcus sdk get <FQN> -o /tmp/out.go --json`
- 自动化：在示例目录 `go test -tags=integration -count=1 -timeout 20m .`（须已登录 Vault；跳过设 `VIVARCUS_SKIP_EXAMPLES_INTEGRATION=1`）。全量：`make test-sdk-examples-integration`（monorepo 根目录）

模块布局：[docs/03-build.md](../docs/03-build.md)。编码 Agent：[AGENTS.md](../AGENTS.md)。

## 场景索引

| 目录 | 场景 |
|------|------|
| **[multi-component](multi-component)** | 多 Action/Trigger + `actions/`/`triggers/`/`shared/`，**combo VPK** |
| [01-hello-action](01-hello-action) | 最小空操作（`sdk put`） |
| [02-update-field](02-update-field) | 记录页按钮 + `go test`（`sdk put`） |
| [03-confirm-dialog](03-confirm-dialog) | 确认框 + 横幅（`sdk put`） |
| [04-lifecycle-entry](04-lifecycle-entry) | entry / event / workflow / cancel（`sdk put` + 生命周期 MDL） |
| [05-stamp-trigger](05-stamp-trigger) | Record Trigger（`sdk put`） |
| [06-query-field](06-query-field) | 记录页按钮 + 只读 VQL（`sdk put`） |
| [07-job-processor](07-job-processor) | Job Processor（`sdk put` + 对象 MDL） |
| [08-hello-webapi](08-hello-webapi) | Custom Web API（`sdk put` + Webapigroup MDL） |
| [09-record-workflow-action](09-record-workflow-action) | Record Workflow Action（`sdk put` + 生命周期/工作流 MDL） |

编号示例把入口放在命名子目录（`actions/`、`triggers/`、`jobs/`、`webapis/`、`workflowactions/`），根目录不放客户 `.go`。`sdk put` 把该文件合并进同一 module 的源码树并整树重编译。
