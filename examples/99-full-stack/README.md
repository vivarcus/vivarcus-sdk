# 99-full-stack

**多 Action 端到端**：同一 Go module 两个类型，**一份 wasm、一个 combo VPK**。

```bash
vivarcus auth login && vivarcus config set default_vault <uuid>
./scripts/dev-demo.sh
```

| 路径 | 作用 |
|------|------|
| [main.go](main.go) | Set Title + Clear Title |
| [mdl/](mdl/) | 对象 MDL |

导入扫描 `gosdk/action.wasm`，`__sdk_describe` 投影两条 Recordaction。单 Action 精读：[02-update-field](../02-update-field)。
