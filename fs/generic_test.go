package fs

import (
	"dotpip"
	"os"
	"path/filepath"
	"testing"
)

func TestFlushAll(t *testing.T) {
	err := dotfs.FlushAll()
	if err != nil {
		t.Errorf("Should not have returned an error when flushing all data")
	}
}

func TestCopy(t *testing.T) {
	srcKey := dotpip.NewKey("source")
	dstKey := dotpip.NewKey("destination")

	_, err := dotfs.Set(srcKey, "source value")
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", srcKey)
	}

	copied := dotfs.Copy(srcKey, dstKey, dotpip.WithReplace())
	if copied <= 0 {
		t.Errorf("cannot copy key %s", srcKey)
	}
}

func TestDel(t *testing.T) {
	srcKey := dotpip.NewKey("source")
	dstKey := dotpip.NewKey("destination")

	deleted := dotfs.Del(srcKey, dstKey)
	if deleted <= 0 {
		t.Errorf("cannot delete %s %s", srcKey, dstKey)
	}
}

func TestExists(t *testing.T) {
	srcKey := dotpip.NewKey("source")
	dstKey := dotpip.NewKey("destination")

	keys := []dotpip.Key{srcKey, dstKey}
	exists, err := dotfs.Exists(srcKey, dstKey)
	if err != nil {
		t.Errorf("Should have returned an error for key %s %s", srcKey, dstKey)
	}

	if len(exists) == 0 {
		t.Errorf("Exists should have len %d", len(keys))
	}

	for key, exist := range exists {
		if exist {
			t.Errorf("key %s should exist", keys[key])
		}
	}

	newKey := dotpip.NewKey("newkey")
	_, _ = dotfs.Set(newKey, "new value")
	exists, err = dotfs.Exists(newKey)
	if err != nil {
		t.Errorf("Should have returned an error for key %s %s", srcKey, dstKey)
	}

	if len(exists) == 0 {
		t.Errorf("Exists should have len %d", len(keys))
	}

	if exists[0] == false {
		t.Errorf("%s key should exist", newKey)
	}
}

func BenchmarkFlushAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = dotfs.FlushAll()
	}
}

func TestRename(t *testing.T) {
	_ = dotfs.FlushAll()
	src := dotpip.NewKey("rename_src")
	dst := dotpip.NewKey("rename_dst")

	_, _ = dotfs.Set(src, "value")
	err := dotfs.Rename(src, dst)
	if err != nil {
		t.Errorf("Rename failed: %v", err)
	}

	exists, _ := dotfs.Exists(src)
	if exists[0] {
		t.Errorf("Source key should not exist after rename")
	}

	exists, _ = dotfs.Exists(dst)
	if !exists[0] {
		t.Errorf("Destination key should exist after rename")
	}

	val, _ := dotfs.Get(dst)
	if val != "value" {
		t.Errorf("Value mismatch after rename, expected 'value' got '%s'", val)
	}
}

