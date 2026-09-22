# 03-confirm-dialog

`OnPreExecute` 确认框 + `OnPostExecute` 横幅（需 UI 手测横幅）。

## Quick start

```bash
cd examples/03-confirm-dialog
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/confirm_dialog.go --json
```

FQN：`acme.corp.confirmdialog.ConfirmDialog`。按钮 `confirm_update__c` 出现在 `sdk_demo__c` 记录页 **All Actions**。

模块布局：[docs/03-build.md](../../docs/03-build.md)。
