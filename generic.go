package dotpip

// CopyCommand represents options for COPY.
// CopyCommand represents options for COPY command.
type CopyCommand struct {
	Destination DotPip
	Replace     bool
}

// CopyOption configures CopyCommand.
// CopyOption configures CopyCommand.
type CopyOption func(*CopyCommand)

// WithDestination sets DESTINATION for COPY.
// WithDestination sets DESTINATION for COPY.
func WithDestination(destination DotPip) CopyOption {
	return func(cmd *CopyCommand) {
		cmd.Destination = destination
	}
}

// WithReplace sets REPLACE for COPY.
// WithReplace sets REPLACE for COPY.
func WithReplace() CopyOption {
	return func(cmd *CopyCommand) {
		cmd.Replace = true
	}
}

// ExpireCommand represents options for EXPIRE.
// ExpireCommand represents options for EXPIRE.
type ExpireCommand struct {
	NX bool
	XX bool
	GT bool
	LT bool
}

// ExpireOption configures ExpireCommand.
// ExpireOption configures ExpireCommand.
type ExpireOption func(*ExpireCommand)

// WithExpireNX sets NX for EXPIRE.
// WithExpireNX sets NX for EXPIRE.
func WithExpireNX() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.NX = true
	}
}

// WithExpireXX sets XX for EXPIRE.
// WithExpireXX sets XX for EXPIRE.
func WithExpireXX() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.XX = true
	}
}

// WithExpireGT sets GT for EXPIRE.
// WithExpireGT sets GT for EXPIRE.
func WithExpireGT() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.GT = true
	}
}

// WithExpireLT sets LT for EXPIRE.
// WithExpireLT sets LT for EXPIRE.
func WithExpireLT() ExpireOption {
	return func(cmd *ExpireCommand) {
		cmd.LT = true
	}
}

// RestoreCommand represents options for RESTORE.
// RestoreCommand represents options for RESTORE.
type RestoreCommand struct {
	Replace  bool
	AbsTTL   bool
	IdleTime int
	Freq     int
}

// RestoreOption configures RestoreCommand.
// RestoreOption configures RestoreCommand.
type RestoreOption func(*RestoreCommand)

// WithRestoreReplace sets REPLACE for RESTORE.
// WithRestoreReplace sets REPLACE for RESTORE.
func WithRestoreReplace() RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.Replace = true
	}
}

// WithRestoreAbsTTL sets ABSTTL for RESTORE.
// WithRestoreAbsTTL sets ABSTTL for RESTORE.
func WithRestoreAbsTTL() RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.AbsTTL = true
	}
}

// WithRestoreIdleTime sets IDLETIME for RESTORE.
// WithRestoreIdleTime sets IDLETIME for RESTORE.
func WithRestoreIdleTime(idleTime int) RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.IdleTime = idleTime
	}
}

// WithRestoreFreq sets FREQ for RESTORE.
// WithRestoreFreq sets FREQ for RESTORE.
func WithRestoreFreq(freq int) RestoreOption {
	return func(cmd *RestoreCommand) {
		cmd.Freq = freq
	}
}

// ScanCommand represents options for SCAN.
// ScanCommand represents options for SCAN.
type ScanCommand struct {
	Match string
	Count int
	Type  string
}

// ScanOption configures ScanCommand.
// ScanOption configures ScanCommand.
type ScanOption func(*ScanCommand)

// WithScanMatch sets MATCH for SCAN.
// WithScanMatch sets MATCH for SCAN.
func WithScanMatch(pattern string) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Match = pattern
	}
}

// WithScanCount sets COUNT for SCAN.
// WithScanCount sets COUNT for SCAN.
func WithScanCount(count int) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Count = count
	}
}

// WithScanType sets TYPE for SCAN.
// WithScanType sets TYPE for SCAN.
func WithScanType(typ string) ScanOption {
	return func(cmd *ScanCommand) {
		cmd.Type = typ
	}
}

// MigrateCommand represents options for MIGRATE.
// MigrateCommand represents options for MIGRATE.
type MigrateCommand struct {
	Copy    bool
	Replace bool
	Keys    []Key
}

// MigrateOption configures MigrateCommand.
// MigrateOption configures MigrateCommand.
type MigrateOption func(*MigrateCommand)

// WithMigrateCopy sets COPY for MIGRATE.
// WithMigrateCopy sets COPY for MIGRATE.
func WithMigrateCopy() MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Copy = true
	}
}

// WithMigrateReplace sets REPLACE for MIGRATE.
// WithMigrateReplace sets REPLACE for MIGRATE.
func WithMigrateReplace() MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Replace = true
	}
}

// WithMigrateKeys sets KEYS for MIGRATE.
// WithMigrateKeys sets KEYS for MIGRATE.
func WithMigrateKeys(keys ...Key) MigrateOption {
	return func(cmd *MigrateCommand) {
		cmd.Keys = keys
	}
}
