# 03-confirm-dialog

`OnPreExecute` 确认框 + `OnPostExecute` 横幅（需 UI 手测横幅）。

**不要打 VPK。** 对象 MDL 单独 apply，Go 用 `sdk put`。多文件 combo 见 [multi-component](../multi-component)。

## Quick start

```bash
cd examples/03-confirm-dialog
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f main.go --json
```

FQN：`acme.corp.confirmdialog.ConfirmDialog`。按钮 `confirm_update__c` 出现在 `sdk_demo__c` 记录页 **All Actions**。

模块布局：[docs/03-build.md](../../docs/03-build.md)。
