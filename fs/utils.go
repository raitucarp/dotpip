package fs

import (
	"dotpip"
	"errors"
	"os"
	"path/filepath"
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

func (f *fileSystem) readFileByKey(key dotpip.Key) (content []byte, err error) {
	content, err = os.ReadFile(f.keyToAbsoluteFilePath(key))
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *fileSystem) writeFileByKey(key dotpip.Key, content []byte) (err error) {
	err = os.WriteFile(f.keyToAbsoluteFilePath(key), content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileSystem) checkExistByKey(key dotpip.Key) (exist bool, err error) {
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

func (f *fileSystem) readExByKey(key dotpip.Key) ([]byte, error) {
	return os.ReadFile(f.keyToAbsoluteExFilePath(key))
}

func (f *fileSystem) removeFileByKey(key dotpip.Key) (err error) {
	keyFileName := f.keyToAbsoluteFilePath(key)

	exist, err := f.checkExistByKey(key)
	if !exist {
		return
	}

	err = os.Remove(keyFileName)
	if err != nil {
		return
	}

	return
}
