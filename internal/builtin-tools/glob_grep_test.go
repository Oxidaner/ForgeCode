package builtin

import (
	"context"
	"strings"
	"testing"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestGlobReturnsSortedLimitedMatches(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "b.go", "package b\n")
	writeTestFile(t, root, "nested/a.go", "package a\n")
	writeTestFile(t, root, "notes.txt", "skip\n")
	tool := NewGlobTool(Deps{WorkspaceRoot: root, Limits: Limits{GlobMaxResults: 10}})

	result, err := tool.Execute(context.Background(), mustJSON(t, GlobInput{Pattern: "*.go"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.Output != "b.go\nnested/a.go" {
		t.Fatalf("unexpected glob output: %q", result.Output)
	}
}

func TestGlobMarksTruncatedAtLimit(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.go", "")
	writeTestFile(t, root, "b.go", "")
	tool := NewGlobTool(Deps{WorkspaceRoot: root, Limits: Limits{GlobMaxResults: 1}})

	result, err := tool.Execute(context.Background(), mustJSON(t, GlobInput{Pattern: "*.go"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.Truncated {
		t.Fatal("expected truncated glob result")
	}
}

func TestGrepSupportsStringRegexAndLimit(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "needle one\nnone\nneedle two\n")
	writeTestFile(t, root, "b.txt", "needle three\n")
	tool := NewGrepTool(Deps{WorkspaceRoot: root, Limits: Limits{GrepMaxMatches: 2}})

	result, err := tool.Execute(context.Background(), mustJSON(t, GrepInput{Pattern: "needle"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.Truncated {
		t.Fatal("expected grep truncation at max matches")
	}
	if !strings.Contains(result.Output, "a.txt:1:needle one") || !strings.Contains(result.Output, "a.txt:3:needle two") {
		t.Fatalf("unexpected grep output: %q", result.Output)
	}

	regexTool := NewGrepTool(Deps{WorkspaceRoot: root})
	regexResult, err := regexTool.Execute(context.Background(), mustJSON(t, GrepInput{Pattern: `needle\s+three`, Regex: true}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(regexResult.Output, "b.txt:1:needle three") {
		t.Fatalf("unexpected regex grep output: %q", regexResult.Output)
	}
}

func TestGrepRejectsInvalidRegex(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "content\n")
	tool := NewGrepTool(Deps{WorkspaceRoot: root})

	_, err := tool.Execute(context.Background(), mustJSON(t, GrepInput{Pattern: "[", Regex: true}))
	if !toolruntime.IsCategory(err, toolruntime.ValidationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
