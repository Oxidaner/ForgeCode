package permission

import "context"

type StaticDecider struct {
	Decision Decision
}

func (d StaticDecider) Decide(ctx context.Context, req DecisionRequest) (Decision, error) {
	if err := ctx.Err(); err != nil {
		return Decision{}, err
	}
	decision := d.Decision
	if decision.Effect == "" {
		decision.Effect = Allow
	}
	if decision.Risk == "" {
		decision.Risk = RiskLow
	}
	for _, source := range req.Overrides {
		sourceDecision := DecisionFromPolicySource(source)
		if !IsMoreRestrictive(sourceDecision.Effect, decision.Effect) && sourceDecision.Effect == Allow {
			continue
		}
		decision = MergeDecisions(decision, sourceDecision)
	}
	return decision, nil
}
