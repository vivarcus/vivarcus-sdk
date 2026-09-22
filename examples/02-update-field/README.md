# 02-update-field

记录页 **All Actions** 按钮：通过 `SetValue` + `platform.Update` 写入 `title__c`。

## Quick start

```bash
cd examples/02-update-field
go test ./...
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/set_title.go --json
```

FQN：`acme.corp.updatefield.SetTitle`。对象 / 按钮名写在 [actions/set_title.go](actions/set_title.go) 的 `Meta()` 与 [`mdl/`](mdl/) 里。`go test` mock `platform.UpdateRecordFunc`，见 [actions/set_title_test.go](actions/set_title_test.go)。

## 手工跟做

| 步骤 | 命令 |
|------|------|
| 创建对象 `sdk_demo__c` | `vivarcus component apply-mdl --confirm -f mdl/01-object.mdl` |
| 部署 `actions/set_title.go` | `vivarcus sdk put -f actions/set_title.go --json` |
| 核对投影 | `vivarcus sdk get acme.corp.updatefield.SetTitle -o /tmp/SetTitle.go --json` |

模块布局：[docs/03-build.md](../../docs/03-build.md)。
