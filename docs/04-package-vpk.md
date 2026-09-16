# 打包 VPK

客户 Record Action 通过 **Inbound VPK** 的 `gosdk/` 步骤部署。

## 目录结构

```text
my-action.vpk (zip)
├── vaultpackage.xml
└── gosdk/
    ├── sdk_manifest.json
    └── action.wasm
```

## vaultpackage.xml

复制 [templates/vpk/vaultpackage.xml.tpl](../templates/vpk/vaultpackage.xml.tpl)，替换占位符：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<VaultPackage xmlns="https://www.veevavault.com/schema/vaultpackage">
  <name>My Demo Action</name>
  <summary>Deploy demo Record Action</summary>
  <packagetype>migration__v</packagetype>
</VaultPackage>
```

## sdk_manifest.json

由 `ov-sdk build` 生成，或从 [templates/vpk/sdk_manifest.json.tpl](../templates/vpk/sdk_manifest.json.tpl) 填写：

```json
{
  "api_version": "1",
  "component_name": "com.acme.actions.Approve",
  "label": "Approve",
  "object": "demo_request__c",
  "object_action": "demo_request__c.approve__c",
  "usages": ["UserAction"],
  "wasm_file": "action.wasm",
  "sha256": "<与 action.wasm 一致的 SHA-256 hex>"
}
```

**关键对齐规则：**

- `sha256` 必须与 `gosdk/action.wasm` 字节完全一致
- `object` 须为 Vault 中已存在的对象 api_name
- `object_action` 通常为 `<object>.<verb>__c` 格式

## 打包命令

手动：

```bash
mkdir -p dist/gosdk
cp action.wasm action.sdk_manifest.json dist/gosdk/
cp sdk_manifest.json dist/gosdk/   # 或重命名 action.sdk_manifest.json
cp vaultpackage.xml dist/
cd dist && zip -r ../my-action.vpk vaultpackage.xml gosdk/
```

或使用脚本：[examples/99-full-stack/scripts/package-vpk.sh](../examples/99-full-stack/scripts/package-vpk.sh)

## validate 阶段校验

平台在 `ov package validate` 时检查：

- `api_version` 是否支持
- `sha256` 与 wasm 是否匹配
- wasm import 是否在白名单内（禁止私自 import 网络/文件系统等）

失败时 `deployment_status` 为 `not_verified__v`，issues 含 `gosdk_invalid`。

下一步：[05-deploy](05-deploy.md)
