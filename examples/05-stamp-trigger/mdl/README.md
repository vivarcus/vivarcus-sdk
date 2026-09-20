# MDL（05-stamp-trigger）

`01-object.mdl` 定义 `sdk_demo__c`（`name__v` 由平台标准字段注入），与 [main.go](../main.go) 里 `Meta.Object` 对齐。

Recordtrigger 组件由 gosdk VPK deploy 从 `__sdk_describe` 的 `triggers[]` 投影，不必手写进 `components/`。端到端见 [multi-component](../../multi-component)。
