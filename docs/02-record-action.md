# Record Action API

每个 wasm 模块包含 **一个** Record Action 入口类型，实现以下接口。

## 必选方法

```go
type RecordAction interface {
    Meta() Meta
    IsExecutable(ctx RecordActionContext) bool
    Execute(ctx RecordActionContext) (ExecuteResult, error)
}
```

### Meta()

声明组件级静态元数据（对标 Java `@RecordActionInfo`）：

| 字段 | 必填 | 说明 |
|------|------|------|
| `Label` | 是 | 按钮显示名称 |
| `Object` | 否 | 绑定对象 api_name；空表示全对象 |
| `Usages` | 否 | 可用场景，默认 `Unspecified` |
| `Icon` | 否 | 按钮图标 |
| `RunAs` | 否 | `SystemUser`（默认）或 `RequestOwner` |

`Usages` 枚举：`UserAction`、`UserBulkAction`、`LifecycleUserAction`、`LifecycleEntryAction`、`EventAction`、`WorkflowStep`、`WorkflowCancel`、`Unspecified`。

### IsExecutable(ctx)

返回 `false` 时按钮隐藏或置灰。典型用法：检查记录状态、字段值、用户上下文。

### Execute(ctx)

核心业务逻辑。通过 `ctx.Records[0]` 访问目标记录；用 `rec.SetValue` 暂存字段变更，再调用 `platform.Update` 持久化。

`platform.Get` / `platform.Update` 的对象与记录 ID **必须**是当前这条上下文记录；传入其他 ID 会失败（`record_action_host_call_failed`）。这与 Veeva Java `RecordService` 允许任意 ID 不同，见 [07-limits](07-limits.md)。

返回 `error` 时事务回滚。

## 可选方法

实现以下方法时，平台编译会导出对应 wasm 钩子：

```go
// 执行前确认对话框（仅 UserAction）
func (T) OnPreExecute(ctx RecordActionContext) (PreExecuteResult, error)

// 执行后结果横幅（仅 UserAction）
func (T) OnPostExecute(ctx RecordActionContext) (PostExecuteResult, error)
```

`PreExecuteResult.ConfirmMessage` 非空时弹出确认框；用户取消则终止执行链。

## 上下文 RecordActionContext

| 字段 | 说明 |
|------|------|
| `Records` | 目标记录（UserAction 通常 1 条） |
| `UserInputRecord` | 用户输入表单（若有） |
| `Configuration` | 管理员配置的键值参数 |
| `VaultID` / `CurrentUserID` / `InitiatingUserID` | 运行时身份 |

## Record 便捷方法

```go
rec.GetString("title__c")
rec.GetInt("count__c")
rec.GetBool("active__c")
rec.SetValue("title__c", "new value")
```

## 组件命名

`component_name`（FQN）由包路径 + 类型名推导，例如：

- 包 `com.acme.actions`，类型 `Approve` → `com.acme.actions.Approve`

FQN 由本地 `go.mod` 的 `module` 路径 + Go 类型名派生（写入 `vaultpackage.xml` `<gosdk><module>`）。需要固定名字时在 `Meta.Name` 写出。

下一步：单个 `.go` 用 [05-deploy](05-deploy.md) 的 `sdk put`；多文件才 [04-package-vpk](04-package-vpk.md)（模块布局：[03-build](03-build.md)）
