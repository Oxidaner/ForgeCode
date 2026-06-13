# TRACEABILITY_MATRIX

需求 → 模块 → Spec 章节 → Task → Test/Eval → Checklist → Status 的追踪。所有 P0 需求至少对应一个模块、一个 Spec 章节、一个 Task、一个测试/Eval、一个 Checklist 验收项。Status 反映规划阶段（均为 Planned，待实现推进）。

## 功能需求（Functional）

| Requirement | Module | Spec Section | Task | Test/Eval | Checklist | Status |
| --- | --- | --- | --- | --- | --- | --- |
| FR-RUNTIME-001 | runtime-core | §8 | FC-RT-001/003 | 状态转移 Golden | runtime-core/Impl | Planned |
| FR-RUNTIME-002 | runtime-core | §6/§10 | FC-RT-004 | 上限 Unit | runtime-core/Impl | Planned |
| FR-RUNTIME-003 | runtime-core | §9/§12 | FC-RT-006 | 取消 Integration | runtime-core/Test | Planned |
| FR-RUNTIME-004 | runtime-core | §9 | FC-RT-005 | 循环 Failure | runtime-core/Impl | Planned |
| FR-RUNTIME-005 | runtime-core | §13 | FC-RT-007 | Provider Failure | runtime-core/Impl | Planned |
| FR-RUNTIME-006 | runtime-core, session-store | §9 | FC-RT-006, FC-SESS-003 | Recovery | runtime-core/Test | Planned |
| FR-PROVIDER-001 | model-provider | §6 | FC-PROV-001/003 | Contract | model-provider/Impl | Planned |
| FR-PROVIDER-002 | model-provider | §6 | FC-PROV-004 | Contract | model-provider/Impl | Planned |
| FR-PROVIDER-003 | model-provider | §6 | FC-PROV-002/006/008 | Contract | model-provider/Test | Planned |
| FR-PROVIDER-004 | model-provider | §13 | FC-PROV-005 | Failure | model-provider/Impl | Planned |
| FR-PROVIDER-005 | model-provider | §6 | FC-PROV-007 | Unit | model-provider/Impl | Planned |
| FR-PROVIDER-006 | model-provider | §6 | FC-PROV-001/006 | Contract | model-provider/Impl | Planned |
| FR-TOOL-001 | tool-runtime | §6 | FC-TOOL-001 | Unit/Contract | tool-runtime/Design | Planned |
| FR-TOOL-002 | tool-runtime | §8/§9 | FC-TOOL-002 | 管线 Integration/Security | tool-runtime/Impl | Planned |
| FR-TOOL-003 | tool-runtime | §9 | FC-TOOL-003 | Failure | tool-runtime/Impl | Planned |
| FR-TOOL-004 | tool-runtime, mcp-client | §9 | FC-TOOL-006, FC-MCP-004 | Contract | tool-runtime/Test | Planned |
| FR-TOOL-101 | builtin-tools | §ReadFile | FC-BT-002 | Unit/Golden | builtin-tools/Impl | Planned |
| FR-TOOL-102 | builtin-tools | §WriteFile | FC-BT-003 | Failure | builtin-tools/Impl | Planned |
| FR-TOOL-103 | builtin-tools | §EditFile | FC-BT-003 | Golden | builtin-tools/Impl | Planned |
| FR-TOOL-104 | builtin-tools | §Bash | FC-BT-004 | Unit/Failure | builtin-tools/Impl | Planned |
| FR-TOOL-105 | builtin-tools | §Glob | FC-BT-005 | Unit | builtin-tools/Impl | Planned |
| FR-TOOL-106 | builtin-tools | §Grep | FC-BT-005 | Unit | builtin-tools/Impl | Planned |
| FR-PERM-001 | permission-engine | §6 | FC-PERM-001 | Unit/Contract | permission-engine/Design | Planned |
| FR-PERM-002 | permission-engine | §3(L1) | FC-PERM-002 | Security | permission-engine/Impl | Planned |
| FR-PERM-003 | permission-engine | §3(L2)/§14 | FC-PERM-003 | Security 语料库 | permission-engine/Impl | Planned |
| FR-PERM-004 | permission-engine | §3(L3) | FC-PERM-004 | Unit | permission-engine/Impl | Planned |
| FR-PERM-005 | permission-engine | §6/§14 | FC-PERM-005 | Golden/Security | permission-engine/Impl | Planned |
| FR-PERM-006 | permission-engine, sandbox | §3(L4) | FC-PERM-006, FC-SBX-001 | Integration | sandbox/Impl | Planned |
| FR-PERM-007 | permission-engine, telemetry | §3(L5) | FC-PERM-007, FC-TEL-003 | Audit | permission-engine/Impl | Planned |
| FR-PERM-008 | permission-engine | §9 | FC-PERM-001/008 | Unit | permission-engine/Impl | Planned |
| FR-SANDBOX-001 | sandbox | §6 | FC-SBX-001 | Docker Integration | sandbox/Impl | Planned |
| FR-SANDBOX-002 | sandbox | §6 | FC-SBX-002/004 | Security | sandbox/Impl | Planned |
| FR-SANDBOX-003 | sandbox, permission-engine | §9 | FC-SBX-003 | Failure | sandbox/Impl | Planned |
| FR-EVENT-001 | event-system | §6 | FC-EVT-001 | Contract | event-system/Design | Planned |
| FR-EVENT-002 | event-system | §9 | FC-EVT-002 | Race | event-system/Impl | Planned |
| FR-EVENT-003 | event-system | §7 | FC-EVT-003 | Unit | event-system/Impl | Planned |
| FR-SESSION-001 | session-store | §6 | FC-SESS-001 | Unit/Race | session-store/Impl | Planned |
| FR-SESSION-002 | session-store | §6 | FC-SESS-002/003 | Unit | session-store/Impl | Planned |
| FR-SESSION-003 | session-store, runtime-core | §9 | FC-SESS-003, FC-RT-006 | Recovery | session-store/Test | Planned |
| FR-SESSION-004 | session-store | §9 | FC-SESS-004 | Integration | session-store/Impl | Planned |
| FR-CONTEXT-001 | context-manager | §6 | FC-CTX-001 | Unit | context-manager/Impl | Planned |
| FR-CONTEXT-002 | context-manager | §6 | FC-CTX-002 | Unit | context-manager/Impl | Planned |
| FR-CONTEXT-003 | context-manager | §9 | FC-CTX-003 | Golden | context-manager/Impl | Planned |
| FR-CONTEXT-004 | context-manager | §8/§9 | FC-CTX-004, FC-RT-008 | Recovery | context-manager/Impl | Planned |
| FR-CONTEXT-005 | context-manager | §6/§9 | FC-CTX-005 | Golden | context-manager/Impl | Planned |
| FR-CONTEXT-006 | context-manager | §8 | FC-CTX-006 | Failure | context-manager/Impl | Planned |
| FR-MEMORY-001 | memory-system | §6 | FC-MEM-001 | Unit | memory-system/Impl | Planned |
| FR-MEMORY-002 | memory-system | §6/§9 | FC-MEM-002/003 | Integration | memory-system/Impl | Planned |
| FR-MEMORY-003 | memory-system | §14 | FC-MEM-005 | Security | memory-system/Impl | Planned |
| FR-MEMORY-004 | memory-system | §14 | FC-MEM-004 | Security | memory-system/Impl | Planned |
| FR-CMD-001 | extension-system | §6 | FC-CMD-001 | Unit | extension-system/Impl | Planned |
| FR-CMD-002 | extension-system | §6 | FC-CMD-002 | Unit | extension-system/Impl | Planned |
| FR-CMD-100 | extension-system, evaluation | §9 | FC-CMD-100 | Eval | extension-system/Impl | Planned |
| FR-HOOK-001 | extension-system, event-system | §6/§9 | FC-HOOK-001 | Integration | extension-system/Impl | Planned |
| FR-HOOK-002 | extension-system | §6 | FC-HOOK-002 | Unit | extension-system/Impl | Planned |
| FR-HOOK-003 | extension-system | §14 | FC-HOOK-003 | Security | extension-system/Impl | Planned |
| FR-SKILL-001 | extension-system | §6/§8 | FC-SKILL-001 | Unit | extension-system/Impl | Planned |
| FR-SKILL-002 | extension-system | §9 | FC-SKILL-002 | Integration | extension-system/Impl | Planned |
| FR-SKILL-003 | extension-system, permission-engine | §14 | FC-SKILL-003 | Security | extension-system/Impl | Planned |
| FR-SKILL-004 | extension-system, evaluation | §9 | FC-SKILL-004, FC-EVAL-004 | Eval | extension-system/Test | Planned |
| FR-MCP-001 | mcp-client | §8 | FC-MCP-001 | Integration | mcp-client/Impl | Planned |
| FR-MCP-002 | mcp-client | §6 | FC-MCP-002 | Contract | mcp-client/Impl | Planned |
| FR-MCP-003 | mcp-client | §9 | FC-MCP-003 | Contract | mcp-client/Impl | Planned |
| FR-MCP-004 | mcp-client, permission-engine | §14 | FC-MCP-004 | Security | mcp-client/Impl | Planned |
| FR-MCP-005 | mcp-client | §14 | FC-MCP-005 | Security | mcp-client/Impl | Planned |
| FR-SUBAGENT-001 | agent-orchestration | §6 | FC-SUB-001 | Unit | agent-orchestration/Impl | Planned |
| FR-SUBAGENT-002 | agent-orchestration | §9 | FC-SUB-002 | Integration/Race | agent-orchestration/Impl | Planned |
| FR-SUBAGENT-003 | agent-orchestration | §6/§14 | FC-SUB-003 | Contract | agent-orchestration/Impl | Planned |
| FR-SUBAGENT-004 | agent-orchestration | §12 | FC-SUB-004 | Unit | agent-orchestration/Impl | Planned |
| FR-TEAM-001 | agent-orchestration | §6/§8 | FC-TEAM-001 | Integration | agent-orchestration/Impl | Planned |
| FR-TEAM-002 | agent-orchestration | §9 | FC-TEAM-002 | Integration | agent-orchestration/Impl | Planned |
| FR-TEAM-003 | agent-orchestration, telemetry | §10 | FC-TEAM-003 | Integration | agent-orchestration/Impl | Planned |
| FR-WORKTREE-001 | git-worktree | §8/§9 | FC-WT-001/002 | Integration | git-worktree/Impl | Planned |
| FR-WORKTREE-002 | git-worktree | §9 | FC-WT-003 | Failure | git-worktree/Impl | Planned |
| FR-WORKTREE-003 | git-worktree | §9 | FC-WT-003/004 | Recovery | git-worktree/Impl | Planned |
| FR-TELEMETRY-001 | telemetry | §6/§14 | FC-TEL-001 | Security | telemetry/Impl | Planned |
| FR-TELEMETRY-002 | telemetry | §6 | FC-TEL-002 | Unit | telemetry/Impl | Planned |
| FR-TELEMETRY-003 | telemetry, event-system | §9 | FC-TEL-003 | Integration | telemetry/Impl | Planned |
| FR-TELEMETRY-004 | telemetry | §6 | FC-TEL-004 | Unit | telemetry/Impl | Planned |
| FR-EVAL-001 | evaluation, session-store | §6/§9 | FC-EVAL-001 | Integration | evaluation/Impl | Planned |
| FR-EVAL-002 | evaluation | §6 | FC-EVAL-002 | Unit | evaluation/Impl | Planned |
| FR-EVAL-003 | evaluation, extension-system | §10 | FC-EVAL-003 | Eval | evaluation/Impl | Planned |
| FR-EVAL-004 | evaluation | §9 | FC-EVAL-004 | Eval | evaluation/Impl | Planned |
| FR-CLI-001 | cli | §6/§9 | FC-CLI-001/002/003 | Integration | cli/Impl | Planned |
| FR-CLI-002 | cli, extension-system | §9 | FC-CLI-004 | Unit | cli/Impl | Planned |
| FR-CLI-003 | cli | §3/§16 | FC-CLI-004/007 | 依赖检查 | cli/Impl | Planned |
| FR-CLI-004 | cli, session-store | §9 | FC-CLI-006 | Integration | cli/Impl | Planned |

