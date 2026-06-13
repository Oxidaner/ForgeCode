# memory-system Checklist

模块：`memory-system`。相关需求 FR-MEMORY-001..004。RISK-010。

## Design Ready
- [ ] Memory 模型与四类记忆已定义（FC-MEM-001）
- [ ] MemoryStore/CandidateReview 接口已定义
- [ ] FTS5 检索（非向量库）方案已定义（ADR-0007）
- [ ] 与 session-store 表所有权分离
- [ ] 候选状态机已定义（与 GLOSSARY 一致）
- [ ] 污染控制与项目隔离设计已定义

## Implementation Ready
- [ ] 任务已拆分（存储/检索/元数据/审批/敏感）
- [ ] FTS5 可用性已确认（FC-SESS-000）
- [ ] 默认置信度/TTL/审批策略已定义
- [ ] Migration 策略已定义

## Implementation Complete
- [ ] 四类记忆可存储检索（FR-MEMORY-001/002）
- [ ] FTS5 按置信度/时效过滤
- [ ] 候选未审批不 Active（FR-MEMORY-004, RISK-010）
- [ ] 跨项目访问被拒（项目隔离）
- [ ] 敏感信息脱敏/拒绝（FR-MEMORY-003）
- [ ] MemoryRead/MemoryWrite 事件记录

## Test Complete
- [ ] 检索/隔离 Unit Test
- [ ] 污染控制 Security Test（未审批不生效）
- [ ] 敏感信息拒绝 Security Test
- [ ] 与 context-manager Integration
- [ ] Schema 迁移测试
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0007 与实现一致
- [ ] 记忆配置示例

## Release Ready
- [ ] P0 安全验收通过（污染/隔离/敏感）
- [ ] 无未处理 Critical 风险（RISK-010 缓解）
- [ ] 检索命中/候选审批指标可观察
- [ ] 跨会话记忆 Demo 可复现
