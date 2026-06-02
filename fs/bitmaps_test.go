package fs

import (
	"dotpip"
	"os"
	"testing"
)

func runBitmapTests(t *testing.T, f dotpip.DotPip) {
	// SetBit & GetBit
	key := dotpip.NewKey("mybitmap")

	// Default GetBit is 0
	val, err := f.GetBit(key, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected 0, got %d", val)
	}

	// SetBit to 1
	old, err := f.SetBit(key, 100, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if old != 0 {
		t.Fatalf("expected old 0, got %d", old)
	}

	// GetBit
	val, err = f.GetBit(key, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}

	// BitCount
	count, err := f.BitCount(key, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	_, _ = f.SetBit(key, 101, 1)
	count, err = f.BitCount(key, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}

	// BitPos
	pos, err := f.BitPos(key, 1, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// offset 100 is bit 4 of byte 12 (100 / 8 = 12, 100 % 8 = 4)
	if pos != 100 {
		t.Fatalf("expected bitpos 100, got %d", pos)
	}

	// BitOp
	key2 := dotpip.NewKey("mybitmap2")
	_, _ = f.SetBit(key2, 100, 1) // 100 is 1
	_, _ = f.SetBit(key2, 102, 1) // 102 is 1

	destKey := dotpip.NewKey("dest")
	opLen, err := f.BitOp(dotpip.BitOpAnd, destKey, key, key2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opLen != 13 { // byte 12 is max index (13 bytes total)
		t.Fatalf("expected oplen 13, got %d", opLen)
	}

	val, _ = f.GetBit(destKey, 100)
	if val != 1 {
		t.Fatalf("expected bit 100 to be 1 in dest, got %d", val)
	}
	val, _ = f.GetBit(destKey, 101)
	if val != 0 {
		t.Fatalf("expected bit 101 to be 0 in dest, got %d", val)
	}
}

func TestBitmap(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}

	for _, encoding := range encodings {
		t.Run(string(encoding), func(t *testing.T) {
			path := "./test_data_" + string(encoding)
			_ = os.MkdirAll(path, 0755)
			defer os.RemoveAll(path)

			fs := FileSystem(path)
			fs.EncodeType(encoding)

			runBitmapTests(t, fs)
		})
	}
}

// BitField test
func TestBitField(t *testing.T) {
	_ = os.MkdirAll("./test_bitfield", 0755)
	fs := FileSystem("./test_bitfield")
	fs.EncodeType(RAW) // Use RAW to avoid JSON invalid UTF-8 string encoding corruption
	defer os.RemoveAll("./test_bitfield")

	key := dotpip.NewKey("foo")

	// Test SET and GET
	res, err := fs.BitField(key, "SET", "i8", 0, 100)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(res) != 1 || res[0] != 0 {
		t.Fatalf("expected old val 0, got %v", res)
	}

	res, err = fs.BitField(key, "GET", "i8", 0)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(res) != 1 || res[0] != 100 {
		t.Fatalf("expected GET to return 100, got %v", res)
	}

	// Test INCRBY OVERFLOW WRAP
	res, err = fs.BitField(key, "OVERFLOW", "WRAP", "INCRBY", "i8", 0, 50)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	// 100 + 50 = 150 -> WRAP for i8 is 150 - 256 = -106
	if len(res) != 1 || res[0] != -106 {
		t.Fatalf("expected INCRBY to wrap to -106, got %v", res)
	}

	// Test INCRBY OVERFLOW SAT
	res, err = fs.BitField(key, "OVERFLOW", "SAT", "INCRBY", "i8", 0, -50)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	// -106 - 50 = -156 -> SAT for i8 is -128
	if len(res) != 1 || res[0] != -128 {
		t.Fatalf("expected INCRBY to sat to -128, got %v", res)
	}
}
