# session-store Checklist

模块：`session-store`。相关需求 FR-SESSION-001..004。

## Design Ready
- [ ] EventLog/Store/Checkpointer 接口已定义（FC-SESS-001/002/003）
- [ ] Append-only + Sequence + EventID 去重语义已定义（ADR-0002）
- [ ] Session State 枚举与 GLOSSARY 一致
- [ ] 与 event-system 的"格式 vs 存储"边界明确
- [ ] 与 memory-system 表所有权分离（共用 DB）
- [ ] 错误模型映射 GLOSSARY（PersistenceError/RecoveryError）

## Implementation Ready
- [ ] SQLite 驱动已选型（FC-SESS-000, OPEN_QUESTIONS Q2）
- [ ] 任务已拆分（Store/Checkpoint/归档/迁移/失败）
- [ ] WAL 与写串行化策略已定义（RISK-017）
- [ ] 迁移与回滚策略已定义

## Implementation Complete
- [ ] 单 Session Sequence 严格单调连续（FR-SESSION-001）
- [ ] 重复 EventID 不重复写
- [ ] 可列出未完成 Session 供 Resume（FR-SESSION-002）
- [ ] 读流支持重放且不触发副作用（FR-SESSION-003）
- [ ] 写失败重试、不静默丢事件
- [ ] 事件 Payload 不入普通日志

## Test Complete
- [ ] Sequence/去重 Unit Test
- [ ] 并发写 Race Test
- [ ] 崩溃后恢复 Recovery Test（NFR-REL-001）
- [ ] 写失败/锁 Failure Injection
- [ ] Schema 迁移测试（NFR-COMPAT-001）
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0002/0006 与实现一致
- [ ] DB 路径/归档配置示例

## Release Ready
- [ ] P0 验收通过
- [ ] 审计 append-only 防篡改可验证
- [ ] DB 大小/写失败指标可观察
- [ ] 升级迁移经过验证
- [ ] Session 恢复 Demo 可复现
