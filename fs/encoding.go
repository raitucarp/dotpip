package fs

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
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
		v, marshalErr := yaml.Marshal(value)
		err = marshalErr
		finalValue = v
	case TOML:
		v, marshalErr := toml.Marshal(map[string]string{"value": value})
		err = marshalErr
		finalValue = v
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
		err = yaml.Unmarshal(value.([]byte), &finalValue)
	case TOML:
		var wrap map[string]string
		err = toml.Unmarshal(value.([]byte), &wrap)
		if err == nil {
			finalValue = wrap["value"]
		}
	case RAW:
		finalValue = value.(string)
	default:
		return "", fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}

	return finalValue, err
}
