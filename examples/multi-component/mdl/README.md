# MDL（multi-component）

| 文件 | 作用 |
|------|------|
| [01-object.mdl](01-object.mdl) | `RECREATE` `multi_demo__c`（本示例专用，不共用 `sdk_demo__c`；含 `stamped__c`，`name__v` 由平台注入） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | `RECREATE` `multi_demo_lc__c`：`submit__c` + `in_review__c` **Entry Action**（`{{ACTION_FQN}}` 须 render） |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | 对 `multi_demo__c` 补绑 lifecycle（VPK 内 `01-object.mdl` 已含 `available_lifecycles`） |

Recordtrigger / Recordaction 由 gosdk VPK deploy 从组件清单投影，不必手写进本目录。`{{ACTION_FQN}}` 在跑 demo 时由 `go.mod` module + `StampDemoOnEnter` 类型名派生。

与 [Go 入口](../actions/set_title_shared.go) 里 `Meta.Object`（`multi_demo__c`）对齐。
