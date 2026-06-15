package builtin

import (
	"encoding/json"

	toolruntime "forgecode/internal/tool-runtime"
)

const (
	ToolReadFile  = "ReadFile"
	ToolWriteFile = "WriteFile"
	ToolEditFile  = "EditFile"
	ToolBash      = "Bash"
	ToolGlob      = "Glob"
	ToolGrep      = "Grep"
)

func Descriptors() []toolruntime.ToolDescriptor {
	return []toolruntime.ToolDescriptor{
		{
			Name:        ToolReadFile,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(readFileSchema),
			Risk:        toolruntime.RiskLow,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionRead}},
		},
		{
			Name:        ToolWriteFile,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(writeFileSchema),
			Risk:        toolruntime.RiskMedium,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionWrite}},
		},
		{
			Name:        ToolEditFile,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(editFileSchema),
			Risk:        toolruntime.RiskMedium,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionWrite}},
		},
		{
			Name:        ToolBash,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(bashSchema),
			Risk:        toolruntime.RiskHigh,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionExecute}},
		},
		{
			Name:        ToolGlob,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(globSchema),
			Risk:        toolruntime.RiskLow,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionRead, toolruntime.PermissionSearch}},
		},
		{
			Name:        ToolGrep,
			Source:      toolruntime.ToolSourceBuiltin,
			InputSchema: mustSchema(grepSchema),
			Risk:        toolruntime.RiskLow,
			Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionRead, toolruntime.PermissionSearch}},
		},
	}
}

func mustSchema(schema string) json.RawMessage {
	raw := json.RawMessage(schema)
	if !json.Valid(raw) {
		panic("invalid built-in tool schema")
	}
	return raw
}

const readFileSchema = `{
  "type": "object",
  "required": ["path"],
  "properties": {
    "path": {"type": "string", "minLength": 1},
    "offset": {"type": "integer", "minimum": 1},
    "limit": {"type": "integer", "minimum": 1}
  },
  "additionalProperties": false
}`

const writeFileSchema = `{
  "type": "object",
  "required": ["path", "content"],
  "properties": {
    "path": {"type": "string", "minLength": 1},
    "content": {"type": "string"}
  },
  "additionalProperties": false
}`

const editFileSchema = `{
  "type": "object",
  "required": ["path", "old_string", "new_string"],
  "properties": {
    "path": {"type": "string", "minLength": 1},
    "old_string": {"type": "string", "minLength": 1},
    "new_string": {"type": "string"}
  },
  "additionalProperties": false
}`

const bashSchema = `{
  "type": "object",
  "required": ["command"],
  "properties": {
    "command": {"type": "string", "minLength": 1},
    "timeout_ms": {"type": "integer", "minimum": 1}
  },
  "additionalProperties": false
}`

const globSchema = `{
  "type": "object",
  "required": ["pattern"],
  "properties": {
    "pattern": {"type": "string", "minLength": 1},
    "root": {"type": "string"}
  },
  "additionalProperties": false
}`

const grepSchema = `{
  "type": "object",
  "required": ["pattern"],
  "properties": {
    "pattern": {"type": "string", "minLength": 1},
    "path": {"type": "string"},
    "regex": {"type": "boolean"}
  },
  "additionalProperties": false
}`
