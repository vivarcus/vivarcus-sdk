# 08-hello-webapi

**Custom Web API** 示例：`POST /api/{version}/custom/hello_world`，JSON 体 `{ "name": "Ada" }` 返回 `message`。

对标 Veeva Java SDK `com.veeva.vault.sdk.api.webapi.WebApi`（`@WebApiInfo` + `execute`）。调用前须有管理员创建的 `Webapigroup`（本例 `integration__c`），并在 Permission Set 上授予 `webapigroup.integration__c.actions('execute')`。`vault_actions.api('all_api')` **不**打开该 Group。

**不要打 VPK。** 多文件 combo 见 [multi-component](../multi-component)。

## Quick start

```bash
cd examples/08-hello-webapi
# 若 vault 里还没有该 Group：CREATE Webapigroup integration__c（MDL / Admin）
vivarcus sdk put -f main.go --json
```

FQN：`acme.corp.hellowebapi.HelloAPI`。先有 `Webapigroup`，再 `put`，然后：

```bash
curl -X POST "$VAULT/api/v22.3/custom/hello_world" \
  -H "Authorization: $SESSION" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada"}'
```
