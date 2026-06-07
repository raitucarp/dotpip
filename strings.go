package dotpip

// SetCommand represents options for SET.
type SetCommand struct {
	NX      bool
	XX      bool
	IfEq    string
	IfNe    string
	IfDeq   string
	IfDne   string
	Get     bool
	Ex      int
	Px      int
	ExAt    int
	PxAt    int
	KeepTTL bool
}

// SetOption configures a SetCommand.
type SetOption func(*SetCommand)

// WithNX configures NX for SET.
func WithNX() SetOption {
	return func(cmd *SetCommand) {
		cmd.NX = true
	}
}

// WithXX configures XX for SET.
func WithXX() SetOption {
	return func(cmd *SetCommand) {
		cmd.XX = true
	}
}

// WithIfEq configures IFEQ for SET.
func WithIfEq(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfEq = value
	}
}

// WithIfNe configures IFNE for SET.
func WithIfNe(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfNe = value
	}
}

// WithIfDeq configures IFDEQ for SET.
func WithIfDeq(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfDeq = value
	}
}

// WithIfDne configures IFDNE for SET.
func WithIfDne(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfDne = value
	}
}

// WithGet configures GET for SET.
func WithGet() SetOption {
	return func(cmd *SetCommand) {
		cmd.Get = true
	}
}

// WithEx configures EX for SET.
func WithEx(seconds int) SetOption {
	return func(cmd *SetCommand) {
		cmd.Ex = seconds
	}
}

// WithPx configures PX for SET.
func WithPx(milliseconds int) SetOption {
	return func(cmd *SetCommand) {
		cmd.Px = milliseconds
	}
}

// WithExAt configures EXAT for SET.
func WithExAt(timestamp int) SetOption {
	return func(cmd *SetCommand) {
		cmd.ExAt = timestamp
	}
}

// WithPxAt configures PXAT for SET.
func WithPxAt(timestamp int) SetOption {
	return func(cmd *SetCommand) {
		cmd.PxAt = timestamp
	}
}

// WithKeepTTL configures KEEPTTL for SET.
func WithKeepTTL() SetOption {
	return func(cmd *SetCommand) {
		cmd.KeepTTL = true
	}
}
