package fs

import (
	"encoding/json"
	"fmt"

	"dotpip"

	yaml "github.com/goccy/go-yaml"
	toml "github.com/pelletier/go-toml/v2"
)

type FileEncodeType string

const (
	JSON FileEncodeType = "json"
	YAML FileEncodeType = "yaml"
	TOML FileEncodeType = "toml"
	RAW  FileEncodeType = "raw"
)

func (f *FileSystem) EncodeType(typeName FileEncodeType) {
	f.encodeType = typeName
}

func (f *FileSystem) stringEncode(value string) (finalValue any, err error) {
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
		finalValue = []byte(value)
	default:
		return "", fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}

	return finalValue, err
}

func (f *FileSystem) listEncode(value []any) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// Pelletier's go-toml requires a map at the root level.
		return toml.Marshal(map[string][]any{"value": value})
	case RAW:
		// For RAW, we might just store a JSON byte slice or newline delimited.
		// For simplicity, falling back to JSON.
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) bitmapEncode(value []uint) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		return toml.Marshal(map[string][]uint{"value": value})
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) bitmapDecode(value any) ([]uint, error) {
	var finalValue []uint
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		var wrap map[string][]uint
		err := toml.Unmarshal(value.([]byte), &wrap)
		if err == nil {
			finalValue = wrap["value"]
		}
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) listDecode(value any) ([]any, error) {
	var finalValue []any
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		var wrap map[string][]any
		err := toml.Unmarshal(value.([]byte), &wrap)
		if err == nil {
			finalValue = wrap["value"]
		}
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) hashEncode(value map[string]string) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// toml marshaller can handle top-level map, but to be safe and consistent, we can wrap it or just use it directly since it is a map[string]string
		return toml.Marshal(value)
	case RAW:
		// For RAW, we might just store a JSON byte slice or fallback.
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) hashDecode(value any) (map[string]string, error) {
	finalValue := make(map[string]string)
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) stringDecode(value any) (v string, err error) {

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
		if b, ok := value.([]byte); ok {
			finalValue = string(b)
		} else if s, ok := value.(string); ok {
			finalValue = s
		} else {
			return "", fmt.Errorf("RAW stringDecode expected []byte or string, got %T", value)
		}
	default:
		return "", fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}

	return finalValue, err
}

func (f *FileSystem) setEncode(value map[string]any) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// toml marshaller can handle top-level map, but we'll wrap it to be consistent with others or just use it directly since it is a map[string]any
		return toml.Marshal(value)
	case RAW:
		// For RAW, we might just store a JSON byte slice or fallback.
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) setDecode(value any) (map[string]any, error) {
	finalValue := make(map[string]any)
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) sortedSetEncode(value map[string]float64) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		return toml.Marshal(value)
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) sortedSetDecode(value any) (map[string]float64, error) {
	finalValue := make(map[string]float64)
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) hyperLogLogEncode(value []byte) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		return toml.Marshal(map[string][]byte{"value": value})
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) hyperLogLogDecode(value any) ([]byte, error) {
	var finalValue []byte
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		var wrap map[string][]byte
		err := toml.Unmarshal(value.([]byte), &wrap)
		if err == nil {
			finalValue = wrap["value"]
		}
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) streamEncode(value dotpip.Stream) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// Pelletier's go-toml requires a map at the root level.
		return toml.Marshal(value)
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) streamDecode(value any) (dotpip.Stream, error) {
	var finalValue dotpip.Stream
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return finalValue, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) JSONEncode(value any) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// TOML requires map[string]any at top level, for a general JSON value this might fail if it's a scalar or slice.
		// A safe fallback would be wrapping it
		if _, ok := value.(map[string]any); ok {
			return toml.Marshal(value)
		}
		return toml.Marshal(map[string]any{"value": value})
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *FileSystem) JSONDecode(value any) (any, error) {
	var finalValue any
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		if err == nil {
			if m, ok := finalValue.(map[string]any); ok {
				if v, hasValue := m["value"]; hasValue && len(m) == 1 {
					finalValue = v
				}
			}
		}
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}
