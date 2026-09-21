# 模块布局与平台编译

客户**不**在本机编译 wasm。编号示例 `sdk put` 一个 `.go`；[multi-component](../examples/multi-component) 才上传 VPK 里的 `gosdk/`。Vault 扫描入口、生成 reactor、用镜像内 TinyGo 出 wasm。本地反馈用 `go test`（见 [examples/02-update-field](../examples/02-update-field)）。

本仓库是 guest API，不含编译器。

## 平台做什么

1. **扫描** — 拒绝 `os`、`net` 等禁止 import
2. **发现入口** — 递归扫描模块（含 `actions/`、`triggers/`、`entries/` 等子目录）中所有实现对应接口的类型（可多个）
3. **Codegen** — 生成 reactor shim（开发者不可见、不进 VPK）
4. **编译** — `tinygo build -target=wasi`
5. **Describe** — `__sdk_describe` 列出 Action / Trigger / Job / Web API / Workflow Action（FQN / label / usages）

名称由 `go.mod` 的 `module`（`sdk put --module` 或 VPK `vaultpackage.xml` `<gosdk><module>`）+ Go 类型名派生；需要固定名字时在 `Meta.Name` 写出。

一个 module 的全部入口编译进 **同一份** wasm。`Meta.ObjectAction` 非空时，部署会同时创建 active 的 Objectaction（记录页按钮）。

## 限制

- wasm 体积 ≤ **2 MB**
- 入口类型须有唯一 Go 类型名
- `shared/` 等 helper 的每个 `.go` 投影为一条 `Sdkcode`，**不能**放 `Meta()` 调度入口；一个文件最多一个入口类型

## 工程化布局

仓库示例 [multi-component](../examples/multi-component) 展示推荐形态。

```
my-vault-sdk/
├── go.mod              # 本地 DX；module path → FQN 前缀（打 VPK 时写入 xml）
├── main.go             # 薄入口（package main）
├── actions/            # package actions — Record Action 入口
│   ├── foo.go
│   └── bar.go
├── triggers/           # package triggers — Record Trigger 入口
│   └── baz.go
└── shared/             # 跨 Action/Trigger 复用（非调度入口）
```

也支持单一入口包 `entries/`（`package entries`），或把入口平铺在模块根（`package main`）。

约定：

- 子目录入口必须是**命名 package**（不能 `package main`；Go 无法 import `package main`）
- 同一 module 内 Action / Trigger **类型名不可重复**（即使在不同子包）
- 把入口从根目录挪到 `actions/` **不会改 FQN**（仍是 `<module>.<Type>`）
- `shared/` 出现 `Meta()` 入口会构建失败

本地 `go.mod` 可临时 `replace github.com/vivarcus/vivarcus-sdk => ../..`。`sdk put` 与 VPK 都只传 `.go`（**不要** `.wasm`、`go.mod`）。VPK 的 `module` 行由 `package-vpk.sh` 写入 `vaultpackage.xml`；`sdk put` 读邻近 `go.mod` 或 `--module`。

编号示例下一步：[05-deploy](05-deploy.md) 的「单文件」。多文件 VPK：[04-package-vpk](04-package-vpk.md)
