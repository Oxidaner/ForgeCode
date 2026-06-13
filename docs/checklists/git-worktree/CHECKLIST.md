# git-worktree Checklist

模块：`git-worktree`。相关需求 FR-WORKTREE-001..003。RISK-014。

## Design Ready
- [ ] WorktreeManager 接口已定义（FC-WT-001）
- [ ] Worktree State 枚举与 GLOSSARY 一致
- [ ] 登记表（孤儿回收支撑）设计已定义
- [ ] 与 agent-orchestration 隔离契约已对齐（ADR-0008）
- [ ] Worktree 逻辑独立于 Git 工具（无混入）
- [ ] 错误模型映射 GLOSSARY（ConflictError 等）

## Implementation Ready
- [ ] git CLI vs go-git 已决策（FC-WT-000, OPEN_QUESTIONS Q3）
- [ ] 任务已拆分（生命周期/合并/边界/审计）
- [ ] 临时分支/路径命名策略已定义
- [ ] 主仓未提交处理策略已定义

## Implementation Complete
- [ ] 临时分支与 Worktree 创建，冲突自动处理（FR-WORKTREE-001）
- [ ] Diff/Merge/Cherry-pick/Discard 可用
- [ ] 全部边界条件被处理（FR-WORKTREE-002）
- [ ] 崩溃孤儿被扫描清理（FR-WORKTREE-003, RISK-014）
- [ ] 取消时回收 Worktree
- [ ] Worktree 在受控根目录（权限边界）

## Test Complete
- [ ] 生命周期 Integration（真实 git）
- [ ] 边界条件 Failure Injection（冲突/占用/清理失败等）
- [ ] 孤儿回收 Recovery Test
- [ ] 合并冲突不丢变更 Test
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0008 与实现一致
- [ ] Worktree 配置示例

## Release Ready
- [ ] P0 边界与孤儿验收通过
- [ ] 无未处理 Critical 风险（RISK-014 缓解）
- [ ] Worktree 审计与指标可观察
- [ ] 并行写隔离 Demo 可复现
