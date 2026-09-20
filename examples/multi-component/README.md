# multi-component

**推荐入口**：`go.mod` + `shared/` + 多个 Action/Trigger 分文件入口，**一个** combo VPK 端到端部署。

```bash
vivarcus auth login && vivarcus config set default_vault <uuid>
cd examples/multi-component
./scripts/dev-demo.sh
```

| 路径 | 作用 |
|------|------|
| [go.mod](go.mod) | 客户 module path → Recordaction FQN 前缀 |
| [shared/](shared/) | 跨入口 helper（UDC 模式，非调度入口） |
| `action_*.go` / `trigger_*.go` | 分文件入口（同一 `package main`） |
| [mdl/](mdl/) | 对象 + lifecycle entry_action MDL |

本地 build 校验：

```bash
cd examples/multi-component
vivarcus-sdk build .
```

`dev-demo.sh` 会 API 验收：

- CREATE 时 `StampName` Record Trigger 给 `name__v` 追加 `-trig`
- 用户按钮 `set_title__c` / `clear_title__c`
- lifecycle `submit__c` 进入 `in_review__c` 时 `StampOnEnter` entry_action 写入 `title__c`

分场景精读：[04-lifecycle-entry](../04-lifecycle-entry)（entry/event/workflow）、[05-stamp-trigger](../05-stamp-trigger)（单一 Trigger）。
