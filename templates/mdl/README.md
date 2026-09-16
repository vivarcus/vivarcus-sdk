# MDL 模板

按顺序执行（替换 `{{...}}` 占位符后）：

| 顺序 | 文件 | 说明 |
|------|------|------|
| 1 | [01-object.mdl](01-object.mdl) | 创建自定义对象与业务字段 |
| 2 | [02-lifecycle.mdl](02-lifecycle.mdl) | 可选：生命周期与状态 |
| 3 | — | 部署 gosdk VPK（见 [docs/05-deploy.md](../../docs/05-deploy.md)） |
| 4 | [03-recordaction-active.mdl](03-recordaction-active.mdl) | 激活 Recordaction / Objectaction |

```bash
ov mdl run templates/mdl/01-object.mdl
# deploy VPK ...
ov mdl run templates/mdl/03-recordaction-active.mdl
```

占位符须与 `sdk_manifest.json` 及 Go `Meta()` 一致。
