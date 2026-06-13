# telemetry Checklist

模块：`telemetry`。相关需求 FR-TELEMETRY-001..004。

## Design Ready
- [ ] Logger/Metrics/AuditSink/UsageMeter/Redactor 接口已定义（FC-TEL-001）
- [ ] 脱敏规则与默认敏感字段已定义（NFR-SEC-002）
- [ ] AuditSink 依赖反转设计（实现 Subscriber，不 import Bus）
- [ ] UsageRecord 模型与聚合作用域已定义
- [ ] 无业务依赖（DEPENDENCY_GRAPH 一致）

## Implementation Ready
- [ ] 任务已拆分（Logger/Metrics/Audit/Usage/测试）
- [ ] 指标导出后端已决策（OPEN_QUESTIONS）
- [ ] UsageRecord 表与 session-store 分离
- [ ] best-effort 失败策略已定义

## Implementation Complete
- [ ] 密钥/Token/环境变量不入普通日志（FR-TELEMETRY-001, NFR-SEC-002）
- [ ] 关键指标可记录查询（FR-TELEMETRY-002）
- [ ] 审计事件被 AuditSink 完整落地（FR-TELEMETRY-003）
- [ ] UsageRecord 可按作用域聚合（FR-TELEMETRY-004）
- [ ] 日志/指标失败不阻断主流程

## Test Complete
- [ ] 脱敏 Security Test（含密钥输入）
- [ ] Usage 聚合 Unit Test
- [ ] AuditSink Integration
- [ ] 依赖检查：telemetry 不 import 业务模块
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] 脱敏字段配置示例
- [ ] 已知限制（脱敏误杀）记录

## Release Ready
- [ ] P0 验收通过
- [ ] 无敏感数据泄露（Evidence）
- [ ] 审计 append-only 防篡改可验证
- [ ] Usage/Cost 指标可观察
