package fs

import (
	"fmt"
	"os"
	"testing"

	"dotpip"
)

func TestStreams(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}

	for _, encoding := range encodings {
		t.Run(fmt.Sprintf("Encoding_%s", encoding), func(t *testing.T) {
			path := fmt.Sprintf("/tmp/dotpip_streams_%s", encoding)
			_ = os.RemoveAll(path)
			_ = os.MkdirAll(path, 0755)

			fs := NewFileSystem(path)
			fs.EncodeType(encoding)

			key := dotpip.NewKey("mystream")

			// Test XAdd
			id1, err := fs.XAdd(key, "1000-0", map[string]string{"name": "Alice", "age": "30"})
			if err != nil {
				t.Fatalf("XAdd failed: %v", err)
			}
			if id1 != "1000-0" {
				t.Errorf("Expected ID 1000-0, got %s", id1)
			}

			_, err = fs.XAdd(key, "1000-1", map[string]string{"name": "Bob", "age": "35"})
			if err != nil {
				t.Fatalf("XAdd failed: %v", err)
			}

			// Test XLen
			length, err := fs.XLen(key)
			if err != nil {
				t.Fatalf("XLen failed: %v", err)
			}
			if length != 2 {
				t.Errorf("Expected length 2, got %d", length)
			}

			// Test XRange
			entries, err := fs.XRange(key, "-", "+", 0)
			if err != nil {
				t.Fatalf("XRange failed: %v", err)
			}
			if len(entries) != 2 {
				t.Errorf("Expected 2 entries, got %d", len(entries))
			}
			if entries[0].ID != "1000-0" || entries[0].Values["name"] != "Alice" {
				t.Errorf("Unexpected entry 0: %v", entries[0])
			}

			// Test XRevRange
			revEntries, err := fs.XRevRange(key, "+", "-", 1)
			if err != nil {
				t.Fatalf("XRevRange failed: %v", err)
			}
			if len(revEntries) != 1 || revEntries[0].ID != "1000-1" {
				t.Errorf("Unexpected rev entry 0: %v", revEntries)
			}

			// Test XRead
			readResult, err := fs.XRead([]dotpip.Key{key}, []string{"1000-0"})
			if err != nil {
				t.Fatalf("XRead failed: %v", err)
			}
			if len(readResult["mystream"]) != 1 || readResult["mystream"][0].ID != "1000-1" {
				t.Errorf("Unexpected read result: %v", readResult)
			}

			// Test XGroupCreate
			groupStatus, err := fs.XGroupCreate(key, "mygroup", "0", true)
			if err != nil {
				t.Fatalf("XGroupCreate failed: %v", err)
			}
			if groupStatus != string(dotpip.StatusOK) {
				t.Errorf("Expected OK, got %s", groupStatus)
			}

			// Test XGroupCreateConsumer
			consStatus, err := fs.XGroupCreateConsumer(key, "mygroup", "cons1")
			if err != nil {
				t.Fatalf("XGroupCreateConsumer failed: %v", err)
			}
			if consStatus != 1 {
				t.Errorf("Expected 1, got %d", consStatus)
			}

			// Test XReadGroup
			groupReadResult, err := fs.XReadGroup("mygroup", "cons1", []dotpip.Key{key}, []string{">"})
			if err != nil {
				t.Fatalf("XReadGroup failed: %v", err)
			}
			if len(groupReadResult["mystream"]) != 2 {
				t.Errorf("Expected 2 messages in read group, got %d", len(groupReadResult["mystream"]))
			}

			// Test XAck
			ackCount, err := fs.XAck(key, "mygroup", "1000-0")
			if err != nil {
				t.Fatalf("XAck failed: %v", err)
			}
			if ackCount != 1 {
				t.Errorf("Expected 1 ack, got %d", ackCount)
			}

			// Verify Pending
			groupReadResult2, err := fs.XReadGroup("mygroup", "cons1", []dotpip.Key{key}, []string{"0"})
			if err != nil {
				t.Fatalf("XReadGroup failed: %v", err)
			}
			if len(groupReadResult2["mystream"]) != 1 || groupReadResult2["mystream"][0].ID != "1000-1" {
				t.Errorf("Expected 1 pending message, got %v", groupReadResult2["mystream"])
			}

			// Test XTrim
			trimCount, err := fs.XTrim(key, dotpip.WithXTrimMaxLen(1, false))
			if err != nil {
				t.Fatalf("XTrim failed: %v", err)
			}
			if trimCount != 1 {
				t.Errorf("Expected 1 trimmed, got %d", trimCount)
			}

			lengthAfterTrim, _ := fs.XLen(key)
			if lengthAfterTrim != 1 {
				t.Errorf("Expected length 1 after trim, got %d", lengthAfterTrim)
			}

			// Test XDel
			delCount, err := fs.XDel(key, "1000-1")
			if err != nil {
				t.Fatalf("XDel failed: %v", err)
			}
			if delCount != 1 {
				t.Errorf("Expected 1 deleted, got %d", delCount)
			}

			// Note: 1000-1 was deleted but it's still in the pending list for group because XDel doesn't automatically ack or remove from pending. We should ack it.
			_, _ = fs.XAck(key, "mygroup", "1000-1")

			lengthAfterDel, _ := fs.XLen(key)
			if lengthAfterDel != 0 {
				t.Errorf("Expected length 0 after del, got %d", lengthAfterDel)
			}

			// Add items back for next tests
			id1, _ = fs.XAdd(key, "2000-0", map[string]string{"foo": "bar"})
			id2, _ := fs.XAdd(key, "2000-1", map[string]string{"foo": "baz"})

			// Test XPending summary
			_, _ = fs.XReadGroup("mygroup", "cons1", []dotpip.Key{key}, []string{">"})
			pendingResult, err := fs.XPending(key, "mygroup")
			if err != nil {
				t.Fatalf("XPending failed: %v", err)
			}
			if len(pendingResult) != 4 || pendingResult[0].(int) != 2 {
				t.Errorf("Expected 4 items in XPending summary with 2 total pending, got: %v", pendingResult)
			}

			// Test XClaim
			claimed, err := fs.XClaim(key, "mygroup", "cons2", 0, []string{id1})
			if err != nil {
				t.Fatalf("XClaim failed: %v", err)
			}
			if len(claimed) != 1 || claimed[0].ID != id1 {
				t.Errorf("Expected XClaim to claim id1, got: %v", claimed)
			}

			// Test XAutoClaim
			nextID, autoClaimed, err := fs.XAutoClaim(key, "mygroup", "cons2", 0, id2, dotpip.WithXAutoClaimCount(10))
			if err != nil {
				t.Fatalf("XAutoClaim failed: %v", err)
			}
			if len(autoClaimed) != 1 || autoClaimed[0].ID != id2 {
				t.Errorf("Expected XAutoClaim to claim id2, got: %v", autoClaimed)
			}
			if nextID == "" {
				t.Errorf("Expected nextID from XAutoClaim, got empty")
			}

			// Test XInfoStream
			infoStream, err := fs.XInfoStream(key)
			if err != nil {
				t.Fatalf("XInfoStream failed: %v", err)
			}
			if infoStream["length"].(int) != 2 || infoStream["groups"].(int) != 1 {
				t.Errorf("Unexpected XInfoStream result: %v", infoStream)
			}

			// Cleanup
			_ = os.RemoveAll(path)
		})
	}
}

func TestStreamIDParsing(t *testing.T) {
	fs := &FileSystem{}

	ms, seq, err := fs.parseStreamID("12345-67")
	if err != nil || ms != 12345 || seq != 67 {
		t.Errorf("Failed parsing: %v, %d, %d", err, ms, seq)
	}

	ms, seq, err = fs.parseStreamID("12345")
	if err != nil || ms != 12345 || seq != 0 {
		t.Errorf("Failed parsing implicit seq: %v, %d, %d", err, ms, seq)
	}
}
