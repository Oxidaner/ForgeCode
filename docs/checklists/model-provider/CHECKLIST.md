# model-provider Checklist

模块：`model-provider`。相关需求 FR-PROVIDER-001..006。ADR-0003。

## Design Ready
- [ ] 中立 `Provider` 接口已定义，不含 Provider 私有字段（FC-PROV-001）
- [ ] 中立 Message/ToolCall/Usage/StopReason 模型已定义
- [ ] runtime-core 仅依赖接口（依赖反转，DEPENDENCY_GRAPH 一致）
- [ ] 错误归一为 `ProviderError`（含 RateLimit 子类），分类符合 GLOSSARY
- [ ] 能力差异矩阵（RISK-004）记录策略已定义

## Implementation Ready
- [ ] 任务已拆分（接口/Mock/流式/工具调用/错误/适配器/契约）
- [ ] Mock Provider 边界已定义（FC-PROV-002）
- [ ] 重试/超时/退避默认参数已定义
- [ ] Contract Test 套件结构已定义（FC-PROV-009）

## Implementation Complete
- [ ] 普通响应与 Streaming 行为一致（FR-PROVIDER-001）
- [ ] 多 Tool Call 顺序与 ID 保留（FR-PROVIDER-002）
- [ ] Stop Reason 与 Token Usage 正确解析
- [ ] 瞬时错误退避重试、不可重试快速失败（FR-PROVIDER-004）
- [ ] provider-specific 结构未泄漏到 runtime-core（FR-PROVIDER-006）
- [ ] Context Cancellation 中断请求/流

## Test Complete
- [ ] Mock + OpenAI 通过统一 Contract Test
- [ ] 流式中断 Failure Injection
- [ ] 限流/超时/5xx Failure Injection
- [ ] `go test -race` 通过
- [ ] Anthropic/OpenAI-Compatible 通过 Contract Test（V0.2，FC-PROV-008）

## Documentation Complete
- [ ] SPEC 接口与实现一致
- [ ] 能力元数据来源记录（OPEN_QUESTIONS Q3-PROV）
- [ ] 配置示例（API key 经环境变量，不入日志）

## Release Ready
- [ ] P0 Provider（Mock+OpenAI）验收通过
- [ ] API key 等敏感信息不写入普通日志（NFR-SEC-002）
- [ ] Usage/Cost 上报 telemetry
- [ ] 能力差异已文档化
