# model-provider Tasks

模块：`model-provider`。Task 前缀 `FC-PROV`。相关需求 FR-PROVIDER-001..006。ADR-0003。RISK-004。

## FC-PROV-001 — 定义中立 Provider 接口与消息模型
| 字段 | 值 |
| --- | --- |
| ID | FC-PROV-001 | 
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Ready | Size | M |
| Dependencies | - | Related Requirements | FR-PROVIDER-001, FR-PROVIDER-006 | Spec | §6 |

**Description**：定义 `Provider` 接口、`ChatRequest/ChatResponse`、中立 Message/ToolCall/Usage/StopReason 结构。
**Implementation Notes**：私有结构禁止泄漏到 runtime-core；接口在本模块定义、构造注入。
**Files**：`internal/model-provider/provider.go`, `types.go`。
**Tests Required**：接口契约骨架、类型序列化 Unit。
**Security Considerations**：响应视为不可信输入。
**Acceptance Criteria**：
- [ ] 接口不含任何 Provider 私有字段
- [ ] runtime-core 仅 import 本接口
**Definition of Done**：接口评审通过 + Contract Test 骨架。
**Evidence**：

## FC-PROV-002 — Mock Provider
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Ready | Size | S |
| Dependencies | FC-PROV-001 | Related Requirements | FR-PROVIDER-003 |

**Description**：可编程 Mock Provider，支持预设响应/工具调用/错误/延迟，供全栈测试。
**Tests Required**：Mock 自身 Unit。
**Acceptance Criteria**：
- [ ] 可注入普通响应、多 Tool Call、错误、超时
- [ ] 确定性可重放
**Definition of Done**：Mock 可被其他模块测试复用。

## FC-PROV-003 — 普通响应与 Streaming
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-001 | Related Requirements | FR-PROVIDER-001, FR-PROVIDER-002 |

**Description**：统一非流式与流式（`StreamChunk`）输出，解析 Stop Reason 与 Token Usage。
**Implementation Notes**：流式中断丢弃半成品按轮重试（见 FAILURE_AND_RECOVERY）。
**Tests Required**：流式聚合 Unit、中断 Failure Injection。
**Acceptance Criteria**：
- [ ] 流式可聚合为完整响应
- [ ] Usage 与 Stop Reason 正确解析
**Definition of Done**：流式与非流式行为一致。

## FC-PROV-004 — Tool Calling 与多 Tool Call
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-003 | Related Requirements | FR-PROVIDER-002 |

**Description**：统一工具调用请求/结果在消息序列中的表达，支持单次多 Tool Call。
**Tests Required**：多 Tool Call 解析 Contract Test。
**Acceptance Criteria**：
- [ ] 多个 Tool Call 顺序与 ID 保留
- [ ] 工具结果回填格式中立
**Definition of Done**：与 tool-runtime 契约对齐。

## FC-PROV-005 — 错误归一化、重试、超时、限流
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-001 | Related Requirements | FR-PROVIDER-004, NFR-REL-002 |

**Description**：将各 Provider 错误归一为 `ProviderError`（含 RateLimit 子类），可重试/不可重试区分，指数退避，超时控制。
**Tests Required**：Failure Injection（429/5xx/超时/网络中断）。
**Acceptance Criteria**：
- [ ] 瞬时错误退避重试（默认 ≤3）
- [ ] 不可重试错误快速失败
**Definition of Done**：错误分类覆盖主流 Provider。

## FC-PROV-006 — OpenAI 适配器与消息转换
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-003, FC-PROV-004, FC-PROV-005 | Related Requirements | FR-PROVIDER-003, FR-PROVIDER-006 |

**Description**：OpenAI provider-specific 消息/工具转换，私有结构封闭在适配器内。
**Tests Required**：Contract Test（与 Mock 同套用例）。
**Acceptance Criteria**：
- [ ] 通过统一 Contract Test 套件
**Definition of Done**：可与 runtime-core 跑通真实/录制响应。

## FC-PROV-007 — 能力元数据与 Structured Output
| Type | Implementation | Priority | P1 | Milestone | M3 | Status | Backlog | Size | S |
| Dependencies | FC-PROV-001 | Related Requirements | FR-PROVIDER-005 |

**Description**：Model Capability、Context Window 元数据；Structured Output（经 schema 约束，见 OPEN_QUESTIONS Q11）。
**Acceptance Criteria**：
- [ ] context-manager 可查询窗口大小
- [ ] 不支持原生 JSON Schema 时回退提示+解析校验

## FC-PROV-008 — Anthropic 与 OpenAI-Compatible 适配器
| Type | Implementation | Priority | P1 | Milestone | M6 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-006 | Related Requirements | FR-PROVIDER-003 |

**Description**：补充 Anthropic 与 OpenAI-Compatible 适配器，通过同套 Contract Test。
**Acceptance Criteria**：
- [ ] 三个真实 Provider 通过 Contract Test（RISK-004）

## FC-PROV-009 — Provider Contract Test 套件
| Type | Test | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-PROV-002 | Related Requirements | FR-PROVIDER-001..006 |

**Description**：一套对所有适配器统一运行的 Contract Test，验证接口行为一致。
**Acceptance Criteria**：
- [ ] 所有适配器（含 Mock）共享同套用例
- [ ] CI 中可对 Mock 运行，真实 Provider 可选
