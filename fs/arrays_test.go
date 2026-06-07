package fs

import (
	"dotpip"
	"os"
	"testing"
)

func TestArrayCommands(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}

	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "dotpip-arrays-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer func() { _ = os.RemoveAll(tempDir) }()

			fsys := NewFileSystem(tempDir)
			defer fsys.Close()
			fsys.EncodeType(enc)

			key := dotpip.Key{"myarray"}

			// Test ARSET and ARGET
			n, err := fsys.ARSet(key, 0, "hello")
			if err != nil {
				t.Fatalf("ARSet error: %v", err)
			}
			if n != 1 {
				t.Errorf("expected 1 slot filled, got %d", n)
			}

			val, err := fsys.ARGet(key, 0)
			if err != nil {
				t.Fatalf("ARGet error: %v", err)
			}
			if val != "hello" {
				t.Errorf("expected 'hello', got '%s'", val)
			}

			// Test multiple values
			n, err = fsys.ARSet(key, 2, "a", "b", "c")
			if err != nil {
				t.Fatalf("ARSet multiple error: %v", err)
			}
			if n != 3 {
				t.Errorf("expected 3 slots filled, got %d", n)
			}

			// Expected array state: ["hello", "", "a", "b", "c"]

			// Test ARLEN and ARCOUNT
			l, err := fsys.ARLen(key)
			if err != nil || l != 5 {
				t.Errorf("ARLen expected 5, got %d (err: %v)", l, err)
			}

			c, err := fsys.ARCount(key)
			if err != nil || c != 4 {
				t.Errorf("ARCount expected 4, got %d (err: %v)", c, err)
			}

			// Test ARDEL
			d, err := fsys.ARDel(key, 1, 2)
			if err != nil {
				t.Fatalf("ARDel error: %v", err)
			}
			if d != 1 { // index 1 was already empty, so only index 2 ("a") deleted
				t.Errorf("ARDel expected 1 element deleted, got %d", d)
			}

			// Expected array state: ["hello", "", "", "b", "c"]

			// Test ARGETRANGE
			vals, err := fsys.ARGetRange(key, 0, 4)
			if err != nil {
				t.Fatalf("ARGetRange error: %v", err)
			}
			if len(vals) != 5 || vals[0] != "hello" || vals[3] != "b" {
				t.Errorf("ARGetRange unexpected result: %v", vals)
			}

			// Test ARINSERT
			i, err := fsys.ARInsert(key, "d")
			if err != nil || i != 5 {
				t.Errorf("ARInsert expected index 5, got %d (err: %v)", i, err)
			}

			// Expected array state: ["hello", "", "", "b", "c", "d"]
			c, _ = fsys.ARCount(key)
			if c != 4 {
				t.Errorf("ARCount after insert expected 4, got %d", c)
			}

			// Test AROP (MATCH)
			res, err := fsys.AROp(key, 0, 10, "MATCH", dotpipStringPtr("b"))
			if err != nil || res.(int) != 1 {
				t.Errorf("AROp MATCH expected 1, got %v (err: %v)", res, err)
			}

			// Test ARDELRANGE
			d, err = fsys.ARDelRange(key, [2]int{3, 5})
			if err != nil || d != 3 {
				t.Errorf("ARDelRange expected 3 elements deleted, got %d (err: %v)", d, err)
			}

			// Test ARRING
			ringKey := dotpip.Key{"ring"}
			_, _ = fsys.ARRing(ringKey, 3, "1", "2")
			_, _ = fsys.ARRing(ringKey, 3, "3", "4")

			ringVals, _ := fsys.ARGetRange(ringKey, 0, 10)
			if len(ringVals) != 3 || ringVals[0] != "2" || ringVals[1] != "3" || ringVals[2] != "4" {
				t.Errorf("ARRing unexpected result: %v", ringVals)
			}
		})
	}
}

func dotpipStringPtr(s string) *string {
	return &s
}

