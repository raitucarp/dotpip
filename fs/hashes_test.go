package fs

import (
	"dotpip"
	"os"
	"testing"
)

func runHashTests(t *testing.T, encodeType FileEncodeType) {
	tempDir, err := os.MkdirTemp("", "dotpip-hash-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := FileSystem(tempDir)
	defer db.Close()
	db.EncodeType(encodeType)

	key := dotpip.NewKey("myhash")

	t.Run("HSet and HGet", func(t *testing.T) {
		added, err := db.HSet(key, map[string]string{"field1": "value1", "field2": "value2"})
		if err != nil {
			t.Fatalf("HSet error: %v", err)
		}
		if added != 2 {
			t.Errorf("Expected 2 added, got %d", added)
		}

		val, err := db.HGet(key, "field1")
		if err != nil {
			t.Fatalf("HGet error: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected value1, got %s", val)
		}
	})

	t.Run("HExists", func(t *testing.T) {
		exists, err := db.HExists(key, "field1")
		if err != nil || !exists {
			t.Errorf("Expected field1 to exist")
		}

		exists, err = db.HExists(key, "field3")
		if err != nil || exists {
			t.Errorf("Expected field3 to not exist")
		}
	})

	t.Run("HGetAll", func(t *testing.T) {
		all, err := db.HGetAll(key)
		if err != nil {
			t.Fatalf("HGetAll error: %v", err)
		}
		if len(all) != 2 || all["field1"] != "value1" || all["field2"] != "value2" {
			t.Errorf("Unexpected HGetAll result: %v", all)
		}
	})

	t.Run("HDel", func(t *testing.T) {
		deleted, err := db.HDel(key, "field1", "field3")
		if err != nil {
			t.Fatalf("HDel error: %v", err)
		}
		if deleted != 1 {
			t.Errorf("Expected 1 deleted, got %d", deleted)
		}

		exists, _ := db.HExists(key, "field1")
		if exists {
			t.Errorf("Expected field1 to be deleted")
		}
	})

	t.Run("HIncrBy", func(t *testing.T) {
		db.HSet(key, map[string]string{"counter": "10"})
		val, err := db.HIncrBy(key, "counter", 5)
		if err != nil {
			t.Fatalf("HIncrBy error: %v", err)
		}
		if val != 15 {
			t.Errorf("Expected 15, got %d", val)
		}
	})

	t.Run("HIncrByFloat", func(t *testing.T) {
		db.HSet(key, map[string]string{"fcounter": "10.5"})
		val, err := db.HIncrByFloat(key, "fcounter", 0.5)
		if err != nil {
			t.Fatalf("HIncrByFloat error: %v", err)
		}
		if val != 11.0 {
			t.Errorf("Expected 11.0, got %f", val)
		}
	})

	t.Run("HKeys and HVals", func(t *testing.T) {
		db.Del(key)
		db.HSet(key, map[string]string{"k1": "v1", "k2": "v2"})

		keys, err := db.HKeys(key)
		if err != nil || len(keys) != 2 {
			t.Errorf("Unexpected HKeys result")
		}

		vals, err := db.HVals(key)
		if err != nil || len(vals) != 2 {
			t.Errorf("Unexpected HVals result")
		}
	})

	t.Run("HLen", func(t *testing.T) {
		l, err := db.HLen(key)
		if err != nil || l != 2 {
			t.Errorf("Expected length 2, got %d", l)
		}
	})

	t.Run("HMGet", func(t *testing.T) {
		vals, err := db.HMGet(key, "k1", "nonexistent", "k2")
		if err != nil {
			t.Fatalf("HMGet error: %v", err)
		}
		if len(vals) != 3 || vals[0] != "v1" || vals[1] != "" || vals[2] != "v2" {
			t.Errorf("Unexpected HMGet result: %v", vals)
		}
	})

	t.Run("HSetNX", func(t *testing.T) {
		set, err := db.HSetNX(key, "k1", "newval")
		if err != nil || set {
			t.Errorf("Expected false for existing key")
		}

		set, err = db.HSetNX(key, "k3", "v3")
		if err != nil || !set {
			t.Errorf("Expected true for new key")
		}
	})

	t.Run("HStrLen", func(t *testing.T) {
		l, err := db.HStrLen(key, "k3") // "v3"
		if err != nil || l != 2 {
			t.Errorf("Expected length 2, got %d", l)
		}
	})

	t.Run("HRandField", func(t *testing.T) {
		fields, err := db.HRandField(key, 2)
		if err != nil || len(fields) != 2 {
			t.Errorf("Expected 2 fields, got %v", fields)
		}

		fieldsVals, err := db.HRandField(key, 2, dotpip.WithHRandFieldWithValues())
		if err != nil || len(fieldsVals) != 4 {
			t.Errorf("Expected 4 items (2 pairs), got %v", fieldsVals)
		}
	})
}

func TestHashCommands(t *testing.T) {
	for _, enc := range []FileEncodeType{JSON, YAML, TOML, RAW} {
		t.Run(string(enc), func(t *testing.T) {
			runHashTests(t, enc)
		})
	}
}
