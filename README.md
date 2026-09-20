# Vivarcus SDK

在 Vivarcus Vault 上开发 **Record Action**（记录页自定义按钮）与 **Record Trigger** 的 Go 开发套件。通过 Inbound VPK 部署，deploy 后在 Vault 中可见。

> 对标 Veeva Vault Java SDK 的 `RecordAction` / `RecordTrigger`；Phase 1 使用 **Go**，`javasdk/` 暂不支持。

## 5 分钟 Quickstart

### 1. 准备工具

| 工具 | 用途 |
|------|------|
| [Go](https://go.dev/dl/) 1.22+ | 编写 Action / Trigger |
| [TinyGo](https://tinygo.org/getting-started/) | 本地 build 校验（VPK 只上传 Go 源码） |
| [vivarcus-sdk](docs/03-build.md) | 扫描、codegen、`vivarcus-sdk build` |
| [vivarcus CLI](docs/01-prerequisites.md) | 部署 VPK 到 Vault |

```bash
git clone https://github.com/vivarcus/vivarcus-sdk.git
cd vivarcus-sdk
```

### 2. 推荐：端到端示例

从 **[examples/multi-component](examples/multi-component)** 开始——一个 module、多个 Action + Trigger、一个 combo VPK：

```bash
cd examples/multi-component
GOTOOLCHAIN=go1.22.12 vivarcus-sdk build . -o action.wasm
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
# vivarcus package import → validate → deploy（见 docs/05-deploy.md）
```

涵盖 `go.mod` + `shared/` 布局、对象 MDL、VPK 打包、`import` → `validate` → `deploy`。

### 3. 写 Action

从模板开始：

```bash
cp templates/action/main.go.tpl my-action/main.go
# 编辑 my-action/main.go，实现 Meta / IsExecutable / Execute
```

分场景精读示例：

| 示例 | 路径 | 说明 |
|------|------|------|
| Hello Action | [examples/01-hello-action](examples/01-hello-action) | 最小空操作 |
| 更新字段 | [examples/02-update-field](examples/02-update-field) | 记录页按钮：`SetValue` + `platform.Update` |
| 确认对话框 | [examples/03-confirm-dialog](examples/03-confirm-dialog) | `OnPreExecute` / `OnPostExecute` |
| **系统自动执行** | [examples/04-lifecycle-entry](examples/04-lifecycle-entry) | entry / event / workflow / cancel 分步跟做 |
| Record Trigger | [examples/05-stamp-trigger](examples/05-stamp-trigger) | 单一 Trigger 入口 |

### 4. 构建

```bash
vivarcus-sdk build ./my-action
vivarcus-sdk describe ./my-action/action.wasm   # Recordaction 名由 go.mod + 类型名派生
```

### 5. 准备 Vault 元数据（MDL）

部署前，目标对象须已存在。复制并修改模板：

```bash
# 见 templates/mdl/README.md
vivarcus mdl run templates/mdl/01-object.mdl
```

### 6. 打包并部署 VPK

见 [multi-component](examples/multi-component) 或 [docs/05-deploy.md](docs/05-deploy.md)。VPK `gosdk/` 只放 Go 源码，可含多个 Action/Trigger。

### 7. 验证

- 按钮：记录详情 **All Actions**
- Trigger：创建记录时自动执行（multi-component demo 会 API 验收）

## 文档

| 文档 | 内容 |
|------|------|
| [01-prerequisites](docs/01-prerequisites.md) | 环境、权限、CLI 配置 |
| [02-record-action](docs/02-record-action.md) | API：Meta、Execute、可选钩子 |
| [03-build](docs/03-build.md) | `vivarcus-sdk build`、manifest 字段 |
| [04-package-vpk](docs/04-package-vpk.md) | VPK 目录结构与 `vaultpackage.xml` |
| [05-deploy](docs/05-deploy.md) | import → validate → deploy |
| [06-lifecycle](docs/06-lifecycle.md) | entry / workflow / 按钮：Usages 与两种 rule 挂法 |
| [07-limits](docs/07-limits.md) | 标准库边界、配额、Phase 1 能力 |
| [troubleshooting](docs/troubleshooting.md) | 常见错误 |

## Go API 包

| 包 | 说明 |
|----|------|
| `action` | Record Action 接口与上下文 |
| `trigger` | Record Trigger 接口与上下文 |
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
