package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runJSONTests(t *testing.T, encoding fs.FileEncodeType) {
	tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()
	dotfs.EncodeType(encoding)

	key := dotpip.NewKey("myjson")

	// Set
	res, err := dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": []any{1, 2, 3}})
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), res)

	// Get
	val, err := dotfs.JSONGet(key, "$.b")
	assert.NoError(t, err)

	// Normalize val to []any with float64 to bypass decoder-specific number types
	b, _ := json.Marshal(val)
	var v []any
	_ = json.Unmarshal(b, &v)

	assert.Equal(t, []any{float64(1), float64(2), float64(3)}, v)
}

func TestJSON(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			runJSONTests(t, enc)
		})
	}
}

func TestJSONGet(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")

			// Set
			res, err := dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": []any{1, 2, 3}})
			assert.NoError(t, err)
			assert.Equal(t, string(dotpip.StatusOK), res)

			// Get Multiple Paths
			val, err := dotfs.JSONGet(key, "$.a", "$.b")
			assert.NoError(t, err)

			// bypass decoder diffs
			b, _ := json.Marshal(val)
			var v map[string]any
			_ = json.Unmarshal(b, &v)

			expected := map[string]any{"$.a": float64(1), "$.b": []any{float64(1), float64(2), float64(3)}}
			assert.Equal(t, expected, v)
		})
	}
}

func TestJSONArrAppend(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": []any{1}})

			res, err := dotfs.JSONArrAppend(key, "$.a", 2, 3)
			assert.NoError(t, err)
			assert.Equal(t, []any{3}, res)

			val, _ := dotfs.JSONGet(key, "$.a")
			b, _ := json.Marshal(val)
			var v []any
			_ = json.Unmarshal(b, &v)
			assert.Equal(t, []any{float64(1), float64(2), float64(3)}, v)
		})
	}
}

func TestJSONNumIncrBy(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1.5, "b": 10})

			res, err := dotfs.JSONNumIncrBy(key, "$.a", 2.0)
			assert.NoError(t, err)
			assert.Equal(t, []any{3.5}, res)

			res, err = dotfs.JSONNumIncrBy(key, "$.b", -5)
			assert.NoError(t, err)
			assert.Equal(t, []any{float64(5)}, res)
		})
	}
}

func TestJSONObjKeys(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": 3}})

			res, err := dotfs.JSONObjKeys(key, "$.b")
			assert.NoError(t, err)

			// Can be ["c", "d"] or ["d", "c"]
			keys := res[0].([]string)
			assert.Contains(t, keys, "c")
			assert.Contains(t, keys, "d")
		})
	}
}

func TestJSONMerge(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": 3}})

			// Update 'c' to 4, add 'e': 5, remove 'd' by setting it to nil
			res, err := dotfs.JSONMerge(key, "$.b", map[string]any{"c": 4, "e": 5, "d": nil})
			assert.NoError(t, err)
			assert.Equal(t, string(dotpip.StatusOK), res)

			val, _ := dotfs.JSONGet(key, "$.b")
			b, _ := json.Marshal(val)
			var v map[string]any
			_ = json.Unmarshal(b, &v)

			expected := map[string]any{"c": float64(4), "e": float64(5)}
			assert.Equal(t, expected, v)
		})
	}
}

func TestJSONNumMultBy(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1.5, "b": 10})

			res, err := dotfs.JSONNumMultBy(key, "$.a", 2.0)
			assert.NoError(t, err)
			assert.Equal(t, []any{3.0}, res)

			res, err = dotfs.JSONNumMultBy(key, "$.b", -5)
			assert.NoError(t, err)
			assert.Equal(t, []any{float64(-50)}, res)
		})
	}
}

