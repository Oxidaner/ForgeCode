# ForgeCode

ForgeCode is a model-agnostic coding agent runtime written in Go.

The project aims to implement the core control plane of a coding agent runtime from first principles, instead of wrapping a single provider SDK or delegating the runtime model to an existing agent framework.

## Goals

ForgeCode is designed to be:

- Model agnostic, with provider adapters for OpenAI, Anthropic, OpenAI-compatible APIs, and mock providers.
- Tool oriented, with built-in coding tools such as `ReadFile`, `WriteFile`, `EditFile`, `Bash`, `Glob`, and `Grep`.
- Safe by default, with explicit permission checks, approval flows, sandboxing, checkpoints, and audit logs.
- Recoverable, using an explicit agent state machine and append-only event records.
- Extensible, with MCP client support, lifecycle hooks, slash commands, and skill packages.
- Parallel-capable, with SubAgents, Git worktree isolation, and Agent Teams.

## Current Stage

ForgeCode is currently in the architecture and planning phase.

The active master plan is maintained in [docs/master-plan.md](docs/master-plan.md). This phase is limited to Markdown, YAML, JSON, and other planning documents. Production code and dependency installation are intentionally out of scope until the architecture, module specs, task breakdowns, checklists, ADRs, traceability matrix, roadmap, and risk register are ready.

## Planned Capability Areas

- Agent runtime and autonomous tool-use loop
- Model provider abstraction
- Built-in tool registry and executor
- MCP client integration
- Permission engine and approval policy
- Context compression and token budget management
- Cross-session memory
- Skill package system
- Slash command framework
- Hook lifecycle system
- SubAgent orchestration
- Git worktree-based parallel isolation
- Agent Team coordination
- PR, SQL, and Kubernetes change review scenarios

## Documentation Roadmap

The master plan calls for the repository to grow a planning document set that includes:

- Project overview and system architecture specs
- Module-level specs, tasks, and checklists
- Cross-module dependency graph and implementation order
- Architecture Decision Records
- Requirement-to-module-to-task traceability matrix
- 12-week implementation roadmap
- Risk register and open validation questions

## Status

No production runtime code has been implemented yet. The next step is to turn the master plan into the first set of architecture and planning documents.