func TestARGrep(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"grep_array"}
	_, _ = fsys.ARSet(key, 0, "apple", "banana", "cherry", "apricot", "blueberry")

	// Test EXACT
	res, _ := fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "EXACT", Value: "banana"}}, dotpip.ARGrepOptions{})
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("ARGrep EXACT error: %v", res)
	}

	// Test EXACT NoCase
	res, _ = fsys.ARGrep(key, "0", "4", []dotpip.ARGrepPredicate{{Type: "EXACT", Value: "BANANA"}}, dotpip.ARGrepOptions{NoCase: true})
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("ARGrep EXACT NoCase error: %v", res)
	}

	// Test MATCH
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "pp"}}, dotpip.ARGrepOptions{})
	if len(res) != 1 || res[0] != 0 {
		t.Errorf("ARGrep MATCH error: %v", res)
	}

	// Test MATCH NoCase
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "PP"}}, dotpip.ARGrepOptions{NoCase: true})
	if len(res) != 1 || res[0] != 0 {
		t.Errorf("ARGrep MATCH NoCase error: %v", res)
	}

	// Test GLOB
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "GLOB", Value: "a*"}}, dotpip.ARGrepOptions{})
	if len(res) != 2 || res[0] != 0 || res[1] != 3 {
		t.Errorf("ARGrep GLOB error: %v", res)
	}

	// Test GLOB NoCase
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "GLOB", Value: "A*"}}, dotpip.ARGrepOptions{NoCase: true})
	if len(res) != 2 || res[0] != 0 || res[1] != 3 {
		t.Errorf("ARGrep GLOB NoCase error: %v", res)
	}

	// Test RE
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "RE", Value: "^c.*"}}, dotpip.ARGrepOptions{})
	if len(res) != 1 || res[0] != 2 {
		t.Errorf("ARGrep RE error: %v", res)
	}

	// Test RE NoCase
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "RE", Value: "^C.*"}}, dotpip.ARGrepOptions{NoCase: true})
	if len(res) != 1 || res[0] != 2 {
		t.Errorf("ARGrep RE NoCase error: %v", res)
	}

	// Test AND
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "a"}, {Type: "MATCH", Value: "e"}}, dotpip.ARGrepOptions{And: true})
	if len(res) != 1 || res[0] != 0 {
		t.Errorf("ARGrep AND error: %v", res)
	}

	// Test OR
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "EXACT", Value: "apple"}, {Type: "EXACT", Value: "cherry"}}, dotpip.ARGrepOptions{})
	if len(res) != 2 || res[0] != 0 || res[1] != 2 {
		t.Errorf("ARGrep OR error: %v", res)
	}

	// Test Limit
	limit := 1
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "a"}}, dotpip.ARGrepOptions{Limit: &limit})
	if len(res) != 1 || res[0] != 0 {
		t.Errorf("ARGrep Limit error: %v", res)
	}

	// Test WithValues
	res, _ = fsys.ARGrep(key, "-", "+", []dotpip.ARGrepPredicate{{Type: "EXACT", Value: "banana"}}, dotpip.ARGrepOptions{WithValues: true})
	if len(res) != 2 || res[0] != 1 || res[1] != "banana" {
		t.Errorf("ARGrep WithValues error: %v", res)
	}

	// Test Reverse Range
	res, _ = fsys.ARGrep(key, "3", "0", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "a"}}, dotpip.ARGrepOptions{})
	if len(res) != 3 || res[0] != 3 || res[1] != 1 || res[2] != 0 {
		t.Errorf("ARGrep Reverse Range error: %v", res)
	}

	// Test Limit Reverse Range
	res, _ = fsys.ARGrep(key, "3", "0", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "a"}}, dotpip.ARGrepOptions{Limit: &limit})
	if len(res) != 1 || res[0] != 3 {
		t.Errorf("ARGrep Limit Reverse Range error: %v", res)
	}

	// Test WithValues Reverse Range
	res, _ = fsys.ARGrep(key, "3", "0", []dotpip.ARGrepPredicate{{Type: "MATCH", Value: "p"}}, dotpip.ARGrepOptions{WithValues: true})
	if len(res) != 4 || res[0] != 3 || res[1] != "apricot" || res[2] != 0 || res[3] != "apple" {
		t.Errorf("ARGrep WithValues Reverse Range error: %v", res)
	}
}

