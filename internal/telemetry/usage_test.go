package telemetry

import (
	"context"
	"testing"
)

func TestMemoryUsageMeterAggregatesByScope(t *testing.T) {
	meter := NewMemoryUsageMeter()
	ctx := context.Background()

	_ = meter.Record(ctx, UsageRecord{SessionID: "s1", AgentID: "a1", InputTokens: 10, OutputTokens: 5, CostUSD: 0.1})
	_ = meter.Record(ctx, UsageRecord{SessionID: "s1", AgentID: "a2", InputTokens: 7, OutputTokens: 3, CostUSD: 0.2})
	_ = meter.Record(ctx, UsageRecord{SessionID: "s2", AgentID: "a1", InputTokens: 100, OutputTokens: 50, CostUSD: 1})

	summary, err := meter.Aggregate(ctx, UsageScope{SessionID: "s1"})
	if err != nil {
		t.Fatal(err)
	}
	if summary.InputTokens != 17 || summary.OutputTokens != 8 || summary.Records != 2 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}
