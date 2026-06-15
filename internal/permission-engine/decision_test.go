package permission

import (
	"context"
	"reflect"
	"testing"
)

func TestMergeDecisionsUsesMostRestrictiveEffect(t *testing.T) {
	merged := MergeDecisions(
		Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL1},
		Decision{Effect: AskOnce, Risk: RiskMedium, Layer: LayerL3},
		Decision{Effect: AskAlways, Risk: RiskHigh, Layer: LayerL5},
		Decision{Effect: Deny, Risk: RiskCritical, Layer: LayerL2},
	)

	if merged.Effect != Deny {
		t.Fatalf("expected Deny to win, got %s", merged.Effect)
	}
	if merged.Risk != RiskCritical {
		t.Fatalf("expected Critical risk to win, got %s", merged.Risk)
	}
	if merged.Layer != LayerL2 {
		t.Fatalf("expected winning layer L2, got %s", merged.Layer)
	}
}

func TestStaticDeciderOverridesCanOnlyTighten(t *testing.T) {
	decider := StaticDecider{Decision: Decision{Effect: AskOnce, Risk: RiskMedium, Layer: LayerL3}}

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Overrides: []PolicySource{
			{Name: "skill", Effect: Allow, Risk: RiskLow},
			{Name: "hook", Effect: Deny, Risk: RiskHigh, RuleHits: []RuleHit{{Layer: LayerL5, RuleID: "hook-deny"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != Deny {
		t.Fatalf("expected deny override to tighten decision, got %s", decision.Effect)
	}
}

func TestDeciderInterfaceDoesNotExposeExecution(t *testing.T) {
	deciderType := reflect.TypeOf((*Decider)(nil)).Elem()
	if deciderType.NumMethod() != 1 || deciderType.Method(0).Name != "Decide" {
		t.Fatalf("Decider should only expose Decide, got %v methods", deciderType.NumMethod())
	}
}
