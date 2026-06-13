# Checklist 模板

> 复制本模板创建 `docs/checklists/<module-id>/CHECKLIST.md`。Checklist 不得只是重复 Task 标题，应描述可验证的状态。

每项格式：

```text
- [ ] 检查项描述
  - Evidence:
  - Related Task:
  - Related Requirement:
  - Blocking:
```

## Design Ready
- [ ] 职责边界已明确
- [ ] 非目标已明确
- [ ] 接口已定义
- [ ] 状态机已定义（如适用）
- [ ] 错误模型已定义
- [ ] 安全边界已定义
- [ ] 并发语义已定义
- [ ] 持久化所有权已定义
- [ ] 依赖方向无环
- [ ] 相关 ADR 已完成

## Implementation Ready
- [ ] Task 已拆分
- [ ] P0 依赖已满足
- [ ] 配置已定义
- [ ] 测试策略已定义
- [ ] Mock/Fake 边界已定义
- [ ] Migration 策略已定义（如适用）
- [ ] 回滚策略已定义（如适用）

## Implementation Complete
- [ ] 核心路径完成
- [ ] 异常路径完成
- [ ] Context Cancellation 生效
- [ ] 错误可识别（标准错误分类）
- [ ] 事件已记录
- [ ] 配置有默认值
- [ ] 敏感数据未写入普通日志

## Test Complete
- [ ] Unit Test
- [ ] Integration Test
- [ ] Race Test
- [ ] Failure Test
- [ ] Security Test
- [ ] Contract Test（如适用）
- [ ] Eval Case（如适用）
- [ ] 覆盖关键恢复路径

## Documentation Complete
- [ ] Spec 更新
- [ ] Task 状态更新
- [ ] ADR 更新
- [ ] README 更新
- [ ] 配置示例更新
- [ ] 操作手册更新（如适用）
- [ ] 已知限制更新

## Release Ready
- [ ] 所有 P0 验收通过
- [ ] 没有未处理 Critical 风险
- [ ] 关键指标可观察
- [ ] 升级与回滚经过验证（如适用）
- [ ] Demo 可复现
- [ ] Evidence 已记录
