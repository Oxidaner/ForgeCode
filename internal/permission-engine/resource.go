package permission

import (
	"os"
	"path/filepath"
	"strings"

	toolruntime "forgecode/internal/tool-runtime"
)

func evaluateResources(config PolicyConfig, req DecisionRequest, input map[string]any) (Decision, error) {
	paths := extractPathValues(input)
	if len(paths) == 0 {
		return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL2}, nil
	}

	for _, requested := range paths {
		decision, err := evaluatePath(config, req.Descriptor, requested)
		if err != nil {
			return Decision{}, err
		}
		if decision.Effect == Deny {
			return decision, nil
		}
	}
	return Decision{
		Effect: Allow,
		Risk:   RiskLow,
		Layer:  LayerL2,
		Reasons: []RuleHit{{
			Layer:  LayerL2,
			RuleID: "workspace-boundary",
			Reason: "all path inputs are within allowed workspace boundaries",
		}},
	}, nil
}

func evaluatePath(config PolicyConfig, descriptor toolruntime.ToolDescriptor, requested string) (Decision, error) {
	root, err := filepath.Abs(config.WorkspaceRoot)
	if err != nil {
		return Decision{}, toolruntime.WrapError(toolruntime.ValidationError, "resolve workspace root", err)
	}
	target := requested
	if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return Decision{}, toolruntime.WrapError(toolruntime.ValidationError, "resolve target path", err)
	}

	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return Decision{}, toolruntime.WrapError(toolruntime.ValidationError, "resolve workspace symlinks", err)
	}
	realTarget, err := resolveExistingOrParent(target)
	if err != nil {
		return Decision{}, err
	}
	if !isWithin(realRoot, realTarget) {
		return denyPath("path-escape", "path escapes workspace root")
	}
	if matchesSensitive(config, realRoot, realTarget) {
		return denyPath("sensitive-path", "path matches sensitive path policy")
	}
	if requiresWrite(descriptor) && !isWithinAny(config.WritablePaths, target) {
		return denyPath("write-boundary", "write path is outside writable paths")
	}
	return Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL2}, nil
}

func resolveExistingOrParent(target string) (string, error) {
	if resolved, err := filepath.EvalSymlinks(target); err == nil {
		return resolved, nil
	}
	parent := filepath.Dir(target)
	for parent != "." && parent != string(filepath.Separator) && parent != filepath.Dir(parent) {
		if info, err := os.Stat(parent); err == nil && info.IsDir() {
			resolved, err := filepath.EvalSymlinks(parent)
			if err != nil {
				return "", toolruntime.WrapError(toolruntime.ValidationError, "resolve parent symlinks", err)
			}
			return filepath.Join(resolved, filepath.Base(target)), nil
		}
		parent = filepath.Dir(parent)
	}
	return target, nil
}

func denyPath(ruleID, reason string) (Decision, error) {
	return Decision{
		Effect: Deny,
		Risk:   RiskCritical,
		Layer:  LayerL2,
		Reasons: []RuleHit{{
			Layer:  LayerL2,
			RuleID: ruleID,
			Reason: reason,
		}},
	}, nil
}

func extractPathValues(input map[string]any) []string {
	keys := map[string]bool{
		"path": true, "root": true, "file": true, "dir": true, "directory": true,
	}
	var paths []string
	var walk func(any, string)
	walk = func(value any, key string) {
		switch typed := value.(type) {
		case string:
			if keys[strings.ToLower(key)] && typed != "" {
				paths = append(paths, typed)
			}
		case map[string]any:
			for childKey, child := range typed {
				walk(child, childKey)
			}
		case []any:
			for _, child := range typed {
				walk(child, key)
			}
		}
	}
	walk(input, "")
	return paths
}

func requiresWrite(descriptor toolruntime.ToolDescriptor) bool {
	for _, action := range descriptor.Permission.Actions {
		if action == toolruntime.PermissionWrite {
			return true
		}
	}
	return false
}

func isWithinAny(roots []string, target string) bool {
	for _, root := range roots {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		if isWithin(absRoot, target) {
			return true
		}
	}
	return false
}

func isWithin(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel))
}

func matchesSensitive(config PolicyConfig, realRoot, realTarget string) bool {
	rel, err := filepath.Rel(realRoot, realTarget)
	if err != nil {
		return false
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	for _, sensitive := range config.SensitivePaths {
		sensitive = strings.Trim(strings.ToLower(filepath.ToSlash(sensitive)), "/")
		if sensitive == "" {
			continue
		}
		for _, part := range parts {
			if strings.ToLower(part) == sensitive {
				return true
			}
		}
		if strings.HasSuffix(strings.ToLower(filepath.ToSlash(rel)), sensitive) {
			return true
		}
	}
	return false
}

func evaluatePathPolicies(config PolicyConfig, input map[string]any) Decision {
	paths := extractPathValues(input)
	decision := Decision{Effect: Allow, Risk: RiskLow, Layer: LayerL3}
	for _, p := range paths {
		for _, policy := range config.PathEffects {
			ok, _ := filepath.Match(policy.Pattern, filepath.ToSlash(p))
			if !ok {
				continue
			}
			risk := policy.Risk
			if risk == "" {
				risk = RiskMedium
			}
			reason := policy.Reason
			if reason == "" {
				reason = "path-specific policy matched " + policy.Pattern
			}
			decision = MergeDecisions(decision, Decision{
				Effect: policy.Effect,
				Risk:   risk,
				Layer:  LayerL3,
				Reasons: []RuleHit{{
					Layer:  LayerL3,
					RuleID: "path-policy",
					Reason: reason,
				}},
			})
		}
	}
	return decision
}
