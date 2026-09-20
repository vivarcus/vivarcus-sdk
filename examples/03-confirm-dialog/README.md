# 03-confirm-dialog

带 **OnPreExecute** 确认框与 **OnPostExecute** 成功横幅的 Record Action。

## Quick start

```bash
vivarcus auth login && vivarcus config set default_vault <uuid>
./scripts/dev-demo.sh
```

- 组件 MDL：[`mdl/`](mdl/)（对象随 VPK `components/` 部署）
- 脚本验证 Execute 改写字段；**确认框 / 横幅须 UI 手测**

## Build

```bash
cd examples/03-confirm-dialog
vivarcus-sdk build . -o action.wasm
```

按钮 api_name 在 `Meta.ObjectAction`（`confirm_update__c`）。

## Expected UI

1. All Actions → **Confirm Update**
2. 确认对话框
3. `title__c` → `confirmed-by-sdk`
4. 成功横幅："Title updated successfully."
