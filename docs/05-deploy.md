# 部署

编号示例（01–09）用 **单文件 `sdk put`**。[multi-component](../examples/multi-component) 才用 **Inbound VPK**。不要把两种流程套在同一个示例上。

## 单文件（01–09）

对象 / 生命周期等配置先 `apply-mdl`（没有对象的示例可跳过），再 PUT 一个 `.go`。平台增量 merge 后**整树重编译**，入口默认 **active**。对标 Veeva `PUT /api/{version}/code`。

```bash
cd examples/01-hello-action
vivarcus sdk put -f actions/noop_action.go --json
```

有对象时（02 / 03 / 05 / 06）：

```bash
cd examples/02-update-field
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f actions/set_title.go --json
```

在客户 module 目录内 `--path` / `--module` 从 `go.mod` 推断。文件不在树里时显式指定：

```bash
vivarcus sdk put -f /tmp/title.go --path shared/title.go --module github.com/acme.corp.hello --json
```

成功 `responseMessage` 为 `Modified file`。`DELETE /code/{name}` **不支持**。不要上传 `*_test.go` 或 `zz_generated_reactor.go`。

生命周期规则若引用 `Recordaction.<FQN>`：先 `sdk put` 投影组件，再 apply 那些 MDL（见 [04-lifecycle-entry](../examples/04-lifecycle-entry)）。

## Combo VPK（仅 multi-component）

对象 MDL 放进 `components/00010/`，与 `gosdk/` 同包 import/deploy。见 [04-package-vpk](04-package-vpk.md)。

也可以先单独 `apply-mdl`，再打纯 `gosdk/` VPK（生命周期、workflow 等复杂配置仍常用此方式）：

```bash
vivarcus component apply-mdl --confirm -f templates/mdl/01-object.mdl
```

若 Action 用于生命周期状态页，还需 [06-lifecycle](06-lifecycle.md)。

### 部署三步

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

- **平台编译** `gosdk/` 源码为 wasm 并写入 blob store
- 创建 **active** 的 `Recordaction` 与 `Objectaction`（`Meta.ObjectAction` 非空时）
- `describe` 中的 Trigger 创建 **active** 的 `Recordtrigger`（`BEFORE_*` / `AFTER_*` 随 DML 自动执行）
- `source_code` 格式为 `<blob_id>@<sha256hex>`

## 验证

1. 打开绑定对象的记录详情
2. 展开 **All Actions**
3. 应看到 `Meta.Label`
4. 点击执行，确认业务效果（字段变更、横幅等）

Agent 另可用 CLI 确认源码已投影（不要 `curl` `/code`，不要和 `--json` 把源码打到 stdout）：

```bash
# FQN = go.mod 的 module 去掉 github.com/ 前缀 + 入口类型名
# 例：examples/01-hello-action → acme.corp.hello.NoopAction
vivarcus sdk get acme.corp.hello.NoopAction -o /tmp/NoopAction.go --json
```

`--json` 成功时 stdout 是 `{"class_name","path","size"}`，源码在 `-o` 文件。内置/标准组件会 403（D-6）。

## 启停（Operational Status）

客户入口默认 **active**。止血或对照 UI 时用 CLI（对标 Veeva `PUT /code/{FQN}/enable|disable`）：

```bash
vivarcus sdk disable acme.corp.hello.NoopAction --json
vivarcus sdk enable acme.corp.hello.NoopAction --json
```

仅 **Recordaction / Recordtrigger / Customwebapi** 可启停。`Sdkcode` helper、`Sdkjob` 会失败。内置/标准组件会 403。按钮随 `disable` 从菜单消失；正在执行的实例不受影响。

等价 MDL（一般不必，Agent 优先 CLI）：

```mdl
ALTER Recordaction acme.corp.hello.NoopAction (active(false));
```

## 权限

- import / validate / deploy：Vault Owner 或 `configuration.deployment`
- `sdk get` / `put` / `enable` / `disable`：与 Admin Configuration 元数据权限相同（查看 + 编辑字段）

见 [troubleshooting](troubleshooting.md) 处理常见失败。