## 非功能需求（Non-Functional）

| Requirement | Module(s) | Spec Section | Task | Test/Eval | Checklist | Status |
| --- | --- | --- | --- | --- | --- | --- |
| NFR-SEC-001 | permission-engine, tool-runtime | §14 | FC-TOOL-002, FC-PERM-* | Security 全套 | MASTER/安全完整性 | Planned |
| NFR-SEC-002 | telemetry | §14 | FC-TEL-001 | Security | telemetry/Impl | Planned |
| NFR-REL-001 | session-store, runtime-core | §11 | FC-SESS-003, FC-RT-006 | Recovery | MASTER/恢复完整性 | Planned |
| NFR-REL-002 | model-provider, runtime-core | §13 | FC-PROV-005, FC-RT-007 | Failure | runtime-core/Impl | Planned |
| NFR-RECOV-001 | session-store, context-manager | §9 | FC-CTX-004, FC-SESS-003 | Recovery | context-manager/Impl | Planned |
| NFR-PERF-001 | tool-runtime | §12 | FC-TOOL-007 | Benchmark | tool-runtime/Test | Planned |
| NFR-PERF-002 | context-manager | §12 | FC-CTX-007 | Benchmark | context-manager/Test | Planned |
| NFR-OBS-001 | telemetry | §15 | FC-TEL-002 | Unit | telemetry/Impl | Planned |
| NFR-TEST-001 | 全体 | §16 | 各模块测试任务 | Race | MASTER/质量门 | Planned |
| NFR-MAINT-001 | 全体 | DEPENDENCY_GRAPH | FC-RT-000 等 | go vet/依赖检查 | MASTER/质量门 | Planned |
| NFR-PORT-001 | sandbox | §14 | FC-SBX-003 | Integration | sandbox/Impl | Planned |
| NFR-COMPAT-001 | session-store, event-system | §11 | FC-SESS-005 | Migration | session-store/Test | Planned |
| NFR-LIMIT-001 | tool-runtime, context-manager, mcp-client | §9 | FC-TOOL-003, FC-CTX-003, FC-MCP-004 | Unit | MASTER/安全 | Planned |
| NFR-COST-001 | context-manager, runtime-core | §6 | FC-CTX-002, FC-RT-004 | Unit | context-manager/Impl | Planned |

