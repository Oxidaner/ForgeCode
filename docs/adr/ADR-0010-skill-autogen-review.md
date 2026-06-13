# ADR-0010：Skill 自动生成需评测与人工审批

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-SKILL-004, FR-SKILL-003, FR-MEMORY-004 |
| Related Modules | extension-system, evaluation |

## Context
从成功任务轨迹自动沉淀 Skill 有价值，但自动生成的 Skill 可能携带错误、过拟合或恶意内容，若直接生效将污染系统。

## Decision
自动生成 Skill 走固定流水线：**轨迹 → 提取可复用步骤 → 候选 Skill → 静态检查 → 回放评测 → 人工审批 → 安装到 Registry**。候选 Skill 在审批前为 `Candidate/StaticChecked/Replayed/Approved`，**未审批不得 Active**。Skill 权限声明只能收窄、不能自我扩权。

## Alternatives Considered
- **自动生成直接生效**：高风险污染与安全问题——拒绝。
- **完全禁止自动生成**：丧失经验沉淀能力——拒绝。

## Consequences
- 正面：经验可沉淀且可控；防止污染与提权。
- 负面：流水线与 Eval 成本；需人工环节。

## Security Impact
阻断恶意/错误 Skill 自动生效；与 Malicious Skill 威胁对应。

## Operational Impact
依赖 evaluation 的回放评测能力（FR-EVAL-004）。

## Revisit Conditions
当回放评测足够可信时，评估对低风险 Skill 的自动审批门槛。
