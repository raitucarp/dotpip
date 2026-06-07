package fs

import (
	"dotpip"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileSystem_SubkeyspaceNotifications(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip-subkeyspace-test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fsys := NewFileSystem(tmpDir)
	defer fsys.Close()

	// Enable Subkeyspace, Subkeyevent, Subkeyspaceitem, Subkeyspaceevent and hash commands
	_ = fsys.ConfigSet("notify-keyspace-events", "STIVh")

	sub, err := fsys.PSubscribe("__subkeyspace@0__:*", "__subkeyevent@0__:*", "__subkeyspaceitem@0__:*", "__subkeyspaceevent@0__:*")
	assert.NoError(t, err)
	defer func() { _ = sub.Close() }()

	// Wait for subscription to settle
	time.Sleep(10 * time.Millisecond)

	// Perform an HSET operation
	_, err = fsys.HSet(dotpip.NewKey("myhash"), map[string]string{"field1": "val1", "field2": "val2"})
	assert.NoError(t, err)

	expectedMessages := 5
	received := 0

	gotSubkeyspace := false
	gotSubkeyevent := false
	gotSubkeyspaceitem1 := false
	gotSubkeyspaceitem2 := false
	gotSubkeyspaceevent := false

	timeout := time.After(2 * time.Second)

loop:
	for {
		select {
		case msg := <-sub.Channel():
			received++
			if msg.Channel == "__subkeyspace@0__:myhash" && strings.HasPrefix(msg.Payload, "hset|") {
				gotSubkeyspace = true
			}
			if msg.Channel == "__subkeyevent@0__:hset" && strings.HasPrefix(msg.Payload, "6:myhash|") {
				gotSubkeyevent = true
			}
			if msg.Channel == "__subkeyspaceitem@0__:myhash\nfield1" && msg.Payload == "hset" {
				gotSubkeyspaceitem1 = true
			}
			if msg.Channel == "__subkeyspaceitem@0__:myhash\nfield2" && msg.Payload == "hset" {
				gotSubkeyspaceitem2 = true
			}
			if msg.Channel == "__subkeyspaceevent@0__:hset|myhash" {
				gotSubkeyspaceevent = true
			}
			if received == expectedMessages {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	assert.True(t, gotSubkeyspace, "missing subkeyspace notification")
	assert.True(t, gotSubkeyevent, "missing subkeyevent notification")
	assert.True(t, gotSubkeyspaceitem1, "missing subkeyspaceitem1 notification")
	assert.True(t, gotSubkeyspaceitem2, "missing subkeyspaceitem2 notification")
	assert.True(t, gotSubkeyspaceevent, "missing subkeyspaceevent notification")
}
