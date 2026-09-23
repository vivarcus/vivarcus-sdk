# 05-stamp-trigger

**Record Trigger** 示例：`BEFORE_INSERT` 给 `name__v` 追加 `-trig`。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/05-stamp-trigger`（monorepo：`cd sdk/examples/05-stamp-trigger`）。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f triggers/stamp_name.go --json
```

FQN：`acme.corp.sdkdemo.StampName`。创建一条 `sdk_demo__c` 时 Trigger 自动跑。

模块布局：[docs/03-build.md](../../docs/03-build.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
