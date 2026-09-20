# MDL（multi-component）

| 文件 | 作用 |
|------|------|
| [01-object.mdl](01-object.mdl) | `sdk_demo__c`：`name__v`（Trigger 写入）与 `title__c`（Action 写入） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | `sdk_demo_lc__c`：`submit__c` 状态迁移 + `in_review__c` **Entry Action** 绑 `Recordaction.{{ACTION_FQN}}` |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | 把 lifecycle 挂到 `sdk_demo__c` |

Recordtrigger / Recordaction 由 gosdk VPK deploy 从组件清单投影，不必手写进本目录。`{{ACTION_FQN}}` 在跑 demo 时由 `vivarcus-sdk describe` 填入 `StampOnEnter` 的 FQN。

与 [Go 入口](../action_set_title.go) 里 `Meta.Object`（`sdk_demo__c`）对齐。
