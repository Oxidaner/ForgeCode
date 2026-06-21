package permission

import (
	"encoding/json"
	"strings"
	"time"
)

func attachApproval(req DecisionRequest, decision Decision) Decision {
	if decision.Effect != AskOnce && decision.Effect != AskAlways {
		decision.Approval = nil
		return decision
	}
	decision.Approval = &ApprovalRequest{
		SessionID:   req.Inv.SessionID,
		AgentID:     req.Inv.AgentID,
		TeamID:      req.Inv.TeamID,
		Source:      req.Inv.Source,
		ToolName:    req.Descriptor.Name,
		ToolSource:  req.Descriptor.Source,
		Input:       redactApprovalInput(req.Input),
		Effect:      decision.Effect,
		Risk:        decision.Risk,
		Reasons:     append([]RuleHit{}, decision.Reasons...),
		RequestedAt: time.Now().UTC(),
	}
	return decision
}

func redactApprovalInput(raw json.RawMessage) json.RawMessage {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return append(json.RawMessage{}, raw...)
	}
	redacted := redactApprovalValue(value, "")
	encoded, err := json.Marshal(redacted)
	if err != nil {
		return json.RawMessage(`{"redaction_error":true}`)
	}
	return encoded
}

func redactApprovalValue(value any, key string) any {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for childKey, child := range typed {
			out[childKey] = redactApprovalValue(child, childKey)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, child := range typed {
			out[i] = redactApprovalValue(child, key)
		}
		return out
	case string:
		if isSensitiveApprovalKey(key) {
			return "[REDACTED]"
		}
		return typed
	default:
		return typed
	}
}

func isSensitiveApprovalKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
	for _, part := range []string{"api_key", "apikey", "token", "password", "passwd", "secret", "authorization", "credential"} {
		if strings.Contains(normalized, part) {
			return true
		}
	}
	return false
}
