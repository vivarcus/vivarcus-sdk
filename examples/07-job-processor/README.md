# 07-job-processor

**Job Processor** 示例：`Init` 从作业参数 `record_ids` 产出 JobItem，`Process` 给每条 `sdk_demo__c` 记录写 `stamped__c=true`。

对标 Veeva Java SDK `com.veeva.vault.sdk.api.job.Job`（`@JobInfo` + `init` / `process`）。`completeWithSuccess` / `completeWithError` 由平台 job engine 根据 Process 结果写入，客户代码不实现。

**不要打 VPK。** 多文件 combo 见 [multi-component](../multi-component)。

## Quick start

```bash
cd examples/07-job-processor
vivarcus sdk put -f main.go --json
```

FQN：`acme.corp.jobprocessor.StampRecords`。`Process` 写的对象须已存在（可先跑 [02-update-field](../02-update-field) 的 `mdl/01-object.mdl`）。

调度：Admin > Operations > **Job Definitions**，类型 **SDK Job**，绑定 SDK Job Metadata（`job_code` 为本 FQN）。
