# 03-confirm-dialog

`OnPreExecute` 确认框 + `OnPostExecute` 横幅（需 UI 手测横幅）。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/03-confirm-dialog`（monorepo：`cd sdk/examples/03-confirm-dialog`）。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/confirm_dialog.go --json
```

FQN：`acme.corp.sdkdemo.ConfirmDialog`。按钮 `confirm_update__c` 出现在 `sdk_demo__c` 记录页 **All Actions**。

模块布局：[docs/03-build.md](../../docs/03-build.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
