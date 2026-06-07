package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLOtherCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_yaml_other_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()
	dotfs.EncodeType(fs.YAML)

	key := dotpip.NewKey("yaml2")
	_, err = dotfs.YAMLSet(key, "$", map[string]any{"a": []any{1, 2, 3}, "b": "str", "c": map[string]any{"d": 4}, "e": 10.5, "f": true})
	assert.NoError(t, err)

	resIndex, err := dotfs.YAMLArrIndex(key, "$.a", float64(2))
	assert.NoError(t, err)
	assert.Equal(t, []any{1}, resIndex)

	resInsert, err := dotfs.YAMLArrInsert(key, "$.a", 1, 99)
	assert.NoError(t, err)
	assert.Equal(t, []any{4}, resInsert)

	resLen, err := dotfs.YAMLArrLen(key, "$.a")
	assert.NoError(t, err)
	assert.Equal(t, []any{4}, resLen)

	resPop, err := dotfs.YAMLArrPop(key, "$.a", -1)
	assert.NoError(t, err)
	assert.NotNil(t, resPop)

	resTrim, err := dotfs.YAMLArrTrim(key, "$.a", 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []any{2}, resTrim)

	resClear, err := dotfs.YAMLClear(key, "$.c")
	assert.NoError(t, err)
	assert.Equal(t, 1, resClear)

	resDebug, err := dotfs.YAMLDebug("MEMORY", key, "$.b")
	assert.NoError(t, err)
	assert.NotNil(t, resDebug)

	resForget, err := dotfs.YAMLForget(key, "$.b")
	assert.NoError(t, err)
	assert.Equal(t, 1, resForget)

	resMerge, err := dotfs.YAMLMerge(key, "$.c", map[string]any{"x": 1})
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), resMerge)

	resDel, err := dotfs.YAMLDel(key, "$.a")
	assert.NoError(t, err)
	assert.Equal(t, 1, resDel)

	resMGet, err := dotfs.YAMLMGet("$.e", key)
	assert.NoError(t, err)
	assert.Equal(t, []any{10.5}, resMGet)

	k2 := dotpip.NewKey("yaml3")
	resMSet, err := dotfs.YAMLMSet(dotpip.JSONMSetArg{Key: k2, Path: "$", Value: map[string]any{"m": 1}})
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), resMSet)

	resIncr, err := dotfs.YAMLNumIncrBy(key, "$.e", 2)
	assert.NoError(t, err)
	assert.Equal(t, []any{12.5}, resIncr)

	resMult, err := dotfs.YAMLNumMultBy(key, "$.e", 2)
	assert.NoError(t, err)
	assert.Equal(t, []any{25.0}, resMult)

	resKeys, err := dotfs.YAMLObjKeys(key, "$")
	assert.NoError(t, err)
	assert.NotNil(t, resKeys)

	resObjLen, err := dotfs.YAMLObjLen(key, "$")
	assert.NoError(t, err)
	assert.NotNil(t, resObjLen)

	resResp, err := dotfs.YAMLResp(key, "$.e")
	assert.NoError(t, err)
	assert.NotNil(t, resResp)

	_, err = dotfs.YAMLSet(key, "$.str", "hello")
	assert.NoError(t, err)

	resStrApp, err := dotfs.YAMLStrAppend(key, "$.str", " world")
	assert.NoError(t, err)
	assert.Equal(t, []any{11}, resStrApp)

	resStrLen, err := dotfs.YAMLStrLen(key, "$.str")
	assert.NoError(t, err)
	assert.Equal(t, []any{11}, resStrLen)

	resToggle, err := dotfs.YAMLToggle(key, "$.f")
	assert.NoError(t, err)
	assert.Equal(t, []any{false}, resToggle)

	resType, err := dotfs.YAMLType(key, "$.str")
	assert.NoError(t, err)
	assert.Equal(t, []any{"string"}, resType)
}
