package fs_test

import (
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubSubMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_pubsub_more_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	// PSubscribe
	psub, err := dotfs.PSubscribe("channel*")
	assert.NoError(t, err)
	assert.NotNil(t, psub)

	// PUnsubscribe
	err = psub.PUnsubscribe("channel*")
	assert.NoError(t, err)

	// SSubscribe
	ssub, err := dotfs.SSubscribe("channel")
	assert.NoError(t, err)
	assert.NotNil(t, ssub)

	// SUnsubscribe
	err = ssub.SUnsubscribe("channel")
	assert.NoError(t, err)

	// PubSubChannels
	channels, err := dotfs.PubSubChannels("channel*")
	assert.NoError(t, err)
	_ = channels

	// PubSubNumSub
	numSub, err := dotfs.PubSubNumSub("channel1", "channel2")
	assert.NoError(t, err)
	assert.NotNil(t, numSub)

	// PubSubNumPat
	numPat, err := dotfs.PubSubNumPat()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, numPat, 0)

	// PubSubShardChannels
	schan, err := dotfs.PubSubShardChannels("channel*")
	assert.NoError(t, err)
	_ = schan

	// PubSubShardNumSub
	snumSub, err := dotfs.PubSubShardNumSub("channel1")
	assert.NoError(t, err)
	assert.NotNil(t, snumSub)

	// Unsubscribe all
	err = psub.Unsubscribe()
	assert.NoError(t, err)
	err = ssub.Unsubscribe()
	assert.NoError(t, err)
}

func TestPubSubShardChannels(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_pubsub_shard_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	sub, _ := dotfs.SSubscribe("channel1", "channel2")
	defer func() { _ = sub.Close() }()

	chans, err := dotfs.PubSubShardChannels("")
	assert.NoError(t, err)
	_ = chans

	chans, err = dotfs.PubSubShardChannels("channel*")
	assert.NoError(t, err)
	_ = chans

	chans, err = dotfs.PubSubShardChannels("nomatch")
	assert.NoError(t, err)
	_ = chans
}

func TestPubSubUnsubscribeEmptyArgs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_pubsub_empty_args_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	psub, _ := dotfs.PSubscribe("p1", "p2")
	_ = psub.PUnsubscribe() // Unsubscribe all patterns

	ssub, _ := dotfs.SSubscribe("s1", "s2")
	_ = ssub.SUnsubscribe() // Unsubscribe all shard channels
}
