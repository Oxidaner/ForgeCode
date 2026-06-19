package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	toolruntime "forgecode/internal/tool-runtime"
)

type bashTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type BashInput struct {
	Command   string `json:"command"`
	TimeoutMs int    `json:"timeout_ms,omitempty"`
}

func (t bashTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t bashTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input BashInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid Bash input", err)
	}
	if input.Command == "" {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "command is required")
	}

	timeout := t.deps.Limits.BashTimeout
	if input.TimeoutMs > 0 {
		timeout = time.Duration(input.TimeoutMs) * time.Millisecond
		if timeout > t.deps.Limits.BashMaxTimeout {
			timeout = t.deps.Limits.BashMaxTimeout
		}
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := shellCommand(runCtx, input.Command)
	if t.deps.WorkspaceRoot != "" {
		cmd.Dir = t.deps.WorkspaceRoot
	}
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined

	err := cmd.Run()
	output, truncated := truncateHeadTail(combined.Bytes(), t.deps.Limits.BashHeadBytes, t.deps.Limits.BashTailBytes)
	result := toolruntime.ToolResult{
		Output:    output,
		Truncated: truncated,
		Meta: map[string]any{
			"timeout_ms": timeout.Milliseconds(),
		},
	}

	if runCtx.Err() != nil {
		result.IsError = true
		if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
			result.Category = toolruntime.TimeoutError
			return result, toolruntime.WrapError(toolruntime.TimeoutError, "Bash command timed out", runCtx.Err())
		}
		result.Category = toolruntime.CancelledError
		return result, toolruntime.WrapError(toolruntime.CancelledError, "Bash command cancelled", runCtx.Err())
	}
	if err != nil {
		result.IsError = true
		result.Category = toolruntime.ToolExecutionError
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.Meta["exit_code"] = exitErr.ExitCode()
		}
		return result, toolruntime.WrapError(toolruntime.ToolExecutionError, "Bash command failed", err)
	}

	result.Meta["exit_code"] = 0
	return result, nil
}

func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

func truncateHeadTail(data []byte, headBytes, tailBytes int) (string, bool) {
	limit := headBytes + tailBytes
	if limit <= 0 || len(data) <= limit {
		return string(data), false
	}
	var out bytes.Buffer
	out.Write(data[:headBytes])
	out.WriteString("\n...[truncated " + strconv.Itoa(len(data)-limit) + " bytes]...\n")
	out.Write(data[len(data)-tailBytes:])
	return out.String(), true
}
