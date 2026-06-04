package fs

import (
	"dotpip"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (f *FileSystem) Copy(source dotpip.Key, destination dotpip.Key, options ...dotpip.CopyOption) int {
	cmd := &dotpip.CopyCommand{}
	for _, option := range options {
		option(cmd)
	}

	if cmd.Destination != nil {
		content, err := f.Get(source)
		if err != nil {
			return 0
		}

		_, err = cmd.Destination.Set(destination, content)
		if err != nil {
			return 0
		}

		return 1
	}

	content, err := f.readFileByKey(source)
	if err != nil {
		return 0
	}

	if cmd.Replace {
		err = f.removeFileByKey(destination)
		if err != nil {
			return 0
		}
	}

	err = f.writeFileByKey(destination, content)
	if err == nil {
		f.emitKeyspaceEvent(destination, "copy_to", 'g')
	}
	if err != nil {
		return 0
	}

	return 1
}

func (f *FileSystem) Del(keys ...dotpip.Key) int {
	count := 0
	for _, key := range keys {
		err := f.removeFileByKey(key)
		if err == nil {
			f.emitKeyspaceEvent(key, "del", 'g')
		}
		if err != nil {
			return count
		}
		count++
	}

	return count
}

func (f *FileSystem) Exists(keys ...dotpip.Key) ([]bool, error) {
	results := make([]bool, len(keys))
	for i, key := range keys {
		exist, err := f.checkExistByKey(key)
		results[i] = exist && err == nil
	}

	return results, nil
}

func (f *FileSystem) FlushAll() (err error) {
	err = os.RemoveAll(f.pathRoot)
	if err != nil {
		return err
	}

	return os.MkdirAll(f.pathRoot, 0755)
}

func (f *FileSystem) Rename(key dotpip.Key, newKey dotpip.Key) error {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return err
	}
	if !exist {
		return os.ErrNotExist
	}

	oldPath := f.keyToAbsoluteFilePath(key)
	newPath := f.keyToAbsoluteFilePath(newKey)

	// ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return err
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	// Rename expiration if it exists
	oldExPath := f.keyToAbsoluteExFilePath(key)
	newExPath := f.keyToAbsoluteExFilePath(newKey)
	if _, err := os.Stat(oldExPath); err == nil {
		_ = os.Rename(oldExPath, newExPath)
		f.expMutex.Lock()
		if expireAt, ok := f.expirations[oldPath]; ok {
			delete(f.expirations, oldPath)
			f.expirations[newPath] = expireAt
		}
		f.expMutex.Unlock()
	}

	f.emitKeyspaceEvent(key, "rename_from", 'g')
	f.emitKeyspaceEvent(newKey, "rename_to", 'g')

	return nil
}

func (f *FileSystem) RenameNX(key dotpip.Key, newKey dotpip.Key) (bool, error) {
	exist, err := f.checkExistByKey(newKey)
	if err != nil {
		return false, err
	}
	if exist {
		return false, nil
	}
	err = f.Rename(key, newKey)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *FileSystem) Keys(pattern string) ([]dotpip.Key, error) {
	var keys []dotpip.Key
	err := filepath.Walk(f.pathRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".ex") && !strings.HasPrefix(info.Name(), ".pubsub_") {
			// Extract key from path
			relPath, err := filepath.Rel(f.pathRoot, path)
			if err != nil {
				return err
			}

			// Strip extension based on encodeType
			ext := filepath.Ext(relPath)
			if ext == ".json" || ext == ".yaml" || ext == ".toml" {
				relPath = relPath[:len(relPath)-len(ext)]
			}

			// Match pattern
			keyStr := strings.ReplaceAll(relPath, string(filepath.Separator), ":") // use a standard delimiter for glob matching
			if matchPattern(pattern, keyStr) {
				keys = append(keys, dotpip.NewKeyWithDelimiter(relPath, string(filepath.Separator)))
			}
		}
		return nil
	})
	return keys, err
}

