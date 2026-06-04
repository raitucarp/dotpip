package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONMore(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_more_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")

			// Forget
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1})
			resForget, err := dotfs.JSONForget(key, "$.a")
			assert.NoError(t, err)
			assert.Equal(t, 1, resForget)

			// Resp
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": "str", "c": []any{1}})
			resResp, err := dotfs.JSONResp(key, "$.b")
			assert.NoError(t, err)
			_ = resResp

			// Debug
			resDebug, err := dotfs.JSONDebug("MEMORY", key, "$")
			assert.NoError(t, err)
			assert.NotNil(t, resDebug)

			resDebug2, err := dotfs.JSONDebug("HELP", key, "$")
			assert.NoError(t, err)
			assert.Equal(t, "OK", resDebug2)
		})
	}
}

func TestJSONObjLenAndDel(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_more2_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson2")

			// JSONObjLen
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": map[string]any{"x": 1, "y": 2}})
			resLen, err := dotfs.JSONObjLen(key, "$.a")
			assert.NoError(t, err)

			// The current code might return []any depending on match.
			assert.NotNil(t, resLen)

			// JSONDel
			resDel, err := dotfs.JSONDel(key, "$.a.x")
			assert.NoError(t, err)
			assert.Equal(t, 1, resDel)

			resDel2, err := dotfs.JSONDel(key, "$")
			assert.NoError(t, err)
			assert.Equal(t, 1, resDel2)
		})
	}
}

func TestJSONTypeGet(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_type_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson3")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": "str", "c": []any{1}})

			// Type
			resType, err := dotfs.JSONType(key, "$.c")
			assert.NoError(t, err)
			assert.NotNil(t, resType)

			resType2, err := dotfs.JSONType(key, "$.b")
			assert.NoError(t, err)
			assert.NotNil(t, resType2)

			// Get
			resGet, err := dotfs.JSONGet(key, "$.a")
			assert.NoError(t, err)
			assert.NotNil(t, resGet)

			// Get with no paths
			resGet2, err := dotfs.JSONGet(key)
			assert.NoError(t, err)
			assert.NotNil(t, resGet2)

			// MGet
			key2 := dotpip.NewKey("myjson4")
			_, _ = dotfs.JSONSet(key2, "$", map[string]any{"a": 2})
			resMGet, err := dotfs.JSONMGet("$.a", key, key2)
			assert.NoError(t, err)
			assert.Len(t, resMGet, 2)
		})
	}
}

func TestJSONGetMorePaths(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_more_paths_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": 2, "c": 3})

			// JSONGet with multiple paths
			resGet, err := dotfs.JSONGet(key, "$.a", "$.b")
			assert.NoError(t, err)
			assert.NotNil(t, resGet)
			// Return type is map[string]any{"$.a": [1], "$.b": [2]}

			// JSONGet with no paths defaults to $
			resGet2, err := dotfs.JSONGet(key)
			assert.NoError(t, err)
			assert.NotNil(t, resGet2)

			// JSONDel with $
			resDel, err := dotfs.JSONDel(key, "$")
			assert.NoError(t, err)
			assert.Equal(t, 1, resDel)
		})
	}
}

func TestJSONClearAndMSet(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_clearmset_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			// MSet with multiple missing keys
			k1 := dotpip.NewKey("km1")
			k2 := dotpip.NewKey("km2")

			resMSet, err := dotfs.JSONMSet(
				dotpip.JSONMSetArg{Key: k1, Path: "$", Value: 1},
				dotpip.JSONMSetArg{Key: k2, Path: "$", Value: 2},
			)
			assert.NoError(t, err)
			assert.Equal(t, "OK", resMSet)

			// Clear
			_, _ = dotfs.JSONSet(k1, "$", []any{1, 2, 3})
			resClear, err := dotfs.JSONClear(k1, "$")
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, resClear, 0)

			// Type on missing key
			kmiss := dotpip.NewKey("missing_key")
			resMiss, err := dotfs.JSONType(kmiss, "$")
			assert.NoError(t, err)
			assert.Nil(t, resMiss)
		})
	}
}

func TestJSONClearTrimErrors(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_json_more_cov_test_")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("myjson_errs")
			_, _ = dotfs.JSONSet(key, "$", map[string]any{"a": 1, "b": "str"})

			// Trim error (not array)
			resTrim, err := dotfs.JSONArrTrim(key, "$.a", 0, 1)
			assert.NoError(t, err)
			assert.NotNil(t, resTrim)

			// Index error (not array)
			resIdx, err := dotfs.JSONArrIndex(key, "$.a", 1)
			assert.NoError(t, err)
			assert.NotNil(t, resIdx)

			// Type error handling in JSONArrInsert (not array)
			resIns, err := dotfs.JSONArrInsert(key, "$.a", 0, 1)
			assert.NoError(t, err)
			assert.NotNil(t, resIns)

			// JSONArrAppend error
			resApp, err := dotfs.JSONArrAppend(key, "$.a", 1)
			assert.NoError(t, err)
			assert.NotNil(t, resApp)

			// JSONGet bad path
			resGetBad, err := dotfs.JSONGet(key, "$.a.invalid")
			assert.NoError(t, err)
			assert.Nil(t, resGetBad)

			// JSONDel error (root)
			resDelRoot, _ := dotfs.JSONDel(key, "$")
			assert.Equal(t, 1, resDelRoot)
		})
	}
}
