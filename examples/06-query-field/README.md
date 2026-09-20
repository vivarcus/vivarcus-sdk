# 06-query-field

记录页按钮：用 `platform.Query` 只读 VQL 统计同对象记录数，再把 `title__c` 写成 `query:<N>`。

## Quick start

```bash
cd examples/06-query-field
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
```

对象 / 按钮名见 [main.go](main.go) 与 [`mdl/`](mdl/)。验收期望 `title__c=query:1`（单条演示记录）。

## Build only

```bash
cd examples/06-query-field
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
vivarcus-sdk test . --context sdk_demo__c/<record_id> --field title__c=old
```
