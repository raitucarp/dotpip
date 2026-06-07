package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"
)

func TestFormatter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_fmt_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	df := dotpip.DataTypeFormatter{}
	dotfs.Formatter(df)
}
