# cli Checklist

模块：`cli`。相关需求 FR-CLI-001..004。

## Design Ready
- [ ] App/REPL/Renderer 接口已定义（FC-CLI-001）
- [ ] 任务 vs Slash Command 解析策略已定义
- [ ] 审批交互 UX 已定义
- [ ] "无核心业务逻辑"约束已写入设计（FR-CLI-003）
- [ ] 依赖方向：仅依赖接口，无环

## Implementation Ready
- [ ] 任务已拆分（REPL/渲染/审批/命令/取消/恢复）
- [ ] CLI 框架已决策（OPEN_QUESTIONS Q5）
- [ ] 固定命令清单已定义
- [ ] Fake Runtime 测试边界已定义

## Implementation Complete
- [ ] 可提交任务并流式渲染（FR-CLI-001）
- [ ] 固定命令本地执行不经模型（FR-CLI-002）
- [ ] 审批提示展示风险/命中原因并脱敏（FR-PERM-007）
- [ ] Ctrl-C 传播取消（FR-RUNTIME-003）
- [ ] 命令逻辑在 extension-system/runtime-core 而非 cli（FR-CLI-003）
- [ ] 终端不回显密钥/Token

## Test Complete
- [ ] 输入解析 Unit Test
- [ ] 提交→流式→审批→完成 Integration
- [ ] 渲染 Golden Test
- [ ] 依赖检查：cli 不含核心业务逻辑
- [ ] `go test` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] 命令帮助文档更新
- [ ] 配置示例更新

## Release Ready
- [ ] P0 验收通过
- [ ] 无敏感信息回显
- [ ] 用户操作指标可观察
- [ ] MVP 六条流程经 CLI 可演示
