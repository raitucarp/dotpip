package fs

import (
	"dotpip"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type fileSystem struct {
	pathRoot   string
	formatter  *dotpip.DataTypeFormatter
	encodeType FileEncodeType

	expirations map[string]int64
	expMutex    sync.RWMutex
	expStop     chan struct{}
}

func FileSystem(pathRoot string) *fileSystem {
	f := fileSystem{
		pathRoot:    pathRoot,
		formatter:   &dotpip.DataTypeFormatter{},
		encodeType:  JSON,
		expirations: make(map[string]int64),
		expStop:     make(chan struct{}),
	}

	f.formatter.StringEncode = f.stringEncode
	f.formatter.StringDecode = f.stringDecode
	f.formatter.HashEncode = f.hashEncode
	f.formatter.HashDecode = f.hashDecode
	f.formatter.ListEncode = f.listEncode
	f.formatter.ListDecode = f.listDecode
	f.formatter.SetEncode = f.setEncode
	f.formatter.SetDecode = f.setDecode
	f.formatter.SortedSetEncode = f.sortedSetEncode
	f.formatter.SortedSetDecode = f.sortedSetDecode
	f.formatter.BitmapEncode = f.bitmapEncode
	f.formatter.BitmapDecode = f.bitmapDecode

	f.loadExpirations()
	go f.scanExpirations()

	return &f
}

func (f *fileSystem) Close() {
	if f.expStop != nil {
		close(f.expStop)
	}
}

func (f *fileSystem) loadExpirations() {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()

	filepath.WalkDir(f.pathRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(path, ".ex") {
			content, err := f.readExByPath(path)
			if err == nil {
				expireAt, err := strconv.ParseInt(string(content), 10, 64)
				if err == nil {
					// We use the absolute .ex file path or the relative key as map key.
					// Let's use the key path relative to pathRoot.
					// For simplicity, we can use the absolute path of the data file as the map key.
					dataPath := strings.TrimSuffix(path, ".ex")
					f.expirations[dataPath] = expireAt
				}
			}
		}
		return nil
	})
}

func (f *fileSystem) scanExpirations() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-f.expStop:
			return
		case <-ticker.C:
			f.processExpirations()
		}
	}
}

func (f *fileSystem) processExpirations() {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()

	now := time.Now().UnixMilli()
	for dataPath, expireAt := range f.expirations {
		if now >= expireAt {
			// Expired!
			f.removeFileByPath(dataPath)
			f.removeExByPath(dataPath + ".ex")
			delete(f.expirations, dataPath)
		}
	}
}
