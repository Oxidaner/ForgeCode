# permission-engine Checklist

模块：`permission-engine`。相关需求 FR-PERM-001..008。安全关键模块（风险等级 Critical）。

## Design Ready
- [x] `Decider`/`Decision`/`Effect`/`RiskLevel` 已定义（与 GLOSSARY 一致，FC-PERM-001）
  Evidence: `internal/permission-engine/types.go`
- [x] 五层职责边界与短路顺序已定义
  Evidence: `internal/permission-engine/policy_decider.go`
- [x] 决策与执行分离明确（ADR-0005，不含执行能力）
  Evidence: `internal/permission-engine/decision_test.go`
- [x] 决策合并优先级规则已定义（Deny 优先、最严格生效）
  Evidence: `internal/permission-engine/decision.go`
- [x] Approval 契约字段已定义，存储归 session-store
  Evidence: `ApprovalRequest` 仅为决策契约，持久化仍由后续 session-store/event-system 接入
- [x] 错误模型映射 GLOSSARY（ValidationError/PermissionDenied/ApprovalRequired）
  Evidence: `internal/tool-runtime/types.go`

## Implementation Ready
- [x] 任务已拆分（L1–L5 + Bash 分析 + 冲突测试）
  Evidence: `docs/tasks/permission-engine/TASKS.md`
- [x] 敏感目录默认清单已定义（跨平台考虑）
  Evidence: `internal/permission-engine/config.go`
- [x] Bash 分析实现路径已决定（OPEN_QUESTIONS Q16）
  Evidence: `docs/planning/OPEN_QUESTIONS.md`, `internal/permission-engine/bash_analyzer.go`
- [x] Security 测试语料库结构已定义
  Evidence: `internal/permission-engine/policy_decider_test.go`

## Implementation Complete
- [x] L1 拒绝畸形/注入输入（FR-PERM-002）
  Evidence: `internal/permission-engine/schema.go`, `policy_decider_test.go`
- [x] L2 拦截路径穿越与符号链接逃逸（FR-PERM-003, RISK-007）
  Evidence: `internal/permission-engine/resource.go`, `policy_decider_test.go`
- [x] L3 四级风险映射决策（FR-PERM-004）
  Evidence: `internal/permission-engine/policy_decider.go`, `policy_decider_test.go`
- [x] L3 Bash 结构化分析识别危险操作（FR-PERM-005, RISK-006）
  Evidence: `internal/permission-engine/bash_analyzer.go`, `bash_analyzer_test.go`
- [x] L5 审批请求含全部审计字段（FR-PERM-007）
  Evidence: `internal/permission-engine/approval.go`, `policy_decider_test.go`
- [x] 引擎不执行任何外部操作
  Evidence: `Decider` 仅返回 `Decision`，无命令/工具执行调用
- [x] Skill/Hook 无法扩权（FR-PERM-008, RISK-008）
  Evidence: `policy_decider_test.go`

## Test Complete
- [x] 各层规则 Unit Test
  Evidence: `go test ./...`
- [x] 路径穿越/symlink Security 语料库通过
  Evidence: `internal/permission-engine/policy_decider_test.go`
- [x] Bash 分析 Golden Test
  Evidence: `internal/permission-engine/bash_analyzer_test.go`
- [x] Bash 危险命令 Security 语料库通过
  Evidence: `internal/permission-engine/bash_analyzer_test.go`
- [x] 决策合并优先级 Unit Test
  Evidence: `go test ./...`
- [x] 与 tool-runtime Contract Test
  Evidence: `PolicyDecider.DecideTool`, `tool-runtime/integration_test.go`
- [x] `go test -race` 通过
  Evidence: `go test -race ./...`

## Documentation Complete
- [ ] SPEC 五层与实现一致
- [x] TASKS 状态更新
  Evidence: `docs/tasks/permission-engine/TASKS.md` FC-PERM-001/002/003/004/005/007/008 Done
- [ ] ADR-0005 与实现一致
- [ ] 默认敏感目录/策略配置示例

## Release Ready
- [ ] 所有 P0 验收通过
- [ ] 无未处理 Critical 安全风险（RISK-006/007/008 缓解到位）
- [ ] 审批与决策审计可观察
- [ ] 危险 Bash 审批流 Demo 可复现
