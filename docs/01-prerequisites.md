# 环境准备

## 必需软件

写代码并部署到 Vault **只需要**：

| 软件 | 版本 | 用途 |
|------|------|------|
| Go | 与 `sdk/go.mod` 一致（当前 1.26.2） | 编写 Action 源码、本地 `go test` |
| vivarcus CLI | 与目标 Vault 版本对齐 | MDL、单文件 `sdk put`、VPK 部署 |

**CLI 与文档对照**（编号示例 01–09 用 `sdk put`，见 [05-deploy](05-deploy.md)）：

| vivarcus-cli Release | `vivarcus sdk` 子命令 |
|----------------------|------------------------|
| `v26R3.3-13318` 及更早 | 仅 `status` |
| 晚于 `v26R3.3-13318` 的首个 Release 起 | `status`、`logs`、`get`、`put`、`enable`、`disable` |

安装后执行 `vivarcus sdk --help` 确认；若缺少 `put`，请升级 CLI 或改用 [04-package-vpk](04-package-vpk.md) / multi-component VPK。

VPK 只上传 `.go`；wasm 由 Vault 镜像编译。本仓库（`github.com/vivarcus/vivarcus-sdk`）是 guest API，不含编译器。模块布局见 [03-build](03-build.md)。

### go.mod 与平台版本（双 tag）

平台发版使用 Veeva 风格 tag（`v26R3.3-13316`），Go modules 不能识别 `26R3.3` 这种 ADCV 字符串。公开仓在**同一 commit** 上会再打 Go 兼容 tag（如 `v1.26.3-3.13316`），供 `go.mod` / `go get` 使用。

| 用途 | tag 示例 |
|------|----------|
| Release / 镜像 | `v26R3.3-13316` |
| `go.mod` `require` / `go get` | `v1.26.3-3.13316` |

映射：`26R3.3` + assembly `13316` → `v1.26.3-3.13316`（Go module major 固定为 `v1.`，ADCV 编在 minor/patch/pre-release 里）。

本地 clone 开发：`require v1.26.3-3.13317`（Go module tag，对齐当前 train）+ `replace => ../..`。打 VPK 时 `package-vpk.sh` **不**打包 `go.mod`；只把 `module` 行写入 `vaultpackage.xml`。平台编译始终用 Vault 镜像内的 SDK。

## 安装 vivarcus CLI

从 [vivarcus/vivarcus-cli Releases](https://github.com/vivarcus/vivarcus-cli/releases) 下载与 Vault 版本匹配的 `vivarcus` 二进制，或使用安装脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/vivarcus/vivarcus-cli/main/scripts/install-vivarcus.sh | bash
```

将 `vivarcus` 放入 `PATH`（例如 `~/.local/bin`）。

## 配置 vivarcus CLI

人工首次登录（OAuth Device Flow）：

```bash
vivarcus auth login --endpoint https://<your-vault>.vivarcus.com
vivarcus config set default_vault <vault_id>
```

### Agent / CI：复用 token，避免登录限流

密码登录限流：**同一 IP + 用户名，1 分钟最多 10 次**（`POST /ui/auth/login`、Vault REST `POST /api/{version}/auth`）。编码 Agent 与自动化脚本**不要在每条 `vivarcus` 命令前重新 login**，也不要在循环里 `curl` 密码登录。

整段构建/部署任务复用同一 session token（最长 48 小时）：

```bash
export VIVARCUS_TOKEN=<session-token>
export VIVARCUS_ENDPOINT=https://<your-vault>.vivarcus.com
export VIVARCUS_VAULT=<vault_id>
vivarcus auth status --json
```

token 可从浏览器 `localStorage`、一次性 `auth login` 或 CI secret 注入。已 429 时按响应头 `Retry-After` 等待后再试。

## Vault 权限

| 操作 | 所需权限 |
|------|----------|
| 执行 MDL（CREATE Object 等） | Vault Owner 或等效 metadata 权限 |
| `vivarcus package import/validate/deploy` | `configuration.deployment` |
| `ALTER Recordaction ... (active(true))` | metadata 编辑权限 |

## 验证环境

```bash
# 1. Go 可用（版本见 sdk/go.mod）
go version

# 2. vivarcus 已登录
vivarcus auth status
```

下一步：[02-record-action](02-record-action.md)