func (f *FileSystem) Type(key dotpip.Key) (string, error) {
	// Determine type by trying to decode it
	// First, check if it exists
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return "none", err
	}
	if !exist {
		return "none", nil
	}

	content, err := f.readFileByKey(key)
	if err != nil {
		return "none", err
	}

	// Try decoding in order of complexity/specificity
	if _, err := f.formatter.StringDecode(content); err == nil {
		return "string", nil
	}
	if _, err := f.formatter.HashDecode(content); err == nil {
		return "hash", nil
	}
	if _, err := f.formatter.ListDecode(content); err == nil {
		return "list", nil
	}
	if _, err := f.formatter.SetDecode(content); err == nil {
		return "set", nil
	}
	if _, err := f.formatter.SortedSetDecode(content); err == nil {
		return "zset", nil
	}
	if _, err := f.formatter.StreamDecode(content); err == nil {
		return "stream", nil
	}

	return "unknown", nil
}

func (f *FileSystem) RandomKey() (dotpip.Key, error) {
	var keys []dotpip.Key
	err := filepath.Walk(f.pathRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".ex") && !strings.HasPrefix(info.Name(), ".pubsub_") {
			relPath, err := filepath.Rel(f.pathRoot, path)
			if err != nil {
				return err
			}
			ext := filepath.Ext(relPath)
			if ext == ".json" || ext == ".yaml" || ext == ".toml" {
				relPath = relPath[:len(relPath)-len(ext)]
			}
			keys = append(keys, dotpip.NewKeyWithDelimiter(relPath, string(filepath.Separator)))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil // Redis returns nil when DB is empty
	}

	// Pick a pseudo-random key
	// time.Now().UnixNano()
	idx := time.Now().UnixNano() % int64(len(keys))
	return keys[idx], nil
}

