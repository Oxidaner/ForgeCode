package telemetry

import "sync"

type Metrics interface {
	Counter(name string, tags ...Tag) Counter
	Histogram(name string, tags ...Tag) Histogram
}

type Counter interface {
	Add(delta int64)
	Value() int64
}

type Histogram interface {
	Observe(value float64)
	Values() []float64
}

type MemoryMetrics struct {
	mu        sync.Mutex
	counters  map[string]*memoryCounter
	histogram map[string]*memoryHistogram
}

func NewMemoryMetrics() *MemoryMetrics {
	return &MemoryMetrics{
		counters:  make(map[string]*memoryCounter),
		histogram: make(map[string]*memoryHistogram),
	}
}

func (m *MemoryMetrics) Counter(name string, tags ...Tag) Counter {
	key := metricKey(name, tags)
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.counters[key]; ok {
		return c
	}
	c := &memoryCounter{}
	m.counters[key] = c
	return c
}

func (m *MemoryMetrics) Histogram(name string, tags ...Tag) Histogram {
	key := metricKey(name, tags)
	m.mu.Lock()
	defer m.mu.Unlock()
	if h, ok := m.histogram[key]; ok {
		return h
	}
	h := &memoryHistogram{}
	m.histogram[key] = h
	return h
}

type memoryCounter struct {
	mu    sync.Mutex
	value int64
}

func (c *memoryCounter) Add(delta int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

func (c *memoryCounter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type memoryHistogram struct {
	mu     sync.Mutex
	values []float64
}

func (h *memoryHistogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.values = append(h.values, value)
}

func (h *memoryHistogram) Values() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]float64{}, h.values...)
}

func metricKey(name string, tags []Tag) string {
	key := name
	for _, tag := range tags {
		key += "|" + tag.Key + "=" + tag.Value
	}
	return key
}
