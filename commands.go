package dotpip

import "strings"

type Key []string

func NewKey(s ...string) Key {
	return Key(s)
}

func NewKeyWithDelimiter(keyString string, delimiter string) Key {
	return Key(strings.Split(keyString, delimiter))
}

type DataTypeFormatter struct {
	StringEncode func(value string) (any, error)
	StringDecode func(value any) (string, error)

	Hash       func(value any) (map[string]string, error)
	List       func(value any) ([]any, error)
	Set        func(value any) (map[string]any, error)
	SortedSet  func(value any) (map[string]float64, error)
	VectorSet  func(value any) ([]float64, error)
	Stream     func(value any) ([]any, error)
	Bitmap     func(value any) ([]uint, error)
	Bitfield   func(value any) ([]any, error)
	Geospatial func(value any) ([]float64, error)
	JSON       func(value any) (any, error)
}

type DotPip interface {
	Get(key Key) (string, error)
	Set(key Key, value string) (string, error)

	FlushAll() error

	Formatter(fmap DataTypeFormatter)
}

func New(dotface DotPip) DotPip {
	return dotface
}
