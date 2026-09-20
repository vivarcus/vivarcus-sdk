# 部署

## 前置：对象与字段

Action 绑定的对象须已存在。两种方式：

1. **Combo VPK**（推荐，见 [04-package-vpk](04-package-vpk.md)）：对象 MDL 放在 `components/00010/`，与 `gosdk/` 同包 import/deploy。
2. **先 apply-mdl**：单独执行对象 MDL，再打纯 `gosdk/` VPK（生命周期、workflow 等复杂配置仍常用此方式）。

```bash
# 方式 2 示例
vivarcus component apply-mdl --confirm -f templates/mdl/01-object.mdl
```

若 Action 用于生命周期状态页，还需 [06-lifecycle](06-lifecycle.md)。

## 部署三步

```bash
# 1. 导入
vivarcus package import ./my-action.vpk
# 返回 package_id

# 2. 校验
vivarcus package validate <package_id> --json
# 期望 deployment_status 非 not_verified__v

# 3. 部署（非 TTY 须 --confirm）
vivarcus package deploy <package_id> --confirm --json
# 期望 deployment_status 为 deployed__v
```

Combo VPK deploy 时：先应用 `components/` 内 MDL（如创建对象），再安装 gosdk（创建 **active** Recordaction；`describe.triggers[]` 同时投影 **active** Recordtrigger）。

部署成功后：

- 平台编译 wasm 并写入 blob store
- 创建 **active** 的 `Recordaction` 与 `Objectaction`（`Meta.ObjectAction` 非空时）
- `describe` 中的 Trigger 创建 **active** 的 `Recordtrigger`（`BEFORE_*` / `AFTER_*` 随 DML 自动执行）
- `source_code` 格式为 `<blob_id>@<sha256hex>`

## 验证

1. 打开绑定对象的记录详情
2. 展开 **All Actions**
3. 应看到 manifest 中的 `label`
4. 点击执行，确认业务效果（字段变更、横幅等）

## 停用

```mdl
ALTER Recordaction com.acme.actions.Approve (active(false));
```

按钮随即从菜单消失；正在执行的实例不受影响。

## 权限

- import / validate / deploy：Vault Owner 或 `configuration.deployment`
- ALTER active：metadata 编辑权限

见 [troubleshooting](troubleshooting.md) 处理常见失败。
