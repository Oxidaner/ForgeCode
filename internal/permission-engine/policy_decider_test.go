package permission

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestPolicyDeciderRejectsInvalidInput(t *testing.T) {
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: t.TempDir()})
	descriptor := readDescriptor()

	cases := map[string]json.RawMessage{
		"invalid-json":     json.RawMessage(`{"path":`),
		"missing-required": json.RawMessage(`{}`),
		"unexpected-field": json.RawMessage(`{"path":"file.txt","extra":true}`),
		"null-byte":        json.RawMessage(`{"path":"bad\u0000file"}`),
		"wrong-type":       json.RawMessage(`{"path":1}`),
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := decider.Decide(context.Background(), DecisionRequest{Descriptor: descriptor, Input: input})
			if !toolruntime.IsCategory(err, toolruntime.ValidationError) {
				t.Fatalf("expected ValidationError, got %v", err)
			}
		})
	}
}

func TestPolicyDeciderDeniesPathTraversal(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "../outside.txt"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, decision, Deny, RiskCritical, LayerL2)
}

func TestPolicyDeciderDeniesSensitivePath(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o700); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, root, ".git/config", "secret")
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": ".git/config"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, decision, Deny, RiskCritical, LayerL2)
}

func TestPolicyDeciderDeniesSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeTestFile(t, outside, "secret.txt", "secret")
	linkPath := filepath.Join(root, "link.txt")
	if err := os.Symlink(filepath.Join(outside, "secret.txt"), linkPath); err != nil {
		t.Skipf("symlink creation is unavailable in this environment: %v", err)
	}
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "link.txt"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, decision, Deny, RiskCritical, LayerL2)
}

func TestPolicyDeciderDeniesWriteOutsideWritablePaths(t *testing.T) {
	root := t.TempDir()
	allowed := filepath.Join(root, "allowed")
	if err := os.MkdirAll(allowed, 0o700); err != nil {
		t.Fatal(err)
	}
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root, WritablePaths: []string{allowed}})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: writeDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "other/file.txt", "content": "x"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, decision, Deny, RiskCritical, LayerL2)
}

func TestPolicyDeciderMapsRiskLevels(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	cases := []struct {
		name       string
		descriptor toolruntime.ToolDescriptor
		want       Effect
	}{
		{name: "low", descriptor: readDescriptor(), want: Allow},
		{name: "medium", descriptor: writeDescriptor(), want: AskOnce},
		{name: "high", descriptor: bashDescriptor(), want: AskAlways},
		{name: "critical", descriptor: criticalDescriptor(), want: Deny},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			decision, err := decider.Decide(context.Background(), DecisionRequest{
				Descriptor: tc.descriptor,
				Input:      validInputFor(tc.descriptor),
			})
			if err != nil {
				t.Fatal(err)
			}
			if decision.Effect != tc.want {
				t.Fatalf("effect = %s, want %s", decision.Effect, tc.want)
			}
		})
	}
}

func TestPolicyDeciderSupportsToolAndPathOverrides(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{
		WorkspaceRoot: root,
		ToolEffects:   map[string]Effect{"ReadFile": AskOnce},
		PathEffects: []PathPolicy{{
			Pattern: "*.secret",
			Effect:  Deny,
			Risk:    RiskCritical,
			Reason:  "secret suffix",
		}},
	})

	toolDecision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "notes.txt"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if toolDecision.Effect != AskOnce {
		t.Fatalf("tool override effect = %s, want AskOnce", toolDecision.Effect)
	}

	pathDecision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "token.secret"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, pathDecision, Deny, RiskCritical, LayerL3)
}

