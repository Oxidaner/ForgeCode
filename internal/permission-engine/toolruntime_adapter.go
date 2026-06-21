package permission

import (
	"context"

	toolruntime "forgecode/internal/tool-runtime"
)

func (d *PolicyDecider) DecideTool(ctx context.Context, req toolruntime.PermissionRequest) (toolruntime.PermissionDecision, error) {
	decision, err := d.Decide(ctx, DecisionRequest{
		Descriptor: req.Descriptor,
		Input:      req.Input,
		Inv:        req.Inv,
	})
	if err != nil {
		return toolruntime.PermissionDecision{}, err
	}
	return toToolRuntimeDecision(decision), nil
}

func toToolRuntimeDecision(decision Decision) toolruntime.PermissionDecision {
	reasons := make([]toolruntime.PermissionReason, 0, len(decision.Reasons))
	for _, reason := range decision.Reasons {
		reasons = append(reasons, toolruntime.PermissionReason{
			Layer:  string(reason.Layer),
			RuleID: reason.RuleID,
			Reason: reason.Reason,
		})
	}
	return toolruntime.PermissionDecision{
		Effect:  toToolRuntimeEffect(decision.Effect),
		Risk:    decision.Risk,
		Reasons: reasons,
	}
}

func toToolRuntimeEffect(effect Effect) toolruntime.PermissionEffect {
	switch effect {
	case Allow:
		return toolruntime.PermissionAllow
	case AskOnce:
		return toolruntime.PermissionAskOnce
	case AskAlways:
		return toolruntime.PermissionAskAlways
	case Deny:
		return toolruntime.PermissionDeny
	default:
		return ""
	}
}
