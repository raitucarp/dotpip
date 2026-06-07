package dotpip

// Z represents a sorted set member and its score.
type Z struct {
	Score  float64
	Member string
}

// ZAddCommand represents options for ZADD.
type ZAddCommand struct {
	NX   bool
	XX   bool
	GT   bool
	LT   bool
	CH   bool
	INCR bool
}

// ZAddOption configures a ZAddCommand.
type ZAddOption func(*ZAddCommand)

// WithZAddNX configures NX for ZADD.
func WithZAddNX() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.NX = true
	}
}

// WithZAddXX configures XX for ZADD.
func WithZAddXX() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.XX = true
	}
}

// WithZAddGT configures GT for ZADD.
func WithZAddGT() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.GT = true
	}
}

// WithZAddLT configures LT for ZADD.
func WithZAddLT() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.LT = true
	}
}

// WithZAddCH configures CH for ZADD.
func WithZAddCH() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.CH = true
	}
}

// WithZAddINCR configures INCR for ZADD.
func WithZAddINCR() ZAddOption {
	return func(cmd *ZAddCommand) {
		cmd.INCR = true
	}
}

// ZRangeCommand represents options for ZRANGE.
type ZRangeCommand struct {
	ByScore bool
	ByLex   bool
	Rev     bool
	Limit   bool
	Offset  int
	Count   int
}

// ZRangeOption configures a ZRangeCommand.
type ZRangeOption func(*ZRangeCommand)

// WithZRangeByScore configures BYSCORE for ZRANGE.
func WithZRangeByScore() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.ByScore = true
	}
}

// WithZRangeByLex configures BYLEX for ZRANGE.
func WithZRangeByLex() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.ByLex = true
	}
}

// WithZRangeRev configures REV for ZRANGE.
func WithZRangeRev() ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.Rev = true
	}
}

// WithZRangeLimit configures LIMIT for ZRANGE.
func WithZRangeLimit(offset int, count int) ZRangeOption {
	return func(cmd *ZRangeCommand) {
		cmd.Limit = true
		cmd.Offset = offset
		cmd.Count = count
	}
}
