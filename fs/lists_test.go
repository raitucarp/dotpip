package fs

import (
	"dotpip"
	"testing"
)

func TestListPushPop(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist")

	// LPush
	n, err := dotfs.LPush(lk, "a", "b", "c")
	if err != nil || n != 3 {
		t.Errorf("LPush expected 3, got %d, err %v", n, err)
	}

	// LPop
	vals, err := dotfs.LPop(lk, 1)
	if err != nil || len(vals) != 1 || vals[0] != "c" {
		t.Errorf("LPop expected [c], got %v, err %v", vals, err)
	}

	// RPush
	n, err = dotfs.RPush(lk, "d", "e")
	if err != nil || n != 4 { // "b", "a", "d", "e"
		t.Errorf("RPush expected 4, got %d", n)
	}

	// RPop
	vals, err = dotfs.RPop(lk, 2)
	if err != nil || len(vals) != 2 || vals[0] != "e" || vals[1] != "d" {
		t.Errorf("RPop expected [e d], got %v", vals)
	}

	// LLen
	n, err = dotfs.LLen(lk)
	if err != nil || n != 2 {
		t.Errorf("LLen expected 2, got %d", n)
	}
}

func TestListRange(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist2")

	_, _ = dotfs.RPush(lk, "one", "two", "three", "four")

	vals, _ := dotfs.LRange(lk, 0, 2)
	if len(vals) != 3 || vals[0] != "one" || vals[1] != "two" || vals[2] != "three" {
		t.Errorf("LRange expected [one two three], got %v", vals)
	}

	vals, _ = dotfs.LRange(lk, -2, -1)
	if len(vals) != 2 || vals[0] != "three" || vals[1] != "four" {
		t.Errorf("LRange negative indices expected [three four], got %v", vals)
	}
}

func TestListIndexSetInsertTrim(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist3")

	_, _ = dotfs.RPush(lk, "a", "b", "c")

	// LIndex
	val, err := dotfs.LIndex(lk, 1)
	if err != nil || val != "b" {
		t.Errorf("LIndex expected 'b', got '%v'", val)
	}

	// LSet
	err = dotfs.LSet(lk, 1, "B")
	if err != nil {
		t.Errorf("LSet error: %v", err)
	}

	val, _ = dotfs.LIndex(lk, 1)
	if val != "B" {
		t.Errorf("LIndex after LSet expected 'B', got '%v'", val)
	}

	// LInsert
	n, err := dotfs.LInsert(lk, dotpip.Before, "B", "A2")
	if err != nil || n != 4 {
		t.Errorf("LInsert expected length 4, got %d", n)
	}

	vals, _ := dotfs.LRange(lk, 0, -1)
	if len(vals) != 4 || vals[0] != "a" || vals[1] != "A2" || vals[2] != "B" || vals[3] != "c" {
		t.Errorf("LInsert contents wrong: %v", vals)
	}

	// LTrim
	err = dotfs.LTrim(lk, 1, -2)
	if err != nil {
		t.Errorf("LTrim error: %v", err)
	}

	n, _ = dotfs.LLen(lk)
	if n != 2 {
		t.Errorf("LTrim expected length 2, got %d", n)
	}
}

func TestListRemPosMove(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist4")

	_, _ = dotfs.RPush(lk, "a", "b", "c", "b", "d")

	// LRem
	n, err := dotfs.LRem(lk, 1, "b")
	if err != nil || n != 1 {
		t.Errorf("LRem expected 1 removed, got %d", n)
	}

	vals, _ := dotfs.LRange(lk, 0, -1)
	if len(vals) != 4 || vals[0] != "a" || vals[1] != "c" || vals[2] != "b" || vals[3] != "d" {
		t.Errorf("LRem contents wrong: %v", vals)
	}

	// LPos
	pos, err := dotfs.LPos(lk, "b")
	if err != nil || len(pos) != 1 || pos[0] != 2 {
		t.Errorf("LPos expected [2], got %v", pos)
	}

	// LMove
	lkDest := dotpip.NewKey("mylist4_dest")
	val, err := dotfs.LMove(lk, lkDest, dotpip.Right, dotpip.Left)
	if err != nil || val != "d" {
		t.Errorf("LMove expected 'd', got '%s', err: %v", val, err)
	}

	val, _ = dotfs.LIndex(lkDest, 0)
	if val != "d" {
		t.Errorf("LMove dest contents wrong: %s", val)
	}

	n, _ = dotfs.LLen(lk)
	if n != 3 {
		t.Errorf("LMove src expected length 3, got %d", n)
	}
}

func TestListEncodings(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}

	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			testFS := NewFileSystem("../data")
			testFS.EncodeType(enc)
			dfs := dotpip.New(testFS)
			_ = dfs.FlushAll()

			key := dotpip.NewKey("test_list_enc")

			_, err := dfs.RPush(key, "item1", "item2")
			if err != nil {
				t.Fatalf("Failed to RPush with %s encoding: %v", enc, err)
			}

			getVal, err := dfs.LRange(key, 0, -1)
			if err != nil {
				t.Fatalf("Failed to LRange with %s encoding: %v", enc, err)
			}

			if len(getVal) != 2 || getVal[0] != "item1" || getVal[1] != "item2" {
				t.Errorf("Expected [item1 item2], got %v with %s encoding", getVal, enc)
			}
		})
	}
}

func TestPushX(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist5")

	n, err := dotfs.LPushX(lk, "a")
	if err != nil || n != 0 {
		t.Errorf("LPushX expected 0, got %d", n)
	}

	_, _ = dotfs.RPush(lk, "b")

	n, err = dotfs.LPushX(lk, "a")
	if err != nil || n != 2 {
		t.Errorf("LPushX expected 2, got %d", n)
	}

	val, _ := dotfs.LIndex(lk, 0)
	if val != "a" {
		t.Errorf("LPushX expected 'a' at 0, got %s", val)
	}
}

func TestLMoveSameKey(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist_same")

	_, _ = dotfs.RPush(lk, "a", "b", "c")

	val, err := dotfs.LMove(lk, lk, dotpip.Right, dotpip.Left)
	if err != nil || val != "c" {
		t.Errorf("LMove same key expected 'c', got '%s', err: %v", val, err)
	}

	vals, _ := dotfs.LRange(lk, 0, -1)
	if len(vals) != 3 || vals[0] != "c" || vals[1] != "a" || vals[2] != "b" {
		t.Errorf("LMove same key contents wrong: %v", vals)
	}
}

func TestLPosMaxLenNegative(t *testing.T) {
	_ = dotfs.FlushAll()
	lk := dotpip.NewKey("mylist_pos")

	_, _ = dotfs.RPush(lk, "a", "b", "c", "b", "d")

	pos, err := dotfs.LPos(lk, "b", dotpip.WithLPosRank(-1), dotpip.WithLPosMaxLen(3))
	if err != nil || len(pos) != 1 || pos[0] != 3 {
		t.Errorf("LPos negative rank with max len expected [3], got %v, err %v", pos, err)
	}
}
