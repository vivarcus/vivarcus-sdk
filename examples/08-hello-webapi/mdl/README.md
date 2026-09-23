| 文件 | 说明 |
|------|------|
| [01-webapigroup.mdl](01-webapigroup.mdl) | `RECREATE` `integration__c`（可重复 `apply-mdl`）— `vivarcus component apply-mdl --confirm -f mdl/01-webapigroup.mdl`，**不要**打进 VPK |

与 [webapis/hello_api.go](../webapis/hello_api.go) 的 `Meta.APIGroup` 保持一致。`Customwebapi` 由 `sdk put` 投影；Group 不存在时整次部署失败。
