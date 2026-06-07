package dotpip

type HRandFieldCommand struct {
	WithValues bool
}

type HRandFieldOption func(*HRandFieldCommand)

func WithHRandFieldWithValues() HRandFieldOption {
	return func(cmd *HRandFieldCommand) {
		cmd.WithValues = true
	}
}
