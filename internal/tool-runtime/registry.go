package toolruntime

import (
	"sort"
	"sync"
)

type MemoryRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *MemoryRegistry {
	return &MemoryRegistry{tools: make(map[string]Tool)}
}

func (r *MemoryRegistry) Register(t Tool) error {
	if t == nil {
		return NewError(ValidationError, "tool is required")
	}
	descriptor := t.Descriptor()
	if err := descriptor.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[descriptor.Name]; exists {
		return NewError(ConflictError, "tool name already registered: "+descriptor.Name)
	}
	r.tools[descriptor.Name] = t
	return nil
}

func (r *MemoryRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

func (r *MemoryRegistry) List() []ToolDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptors := make([]ToolDescriptor, 0, len(r.tools))
	for _, tool := range r.tools {
		descriptors = append(descriptors, tool.Descriptor())
	}
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Name < descriptors[j].Name
	})
	return descriptors
}
