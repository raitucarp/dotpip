package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_more_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey")
	_, _ = dotfs.Set(key, "val")

	key2 := dotpip.NewKey("mykey2")

	// RenameNX
	resRen, err := dotfs.RenameNX(key, key2)
	assert.NoError(t, err)
	assert.Equal(t, true, resRen) // OK since key2 does not exist

	// RenameNX when destination exists
	_, _ = dotfs.Set(key, "val")
	resRen2, err := dotfs.RenameNX(key, key2)
	assert.NoError(t, err)
	assert.Equal(t, false, resRen2)

	// ObjectEncoding
	resEnc, err := dotfs.ObjectEncoding(key)
	assert.NoError(t, err)
	assert.Equal(t, dotpip.ObjectEncodingJSON, resEnc)

	// ObjectEncoding nonexistent
	resEnc, err = dotfs.ObjectEncoding(dotpip.NewKey("no"))
	assert.NoError(t, err)
	assert.Equal(t, dotpip.ObjectEncoding(""), resEnc)

	// Type nonexistent
	resType, err := dotfs.Type(dotpip.NewKey("no"))
	assert.NoError(t, err)
	assert.Equal(t, dotpip.ObjectTypeNone, resType)
}

func TestZSetsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_zsets_more_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key1 := dotpip.NewKey("zset1")
	key2 := dotpip.NewKey("zset2")

	_, _ = dotfs.ZAdd(key1, []dotpip.Z{{Score: 1, Member: "a"}, {Score: 2, Member: "b"}})
	_, _ = dotfs.ZAdd(key2, []dotpip.Z{{Score: 2, Member: "b"}, {Score: 3, Member: "c"}})

	// ZDiff
	resDiff, err := dotfs.ZDiff(key1, key2)
	assert.NoError(t, err)
	assert.Len(t, resDiff, 1)
	assert.Equal(t, "a", resDiff[0])

	// ZDiffWithScores
	resDiffScores, err := dotfs.ZDiffWithScores(key1, key2)
	assert.NoError(t, err)
	assert.Len(t, resDiffScores, 1)
	assert.Equal(t, "a", resDiffScores[0].Member)

	// ZUnion
	resUnion, err := dotfs.ZUnion(key1, key2)
	assert.NoError(t, err)
	assert.Len(t, resUnion, 3) // a, b, c
}

func TestGenericRestoreScan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_more_cov_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	// Rename missing source
	errRen := dotfs.Rename(dotpip.NewKey("notsource"), dotpip.NewKey("dest"))
	assert.Error(t, errRen)

	// Rename to same key
	_, _ = dotfs.Set(dotpip.NewKey("same"), "val")
	errRen2 := dotfs.Rename(dotpip.NewKey("same"), dotpip.NewKey("same"))
	assert.NoError(t, errRen2)

	// Copy missing
	resCopy := dotfs.Copy(dotpip.NewKey("notcopy"), dotpip.NewKey("dest2"))
	_ = resCopy

	// Copy same key
	resCopy2 := dotfs.Copy(dotpip.NewKey("same"), dotpip.NewKey("same"))
	assert.Equal(t, 1, resCopy2)

	// Scan
	for i := 0; i < 5; i++ {
		_, _ = dotfs.Set(dotpip.NewKey("scan_key"), "val")
	}

	_, keys, err := dotfs.Scan(0, dotpip.WithScanMatch("*"), dotpip.WithScanCount(100), dotpip.WithScanType(string(dotpip.ObjectTypeString)))
	assert.NoError(t, err)
	assert.NotNil(t, keys)
}

func TestGenericTypeCopyRename(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_more2_cov_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey_tcr")
	_, _ = dotfs.Set(key, "val")

	// Type
	typ, _ := dotfs.Type(key)
	assert.Equal(t, dotpip.ObjectTypeString, typ)

	// Rename to self
	errRename := dotfs.Rename(key, key)
	assert.NoError(t, errRename)

	// Copy to self
	resCopy := dotfs.Copy(key, key)
	_ = resCopy

	// Copy to existing key
	key2 := dotpip.NewKey("mykey2_tcr")
	_, _ = dotfs.Set(key2, "val2")
	resCopy2 := dotfs.Copy(key, key2)
	assert.Equal(t, 1, resCopy2)

	// Copy with REPLACE
	resCopy3 := dotfs.Copy(key, key2, dotpip.WithReplace())
	assert.Equal(t, 1, resCopy3)

	// RandomKey empty db
	dotfs2 := fs.NewFileSystem(tmpDir + "_empty")
	defer dotfs2.Close()
	defer func() { _ = os.RemoveAll(tmpDir + "_empty") }()
	rk, _ := dotfs2.RandomKey()
	assert.Nil(t, rk)
}
