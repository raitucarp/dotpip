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

type ExpireCommand struct {
	NX bool
	XX bool
	GT bool
	LT bool
}

type ExpireOption func(*ExpireCommand)

func WithExpireNX() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.NX = true
	}
}

func WithExpireXX() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.XX = true
	}
}

func WithExpireGT() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.GT = true
	}
}

func WithExpireLT() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.LT = true
	}
}
