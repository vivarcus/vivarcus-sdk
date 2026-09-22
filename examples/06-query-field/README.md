# 06-query-field

记录页按钮：用 `platform.Query` 只读 VQL 统计同对象记录数，再把 `title__c` 写成 `query:<N>`。

## Quick start

```bash
cd examples/06-query-field
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/query_and_stamp.go --json
```

FQN：`acme.corp.queryfield.QueryAndStamp`。对象 / 按钮名见 [actions/query_and_stamp.go](actions/query_and_stamp.go) 与 [`mdl/`](mdl/)。验收期望 `title__c=query:1`（单条演示记录）。
