# permission-engine Tasks

模块：`permission-engine`。Task 前缀 `FC-PERM`。相关需求 FR-PERM-001..008。ADR-0005。RISK-006/007/008。

## FC-PERM-001 — Decider 接口与决策模型
| Type | Architecture | Priority | P0 | Milestone | M2 | Status | Done | Size | M |
| Dependencies | FC-TOOL-001 | Related Requirements | FR-PERM-001, FR-PERM-008 | Spec | §6 |

**Description**：定义 `Decider`、`Decision`、`Effect`、`RiskLevel`、`RuleHit`、`Layer`、决策合并优先级。
**Security Considerations**：纯决策，不执行（ADR-0005）。
**Tests Required**：合并优先级 Unit。
**Acceptance Criteria**：
- [x] Deny>AskAlways>AskOnce>Allow 生效
- [x] 接口不含执行能力
**Definition of Done**：接口评审通过。
**Evidence**：实现 `internal/permission-engine` 的 `Decider`、`Decision`、`Effect`、`RiskLevel`、`RuleHit`、`Layer`、`PolicySource`、`BashAnalysis` 与决策合并测试。`go build ./...`、`go test ./...`、`go vet ./...` 通过；race 因缺 `gcc` 未执行。

## FC-PERM-002 — L1 Schema 与输入校验
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-PERM-001 | Related Requirements | FR-PERM-002 |

**Description**：必填、长度、路径格式、空字节、非法编码、JSON Schema、Tool Call 注入风险。
**Tests Required**：Security Test（畸形输入）。
**Acceptance Criteria**：
- [ ] 不合法输入返回 ValidationError
- [ ] 空字节/非法编码被拒

## FC-PERM-003 — L2 资源边界
| Type | Security | Priority | P0 | Milestone | M2 | Status | Backlog | Size | L |
| Dependencies | FC-PERM-001 | Related Requirements | FR-PERM-003 |

**Description**：Workspace Root、读写/敏感/密钥目录、路径穿越、符号链接逃逸、隐藏文件、环境变量。
**Implementation Notes**：规范化 + 解析符号链接真实路径后判断。
**Tests Required**：路径穿越/symlink 逃逸 Security 语料库。
**Security Considerations**：RISK-007。
**Acceptance Criteria**：
- [ ] `..` 逃逸被拒
- [ ] 符号链接指向外部被拒
- [ ] 密钥目录读取被拒/降级
**Definition of Done**：Security 语料库全过。

## FC-PERM-004 — L3 风险策略与决策
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-PERM-001 | Related Requirements | FR-PERM-004 |

**Description**：RiskLevel 评估与 Effect 决策映射、可配置策略。
**Acceptance Criteria**：
- [ ] 四级风险映射到决策
- [ ] 策略可按工具/路径覆盖

## FC-PERM-005 — Bash 结构化分析器
| Type | Security | Priority | P0 | Milestone | M2 | Status | Backlog | Size | L |
| Dependencies | FC-PERM-004 | Related Requirements | FR-PERM-005 |

**Description**：解析程序/参数/管道/重定向/子 Shell/命令替换；识别网络、删除、force-push、Docker/K8s、DB 写、下载后执行、提权。
**Implementation Notes**：非整串匹配（OPEN_QUESTIONS：自研 lexer vs shell 解析库）。
**Tests Required**：Golden（命令→分析）+ Security 语料库。
**Security Considerations**：RISK-006。
**Acceptance Criteria**：
- [ ] 危险模式识别召回达标
- [ ] 绕过尝试（变量拼接/编码）被覆盖或保守降级
**Definition of Done**：语料库通过。

## FC-PERM-006 — L4 沙箱委托挂钩
| Type | Architecture | Priority | P1 | Milestone | M9 | Status | Backlog | Size | S |
| Dependencies | FC-PERM-004 | Related Requirements | FR-PERM-006, FR-SANDBOX-003 |

**Description**：定义委托 sandbox 执行的接口与降级策略（Docker 不可用时拒绝高风险/受限本地）。
**Acceptance Criteria**：
- [ ] 高风险可路由到沙箱
- [ ] 沙箱不可用按策略降级

## FC-PERM-007 — L5 审批请求与审计记录
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-PERM-001, FC-EVT-002 | Related Requirements | FR-PERM-007 |

**Description**：AskOnce/AskAlways 产生 ApprovalRequested；记录发起者/Session/Agent/Tool/原始参数（脱敏）/风险/命中原因/结果。
**Security Considerations**：审批绕过防护 RISK-008。
**Acceptance Criteria**：
- [ ] 审批事件含全部审计字段
- [ ] 决策与执行分离可验证

## FC-PERM-008 — 决策冲突与优先级测试
| Type | Test | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-PERM-002, FC-PERM-003, FC-PERM-004, FC-PERM-005, FC-PERM-007 | Related Requirements | FR-PERM-008 |

**Description**：多来源（默认/Skill/Hook/用户）冲突合并、Skill/Hook 无法扩权的测试。
**Acceptance Criteria**：
- [ ] Skill/Hook 仅能收窄
- [ ] 最严格生效
- [ ] `go test -race` 通过
