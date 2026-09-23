# 07-job-processor

**Job Processor** 示例：`Init` 从作业参数 `record_ids` 产出 JobItem，`Process` 给每条 `sdk_demo__c` 记录写 `stamped__c=true`。

实现 `Init` / `Process` 处理批量作业项。`completeWithSuccess` / `completeWithError` 由平台 job engine 根据 Process 结果写入，客户代码不实现。

## 部署

**前置**：Vault 地址、Vault ID、session（`auth login` 或 `examples/deploy.env`），见 [examples/README](../README.md)。目录：`cd examples/07-job-processor`（monorepo：`cd sdk/examples/07-job-processor`）。

**一键**（推荐）：`./deploy.sh`

手动等价：

```bash
vivarcus component apply-mdl --confirm -f mdl/01-object.mdl
vivarcus sdk put -f jobs/stamp_records.go --json
```

FQN：`acme.corp.sdkdemo.StampRecords`。`Process` 写 `sdk_demo__c.stamped__c`，须先 `apply-mdl`（与 [02-update-field](../02-update-field) 同对象定义）。

调度：Admin > Operations > **Job Definitions**，类型 **SDK Job**，绑定 SDK Job Metadata（`job_code` 为本 FQN）。

集成测试（已登录 Vault）：`go test -tags=integration -count=1 -timeout 20m .`
