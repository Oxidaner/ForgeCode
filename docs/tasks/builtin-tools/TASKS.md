# Built-in Tools Tasks

> 模块 `builtin-tools` 任务清单。任务遵循 `docs/templates/TASK_TEMPLATE.md`。前缀 `FC-BT`。

## Task: FC-BT-001 — Built-in Tool 装配与 Descriptor 骨架

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-001 |
| Title | Built-in Tool 装配与 Descriptor 骨架 |
| Module | `builtin-tools` |
| Type | Architecture |
| Priority | P0 |
| Milestone | M1 |
| Status | Ready |
| Size | M |
| Dependencies | - |
| Related Requirements | FR-TOOL-001, FR-TOOL-004 |
| Related Spec Sections | `builtin-tools` §3, §6, §7 |

**Description**：搭建本模块骨架：定义六个工具的 `ToolDescriptor`（名称、JSON Schema、来源 `builtin`、默认风险、权限要求声明），实现 `RegisterBuiltins(reg, deps)` 把工具注册进 `tool-runtime.Registry`，定义 `Deps`（Checkpointer / WorkspaceRoot / Clock / Limits）注入结构。输入：tool-runtime 的 `Tool`/`Registry` 接口；输出：可注册、可发现的内置工具集合（Execute 可为占位）。

**Implementation Notes**：不重定义 `Tool`/`Registry`/管线（属 tool-runtime）。Descriptor 风险标注：Read/Glob/Grep=Low，Write/Edit=Medium，Bash=High。JSON Schema 与 §6 输入结构一致。Clock 注入便于超时测试。

**Files or Packages Likely Affected**：`internal/builtin/registry.go`、`internal/builtin/descriptors.go`、`internal/builtin/deps.go`。

**Tests Required**：Contract Test（六工具满足 Tool 接口、Descriptor 通过 schema 校验）；Unit（RegisterBuiltins 注册数量与名称）。

**Security Considerations**：确保 Bash Descriptor 默认 High 风险，保证进入审批路径；不在装配层绕过管线。

**Acceptance Criteria**：
- [ ] 六个工具 Descriptor 定义完整且通过 schema 校验
- [ ] `RegisterBuiltins` 成功注册六工具且名称唯一
- [ ] Bash 默认风险标注为 High

**Definition of Done**：代码 + Contract/Unit 测试 + Spec §6 对齐 + Checklist 勾选 + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-002 — ReadFile 工具（分页 / 二进制识别 / 超大文件保护）

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-002 |
| Title | ReadFile 工具（分页 / 二进制识别 / 超大文件保护） |
| Module | `builtin-tools` |
| Type | Implementation |
| Priority | P0 |
| Milestone | M1 |
| Status | Backlog |
| Size | M |
| Dependencies | FC-BT-001 |
| Related Requirements | FR-TOOL-101, NFR-LIMIT-001 |
| Related Spec Sections | `builtin-tools` §6, §9, §10 |

**Description**：实现 `ReadFile`：按 `offset`/`limit`（行）分页读取，输出带行号；通过采样字节探测 NUL/非法 UTF-8 识别二进制并拒绝；文件超过 `maxBytes` 时返回保护性 `ValidationError` 提示分页。输出 `ToolResult`，Meta 含起止行号与总行数。

**Implementation Notes**：探测采样 `binaryProbeBytes` 字节；未指定 limit 用 `defaultLimit`；超大文件不读入全文。尊重 ctx 取消。

**Files or Packages Likely Affected**：`internal/builtin/readfile.go`。

**Tests Required**：Unit（空文件、单行、offset 越界、二进制拒绝）；Golden（分页输出带行号）；Failure（超大文件保护）。

**Security Considerations**：拒绝二进制/超大文件，降低密钥/二进制整体回灌上下文风险（边界判定仍属 permission-engine）。

**Acceptance Criteria**：
- [ ] offset/limit 分页正确，输出带行号
- [ ] 二进制文件返回 ValidationError
- [ ] 超大文件触发保护并提示分页

**Definition of Done**：代码 + Unit/Golden/Failure 测试 + Checklist + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-003 — WriteFile / EditFile 工具（Checkpoint + 原子写 + 唯一匹配 + Diff）

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-003 |
| Title | WriteFile / EditFile 工具（Checkpoint + 原子写 + 唯一匹配 + Diff） |
| Module | `builtin-tools` |
| Type | Implementation |
| Priority | P0 |
| Milestone | M2 |
| Status | Backlog |
| Size | L |
| Dependencies | FC-BT-001 |
| Related Requirements | FR-TOOL-102, FR-TOOL-103, NFR-RECOV-001 |
| Related Spec Sections | `builtin-tools` §6, §9, §11 |

