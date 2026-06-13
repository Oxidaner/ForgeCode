# evaluation Checklist

模块：`evaluation`。相关需求 FR-EVAL-001..004。RISK-018/019。

## Design Ready
- [ ] Replayer/EvalRunner/Scorer 接口已定义（FC-EVAL-001/002）
- [ ] Replay 不重执行副作用语义已定义
- [ ] Eval Case 模型（规则断言 + 可选模型评审）已定义
- [ ] 三标杆场景检查项已定义（SPEC §10）
- [ ] 候选 Skill 回放（ADR-0010 安全门）已对齐
- [ ] 错误模型映射 GLOSSARY

## Implementation Ready
- [ ] 任务已拆分（Replay/框架/场景/候选/导出）
- [ ] 评分方式已决策（OPEN_QUESTIONS Q15）
- [ ] 标杆案例来源已确定（RISK-019）
- [ ] Mock Provider 测试边界已定义

## Implementation Complete
- [ ] Replay 不触发外部副作用（FR-EVAL-001）
- [ ] Eval 框架支持断言 + 可选评审（FR-EVAL-002）
- [ ] 三场景检查项均有断言并产出得分（FR-EVAL-003）
- [ ] 候选 Skill 回放产出建议（FR-EVAL-004）
- [ ] EvalReport 可导出

## Test Complete
- [ ] Scorer 断言 Unit Test
- [ ] Replay/场景 Integration
- [ ] 场景期望 Golden 快照
- [ ] CI 可运行场景 Eval
- [ ] `go test` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] 标杆场景检查项文档更新
- [ ] Eval 用例编写指南

## Release Ready
- [ ] 三场景 Eval 验收通过（RISK-018 缓解）
- [ ] Demo 案例有足够难度（RISK-019 缓解）
- [ ] Eval 通过率/得分可观察
- [ ] 候选 Skill 安全门可复现
