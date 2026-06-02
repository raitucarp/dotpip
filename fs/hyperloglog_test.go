package fs

import (
	"dotpip"
	"os"
	"testing"
)

func testHLL(t *testing.T, encodeType FileEncodeType) {
	t.Run(string(encodeType), func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "dotpip-hll-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		f := FileSystem(tempDir)
		f.EncodeType(encodeType)

		// PFAdd
		key1 := dotpip.NewKey("hll1")
		added, err := f.PFAdd(key1, "a", "b", "c")
		if err != nil {
			t.Fatalf("PFAdd failed: %v", err)
		}
		if added != 1 {
			t.Errorf("Expected PFAdd to return 1, got %d", added)
		}

		added, err = f.PFAdd(key1, "a")
		if err != nil {
			t.Fatalf("PFAdd failed: %v", err)
		}
		if added != 0 {
			t.Errorf("Expected PFAdd to return 0, got %d", added)
		}

		// PFCount
		count, err := f.PFCount(key1)
		if err != nil {
			t.Fatalf("PFCount failed: %v", err)
		}
		if count != 3 {
			t.Errorf("Expected PFCount to return 3, got %d", count)
		}

		key2 := dotpip.NewKey("hll2")
		_, err = f.PFAdd(key2, "c", "d", "e")
		if err != nil {
			t.Fatalf("PFAdd failed: %v", err)
		}

		// PFCount multiple
		count, err = f.PFCount(key1, key2)
		if err != nil {
			t.Fatalf("PFCount failed: %v", err)
		}
		if count != 5 {
			t.Errorf("Expected PFCount to return 5, got %d", count)
		}

		// PFMerge
		key3 := dotpip.NewKey("hll3")
		err = f.PFMerge(key3, key1, key2)
		if err != nil {
			t.Fatalf("PFMerge failed: %v", err)
		}

		count, err = f.PFCount(key3)
		if err != nil {
			t.Fatalf("PFCount failed: %v", err)
		}
		if count != 5 {
			t.Errorf("Expected PFCount on merged key to return 5, got %d", count)
		}
	})
}

func TestHyperLogLog(t *testing.T) {
	testHLL(t, JSON)
	testHLL(t, YAML)
	testHLL(t, TOML)
	testHLL(t, RAW)
}
