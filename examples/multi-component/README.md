# multi-component

**本仓库唯一用 Inbound VPK 部署的示例**：一个 Go module 里多种 Record Action + Record Trigger，打进 **一个** combo VPK。

编号示例（01–09）全部改用 `vivarcus sdk put`，不要拿本目录的 `package-vpk.sh` 去部署它们。

对标 ADR-21 客户源码树：`go.mod`、共享包（UDC 模式）、按子目录组织入口。

## 目录

```
multi-component/
├── go.mod                      # 本地 DX（FQN 前缀来源；不进 VPK）
├── actions/                    # package actions — Record Action 调度入口
│   ├── set_title.go
│   ├── clear_title.go
│   └── stamp_on_enter.go       # 生命周期 entry_action（进状态时自动跑）
├── triggers/                   # package triggers — Record Trigger 调度入口
│   └── stamp_name.go
├── shared/                     # 跨入口复用（非调度入口，ADR-22）
│   ├── title.go
│   └── name.go
└── mdl/
    ├── 01-object.mdl
    ├── 02-lifecycle.mdl        # entry_action rule → Recordaction.<FQN>
    └── 03-bind-object-lifecycle.mdl
```

Deploy 时平台从 `gosdk/` Go 源码编译并扫描组件清单，投影：

| 入口 | 类型 | 触发方式 |
|------|------|----------|
| `SetTitle` / `ClearTitle` | Recordaction + Objectaction | 记录页按钮（`UsageUserAction`） |
| `StampOnEnter` | Recordaction | 生命周期 **Entry Action**（`UsageLifecycleEntryAction`） |
| `StampName` | Recordtrigger | `BEFORE_INSERT`（创建记录时自动跑） |

## Quick start

```bash
cd examples/multi-component
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
vivarcus package import ./action.vpk
# 记下 package_id
vivarcus package validate <package_id> --json
vivarcus package deploy <package_id> --confirm --json
```

不必本机编 wasm。模块布局见 [docs/03-build.md](../../docs/03-build.md)；三步细节见 [docs/05-deploy.md](../../docs/05-deploy.md)。

## go.mod 说明

- 本地 clone 开发：`replace github.com/vivarcus/vivarcus-sdk => ../..`（见 [go.mod](go.mod)）。
- **不要把 `go.mod` 打进 VPK**，也不要在业务 `.go` 里写 `func main`，不要再放只做空导入的根 `main.go`。`package-vpk.sh` 只打包业务 `.go`，并把 `module` 行写入 `vaultpackage.xml`。平台 reactor 生成 `func main` 并 import 入口包。

## 流程（端到端）

1. 应用 lifecycle MDL（`02-lifecycle.mdl` 把 entry_action 绑到 `Recordaction.<FQN>`；须在 gosdk deploy **之后**，或与对象 MDL 分步 apply）
2. Combo VPK：对象 MDL + `gosdk/` 全部 `.go`（`package-vpk.sh`）
3. `import` → `validate` → `deploy`（Vault 编译；投影 active Recordaction / Objectaction / Recordtrigger）
4. 验证：
   - CREATE 时 `StampName` Trigger 给 `name__v` 追加 `-trig`
   - 按钮 `set_title__c` → `title__c=from-sdk`，再 `clear_title__c` 清空
   - lifecycle `submit__c` 进入 `in_review__c` → `StampOnEnter` 把 `title__c` 写成 `stamped-on-enter`
   - Agent 用 CLI 核对投影（FQN = `acme.corp.sdkdemo.<Type>`）：

```bash
vivarcus sdk get acme.corp.sdkdemo.SetTitle -o /tmp/SetTitle.go --json
# 可启停：SetTitle / ClearTitle / StampOnEnter / StampName
# 不可启停：shared/title.go、shared/name.go 投影的 Sdkcode
vivarcus sdk disable acme.corp.sdkdemo.SetTitle --json
vivarcus sdk enable acme.corp.sdkdemo.SetTitle --json
```

本树已用 VPK 一次部署齐全。不要对编号示例打 VPK；单文件迭代见 [01-hello-action](../01-hello-action)。

## 分场景精读（编号示例用 sdk put）

| 主题 | 本示例（VPK） | 单一入口（`sdk put`） |
|------|--------|--------------|
| 用户按钮改字段 | `actions/set_title.go` | [02-update-field](../02-update-field) |
| 生命周期 entry_action | `actions/stamp_on_enter.go` + `mdl/02-lifecycle.mdl` | [04-lifecycle-entry](../04-lifecycle-entry) |
| Record Trigger | `triggers/stamp_name.go` | [05-stamp-trigger](../05-stamp-trigger) |

entry / event / workflow 更多系统路径见 [04-lifecycle-entry](../04-lifecycle-entry) 与 [docs/06-lifecycle.md](../../docs/06-lifecycle.md)。

Agent 清单：[AGENTS.md](../../AGENTS.md)。
