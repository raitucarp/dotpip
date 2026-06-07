package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStreamsXClaim(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xclaim_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")

	// Create stream with some elements
	id1, _ := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})

	// Create Group
	_, _ = dotfs.XGroupCreate(key, "mygroup", "0-0", false)

	// Read from group to put messages in pending state for consumer Alice
	_, _ = dotfs.XReadGroup("mygroup", "alice", []dotpip.Key{key}, []string{">"}, dotpip.WithXReadGroupCount(1))

	// Claim by consumer Bob
	// Use idle time 0 so it immediately qualifies
	claimed, err := dotfs.XClaim(key, "mygroup", "bob", 0, []string{id1})
	assert.NoError(t, err)
	assert.Len(t, claimed, 1)
	assert.Equal(t, id1, claimed[0].ID)

	// Verify Bob now has it in pending
	resRange, _ := dotfs.XPending(key, "mygroup", dotpip.WithXPendingRange("-", "+", 10))
	assert.Len(t, resRange, 1)
	msg := resRange[0].([]any)
	assert.Equal(t, id1, msg[0])
	assert.Equal(t, "bob", msg[1])
}

func TestStreamsXAutoClaim(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xautoclaim_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")
	id1, _ := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})
	id2, _ := dotfs.XAdd(key, "*", map[string]string{"k2": "v2"})

	_, _ = dotfs.XGroupCreate(key, "mygroup", "0-0", false)

	// Alice reads both
	_, _ = dotfs.XReadGroup("mygroup", "alice", []dotpip.Key{key}, []string{">"}, dotpip.WithXReadGroupCount(2))

	// Wait a tiny bit so idle time is non-zero (even 1ms is fine if we use 0)
	time.Sleep(10 * time.Millisecond)

	// Bob auto-claims from start "0-0"
	nextStart, claimed, err := dotfs.XAutoClaim(key, "mygroup", "bob", 0, "0-0", dotpip.WithXAutoClaimCount(1))
	assert.NoError(t, err)
	assert.Len(t, claimed, 1)
	assert.Equal(t, id1, claimed[0].ID)
	// Next start will be id2 or > id1
	assert.NotEqual(t, "0-0", nextStart)

	// Claim next
	_, claimed2, err := dotfs.XAutoClaim(key, "mygroup", "bob", 0, nextStart, dotpip.WithXAutoClaimCount(1))
	assert.NoError(t, err)
	assert.Len(t, claimed2, 1)
	assert.Equal(t, id2, claimed2[0].ID)
}
