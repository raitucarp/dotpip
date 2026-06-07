package dotpip

type StreamEntry struct {
	ID     string            `json:"id" yaml:"id" toml:"id"`
	Values map[string]string `json:"values" yaml:"values" toml:"values"`
}

type StreamPendingEntry struct {
	Consumer      string `json:"consumer" yaml:"consumer" toml:"consumer"`
	DeliveryTime  int64  `json:"delivery_time" yaml:"delivery_time" toml:"delivery_time"`
	DeliveryCount int    `json:"delivery_count" yaml:"delivery_count" toml:"delivery_count"`
}

type StreamConsumer struct {
	Pending map[string]int64 `json:"pending" yaml:"pending" toml:"pending"`
}

type StreamGroup struct {
	LastDeliveredID string                        `json:"last_delivered_id" yaml:"last_delivered_id" toml:"last_delivered_id"`
	Pending         map[string]StreamPendingEntry `json:"pending" yaml:"pending" toml:"pending"`
	Consumers       map[string]StreamConsumer     `json:"consumers" yaml:"consumers" toml:"consumers"`
}

type Stream struct {
	Entries []StreamEntry          `json:"entries" yaml:"entries" toml:"entries"`
	Groups  map[string]StreamGroup `json:"groups" yaml:"groups" toml:"groups"`
}

type XAddCommand struct {
	NoMkStream bool
	MaxLen     int
	MinID      string
	Limit      int
	Approx     bool
}

type XAddOption func(*XAddCommand)

func WithXAddNoMkStream() XAddOption {
	return func(cmd *XAddCommand) {
		cmd.NoMkStream = true
	}
}

func WithXAddMaxLen(maxLen int, approx bool) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.MaxLen = maxLen
		cmd.Approx = approx
	}
}

func WithXAddMinID(minID string, approx bool) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.MinID = minID
		cmd.Approx = approx
	}
}

func WithXAddLimit(limit int) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.Limit = limit
	}
}

type XTrimCommand struct {
	MaxLen int
	MinID  string
	Limit  int
	Approx bool
}

type XTrimOption func(*XTrimCommand)

func WithXTrimMaxLen(maxLen int, approx bool) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.MaxLen = maxLen
		cmd.Approx = approx
	}
}

func WithXTrimMinID(minID string, approx bool) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.MinID = minID
		cmd.Approx = approx
	}
}

func WithXTrimLimit(limit int) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.Limit = limit
	}
}

type XReadCommand struct {
	Count int
	Block int // milliseconds
}

type XReadOption func(*XReadCommand)

func WithXReadCount(count int) XReadOption {
	return func(cmd *XReadCommand) {
		cmd.Count = count
	}
}

func WithXReadBlock(block int) XReadOption {
	return func(cmd *XReadCommand) {
		cmd.Block = block
	}
}

type XReadGroupCommand struct {
	Count int
	Block int // milliseconds
	NoAck bool
}

type XReadGroupOption func(*XReadGroupCommand)

func WithXReadGroupCount(count int) XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.Count = count
	}
}

func WithXReadGroupBlock(block int) XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.Block = block
	}
}

func WithXReadGroupNoAck() XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.NoAck = true
	}
}

type XPendingCommand struct {
	Idle  int
	Start string
	End   string
	Count int
}

type XPendingOption func(*XPendingCommand)

func WithXPendingIdle(idle int) XPendingOption {
	return func(cmd *XPendingCommand) {
		cmd.Idle = idle
	}
}

func WithXPendingRange(start, end string, count int) XPendingOption {
	return func(cmd *XPendingCommand) {
		cmd.Start = start
		cmd.End = end
		cmd.Count = count
	}
}

type XClaimCommand struct {
	Idle       int
	Time       int64
	RetryCount int
	Force      bool
	JustID     bool
}

type XClaimOption func(*XClaimCommand)

func WithXClaimIdle(idle int) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Idle = idle
	}
}

func WithXClaimTime(time int64) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Time = time
	}
}

func WithXClaimRetryCount(count int) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.RetryCount = count
	}
}

func WithXClaimForce() XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Force = true
	}
}

func WithXClaimJustID() XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.JustID = true
	}
}

type XAutoClaimCommand struct {
	Count  int
	JustID bool
}

type XAutoClaimOption func(*XAutoClaimCommand)

func WithXAutoClaimCount(count int) XAutoClaimOption {
	return func(cmd *XAutoClaimCommand) {
		cmd.Count = count
	}
}

func WithXAutoClaimJustID() XAutoClaimOption {
	return func(cmd *XAutoClaimCommand) {
		cmd.JustID = true
	}
}
