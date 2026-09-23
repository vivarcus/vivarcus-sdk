# 01-hello-action

最小 Record Action — `Execute` 为空操作。用来确认「一个 `.go` → 部署 → Vault 编译」通路，再做 [02-update-field](../02-update-field) 或 [04-lifecycle-entry](../04-lifecycle-entry)。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/01-hello-action`（monorepo：`cd sdk/examples/01-hello-action`）。

**一键**（推荐）：

```bash
./deploy.sh
```

本示例**不需要** `apply-mdl`。手动等价：

```bash
vivarcus sdk put -f actions/noop_action.go --json
```

`--path` / `--module` 从本目录 `go.mod` 推断（`github.com/acme.corp.sdkdemo`，与其它示例相同）。成功 `responseMessage` 为 `Modified file`。平台整树重编译，入口默认 **active**。

FQN：`acme.corp.sdkdemo.NoopAction`。Agent 核对（整段复用 `VIVARCUS_TOKEN`，不要每条命令 login）：

```bash
vivarcus sdk get acme.corp.sdkdemo.NoopAction -o /tmp/NoopAction.go --json
vivarcus sdk disable acme.corp.sdkdemo.NoopAction --json
vivarcus sdk enable acme.corp.sdkdemo.NoopAction --json
```

`get` 应能看到 `package actions` / `NoopAction`。`disable` 后不要停住，验证完再 `enable`。

本示例 `Meta` 未绑对象，记录页不会出现按钮；只验证编译与投影。绑对象 + 按钮见 [02-update-field](../02-update-field)。

模块布局：[docs/03-build.md](../../docs/03-build.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
