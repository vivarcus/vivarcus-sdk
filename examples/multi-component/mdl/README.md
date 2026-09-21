# MDL（multi-component）

| 文件 | 作用 |
|------|------|
| [01-object.mdl](01-object.mdl) | `sdk_demo__c`：`title__c`（Action 写入）；`name__v` 由平台标准字段注入（Trigger 写入） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | `sdk_demo_lc__c`：`submit__c` 状态迁移 + `in_review__c` **Entry Action** 绑 `Recordaction.{{ACTION_FQN}}` |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | 可选：对已有 `sdk_demo__c` 补绑 lifecycle（VPK 内 `01-object.mdl` 已含 `available_lifecycles`） |

Recordtrigger / Recordaction 由 gosdk VPK deploy 从组件清单投影，不必手写进本目录。`{{ACTION_FQN}}` 在跑 demo 时由 `go.mod` module + `StampOnEnter` 类型名派生。

与 [Go 入口](../actions/set_title.go) 里 `Meta.Object`（`sdk_demo__c`）对齐。
