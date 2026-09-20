# 03-confirm-dialog

`OnPreExecute` 确认框 + `OnPostExecute` 横幅（需 UI 手测横幅）。

## Quick start

```bash
cd examples/03-confirm-dialog
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
```

## Build only

```bash
cd examples/03-confirm-dialog
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
```
