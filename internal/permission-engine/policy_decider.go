package permission

import (
	"context"

	toolruntime "forgecode/internal/tool-runtime"
)

type PolicyDecider struct {
	Config PolicyConfig
}

func NewPolicyDecider(config PolicyConfig) *PolicyDecider {
	return &PolicyDecider{Config: config}
}

func (d *PolicyDecider) Decide(ctx context.Context, req DecisionRequest) (Decision, error) {
	if err := ctx.Err(); err != nil {
		return Decision{}, err
	}

	config := d.Config.withDefaults(req.Inv)
	input, err := validateInput(req.Descriptor, req.Input)
	if err != nil {
		return Decision{}, err
	}

	resourceDecision, err := evaluateResources(config, req, input)
	if err != nil {
		return Decision{}, err
	}
	if resourceDecision.Effect == Deny {
		decision := mergeWithOverrides(resourceDecision, req.Overrides)
		return attachApproval(req, decision), nil
	}

	riskDecision, err := evaluateRisk(config, req.Descriptor, input)
	if err != nil {
		return Decision{}, err
	}
	decision := MergeDecisions(resourceDecision, riskDecision)
	decision = MergeDecisions(decision, evaluateToolPolicy(config, req.Descriptor))
	decision = MergeDecisions(decision, evaluatePathPolicies(config, input))
	decision = mergeWithOverrides(decision, req.Overrides)
	return attachApproval(req, decision), nil
}

func mergeWithOverrides(base Decision, overrides []PolicySource) Decision {
	decision := base
	for _, source := range overrides {
		sourceDecision := DecisionFromPolicySource(source)
		if !IsMoreRestrictive(sourceDecision.Effect, decision.Effect) && sourceDecision.Effect == Allow {
			continue
		}
		decision = MergeDecisions(decision, sourceDecision)
	}
	return decision
}

func evaluateRisk(config PolicyConfig, descriptor toolruntime.ToolDescriptor, input map[string]any) (Decision, error) {
	effect, ok := config.RiskEffects[descriptor.Risk]
	if !ok {
		effect = AskAlways
	}
	base := Decision{
		Effect: effect,
		Risk:   descriptor.Risk,
		Layer:  LayerL3,
		Reasons: []RuleHit{{
			Layer:  LayerL3,
			RuleID: "risk-policy",
			Reason: "tool default risk maps to " + string(effect),
		}},
	}
	bashDecision, err := evaluateBashRisk(config, descriptor, input)
	if err != nil {
		return Decision{}, err
	}
	return MergeDecisions(base, bashDecision), nil
}

func evaluateToolPolicy(config PolicyConfig, descriptor toolruntime.ToolDescriptor) Decision {
	effect, ok := config.ToolEffects[descriptor.Name]
	if !ok {
		return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL3}
	}
	return Decision{
		Effect: effect,
		Risk:   descriptor.Risk,
		Layer:  LayerL3,
		Reasons: []RuleHit{{
			Layer:  LayerL3,
			RuleID: "tool-policy",
			Reason: "tool-specific policy for " + descriptor.Name,
		}},
	}
}

func evaluateBashRisk(config PolicyConfig, descriptor toolruntime.ToolDescriptor, input map[string]any) (Decision, error) {
	if !isBashLike(descriptor) {
		return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL3}, nil
	}
	command, ok := input["command"].(string)
	if !ok {
		return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL3}, nil
	}
	analysis, err := config.BashAnalyzer.Analyze(command)
	if err != nil {
		return Decision{}, err
	}
	return decisionFromBashAnalysis(config, analysis), nil
}

func isBashLike(descriptor toolruntime.ToolDescriptor) bool {
	if descriptor.Name == "Bash" {
		return true
	}
	for _, action := range descriptor.Permission.Actions {
		if action == toolruntime.PermissionExecute {
			return true
		}
	}
	return false
}

func decisionFromBashAnalysis(config PolicyConfig, analysis BashAnalysis) Decision {
	decision := Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL3}
	add := func(risk RiskLevel, ruleID, reason string) {
		effect, ok := config.RiskEffects[risk]
		if !ok {
			effect = AskAlways
		}
		decision = MergeDecisions(decision, Decision{
			Effect: effect,
			Risk:   risk,
			Layer:  LayerL3,
			Reasons: []RuleHit{{
				Layer:  LayerL3,
				RuleID: ruleID,
				Reason: reason,
			}},
		})
	}

	if analysis.DownloadThenExec {
		add(RiskCritical, "bash-download-then-exec", "downloaded content is executed by another interpreter")
	}
	if analysis.PrivilegeEscalation {
		add(RiskCritical, "bash-privilege-escalation", "command attempts privilege escalation")
	}
	if analysis.FileDeletion {
		add(RiskHigh, "bash-file-deletion", "command deletes files or directories")
	}
	if analysis.ForcePush {
		add(RiskHigh, "bash-force-push", "git push uses force semantics")
	}
	if analysis.Docker {
		add(RiskHigh, "bash-docker", "command controls local container runtime")
	}
	if analysis.Kubernetes {
		add(RiskHigh, "bash-kubernetes", "command controls Kubernetes resources")
	}
	if analysis.DBWrite {
		add(RiskHigh, "bash-db-write", "database command contains write operation")
	}
	if analysis.NetworkAccess {
		add(RiskHigh, "bash-network-access", "command performs network access")
	}
	return decision
}
