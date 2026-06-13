# mcp-client Checklist

模块：`mcp-client`。相关需求 FR-MCP-001..005。RISK-011。

## Design Ready
- [ ] MCPClient/ServerHandle 接口已定义（FC-MCP-001）
- [ ] MCP Server State 枚举与 GLOSSARY 一致
- [ ] MCP 工具映射统一 ToolDescriptor 经 tool-runtime（ADR-0004）
- [ ] 信任级别（Trusted/Limited/Untrusted）与权限映射已定义
- [ ] "默认不可信"原则已写入设计（RISK-011）
- [ ] 错误模型映射 GLOSSARY

## Implementation Ready
- [ ] 任务已拆分（生命周期/transport/能力/接入/安全边界）
- [ ] 输出大小/超时/重连默认值已定义
- [ ] Mock MCP Server 测试边界已定义
- [ ] 信任→权限映射表已确定（OPEN_QUESTIONS）

## Implementation Complete
- [ ] Server 生命周期与重连（FR-MCP-001）
- [ ] stdio + Streamable HTTP（FR-MCP-002，凭证脱敏）
- [ ] tools/resources/prompts 接入（FR-MCP-003）
- [ ] MCP 工具经统一管线、Namespace 防冒充、输出限制（FR-MCP-004）
- [ ] 外部 Prompt/Resource 隔离标注（FR-MCP-005）
- [ ] 断开不阻塞其他 Server

## Test Complete
- [ ] transport Contract Test
- [ ] MCP 工具经 tool-runtime + 权限 Integration
- [ ] 越权/冒充/超大输出 Security Test（RISK-011）
- [ ] 连接断开/重连 Failure Injection
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0004 与实现一致
- [ ] Server 配置与信任级别示例

## Release Ready
- [ ] P0 安全验收通过（无权限绕过）
- [ ] 不可信 Server 受控
- [ ] Server 健康/调用指标可观察
- [ ] MCP 工具调用 Demo 可复现
