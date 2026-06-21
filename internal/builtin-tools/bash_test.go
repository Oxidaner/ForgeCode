package builtin

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestBashRunsCommandAndCapturesOutput(t *testing.T) {
	tool := NewBashTool(Deps{WorkspaceRoot: t.TempDir()})

	result, err := tool.Execute(context.Background(), mustJSON(t, BashInput{Command: "echo hello"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Fatalf("expected output to contain hello, got %q", result.Output)
	}
	if result.Meta["exit_code"] != 0 {
		t.Fatalf("unexpected metadata: %#v", result.Meta)
	}
}

func TestBashNonZeroExitIsToolExecutionError(t *testing.T) {
	tool := NewBashTool(Deps{WorkspaceRoot: t.TempDir()})

	result, err := tool.Execute(context.Background(), mustJSON(t, BashInput{Command: exitCommand(7)}))
	if !toolruntime.IsCategory(err, toolruntime.ToolExecutionError) {
		t.Fatalf("expected ToolExecutionError, got %v", err)
	}
	if !result.IsError || result.Category != toolruntime.ToolExecutionError || result.Meta["exit_code"] != 7 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestBashTimeoutReturnsTimeoutError(t *testing.T) {
	tool := NewBashTool(Deps{
		WorkspaceRoot: t.TempDir(),
		Limits: Limits{
			BashTimeout:    10 * time.Millisecond,
			BashMaxTimeout: 20 * time.Millisecond,
			BashHeadBytes:  32,
			BashTailBytes:  32,
		},
	})

	result, err := tool.Execute(context.Background(), mustJSON(t, BashInput{Command: sleepCommand(), TimeoutMs: 10}))
	if !toolruntime.IsCategory(err, toolruntime.TimeoutError) {
		t.Fatalf("expected TimeoutError, got %v", err)
	}
	if !result.IsError || result.Category != toolruntime.TimeoutError {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestBashTruncatesOutputHeadTail(t *testing.T) {
	tool := NewBashTool(Deps{
		WorkspaceRoot: t.TempDir(),
		Limits:        Limits{BashHeadBytes: 5, BashTailBytes: 5},
	})

	result, err := tool.Execute(context.Background(), mustJSON(t, BashInput{Command: longOutputCommand()}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.Truncated {
		t.Fatal("expected truncated output")
	}
	if !strings.Contains(result.Output, "truncated") {
		t.Fatalf("expected truncation marker, got %q", result.Output)
	}
}

func exitCommand(code int) string {
	if runtime.GOOS == "windows" {
		return "exit /b " + strconv.Itoa(code)
	}
	return "exit " + strconv.Itoa(code)
}

func sleepCommand() string {
	if runtime.GOOS == "windows" {
		return "ping -n 3 127.0.0.1 > nul"
	}
	return "sleep 2"
}

func longOutputCommand() string {
	if runtime.GOOS == "windows" {
		return "echo 1234567890abcdefghijklmnopqrstuvwxyz"
	}
	return "printf '1234567890abcdefghijklmnopqrstuvwxyz'"
}
