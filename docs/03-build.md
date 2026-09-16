# 构建 wasm

## 命令

```bash
vivarcus-sdk build <module-dir> [-o action.wasm] [--skip-compile]
```

| 标志 | 说明 |
|------|------|
| `-o` | 输出 wasm 路径（默认 `<module-dir>/action.wasm`） |
| `--skip-compile` | 仅扫描 + codegen reactor，不调用 tinygo |
| `--compiler tinygo` | 默认；Phase 1 仅支持 tinygo |

## 构建流程

1. **扫描** — 拒绝 `os`、`net` 等禁止 import
2. **发现入口** — 模块内唯一实现 `Meta`/`IsExecutable`/`Execute` 的类型
3. **Codegen** — 生成 `zz_generated_reactor.go`（勿手改）
4. **编译** — `tinygo build -target=wasi` 产出 wasm
5. **Manifest** — 写入 `action.sdk_manifest.json`（含 `sha256`、`exports`）

## 产出物

```
my-action/
├── main.go
├── zz_generated_reactor.go   # 自动生成
├── action.wasm               # 部署用
└── action.sdk_manifest.json  # 放入 VPK gosdk/
```

## manifest 字段

| 字段 | 说明 |
|------|------|
| `api_version` | 固定 `"1"` |
| `component_name` | Recordaction FQN |
| `label` | 按钮名称 |
| `object` | 对象 api_name |
| `object_action` | Objectaction api_name（通常 `<object>.<action>__c`） |
| `usages` | 如 `["UserAction"]` |
| `wasm_file` | wasm 文件名 |
| `sha256` | wasm 文件 SHA-256 十六进制 |
| `exports` | wasm 导出列表（校验用） |

部署前请核对 `component_name`、`object`、`object_action` 与 Vault MDL 一致。

## 本地验证（无需 Vault）

```bash
vivarcus-sdk build ./examples/01-hello-action --skip-compile   # 验证扫描/codegen
vivarcus-sdk build ./examples/01-hello-action -o /tmp/action.wasm # 完整编译
sha256sum /tmp/action.wasm
```

## 限制

- wasm 体积 ≤ **2 MB**
- 单模块仅 **一个** 入口类型

下一步：[04-package-vpk](04-package-vpk.md)
