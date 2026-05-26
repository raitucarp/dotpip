package fs

import (
	"encoding/json"
	"fmt"
)

type FileEncodeType string

const (
	JSON FileEncodeType = "json"
	YAML FileEncodeType = "yaml"
	TOML FileEncodeType = "toml"
	RAW  FileEncodeType = "raw"
)

func (f *fileSystem) EncodeType(typeName FileEncodeType) {
	f.encodeType = typeName
}

func (f *fileSystem) stringEncode(value string) (finalValue any, err error) {
	switch f.encodeType {
	case JSON:
		
		v, marshalErr := json.Marshal(value)
		err = marshalErr
		finalValue = v
	case YAML:
	case TOML:
	case RAW:
		finalValue = value
	default:
		return "", fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}

	return finalValue, err
}

func (f *fileSystem) stringDecode(value any) (v string, err error) {

	finalValue := ""
	switch f.encodeType {
	case JSON:
		err = json.Unmarshal(value.([]byte), &finalValue)
	case YAML:
	case TOML:
	case RAW:
		finalValue = value.(string)
	default:
		return "", fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}

	return finalValue, err
}
