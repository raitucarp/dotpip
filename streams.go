package dotpip

// StreamEntry represents a single entry in a stream.
type StreamEntry struct {
	ID     string            `json:"id" yaml:"id" toml:"id"`
	Values map[string]string `json:"values" yaml:"values" toml:"values"`
}

// StreamPendingEntry represents a pending entry in a stream.
type StreamPendingEntry struct {
	Consumer      string `json:"consumer" yaml:"consumer" toml:"consumer"`
	DeliveryTime  int64  `json:"delivery_time" yaml:"delivery_time" toml:"delivery_time"`
	DeliveryCount int    `json:"delivery_count" yaml:"delivery_count" toml:"delivery_count"`
}

// StreamConsumer represents a consumer in a stream group.
type StreamConsumer struct {
	Pending map[string]int64 `json:"pending" yaml:"pending" toml:"pending"`
}

// StreamGroup represents a consumer group in a stream.
type StreamGroup struct {
	LastDeliveredID string                        `json:"last_delivered_id" yaml:"last_delivered_id" toml:"last_delivered_id"`
	Pending         map[string]StreamPendingEntry `json:"pending" yaml:"pending" toml:"pending"`
	Consumers       map[string]StreamConsumer     `json:"consumers" yaml:"consumers" toml:"consumers"`
}

// Stream represents a stream object.
type Stream struct {
	Entries []StreamEntry          `json:"entries" yaml:"entries" toml:"entries"`
	Groups  map[string]StreamGroup `json:"groups" yaml:"groups" toml:"groups"`
}

// XAddCommand options.
type XAddCommand struct {
	NoMkStream bool
	MaxLen     int
	MinID      string
	Limit      int
	Approx     bool
}

// XAddOption configures an XAddCommand.
type XAddOption func(*XAddCommand)

// WithXAddNoMkStream configures NOMKSTREAM for XADD.
func WithXAddNoMkStream() XAddOption {
	return func(cmd *XAddCommand) {
		cmd.NoMkStream = true
	}
}

// WithXAddMaxLen configures MAXLEN for XADD.
func WithXAddMaxLen(maxLen int, approx bool) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.MaxLen = maxLen
		cmd.Approx = approx
	}
}

// WithXAddMinID configures MINID for XADD.
func WithXAddMinID(minID string, approx bool) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.MinID = minID
		cmd.Approx = approx
	}
}

// WithXAddLimit configures LIMIT for XADD.
func WithXAddLimit(limit int) XAddOption {
	return func(cmd *XAddCommand) {
		cmd.Limit = limit
	}
}

// XTrimCommand options.
type XTrimCommand struct {
	MaxLen int
	MinID  string
	Limit  int
	Approx bool
}

// XTrimOption configures an XTrimCommand.
type XTrimOption func(*XTrimCommand)

// WithXTrimMaxLen configures MAXLEN for XTRIM.
func WithXTrimMaxLen(maxLen int, approx bool) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.MaxLen = maxLen
		cmd.Approx = approx
	}
}

// WithXTrimMinID configures MINID for XTRIM.
func WithXTrimMinID(minID string, approx bool) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.MinID = minID
		cmd.Approx = approx
	}
}

// WithXTrimLimit configures LIMIT for XTRIM.
func WithXTrimLimit(limit int) XTrimOption {
	return func(cmd *XTrimCommand) {
		cmd.Limit = limit
	}
}

// XReadCommand options.
type XReadCommand struct {
	Count int
	Block int // milliseconds
}

// XReadOption configures an XReadCommand.
type XReadOption func(*XReadCommand)

// WithXReadCount configures COUNT for XREAD.
func WithXReadCount(count int) XReadOption {
	return func(cmd *XReadCommand) {
		cmd.Count = count
	}
}

// WithXReadBlock configures BLOCK for XREAD.
func WithXReadBlock(block int) XReadOption {
	return func(cmd *XReadCommand) {
		cmd.Block = block
	}
}

// XReadGroupCommand options.
type XReadGroupCommand struct {
	Count int
	Block int // milliseconds
	NoAck bool
}

// XReadGroupOption configures an XReadGroupCommand.
type XReadGroupOption func(*XReadGroupCommand)

// WithXReadGroupCount configures COUNT for XREADGROUP.
func WithXReadGroupCount(count int) XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.Count = count
	}
}

// WithXReadGroupBlock configures blocking for XREADGROUP.
func WithXReadGroupBlock(block int) XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.Block = block
	}
}

// WithXReadGroupNoAck configures noack for XREADGROUP.
func WithXReadGroupNoAck() XReadGroupOption {
	return func(cmd *XReadGroupCommand) {
		cmd.NoAck = true
	}
}

// XPendingCommand options.
type XPendingCommand struct {
	Idle  int
	Start string
	End   string
	Count int
}

// XPendingOption configures an XPendingCommand.
type XPendingOption func(*XPendingCommand)

// WithXPendingIdle configures idle for XPENDING.
func WithXPendingIdle(idle int) XPendingOption {
	return func(cmd *XPendingCommand) {
		cmd.Idle = idle
	}
}

// WithXPendingRange configures range for XPENDING.
func WithXPendingRange(start, end string, count int) XPendingOption {
	return func(cmd *XPendingCommand) {
		cmd.Start = start
		cmd.End = end
		cmd.Count = count
	}
}

// XClaimCommand options.
type XClaimCommand struct {
	Idle       int
	Time       int64
	RetryCount int
	Force      bool
	JustID     bool
}

// XClaimOption configures an XClaimCommand.
type XClaimOption func(*XClaimCommand)

// WithXClaimIdle configures idle for XCLAIM.
func WithXClaimIdle(idle int) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Idle = idle
	}
}

// WithXClaimTime configures time for XCLAIM.
func WithXClaimTime(time int64) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Time = time
	}
}

// WithXClaimRetryCount configures retry count for XCLAIM.
func WithXClaimRetryCount(count int) XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.RetryCount = count
	}
}

// WithXClaimForce configures force for XCLAIM.
func WithXClaimForce() XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.Force = true
	}
}

// WithXClaimJustID configures just id for XCLAIM.
func WithXClaimJustID() XClaimOption {
	return func(cmd *XClaimCommand) {
		cmd.JustID = true
	}
}

// XAutoClaimCommand options.
type XAutoClaimCommand struct {
	Count  int
	JustID bool
}

// XAutoClaimOption configures an XAutoClaimCommand.
type XAutoClaimOption func(*XAutoClaimCommand)

// WithXAutoClaimCount configures count for XAUTOCLAIM.
func WithXAutoClaimCount(count int) XAutoClaimOption {
	return func(cmd *XAutoClaimCommand) {
		cmd.Count = count
	}
}

// WithXAutoClaimJustID configures just id for XAUTOCLAIM.
func WithXAutoClaimJustID() XAutoClaimOption {
	return func(cmd *XAutoClaimCommand) {
		cmd.JustID = true
	}
}
