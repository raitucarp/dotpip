package dotpip

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type fileSystem struct {
	pathRoot   string
	formatter  *DataTypeFormatter
	encodeType EncodeType
}

func (f *fileSystem) Get(key Key) (result string, err error) {
	content, err := os.ReadFile(f.keyToAbsoluteFilePath(key))
	if err != nil {
		return "", err
	}

	value, err := f.formatter.StringDecode(content)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (f *fileSystem) Set(key Key, value string) (result string, err error) {

	// finalValue := []byte{}

	finalValue, err := f.formatter.StringEncode(value)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(f.keyToAbsoluteFilePath(key), finalValue.([]byte), 0644)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (f *fileSystem) FlushAll() (err error) {

	_, err = os.Stat(f.pathRoot)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(f.pathRoot, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		err = os.RemoveAll(f.pathRoot)
		if err != nil {
			return err
		}

		err = os.MkdirAll(f.pathRoot, 0755)
		if err != nil {
			return err
		}
	}

	return
}

func (f *fileSystem) Formatter(formatter DataTypeFormatter) {
	f.formatter = &formatter
}

type EncodeType string

const (
	JSON EncodeType = "json"
	YAML EncodeType = "yaml"
	TOML EncodeType = "toml"
)

func (f *fileSystem) EncodeType(typeName EncodeType) {
	f.encodeType = typeName
}

func (f *fileSystem) keyToAbsoluteFilePath(key Key) string {
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

func (f *fileSystem) stringDecode(value any) (string, error) {
	var err error
	finalValue := ""
	switch f.encodeType {
	case JSON:
		err = json.Unmarshal(value.([]byte), &finalValue)
		// Implement JSON encoding logic here
		// For example, you could marshal the value to JSON format
		// and return the resulting string.
	case YAML:
		// Implement YAML encoding logic here
	case TOML:
		// Implement TOML encoding logic here
	default:
		// Handle unsupported encoding types if necessary
	}

	return finalValue, err
}

func (f *fileSystem) stringEncode(value string) (any, error) {
	var err error
	finalValue := []byte{}
	switch f.encodeType {
	case JSON:
		finalValue, err = json.Marshal(value)
		// Implement JSON encoding logic here
		// For example, you could marshal the value to JSON format
		// and return the resulting byte slice.
	case YAML:
		// Implement YAML encoding logic here
	case TOML:
		// Implement TOML encoding logic here
	default:
		// Handle unsupported encoding types if necessary
	}

	return finalValue, err
}

func FileSystem(pathRoot string) *fileSystem {
	f := fileSystem{
		pathRoot:   pathRoot,
		formatter:  &DataTypeFormatter{},
		encodeType: JSON,
	}

	f.formatter.StringEncode = f.stringEncode
	f.formatter.StringDecode = f.stringDecode

	return &f
}