func (f *FileSystem) Touch(keys ...dotpip.Key) (int, error) {
	count := 0
	now := time.Now().Local()
	for _, key := range keys {
		exist, err := f.checkExistByKey(key)
		if exist && err == nil {
			path := f.keyToAbsoluteFilePath(key)
			if err := os.Chtimes(path, now, now); err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (f *FileSystem) Unlink(keys ...dotpip.Key) int {
	// Unlink is usually async, but for fs backend we can just use Del
	return f.Del(keys...)
}

func (f *FileSystem) Dump(key dotpip.Key) ([]byte, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return nil, err // Redis returns nil if key doesn't exist
	}

	content, err := f.readFileByKey(key)
	if err != nil {
		return nil, err
	}

	// Simplest serialization for fs dump is just the raw file bytes
	// Real Redis uses RDB format, but we just need it to be restorable
	return content, nil
}

func (f *FileSystem) Restore(key dotpip.Key, ttl int, serializedValue []byte, options ...dotpip.RestoreOption) error {
	cmd := &dotpip.RestoreCommand{}
	for _, option := range options {
		option(cmd)
	}

	exist, err := f.checkExistByKey(key)
	if err != nil {
		return err
	}

	if exist && !cmd.Replace {
		return errors.New("BUSYKEY Target key name already exists")
	}

	if err := f.writeFileByKey(key, serializedValue); err != nil {
		return err
	}

	if ttl > 0 {
		var expireAt int64
		if cmd.AbsTTL {
			expireAt = int64(ttl)
		} else {
			expireAt = time.Now().UnixMilli() + int64(ttl)
		}
		f.setExpiration(f.keyToAbsoluteFilePath(key), expireAt)
		expireContent := strconv.FormatInt(expireAt, 10)
		_ = f.writeExByKey(key, []byte(expireContent))
	}

	return nil
}

func (f *FileSystem) Sort(key dotpip.Key) ([]string, error) {
	typ, err := f.Type(key)
	if err != nil {
		return nil, err
	}

	var items []string
	switch typ {
	case "none":
		return []string{}, nil
	case "list":
		items, err = f.LRange(key, 0, -1)
	case "set":
		items, err = f.SMembers(key)
	case "zset":
		items, err = f.ZRange(key, "-inf", "+inf")
	default:
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if err != nil {
		return nil, err
	}

	sort.SliceStable(items, func(i, j int) bool {
		f1, err1 := strconv.ParseFloat(items[i], 64)
		f2, err2 := strconv.ParseFloat(items[j], 64)
		if err1 == nil && err2 == nil {
			return f1 < f2
		}
		return items[i] < items[j]
	})

	return items, nil
}

func (f *FileSystem) Scan(cursor uint64, options ...dotpip.ScanOption) (uint64, []dotpip.Key, error) {
	cmd := &dotpip.ScanCommand{Count: 10}
	for _, option := range options {
		option(cmd)
	}

	var allKeys []dotpip.Key
	err := filepath.Walk(f.pathRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".ex") && !strings.HasPrefix(info.Name(), ".pubsub_") {
			relPath, err := filepath.Rel(f.pathRoot, path)
			if err != nil {
				return err
			}
			ext := filepath.Ext(relPath)
			if ext == ".json" || ext == ".yaml" || ext == ".toml" {
				relPath = relPath[:len(relPath)-len(ext)]
			}

			keyStr := strings.ReplaceAll(relPath, string(filepath.Separator), ":")
			key := dotpip.NewKeyWithDelimiter(relPath, string(filepath.Separator))

			if cmd.Match != "" && !matchPattern(cmd.Match, keyStr) {
				return nil
			}

			if cmd.Type != "" {
				typ, err := f.Type(key)
				if err != nil || typ != cmd.Type {
					return nil
				}
			}

			allKeys = append(allKeys, key)
		}
		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	if cursor >= uint64(len(allKeys)) {
		return 0, []dotpip.Key{}, nil
	}

	end := cursor + uint64(cmd.Count)
	nextCursor := end
	if end >= uint64(len(allKeys)) {
		end = uint64(len(allKeys))
		nextCursor = 0
	}

	return nextCursor, allKeys[cursor:end], nil
}

func (f *FileSystem) Move(key dotpip.Key, db dotpip.DotPip) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, nil
	}

	dbExist, err := db.Exists(key)
	if err != nil {
		return 0, err
	}
	if len(dbExist) > 0 && dbExist[0] {
		return 0, nil // Target exists
	}

	dump, err := f.Dump(key)
	if err != nil {
		return 0, err
	}

	ttl, _ := f.TTL(key)
	if ttl < 0 {
		ttl = 0
	}

	err = db.Restore(key, int(ttl), dump)
	if err != nil {
		return 0, err
	}

	f.Del(key)
	return 1, nil
}

func (f *FileSystem) Wait(_ int, _ int) (int, error) {
	// No replication supported in FS backend
	return 0, nil
}

func (f *FileSystem) WaitAOF(_ int, _ int, _ int) (int, int, error) {
	// No AOF/replication supported in FS backend
	return 0, 0, nil
}

func (f *FileSystem) DBSize() (int, error) {
	count := 0
	err := filepath.Walk(f.pathRoot, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".ex") && !strings.HasPrefix(info.Name(), ".pubsub_") {
			count++
		}
		return nil
	})
	return count, err
}

func (f *FileSystem) ObjectEncoding(key dotpip.Key) (string, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return "", err
	}
	// Simplified representation since everything is serialized as defined by f.encodeType
	switch f.encodeType {
	case JSON:
		return "json", nil
	case YAML:
		return "yaml", nil
	case TOML:
		return "toml", nil
	default:
		return "raw", nil
	}
}

func (f *FileSystem) ObjectFreq(key dotpip.Key) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return 0, err
	}
	return 0, errors.New("ERR LFU eviction not supported")
}

func (f *FileSystem) ObjectIdletime(key dotpip.Key) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return 0, err
	}

	path := f.keyToAbsoluteFilePath(key)
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	// Calculate idle time using modtime since we updated it on Touch, Get, etc. (we didn't actually update it on Get, but simplified)
	idle := time.Since(info.ModTime()).Seconds()
	return int(idle), nil
}

func (f *FileSystem) ObjectRefcount(key dotpip.Key) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return 0, err
	}
	return 1, nil // Refcount is typically 1 since we don't share objects
}

func (f *FileSystem) Migrate(_ string, _ int, _ dotpip.Key, _ dotpip.DotPip, _ int, _ ...dotpip.MigrateOption) error {
	return errors.New("MIGRATE is not supported in fs mode over network")
}
