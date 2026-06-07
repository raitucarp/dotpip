package fs

import (
	"dotpip"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileSystem_PubSub(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip-pubsub-test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fsys := NewFileSystem(tmpDir)
	defer fsys.Close()

	// Ensure subscriptions are working
	sub, err := fsys.Subscribe("channel1", "channel2")
	assert.NoError(t, err)
	defer func() { _ = sub.Close() }()

	// Wait a bit to ensure subscription is registered
	time.Sleep(10 * time.Millisecond)

	count, err := fsys.Publish("channel1", "msg1")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	count, err = fsys.Publish("channel2", "msg2")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	select {
	case msg := <-sub.Channel():
		assert.Equal(t, "channel1", msg.Channel)
		assert.Equal(t, "msg1", msg.Payload)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for msg1")
	}

	select {
	case msg := <-sub.Channel():
		assert.Contains(t, []string{"channel1", "channel2"}, msg.Channel)
		assert.Contains(t, []string{"msg1", "msg2"}, msg.Payload)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for msg2")
	}

	// Unsubscribe
	err = sub.Unsubscribe("channel1")
	assert.NoError(t, err)

	count, err = fsys.Publish("channel1", "msg3")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestFileSystem_KeyspaceNotifications(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip-keyspace-test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fsys := NewFileSystem(tmpDir)
	defer fsys.Close()
	_ = fsys.ConfigSet("notify-keyspace-events", "KEA")

	sub, err := fsys.PSubscribe("__keyspace@0__:*", "__keyevent@0__:*")
	assert.NoError(t, err)
	defer func() { _ = sub.Close() }()

	// Wait for fsnotify to settle
	time.Sleep(100 * time.Millisecond)

	// Perform a SET operation
	_, err = fsys.Set(dotpip.NewKey("mykey"), "myval")
	assert.NoError(t, err)

	// We expect fsnotify to emit CREATE and/or WRITE, resulting in 'set' action
	gotKeyspace := false
	gotKeyevent := false

	timeout := time.After(2 * time.Second)

loop:
	for {
		select {
		case msg := <-sub.Channel():
			if msg.Channel == "__keyspace@0__:mykey" && msg.Payload == "set" {
				gotKeyspace = true
			}
			if msg.Channel == "__keyevent@0__:set" && msg.Payload == "mykey" {
				gotKeyevent = true
			}
			if gotKeyspace && gotKeyevent {
				break loop
			}
		case <-timeout:
			t.Fatalf("Timeout waiting for keyspace notifications")
		}
	}

	assert.True(t, gotKeyspace)
	assert.True(t, gotKeyevent)
}

func TestFileSystem_PubSub_Encodings(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}

	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip-pubsub-enc-test")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			fsys := NewFileSystem(tmpDir)
			fsys.EncodeType(enc)
			defer fsys.Close()

			sub, err := fsys.Subscribe("channel_enc")
			assert.NoError(t, err)
			defer func() { _ = sub.Close() }()

			time.Sleep(10 * time.Millisecond)

			count, err := fsys.Publish("channel_enc", "msg_"+string(enc))
			assert.NoError(t, err)
			assert.Equal(t, 1, count)

			select {
			case msg := <-sub.Channel():
				assert.Equal(t, "channel_enc", msg.Channel)
				assert.Equal(t, "msg_"+string(enc), msg.Payload)
			case <-time.After(1 * time.Second):
				t.Fatal("timeout waiting for msg")
			}
		})
	}
}

func TestFileSystem_PSubscribe(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip-psubscribe-test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fsys := NewFileSystem(tmpDir)
	defer fsys.Close()

	sub, err := fsys.PSubscribe("user_*")
	assert.NoError(t, err)
	defer func() { _ = sub.Close() }()

	time.Sleep(10 * time.Millisecond)

	count, err := fsys.Publish("user_123", "hello")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	select {
	case msg := <-sub.Channel():
		assert.Equal(t, "user_123", msg.Channel)
		assert.Equal(t, "hello", msg.Payload)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for msg")
	}

	err = sub.PUnsubscribe("user_*")
	assert.NoError(t, err)

	count, err = fsys.Publish("user_123", "hello2")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestFileSystem_PubSubStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip-pubsubstats-test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fsys := NewFileSystem(tmpDir)
	defer fsys.Close()

	sub1, _ := fsys.Subscribe("ch1", "ch2")
	defer func() { _ = sub1.Close() }()

	sub2, _ := fsys.Subscribe("ch2")
	defer func() { _ = sub2.Close() }()

	sub3, _ := fsys.PSubscribe("ch*")
	defer func() { _ = sub3.Close() }()

	time.Sleep(10 * time.Millisecond)

	channels, err := fsys.PubSubChannels("")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"ch1", "ch2"}, channels)

	numSub, err := fsys.PubSubNumSub("ch1", "ch2", "ch3")
	assert.NoError(t, err)
	assert.Equal(t, 1, numSub["ch1"])
	assert.Equal(t, 2, numSub["ch2"])
	assert.Equal(t, 0, numSub["ch3"])

	numPat, err := fsys.PubSubNumPat()
	assert.NoError(t, err)
	assert.Equal(t, 1, numPat)
}
