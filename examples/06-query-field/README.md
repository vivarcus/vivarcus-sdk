# 06-query-field

记录页按钮：用 `platform.Query` 只读 VQL 统计同对象记录数，再把 `title__c` 写成 `query:<N>`。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/06-query-field`（monorepo：`cd sdk/examples/06-query-field`）。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/query_and_stamp.go --json
```

FQN：`acme.corp.sdkdemo.QueryAndStamp`。对象 / 按钮名见 [actions/query_and_stamp.go](actions/query_and_stamp.go) 与 [`mdl/`](mdl/)。验收期望 `title__c=query:1`（单条演示记录）。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
