# 打包 VPK

客户 Record Action 通过 **Inbound VPK** 部署。结构：**可选 `components/` MDL 步** + 末尾 **`gosdk/`**（Go 源码树）。

## 目录结构

纯 gosdk：

```text
my-action.vpk (zip)
├── vaultpackage.xml
└── gosdk/
    ├── go.mod
    └── main.go
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
    ├── go.mod
    ├── main.go
    └── shared/...
```

平台 deploy 顺序：按 step 编号执行所有 `components/<step>/` MDL，**最后**执行 `gosdk/`：

- 校验 `gosdk/` 下 `.go` / `go.mod`（**拒绝** `.wasm`）
- 合并到 Vault 源码树后由平台 `tinygo` 整树编译
- `__sdk_describe` 投影 Action / Trigger 列表
- 每个 Action 创建 active `Recordaction`（及 `Meta.ObjectAction` 对应的 Objectaction）
- 每个 Trigger 创建 active `Recordtrigger`

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
cp go.mod main.go dist/gosdk/
cd dist && zip -r ../my-action.vpk vaultpackage.xml components/ gosdk/
```

`package-vpk.sh` 会去掉 `replace`、把 SDK `require` 写成 **Go module tag**（由平台 tag 映射，如 `v26R3.3-13316` → `v1.26.3-3.13316`），并尽量写入 `go.sum`。可用平台镜像 tag（如 `v26R3.3-13316`）或 Go module tag 通过 `VIVARCUS_SDK_MODULE_REF` 覆盖。**不要**使用裸 `v0.0.0`。

平台 validate/deploy 时会把 `require` 里的 assembly 后缀与**当前 Vault 镜像**的 assembly 比对，并把 `go` 行与镜像内 `sdk/go.mod` 的 Go 版本（当前 `1.26.2`）比对；不一致则 `gosdk_invalid`（实际编译仍走镜像内 `/app/sdk/`）。

或示例脚本：

```bash
./_shared/scripts/package-vpk.sh \
  --component 10:Object:sdk_demo__c:./rendered/01-object.mdl \
  ./my-action ./demo-action.vpk
```

失败时 `deployment_status` 为 `not_verified__v`，issues 含 `gosdk_invalid`。
