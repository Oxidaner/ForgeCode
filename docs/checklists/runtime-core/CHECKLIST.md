# runtime-core Checklist

模块：`runtime-core`。相关需求 FR-RUNTIME-001..006。

## Design Ready
- [x] Agent State 枚举与合法转移表已定义（与 GLOSSARY 一致，FC-RT-001；Evidence: `go test ./internal/runtime-core`）
- [ ] Coordinator 仅依赖接口，无具体 Provider/Tool 类型（ADR-0003）
- [ ] 错误模型映射到 GLOSSARY 分类
- [ ] 恢复语义（事件重放、不重放副作用）已定义
- [ ] 强制工具调用经统一管线的约束已写入设计
- [ ] 依赖方向无环（DEPENDENCY_GRAPH 一致）

## Implementation Ready
- [ ] 任务已拆分（状态机/Coordinator/Loop/预算/循环/恢复）
- [x] Go 版本基线已确定（FC-RT-000；Evidence: `go.mod` Go 1.22, Spike 环境 Go 1.26.1, `go test ./internal/runtime-core`）
- [ ] 预算/上限默认值已定义
- [ ] Mock Provider 与 Fake Store 边界已定义

## Implementation Complete
- [x] 状态机非法转移被拒并记录（FR-RUNTIME-001；Evidence: `internal/runtime-core/state_test.go`）
- [ ] 五类上限均可触发安全终止（FR-RUNTIME-002）
- [ ] 取消传播到正在执行的工具（FR-RUNTIME-003）
- [ ] 非法 Tool Call 作为 Observation 反馈不直接执行（FR-RUNTIME-004）
- [ ] 重复工具调用/相同错误循环被检测（RISK-003）
- [x] AgentStateChanged 等事件完整记录（FC-RT-001；Evidence: `AgentStateChangedPayload` 测试）
- [ ] 完整提示/密钥不入普通日志

## Test Complete
- [x] 状态转移 Golden Test（Evidence: `go test ./internal/runtime-core`）
- [ ] 预算上限 Unit Test
- [ ] 循环检测 Failure Injection
- [ ] 取消传播 Integration
- [ ] 崩溃后 Resume Recovery Test（NFR-REL-001）
- [ ] `go test -race` 通过

## Documentation Complete
- [x] SPEC 状态机与实现一致
- [x] TASKS 状态更新
- [ ] ADR-0001/0002 与实现一致
- [ ] 配置示例更新

## Release Ready
- [ ] 六条核心流程（§6.7 1–5,12）可演示
- [ ] 无未处理 Critical 风险
- [ ] 轮次/工具/Token/Cost 指标可观察
- [ ] 审批后崩溃恢复不自动执行高危操作（Evidence）
- [ ] Demo 可复现
