# ForgeCode 文档导航

本目录是 ForgeCode 的架构与规划文档集。当前已从 **架构规划阶段** 进入首批 Ready P0 Task 的实现阶段；实现范围以对应 Task 和 Spec 为准。

## 阅读顺序
1. [master-plan.md](master-plan.md) — 原始规划指令（输入）。
2. [specs/00-master/SPEC.md](specs/00-master/SPEC.md) — 总体 Spec（目标、需求、流程、数据、版本）。
3. [architecture/SYSTEM_OVERVIEW.md](architecture/SYSTEM_OVERVIEW.md) — 系统总览。
4. [architecture/MODULE_MAP.md](architecture/MODULE_MAP.md) — 权威模块清单。
5. [architecture/DEPENDENCY_GRAPH.md](architecture/DEPENDENCY_GRAPH.md) — 依赖图与分层。
6. 模块 Spec：[specs/&lt;module-id&gt;/SPEC.md](specs/)。
7. 任务与清单：[tasks/MASTER_TASKS.md](tasks/MASTER_TASKS.md)、[checklists/MASTER_CHECKLIST.md](checklists/MASTER_CHECKLIST.md)。

## 目录结构

```text
docs/
├── README.md                  本导航
├── master-plan.md             原始规划指令
├── architecture/              SYSTEM_OVERVIEW / MODULE_MAP / DEPENDENCY_GRAPH /
│                              EVENT_MODEL / DATA_OWNERSHIP / SECURITY_MODEL / FAILURE_AND_RECOVERY
├── specs/                     00-master 总 Spec + 每模块 SPEC.md
├── tasks/                     MASTER_TASKS + 每模块 TASKS.md
├── checklists/                MASTER_CHECKLIST + 每模块 CHECKLIST.md
├── adr/                       ADR-0001 … ADR-0012 + README
├── planning/                  ROADMAP / MILESTONES / TRACEABILITY_MATRIX /
│                              RISK_REGISTER / OPEN_QUESTIONS / GLOSSARY
└── templates/                 SPEC / TASK / CHECKLIST / ADR 模板
```

## 权威来源（Single Source of Truth）
- **模块命名**：`architecture/MODULE_MAP.md`
- **事件格式**：`architecture/EVENT_MODEL.md`
- **数据所有权**：`architecture/DATA_OWNERSHIP.md`
- **状态/错误/枚举术语**：`planning/GLOSSARY.md`
- **需求 ID**：`specs/00-master/SPEC.md` §6.4/§6.5
- **任务 ID 与 DAG**：`tasks/MASTER_TASKS.md`
- **追踪矩阵**：`planning/TRACEABILITY_MATRIX.md`

修改任一权威来源时，须同步相关文档（见根目录 `AGENTS.md`）。
