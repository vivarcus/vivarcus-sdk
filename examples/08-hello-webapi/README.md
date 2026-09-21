# 08-hello-webapi

**Custom Web API** 示例：`POST /api/{version}/custom/hello_world`，JSON 体 `{ "name": "Ada" }` 返回 `message`。

对标 Veeva Java SDK `com.veeva.vault.sdk.api.webapi.WebApi`（`@WebApiInfo` + `execute`）。调用前须有管理员创建的 `Webapigroup`（本例 `integration__c`），并在 Permission Set 上授予 `webapigroup.integration__c.actions('execute')`。`vault_actions.api('all_api')` **不**打开该 Group。

## Build

```bash
cd examples/08-hello-webapi
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
vivarcus-sdk describe action.wasm
```

Deploy 见 [multi-component](../multi-component)。先 `CREATE Webapigroup integration__c`，再 gosdk 部署，然后：

```bash
curl -X POST "$VAULT/api/v22.3/custom/hello_world" \
  -H "Authorization: $SESSION" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada"}'
```