**Description**：实现 `WriteFile`（创建/覆盖：覆盖既有文件前经 `Checkpointer` 创建 Checkpoint，临时文件 + 原子 rename，写失败保持原文件不变）与 `EditFile`（要求 `old_string` 在文件中唯一匹配，0/多匹配返回 `ValidationError`；成功产出 unified diff，写前同样 Checkpoint）。

**Implementation Notes**：rename 前 fsync 临时文件；Checkpoint 失败即中止写（PersistenceError）。EditFile 命中次数写入错误信息；diff 放入 Meta。创建新文件无需 Checkpoint。

**Files or Packages Likely Affected**：`internal/builtin/writefile.go`、`internal/builtin/editfile.go`、`internal/builtin/diff.go`。

**Tests Required**：Unit（创建 vs 覆盖、EditFile 0/1/多匹配）；Golden（diff 输出）；Failure（写临时文件失败 / rename 失败 / Checkpoint 失败 → 原文件字节级不变）。

**Security Considerations**：写前 Checkpoint 保证危险写可 `/rewind`；不绕过 permission-engine 的写边界决策。

**Acceptance Criteria**：
- [ ] 覆盖前创建 Checkpoint，写失败原文件不变
- [ ] EditFile 仅唯一匹配时替换，否则 ValidationError
- [ ] EditFile 成功产出正确 Diff

**Definition of Done**：代码 + Unit/Golden/Failure 测试 + Checklist + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-004 — Bash 工具（执行 / 超时 / 头尾保留 / 退出码与错误分类）

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-004 |
| Title | Bash 工具（执行 / 超时 / 头尾保留 / 退出码与错误分类） |
| Module | `builtin-tools` |
| Type | Implementation |
| Priority | P0 |
| Milestone | M1 |
| Status | Backlog |
| Size | L |
| Dependencies | FC-BT-001 |
| Related Requirements | FR-TOOL-104, NFR-LIMIT-001 |
| Related Spec Sections | `builtin-tools` §6, §9, §12, §13 |

**Description**：实现 `Bash`：以子进程组执行命令，强制 Deadline 超时（入参可收窄，受 `maxTimeoutMs` 上限），输出超 `headBytes+tailBytes` 时保留头尾并标记截断，返回退出码并据此分类错误（非零→`ToolExecutionError`，超时→`TimeoutError`，取消→`CancelledError`）。

**Implementation Notes**：超时/取消终止整个进程组避免孤儿；合并 stdout/stderr 或分别保留头尾（择一并记录）；退出码写 Meta。MVP 本地执行，沙箱化由 sandbox 经管线接管。

**Files or Packages Likely Affected**：`internal/builtin/bash.go`、`internal/builtin/output_buffer.go`。

**Tests Required**：Unit（退出码映射、头尾截断）；Failure Injection（超时杀进程组、被取消、孤儿进程检查）；Race（并发执行）。

**Security Considerations**：默认 High 风险进入审批；Bash 危险分析属 permission-engine（RISK-006），本模块不内联放行；输出可能含密钥，不写普通日志。

**Acceptance Criteria**：
- [ ] 超时终止进程组并返回 TimeoutError 且保留头尾输出
- [ ] 非零退出码归类 ToolExecutionError 且退出码入 Meta
- [ ] ctx 取消时无孤儿进程残留

**Definition of Done**：代码 + Unit/Failure/Race 测试 + Checklist + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-005 — Glob / Grep 工具（模式匹配 / 正则 / 去重 / 上限）

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-005 |
| Title | Glob / Grep 工具（模式匹配 / 正则 / 去重 / 上限） |
| Module | `builtin-tools` |
| Type | Implementation |
| Priority | P0 |
| Milestone | M1 |
| Status | Backlog |
| Size | M |
| Dependencies | FC-BT-001 |
| Related Requirements | FR-TOOL-105, FR-TOOL-106, NFR-LIMIT-001 |
| Related Spec Sections | `builtin-tools` §6, §9, §10 |

**Description**：实现 `Glob`（按 glob 模式匹配 Workspace 内文件，返回排序路径，受 `maxResults` 限制）与 `Grep`（字符串/正则搜索文件集合，结果按 `file:line` 去重，受 `maxMatches` 限制，超出标记截断）。

