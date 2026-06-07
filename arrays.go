// Package dotpip provides an interface and structs for handling redis-like commands.
package dotpip

import "errors"

var ErrArrayEmpty = errors.New("array is empty")

type ARIndexValue struct {
	Index int
	Value string
}

type ARGrepPredicate struct {
	Type  string
	Value string
}

type ARGrepOptions struct {
	And        bool
	Or         bool
	Limit      *int
	WithValues bool
	NoCase     bool
}
