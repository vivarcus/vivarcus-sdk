# 打包 VPK

客户 Record Action 通过 **Inbound VPK** 部署。结构：**可选 `components/` MDL 步** + 末尾 **`gosdk/`**（一份或多份 `.wasm`）。

## 目录结构

纯 gosdk：

```text
my-action.vpk (zip)
├── vaultpackage.xml
└── gosdk/
    └── action.wasm
```

**Combo VPK**（对象 MDL 与 Action 同包）：

```text
demo-action.vpk (zip)
├── vaultpackage.xml
├── components/
│   └── 00010/
│       ├── Object.sdk_demo__c.mdl
│       └── Object.sdk_demo__c.md5
└── gosdk/
    └── action.wasm
```

平台 deploy 顺序：按 step 编号执行所有 `components/<step>/` MDL，**最后**执行 `gosdk/`：

- 扫描 `gosdk/*.wasm`
- 实例化并调用 `__sdk_describe`，得到 Action 列表
- 每个 Action 创建 active `Recordaction`（及 `Meta.ObjectAction` 对应的 Objectaction）
- 同一 wasm 的多条 Recordaction 共享同一个 blob

## vaultpackage.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<VaultPackage xmlns="https://www.veevavault.com/schema/vaultpackage">
  <name>My Demo Action</name>
  <summary>Deploy demo Record Action</summary>
  <packagetype>migration__v</packagetype>
</VaultPackage>
```

## 打包命令

```bash
mkdir -p dist/components/00010 dist/gosdk
cp rendered/01-object.mdl dist/components/00010/Object.demo_request__c.mdl
md5sum dist/components/00010/Object.demo_request__c.mdl | awk '{print $1" Object.demo_request__c"}' \
  > dist/components/00010/Object.demo_request__c.md5
cp action.wasm dist/gosdk/
cd dist && zip -r ../my-action.vpk vaultpackage.xml components/ gosdk/
```

或示例脚本：

```bash
./_shared/scripts/package-vpk.sh \
  --component 10:Object:sdk_demo__c:./rendered/01-object.mdl \
  ./action.wasm ./demo-action.vpk
```

失败时 `deployment_status` 为 `not_verified__v`，issues 含 `gosdk_invalid`。
