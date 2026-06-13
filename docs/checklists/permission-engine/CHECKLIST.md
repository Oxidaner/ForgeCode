# permission-engine Checklist

模块：`permission-engine`。相关需求 FR-PERM-001..008。安全关键模块（风险等级 Critical）。

## Design Ready
- [ ] `Decider`/`Decision`/`Effect`/`RiskLevel` 已定义（与 GLOSSARY 一致，FC-PERM-001）
- [ ] 五层职责边界与短路顺序已定义
- [ ] 决策与执行分离明确（ADR-0005，不含执行能力）
- [ ] 决策合并优先级规则已定义（Deny 优先、最严格生效）
- [ ] Approval 契约字段已定义，存储归 session-store
- [ ] 错误模型映射 GLOSSARY（ValidationError/PermissionDenied/ApprovalRequired）

## Implementation Ready
- [ ] 任务已拆分（L1–L5 + Bash 分析 + 冲突测试）
- [ ] 敏感目录默认清单已定义（跨平台考虑）
- [ ] Bash 分析实现路径已决定（OPEN_QUESTIONS）
- [ ] Security 测试语料库结构已定义

## Implementation Complete
- [ ] L1 拒绝畸形/注入输入（FR-PERM-002）
- [ ] L2 拦截路径穿越与符号链接逃逸（FR-PERM-003, RISK-007）
- [ ] L3 四级风险映射决策（FR-PERM-004）
- [ ] L4 Bash 结构化分析识别危险操作（FR-PERM-005, RISK-006）
- [ ] L5 审批事件含全部审计字段（FR-PERM-007）
- [ ] 引擎不执行任何外部操作
- [ ] Skill/Hook 无法扩权（FR-PERM-008, RISK-008）

## Test Complete
- [ ] 各层规则 Unit Test
- [ ] 路径穿越/symlink/命令注入/提权 Security 语料库通过
- [ ] Bash 分析 Golden Test
- [ ] 决策合并优先级 Unit Test
- [ ] 与 tool-runtime Contract Test
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 五层与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0005 与实现一致
- [ ] 默认敏感目录/策略配置示例

## Release Ready
- [ ] 所有 P0 验收通过
- [ ] 无未处理 Critical 安全风险（RISK-006/007/008 缓解到位）
- [ ] 审批与决策审计可观察
- [ ] 危险 Bash 审批流 Demo 可复现