func TestARScan(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"scan_array"}
	_, _ = fsys.ARSet(key, 0, "0", "1", "2", "3", "4")

	res, _ := fsys.ARScan(key, 1, 3, nil)
	if len(res) != 6 || res[0] != 1 || res[1] != "1" || res[4] != 3 || res[5] != "3" {
		t.Errorf("ARScan error: %v", res)
	}

	res, _ = fsys.ARScan(key, 3, 1, nil)
	if len(res) != 6 || res[0] != 3 || res[1] != "3" || res[4] != 1 || res[5] != "1" {
		t.Errorf("ARScan reverse error: %v", res)
	}

	limit := 1
	res, _ = fsys.ARScan(key, 1, 3, &limit)
	if len(res) != 2 || res[0] != 1 || res[1] != "1" {
		t.Errorf("ARScan limit error: %v", res)
	}

	res, _ = fsys.ARScan(key, 3, 1, &limit)
	if len(res) != 2 || res[0] != 3 || res[1] != "3" {
		t.Errorf("ARScan reverse limit error: %v", res)
	}
}

func TestARMSetARMGet(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"mset_array"}

	n, _ := fsys.ARMSet(key, []dotpip.ARIndexValue{{Index: 1, Value: "a"}, {Index: 3, Value: "b"}})
	if n != 2 {
		t.Errorf("ARMSet expected 2, got %d", n)
	}

	res, _ := fsys.ARMGet(key, 0, 1, 2, 3, 4)
	if len(res) != 5 || res[0] != "" || res[1] != "a" || res[2] != "" || res[3] != "b" || res[4] != "" {
		t.Errorf("ARMGet error: %v", res)
	}
}

func TestAROp(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"op_array"}
	_, _ = fsys.ARSet(key, 0, "1", "2", "3", "notanumber")

	res, _ := fsys.AROp(key, 0, 3, "SUM", nil)
	if res.(string) != "6" {
		t.Errorf("AROp SUM error: %v", res)
	}

	res, _ = fsys.AROp(key, 0, 3, "MIN", nil)
	if res.(string) != "1" {
		t.Errorf("AROp MIN error: %v", res)
	}

	res, _ = fsys.AROp(key, 0, 3, "MAX", nil)
	if res.(string) != "3" {
		t.Errorf("AROp MAX error: %v", res)
	}

	key2 := dotpip.Key{"op_array_2"}
	_, _ = fsys.ARSet(key2, 0, "2", "6") // 0010, 0110

	res, _ = fsys.AROp(key2, 0, 1, "AND", nil)
	if res.(int) != 2 {
		t.Errorf("AROp AND error: %v", res)
	}

	res, _ = fsys.AROp(key2, 0, 1, "OR", nil)
	if res.(int) != 6 {
		t.Errorf("AROp OR error: %v", res)
	}

	res, _ = fsys.AROp(key2, 0, 1, "XOR", nil)
	if res.(int) != 4 {
		t.Errorf("AROp XOR error: %v", res)
	}

	res, _ = fsys.AROp(key, 0, 3, "USED", nil)
	if res.(int) != 4 {
		t.Errorf("AROp USED error: %v", res)
	}

	res, _ = fsys.AROp(key, 0, 3, "MATCH", dotpipStringPtr("2"))
	if res.(int) != 1 {
		t.Errorf("AROp MATCH error: %v", res)
	}
}

func TestARInfo(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"info_array"}
	_, _ = fsys.ARSet(key, 0, "a", "b", "c")

	res, _ := fsys.ARInfo(key, false)
	if res["length"].(int) != 3 {
		t.Errorf("ARInfo length error: %v", res)
	}
}

func TestARSeek(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()
	key := dotpip.Key{"seek_array"}
	idx, _ := fsys.ARSeek(key, 5)
	if idx != 5 {
		t.Errorf("ARSeek expected 5, got %d", idx)
	}
}
