| 文件 | 说明 |
|------|------|
| [01-object.mdl](01-object.mdl) | `RECREATE` `sdk_demo__c`（含 `stamped__c`，供 `Process` 写入）— `vivarcus component apply-mdl --confirm -f mdl/01-object.mdl`，**不要**打进 VPK |

与 [jobs/stamp_records.go](../jobs/stamp_records.go) 中 `platform.Update(..., "stamped__c")` 对齐。与 [02-update-field](../../02-update-field/mdl/01-object.mdl) 等同对象定义，可重复 `apply-mdl`。
