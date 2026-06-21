package builtin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestReadFilePagesWithLineNumbers(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "notes.txt", "alpha\nbeta\ngamma\n")
	tool := NewReadFileTool(Deps{WorkspaceRoot: root, Limits: Limits{ReadDefaultLimit: 2}})

	result, err := tool.Execute(context.Background(), mustJSON(t, ReadFileInput{Path: "notes.txt", Offset: 2, Limit: 2}))
	if err != nil {
		t.Fatal(err)
	}
	if result.Output != "2: beta\n3: gamma\n" {
		t.Fatalf("unexpected output:\n%s", result.Output)
	}
	if result.Truncated {
		t.Fatal("did not expect truncated page")
	}
	if result.Meta["start_line"] != 2 || result.Meta["end_line"] != 3 || result.Meta["total_lines"] != 3 {
		t.Fatalf("unexpected metadata: %#v", result.Meta)
	}
}

func TestReadFileRejectsBinary(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "bin.dat"), []byte{'a', 0, 'b'}, 0o600); err != nil {
		t.Fatal(err)
	}
	tool := NewReadFileTool(Deps{WorkspaceRoot: root})

	_, err := tool.Execute(context.Background(), mustJSON(t, ReadFileInput{Path: "bin.dat"}))
	if !toolruntime.IsCategory(err, toolruntime.ValidationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestReadFileProtectsLargeFileWithoutPaging(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "large.txt", strings.Repeat("x", 20))
	tool := NewReadFileTool(Deps{WorkspaceRoot: root, Limits: Limits{ReadMaxBytes: 4}})

	_, err := tool.Execute(context.Background(), mustJSON(t, ReadFileInput{Path: "large.txt"}))
	if !toolruntime.IsCategory(err, toolruntime.ValidationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func mustJSON(t *testing.T, value any) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func writeTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
