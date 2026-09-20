# multi-component

**推荐入口**：一个 Go module 里多种 Record Action + Record Trigger，打进 **一个** combo VPK，端到端部署到 Vault。

对标 ADR-21 客户源码树：`go.mod`、共享包（UDC 模式）、分文件入口。

## 目录

```
multi-component/
├── go.mod                      # 客户 module path（FQN 前缀来源）
├── shared/                     # 跨入口复用（非调度入口，ADR-22）
│   ├── title.go
│   └── name.go
├── action_set_title.go         # package main — 用户按钮 Record Action
├── action_clear_title.go
├── action_stamp_on_enter.go    # 生命周期 entry_action（进状态时自动跑）
├── trigger_stamp_name.go       # BEFORE_INSERT Record Trigger
├── main.go
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
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
bash ../_shared/scripts/package-vpk.sh . action.vpk \
  --component 10:Object:sdk_demo__c:mdl/01-object.mdl
# vivarcus package import → validate → deploy（见 docs/05-deploy.md）
```

## Build only

```bash
cd examples/multi-component
GOTOOLCHAIN=go1.26.2 vivarcus-sdk build . -o action.wasm
```

## go.mod 说明

- 本地 clone 开发：`replace github.com/vivarcus/vivarcus-sdk => ../..`（见 [go.mod](go.mod)）。
- **打进 VPK 的 `go.mod` 不得含 `replace`**（见 ADR-21 validate）。

## 流程（端到端）

1. `vivarcus-sdk build` — 扫描 module、生成 glue、本地校验
2. 应用 lifecycle MDL（`02-lifecycle.mdl` 把 entry_action 绑到 `Recordaction.<FQN>`）
3. Combo VPK：对象 MDL + `gosdk/` Go 源码
4. `import` → `validate` → `deploy`（投影 active Recordaction / Objectaction / Recordtrigger）
5. 验证：
   - CREATE 时 `StampName` Trigger 给 `name__v` 追加 `-trig`
   - 按钮 `set_title__c` → `title__c=from-sdk`，再 `clear_title__c` 清空
   - lifecycle `submit__c` 进入 `in_review__c` → `StampOnEnter` 把 `title__c` 写成 `stamped-on-enter`

## 分场景精读

| 主题 | 本示例 | 单一入口精读 |
|------|--------|--------------|
| 用户按钮改字段 | `action_set_title.go` | [02-update-field](../02-update-field) |
| 生命周期 entry_action | `action_stamp_on_enter.go` + `mdl/02-lifecycle.mdl` | [04-lifecycle-entry](../04-lifecycle-entry) |
| Record Trigger | `trigger_stamp_name.go` | [05-stamp-trigger](../05-stamp-trigger) |

entry / event / workflow 更多系统路径见 [04-lifecycle-entry](../04-lifecycle-entry) 与 [docs/06-lifecycle.md](../../docs/06-lifecycle.md)。

Agent 清单：[AGENTS.md](../../AGENTS.md)。
