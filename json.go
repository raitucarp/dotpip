package dotpip

// JSONMSetArg represents arguments for JSON.MSET.
type JSONMSetArg struct {
	Key   Key
	Path  string
	Value any
}
