package builtin

import (
	"path/filepath"

	toolruntime "forgecode/internal/tool-runtime"
)

func resolveWorkspacePath(root, requested string) (string, error) {
	if requested == "" {
		return "", toolruntime.NewError(toolruntime.ValidationError, "path is required")
	}
	if filepath.IsAbs(requested) {
		return filepath.Clean(requested), nil
	}
	base := root
	if base == "" {
		base = "."
	}
	return filepath.Clean(filepath.Join(base, requested)), nil
}

func resolveOptionalWorkspacePath(root, requested string) string {
	if requested == "" {
		if root == "" {
			return "."
		}
		return filepath.Clean(root)
	}
	if filepath.IsAbs(requested) {
		return filepath.Clean(requested)
	}
	base := root
	if base == "" {
		base = "."
	}
	return filepath.Clean(filepath.Join(base, requested))
}

func relativeToWorkspace(root, path string) string {
	if root == "" {
		return filepath.Clean(path)
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.ToSlash(rel)
}