func TestKeys(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("user:1")
	k2 := dotpip.NewKey("user:2")
	k3 := dotpip.NewKey("other:1")

	_, _ = dotfs.Set(k1, "v1")
	_, _ = dotfs.Set(k2, "v2")
	_, _ = dotfs.Set(k3, "v3")

	keys, _ := dotfs.Keys("user:*")
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestType(t *testing.T) {
	_ = dotfs.FlushAll()
	strKey := dotpip.NewKey("str_type")
	_, _ = dotfs.Set(strKey, "val")

	listKey := dotpip.NewKey("list_type")
	_, _ = dotfs.LPush(listKey, "val1")

	typ, _ := dotfs.Type(strKey)
	if typ != "string" {
		t.Errorf("Expected type string, got %s", typ)
	}

	typ, _ = dotfs.Type(listKey)
	if typ != "list" {
		t.Errorf("Expected type list, got %s", typ)
	}
}

func TestRandomKey(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("user:1")
	k2 := dotpip.NewKey("user:2")
	k3 := dotpip.NewKey("other:1")

	_, _ = dotfs.Set(k1, "v1")
	_, _ = dotfs.Set(k2, "v2")
	_, _ = dotfs.Set(k3, "v3")

	key, err := dotfs.RandomKey()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if key == nil {
		t.Errorf("Expected a key, got nil")
	}

	_ = dotfs.FlushAll()
	key, err = dotfs.RandomKey()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if key != nil {
		t.Errorf("Expected nil key, got %v", key)
	}
}

func TestTouch(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("touch_test1")
	k2 := dotpip.NewKey("touch_test2")

	_, _ = dotfs.Set(k1, "v1")

	count, _ := dotfs.Touch(k1, k2)
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestUnlink(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("unlink_test1")

	_, _ = dotfs.Set(k1, "v1")
	count := dotfs.Unlink(k1)
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
	exist, _ := dotfs.Exists(k1)
	if exist[0] {
		t.Errorf("Expected key to not exist")
	}
}

func TestDumpRestore(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("dump_test1")
	_, _ = dotfs.Set(k1, "v1")

	dump, err := dotfs.Dump(k1)
	if err != nil {
		t.Errorf("Expected no error from Dump, got %v", err)
	}
	if len(dump) == 0 {
		t.Errorf("Expected non-empty dump")
	}

	k2 := dotpip.NewKey("restore_test1")
	err = dotfs.Restore(k2, 0, dump)
	if err != nil {
		t.Errorf("Expected no error from Restore, got %v", err)
	}

	exist, _ := dotfs.Exists(k2)
	if !exist[0] {
		t.Errorf("Expected key to exist after restore")
	}

	val, _ := dotfs.Get(k2)
	if val != "v1" {
		t.Errorf("Expected value 'v1', got '%s'", val)
	}
}

func TestScan(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("user:1")
	k2 := dotpip.NewKey("user:2")
	k3 := dotpip.NewKey("other:1")

	_, _ = dotfs.Set(k1, "v1")
	_, _ = dotfs.Set(k2, "v2")
	_, _ = dotfs.Set(k3, "v3")

	cursor, keys, err := dotfs.Scan(0, dotpip.WithScanMatch("user:*"), dotpip.WithScanCount(10))
	if err != nil {
		t.Errorf("Scan error: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
	if cursor != 0 {
		t.Errorf("Expected cursor to be 0 at end of iteration")
	}
}

func TestObjectCommands(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("obj:1")
	_, _ = dotfs.Set(k1, "v1")

	enc, err := dotfs.ObjectEncoding(k1)
	if err != nil {
		t.Errorf("ObjectEncoding error: %v", err)
	}
	if enc != "json" && enc != "yaml" && enc != "toml" && enc != "raw" {
		t.Errorf("Unexpected encoding: %s", enc)
	}

	freq, err := dotfs.ObjectFreq(k1)
	if err == nil {
		t.Errorf("Expected LFU eviction not supported error")
	}
	_ = freq

	idle, err := dotfs.ObjectIdletime(k1)
	if err != nil {
		t.Errorf("ObjectIdletime error: %v", err)
	}
	if idle < 0 {
		t.Errorf("Expected idle time >= 0, got %d", idle)
	}

	refcount, err := dotfs.ObjectRefcount(k1)
	if err != nil {
		t.Errorf("ObjectRefcount error: %v", err)
	}
	if refcount != 1 {
		t.Errorf("Expected refcount 1, got %d", refcount)
	}
}

func TestDBSizeWaitMove(t *testing.T) {
	_ = dotfs.FlushAll()
	k1 := dotpip.NewKey("db:1")
	_, _ = dotfs.Set(k1, "v1")

	size, err := dotfs.DBSize()
	if err != nil || size != 1 {
		t.Errorf("DBSize error: %v, size: %d", err, size)
	}

	wait, _ := dotfs.Wait(1, 100)
	if wait != 0 {
		t.Errorf("Wait should be 0")
	}

	waitAOF, _, _ := dotfs.WaitAOF(1, 1, 100)
	if waitAOF != 0 {
		t.Errorf("WaitAOF should be 0")
	}

	moved, _ := dotfs.Move(k1, 1)
	if moved != 0 {
		t.Errorf("Move should return 0")
	}
}

func TestFileExistenceAndEncoding(t *testing.T) {
	_ = dotfs.FlushAll()

	// 1. NewKey("a:b:c") should create a/b/c.json
	k1 := dotpip.NewKey("a:b:c")
	_, err := dotfs.Set(k1, "value1")
	if err != nil {
		t.Fatalf("Failed to Set k1: %v", err)
	}

	expectedPath1 := filepath.Join(fssa.pathRoot, "a", "b", "c.json")
	if _, err := os.Stat(expectedPath1); os.IsNotExist(err) {
		t.Errorf("File for NewKey(\"a:b:c\") was not created at expected path: %s", expectedPath1)
	}

	// 2. NewKey("x", "y", "z") should create x/y/z.json
	k2 := dotpip.NewKey("x", "y", "z")
	_, err = dotfs.Set(k2, "value2")
	if err != nil {
		t.Fatalf("Failed to Set k2: %v", err)
	}

	expectedPath2 := filepath.Join(fssa.pathRoot, "x", "y", "z.json")
	if _, err := os.Stat(expectedPath2); os.IsNotExist(err) {
		t.Errorf("File for NewKey(\"x\", \"y\", \"z\") was not created at expected path: %s", expectedPath2)
	}
}

func TestHScan(t *testing.T) {
	key := dotpip.NewKey("myhash")
	dotfs.HSet(key, map[string]string{
		"field1": "val1",
		"field2": "val2",
		"field3": "val3",
		"other":  "val4",
	})

	// Scan all
	cursor, result, err := dotfs.HScan(key, 0, dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 4 {
		t.Errorf("result should have 4 fields")
	}

	// Scan match
	cursor, result, err = dotfs.HScan(key, 0, dotpip.WithScanMatch("field*"), dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 3 {
		t.Errorf("result should have 3 fields")
	}

	// Scan with count
	cursor, result, err = dotfs.HScan(key, 0, dotpip.WithScanCount(2))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 2 {
		t.Errorf("cursor should be 2")
	}
	if len(result) != 2 {
		t.Errorf("result should have 2 fields")
	}
}

func TestSScan(t *testing.T) {
	key := dotpip.NewKey("myset")
	dotfs.SAdd(key, "mem1", "mem2", "mem3", "other")

	// Scan all
	cursor, result, err := dotfs.SScan(key, 0, dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 4 {
		t.Errorf("result should have 4 members")
	}

	// Scan match
	cursor, result, err = dotfs.SScan(key, 0, dotpip.WithScanMatch("mem*"), dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 3 {
		t.Errorf("result should have 3 members")
	}

	// Scan with count
	cursor, result, err = dotfs.SScan(key, 0, dotpip.WithScanCount(2))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 2 {
		t.Errorf("cursor should be 2")
	}
	if len(result) != 2 {
		t.Errorf("result should have 2 members")
	}
}

func TestZScan(t *testing.T) {
	key := dotpip.NewKey("myzset")
	dotfs.ZAdd(key, []dotpip.Z{
		{Score: 1, Member: "mem1"},
		{Score: 2, Member: "mem2"},
		{Score: 3, Member: "mem3"},
		{Score: 4, Member: "other"},
	})

	// Scan all
	cursor, result, err := dotfs.ZScan(key, 0, dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 4 {
		t.Errorf("result should have 4 members")
	}

	// Scan match
	cursor, result, err = dotfs.ZScan(key, 0, dotpip.WithScanMatch("mem*"), dotpip.WithScanCount(100))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 0 {
		t.Errorf("cursor should be 0")
	}
	if len(result) != 3 {
		t.Errorf("result should have 3 members")
	}

	// Scan with count
	cursor, result, err = dotfs.ZScan(key, 0, dotpip.WithScanCount(2))
	if err != nil {
		t.Errorf("Should not have returned an error for key %s", key)
	}
	if cursor != 2 {
		t.Errorf("cursor should be 2")
	}
	if len(result) != 2 {
		t.Errorf("result should have 2 members")
	}
}