func TestPolicyOverridesCannotLoosenDecision(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: bashDescriptor(),
		Input:      mustRaw(t, map[string]any{"command": "echo hello"}),
		Overrides:  []PolicySource{{Name: "skill", Effect: Allow, Risk: RiskLow}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != AskAlways {
		t.Fatalf("allow override loosened decision to %s", decision.Effect)
	}
}

func TestPolicyOverridesCanOnlyTightenDecision(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "notes.txt"}),
		Overrides: []PolicySource{
			{Name: "skill", Effect: Allow, Risk: RiskLow},
			{Name: "hook", Effect: Deny, Risk: RiskCritical, RuleHits: []RuleHit{{Layer: LayerL5, RuleID: "hook-deny"}}},
			{Name: "user", Effect: AskOnce, Risk: RiskMedium, RuleHits: []RuleHit{{Layer: LayerL3, RuleID: "user-ask"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != Deny || decision.Risk != RiskCritical {
		t.Fatalf("decision = %#v, want strictest Deny Critical", decision)
	}
	if !hasRule(decision, "hook-deny") {
		t.Fatalf("decision reasons = %#v, want hook-deny", decision.Reasons)
	}
}

func TestPolicyDeciderImplementsToolRuntimePermissionChecker(t *testing.T) {
	var checker toolruntime.PermissionChecker = NewPolicyDecider(PolicyConfig{WorkspaceRoot: t.TempDir()})

	decision, err := checker.DecideTool(context.Background(), toolruntime.PermissionRequest{
		Descriptor: bashDescriptor(),
		Input:      mustRaw(t, map[string]any{"command": "git push --force origin main"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != toolruntime.PermissionAskAlways || decision.Risk != RiskHigh {
		t.Fatalf("decision = %#v, want AskAlways High", decision)
	}
	if len(decision.Reasons) == 0 || decision.Reasons[0].RuleID == "" {
		t.Fatalf("decision reasons = %#v, want rule IDs", decision.Reasons)
	}
}

func TestPolicyDeciderAddsApprovalRequestForAskDecision(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: sensitiveWriteDescriptor(),
		Input:      json.RawMessage(`{"path":"file.txt","api_key":"sk-secret"}`),
		Inv: toolruntime.InvocationContext{
			SessionID:     "session-1",
			AgentID:       "agent-1",
			TeamID:        "team-1",
			WorkspaceRoot: root,
			Source:        "model",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != AskOnce || decision.Approval == nil {
		t.Fatalf("decision = %#v, want AskOnce with approval request", decision)
	}
	approval := decision.Approval
	if approval.SessionID != "session-1" ||
		approval.AgentID != "agent-1" ||
		approval.TeamID != "team-1" ||
		approval.ToolName != "WriteSecret" ||
		approval.Effect != AskOnce ||
		approval.Risk != RiskMedium ||
		len(approval.Reasons) == 0 {
		t.Fatalf("approval = %#v", approval)
	}
	if strings.Contains(string(approval.Input), "sk-secret") {
		t.Fatalf("approval input was not redacted: %s", approval.Input)
	}
}

func TestPolicyDeciderDoesNotRequestApprovalForDeny(t *testing.T) {
	root := t.TempDir()
	decider := NewPolicyDecider(PolicyConfig{WorkspaceRoot: root})

	decision, err := decider.Decide(context.Background(), DecisionRequest{
		Descriptor: readDescriptor(),
		Input:      mustRaw(t, map[string]any{"path": "../outside.txt"}),
		Inv:        toolruntime.InvocationContext{WorkspaceRoot: root},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Effect != Deny || decision.Approval != nil {
		t.Fatalf("decision = %#v, want Deny without approval request", decision)
	}
}

func readDescriptor() toolruntime.ToolDescriptor {
	return toolruntime.ToolDescriptor{
		Name:        "ReadFile",
		Source:      toolruntime.ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object","required":["path"],"properties":{"path":{"type":"string","minLength":1}},"additionalProperties":false}`),
		Risk:        toolruntime.RiskLow,
		Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionRead}},
	}
}

func sensitiveWriteDescriptor() toolruntime.ToolDescriptor {
	return toolruntime.ToolDescriptor{
		Name:        "WriteSecret",
		Source:      toolruntime.ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object","required":["path","api_key"],"properties":{"path":{"type":"string","minLength":1},"api_key":{"type":"string","minLength":1}},"additionalProperties":false}`),
		Risk:        toolruntime.RiskMedium,
		Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionWrite}},
	}
}

func writeDescriptor() toolruntime.ToolDescriptor {
	return toolruntime.ToolDescriptor{
		Name:        "WriteFile",
		Source:      toolruntime.ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object","required":["path","content"],"properties":{"path":{"type":"string","minLength":1},"content":{"type":"string"}},"additionalProperties":false}`),
		Risk:        toolruntime.RiskMedium,
		Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionWrite}},
	}
}

func bashDescriptor() toolruntime.ToolDescriptor {
	return toolruntime.ToolDescriptor{
		Name:        "Bash",
		Source:      toolruntime.ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object","required":["command"],"properties":{"command":{"type":"string","minLength":1}},"additionalProperties":false}`),
		Risk:        toolruntime.RiskHigh,
		Permission:  toolruntime.PermissionHint{Actions: []toolruntime.PermissionAction{toolruntime.PermissionExecute}},
	}
}

func criticalDescriptor() toolruntime.ToolDescriptor {
	d := bashDescriptor()
	d.Name = "CriticalTool"
	d.Risk = toolruntime.RiskCritical
	return d
}

func validInputFor(descriptor toolruntime.ToolDescriptor) json.RawMessage {
	switch descriptor.Name {
	case "Bash", "CriticalTool":
		return json.RawMessage(`{"command":"echo hello"}`)
	case "WriteFile":
		return json.RawMessage(`{"path":"file.txt","content":"x"}`)
	default:
		return json.RawMessage(`{"path":"file.txt"}`)
	}
}

func mustRaw(t *testing.T, value any) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func writeTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertDecision(t *testing.T, decision Decision, effect Effect, risk RiskLevel, layer Layer) {
	t.Helper()
	if decision.Effect != effect || decision.Risk != risk || decision.Layer != layer {
		t.Fatalf("decision = %#v, want effect=%s risk=%s layer=%s", decision, effect, risk, layer)
	}
}
