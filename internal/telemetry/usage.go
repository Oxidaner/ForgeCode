package telemetry

import (
	"context"
	"sync"
	"time"
)

type UsageRecord struct {
	SessionID    string
	AgentID      string
	TeamID       string
	Model        string
	InputTokens  int64
	OutputTokens int64
	CostUSD      float64
	CreatedAt    time.Time
}

type UsageScope struct {
	SessionID string
	AgentID   string
	TeamID    string
}

type UsageSummary struct {
	InputTokens  int64
	OutputTokens int64
	CostUSD      float64
	Records      int
}

type UsageMeter interface {
	Record(ctx context.Context, u UsageRecord) error
	Aggregate(ctx context.Context, scope UsageScope) (UsageSummary, error)
}

type MemoryUsageMeter struct {
	mu      sync.RWMutex
	records []UsageRecord
}

func NewMemoryUsageMeter() *MemoryUsageMeter {
	return &MemoryUsageMeter{}
}

func (m *MemoryUsageMeter) Record(ctx context.Context, u UsageRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, u)
	return nil
}

func (m *MemoryUsageMeter) Aggregate(ctx context.Context, scope UsageScope) (UsageSummary, error) {
	if err := ctx.Err(); err != nil {
		return UsageSummary{}, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var summary UsageSummary
	for _, record := range m.records {
		if scope.SessionID != "" && record.SessionID != scope.SessionID {
			continue
		}
		if scope.AgentID != "" && record.AgentID != scope.AgentID {
			continue
		}
		if scope.TeamID != "" && record.TeamID != scope.TeamID {
			continue
		}
		summary.InputTokens += record.InputTokens
		summary.OutputTokens += record.OutputTokens
		summary.CostUSD += record.CostUSD
		summary.Records++
	}
	return summary, nil
}
