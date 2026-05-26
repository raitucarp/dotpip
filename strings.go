package dotpip

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

type SetOption func(*SetCommand)

func WithNX() SetOption {
	return func(cmd *SetCommand) {
		cmd.NX = true
	}
}

func WithXX() SetOption {
	return func(cmd *SetCommand) {
		cmd.XX = true
	}
}

func WithIfEq(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfEq = value
	}
}

func WithIfNe(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfNe = value
	}
}

func WithIfDeq(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfDeq = value
	}
}

func WithIfDne(value string) SetOption {
	return func(cmd *SetCommand) {
		cmd.IfDne = value
	}
}

func WithGet() SetOption {
	return func(cmd *SetCommand) {
		cmd.Get = true
	}
}

func WithEx(seconds int) SetOption {
	return func(cmd *SetCommand) {
		cmd.Ex = seconds
	}
}

func WithPx(milliseconds int) SetOption {
	return func(cmd *SetCommand) {
		cmd.Px = milliseconds
	}
}

func WithExAt(timestamp int) SetOption {
	return func(cmd *SetCommand) {
		cmd.ExAt = timestamp
	}
}

func WithPxAt(timestamp int) SetOption {
	return func(cmd *SetCommand) {
		cmd.PxAt = timestamp
	}
}

func WithKeepTTL() SetOption {
	return func(cmd *SetCommand) {
		cmd.KeepTTL = true
	}
}
