# 限制与 Phase 1 边界

## 执行配额

| 项 | 限制 |
|----|------|
| 超时 | 30 秒 |
| 内存 | 128 MB |
| wasm 体积 | 2 MB |
| 调用链深度 | 10 层（Trigger 级联，Action 间接触发） |

## 标准库：禁止使用

平台编译拒绝（打进 `gosdk/` 的源码）：

- `os`、`net`、`net/http`
- `database/sql`、`os/exec`
- `unsafe`、`syscall`、`testing`
- `crypto/rand`

`testing` 只用于本地 `*_test.go`；`sdk put` 与 `package-vpk.sh` 都不部署测试文件。入口代码不要 import `testing`。

## 标准库：受限

- `encoding/json` — 勿对 `map[string]any` / `interface{}` 做反射序列化
- `sync` — wasm 单线程，Mutex 为 no-op
- `time.Sleep` — 不可用

## Phase 1 宿主能力

客户 Record Action 跑在 **wasm 沙箱**里，不是 Veeva Java SDK 的 Vault Owner 服务账号。`platform.Get` / `platform.Update` 的 `object` 与 `recordID` **必须等于**当前 Action 上下文记录；传入其他 ID 返回 `record_action_host_call_failed`。改当前记录请用 `rec.SetValue` 或对上下文 ID 调用 `platform.Update`。

| 能力 | 状态 |
|------|------|
| `platform.Get` | Action：仅当前上下文记录。Job Processor：任意记录 |
| `platform.Update` | Action：仅当前上下文记录。Job Processor：任意记录；未知字段 / 保留键拒绝；单次字段数与单值体积有上限 |
| `platform.LogInfo` 等 | 支持 |
| `platform.Create` / `Delete` | 视宿主实现，可能 NOT_IMPLEMENTED |
| `platform.Query`（只读 VQL） | 支持；`PAGESIZE` 上限 200；不支持 `PAGESIZE 0`（纯 count）、`PAGEOFFSET`、分页游标 |
| 通知 | **不支持** |
| HTTP 出站 | **不支持** |
| 生命周期切换 / 工作流 | **不支持**（host API） |

## 批量 Action

- `UsageUserBulkAction` 不能与其他 Usage 混用
- 批量执行：**逐条串行**，每条独立事务（与 Veeva 500 条一批不同）

## 部署

- `javasdk/` VPK → `not_supported__v`
- 客户 Action 部署后默认 **active**
- 管理员可用 `vivarcus sdk disable <FQN>`（或 `ALTER active(false)`）停用

## Record Trigger

规格见平台文档。客户 wasm Trigger 与 Action 可同 VPK 部署。

## Job Processor

规格见 [08-job-processor](08-job-processor.md)。客户 `Sdkjob` 经 gosdk VPK 部署后，由 Job Definition 类型 **SDK Job** 调度。

## 错误码（用户可见）

| 码 | 含义 |
|----|------|
| `record_action_inactive` | 未激活 |
| `record_action_checksum_mismatch` | wasm 被篡改 |
| `record_action_timeout` | 超时 |
| `record_action_wasm_trap` | wasm panic |
| `record_action_execution_failed` | 业务返回 error |

完整错误码列表见 Vivarcus 平台文档或联系支持获取。
