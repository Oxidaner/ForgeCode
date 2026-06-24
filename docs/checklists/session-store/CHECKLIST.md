# session-store Checklist

模块：`session-store`。相关需求 FR-SESSION-001..004。

## Design Ready
- [x] EventLog/Store/Checkpointer 接口已定义（EventLog/Store 已实现：FC-SESS-001/002；Checkpointer 留 FC-SESS-003）
- [x] Append-only + Sequence + EventID 去重语义已定义（ADR-0002；Evidence: `internal/session-store/event_log.go`）
- [x] Session State 枚举与 GLOSSARY 一致（Evidence: `internal/session-store/session_store.go`）
- [x] 与 event-system 的"格式 vs 存储"边界明确（EventLog 使用 `eventsystem.Event`，Sequence 由 session-store 分配）
- [ ] 与 memory-system 表所有权分离（共用 DB）
- [ ] 错误模型映射 GLOSSARY（PersistenceError/RecoveryError）

## Implementation Ready
- [x] SQLite 驱动已选型（FC-SESS-000, OPEN_QUESTIONS Q2；Evidence: `modernc.org/sqlite v1.34.5`, `go test ./internal/session-store`）
- [ ] 任务已拆分（Store/Checkpoint/归档/迁移/失败）
- [x] WAL 与写串行化策略已定义（RISK-017；Evidence: `SQLiteEventLog` 写路径互斥 + 事务内分配 Sequence）
- [ ] 迁移与回滚策略已定义

## Implementation Complete
- [x] 单 Session Sequence 严格单调连续（FR-SESSION-001；Evidence: `TestSQLiteEventLogAssignsContiguousSequencePerSession`）
- [x] 重复 EventID 不重复写（Evidence: `TestSQLiteEventLogReturnsExistingSequenceForDuplicateEventID`）
- [x] 可列出未完成 Session 供 Resume（FR-SESSION-002；Evidence: `TestSQLiteStoreListsSessionsByState`）
- [x] 读流支持重放且不触发副作用（FC-SESS-001 只读取事件，不执行外部动作；Evidence: `Read` tests）
- [ ] 写失败重试、不静默丢事件
- [x] 事件 Payload 不入普通日志（当前实现无 Payload 日志输出）

## Test Complete
- [x] Sequence/去重 Unit Test（Evidence: `event_log_test.go`）
- [x] 并发写 Race Test（Evidence: `TestSQLiteEventLogConcurrentAppendsStayContiguous`; full race gate pending final verification）
- [ ] 崩溃后恢复 Recovery Test（NFR-REL-001）
- [ ] 写失败/锁 Failure Injection
- [ ] Schema 迁移测试（NFR-COMPAT-001）
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [x] TASKS 状态更新（FC-SESS-000/001/002）
- [ ] ADR-0002/0006 与实现一致
- [ ] DB 路径/归档配置示例

## Release Ready
- [ ] P0 验收通过
- [ ] 审计 append-only 防篡改可验证
- [ ] DB 大小/写失败指标可观察
- [ ] 升级迁移经过验证
- [ ] Session 恢复 Demo 可复现