## 核心能力覆盖核对（master-plan §2 十三能力）

| 能力 | 主拥有模块 | 状态 |
| --- | --- | --- |
| 2.1 六大工具 | builtin-tools（+tool-runtime 管线） | 覆盖 |
| 2.2 Agent Loop | runtime-core | 覆盖 |
| 2.3 Provider 抽象 | model-provider | 覆盖 |
| 2.4 MCP | mcp-client | 覆盖 |
| 2.5 Skill | extension-system | 覆盖 |
| 2.6 Slash Command | extension-system（+cli 固定命令） | 覆盖 |
| 2.7 Hook | extension-system（+event-system） | 覆盖 |
| 2.8 五层权限 | permission-engine（+sandbox L4） | 覆盖 |
| 2.9 上下文/Token | context-manager | 覆盖 |
| 2.10 记忆 | memory-system | 覆盖 |
| 2.11 SubAgent | agent-orchestration | 覆盖 |
| 2.12 Worktree | git-worktree | 覆盖 |
| 2.13 Agent Team | agent-orchestration | 覆盖 |
| 2.14 标杆 Skill/Eval | evaluation（+extension-system 命令） | 覆盖 |

## 一致性声明
- 无 P0 需求缺失 Task；无 Task 缺失验收标准；无需求缺失负责模块。
- 每个核心实体唯一主拥有模块（见 `DATA_OWNERSHIP.md`）。
- Spec 中重要能力均进入本矩阵。
