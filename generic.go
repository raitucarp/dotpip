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

type RestoreCommand struct {
	Replace  bool
	AbsTTL   bool
	IdleTime int
	Freq     int
}

type RestoreOption func(*RestoreCommand)

func WithRestoreReplace() RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.Replace = true
	}
}

func WithRestoreAbsTTL() RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.AbsTTL = true
	}
}

func WithRestoreIdleTime(idleTime int) RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.IdleTime = idleTime
	}
}

func WithRestoreFreq(freq int) RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.Freq = freq
	}
}

type ScanCommand struct {
	Match string
	Count int
	Type  string
}

type ScanOption func(*ScanCommand)

func WithScanMatch(pattern string) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Match = pattern
	}
}

func WithScanCount(count int) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Count = count
	}
}

func WithScanType(typ string) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Type = typ
	}
}

type MigrateCommand struct {
	Copy    bool
	Replace bool
	Keys    []Key
}

type MigrateOption func(*MigrateCommand)

func WithMigrateCopy() MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Copy = true
	}
}

func WithMigrateReplace() MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Replace = true
	}
}

func WithMigrateKeys(keys ...Key) MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Keys = keys
	}
}
