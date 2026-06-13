# GLOSSARY（术语与权威枚举）

本文件是 ForgeCode 的 **权威术语表**。状态枚举、错误分类、风险与决策枚举以此为准，所有 Spec 必须一致引用。

## 状态枚举（State Enums）

### Agent State（runtime-core 拥有）
```text
Created → Initializing → Thinking → ToolRequested → AwaitingApproval
→ ToolExecuting → Observing → Compacting → Paused → Completed → Failed → Cancelled
```
合法转移见 `docs/specs/runtime-core/SPEC.md` §8。`Paused/Completed/Failed/Cancelled` 为可恢复/终止态。

### Session State（session-store 拥有）
```text
Active → Paused → Completed → Failed → Cancelled
```

### Task State（agent-orchestration 拥有，Team 用）
```text
Pending → Blocked → Ready → Assigned → Running → Reviewing → Completed → Failed → Cancelled
```

### Team State（agent-orchestration 拥有）
```text
Forming → Active → Reviewing → Closed → Failed
```

### Worktree State（git-worktree 拥有）
```text
Creating → Ready → Executing → Diffed → Merging → Merged → Discarded → CleanedUp → Orphaned
```

### MCP Server State（mcp-client 拥有）
```text
Configured → Starting → Initializing → Ready → Unhealthy → Reconnecting → Stopped → Failed
```

### Skill State（extension-system 拥有）
```text
Discovered → Installed → Loaded → Active → Disabled → Uninstalled
```
候选 Skill：`Candidate → StaticChecked → Replayed → Approved → Installed`（未审批不得 Active）。

## 错误分类（Error Categories）

所有模块错误必须可归入下列类别之一（不只返回字符串）：

```text
ValidationError      // 输入/Schema 不合法（L1）
PermissionDenied     // 权限拒绝（L2/L3/L5）
ApprovalRequired     // 需人工审批（非终态错误，触发 AwaitingApproval）
TimeoutError         // 超时
CancelledError       // 取消（Context Cancellation / 用户取消）
ProviderError        // 模型 Provider 错误（含 RateLimit 子类）
ToolExecutionError   // 工具执行失败
SandboxError         // 沙箱相关失败
PersistenceError     // SQLite/存储失败
ConflictError        // 冲突（分支/路径/合并/重复）
RecoveryError        // 恢复/重放失败
```

每类错误标注是否可重试、是否可恢复、是否需审计。

## 决策与风险枚举

- **RiskLevel**：`Low / Medium / High / Critical`
- **Decision**：`Allow / AskOnce / AskAlways / Deny`
- **Hook Result**：`Allow / Deny / Ask / Modify / Continue`
- **Hook Type**：`Internal Go / Shell / HTTP`
- **Trust Level（MCP）**：`Trusted / Limited / Untrusted`

## 常用术语

- **Event Envelope**：统一事件信封，见 `architecture/EVENT_MODEL.md`。
- **Checkpoint**：可回退的状态快照点。
- **Compaction**：上下文压缩。
- **Observation**：工具执行结果反馈给模型的内容。
- **Workspace Root**：允许操作的工作区根目录边界。
- **SubAgent**：一次性任务委派，完成即返回父 Agent。
- **Agent Team**：长期角色 + 共享 Task DAG + 成员通信 + 结果集成。
- **Descriptor / ToolDescriptor**：工具的统一描述（名称、schema、风险、权限要求）。
- **Artifact**：Agent 产出的结构化产物。

## ID 命名规范

- **Module ID**：kebab-case，见 `architecture/MODULE_MAP.md`。
- **Requirement ID**：`FR-<DOMAIN>-NNN` / `NFR-<CAT>-NNN`，见 Master SPEC §6.4/§6.5。
- **Task ID**：`FC-<AREA>-NNN`，见 `tasks/MASTER_TASKS.md`。
- **ADR ID**：`ADR-NNNN`。
- **Risk ID**：`RISK-NNN`。
- **Open Question ID**：`Q<N>`。
