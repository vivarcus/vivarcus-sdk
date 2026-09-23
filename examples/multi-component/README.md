# multi-component

**一个 module、多个入口** 的示例，和 01–09 用同一个 Go module（`github.com/acme.corp.sdkdemo`）。入口的路径和类型名与编号示例错开，可以增量部署到同一棵源码树上。

按 ADR-21 组织客户源码树：`go.mod`、共享包、按子目录组织入口。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/multi-component`（monorepo：`cd sdk/examples/multi-component`）。执行 `./deploy.sh`。

## 目录

```
multi-component/
├── go.mod                      # 本地 DX（FQN 前缀来源；不进 VPK）
├── actions/                    # package actions — Record Action 调度入口
│   ├── set_title_shared.go     # 按钮 set_title_shared__c（不与 02 的 SetTitle 撞名）
│   ├── clear_title.go
│   └── stamp_demo_on_enter.go  # 生命周期 entry_action（不与 04 的 StampOnEnter 撞名）
├── triggers/                   # package triggers — Record Trigger 调度入口
│   └── stamp_demo_name.go      # 不与 05 的 StampName 撞名；后缀已有则不再追加
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
| `SetTitleShared` / `ClearTitle` | Recordaction + Objectaction | 记录页按钮（`UsageUserAction`） |
| `StampDemoOnEnter` | Recordaction | 生命周期 **Entry Action**（`UsageLifecycleEntryAction`） |
| `StampDemoName` | Recordtrigger | `BEFORE_INSERT`（创建记录时自动跑） |

**一键**（推荐）：`./deploy.sh`（生命周期 MDL + 打 VPK + import/validate/deploy，与集成测试一致）

### VPK 命令（手动）

```bash
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:multi_demo__c:mdl/01-object.mdl
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

1. 应用 lifecycle MDL，再 apply `01-object.mdl` 创建 `multi_demo__c` 并绑上 `multi_demo_lc__c`（`02-lifecycle.mdl` 把 entry_action 绑到 `Recordaction.<FQN>`；须在 gosdk deploy **之后**，或与对象 MDL 分步 apply）
2. Combo VPK：对象 MDL + `gosdk/` 全部 `.go`（`package-vpk.sh`）
3. `import` → `validate` → `deploy`（Vault 编译；投影 active Recordaction / Objectaction / Recordtrigger）
4. 验证：
   - CREATE 时 `StampDemoName` Trigger 给 `name__v` 追加 `-trig`（已有后缀则跳过）
   - 按钮 `set_title_shared__c` → `title__c=from-sdk`，再 `clear_title__c` 清空
   - lifecycle `submit__c` 进入 `in_review__c` → `StampDemoOnEnter` 把 `title__c` 写成 `stamped-on-enter`
   - Agent 用 CLI 核对投影（FQN = `acme.corp.sdkdemo.<Type>`）：

```bash
vivarcus sdk get acme.corp.sdkdemo.SetTitleShared -o /tmp/SetTitleShared.go --json
# 可启停：SetTitleShared / ClearTitle / StampDemoOnEnter / StampDemoName
# 不可启停：shared/title.go、shared/name.go 投影的 Sdkcode
vivarcus sdk disable acme.corp.sdkdemo.SetTitleShared --json
vivarcus sdk enable acme.corp.sdkdemo.SetTitleShared --json
```

本树已用 VPK 一次部署齐全。不要对编号示例打 VPK；单文件迭代见 [01-hello-action](../01-hello-action)。

## 分场景精读（编号示例用 sdk put）

| 主题 | 本示例（VPK） | 单一入口（`sdk put`） |
|------|--------|--------------|
| 用户按钮改字段 | `actions/set_title_shared.go` | [02-update-field](../02-update-field) |
| 生命周期 entry_action | `actions/stamp_demo_on_enter.go` + `mdl/02-lifecycle.mdl` | [04-lifecycle-entry](../04-lifecycle-entry) |
| Record Trigger | `triggers/stamp_demo_name.go` | [05-stamp-trigger](../05-stamp-trigger) |

entry / event / workflow 更多系统路径见 [04-lifecycle-entry](../04-lifecycle-entry) 与 [docs/06-lifecycle.md](../../docs/06-lifecycle.md)。

Agent 清单：[AGENTS.md](../../AGENTS.md)。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
