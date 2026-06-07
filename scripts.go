package dotpip

// ScriptFlushOption configures a ScriptFlushCommand.
type ScriptFlushOption func(*ScriptFlushCommand)

// ScriptFlushCommand represents options for SCRIPT FLUSH.
type ScriptFlushCommand struct {
	Sync  bool
	Async bool
}

// WithScriptFlushSync configures SYNC for SCRIPT FLUSH.
func WithScriptFlushSync() ScriptFlushOption {
	return func(c *ScriptFlushCommand) {
		c.Sync = true
	}
}

// WithScriptFlushAsync configures ASYNC for SCRIPT FLUSH.
func WithScriptFlushAsync() ScriptFlushOption {
	return func(c *ScriptFlushCommand) {
		c.Async = true
	}
}
