package builtin

import (
	"fmt"

	toolruntime "forgecode/internal/tool-runtime"
)

func RegisterBuiltins(reg toolruntime.Registry, deps Deps) error {
	deps = deps.withDefaults()
	for _, tool := range []toolruntime.Tool{
		NewReadFileTool(deps),
		NewWriteFileTool(deps),
		NewEditFileTool(deps),
		NewBashTool(deps),
		NewGlobTool(deps),
		NewGrepTool(deps),
	} {
		if err := reg.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

func NewReadFileTool(d Deps) toolruntime.Tool {
	return readFileTool{descriptor: descriptorByName(ToolReadFile), deps: d.withDefaults()}
}

func NewWriteFileTool(d Deps) toolruntime.Tool {
	return writeFileTool{descriptor: descriptorByName(ToolWriteFile), deps: d.withDefaults()}
}

func NewEditFileTool(d Deps) toolruntime.Tool {
	return editFileTool{descriptor: descriptorByName(ToolEditFile), deps: d.withDefaults()}
}

func NewBashTool(d Deps) toolruntime.Tool {
	return bashTool{descriptor: descriptorByName(ToolBash), deps: d.withDefaults()}
}

func NewGlobTool(d Deps) toolruntime.Tool {
	return globTool{descriptor: descriptorByName(ToolGlob), deps: d.withDefaults()}
}

func NewGrepTool(d Deps) toolruntime.Tool {
	return grepTool{descriptor: descriptorByName(ToolGrep), deps: d.withDefaults()}
}

func descriptorByName(name string) toolruntime.ToolDescriptor {
	for _, descriptor := range Descriptors() {
		if descriptor.Name == name {
			return descriptor
		}
	}
	panic(fmt.Sprintf("missing built-in tool descriptor %q", name))
}
