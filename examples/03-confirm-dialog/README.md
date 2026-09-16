# 03-confirm-dialog

Record Action with **OnPreExecute** confirm dialog and **OnPostExecute** result banner.

## Prerequisites

- Object `sdk_demo__c` with field `title__c` (see `templates/mdl/01-object.mdl`)

## Build

```bash
ov-sdk build ./examples/03-confirm-dialog -o action.wasm
```

## Deploy

See [docs/05-deploy.md](../../docs/05-deploy.md). Set manifest:

- `component_name`: e.g. `com.example.ConfirmDialog`
- `object`: `sdk_demo__c`
- `object_action`: `sdk_demo__c.confirm_update__c`

## Expected UI

1. Click **Confirm Update** in All Actions
2. Confirm dialog appears
3. After confirm, `title__c` becomes `confirmed-by-sdk`
4. Success banner: "Title updated successfully."
