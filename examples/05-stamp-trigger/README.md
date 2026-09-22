# 05-stamp-trigger

**Record Trigger** 示例：`BEFORE_INSERT` 给 `name__v` 追加 `-trig`。

## Quick start

```bash
cd examples/05-stamp-trigger
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f triggers/stamp_name.go --json
```

FQN：`acme.corp.stamptrigger.StampName`。创建一条 `sdk_demo__c` 时 Trigger 自动跑。

模块布局：[docs/03-build.md](../../docs/03-build.md)。
