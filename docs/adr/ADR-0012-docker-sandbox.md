# ADR-0012：第一版 Sandbox 使用 Docker

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-SANDBOX-001, FR-SANDBOX-002, FR-SANDBOX-003, FR-PERM-006 |
| Related Modules | sandbox, permission-engine |

## Context
高风险命令需运行时隔离。自研容器运行时成本过高且非项目重点（明确非目标）。

## Decision
第一版 Sandbox **基于 Docker**：工作目录挂载、只读挂载、网络控制、CPU/内存/PID/执行时间限制、环境变量过滤、进程回收。Permission Engine L4 通过接口委托 sandbox 执行，不自研容器运行时。MVP 默认关闭（本地受限执行），V0.2 引入可选开启。Docker 不可用时按 FR-SANDBOX-003 降级（拒绝高风险或本地受限执行）。

## Alternatives Considered
- **自研容器运行时**：成本高、非重点——拒绝（非目标）。
- **无沙箱**：高风险命令无隔离——仅 MVP 临时，靠权限+审批兜底。
- **gVisor/Firecracker**：更强隔离但更重——后续评估。

## Consequences
- 正面：复用成熟隔离、资源限制清晰。
- 负面：依赖本地 Docker（NFR-PORT-001 降级处理）；Sandbox Escape 风险仍由 L1–L3 兜底。

## Security Impact
提供运行时隔离层；逃逸假设下仍受输入/资源/风险三层约束。

## Operational Impact
需检测 Docker 可用性与镜像管理；缺失时降级且告警。

## Revisit Conditions
若需更强隔离或无 Docker 环境占比高，评估替代沙箱技术。
