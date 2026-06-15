package builtin

import (
	"context"
	"time"
)

type Checkpointer interface {
	CreateCheckpoint(ctx context.Context, sessionID, reason string) (checkpointID string, err error)
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

type Limits struct {
	ReadMaxBytes         int64
	ReadDefaultLimit     int
	ReadBinaryProbeBytes int
	BashTimeout          time.Duration
	BashMaxTimeout       time.Duration
	BashHeadBytes        int
	BashTailBytes        int
	GrepMaxMatches       int
	GlobMaxResults       int
}

type Deps struct {
	Checkpointer  Checkpointer
	WorkspaceRoot string
	Clock         Clock
	Limits        Limits
}

func (d Deps) withDefaults() Deps {
	if d.Clock == nil {
		d.Clock = RealClock{}
	}
	if d.Limits.ReadMaxBytes == 0 {
		d.Limits.ReadMaxBytes = 2 << 20
	}
	if d.Limits.ReadDefaultLimit == 0 {
		d.Limits.ReadDefaultLimit = 2000
	}
	if d.Limits.ReadBinaryProbeBytes == 0 {
		d.Limits.ReadBinaryProbeBytes = 8192
	}
	if d.Limits.BashTimeout == 0 {
		d.Limits.BashTimeout = 120 * time.Second
	}
	if d.Limits.BashMaxTimeout == 0 {
		d.Limits.BashMaxTimeout = 600 * time.Second
	}
	if d.Limits.BashHeadBytes == 0 {
		d.Limits.BashHeadBytes = 16 << 10
	}
	if d.Limits.BashTailBytes == 0 {
		d.Limits.BashTailBytes = 16 << 10
	}
	if d.Limits.GrepMaxMatches == 0 {
		d.Limits.GrepMaxMatches = 1000
	}
	if d.Limits.GlobMaxResults == 0 {
		d.Limits.GlobMaxResults = 1000
	}
	return d
}
