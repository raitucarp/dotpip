package dotpip

type LPosCommand struct {
	Rank   int
	Count  int
	MaxLen int
}

type LPosOption func(*LPosCommand)

func WithLPosRank(rank int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.Rank = rank
	}
}

func WithLPosCount(count int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.Count = count
	}
}

func WithLPosMaxLen(maxLen int) LPosOption {
	return func(cmd *LPosCommand) {
		cmd.MaxLen = maxLen
	}
}

type LInsertOption string

const (
	Before LInsertOption = "BEFORE"
	After  LInsertOption = "AFTER"
)

type LMoveDir string

const (
	Left  LMoveDir = "LEFT"
	Right LMoveDir = "RIGHT"
)
