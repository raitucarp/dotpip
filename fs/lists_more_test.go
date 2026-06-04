package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_lists_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mylist")

	// Test RPushX on missing list
	resX, err := dotfs.RPushX(key, "v1", "v2")
	assert.NoError(t, err)
	assert.Equal(t, 0, resX)

	// Create list
	_, _ = dotfs.RPush(key, "a", "b", "c")

	// Test RPushX on existing list
	resX, err = dotfs.RPushX(key, "v1", "v2")
	assert.NoError(t, err)
	assert.Equal(t, 5, resX)

	// LTrim
	err = dotfs.LTrim(key, 1, 3)
	assert.NoError(t, err)

	// LSet
	err = dotfs.LSet(key, 0, "x")
	assert.NoError(t, err)

	// LSet out of range
	err = dotfs.LSet(key, 99, "x")
	assert.Error(t, err)

	// LRem
	_, _ = dotfs.LPush(key, "x", "y", "x")
	resRem, err := dotfs.LRem(key, 2, "x")
	assert.NoError(t, err)
	assert.Equal(t, 2, resRem)
}

func TestArraysMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_arrays_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("myarray")

	// ARMSet
	resSet, err := dotfs.ARMSet(key, []dotpip.ARIndexValue{{Index: 0, Value: "a"}, {Index: 1, Value: "b"}})
	assert.NoError(t, err)
	assert.Equal(t, 2, resSet)

	// ARLastItems
	last, err := dotfs.ARLastItems(key, 2)
	assert.NoError(t, err)
	assert.Len(t, last, 2)
	assert.Equal(t, "a", last[0])
	assert.Equal(t, "b", last[1])

	// ARNext
	next, err := dotfs.ARNext(key)
	assert.NoError(t, err)
	assert.Equal(t, 2, next)

	// ARScan (basic)
	limit := 10
	scanRes, err := dotfs.ARScan(key, 0, 10, &limit)
	assert.NoError(t, err)
	assert.Greater(t, len(scanRes), 0)
}

func TestListsLRangeLIndexLInsert(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_lists_more2_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mylist2")
	_, _ = dotfs.RPush(key, "a", "b", "c")

	// LRange
	lrange, err := dotfs.LRange(key, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, lrange, 3)

	lrange2, err := dotfs.LRange(key, -2, -1)
	assert.NoError(t, err)
	assert.Len(t, lrange2, 2)
	assert.Equal(t, "b", lrange2[0])

	// LIndex
	lindex, err := dotfs.LIndex(key, -1)
	assert.NoError(t, err)
	assert.Equal(t, "c", lindex)

	// LInsert BEFORE
	resIns, err := dotfs.LInsert(key, dotpip.Before, "b", "x")
	assert.NoError(t, err)
	assert.Equal(t, 4, resIns)

	lindex2, _ := dotfs.LIndex(key, 1)
	assert.Equal(t, "x", lindex2)

	// LInsert AFTER
	resIns2, err := dotfs.LInsert(key, dotpip.After, "c", "y")
	assert.NoError(t, err)
	assert.Equal(t, 5, resIns2)

	// LInsert non-existent pivot
	resIns3, err := dotfs.LInsert(key, dotpip.Before, "nonexistent", "z")
	assert.NoError(t, err)
	assert.Equal(t, -1, resIns3)
}

func TestListsLPopRPopLMove(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_lists_more3_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key1 := dotpip.NewKey("list1")
	key2 := dotpip.NewKey("list2")
	_, _ = dotfs.RPush(key1, "a", "b", "c")

	// LPop count
	resLPop, err := dotfs.LPop(key1, 2)
	assert.NoError(t, err)
	assert.Len(t, resLPop, 2)
	assert.Equal(t, "a", resLPop[0])

	// RPop count
	_, _ = dotfs.RPush(key1, "d", "e")
	resRPop, err := dotfs.RPop(key1, 2)
	assert.NoError(t, err)
	assert.Len(t, resRPop, 2)
	assert.Equal(t, "e", resRPop[0])

	// RPop > length
	_, _ = dotfs.RPush(key1, "a")
	resRPop2, err := dotfs.RPop(key1, 100)
	assert.NoError(t, err)
	assert.Len(t, resRPop2, 2) // only c and a are left

	// LMove
	_, _ = dotfs.RPush(key1, "1", "2", "3")
	resMove, err := dotfs.LMove(key1, key2, dotpip.Right, dotpip.Left)
	assert.NoError(t, err)
	assert.Equal(t, "3", resMove)

	lindex, _ := dotfs.LIndex(key2, 0)
	assert.Equal(t, "3", lindex)
}

func TestListsLRemAndLRange(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_lists_lrem_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("list3")

	// Test LRange out of bounds
	_, _ = dotfs.RPush(key, "a")
	resLR, _ := dotfs.LRange(key, 0, 5)
	assert.Len(t, resLR, 1)

	// LRem count < 0
	_, _ = dotfs.RPush(key, "x", "y", "x")
	resRem, err := dotfs.LRem(key, -1, "x")
	assert.NoError(t, err)
	assert.Equal(t, 1, resRem)

	// LRem count == 0
	_, _ = dotfs.RPush(key, "z", "z")
	resRem2, err := dotfs.LRem(key, 0, "z")
	assert.NoError(t, err)
	assert.Equal(t, 2, resRem2)
}

func TestArraysScanOp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_arrays_scan_op_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("myarray_op")
	_, _ = dotfs.ARMSet(key, []dotpip.ARIndexValue{{Index: 0, Value: "a"}, {Index: 1, Value: "b"}})

	// AROp
	match := "b"
	resOp, err := dotfs.AROp(key, 0, 1, "SUM", &match)
	assert.NoError(t, err)
	assert.Nil(t, resOp)

	// ARGrep
	resGrep, err := dotfs.ARGrep(key, "a", "b", nil, dotpip.ARGrepOptions{})
	assert.NoError(t, err)
	assert.Nil(t, resGrep)

	// ARDelRange
	resDelRange, err := dotfs.ARDelRange(key, [2]int{0, 1})
	assert.NoError(t, err)
	assert.Equal(t, 2, resDelRange)

	// ARInfo
	_, _ = dotfs.ARMSet(key, []dotpip.ARIndexValue{{Index: 0, Value: "a"}})
	resInfo, err := dotfs.ARInfo(key, false)
	assert.NoError(t, err)
	assert.NotNil(t, resInfo)
}

func TestListsRPushXMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_lists_rpushx_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("myrpushx")

	res1, _ := dotfs.RPushX(key, "1")
	assert.Equal(t, 0, res1) // does not create missing list

	_, _ = dotfs.RPush(key, "a")

	res2, _ := dotfs.RPushX(key, "b", "c")
	assert.Equal(t, 3, res2)
}
