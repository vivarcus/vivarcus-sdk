# Vivarcus SDK

在 Vivarcus Vault 上开发 **Record Action**（记录页自定义按钮）的 Go 开发套件。代码编译为 WebAssembly，通过 Inbound VPK 部署，管理员激活后在记录详情 **All Actions** 中可见。

> 对标 Veeva Vault Java SDK 的 `RecordAction`；Phase 1 使用 **Go + wasm**，`javasdk/` 暂不支持。

## 5 分钟 Quickstart

### 1. 准备工具

| 工具 | 用途 |
|------|------|
| [Go](https://go.dev/dl/) 1.22+ | 编写 Action |
| [TinyGo](https://tinygo.org/getting-started/) | 编译 wasm |
| [ov-sdk](docs/03-build.md) | 扫描、codegen、打包 manifest |
| [ov CLI](docs/01-prerequisites.md) | 部署 VPK 到 Vault |

```bash
git clone https://github.com/vivarcus/vivarcus-sdk.git
cd vivarcus-sdk
```

### 2. 写 Action

从模板开始：

```bash
cp templates/action/main.go.tpl my-action/main.go
# 编辑 my-action/main.go，实现 Meta / IsExecutable / Execute
```

或参考示例：

| 示例 | 路径 | 说明 |
|------|------|------|
| Hello Action | [examples/01-hello-action](examples/01-hello-action) | 最小空操作 |
| 更新字段 | [examples/02-update-field](examples/02-update-field) | `SetValue` + `platform.Update` |
| 确认对话框 | [examples/03-confirm-dialog](examples/03-confirm-dialog) | `OnPreExecute` / `OnPostExecute` |
| 端到端 | [examples/99-full-stack](examples/99-full-stack) | MDL + 构建 + VPK + 部署脚本 |

### 3. 构建 wasm

```bash
ov-sdk build ./my-action -o action.wasm
# 产出 action.wasm 与 action.sdk_manifest.json
```

### 4. 准备 Vault 元数据（MDL）

部署 Action 前，目标对象须已存在。复制并修改模板：

```bash
# 见 templates/mdl/README.md
ov mdl run templates/mdl/01-object.mdl
```

### 5. 打包并部署 VPK

```bash
cd examples/99-full-stack
./scripts/package-vpk.sh ../02-update-field/action.wasm ../02-update-field/action.sdk_manifest.json
ov package import ./dist/my-action.vpk
ov package validate <package_id>
ov package deploy <package_id> --confirm
```

### 6. 激活

```bash
ov mdl run templates/mdl/03-recordaction-active.mdl
```

打开记录详情 → **All Actions** → 点击按钮验证。

## 文档

| 文档 | 内容 |
|------|------|
| [01-prerequisites](docs/01-prerequisites.md) | 环境、权限、CLI 配置 |
| [02-record-action](docs/02-record-action.md) | API：Meta、Execute、可选钩子 |
| [03-build](docs/03-build.md) | `ov-sdk build`、manifest 字段 |
| [04-package-vpk](docs/04-package-vpk.md) | VPK 目录结构与 `vaultpackage.xml` |
| [05-deploy](docs/05-deploy.md) | import → validate → deploy → 激活 |
| [06-lifecycle](docs/06-lifecycle.md) | 生命周期按钮与 MDL |
| [07-limits](docs/07-limits.md) | 标准库边界、配额、Phase 1 能力 |
| [troubleshooting](docs/troubleshooting.md) | 常见错误 |

## Go API 包

| 包 | 说明 |
|----|------|
| `action` | Record Action 接口与上下文 |
| `platform` | 宿主能力：`Get` / `Update` / `Log` |
| `wire` | ABI 编解码（一般由 codegen 使用） |

```go
import (
    "github.com/vivarcus/vivarcus-sdk/action"
    "github.com/vivarcus/vivarcus-sdk/platform"
)
```

## 编码 Agent

使用 Cursor / Copilot 等 Agent 时，请先读根目录 [AGENTS.md](AGENTS.md)。

## License

Apache License 2.0 — 见 [LICENSE](LICENSE)。
