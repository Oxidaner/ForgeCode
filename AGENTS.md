# AGENTS.md

ForgeCode 仓库的长期有效规则。本文件只放稳定规则，不复制 Spec 内容。完整设计见 `docs/`。

## 项目定位
ForgeCode：用 Go 自主实现的、模型无关的 Coding Agent Runtime 控制平面。**不是封装 Agent SDK**。详见 `docs/specs/00-master/SPEC.md`。

## 语言与版本
- 主语言：Go（目标版本见 `docs/planning/OPEN_QUESTIONS.md` Q1，待 Spike 确认；规划假设 Go 1.22+）。
- 文档：简体中文；代码标识符、ID、状态枚举、技术术语用英文。

## 目录约定
- 设计文档在 `docs/`（结构见 `docs/README.md`）。
- 计划实现代码（待启动）：`cmd/`（入口）、`internal/<module-id>/`（按模块）、`pkg/`（可复用公开包，谨慎使用）。
- 不创建 `common`/`utils`/`manager` 垃圾桶包。

## 架构依赖规则
- 依赖方向：CLI → Runtime Coordinator → Core Domain Interfaces → Infrastructure。见 `docs/architecture/DEPENDENCY_GRAPH.md`。
- 核心逻辑（Agent Loop、状态机、工具/权限管线、事件模型、上下文压缩）**不得依赖任何 Agent SDK**。
- `runtime-core` **不得依赖具体 Provider**；Provider 私有结构只在 `model-provider` 适配器内部。
- 禁止循环依赖（`go vet` + 评审保证）。

## 安全规则（不可绕过）
- 所有工具调用必经统一管线：**Validation → Permission → Hook → Execution → Audit**。
- 禁止绕过 Permission Engine。
- 禁止 MCP Tool 绕过统一 Tool Runtime 与权限。
- 禁止 SubAgent / Team Member / Slash Command / Skill 绕过权限。
- Skill/Hook 不得扩大自身权限；Hook 不得静默提权。
- 敏感数据（密钥、Token、完整环境变量）不得写入普通日志（telemetry 强制脱敏）。
- 高风险操作未经审批不得自动执行。

## 命令（待建立）
当前为规划阶段，尚无构建/测试代码。以下命令在代码骨架建立后填入，**现在不要虚构**：
- 构建：待建立。
- 测试：计划 `go test ./...`（含 `-race`），待项目初始化后确认。
- 格式化：计划 `gofmt` / `goimports`，待确认。
- 静态检查：计划 `go vet` + linter（如 `golangci-lint`），待确认。

## ID 规则
- Module ID：见 `docs/architecture/MODULE_MAP.md`。
- Requirement ID：`FR-<DOMAIN>-NNN` / `NFR-<CAT>-NNN`。
- Task ID：`FC-<AREA>-NNN`。
- ADR：`ADR-NNNN`；Risk：`RISK-NNN`；Open Question：`Q<N>`。

## 文档同步规则
- 修改 Spec/接口/事件/状态枚举时，同步更新 `MODULE_MAP`、`EVENT_MODEL`、`DATA_OWNERSHIP`、`GLOSSARY` 与 `TRACEABILITY_MATRIX`。
- 完成 Task 时更新对应模块 `TASKS.md` 状态、`CHECKLIST.md` 勾选项与 Evidence。
- 新增重要能力须进入 Traceability Matrix（不得有无 Task 的 P0 需求、无验收的 Task）。
- 状态机/错误分类以 `docs/planning/GLOSSARY.md` 为准。

## 当前阶段约束
- 规划阶段仅修改 Markdown/YAML/JSON 设计文档。
- 不实现生产代码、不安装依赖、不创建大量空骨架文件，直到架构文档就绪并选定首个 Ready 的 P0 Task。
