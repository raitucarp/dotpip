package fs

import (
	"dotpip"
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

	if len(exists) <= 0 {
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

	if len(exists) <= 0 {
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
