# builtin-tools Checklist

模块：`builtin-tools`。相关需求 FR-TOOL-101..106。检查项针对六个内置工具的安全与正确性，可验证、非任务标题复制。

## Design Ready
- [ ] 六个工具均实现 tool-runtime 的 `Tool` 接口并经 Registry 注册（Related Task: FC-BT-001）
- [ ] 每个工具的 `ToolDescriptor`（schema/风险等级/权限要求）已定义
- [ ] 工具不自行做权限决策（委托 permission-engine），不自行写 Event Store/Approval
- [ ] 错误归入 GLOSSARY 错误分类（ToolExecutionError/ValidationError/TimeoutError）
- [ ] 写类工具的 Checkpoint 依赖 session-store 接口已约定

## Implementation Ready
- [ ] 任务已拆分（读类/写编辑类/Bash/搜索类）
- [ ] P0 依赖（FC-TOOL-001/002 管线、FC-PERM L1–L3）已满足
- [ ] 各工具配置项（分页大小、超时、输出上限、二进制阈值）已定义默认值
- [ ] Golden 测试语料与 Fake FS/Clock 边界已定义

## Implementation Complete
- [ ] ReadFile：分页、二进制识别、超大文件保护生效（FR-TOOL-101）
- [ ] WriteFile：写前 Checkpoint；失败保持原文件不变（FR-TOOL-102）
- [ ] EditFile：唯一匹配校验，非唯一/未命中报错，产出 Diff（FR-TOOL-103）
- [ ] Bash：超时终止、输出头尾保留、退出码与错误分类（FR-TOOL-104）
- [ ] Glob：模式匹配，结果受 Workspace 边界约束（FR-TOOL-105）
- [ ] Grep：正则搜索，结果去重（FR-TOOL-106）
- [ ] 所有工具响应 Context Cancellation
- [ ] 工具输出在返回前可被 context-manager 截断（不内嵌截断策略冲突）

## Test Complete
- [ ] 各工具 Unit Test
- [ ] WriteFile/Bash 的 Failure Injection（写失败、超时、被取消）
- [ ] EditFile/Grep/ReadFile Golden Test
- [ ] Bash 危险命令经 permission-engine 被拦截的 Security/Integration Test（RISK-006）
- [ ] 路径越界读写被拒的 Security Test（配合 FR-PERM-003）
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与工具 Descriptor 一致
- [ ] TASKS 状态更新
- [ ] 配置示例（分页/超时/上限）更新
- [ ] 已知限制（二进制阈值、换行归一化）记录

## Release Ready
- [ ] 六工具 P0 验收通过
- [ ] 无未处理 Critical 安全风险
- [ ] 工具调用计入审计与 Usage
- [ ] Demo 中只读+修改流程可复现
