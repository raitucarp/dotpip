package fs

import (
	"dotpip"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func (f *fileSystem) keyToAbsoluteFilePath(key dotpip.Key) string {
	finalPath := []string{f.pathRoot}
	finalPath = append(finalPath, key...)

	dataPath := filepath.Join(finalPath...)

	switch f.encodeType {
	case JSON:
		dataPath += ".json"
	case YAML:
		dataPath += ".yaml"
	case TOML:
		dataPath += ".toml"
	default:
		dataPath += ""
	}

	return dataPath
}

func (f *fileSystem) keyToAbsoluteExFilePath(key dotpip.Key) string {
	finalPath := []string{f.pathRoot}
	finalPath = append(finalPath, key...)

	dataPath := filepath.Join(finalPath...)

	return dataPath + ".ex"
}

func (f *fileSystem) isExpired(key dotpip.Key) bool {
	dataPath := f.keyToAbsoluteFilePath(key)
	expireAt, hasTTL := f.getExpiration(dataPath)
	if hasTTL {
		if time.Now().UnixMilli() >= expireAt {
			_ = f.removeFileByPath(dataPath)
			_ = f.removeExByPath(dataPath + ".ex")
			f.unsetExpiration(dataPath)
			return true
		}
	}
	return false
}

func (f *fileSystem) readFileByKey(key dotpip.Key) (content []byte, err error) {
	if f.isExpired(key) {
		return nil, os.ErrNotExist
	}
	content, err = os.ReadFile(f.keyToAbsoluteFilePath(key))
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *fileSystem) writeFileByKey(key dotpip.Key, content []byte) (err error) {
	path := f.keyToAbsoluteFilePath(key)
	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileSystem) checkExistByKey(key dotpip.Key) (exist bool, err error) {
	if f.isExpired(key) {
		return false, nil
	}
	keyFileName := f.keyToAbsoluteFilePath(key)
	_, err = os.Stat(keyFileName)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func (f *fileSystem) writeExByKey(key dotpip.Key, content []byte) error {
	return os.WriteFile(f.keyToAbsoluteExFilePath(key), content, 0644)
}

func (f *fileSystem) removeExByKey(key dotpip.Key) error {
	err := os.Remove(f.keyToAbsoluteExFilePath(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (f *fileSystem) removeFileByKey(key dotpip.Key) (err error) {
	keyFileName := f.keyToAbsoluteFilePath(key)

	err = os.Remove(keyFileName)
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func (f *fileSystem) readExByPath(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (f *fileSystem) removeFileByPath(path string) error {
	return os.Remove(path)
}

func (f *fileSystem) removeExByPath(path string) error {
	return os.Remove(path)
}

func (f *fileSystem) setExpiration(dataPath string, expireAt int64) {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()
	f.expirations[dataPath] = expireAt
}

func (f *fileSystem) unsetExpiration(dataPath string) {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()
	delete(f.expirations, dataPath)
}

func (f *fileSystem) getExpiration(dataPath string) (int64, bool) {
	f.expMutex.RLock()
	defer f.expMutex.RUnlock()
	val, ok := f.expirations[dataPath]
	return val, ok
}
