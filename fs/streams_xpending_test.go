package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamsXPending(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xpending_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")

	// Create stream with some elements
	id1, _ := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})
	id2, _ := dotfs.XAdd(key, "*", map[string]string{"k2": "v2"})

	// Create Group
	_, _ = dotfs.XGroupCreate(key, "mygroup", "0-0", false)

	// Read from group to put messages in pending state
	_, _ = dotfs.XReadGroup("mygroup", "alice", []dotpip.Key{key}, []string{">"}, dotpip.WithXReadGroupCount(2))

	// XPending without range (summary)
	res, err := dotfs.XPending(key, "mygroup")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 4) // Count, MinID, MaxID, Consumers

	// Format is [count, min_id, max_id, [[consumer, count]]]
	assert.Equal(t, 2, res[0])   // Count
	assert.Equal(t, id1, res[1]) // MinID
	assert.Equal(t, id2, res[2]) // MaxID

	// XPending with range
	resRange, err := dotfs.XPending(key, "mygroup", dotpip.WithXPendingRange("-", "+", 10))
	assert.NoError(t, err)
	assert.NotNil(t, resRange)
	assert.Len(t, resRange, 2)

	// For range, format is an array of arrays: [id, consumer, idle_time, deliveries]
	firstMsg := resRange[0].([]any)
	assert.Equal(t, id1, firstMsg[0])
	assert.Equal(t, "alice", firstMsg[1])
}

func TestStreamsXRangeRevRange(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xrange_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")
	id1, _ := dotfs.XAdd(key, "1-0", map[string]string{"k": "v1"})
	id2, _ := dotfs.XAdd(key, "2-0", map[string]string{"k": "v2"})
	id3, _ := dotfs.XAdd(key, "3-0", map[string]string{"k": "v3"})

	// XRange
	resRange, err := dotfs.XRange(key, "-", "+", 0) // count 0 means all
	assert.NoError(t, err)
	assert.Len(t, resRange, 3)
	assert.Equal(t, id1, resRange[0].ID)

	resRangeLimit, err := dotfs.XRange(key, "-", "+", 2)
	assert.NoError(t, err)
	assert.Len(t, resRangeLimit, 2)

	// XRevRange
	resRev, err := dotfs.XRevRange(key, "+", "-", 0) // count 0 means all
	assert.NoError(t, err)
	assert.Len(t, resRev, 3)
	assert.Equal(t, id3, resRev[0].ID)

	resRevLimit, err := dotfs.XRevRange(key, "+", "-", 2)
	assert.NoError(t, err)
	assert.Len(t, resRevLimit, 2)
	assert.Equal(t, id3, resRevLimit[0].ID)
	assert.Equal(t, id2, resRevLimit[1].ID)
}

func TestStreamsXTrim(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xtrim_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")
	for i := 0; i < 10; i++ {
		_, _ = dotfs.XAdd(key, "*", map[string]string{"k": "v"})
	}

	// Trim by MAXLEN
	trimmed, err := dotfs.XTrim(key, dotpip.WithXTrimMaxLen(5, false))
	assert.NoError(t, err)
	assert.Equal(t, 5, trimmed)

	l, _ := dotfs.XLen(key)
	assert.Equal(t, 5, l)

	// Trim by MINID
	// Get elements to find a MINID
	entries, _ := dotfs.XRange(key, "-", "+", 0)
	minID := entries[2].ID // 3rd element will become the new minimum

	trimmed2, err := dotfs.XTrim(key, dotpip.WithXTrimMinID(minID, false))
	assert.NoError(t, err)
	assert.Equal(t, 2, trimmed2) // Trimmed the first 2
}
