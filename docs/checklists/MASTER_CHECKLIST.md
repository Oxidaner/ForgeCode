# MASTER_CHECKLIST

跨模块总验收清单。模块级清单见 `docs/checklists/<module-id>/CHECKLIST.md`。本清单聚焦全局一致性、安全完整性、恢复完整性与里程碑就绪。

## 架构就绪（规划阶段）
- [ ] 17 个模块均有 SPEC/TASKS/CHECKLIST（Evidence: `docs/specs|tasks|checklists/*`）
- [ ] 模块命名、事件格式、数据所有权、术语枚举单一真相源已建立
- [ ] 依赖图无循环（Related: `DEPENDENCY_GRAPH.md`）
- [ ] 十三项核心能力均有负责模块（Related: `TRACEABILITY_MATRIX.md`）
- [ ] 所有 P0 需求有 模块/Spec/Task/Test/Checklist（Related: 追踪矩阵）
- [ ] ADR-0001..0012 已完成
- [ ] 风险登记与 Open Questions 已建立

## 安全完整性（全局，NFR-SEC-001）
- [ ] 所有工具调用路径经 Validation→Permission→Hook→Execution→Audit
  - Related Task: FC-TOOL-002
- [ ] MCP Tool 不绕过权限（FC-MCP-004）
- [ ] SubAgent/Team Member 不绕过权限（FC-SUB-002, FC-TEAM-002）
- [ ] Slash Command 不绕过权限（FC-CMD-002）
- [ ] Skill 不可扩大自身权限（FC-SKILL-003）
- [ ] Hook 不可静默提升权限（FC-HOOK-003）
- [ ] 敏感数据不入普通日志（FC-TEL-001, NFR-SEC-002）
- [ ] 路径穿越/符号链接逃逸被拦截（FC-PERM-003）
- [ ] Bash 危险命令结构化识别（FC-PERM-005）
- [ ] 记忆污染控制（FC-MEM-004）
- [ ] 候选 Skill 未审批不生效（FC-SKILL-004, ADR-0010）

## 恢复完整性（NFR-REL-001）
- [ ] 关键状态可由持久化事件/Checkpoint 恢复（FC-SESS-003, FC-RT-006）
- [ ] 恢复不重执行外部副作用
- [ ] 审批后执行前崩溃恢复不自动执行高危操作
- [ ] 压缩失败可回滚 Checkpoint（FC-CTX-006）
- [ ] 崩溃孤儿 Worktree 可清理（FC-WT-003）

## 状态一致性
- [ ] Agent/Session/Task/Team/Worktree/MCP/Skill 状态定义与 GLOSSARY 一致
- [ ] 各模块 SPEC 状态机引用权威枚举

## MVP 可完成性（M1–M5）
- [ ] MVP 模块集（10+扩展 Command/Hook）形成单 Agent 闭环
- [ ] 不依赖 Agent Teams 完整实现即可演示
- [ ] 六条核心流程（只读/修改/审批/压缩/恢复/取消）可复现
- [ ] MVP 各模块 Release Ready 清单通过

## 里程碑就绪门（每个 Milestone 退出条件）
- [ ] M1：单 Agent 只读任务可运行且可恢复
- [ ] M2：编辑 + 权限 + 审批 + Hook 闭环
- [ ] M3：上下文压缩与预算
- [ ] M4：暂停/恢复/取消/崩溃恢复
- [ ] M5：MVP Demo 可复现
- [ ] M6–M13：见 `MILESTONES.md` 各 DoD

## 质量门（实现阶段统一）
- [ ] `go build ./...` 通过
- [ ] `go test ./... -race` 通过
- [ ] `go vet` / lint 通过，无循环依赖
- [ ] 安全测试套件通过
- [ ] 关键恢复路径测试通过
- [ ] 文档与 Traceability Matrix 同步（Evidence 已记录）
