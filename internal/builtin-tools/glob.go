package builtin

import (
	"context"
	"encoding/json"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	toolruntime "forgecode/internal/tool-runtime"
)

type globTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type GlobInput struct {
	Pattern string `json:"pattern"`
	Root    string `json:"root,omitempty"`
}

func (t globTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t globTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input GlobInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid Glob input", err)
	}
	if input.Pattern == "" {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "pattern is required")
	}

	root := resolveOptionalWorkspacePath(t.deps.WorkspaceRoot, input.Root)
	matches := make([]string, 0)
	truncated := false
	err := filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel := relativeToWorkspace(t.deps.WorkspaceRoot, filePath)
		if matchesGlob(input.Pattern, rel) {
			matches = append(matches, rel)
			if len(matches) >= t.deps.Limits.GlobMaxResults {
				truncated = true
				return filepath.SkipAll
			}
		}
		return nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.CancelledError, "Glob cancelled", ctx.Err())
		}
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "walk files", err)
	}

	sort.Strings(matches)
	return toolruntime.ToolResult{
		Output:    strings.Join(matches, "\n"),
		Truncated: truncated,
		Meta: map[string]any{
			"matches": len(matches),
			"pattern": input.Pattern,
		},
	}, nil
}

func matchesGlob(patternValue, rel string) bool {
	normalizedPattern := filepath.ToSlash(patternValue)
	normalizedRel := filepath.ToSlash(rel)
	if ok, _ := path.Match(normalizedPattern, normalizedRel); ok {
		return true
	}
	if !strings.Contains(normalizedPattern, "/") {
		if ok, _ := path.Match(normalizedPattern, path.Base(normalizedRel)); ok {
			return true
		}
	}
	return false
}
