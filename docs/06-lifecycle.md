# 生命周期 Action

在生命周期状态页显示自定义按钮时，使用 `UsageLifecycleUserAction`。

## Go 代码

```go
func (T) Meta() action.Meta {
    return action.Meta{
        Label:  "Submit for Review",
        Object: "demo_request__c",
        Usages: []action.Usage{action.UsageLifecycleUserAction},
    }
}
```

## MDL：创建生命周期

使用 [templates/mdl/02-lifecycle.mdl](../templates/mdl/02-lifecycle.mdl)：

```mdl
CREATE Lifecycle demo_request_lc__c (
  label('Demo Request LC'),
  object('demo_request__c'),
  State draft__c (label('Draft'), active(true)),
  State in_review__c (label('In Review'), active(true))
);
```

## 绑定到状态

在 Vault 配置中将 Action 绑定到目标生命周期状态的 **User Actions**（与 Veeva 生命周期配置类似）。具体 UI 路径因环境而异；MDL 侧确保 `Objectaction` 已部署且 `active(true)`。

## 与 UserAction 的区别

| Usage | 出现位置 |
|-------|----------|
| `UserAction` | 记录详情 **All Actions** |
| `LifecycleUserAction` | 生命周期状态页操作区 |

同一 Action 可同时声明多个 `Usages`（`UserBulkAction` 除外，须单独使用）。

## 部署顺序

1. CREATE Object
2. CREATE Lifecycle + States
3. deploy gosdk VPK
4. ALTER Recordaction / Objectaction `active(true)`
5. 在生命周期配置中绑定 Action 到状态

Record Action 部署与激活步骤见 [05-deploy](05-deploy.md)。
