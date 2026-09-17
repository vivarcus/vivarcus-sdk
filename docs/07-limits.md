# 限制与 Phase 1 边界

## 执行配额

| 项 | 限制 |
|----|------|
| 超时 | 30 秒 |
| 内存 | 128 MB |
| wasm 体积 | 2 MB |
| 调用链深度 | 10 层（Trigger 级联，Action 间接触发） |

## 标准库：禁止使用

编译期拒绝：

- `os`、`net`、`net/http`
- `database/sql`、`os/exec`
- `unsafe`、`syscall`、`testing`
- `crypto/rand`

## 标准库：受限

- `encoding/json` — 勿对 `map[string]any` / `interface{}` 做反射序列化
- `sync` — wasm 单线程，Mutex 为 no-op
- `time.Sleep` — 不可用

## Phase 1 宿主能力

客户 Record Action 跑在 **wasm 沙箱**里，不是 Veeva Java SDK 的 Vault Owner 服务账号。`platform.Get` / `platform.Update` 的 `object` 与 `recordID` **必须等于**当前 Action 上下文记录；传入其他 ID 返回 `record_action_host_call_failed`。改当前记录请用 `rec.SetValue` 或对上下文 ID 调用 `platform.Update`。

| 能力 | 状态 |
|------|------|
| `platform.Get` | 仅当前上下文记录 |
| `platform.Update` | 仅当前上下文记录；未知字段 / 保留键拒绝；单次字段数与单值体积有上限 |
| `platform.LogInfo` 等 | 支持 |
| `platform.Create` / `Delete` | 视宿主实现，可能 NOT_IMPLEMENTED |
| VQL 查询 | **不支持** |
| 通知 | **不支持** |
| HTTP 出站 | **不支持** |
| 生命周期切换 / 工作流 | **不支持**（host API） |

## 批量 Action

- `UsageUserBulkAction` 不能与其他 Usage 混用
- 批量执行：**逐条串行**，每条独立事务（与 Veeva 500 条一批不同）

## 部署

- `javasdk/` VPK → `not_supported__v`
- 客户 Action 部署后默认 **active**
- 管理员可用 `ALTER active(false)` 停用

## Record Trigger

规格见平台文档；**Phase 1 运行时未实现** Trigger 客户 wasm 路径。仅 Record Action 可端到端部署。

## 错误码（用户可见）

| 码 | 含义 |
|----|------|
| `record_action_inactive` | 未激活 |
| `record_action_checksum_mismatch` | wasm 被篡改 |
| `record_action_timeout` | 超时 |
| `record_action_wasm_trap` | wasm panic |
| `record_action_execution_failed` | 业务返回 error |

完整错误码列表见 Vivarcus 平台文档或联系支持获取。
