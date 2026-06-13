# SECURITY_MODEL

ForgeCode 的安全模型围绕 **统一工具调用管线** 与 **五层纵深权限防御** 建立。所有工具（内置 / MCP / Skill 触发）调用路径必须经过：

```text
Validation → Permission → Hook → Execution → Audit
```

任何绕过该管线的执行都视为安全缺陷（见一致性检查 §17.8）。

## Threat Model（资产 / 威胁 / 对手）

- **资产**：用户源代码与密钥、本地文件系统、Git 历史、外部服务凭证、执行主机。
- **对手**：恶意/被劫持的 Model 响应、不可信 MCP Server、恶意 Skill/Hook、被污染的记忆、用户误操作。
- **核心原则**：模型输出、工具输出、MCP 输出、外部 Prompt/Resource 均为 **不可信输入**；高风险操作默认需审批；最严格决策生效。

## Trust Boundaries

| 边界 | 信任级别 | 控制措施 |
| --- | --- | --- |
| User ↔ CLI | 半信任（可发起高危请求） | 审批交互、命令权限要求 |
| Runtime ↔ Model Provider | 响应不可信 | Tool Call schema 校验、注入防护、权限引擎 |
| Runtime ↔ MCP Server | 不可信 | 信任级别、Namespace、输出限制、独立权限等级、审计 |
| Runtime ↔ Tool 执行 | 受控 | 资源边界、Bash 分析、Sandbox |
| Skill/Hook ↔ Runtime | 受控 | 权限声明不可自我扩权、来源追踪、审计 |
| Memory ↔ Runtime | 受控 | 候选审批、污染控制、项目隔离 |

## 威胁清单与缓解

| 威胁 | 缓解 | 责任模块 |
| --- | --- | --- |
| Prompt Injection | 工具输出/外部内容标注来源；不让模型输出直接成为权限决策依据；高危操作需审批 | runtime-core, permission-engine |
| Tool Output Injection | 工具结果截断+来源标注；结果不作为可信指令执行 | tool-runtime, context-manager |
| Path Traversal | 规范化路径、限制在 Workspace Root、拒绝 `..` 逃逸 | permission-engine(L2) |
| Symlink Escape | 解析符号链接真实路径后再做边界判断 | permission-engine(L2) |
| Command Injection | Bash 结构化分析（程序/参数/管道/重定向/子 Shell/命令替换），非整串匹配 | permission-engine(L3) |
| Secret Leakage | 敏感目录/密钥文件读取限制、日志脱敏、环境变量过滤 | permission-engine(L2), telemetry, sandbox |
| MCP Supply Chain | Server 信任级别、能力协商、输出大小限制、独立审计 | mcp-client, permission-engine |
| Malicious Skill | 权限声明上限、依赖检查、来源追踪、自动生成需评测+审批 | extension-system, evaluation |
| Malicious Hook | Hook 不可静默提权、Timeout、失败策略、递归防护、审计 | extension-system |
| Sandbox Escape | 只读挂载、网络控制、资源/PID 限制；逃逸假设下仍受 L1–L3 约束 | sandbox |
| Approval Bypass | 决策与执行分离、Deny 优先、审批后执行前崩溃可恢复重判 | permission-engine, session-store |
| Memory Poisoning | 不自动写入模型输出、候选审批、置信度/过期、项目隔离 | memory-system |
| Audit Tampering | Append-only 审计事件、序号单调、EventID 去重 | session-store, telemetry |

## 五层纵深权限防御（概览）

完整设计见 `docs/specs/permission-engine/SPEC.md`。

1. **L1 Schema/输入校验**：参数合法性、必填、长度、路径格式、空字节、非法编码、命令字段结构、JSON Schema、Tool Call 注入风险。
2. **L2 资源边界**：Workspace Root、可读/可写/敏感目录、路径穿越、符号链接逃逸、隐藏/密钥文件、环境变量、用户/系统目录。
3. **L3 操作风险策略**：风险等级 `Low/Medium/High/Critical`；决策 `Allow/AskOnce/AskAlways/Deny`；Bash 结构化分析（网络访问、文件删除、Git 危险操作、Docker/K8s、DB 写入、Force Push、下载后执行、提权）。
4. **L4 运行时沙箱**：Docker（V0.2）——工作目录挂载、只读挂载、网络/CPU/内存/PID/时间限制、环境变量过滤、进程回收。
5. **L5 人工审批与审计**：记录发起者、Session、Agent、Tool、原始参数、风险等级、规则命中原因、批准/拒绝、执行结果、文件变更、时间、审计事件 ID。

## 决策优先级与冲突规则

- **Deny 优先**：任一层 Deny 即终止。
- **最严格生效**：多来源（用户配置 / Skill 声明 / Hook 返回 / 默认策略）冲突时取最严格。
- **Skill/Hook 不可扩权**：声明只能收窄、不能放宽默认权限。
- **AskAlways > AskOnce > Allow**：同等风险下更保守者生效。
- 决策结果与命中原因落审计事件，可追溯。
