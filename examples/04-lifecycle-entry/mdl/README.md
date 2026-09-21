# 04-lifecycle-entry 演示 MDL

按 **Step 1–7** 顺序放置可执行的 Vault 组件 MDL。对象 / lifecycle / workflow 名写在这些文件里；只有 `{{ACTION_FQN}}` 由 `go.mod` module + 类型名派生（demo 脚本填入）。

| 文件 | 对应 README | 内容 |
|------|-------------|------|
| [01-object.mdl](01-object.mdl) | Prerequisites | 对象与 `title__c` — `vivarcus component apply-mdl`（**不要**打 VPK） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | Step 3–4、6 | 生命周期 + submit user action + entry + `create_record` event |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | — | `ALTER Object` 绑定 lifecycle |
| [04-workflow-action.mdl](04-workflow-action.mdl) | Step 5 | Action 步 workflow（`Recordaction.{{ACTION_FQN}}`） |
| [05-workflow-cancel.mdl](05-workflow-cancel.mdl) | Step 7 | 带 usertask + Cancelation Action 的 workflow |

手工应用：`01-object.mdl` 无占位符，直接 apply。`02` / `04` / `05` 须先 `sdk put` 投影 `Recordaction`，再展开 `{{ACTION_FQN}}`：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f main.go --json
export ACTION_FQN=acme.corp.lifecycle.StampOnEnter
python3 ../_shared/scripts/render-mdl.py mdl/02-lifecycle.mdl /tmp/02-lifecycle.mdl
vivarcus component apply-mdl --confirm -f /tmp/02-lifecycle.mdl
```

共享 `render-mdl.py`：[`../_shared/scripts/render-mdl.py`](../_shared/scripts/render-mdl.py)

规则 XML 片段说明见 [`templates/mdl`](../../templates/mdl/)。
