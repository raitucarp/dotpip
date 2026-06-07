// Package dotpip provides an interface and structs for handling redis-like commands.
package dotpip

import "errors"

// ErrArrayEmpty is returned when an array is empty.
// ErrArrayEmpty is returned when an array is empty.
var ErrArrayEmpty = errors.New("array is empty")

// ARIndexValue represents an indexed value in an array.
// ARIndexValue represents an indexed value in an array.
// ARIndexValue represents an indexed value in an array.
type ARIndexValue struct {
	Index int
	Value string
}

// ARGrepPredicate represents a predicate for ARGrep.
// ARGrepPredicate represents a predicate for ARGrep.
// ARGrepPredicate represents a predicate for ARGrep.
type ARGrepPredicate struct {
	Type  string
	Value string
}

// ARGrepOptions represents options for ARGrep.
// ARGrepOptions represents options for ARGrep.
// ARGrepOptions represents options for ARGrep.
type ARGrepOptions struct {
	And        bool
	Or         bool
	Limit      *int
	WithValues bool
	NoCase     bool
}
