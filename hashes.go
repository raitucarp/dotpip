package dotpip

// HRandFieldCommand represents options for HRANDFIELD.
type HRandFieldCommand struct {
	WithValues bool
}

// HRandFieldOption configures an HRandFieldCommand.
type HRandFieldOption func(*HRandFieldCommand)

// WithHRandFieldWithValues configures WITHVALUES for HRANDFIELD.
func WithHRandFieldWithValues() HRandFieldOption {
	return func(cmd *HRandFieldCommand) {
		cmd.WithValues = true
	}
}
