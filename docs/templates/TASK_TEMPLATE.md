# Task 模板

> 复制本模板填写 `docs/tasks/<module-id>/TASKS.md` 中的单个 Task 条目。

## Task: FC-XXX-000 — <Title>

| 字段 | 值 |
| --- | --- |
| ID | FC-XXX-000 |
| Title | <Title> |
| Module | `<module-id>` |
| Type | Spike / Architecture / Implementation / Test / Security / Documentation / Evaluation / Migration / Refactor |
| Priority | P0 / P1 / P2 |
| Milestone | M1..M12 |
| Status | Backlog / Ready / In Progress / Blocked / Done |
| Size | XS / S / M / L / XL（XL 必须继续拆分） |
| Dependencies | FC-XXX-000, ... |
| Related Requirements | FR-XXX-001, ... |
| Related Spec Sections | `<module-id>` §6, §8 |

**Description**：要做什么，输入与输出是什么。

**Implementation Notes**：关键实现要点、约束、注意事项。

**Files or Packages Likely Affected**：`internal/...`（计划路径，非强制）。

**Tests Required**：需要的测试类型与关键用例。

**Security Considerations**：是否触及权限/审计/敏感数据。

**Acceptance Criteria**：
- [ ] 可测试条件 1
- [ ] 可测试条件 2

**Definition of Done**：代码 + 测试 + 文档 + Checklist + Evidence 的完成标准。

**Evidence**：完成后填写（测试输出、PR、commit、Eval 报告链接）。

---

> Spike 类型额外要求：需回答的问题、最小实验、输出决策、结束条件。
