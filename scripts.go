package dotpip

type ScriptFlushOption func(*ScriptFlushCommand)

type ScriptFlushCommand struct {
	Sync  bool
	Async bool
}

func WithScriptFlushSync() ScriptFlushOption {
	return func(c *ScriptFlushCommand) {
		c.Sync = true
	}
}

func WithScriptFlushAsync() ScriptFlushOption {
	return func(c *ScriptFlushCommand) {
		c.Async = true
	}
}
