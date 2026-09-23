# 08-hello-webapi

**Custom Web API** 示例：`POST /api/{version}/custom/hello_world`，JSON 体 `{ "name": "Ada" }` 返回 `message`。

自定义 HTTP 端点：`Meta()` + `Execute`。`Meta.APIGroup` 必须指向已存在的 `Webapigroup`（本例 `integration__c`，见 [`mdl/`](mdl/)）。非 Vault Owner 还须在 Permission Set 上授予 `webapigroup.integration__c.actions('execute')`，并已有 `vault_actions.api`。`vault_actions.api('all_api')` **不**打开该 Group。Vault Owner 可直接调用已启用端点。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/08-hello-webapi`（monorepo：`cd sdk/examples/08-hello-webapi`）。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-webapigroup.mdl
vivarcus sdk put -f webapis/hello_api.go --json
```

FQN：`acme.corp.sdkdemo.HelloAPI`。非 Vault Owner 调用前，在自己的 Permission Set 上追加（把名字换成实际组件名）：

```mdl
ALTER Permissionset your_permission_set__c (
  webapigroup.integration__c.actions('execute')
);
```

调用：

```bash
curl -X POST "$VAULT/api/v22.3/custom/hello_world" \
  -H "Authorization: $SESSION" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada"}'
```

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
