# MDL（05-stamp-trigger）

`01-object.mdl` 定义 `sdk_demo__c`（`name__v` 由平台标准字段注入），与 [triggers/stamp_name.go](../triggers/stamp_name.go) 里 `Meta.Object` 对齐。用 `vivarcus component apply-mdl --confirm -f mdl/01-object.mdl`，不要打 VPK。

Recordtrigger 由 `sdk put` 从 `__sdk_describe` 的 `triggers[]` 投影。端到端多文件见 [multi-component](../../multi-component)。
