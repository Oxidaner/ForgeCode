package builtin

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	toolruntime "forgecode/internal/tool-runtime"
)

type readFileTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type ReadFileInput struct {
	Path   string `json:"path"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

func (t readFileTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t readFileTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input ReadFileInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid ReadFile input", err)
	}
	if input.Offset < 0 || input.Limit < 0 {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "offset and limit must be positive")
	}
	offset := input.Offset
	if offset == 0 {
		offset = 1
	}
	limit := input.Limit
	if limit == 0 {
		limit = t.deps.Limits.ReadDefaultLimit
	}

	path, err := resolveWorkspacePath(t.deps.WorkspaceRoot, input.Path)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "stat file", err)
	}
	if info.IsDir() {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "ReadFile path is a directory")
	}
	if info.Size() > t.deps.Limits.ReadMaxBytes && input.Limit == 0 {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "file exceeds read max bytes; provide offset and limit to page through it")
	}
	if err := rejectBinaryFile(path, t.deps.Limits.ReadBinaryProbeBytes); err != nil {
		return toolruntime.ToolResult{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "open file", err)
	}
	defer func() { _ = file.Close() }()

	output, startLine, endLine, totalLines, truncated, err := readLinePage(ctx, file, offset, limit)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}

	return toolruntime.ToolResult{
		Output:    output,
		Truncated: truncated,
		Meta: map[string]any{
			"path":        input.Path,
			"start_line":  startLine,
			"end_line":    endLine,
			"total_lines": totalLines,
			"offset":      offset,
			"limit":       limit,
		},
	}, nil
}

func rejectBinaryFile(path string, probeBytes int) error {
	file, err := os.Open(path)
	if err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "open file for binary probe", err)
	}
	defer func() { _ = file.Close() }()

	buf := make([]byte, probeBytes)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "read binary probe", err)
	}
	sample := buf[:n]
	if bytes.Contains(sample, []byte{0}) || !utf8.Valid(sample) {
		return toolruntime.NewError(toolruntime.ValidationError, "binary or non-UTF-8 file is not readable")
	}
	return nil
}

func readLinePage(ctx context.Context, reader io.Reader, offset, limit int) (string, int, int, int, bool, error) {
	bufReader := bufio.NewReader(reader)
	var builder strings.Builder
	lineNo := 0
	written := 0
	startLine := 0
	endLine := 0

	for {
		if err := ctx.Err(); err != nil {
			return "", 0, 0, 0, false, toolruntime.WrapError(toolruntime.CancelledError, "ReadFile cancelled", err)
		}

		line, err := bufReader.ReadString('\n')
		if len(line) > 0 {
			lineNo++
			if lineNo >= offset && written < limit {
				if startLine == 0 {
					startLine = lineNo
				}
				endLine = lineNo
				fmt.Fprintf(&builder, "%d: %s", lineNo, line)
				if !strings.HasSuffix(line, "\n") {
					builder.WriteByte('\n')
				}
				written++
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", 0, 0, 0, false, toolruntime.WrapError(toolruntime.ToolExecutionError, "read file", err)
		}
	}

	truncated := offset+written-1 < lineNo
	return builder.String(), startLine, endLine, lineNo, truncated, nil
}
