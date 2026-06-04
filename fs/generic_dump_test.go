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
	defer os.RemoveAll(tmpDir)

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
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey")
	_, _ = dotfs.Set(key, "val")

	// Migrate
	err = dotfs.Migrate("host", 6379, key, 0, 0)
	assert.Error(t, err)

	// Sort
	resSort, err := dotfs.Sort(key)
	assert.Error(t, err) // It says not implemented usually
	_ = resSort
}
