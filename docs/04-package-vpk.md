# 打包 VPK

**编号示例（01–09）不要走本文。** 它们用 `vivarcus sdk put -f main.go`，见 [05-deploy](05-deploy.md)。本文只给 [multi-component](../examples/multi-component) 这类多文件树。

客户多入口 Record Action 通过 **Inbound VPK** 部署。结构：**可选 `components/` MDL 步** + 末尾 **`gosdk/`**（仅 `.go`）。

本地 `go.mod` 只给 `go get` / `go test` 用，**不要**放进 VPK。客户 module path 写在 `vaultpackage.xml` 的 `<gosdk><module>`；平台编译时按当前镜像生成 `go.mod` 并注入 SDK。`package-vpk.sh` 会跳过 `*_test.go` 与 `zz_generated_reactor.go`。

## 目录结构

纯 gosdk：

```text
my-action.vpk (zip)
├── vaultpackage.xml    # 含 <gosdk><module>…
└── gosdk/
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
    ├── main.go
    ├── actions/...
    ├── triggers/...
    └── shared/...
```

平台 deploy 顺序：按 step 编号执行所有 `components/<step>/` MDL，**最后**执行 `gosdk/`：

- 校验 `gosdk/` 只含 `.go`（**拒绝** `.wasm`、`go.mod`、`go.sum`）
- 读取 `<gosdk><module>` 作为 FQN 前缀
- 合并到 Vault 源码树后由平台 `tinygo` 整树编译（SDK 来自镜像 `/app/sdk/`）
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
  <gosdk>
    <deployment_option>incremental</deployment_option>
    <module>github.com/acme.corp.sdkdemo</module>
  </gosdk>
</VaultPackage>
```

`<module>` 禁止 `com.example.` 前缀。`delete_all` 时可省略。

## 打包命令

推荐用脚本（从本地 `go.mod` 抄 `module` 进 xml，只打包 `.go`）：

```bash
./_shared/scripts/package-vpk.sh \
  --component 10:Object:sdk_demo__c:./rendered/01-object.mdl \
  ./my-action ./demo-action.vpk
```

手打时 **不要** `cp go.mod dist/gosdk/`。

失败时 `deployment_status` 为 `not_verified__v`，issues 含 `gosdk_invalid`。
