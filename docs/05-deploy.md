# 部署与激活

## 前置：对象与字段

Action 绑定的对象须已存在。先执行 MDL：

```bash
# 编辑 templates/mdl/01-object.mdl 中的 {{OBJECT}} 等占位符
ov mdl run templates/mdl/01-object.mdl
```

若 Action 用于生命周期状态页，还需 [06-lifecycle](06-lifecycle.md)。

## 部署三步

```bash
# 1. 导入
ov package import ./my-action.vpk
# 返回 package_id

# 2. 校验
ov package validate <package_id> --json
# 期望 deployment_status 非 not_verified__v

# 3. 部署（非 TTY 须 --confirm）
ov package deploy <package_id> --confirm --json
# 期望 deployment_status 为 deployed__v
```

部署成功后：

- wasm 写入 blob store
- 创建 **inactive** 的 `Recordaction` 与 `Objectaction`
- `source_code` 格式为 `<blob_id>@<sha256hex>`

## 激活

客户组件默认 **inactive**，须管理员激活后按钮才出现：

```bash
# 编辑 templates/mdl/03-recordaction-active.mdl
ov mdl run templates/mdl/03-recordaction-active.mdl
```

或手动 MDL：

```mdl
ALTER Recordaction com.acme.actions.Approve SET active(true);
ALTER Objectaction demo_request__c.approve__c SET active(true);
```

## 验证

1. 打开绑定对象的记录详情
2. 展开 **All Actions**
3. 应看到 manifest 中的 `label`
4. 点击执行，确认业务效果（字段变更、横幅等）

## 停用

```mdl
ALTER Recordaction com.acme.actions.Approve SET active(false);
```

按钮随即从菜单消失；正在执行的实例不受影响。

## 权限

- import / validate / deploy：Vault Owner 或 `configuration.deployment`
- ALTER active：metadata 编辑权限

见 [troubleshooting](troubleshooting.md) 处理常见失败。