func TestJSONType(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": "str", "c": false})

			res, err := dotfs.JSONType(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, []any{"number"}, res)

			res, err = dotfs.JSONType(key, "$.b")
			assert.NoError(t, err)
			assert.Equal(t, []any{"string"}, res)

			res, err = dotfs.JSONType(key, "$.c")
			assert.NoError(t, err)
			assert.Equal(t, []any{"boolean"}, res)
		})
	}
}

func TestJSONClear(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": []any{1, 2}, "c": map[string]any{"d": 4}})

			res, err := dotfs.JSONClear(key, "$.b")
			assert.NoError(t, err)
			assert.Equal(t, 1, res)

			val, _ := dotfs.JSONGet(key, "$.b")
			assert.Equal(t, []any{}, val)

			res, err = dotfs.JSONClear(key, "$.c")
			assert.NoError(t, err)
			assert.Equal(t, 1, res)

			val, _ = dotfs.JSONGet(key, "$.c")
			assert.Equal(t, map[string]any{}, val)
		})
	}
}

func TestJSONDel(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": 2})

			res, err := dotfs.JSONDel(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, 1, res)

			val, _ := dotfs.JSONGet(key, "$")

			b, _ := json.Marshal(val)
			var v map[string]any
			_ = json.Unmarshal(b, &v)

			expected := map[string]any{"b": float64(2)}
			assert.Equal(t, expected, v)
		})
	}
}

// The generic commands logic is all covered by TestJSON. We don't need additional YAML/TOML wrappers in test yet.

func TestJSONStrAppendAndLen(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": "hello"})

			res, err := dotfs.JSONStrAppend(key, "$.a", " world")
			assert.NoError(t, err)
			assert.Equal(t, []any{11}, res)

			resLen, err := dotfs.JSONStrLen(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, []any{11}, resLen)
		})
	}
}

func TestJSONToggle(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": true})

			res, err := dotfs.JSONToggle(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, []any{false}, res)

			res, err = dotfs.JSONToggle(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, []any{true}, res)
		})
	}
}

func TestJSONArrIndexLenPopTrim(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": []any{1, 2, 3, 4, 5}})

			res, err := dotfs.JSONArrLen(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, []any{5}, res)

			res, err = dotfs.JSONArrIndex(key, "$.a", float64(3)) // use float64 for generic number match
			assert.NoError(t, err)
			assert.Equal(t, []any{2}, res)

			res, err = dotfs.JSONArrInsert(key, "$.a", 1, 99)
			assert.NoError(t, err)
			assert.Equal(t, []any{6}, res)

			res, err = dotfs.JSONArrPop(key, "$.a", 1) // pop the 99
			assert.NoError(t, err)

			b, _ := json.Marshal(res)
			var v []any
			_ = json.Unmarshal(b, &v)
			assert.Equal(t, []any{float64(99)}, v)

			res, err = dotfs.JSONArrTrim(key, "$.a", 1, 2) // [1, 2, 3, 4, 5] -> trims to [2, 3] -> length 2
			assert.NoError(t, err)
			assert.Equal(t, []any{2}, res)

			val, _ := dotfs.JSONGet(key, "$.a")
			b, _ = json.Marshal(val)
			var v2 []any
			_ = json.Unmarshal(b, &v2)
			assert.Equal(t, []any{float64(2), float64(3)}, v2)
		})
	}
}

func TestJSONMGetMSet(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			k1 := dotpip.NewKey("k1")
			k2 := dotpip.NewKey("k2")

			_, _ = dotfs.JSONMSet(
				dotpip.JSONMSetArg{Key: k1, Path: "$", Value: map[string]any{"a": 1}},
				dotpip.JSONMSetArg{Key: k2, Path: "$", Value: map[string]any{"b": 2}},
			)

			res, err := dotfs.JSONMGet("$.a", k1, k2)
			assert.NoError(t, err)

			b, _ := json.Marshal(res)
			var v []any
			_ = json.Unmarshal(b, &v)

			assert.Equal(t, []any{float64(1), nil}, v)
		})
	}
}
