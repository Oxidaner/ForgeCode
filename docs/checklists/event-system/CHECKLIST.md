# event-system Checklist

模块：`event-system`。相关需求 FR-EVENT-001..003。

## Design Ready
- [x] Event Envelope 字段与 EVENT_MODEL.md 完全一致（FC-EVT-001；Evidence: `go test ./internal/event-system`）
- [x] EventType 权威枚举与 EventClass 分类已定义（Evidence: `internal/event-system/types.go`）
- [x] Bus 接口（Publish/Subscribe/Filter）已定义（Evidence: `internal/event-system/types.go`）
- [x] "格式 vs 存储"边界明确（Sequence 由 session-store 分配）
- [ ] 错误隔离与有序投递语义已定义
- [x] 依赖方向无环（当前实现仅依赖 Go 标准库）

## Implementation Ready
- [ ] 任务已拆分（Envelope/Bus/过滤/超时/测试）
- [ ] 同步阻断 vs 异步通知订阅接口已决策（OPEN_QUESTIONS）
- [ ] Bus 缓冲/超时默认值已定义
- [ ] Fake Subscriber 测试边界已定义

## Implementation Complete
- [ ] 单订阅者有序投递（FR-EVENT-002）
- [ ] 订阅者失败被隔离、不影响发布方
- [ ] 按 Type/Class 过滤生效（FR-EVENT-003）
- [ ] 订阅者超时被隔离并计数
- [ ] Payload 不在普通日志打印

## Test Complete
- [x] Envelope Contract Test 锁定字段（Evidence: `internal/event-system/event_test.go`）
- [ ] 并发发布/订阅 Race Test
- [ ] 错误/超时隔离 Failure Injection
- [ ] `go test -race` 通过

## Documentation Complete
- [x] SPEC 与 EVENT_MODEL.md 一致
- [x] TASKS 状态更新
- [ ] ADR-0011 与实现一致
- [ ] 已知限制（背压）记录

## Release Ready
- [ ] P0 验收通过
- [ ] 无未处理 Critical 风险
- [ ] Bus 指标可观察
- [ ] 事件驱动闭环 Demo 可复现
