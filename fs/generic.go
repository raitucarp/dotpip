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

// Copy copies a key to another.
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

// Del deletes keys.
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

// Exists checks if a key exists.
func (f *FileSystem) Exists(keys ...dotpip.Key) ([]bool, error) {
	results := make([]bool, len(keys))
	for i, key := range keys {
		exist, err := f.checkExistByKey(key)
		results[i] = exist && err == nil
	}

	return results, nil
}

// FlushAll flushes all keys.
func (f *FileSystem) FlushAll() (err error) {
	err = os.RemoveAll(f.pathRoot)
	if err != nil {
		return err
	}

	return os.MkdirAll(f.pathRoot, 0755)
}

// Rename renames a key.
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

// RenameNX renames a key only if the new key does not exist.
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

// Keys returns all keys matching a pattern.
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

// Type returns the string representation of the type of the value stored at key.
func (f *FileSystem) Type(key dotpip.Key) (dotpip.ObjectType, error) {
	// Determine type by trying to decode it
	// First, check if it exists
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return dotpip.ObjectTypeNone, err
	}
	if !exist {
		return dotpip.ObjectTypeNone, nil
	}

	content, err := f.readFileByKey(key)
	if err != nil {
		return dotpip.ObjectTypeNone, err
	}

	// Try decoding in order of complexity/specificity
	if _, err := f.formatter.StringDecode(content); err == nil {
		return dotpip.ObjectTypeString, nil
	}
	if _, err := f.formatter.HashDecode(content); err == nil {
		return dotpip.ObjectTypeHash, nil
	}
	if _, err := f.formatter.ListDecode(content); err == nil {
		return dotpip.ObjectTypeList, nil
	}
	if _, err := f.formatter.SetDecode(content); err == nil {
		return dotpip.ObjectTypeSet, nil
	}
	if _, err := f.formatter.SortedSetDecode(content); err == nil {
		return dotpip.ObjectTypeZSet, nil
	}
	if _, err := f.formatter.StreamDecode(content); err == nil {
		return dotpip.ObjectTypeStream, nil
	}

	return dotpip.ObjectTypeUnknown, nil
}

// RandomKey returns a random key from the currently selected database.
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

// Touch alters the last access time of a key.
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

// Unlink deletes a key without blocking.
func (f *FileSystem) Unlink(keys ...dotpip.Key) int {
	// Unlink is usually async, but for fs backend we can just use Del
	return f.Del(keys...)
}

// Dump serializes the value stored at key in a dotpip-specific format.
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

// Restore creates a key associated with a value that is obtained by deserializing the provided serialized value.
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
		return errors.New(string(dotpip.ErrMsgBusyKey))
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

// Sort returns or stores the elements contained in the list, set or sorted set at key.
func (f *FileSystem) Sort(key dotpip.Key) ([]string, error) {
	typ, err := f.Type(key)
	if err != nil {
		return nil, err
	}

	var items []string
	switch typ {
	case dotpip.ObjectTypeNone:
		return []string{}, nil
	case dotpip.ObjectTypeList:
		items, err = f.LRange(key, 0, -1)
	case dotpip.ObjectTypeSet:
		items, err = f.SMembers(key)
	case dotpip.ObjectTypeZSet:
		items, err = f.ZRange(key, "-inf", "+inf")
	default:
		return nil, errors.New(string(dotpip.ErrMsgWrongType))
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

// Scan incrementally iterates the keys space.
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
				if err != nil || string(typ) != cmd.Type {
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

// Move moves a key to another dotpip instance.
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

// Wait blocks the current client until all the previous write commands are successfully transferred.
func (f *FileSystem) Wait(_ int, _ int) (int, error) {
	// No replication supported in FS backend
	return 0, nil
}

// WaitAOF blocks the current client until all the previous write commands are synced to the AOF.
func (f *FileSystem) WaitAOF(_ int, _ int, _ int) (int, int, error) {
	// No AOF/replication supported in FS backend
	return 0, 0, nil
}

// DBSize returns the number of keys in the selected database.
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

// ObjectEncoding returns the internal encoding for the object associated with key.
func (f *FileSystem) ObjectEncoding(key dotpip.Key) (dotpip.ObjectEncoding, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return "", err
	}
	// Simplified representation since everything is serialized as defined by f.encodeType
	switch f.encodeType {
	case JSON:
		return dotpip.ObjectEncodingJSON, nil
	case YAML:
		return dotpip.ObjectEncodingYAML, nil
	case TOML:
		return dotpip.ObjectEncodingTOML, nil
	default:
		return dotpip.ObjectEncodingRAW, nil
	}
}

// ObjectFreq returns the logarithmic access frequency counter of the object stored at key.
func (f *FileSystem) ObjectFreq(key dotpip.Key) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return 0, err
	}
	return 0, errors.New(string(dotpip.ErrMsgLFUEviction))
}

// ObjectIdletime returns the time in seconds since the last access to the value stored at key.
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

// ObjectRefcount returns the reference count of the object stored at key.
func (f *FileSystem) ObjectRefcount(key dotpip.Key) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return 0, err
	}
	return 1, nil // Refcount is typically 1 since we don't share objects
}

// Migrate atomically transfers a key from a source dotpip instance to a destination dotpip instance.
func (f *FileSystem) Migrate(_ string, _ int, _ dotpip.Key, _ dotpip.DotPip, _ int, _ ...dotpip.MigrateOption) error {
	return errors.New(string(dotpip.ErrMsgMigrateNotSupported))
}
