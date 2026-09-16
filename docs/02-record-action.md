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

返回 `error` 时事务回滚。

## 可选方法

实现以下方法时，`ov-sdk build` 会自动导出对应 wasm 钩子：

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

`ov-sdk build` 在 manifest 中生成默认值，部署前请改为正式 FQN。

下一步：[03-build](03-build.md)
