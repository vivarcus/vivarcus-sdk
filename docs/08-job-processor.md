# Job Processor API

每个 wasm 模块可包含 **一个或多个** Job Processor 入口类型，实现以下接口。

## 必选方法

```go
type Job interface {
    Meta() Meta
    Init(ctx InitContext) (Input, error)
    Process(ctx ProcessContext) (ProcessResult, error)
}
```

对标 Veeva Java SDK `com.veeva.vault.sdk.api.job.Job`。`completeWithSuccess` / `completeWithError` 由平台 job engine 根据 Process 结果写入，客户代码不实现。

### Meta()

声明组件级静态元数据（对标 Java `@JobInfo`）：

| 字段 | 必填 | 说明 |
|------|------|------|
| `Label` | 是 | 显示名称 |
| `Name` | 否 | Sdkjob FQN；空则由 module + 类型名派生 |
| `Idempotent` | 否 | 幂等 |
| `Visible` | 否 | 是否在 Admin 可见 |
| `AdminConfigurable` | 否 | 管理员可配置 |

### Init(ctx)

根据作业参数产出 JobItem 列表。平台按 Jobmetadata `chunk_size`（1–500，默认 500）切成最多 5000 个 task，再串行调用 `Process`。

常用参数读取：

```go
ctx.ParamString("id")
ctx.ParamStrings("record_ids")
```

### Process(ctx)

处理当前 task 的 `ctx.Items`。返回空 `Results` 且 `error == nil` 时，平台将该批全部记为 success。

`platform.Get` / `platform.Update` 在 Job Processor 中**不绑定** Action 上下文记录，可读写任意对象记录。

## 组件命名

FQN 由 module 路径 + 类型名派生，例如 `acme.corp.jobprocessor.StampRecords`。入口类型可放在模块根、`jobs/` 或 `entries/`，**不要**放在 `shared/`。

调度：先部署 Sdkjob，再在 Admin > Operations 创建 SDK Job Metadata（`job_code` 填该 FQN）与类型为 **SDK Job** 的 Job Definition。
