# 99-full-stack

端到端示例：MDL 预置 → 构建 wasm → 打 VPK → 部署 → 激活。

## Prerequisites

| 项 | 值 |
|----|-----|
| Vault | 你的 Sandbox 或生产 Vault，Vault Owner |
| 对象 | `demo_request__c`（见下方 MDL） |
| 工具 | Go、TinyGo、`ov-sdk`、`ov` CLI |

## Step 1 — 创建对象（MDL）

编辑并执行 [templates/mdl/01-object.mdl](../../templates/mdl/01-object.mdl)：

- `{{OBJECT}}` → `demo_request__c`
- `{{OBJECT_LABEL}}` → `Demo Request`

```bash
ov mdl run templates/mdl/01-object.mdl
```

## Step 2 — 构建 wasm

```bash
ov-sdk build ./examples/02-update-field -o ./action.wasm
```

编辑 `action.sdk_manifest.json`：

| 字段 | 值 |
|------|-----|
| `component_name` | `com.example.SetTitle` |
| `label` | `Set Title` |
| `object` | `demo_request__c` |
| `object_action` | `demo_request__c.set_title__c` |

## Step 3 — 打包 VPK

```bash
./examples/99-full-stack/scripts/package-vpk.sh ./action.wasm ./action.sdk_manifest.json
```

## Step 4 — 部署

```bash
ov package import ./examples/99-full-stack/dist/demo-action.vpk
ov package validate <package_id>
ov package deploy <package_id> --confirm
```

## Step 5 — 激活

编辑 [templates/mdl/03-recordaction-active.mdl](../../templates/mdl/03-recordaction-active.mdl)：

- `{{ACTION_FQN}}` → `com.example.SetTitle`
- `{{OBJECT_ACTION}}` → `demo_request__c.set_title__c`

```bash
ov mdl run templates/mdl/03-recordaction-active.mdl
```

## Step 6 — 验证

1. 创建 `demo_request__c` 记录
2. 记录详情 → **All Actions** → **Set Title**
3. `title__c` 变为 `from-sdk`

## Agent checklist

见 [AGENTS.md](../../AGENTS.md) 端到端检查清单。
