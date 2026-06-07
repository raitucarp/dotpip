package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamsXGroupDestroyAndSetID(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xgroup_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")

	// Create stream by adding an element
	_, err = dotfs.XAdd(key, "*", map[string]string{"field1": "value1"})
	assert.NoError(t, err)

	// Create Group
	res, err := dotfs.XGroupCreate(key, "mygroup", "$", false)
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), res)

	// Set ID
	_, err = dotfs.XGroupSetID(key, "mygroup", "0-0")
	assert.NoError(t, err)
	assert.Equal(t, string(dotpip.StatusOK), res)

	// Set ID with invalid group
	_, err = dotfs.XGroupSetID(key, "mygroup_invalid", "0-0")
	assert.Error(t, err)

	// Destroy Group
	resInt, err := dotfs.XGroupDestroy(key, "mygroup")
	assert.NoError(t, err)
	assert.Equal(t, 1, resInt)

	// Destroy non-existent Group
	resInt, err = dotfs.XGroupDestroy(key, "mygroup")
	assert.NoError(t, err)
	assert.Equal(t, 0, resInt)

	// XGroupDelConsumer
	_, err = dotfs.XGroupCreate(key, "mygroup2", "$", false)
	assert.NoError(t, err)
	_, err = dotfs.XGroupCreateConsumer(key, "mygroup2", "alice")
	assert.NoError(t, err)

	resDelCons, err := dotfs.XGroupDelConsumer(key, "mygroup2", "alice")
	assert.NoError(t, err)
	assert.Equal(t, 0, resDelCons) // returns number of pending messages consumer had, which is 0
}

func TestStreamsXInfo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_xinfo_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream")
	_, _ = dotfs.XAdd(key, "*", map[string]string{"f1": "v1"})
	_, _ = dotfs.XGroupCreate(key, "mygroup", "$", false)
	_, _ = dotfs.XGroupCreateConsumer(key, "mygroup", "alice")

	// XInfoGroups
	groups, err := dotfs.XInfoGroups(key)
	assert.NoError(t, err)
	assert.Len(t, groups, 1)

	// XInfoConsumers
	consumers, err := dotfs.XInfoConsumers(key, "mygroup")
	assert.NoError(t, err)
	assert.Len(t, consumers, 1)
}
