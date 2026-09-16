# 环境准备

## 必需软件

| 软件 | 版本 | 用途 |
|------|------|------|
| Go | 1.22+ | 编写 Action 源码 |
| TinyGo | 0.33+（建议最新稳定版） | 编译 wasm |
| Go（TinyGo 兼容） | TinyGo 当前支持 Go 1.19–1.26 | 本机 Go 过新时 `ov-sdk build` 会失败，见 [troubleshooting](troubleshooting.md) |
| ov-sdk | 与 Vault 版本对齐 | 扫描、codegen、`sdk_manifest.json` |
| ov CLI | 与目标 Vault 版本对齐 | MDL、VPK 部署 |

## 安装 TinyGo

```bash
# Linux（示例）
wget https://github.com/tinygo-org/tinygo/releases/download/v0.33.0/tinygo0.33.0.linux-amd64.tar.gz
tar -xzf tinygo0.33.0.linux-amd64.tar.gz
sudo mv tinygo /usr/local/
export PATH="/usr/local/tinygo/bin:$PATH"
tinygo version
```

## 安装 ov-sdk

从 [GitHub Releases](https://github.com/vivarcus/vivarcus-sdk/releases) 下载与 Vault 版本匹配的 `ov-sdk` 二进制，或使用安装脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/vivarcus/vivarcus-sdk/main/scripts/install-ov-sdk.sh | bash
```

将 `ov-sdk` 放入 `PATH`（例如 `~/.local/bin`）。

若暂无 Release，请联系 Vivarcus 支持获取对应版本的构建工具。

## 安装 ov CLI

从 [vivarcus/vivarcus-cli Releases](https://github.com/vivarcus/vivarcus-cli/releases) 下载与 Vault 版本匹配的 `ov` 二进制，或使用安装脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/vivarcus/vivarcus-cli/main/scripts/install-ov.sh | bash
```

将 `ov` 放入 `PATH`（例如 `~/.local/bin`）。

## 配置 ov CLI

```bash
ov auth login --endpoint https://<your-vault>.vivarcus.com
ov config set default_vault <vault_id>
```

## Vault 权限

| 操作 | 所需权限 |
|------|----------|
| 执行 MDL（CREATE Object 等） | Vault Owner 或等效 metadata 权限 |
| `ov package import/validate/deploy` | `configuration.deployment` |
| `ALTER Recordaction ... active(true)` | metadata 编辑权限 |

## 验证环境

```bash
# 1. tinygo 可用
tinygo version

# 2. ov-sdk 可用
ov-sdk build ./examples/01-hello-action --skip-compile

# 3. ov 已登录
ov whoami
```

下一步：[02-record-action](02-record-action.md)
