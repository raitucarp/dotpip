package dotpip

type CopyCommand struct {
	Destination DotPip
	Replace     bool
}

type CopyOption func(*CopyCommand)

func WithDestination(destination DotPip) CopyOption {
	return func(cmd *CopyCommand) {
		cmd.Destination = destination
	}
}

func WithReplace() CopyOption {
	return func(cmd *CopyCommand) {
		cmd.Replace = true
	}
}
