package permission

import (
	"context"
	"encoding/json"

	toolruntime "forgecode/internal/tool-runtime"
)

type RiskLevel = toolruntime.RiskLevel

const (
	RiskLow      = toolruntime.RiskLow
	RiskMedium   = toolruntime.RiskMedium
	RiskHigh     = toolruntime.RiskHigh
	RiskCritical = toolruntime.RiskCritical
)

type Effect string

const (
	Allow     Effect = "Allow"
	AskOnce   Effect = "AskOnce"
	AskAlways Effect = "AskAlways"
	Deny      Effect = "Deny"
)

type Layer string

const (
	LayerUnknown Layer = "Unknown"
	LayerL1      Layer = "L1"
	LayerL2      Layer = "L2"
	LayerL3      Layer = "L3"
	LayerL4      Layer = "L4"
	LayerL5      Layer = "L5"
)

type RuleHit struct {
	Layer  Layer
	RuleID string
	Reason string
}

type PolicySource struct {
	Name     string
	Effect   Effect
	Risk     RiskLevel
	RuleHits []RuleHit
}

type InvocationContext = toolruntime.InvocationContext

type DecisionRequest struct {
	Descriptor toolruntime.ToolDescriptor
	Input      json.RawMessage
	Inv        InvocationContext
	Overrides  []PolicySource
}

type Decision struct {
	Effect  Effect
	Risk    RiskLevel
	Reasons []RuleHit
	Layer   Layer
}

type Decider interface {
	Decide(ctx context.Context, req DecisionRequest) (Decision, error)
}

type BashAnalysis struct {
	Programs            []string
	Pipes               bool
	Redirects           bool
	Subshells           bool
	CmdSubst            bool
	NetworkAccess       bool
	FileDeletion        bool
	ForcePush           bool
	Docker              bool
	Kubernetes          bool
	DBWrite             bool
	DownloadThenExec    bool
	PrivilegeEscalation bool
}

type BashAnalyzer interface {
	Analyze(cmd string) (BashAnalysis, error)
}
