# 构建 wasm

## 命令

```bash
vivarcus-sdk build <module-dir> [-o action.wasm] [--skip-compile]
vivarcus-sdk describe action.wasm
```

| 标志 | 说明 |
|------|------|
| `-o` | 输出 wasm 路径（默认 `<module-dir>/action.wasm`） |
| `--skip-compile` | 仅扫描 + codegen reactor，不调用 tinygo |
| `--compiler tinygo` | 默认；Phase 1 仅支持 tinygo |

`describe` 打印 `__sdk_describe`：每条 Action 的 `component_name`（Recordaction 名）、label、object、object_action，以及每条 Trigger 的 `component_name`（Recordtrigger 名）、events、event_segment、order。名称由本地 `go.mod` 模块路径 + Go 类型名派生（部署时与 `vaultpackage.xml` `<gosdk><module>` 相同）；需要固定名字时在 `Meta.Name` 写出。

## 构建流程

1. **扫描** — 拒绝 `os`、`net` 等禁止 import
2. **发现入口** — 递归扫描模块（含 `actions/`、`triggers/`、`entries/` 等子目录）中所有实现 `Meta`/`IsExecutable`/`Execute` 的类型（可多个）
3. **Codegen** — 生成 `zz_generated_reactor.go`（勿手改）
4. **编译** — `tinygo build -target=wasi` 产出 wasm
5. **Describe** — 本地校验 `__sdk_describe` 列出的 Action（FQN / label / usages）

## 产出物

```
my-action/
├── main.go
├── zz_generated_reactor.go   # 自动生成
└── action.wasm               # 本地验证；不要放入 VPK
```

一个 module 可以包含多个 Action 类型；它们编译进 **同一份** wasm。导入时平台扫描 wasm，以 describe 为真源，不再需要 `sdk_manifest.json`。

`Meta.ObjectAction` 非空时，部署会同时创建 active 的 Objectaction（记录页按钮）。

## 限制

- wasm 体积 ≤ **2 MB**
- 入口类型须有唯一 Go 类型名（FQN 由本地 `go.mod` / `<gosdk><module>` + 类型名派生；需要固定名字时写 `Meta.Name`）
- `shared/` 等 helper 的每个 `.go` 投影为一条 `Sdkcode`（该文件源码），**不能**放 `Meta()` 调度入口；一个文件最多一个入口类型

## 工程化布局

仓库示例 [multi-component](../examples/multi-component)（`cd examples/multi-component && vivarcus-sdk build .`）展示 ADR-21 推荐形态。`vivarcus-sdk build` 会递归扫描子目录中的入口类型，在模块根生成 `zz_generated_reactor.go`（`package main`）并 import 各入口包。

**推荐（按职责分子目录）**：

```
my-vault-sdk/
├── go.mod              # 本地 DX；module path → FQN 前缀（打 VPK 时写入 xml）
├── main.go             # 薄入口 / codegen 落点（package main）
├── actions/            # package actions — Record Action 入口
│   ├── foo.go
│   └── bar.go
├── triggers/           # package triggers — Record Trigger 入口
│   └── baz.go
└── shared/             # 跨 Action/Trigger 复用（非调度入口）
```

也支持单一入口包 `entries/`（`package entries`），或继续把入口平铺在模块根（`package main`）。

约定：

- 子目录入口必须是**命名 package**（不能 `package main`；Go 无法 import `package main`）
- 同一 module 内 Action / Trigger **类型名不可重复**（即使在不同子包）
- 把入口从根目录挪到 `actions/` **不会改 FQN**（仍是 `<module>.<Type>`）
- `shared/` 出现 `Meta()` 入口会构建失败

本地 `go.mod` 可临时 `replace github.com/vivarcus/vivarcus-sdk => ../..`。VPK `gosdk/` 只上传 `.go`（**不要** `.wasm`、`go.mod`）；`module` 行由 `package-vpk.sh` 写入 `vaultpackage.xml`。

下一步：[04-package-vpk](04-package-vpk.md)
