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

`describe` 打印 `__sdk_describe`：每条 Action 的 `component_name`（Recordaction 名）、label、object、object_action。名称由 `go.mod` 模块路径 + Go 类型名派生；需要固定名字时在 `Meta.Name` 写出。

## 构建流程

1. **扫描** — 拒绝 `os`、`net` 等禁止 import
2. **发现入口** — 模块内所有实现 `Meta`/`IsExecutable`/`Execute` 的类型（可多个）
3. **Codegen** — 生成 `zz_generated_reactor.go`（勿手改）
4. **编译** — `tinygo build -target=wasi` 产出 wasm
5. **Describe** — 本地校验 `__sdk_describe` 列出的 Action（FQN / label / usages）

## 产出物

```
my-action/
├── main.go
├── zz_generated_reactor.go   # 自动生成
└── action.wasm               # 放入 VPK gosdk/
```

一个 module 可以包含多个 Action 类型；它们编译进 **同一份** wasm。导入时平台扫描 wasm，以 describe 为真源，不再需要 `sdk_manifest.json`。

`Meta.ObjectAction` 非空时，部署会同时创建 active 的 Objectaction（记录页按钮）。

## 限制

- wasm 体积 ≤ **2 MB**
- 入口类型须在同一 Go package

下一步：[04-package-vpk](04-package-vpk.md)
