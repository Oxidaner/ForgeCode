# extension-system Checklist

模块：`extension-system`（Skill + Command + Hook）。相关需求 FR-CMD-001/002/100、FR-HOOK-001..003、FR-SKILL-001..004。

## Design Ready
- [ ] CommandRegistry/HookDispatcher/SkillManager 接口已定义
- [ ] CommandKind（Fixed/Prompt/Skill）区分已定义
- [ ] Hook 基于 Event Bus（不散落回调，ADR-0011）
- [ ] Skill State 与候选状态机已定义（与 GLOSSARY 一致）
- [ ] "不可扩权"约束已写入设计（Hook/Skill）
- [ ] 候选 Skill 审批前不得 Active（ADR-0010）

## Implementation Ready
- [ ] 任务已拆分（命令/Hook/Skill/候选/标杆）
- [ ] Shell/HTTP Hook 默认禁用策略已定义
- [ ] Skill 搜索路径与依赖检查策略已定义
- [ ] evaluation 回放接口已约定（FC-EVAL-004）

## Implementation Complete
- [ ] 固定命令本地执行；Prompt/Skill 展开为任务（FR-CMD-001）
- [ ] 参数校验/Alias/冲突/Help 生效（FR-CMD-002）
- [ ] Hook 订阅生命周期事件、按优先级执行（FR-HOOK-001/002）
- [ ] Hook 超时/失败策略/递归防护/不可提权（FR-HOOK-003, RISK-012）
- [ ] Skill 依赖检查 + 延迟加载 + 权限仅收窄（FR-SKILL-001/003）
- [ ] 候选 Skill 未审批不 Active（FR-SKILL-004）
- [ ] 三标杆命令已注册（FR-CMD-100）

## Test Complete
- [ ] 命令解析/冲突 Unit Test
- [ ] Hook 提权/递归/超时 Security Test
- [ ] HookDispatcher 与 event-system/permission-engine Integration
- [ ] Skill 依赖检查 Test
- [ ] 候选 Skill 回放 Eval（FR-SKILL-004）
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0010/0011 与实现一致
- [ ] Hook/Skill 配置与安全说明更新

## Release Ready
- [ ] P0（Command+Hook）验收通过
- [ ] 无未处理 Critical 安全风险（Hook 无法提权/绕过）
- [ ] Hook/命令/Skill 指标与审计可观察
- [ ] 标杆命令 Demo 可复现
