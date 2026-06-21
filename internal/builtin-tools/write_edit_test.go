package builtin

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestWriteFileCreatesNewFileWithoutCheckpoint(t *testing.T) {
	root := t.TempDir()
	checkpointer := &fakeCheckpointer{}
	tool := NewWriteFileTool(Deps{WorkspaceRoot: root, Checkpointer: checkpointer})

	result, err := tool.Execute(context.Background(), mustJSON(t, WriteFileInput{
		Path:    "notes.txt",
		Content: "hello\n",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, root, "notes.txt"); got != "hello\n" {
		t.Fatalf("file content = %q", got)
	}
	if checkpointer.calls != 0 {
		t.Fatalf("checkpoint calls = %d, want 0", checkpointer.calls)
	}
	if result.Meta["bytes"] != 6 || result.Meta["created"] != true {
		t.Fatalf("metadata = %#v", result.Meta)
	}
}

func TestWriteFileOverwritesWithCheckpoint(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "notes.txt", "old\n")
	checkpointer := &fakeCheckpointer{id: "cp-1"}
	tool := NewWriteFileTool(Deps{WorkspaceRoot: root, Checkpointer: checkpointer})

	result, err := tool.Execute(context.Background(), mustJSON(t, WriteFileInput{
		Path:    "notes.txt",
		Content: "new\n",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, root, "notes.txt"); got != "new\n" {
		t.Fatalf("file content = %q", got)
	}
	if checkpointer.calls != 1 || !strings.Contains(checkpointer.reason, "WriteFile") {
		t.Fatalf("checkpoint = %#v", checkpointer)
	}
	if result.Meta["checkpoint_id"] != "cp-1" || result.Meta["created"] != false {
		t.Fatalf("metadata = %#v", result.Meta)
	}
}

func TestWriteFileStopsWhenCheckpointFails(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "notes.txt", "old\n")
	tool := NewWriteFileTool(Deps{
		WorkspaceRoot: root,
		Checkpointer:  &fakeCheckpointer{err: errors.New("store down")},
	})

	_, err := tool.Execute(context.Background(), mustJSON(t, WriteFileInput{
		Path:    "notes.txt",
		Content: "new\n",
	}))
	if !toolruntime.IsCategory(err, toolruntime.PersistenceError) {
		t.Fatalf("expected PersistenceError, got %v", err)
	}
	if got := readTestFile(t, root, "notes.txt"); got != "old\n" {
		t.Fatalf("file changed after checkpoint failure: %q", got)
	}
}

func TestEditFileReplacesUniqueMatchAndReturnsDiff(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "notes.txt", "alpha\nold\nomega\n")
	checkpointer := &fakeCheckpointer{id: "cp-2"}
	tool := NewEditFileTool(Deps{WorkspaceRoot: root, Checkpointer: checkpointer})

	result, err := tool.Execute(context.Background(), mustJSON(t, EditFileInput{
		Path:      "notes.txt",
		OldString: "old\n",
		NewString: "new\n",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, root, "notes.txt"); got != "alpha\nnew\nomega\n" {
		t.Fatalf("file content = %q", got)
	}
	diff, ok := result.Meta["diff"].(string)
	if !ok || !strings.Contains(diff, "-old") || !strings.Contains(diff, "+new") {
		t.Fatalf("diff = %#v", result.Meta["diff"])
	}
	if checkpointer.calls != 1 || result.Meta["checkpoint_id"] != "cp-2" {
		t.Fatalf("checkpoint/meta = %#v %#v", checkpointer, result.Meta)
	}
}

func TestEditFileRejectsZeroAndMultipleMatches(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "notes.txt", "same\nsame\n")
	tool := NewEditFileTool(Deps{WorkspaceRoot: root, Checkpointer: &fakeCheckpointer{id: "cp"}})

	for _, input := range []EditFileInput{
		{Path: "notes.txt", OldString: "missing", NewString: "new"},
		{Path: "notes.txt", OldString: "same", NewString: "new"},
	} {
		_, err := tool.Execute(context.Background(), mustJSON(t, input))
		if !toolruntime.IsCategory(err, toolruntime.ValidationError) {
			t.Fatalf("input %#v expected ValidationError, got %v", input, err)
		}
	}
	if got := readTestFile(t, root, "notes.txt"); got != "same\nsame\n" {
		t.Fatalf("file changed after rejected edit: %q", got)
	}
}

type fakeCheckpointer struct {
	id      string
	err     error
	calls   int
	reason  string
	session string
}

func (f *fakeCheckpointer) CreateCheckpoint(ctx context.Context, sessionID, reason string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	f.calls++
	f.reason = reason
	f.session = sessionID
	if f.err != nil {
		return "", f.err
	}
	if f.id == "" {
		return "checkpoint-1", nil
	}
	return f.id, nil
}

func readTestFile(t *testing.T, root, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
