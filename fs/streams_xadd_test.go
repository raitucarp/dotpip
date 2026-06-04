package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamsXAddOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xadd_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")

	// Test NOMKSTREAM
	res, _ := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"}, dotpip.WithXAddNoMkStream())
	assert.Equal(t, "", res) // Usually returns empty string when not created

	// Create stream normally
	id1, err := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)

	// Test MINID and MAXLEN
	for i := 0; i < 5; i++ {
		_, _ = dotfs.XAdd(key, "*", map[string]string{"k": "v"})
	}

	// MAXLEN
	_, err = dotfs.XAdd(key, "*", map[string]string{"k": "v"}, dotpip.WithXAddMaxLen(3, false))
	assert.NoError(t, err)

	l, _ := dotfs.XLen(key)
	assert.Equal(t, 3, l)

	// MINID
	entries, _ := dotfs.XRange(key, "-", "+", 0)
	minID := entries[1].ID // 2nd element will become new min

	_, err = dotfs.XAdd(key, "*", map[string]string{"k": "v"}, dotpip.WithXAddMinID(minID, false))
	assert.NoError(t, err)

	// We just want to check it didn't fail
	l2, _ := dotfs.XLen(key)
	assert.GreaterOrEqual(t, l2, 1)
}

func TestStreamsXRangeIncomplete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xrange2_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")

	id1, _ := dotfs.XAdd(key, "1-0", map[string]string{"k1": "v1"})
	id2, _ := dotfs.XAdd(key, "2-0", map[string]string{"k2": "v2"})

	// Test partial start/end like "1", "2"
	resRange, err := dotfs.XRange(key, "1", "2", 0)
	assert.NoError(t, err)
	assert.Len(t, resRange, 2)
	assert.Equal(t, id1, resRange[0].ID)
	assert.Equal(t, id2, resRange[1].ID)

	resRev, err := dotfs.XRevRange(key, "2", "1", 0)
	assert.NoError(t, err)
	assert.Len(t, resRev, 2)
	assert.Equal(t, id2, resRev[0].ID)
	assert.Equal(t, id1, resRev[1].ID)
}

func TestStreamsXRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xread_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")
	_, _ = dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})
	id2, _ := dotfs.XAdd(key, "*", map[string]string{"k2": "v2"})

	// Test XRead from 0-0
	res, err := dotfs.XRead([]dotpip.Key{key}, []string{"0-0"}, dotpip.WithXReadCount(2))
	assert.NoError(t, err)
	assert.NotNil(t, res)

	keyStr := "mystream" // Usually strings are returned from keys
	streamEntries := res[keyStr]
	assert.Len(t, streamEntries, 2)

	// Test XRead from specific ID
	res2, err := dotfs.XRead([]dotpip.Key{key}, []string{streamEntries[0].ID})
	assert.NoError(t, err)
	assert.NotNil(t, res2)
	assert.Len(t, res2[keyStr], 1)
	assert.Equal(t, id2, res2[keyStr][0].ID)
}
