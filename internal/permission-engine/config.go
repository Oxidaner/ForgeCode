package permission

import toolruntime "forgecode/internal/tool-runtime"

type PolicyConfig struct {
	WorkspaceRoot  string
	WritablePaths  []string
	SensitivePaths []string
	RiskEffects    map[RiskLevel]Effect
	ToolEffects    map[string]Effect
	PathEffects    []PathPolicy
	BashAnalyzer   BashAnalyzer
}

type PathPolicy struct {
	Pattern string
	Effect  Effect
	Risk    RiskLevel
	Reason  string
}

func (c PolicyConfig) withDefaults(inv InvocationContext) PolicyConfig {
	if c.WorkspaceRoot == "" {
		c.WorkspaceRoot = inv.WorkspaceRoot
	}
	if len(c.WritablePaths) == 0 && c.WorkspaceRoot != "" {
		c.WritablePaths = []string{c.WorkspaceRoot}
	}
	if len(c.SensitivePaths) == 0 {
		c.SensitivePaths = []string{".git", ".env", ".ssh", "id_rsa", "id_ed25519"}
	}
	if c.RiskEffects == nil {
		c.RiskEffects = map[RiskLevel]Effect{
			toolruntime.RiskLow:      Allow,
			toolruntime.RiskMedium:   AskOnce,
			toolruntime.RiskHigh:     AskAlways,
			toolruntime.RiskCritical: Deny,
		}
	}
	if c.ToolEffects == nil {
		c.ToolEffects = map[string]Effect{}
	}
	if c.BashAnalyzer == nil {
		c.BashAnalyzer = NewBashAnalyzer()
	}
	return c
}
