package dotpip

type Z struct {
	Score  float64
	Member string
}

type ZAddCommand struct {
	NX   bool
	XX   bool
	GT   bool
	LT   bool
	CH   bool
	INCR bool
}

type ZAddOption func(*ZAddCommand)

func WithZAddNX() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.NX = true
	}
}

func WithZAddXX() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.XX = true
	}
}

func WithZAddGT() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.GT = true
	}
}

func WithZAddLT() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.LT = true
	}
}

func WithZAddCH() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.CH = true
	}
}

func WithZAddINCR() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.INCR = true
	}
}

type ZRangeCommand struct {
	ByScore bool
	ByLex   bool
	Rev     bool
	Limit   bool
	Offset  int
	Count   int
}

type ZRangeOption func(*ZRangeCommand)

func WithZRangeByScore() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.ByScore = true
	}
}

func WithZRangeByLex() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.ByLex = true
	}
}

func WithZRangeRev() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.Rev = true
	}
}

func WithZRangeLimit(offset int, count int) ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.Limit = true
		cmd.Offset = offset
		cmd.Count = count
	}
}
