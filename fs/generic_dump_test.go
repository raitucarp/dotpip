package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericDumpRestore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_dump_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey")
	_, _ = dotfs.Set(key, "val")

	// Dump
	dumpBytes, err := dotfs.Dump(key)
	assert.NoError(t, err)
	assert.NotNil(t, dumpBytes)

	// Dump non-existent
	dumpBytes2, err := dotfs.Dump(dotpip.NewKey("non"))
	assert.NoError(t, err)
	assert.Nil(t, dumpBytes2)

	// Restore
	key3 := dotpip.NewKey("mykey3")
	err = dotfs.Restore(key3, 0, dumpBytes)
	assert.NoError(t, err)

	// Restore with replace
	err = dotfs.Restore(key3, 0, dumpBytes, dotpip.WithRestoreReplace())
	assert.NoError(t, err)
}

func TestGenericMigrateSort(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_migrate_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey")
	_, _ = dotfs.Set(key, "val")

	// Migrate
	tmpDirDB2, _ := os.MkdirTemp("", "dotpip_generic_migrate_test_db2_")
	defer func() { _ = os.RemoveAll(tmpDirDB2) }()
	db2 := fs.NewFileSystem(tmpDirDB2)
	defer db2.Close()

	err = dotfs.Migrate("host", 6379, key, db2, 0)
	assert.Error(t, err)

	// Sort
	resSort, err := dotfs.Sort(key)
	assert.Error(t, err)
	assert.Nil(t, resSort)

	// Sort success
	lk := dotpip.NewKey("mykey_list")
	_, _ = dotfs.LPush(lk, "3", "1.5", "2", "10")
	resSort, err = dotfs.Sort(lk)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1.5", "2", "3", "10"}, resSort)
}

func TestGenericRestoreMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_generic_restore_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey_restore")
	keyTarget := dotpip.NewKey("mykey_target")
	_, _ = dotfs.Set(key, "val")

	dumpBytes, _ := dotfs.Dump(key)

	// Restore over existing key without replace
	_, _ = dotfs.Set(keyTarget, "val")
	err = dotfs.Restore(keyTarget, 0, dumpBytes)
	assert.Error(t, err)

	// Restore with TTL relative
	keyTTL := dotpip.NewKey("mykey_ttl")
	err = dotfs.Restore(keyTTL, 1000000, dumpBytes)
	assert.NoError(t, err)

	// Restore with TTL absolute
	keyAbs := dotpip.NewKey("mykey_abs")
	err = dotfs.Restore(keyAbs, 1999999999999, dumpBytes, dotpip.WithRestoreAbsTTL())
	assert.NoError(t, err)
}
