package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOMLOtherCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_toml_other_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()
	dotfs.EncodeType(fs.TOML)

	key := dotpip.NewKey("toml2")
	_, err = dotfs.TOMLSet(key, "$", map[string]any{"a": []any{1, 2, 3}, "b": "str", "c": map[string]any{"d": 4}, "e": 10.5, "f": true})
	assert.NoError(t, err)

	resIndex, err := dotfs.TOMLArrIndex(key, "$.a", float64(2))
	assert.NoError(t, err)
	assert.Equal(t, []any{1}, resIndex)

	resInsert, err := dotfs.TOMLArrInsert(key, "$.a", 1, 99)
	assert.NoError(t, err)
	assert.Equal(t, []any{4}, resInsert)

	resLen, err := dotfs.TOMLArrLen(key, "$.a")
	assert.NoError(t, err)
	assert.Equal(t, []any{4}, resLen)

	resPop, err := dotfs.TOMLArrPop(key, "$.a", -1)
	assert.NoError(t, err)
	assert.NotNil(t, resPop)

	resTrim, err := dotfs.TOMLArrTrim(key, "$.a", 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []any{2}, resTrim)

	resClear, err := dotfs.TOMLClear(key, "$.c")
	assert.NoError(t, err)
	assert.Equal(t, 1, resClear)

	resDebug, err := dotfs.TOMLDebug("MEMORY", key, "$.b")
	assert.NoError(t, err)
	assert.NotNil(t, resDebug)

	resForget, err := dotfs.TOMLForget(key, "$.b")
	assert.NoError(t, err)
	assert.Equal(t, 1, resForget)

	resMerge, err := dotfs.TOMLMerge(key, "$.c", map[string]any{"x": 1})
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), resMerge)

	resDel, err := dotfs.TOMLDel(key, "$.a")
	assert.NoError(t, err)
	assert.Equal(t, 1, resDel)

	resMGet, err := dotfs.TOMLMGet("$.e", key)
	assert.NoError(t, err)
	assert.Equal(t, []any{10.5}, resMGet)

	k2 := dotpip.NewKey("toml3")
	resMSet, err := dotfs.TOMLMSet(dotpip.JSONMSetArg{Key: k2, Path: "$", Value: map[string]any{"m": 1}})
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), resMSet)

	resIncr, err := dotfs.TOMLNumIncrBy(key, "$.e", 2)
	assert.NoError(t, err)
	assert.Equal(t, []any{12.5}, resIncr)

	resMult, err := dotfs.TOMLNumMultBy(key, "$.e", 2)
	assert.NoError(t, err)
	assert.Equal(t, []any{25.0}, resMult)

	resKeys, err := dotfs.TOMLObjKeys(key, "$")
	assert.NoError(t, err)
	assert.NotNil(t, resKeys)

	resObjLen, err := dotfs.TOMLObjLen(key, "$")
	assert.NoError(t, err)
	assert.NotNil(t, resObjLen)

	resResp, err := dotfs.TOMLResp(key, "$.e")
	assert.NoError(t, err)
	assert.NotNil(t, resResp)

	_, err = dotfs.TOMLSet(key, "$.str", "hello")
	assert.NoError(t, err)

	resStrApp, err := dotfs.TOMLStrAppend(key, "$.str", " world")
	assert.NoError(t, err)
	assert.Equal(t, []any{11}, resStrApp)

	resStrLen, err := dotfs.TOMLStrLen(key, "$.str")
	assert.NoError(t, err)
	assert.Equal(t, []any{11}, resStrLen)

	resToggle, err := dotfs.TOMLToggle(key, "$.f")
	assert.NoError(t, err)
	assert.Equal(t, []any{false}, resToggle)

	resType, err := dotfs.TOMLType(key, "$.str")
	assert.NoError(t, err)
	assert.Equal(t, []any{"string"}, resType)
}
