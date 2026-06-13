# sandbox Checklist

模块：`sandbox`。相关需求 FR-SANDBOX-001..003, FR-PERM-006。

## Design Ready
- [ ] Sandbox/ExecSpec/ExecResult 接口已定义（FC-SBX-001）
- [ ] Docker 封装（非自研运行时）已定义（ADR-0012）
- [ ] 资源限制/网络/挂载/env 过滤模型已定义
- [ ] 降级策略已定义（FR-SANDBOX-003）
- [ ] 作为 permission-engine L4 的挂钩契约已对齐（FR-PERM-006）
- [ ] 错误模型映射 GLOSSARY（SandboxError）

## Implementation Ready
- [ ] 任务已拆分（执行/限制/降级/回收）
- [ ] 基础镜像已选型（OPEN_QUESTIONS）
- [ ] 默认资源限制/超时已定义
- [ ] Docker 可用性检测策略已定义

## Implementation Complete
- [ ] 命令在容器执行，挂载/网络生效（FR-SANDBOX-001）
- [ ] CPU/内存/PID/时间限制 + env 白名单（FR-SANDBOX-002）
- [ ] Docker 不可用按策略降级并审计（FR-SANDBOX-003）
- [ ] 超时/取消终止并回收，无孤儿容器
- [ ] 容器敏感输出不入普通日志

## Test Complete
- [ ] Docker Integration（有 Docker 时）
- [ ] 降级 Failure Injection
- [ ] 网络隔离/挂载越界/env 泄露 Security Test
- [ ] 超时/取消回收 Test
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0012 与实现一致
- [ ] 沙箱配置示例（资源/网络/降级）

## Release Ready
- [ ] P1 验收通过
- [ ] 隔离与降级安全可验证
- [ ] 执行/降级指标可观察
- [ ] 沙箱执行 Demo 可复现（含降级路径）
