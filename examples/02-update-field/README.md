# 02-update-field

Record Action that writes `title__c` via `SetValue` and `platform.Update`.

## Prerequisites

- Object with field `title__c` (default example uses `sdk_demo__c`; full-stack uses `demo_request__c`)

## Build

```bash
vivarcus-sdk build ./examples/02-update-field -o action.wasm
```

## Manifest

After build, edit `action.sdk_manifest.json` so `object` / `object_action` / `component_name` match your Vault MDL.

## Deploy

See [99-full-stack](../99-full-stack) or [docs/05-deploy.md](../../docs/05-deploy.md).
