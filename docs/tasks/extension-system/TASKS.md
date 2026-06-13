# extension-system Tasks

模块：`extension-system`。Task 前缀：`FC-CMD`（命令）、`FC-HOOK`（Hook）、`FC-SKILL`（Skill）。相关需求 FR-CMD-001/002/100、FR-HOOK-001..003、FR-SKILL-001..004。ADR-0010/0011。RISK-012。

## FC-CMD-001 — Slash Command 框架
| Type | Architecture | Priority | P0 | Milestone | M2 | Status | Ready | Size | M |
| Dependencies | - | Related Requirements | FR-CMD-001 | Spec | §6 |

**Description**：CommandRegistry、Command/CommandKind（Fixed/Prompt/Skill）、Resolve、固定逻辑本地执行。
**Files**：`internal/extension-system/command/`。
**Tests Required**：解析/分派 Unit。
**Acceptance Criteria**：
- [ ] 固定命令本地执行不经模型
- [ ] Prompt/Skill 命令展开为任务输入
**Definition of Done**：与 cli 集成基础命令。

## FC-CMD-002 — 命令级别、参数校验、Alias、冲突、Help
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-CMD-001 | Related Requirements | FR-CMD-002 |

**Description**：内置/用户/项目级命令、ArgSpec 校验、Alias、冲突处理、Help、Command Hook、权限要求。
**Acceptance Criteria**：
- [ ] 参数校验生效
- [ ] 命名冲突返回 ConflictError

## FC-HOOK-001 — HookDispatcher（基于 Event Bus）
| Type | Architecture | Priority | P0 | Milestone | M2 | Status | Backlog | Size | L |
| Dependencies | FC-EVT-003 | Related Requirements | FR-HOOK-001 |

**Description**：实现 eventsystem.Subscriber，订阅生命周期事件，按 Matcher 分发。
**Security Considerations**：不绕过事件总线（ADR-0011）。
**Acceptance Criteria**：
- [ ] 订阅 §2.7 生命周期事件
- [ ] 不存在散落硬编码回调
**Definition of Done**：与 event-system 集成。

## FC-HOOK-002 — Hook 类型与决策结果
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-HOOK-001 | Related Requirements | FR-HOOK-002 |

**Description**：Internal Go/Shell/HTTP Hook；结果 Allow/Deny/Ask/Modify/Continue。
**Security Considerations**：Shell/HTTP 默认禁用。
**Acceptance Criteria**：
- [ ] 三类 Hook 可执行
- [ ] 结果语义正确

## FC-HOOK-003 — Hook 顺序、超时、失败、递归与提权防护
| Type | Security | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-HOOK-002, FC-PERM-001 | Related Requirements | FR-HOOK-003 |

**Description**：优先级/Matcher/Timeout/FailPolicy/决策冲突合并/递归防护/审计；Hook 不可提权。
**Security Considerations**：RISK-012。
**Tests Required**：Security Test（提权尝试、递归、超时）。
**Acceptance Criteria**：
- [ ] Hook 无法扩大权限（提交 permission-engine 仅收窄）
- [ ] 超时按 FailPolicy；递归被防护

## FC-SKILL-001 — Skill 包结构与加载
| Type | Architecture | Priority | P1 | Milestone | M6 | Status | Backlog | Size | L |
| Dependencies | FC-CMD-001 | Related Requirements | FR-SKILL-001 |

**Description**：SKILL.md+manifest.yaml+资源；用户/项目/内置级；SkillManager；延迟加载；Skill State。
**Acceptance Criteria**：
- [ ] manifest 校验
- [ ] 延迟加载注入 ActiveSkill 层

## FC-SKILL-002 — Skill 安装/卸载/升级/Discovery
| Type | Implementation | Priority | P1 | Milestone | M6 | Status | Backlog | Size | M |
| Dependencies | FC-SKILL-001 | Related Requirements | FR-SKILL-002 |

**Description**：安装/卸载/升级/版本锁定/Discovery/显式+自动选择/来源追踪。
**Acceptance Criteria**：
- [ ] 版本锁定生效
- [ ] Discovery 可按任务推荐

## FC-SKILL-003 — Skill 依赖检查与权限边界
| Type | Security | Priority | P1 | Milestone | M6 | Status | Backlog | Size | M |
| Dependencies | FC-SKILL-001, FC-PERM-001 | Related Requirements | FR-SKILL-003 |

**Description**：Tool/MCP/权限依赖检查；Skill 不可扩大自身权限；冲突处理。
**Security Considerations**：越权声明被拒。
**Acceptance Criteria**：
- [ ] 缺依赖拒绝加载
- [ ] 权限声明仅收窄

## FC-SKILL-004 — 候选 Skill 自动生成流水线
| Type | Implementation | Priority | P1 | Milestone | M10 | Status | Backlog | Size | L |
| Dependencies | FC-SKILL-002, FC-EVAL-004 | Related Requirements | FR-SKILL-004 |

**Description**：轨迹→候选→静态检查→回放评测→人工审批→安装；未审批不得 Active。
**Security Considerations**：ADR-0010, RISK-010。
**Acceptance Criteria**：
- [ ] 候选未审批不 Active
- [ ] 静态检查 + 回放评测前置

## FC-CMD-100 — 标杆命令注册（/review-pr /review-sql /review-k8s）
| Type | Implementation | Priority | P1 | Milestone | M6 | Status | Backlog | Size | M |
| Dependencies | FC-CMD-002, FC-SKILL-001 | Related Requirements | FR-CMD-100 |

**Description**：注册三个标杆 Skill 命令，定义输入/输出与所需 Agent（供 evaluation）。
**Acceptance Criteria**：
- [ ] 三命令可触发
- [ ] 与 evaluation 场景对齐

## FC-EXT-900 — extension-system 测试套件
| Type | Test | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-HOOK-003, FC-CMD-002 | Related Requirements | FR-HOOK-001..003, FR-CMD-001/002 |

**Description**：命令/Hook Unit、Hook 安全 Security、与 event-system/permission-engine Integration。
**Acceptance Criteria**：
- [ ] Hook 安全测试通过
- [ ] `go test -race` 通过
