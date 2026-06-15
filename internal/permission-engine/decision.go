package permission

import toolruntime "forgecode/internal/tool-runtime"

func NewAllowDecision(risk RiskLevel, reason RuleHit) Decision {
	return Decision{Effect: Allow, Risk: risk, Reasons: []RuleHit{reason}, Layer: reason.Layer}
}

func MergeDecisions(decisions ...Decision) Decision {
	if len(decisions) == 0 {
		return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerUnknown}
	}

	merged := decisions[0]
	for _, decision := range decisions[1:] {
		if effectRank(decision.Effect) > effectRank(merged.Effect) {
			merged.Effect = decision.Effect
			merged.Layer = decision.Layer
		}
		if riskRank(decision.Risk) > riskRank(merged.Risk) {
			merged.Risk = decision.Risk
		}
		merged.Reasons = append(merged.Reasons, decision.Reasons...)
	}
	return merged
}

func DecisionFromPolicySource(source PolicySource) Decision {
	layer := LayerUnknown
	if len(source.RuleHits) > 0 {
		layer = source.RuleHits[0].Layer
	}
	return Decision{
		Effect:  source.Effect,
		Risk:    source.Risk,
		Reasons: append([]RuleHit{}, source.RuleHits...),
		Layer:   layer,
	}
}

func IsMoreRestrictive(candidate, baseline Effect) bool {
	return effectRank(candidate) > effectRank(baseline)
}

func effectRank(effect Effect) int {
	switch effect {
	case Deny:
		return 4
	case AskAlways:
		return 3
	case AskOnce:
		return 2
	case Allow:
		return 1
	default:
		return 0
	}
}

func riskRank(risk RiskLevel) int {
	switch risk {
	case toolruntime.RiskCritical:
		return 4
	case toolruntime.RiskHigh:
		return 3
	case toolruntime.RiskMedium:
		return 2
	case toolruntime.RiskLow:
		return 1
	default:
		return 0
	}
}
