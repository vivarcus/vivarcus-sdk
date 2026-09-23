# 02-update-field

记录页 **All Actions** 按钮：通过 `SetValue` + `platform.Update` 写入 `title__c`。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/02-update-field`（monorepo：`cd sdk/examples/02-update-field`）。

**一键**（推荐）：`./deploy.sh`

本地单测：`go test ./...`（mock `platform.UpdateRecordFunc`，见 [actions/set_title_test.go](actions/set_title_test.go)）。

`./deploy.sh` 会 apply 对象 MDL，再 `sdk put` `actions/set_title.go`。module 与其它示例相同，类型名 `SetTitle` 不与 `multi-component` 的 `SetTitleShared` 冲突。

FQN：`acme.corp.sdkdemo.SetTitle`。对象 / 按钮名写在 [actions/set_title.go](actions/set_title.go) 的 `Meta()` 与 [`mdl/`](mdl/) 里。

核对投影：`vivarcus sdk get acme.corp.sdkdemo.SetTitle -o /tmp/SetTitle.go --json`

模块布局：[docs/03-build.md](../../docs/03-build.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
