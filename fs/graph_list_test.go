package fs_test

import (
	"dotpip/fs"
	"path/filepath"
	"testing"
)

func TestGraphListError(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db\000invalid"))

	_, err := dotfs.GraphList()
	if err == nil {
		t.Log("Warning: GraphList didn't bubble error")
	}
}