**Implementation Notes**：Glob 结果稳定排序；Grep 正则编译失败 → `ValidationError`；遍历尊重 ctx 取消；去重键 `file:line`。

**Files or Packages Likely Affected**：`internal/builtin/glob.go`、`internal/builtin/grep.go`。

**Tests Required**：Unit（glob 无匹配 / 多匹配排序、grep 正则、grep 多匹配去重、上限截断、非法正则报错）；Golden（grep 去重结果）。

**Security Considerations**：搜索范围受 WorkspaceRoot 与 permission-engine 边界约束；不绕过管线。

**Acceptance Criteria**：
- [ ] Glob 返回排序路径并受 maxResults 限制
- [ ] Grep 支持正则且结果按 file:line 去重
- [ ] 非法正则返回 ValidationError，结果超限标记截断

**Definition of Done**：代码 + Unit/Golden 测试 + Checklist + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-006 — 工具测试矩阵（Unit / Golden / Failure / Race / Contract）

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-006 |
| Title | 工具测试矩阵（Unit / Golden / Failure / Race / Contract） |
| Module | `builtin-tools` |
| Type | Test |
| Priority | P0 |
| Milestone | M2 |
| Status | Backlog |
| Size | M |
| Dependencies | FC-BT-002, FC-BT-003, FC-BT-004, FC-BT-005 |
| Related Requirements | FR-TOOL-101, FR-TOOL-102, FR-TOOL-103, FR-TOOL-104, FR-TOOL-105, FR-TOOL-106, NFR-TEST-001 |
| Related Spec Sections | `builtin-tools` §16, §17 |

**Description**：建立六工具的综合测试矩阵与 Fake 边界（Fake Checkpointer / Fake Clock / 临时目录 FS），覆盖 §16 全部测试类型与 §17 验收条件，纳入 CI 并启用 `-race`。

**Implementation Notes**：Golden 文件集中管理（ReadFile 分页、EditFile diff、Grep 去重）；Failure 用可注入错误的 FS/Checkpointer；Contract 复用 tool-runtime 的接口契约测试。

**Files or Packages Likely Affected**：`internal/builtin/*_test.go`、`internal/builtin/testdata/`。

**Tests Required**：Unit + Golden + Failure Injection + Race + Contract 全量。

**Security Considerations**：包含验证 Bash 进审批、ReadFile 拒绝二进制/超大的安全相关断言。

**Acceptance Criteria**：
- [ ] §17 全部验收条件有对应测试
- [ ] `go test -race ./internal/builtin/...` 通过
- [ ] Golden 与 Failure 用例纳入 CI

**Definition of Done**：测试代码 + CI 配置 + Checklist + Evidence。

**Evidence**：（完成后填写）

---

## Task: FC-BT-007 — 内置工具安全与管线一致性验证

| 字段 | 值 |
| --- | --- |
| ID | FC-BT-007 |
| Title | 内置工具安全与管线一致性验证 |
| Module | `builtin-tools` |
| Type | Security |
| Priority | P0 |
| Milestone | M2 |
| Status | Backlog |
| Size | S |
| Dependencies | FC-BT-001, FC-BT-004 |
| Related Requirements | FR-TOOL-102, FR-TOOL-104, NFR-SEC-001, NFR-RECOV-001 |
| Related Spec Sections | `builtin-tools` §14, §15 |

**Description**：验证内置工具不绕过统一管线：所有工具调用必经 Validation→Permission→Hook→Execute→Audit；Bash 默认进审批；写类工具写前 Checkpoint；工具不直接写 Event Store/Approval。建立 Security Test 防止"工具内联放行"回归。

**Implementation Notes**：用集成测试断言每条工具路径触发 permission-engine 决策与审计事件；断言 WriteFile/EditFile 写前必产生 CheckpointCreated。

**Files or Packages Likely Affected**：`internal/builtin/security_test.go`。

**Tests Required**：Security Test（审批绕过防护、管线一致性）；Integration（Checkpoint 先于写）。

**Security Considerations**：直接对应 RISK-006 与审批绕过防护；敏感输出不入普通日志验证。

**Acceptance Criteria**：
- [ ] 任一工具路径均经过 Permission 与 Audit
- [ ] 写类工具在写入前必有 Checkpoint 事件
- [ ] 工具不直接写 Event Store / Approval

**Definition of Done**：Security/Integration 测试 + Checklist + Evidence。

**Evidence**：（完成后填写）
