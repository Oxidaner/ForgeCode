# runtime-core Tasks

模块：`runtime-core`。Task 前缀 `FC-RT`。相关需求 FR-RUNTIME-001..006。ADR-0001/0002/0003。RISK-001/003/020。

## FC-RT-000 — Spike：Go 版本与运行时基线
| Type | Spike | Priority | P0 | Milestone | M1 | Status | Ready | Size | XS |
| Dependencies | - | Related Requirements | NFR-MAINT-001 |

**需回答的问题**：目标 Go 版本（OPEN_QUESTIONS Q1）、context/slog/errors 用法基线。
**最小实验**：建最小模块编译验证泛型与 `log/slog`。
**输出决策**：确定 Go 版本，更新 AGENTS.md 与 Q1。
**结束条件**：版本写入文档。
**Evidence**：

## FC-RT-001 — Agent 状态机
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Ready | Size | M |
| Dependencies | FC-RT-000, FC-EVT-001 | Related Requirements | FR-RUNTIME-001 | Spec | §8 |

**Description**：实现 Agent State 枚举与合法转移表，非法转移拒绝并记录 AgentStateChanged。
**Implementation Notes**：转移表数据化，便于 Golden 测试。
**Files**：`internal/runtime-core/state.go`。
**Tests Required**：转移 Unit + 非法转移 Golden。
**Acceptance Criteria**：
- [ ] 所有 §8 转移合法、其余被拒
- [ ] 每次转移产生事件
**Definition of Done**：状态机评审 + 测试通过。

## FC-RT-002 — Runtime Coordinator 与依赖装配
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-RT-001 | Related Requirements | FR-RUNTIME-001, FR-PROVIDER-001, FR-TOOL-002 |

**Description**：定义 Coordinator，注入 Provider/Invoker/Builder/Store/Bus 接口，管理单 Session 运行。
**Security Considerations**：强制工具调用经 Invoker（统一管线）。
**Acceptance Criteria**：
- [ ] 仅依赖接口，无具体 Provider/Tool 类型
- [ ] 装配可被 Mock 替换
**Definition of Done**：Integration 跑通空循环。

## FC-RT-003 — Agent Loop 主循环
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | L |
| Dependencies | FC-RT-002 | Related Requirements | FR-RUNTIME-001, FR-RUNTIME-004 |

**Description**：实现 Step：调用模型→解析 Tool Call→经管线执行→读 Observation→判断继续/完成。
**Implementation Notes**：非法 Tool Call 作为 Observation 反馈，不直接执行。
**Tests Required**：Golden（Mock 响应序列→状态轨迹）。
**Acceptance Criteria**：
- [ ] 只读任务可跑完整循环至 Completed
- [ ] 非法 Tool Call 被安全处理
**Definition of Done**：与 builtin-tools 只读工具集成通过。

## FC-RT-004 — 预算、上限与 Deadline
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-RT-003 | Related Requirements | FR-RUNTIME-002 |

**Description**：BudgetController：max turns/tool-calls/token/cost/deadline，触发安全终止并落 BudgetExceeded。
**Acceptance Criteria**：
- [ ] 每类上限均有触发测试
- [ ] 终止为 Failed 且事件完整
**Definition of Done**：Unit 覆盖全部上限。

## FC-RT-005 — 循环与重复检测
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-RT-003 | Related Requirements | FR-RUNTIME-004 |

**Description**：LoopDetector：重复工具调用指纹与相同错误循环检测，命中强制终止 + LoopDetected。
**Tests Required**：Failure Injection 构造循环。
**Acceptance Criteria**：
- [ ] 重复调用窗口命中即终止（RISK-003）
- [ ] 相同错误连续 N 次终止
**Definition of Done**：循环场景测试通过。

## FC-RT-006 — 取消、暂停与恢复
| Type | Implementation | Priority | P0 | Milestone | M4 | Status | Backlog | Size | L |
| Dependencies | FC-RT-003, FC-SESS-003 | Related Requirements | FR-RUNTIME-003, FR-RUNTIME-006, FR-SESSION-003 |

**Description**：Context Cancellation 传播、用户 Pause/Cancel、Resume 从事件重建状态。
**Implementation Notes**：恢复不重放外部副作用，使用已记录 ToolResult。
**Tests Required**：Recovery Test（杀进程后 Resume）、取消传播 Integration。
**Security Considerations**：审批未完成的高危操作恢复后不自动执行。
**Acceptance Criteria**：
- [ ] 取消传播到正在执行工具
- [ ] Resume 后状态与崩溃前一致（NFR-REL-001）
**Definition of Done**：Recovery Test 通过。

## FC-RT-007 — Provider 错误重试集成
| Type | Implementation | Priority | P1 | Milestone | M1 | Status | Backlog | Size | S |
| Dependencies | FC-RT-003, FC-PROV-005 | Related Requirements | FR-RUNTIME-005, NFR-REL-002 |

**Description**：在循环中集成 Provider 重试/退避，区分可重试，超限转 Failed。
**Acceptance Criteria**：
- [ ] 瞬时错误自动重试
- [ ] 超限落 ModelCallFailed + Failed

## FC-RT-008 — 压缩触发与 Checkpoint 协调
| Type | Implementation | Priority | P0 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-RT-003, FC-CTX-004 | Related Requirements | FR-CONTEXT-004 |

**Description**：检测 token 阈值触发 Compacting，压缩前创建 Checkpoint，压缩失败回滚。
**Acceptance Criteria**：
- [ ] 超阈值进入 Compacting
- [ ] 压缩失败回滚 Checkpoint（RISK-009）

## FC-RT-009 — runtime-core 测试套件
| Type | Test | Priority | P0 | Milestone | M4 | Status | Backlog | Size | M |
| Dependencies | FC-RT-006 | Related Requirements | FR-RUNTIME-001..006 |

**Description**：汇总 Golden/Integration/Failure/Recovery/Race 测试。
**Acceptance Criteria**：
- [ ] `go test -race` 通过
- [ ] 六条核心流程覆盖
