# builtin-tools Checklist

模块：`builtin-tools`。相关需求 FR-TOOL-101..106。检查项针对六个内置工具的安全与正确性，可验证、非任务标题复制。

## Design Ready
- [x] 六个工具均实现 tool-runtime 的 `Tool` 接口并经 Registry 注册（Related Task: FC-BT-001）
  Evidence: `internal/builtin-tools/registry.go`, `registry_test.go`
- [x] 每个工具的 `ToolDescriptor`（schema/风险等级/权限要求）已定义
  Evidence: `internal/builtin-tools/descriptors.go`
- [x] 工具不自行做权限决策（委托 permission-engine），不自行写 Event Store/Approval
  Evidence: 内置工具仅实现 `Execute`，权限/audit 经 `tool-runtime.Invoker` 统一处理
- [x] 错误归入 GLOSSARY 错误分类（ToolExecutionError/ValidationError/TimeoutError）
  Evidence: `internal/builtin-tools/registry.go`, `internal/tool-runtime/types.go`
- [x] 写类工具的 Checkpoint 依赖 session-store 接口已约定
  Evidence: `internal/builtin-tools/deps.go`

## Implementation Ready
- [x] 任务已拆分（读类/写编辑类/Bash/搜索类）
  Evidence: `docs/tasks/builtin-tools/TASKS.md`
- [x] P0 依赖（FC-TOOL-001/002 管线、FC-PERM L1–L3）已满足
  Evidence: `tool-runtime/integration_test.go`, `permission-engine` tests
- [x] 各工具配置项（分页大小、超时、输出上限、二进制阈值）已定义默认值
  Evidence: `internal/builtin-tools/deps.go`
- [x] Golden 测试语料与 Fake FS/Clock 边界已定义
  Evidence: `*_test.go`, `fakeCheckpointer`

## Implementation Complete
- [x] ReadFile：分页、二进制识别、超大文件保护生效（FR-TOOL-101）
  Evidence: `internal/builtin-tools/readfile.go`, `readfile_test.go`
- [x] WriteFile：写前 Checkpoint；失败保持原文件不变（FR-TOOL-102）
  Evidence: `writefile.go`, `atomic_write.go`, `write_edit_test.go`
- [x] EditFile：唯一匹配校验，非唯一/未命中报错，产出 Diff（FR-TOOL-103）
  Evidence: `editfile.go`, `write_edit_test.go`
- [x] Bash：超时终止、输出头尾保留、退出码与错误分类（FR-TOOL-104）
  Evidence: `internal/builtin-tools/bash.go`, `bash_test.go`
- [x] Glob：模式匹配，结果受 Workspace 边界约束（FR-TOOL-105）
  Evidence: `internal/builtin-tools/glob.go`, `glob_grep_test.go`
- [x] Grep：正则搜索，结果去重（FR-TOOL-106）
  Evidence: `internal/builtin-tools/grep.go`, `glob_grep_test.go`
- [x] 所有工具响应 Context Cancellation
  Evidence: 各工具 `Execute` 读取 `context.Context`；Bash/Grep/Glob/ReadFile 有取消路径，Write/Edit 写入前检查 ctx
- [x] 工具输出在返回前可被 tool-runtime 硬截断（不内嵌截断策略冲突）
  Evidence: `tool-runtime.PipelineInvoker.truncateResult`

## Test Complete
- [x] ReadFile / Bash / Glob / Grep Unit Test
  Evidence: `go test ./...`
- [x] WriteFile/Bash 的 Failure Injection（Checkpoint 失败、超时、被取消）
  Evidence: `write_edit_test.go`, `bash_test.go`, `tool-runtime/invoker_test.go`
- [x] Bash 超时 Failure Test
  Evidence: `internal/builtin-tools/bash_test.go`
- [x] EditFile Golden Test
  Evidence: `write_edit_test.go`
- [x] Grep/ReadFile 输出测试
  Evidence: `internal/builtin-tools/glob_grep_test.go`, `readfile_test.go`
- [x] Bash 危险命令经 permission-engine 被拦截的 Security/Integration Test（RISK-006）
  Evidence: `permission-engine/bash_analyzer_test.go`, `tool-runtime/integration_test.go`
- [x] 路径越界读写被拒的 Security Test（配合 FR-PERM-003）
  Evidence: `permission-engine/policy_decider_test.go`, `tool-runtime/integration_test.go`
- [x] `go test -race` 通过
  Evidence: `go test -race ./...`

## Documentation Complete
- [x] SPEC 与工具 Descriptor 一致
  Evidence: `docs/specs/builtin-tools/SPEC.md`, `internal/builtin-tools/descriptors.go`
- [x] TASKS 状态更新
  Evidence: `docs/tasks/builtin-tools/TASKS.md` FC-BT-001/002/003/004/005/006/007 Done
- [ ] 配置示例（分页/超时/上限）更新
- [ ] 已知限制（二进制阈值、换行归一化）记录

## Release Ready
- [ ] 六工具 P0 验收通过
- [ ] 无未处理 Critical 安全风险
- [ ] 工具调用计入审计与 Usage
- [ ] Demo 中只读+修改流程可复现
