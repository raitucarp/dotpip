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
	defer os.RemoveAll(tmpDir)

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
	assert.Equal(t, "json", resEnc)

	// ObjectEncoding nonexistent
	resEnc, err = dotfs.ObjectEncoding(dotpip.NewKey("no"))
	assert.NoError(t, err)
	assert.Equal(t, "", resEnc)

	// Type nonexistent
	resType, err := dotfs.Type(dotpip.NewKey("no"))
	assert.NoError(t, err)
	assert.Equal(t, "none", resType)
}

func TestZSetsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_zsets_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

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
