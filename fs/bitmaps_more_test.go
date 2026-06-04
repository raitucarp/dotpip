package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitmapsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_bitmaps_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	k1 := dotpip.NewKey("bm1")
	k2 := dotpip.NewKey("bm2")
	kout := dotpip.NewKey("bmout")

	// SetBit
	res, err := dotfs.SetBit(k1, 0, 1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, res, 0)

	res, err = dotfs.SetBit(k1, 10, 1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, res, 0)

	_, _ = dotfs.SetBit(k2, 0, 1)
	_, _ = dotfs.SetBit(k2, 5, 1)

	// BitCount
	c1, err := dotfs.BitCount(k1, 0, -1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, c1, 0)

	// BitOp AND
	opRes, err := dotfs.BitOp(dotpip.BitOpAnd, kout, k1, k2)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, opRes, 0) // length of longest string

	// BitOp OR
	opRes, err = dotfs.BitOp(dotpip.BitOpOr, kout, k1, k2)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, opRes, 0)

	// BitOp XOR
	opRes, err = dotfs.BitOp(dotpip.BitOpXor, kout, k1, k2)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, opRes, 0)

	// BitOp NOT
	opRes, err = dotfs.BitOp(dotpip.BitOpNot, kout, k1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, opRes, 0)

	// BitPos
	pos, err := dotfs.BitPos(k1, 1, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 0, pos)

	pos0, err := dotfs.BitPos(k1, 0, 0, -1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, pos0, 0) // first bit that is 0

	// BitField
	bfRes, err := dotfs.BitField(k1,
		"GET", "u8", "0",
		"SET", "i8", "8", 10,
		"INCRBY", "u4", "16", 2,
	)
	assert.NoError(t, err)
	assert.Len(t, bfRes, 3)
}

func TestBitmapsBitPosBitField(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_bitmaps_more2_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	k1 := dotpip.NewKey("bm3")
	// Add some bits
	_, _ = dotfs.SetBit(k1, 0, 1)
	_, _ = dotfs.SetBit(k1, 8, 1)

	// BitPos
	pos, err := dotfs.BitPos(k1, 1, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 8, pos)

	// BitField
	bfRes, err := dotfs.BitField(k1,
		"SET", "u4", "0", 15,
		"GET", "u4", "0",
		"INCRBY", "i4", "4", -1,
	)
	assert.NoError(t, err)
	assert.Len(t, bfRes, 3)

	// OVERFLOW FAIL
	bfRes2, err := dotfs.BitField(k1,
		"OVERFLOW", "FAIL",
		"INCRBY", "u4", "0", 100, // this should fail overflow
	)
	assert.NoError(t, err)
	assert.Len(t, bfRes2, 1)
	assert.Nil(t, bfRes2[0]) // Overflow fail returns nil for the operation

	// OVERFLOW SAT
	bfRes3, err := dotfs.BitField(k1,
		"OVERFLOW", "SAT",
		"INCRBY", "u4", "0", 100, // this should saturate
	)
	assert.NoError(t, err)
	assert.Len(t, bfRes3, 1)
}

func TestBitmapsBitPosMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_bitmaps_more3_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("bm_bitpos")

	// Empty string bitpos 0
	pos, _ := dotfs.BitPos(key, 0, 0, -1)
	assert.Equal(t, 0, pos)

	// Empty string bitpos 1
	pos2, _ := dotfs.BitPos(key, 1, 0, -1)
	assert.Equal(t, -1, pos2)

	_, _ = dotfs.SetBit(key, 10, 1)

	// Out of bounds start
	pos3, _ := dotfs.BitPos(key, 1, 100, -1)
	assert.Equal(t, -1, pos3)

	// Negative start/end
	pos4, _ := dotfs.BitPos(key, 1, -1, -1)
	assert.GreaterOrEqual(t, pos4, 0)
}

func TestBitmapsBitCountMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_bitmaps_more4_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("bm_bitcount")

	// Missing key
	c, _ := dotfs.BitCount(key, 0, -1)
	assert.Equal(t, 0, c)

	_, _ = dotfs.SetBit(key, 0, 1)
	_, _ = dotfs.SetBit(key, 8, 1)

	// Negative offsets
	c2, _ := dotfs.BitCount(key, -1, -1)
	assert.GreaterOrEqual(t, c2, 0)

	// Start > End
	c3, _ := dotfs.BitCount(key, 10, 5)
	assert.Equal(t, 0, c3)

	// Negative start < 0
	c4, _ := dotfs.BitCount(key, -100, -1)
	assert.GreaterOrEqual(t, c4, 0)

	// GetBit out of bounds
	bitOut, _ := dotfs.GetBit(key, 100)
	assert.Equal(t, 0, bitOut)

	// BitField parse type error
	_, errBF := dotfs.BitField(key, "GET", "invalid", "0")
	assert.Error(t, errBF)

	// BitField out of bounds
	bfResOut, _ := dotfs.BitField(key, "GET", "u8", "1000")
	assert.Len(t, bfResOut, 1)
	assert.NotNil(t, bfResOut[0])
}
