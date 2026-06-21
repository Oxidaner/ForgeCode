package builtin

import (
	"bufio"
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	toolruntime "forgecode/internal/tool-runtime"
)

type grepTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type GrepInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
	Regex   bool   `json:"regex,omitempty"`
}

func (t grepTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t grepTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input GrepInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid Grep input", err)
	}
	if input.Pattern == "" {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "pattern is required")
	}

	var compiled *regexp.Regexp
	if input.Regex {
		re, err := regexp.Compile(input.Pattern)
		if err != nil {
			return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid regular expression", err)
		}
		compiled = re
	}

	root := resolveOptionalWorkspacePath(t.deps.WorkspaceRoot, input.Path)
	files, err := grepFiles(root)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}

	seen := map[string]bool{}
	results := make([]string, 0)
	truncated := false
	for _, filePath := range files {
		if err := ctx.Err(); err != nil {
			return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.CancelledError, "Grep cancelled", err)
		}
		fileMatches, err := grepFile(ctx, filePath, t.deps.WorkspaceRoot, input.Pattern, compiled)
		if err != nil {
			return toolruntime.ToolResult{}, err
		}
		for _, match := range fileMatches {
			key := match.key
			if seen[key] {
				continue
			}
			seen[key] = true
			results = append(results, match.output)
			if len(results) >= t.deps.Limits.GrepMaxMatches {
				truncated = true
				break
			}
		}
		if truncated {
			break
		}
	}

	return toolruntime.ToolResult{
		Output:    strings.Join(results, "\n"),
		Truncated: truncated,
		Meta: map[string]any{
			"matches": len(results),
			"regex":   input.Regex,
		},
	}, nil
}

func grepFiles(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, toolruntime.WrapError(toolruntime.ToolExecutionError, "stat grep path", err)
	}
	if !info.IsDir() {
		return []string{root}, nil
	}

	files := make([]string, 0)
	err = filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		files = append(files, filePath)
		return nil
	})
	if err != nil {
		return nil, toolruntime.WrapError(toolruntime.ToolExecutionError, "walk grep files", err)
	}
	return files, nil
}

type grepMatch struct {
	key    string
	output string
}

func grepFile(ctx context.Context, filePath, workspaceRoot, patternValue string, re *regexp.Regexp) ([]grepMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, toolruntime.WrapError(toolruntime.ToolExecutionError, "open grep file", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	matches := make([]grepMatch, 0)
	rel := relativeToWorkspace(workspaceRoot, filePath)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, toolruntime.WrapError(toolruntime.CancelledError, "Grep cancelled", err)
		}
		lineNo++
		line := scanner.Text()
		matched := strings.Contains(line, patternValue)
		if re != nil {
			matched = re.MatchString(line)
		}
		if matched {
			key := rel + ":" + strconvLine(lineNo)
			matches = append(matches, grepMatch{
				key:    key,
				output: key + ":" + line,
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, toolruntime.WrapError(toolruntime.ToolExecutionError, "scan grep file", err)
	}
	return matches, nil
}

func strconvLine(line int) string {
	return strconv.Itoa(line)
}
