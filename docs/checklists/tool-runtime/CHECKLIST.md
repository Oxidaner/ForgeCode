# tool-runtime Checklist

模块：`tool-runtime`。相关需求 FR-TOOL-001..004。

## Design Ready
- [x] `Tool`/`ToolDescriptor`/`Registry`/`Invoker` 接口已定义（FC-TOOL-001）
  Evidence: `internal/tool-runtime/tool.go`, `types.go`, `registry.go`
- [ ] 统一管线阶段顺序已定义且不可跳过（ADR-0004/0005）
- [x] ToolCall/ToolResult 契约已定义，存储归 session-store
  Evidence: `internal/tool-runtime/types.go`
- [x] 错误分类映射 GLOSSARY
  Evidence: `internal/tool-runtime/types.go`
- [x] Namespace 与命名冲突策略已定义
  Evidence: `internal/tool-runtime/registry_test.go`
- [ ] 依赖方向无环

## Implementation Ready
- [ ] 任务已拆分（接口/管线/截断/审批/审计/MCP/契约）
- [ ] permission-engine 与 event-system 接口已约定
- [ ] 截断与超时默认值已定义
- [ ] Fake Tool/Fake Decider 测试边界已定义

## Implementation Complete
- [ ] 任何 Invoke 经 Validation→Permission→Hook→Execute→Audit（FR-TOOL-002）
- [x] 命名冲突返回 ConflictError
  Evidence: `internal/tool-runtime/registry_test.go`
- [ ] 输出硬截断并标注 Truncated（FR-TOOL-003, NFR-LIMIT-001）
- [ ] 超时/取消正确终止并归类
- [ ] ApprovalRequired 上抛、批准后续行（FC-TOOL-004）
- [ ] 完整审计事件产生，输入脱敏

## Test Complete
- [ ] 管线顺序 Integration
- [ ] 绕过尝试 Security Test（无旁路）
- [ ] 截断/超时/取消 Failure Injection
- [ ] 内置与 MCP 工具通过同套 Contract Test（RISK-005）
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 接口与实现一致
- [x] TASKS 状态更新
  Evidence: `docs/tasks/tool-runtime/TASKS.md` FC-TOOL-001 Done
- [ ] 配置示例更新
- [ ] 截断/超时已知限制记录

## Release Ready
- [ ] P0 验收通过
- [ ] 无未处理 Critical 安全风险（无权限旁路）
- [ ] 工具调用指标与审计可观察
- [ ] Demo 中工具调用链可复现
