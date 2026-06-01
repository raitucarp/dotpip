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

	HashEncode func(value map[string]string) (any, error)
	HashDecode func(value any) (map[string]string, error)

	ListEncode func(value []any) (any, error)
	ListDecode func(value any) ([]any, error)

	SetEncode func(value map[string]any) (any, error)
	SetDecode func(value any) (map[string]any, error)

	SortedSetEncode func(value map[string]float64) (any, error)
	SortedSetDecode func(value any) (map[string]float64, error)

	VectorSetDecode func(value any) ([]float64, error)
	VectorSetEncode func(value []float64) (any, error)

	StreamEncode func(value []any) (any, error)
	StreamDecode func(value any) ([]any, error)

	BitmapEncode func(value []uint) (any, error)
	BitmapDecode func(value any) ([]uint, error)

	BitfieldEncode func(value any) ([]any, error)
	BitfieldDecode func(value []any) (any, error)

	GeospatialEncode func(value any) ([]float64, error)
	GeospatialDecode func(value []float64) (any, error)

	JSONEncode func(value any) (any, error)
	JSONDecode func(value any) (any, error)
}

type KV struct {
	Key   Key
	Value string
}

type DotPip interface {
	Append(key Key, value string) int
	Get(key Key) (string, error)
	Set(key Key, value string, options ...SetOption) (string, error)
	Digest(key Key) (string, error)

	StrLen(key Key) int
	Incr(key Key) (int, error)
	IncrBy(key Key, increment int) (int, error)
	IncrByFloat(key Key, increment float64) (float64, error)
	Decr(key Key) (int, error)
	DecrBy(key Key, decrement int) (int, error)
	GetDel(key Key) (string, error)
	GetRange(key Key, start int, end int) (string, error)
	SetRange(key Key, offset int, value string) (int, error)
	MGet(keys ...Key) ([]string, error)
	MSet(kvs ...KV) error
	MSetNX(kvs ...KV) (bool, error)

	Del(keys ...Key) int
	Copy(source Key, destination Key, options ...CopyOption) int
	Exists(keys ...Key) ([]bool, error)
	FlushAll() error

	Formatter(fmap DataTypeFormatter)
}

func New(dotface DotPip) DotPip {
	return dotface
}
