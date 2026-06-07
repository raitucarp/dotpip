package dotpip

// LPosCommand represents options for LPOS.
type LPosCommand struct {
	Rank   int
	Count  int
	MaxLen int
}

// LPosOption configures LPosCommand.
type LPosOption func(*LPosCommand)

// WithLPosRank sets RANK for LPOS.
func WithLPosRank(rank int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.Rank = rank
	}
}

// WithLPosCount sets COUNT for LPOS.
func WithLPosCount(count int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.Count = count
	}
}

// WithLPosMaxLen sets MAXLEN for LPOS.
func WithLPosMaxLen(maxLen int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.MaxLen = maxLen
	}
}

// LInsertOption represents BEFORE or AFTER.
type LInsertOption string

const (
	// Before represents BEFORE.
	Before LInsertOption = "BEFORE"
	// After represents AFTER.
	After  LInsertOption = "AFTER"
)

// LMoveDir represents left or right direction.
type LMoveDir string

const (
	// Left represents LEFT.
	Left  LMoveDir = "LEFT"
	// Right represents RIGHT.
	Right LMoveDir = "RIGHT"
)
