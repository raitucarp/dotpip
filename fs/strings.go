package fs

import (
	"dotpip"
	"fmt"

	"github.com/zeebo/xxh3"
)

func (f *fileSystem) Append(key dotpip.Key, value string) (appendedString int) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return 0
	}

	oldValue, err := f.formatter.StringDecode(content)
	if err != nil {
		return 0
	}

	newValue := oldValue + value
	_, err = f.Set(key, newValue)
	if err != nil {
		return 0
	}

	return len(newValue)
}

func (f *fileSystem) Get(key dotpip.Key) (result string, err error) {
	content, err := f.readFileByKey(key)

	if err != nil {
		return "", err
	}

	value, err := f.formatter.StringDecode(content)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (f *fileSystem) Set(key dotpip.Key, value string, options ...dotpip.SetOption) (result string, err error) {
	cmd := &dotpip.SetCommand{}
	for _, option := range options {
		option(cmd)
	}

	finalValue, err := f.formatter.StringEncode(value)
	if err != nil {
		return "", err
	}

	keyExists, err := f.Exists(key)

	if cmd.NX && keyExists[0] {
		return "", nil
	} else if cmd.XX && !keyExists[0] {
		return "", nil
	} else if cmd.IfEq != "" && (!keyExists[0] || value != cmd.IfEq) {
		return "", nil
	} else if cmd.IfNe != "" && (!keyExists[0] || value == cmd.IfNe) {
		return "", nil
	} else if cmd.IfDeq != "" && (keyExists[0] || value != cmd.IfDeq) {
		return "", nil
	} else if cmd.IfDne != "" && (keyExists[0] || value == cmd.IfDne) {
		return "", nil
	}

	err = f.writeFileByKey(key, finalValue.([]byte))
	if err != nil {
		return "", err
	}

	return value, nil
}

func (f *fileSystem) Digest(key dotpip.Key) (hexHash string, err error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return "", err
	}

	xxhash := xxh3.Hash(content)
	hexHash = fmt.Sprintf("%x", xxhash)

	return hexHash, nil
}
