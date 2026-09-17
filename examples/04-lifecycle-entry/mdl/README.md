# 04-lifecycle-entry 演示 MDL

按 **Step 1–7** 顺序放置可执行的 Vault 组件 MDL。对象 / lifecycle / workflow 名写在这些文件里；只有 `{{ACTION_FQN}}` 由 `vivarcus-sdk describe` 在跑 demo 时填入。

| 文件 | 对应 README | 内容 |
|------|-------------|------|
| [01-object.mdl](01-object.mdl) | Prerequisites | 对象与 `title__c` — **VPK `components/00010/`**（非单独 apply-mdl） |
| [02-lifecycle.mdl](02-lifecycle.mdl) | Step 3–4、6 | 生命周期 + submit user action + entry + `create_record` event |
| [03-bind-object-lifecycle.mdl](03-bind-object-lifecycle.mdl) | — | `ALTER Object` 绑定 lifecycle |
| [04-workflow-action.mdl](04-workflow-action.mdl) | Step 5 | Action 步 workflow（`Recordaction.{{ACTION_FQN}}`） |
| [05-workflow-cancel.mdl](05-workflow-cancel.mdl) | Step 7 | 带 usertask + Cancelation Action 的 workflow |

手工应用（需先 `scripts/render-mdl.py mdl/01-object.mdl` 或 export 占位符）：

```bash
vivarcus component apply-mdl --confirm -f /tmp/rendered.mdl
```

共享 `render-mdl.py`：[`../_shared/scripts/render-mdl.py`](../_shared/scripts/render-mdl.py)

规则 XML 片段说明见 [`templates/mdl`](../../templates/mdl/)。
