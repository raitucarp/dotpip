package fs

import (
	"dotpip"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zeebo/xxh3"
)

// Append appends a value to a key.
func (f *FileSystem) Append(key dotpip.Key, value string) (appendedString int) {
	defer func() {
		if appendedString > 0 {
			f.emitKeyspaceEvent(key, "append", '$')
		}
	}()

	content, err := f.readFileByKey(key)
	if err != nil {
		return 0
	}

	oldValue, err := f.formatter.StringDecode(content)
	if err != nil {
		return 0
	}

	newValue := oldValue + value
	_, err = f.internalSet(key, newValue)
	if err != nil {
		return 0
	}

	return len(newValue)
}

// Get gets the value of a key.
func (f *FileSystem) Get(key dotpip.Key) (result string, err error) {
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

func (f *FileSystem) internalSet(key dotpip.Key, value string, options ...dotpip.SetOption) (result string, err error) {

	cmd := &dotpip.SetCommand{}
	for _, option := range options {
		option(cmd)
	}

	finalValue, err := f.formatter.StringEncode(value)
	if err != nil {
		return "", err
	}

	keyExists, _ := f.Exists(key)
	exists := keyExists[0]

	var oldValue string
	if cmd.Get || cmd.IfEq != "" || cmd.IfNe != "" || cmd.IfDeq != "" || cmd.IfDne != "" {
		if exists {
			oldVal, getErr := f.Get(key)
			if getErr == nil {
				oldValue = oldVal
			}
		}
	}

	returnOldValueIfGet := func() (string, error) {
		if cmd.Get {
			if exists {
				return oldValue, nil
			}
			return "", nil
		}
		return "", nil
	}

	switch {
	case cmd.NX && exists:
		return returnOldValueIfGet()
	case cmd.XX && !exists:
		return returnOldValueIfGet()
	case cmd.IfEq != "" && (!exists || oldValue != cmd.IfEq):
		return returnOldValueIfGet()
	case cmd.IfNe != "" && (!exists || oldValue == cmd.IfNe):
		return returnOldValueIfGet()
	case cmd.IfDeq != "" && (exists || oldValue != cmd.IfDeq):
		return returnOldValueIfGet()
	case cmd.IfDne != "" && (exists || oldValue == cmd.IfDne):
		return returnOldValueIfGet()
	}

	err = f.writeFileByKey(key, finalValue.([]byte))

	if err != nil {
		return "", err
	}

	if !cmd.KeepTTL {
		_ = f.removeExByKey(key) // Removing existing ttl first

		var ttlMs int64
		switch {
		case cmd.Ex > 0:
			ttlMs = int64(cmd.Ex) * 1000
		case cmd.Px > 0:
			ttlMs = int64(cmd.Px)
		case cmd.ExAt > 0:
			nowMs := time.Now().UnixMilli()
			ttlMs = int64(cmd.ExAt)*1000 - nowMs
		case cmd.PxAt > 0:
			nowMs := time.Now().UnixMilli()
			ttlMs = int64(cmd.PxAt) - nowMs
		}

		if ttlMs > 0 {
			expireAt := time.Now().UnixMilli() + ttlMs
			expireContent := strconv.FormatInt(expireAt, 10)
			_ = f.writeExByKey(key, []byte(expireContent))
			dataPath := f.keyToAbsoluteFilePath(key)
			f.setExpiration(dataPath, expireAt)
		} else {
			dataPath := f.keyToAbsoluteFilePath(key)
			f.unsetExpiration(dataPath)
		}
	}

	if cmd.Get {
		if exists {
			return oldValue, nil
		}
		return "", nil
	}

	return value, nil
}

// Digest gets the hex hash of a value.
func (f *FileSystem) Digest(key dotpip.Key) (hexHash string, err error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return "", err
	}

	xxhash := xxh3.Hash(content)
	hexHash = fmt.Sprintf("%x", xxhash)

	return hexHash, nil
}

// StrLen returns the string length.
func (f *FileSystem) StrLen(key dotpip.Key) int {
	val, err := f.Get(key)
	if err != nil {
		return 0
	}
	return len(val)
}

// Incr increments the integer value of a key by one.
func (f *FileSystem) Incr(key dotpip.Key) (int, error) {
	return f.IncrBy(key, 1)
}

// IncrBy increments the integer value of a key by a number.
func (f *FileSystem) IncrBy(key dotpip.Key, increment int) (ret int, err error) {
	defer func() {
		if err == nil {
			f.emitKeyspaceEvent(key, "incrby", '$')
		}
	}()

	val, err := f.Get(key)
	if err != nil {
		// If key does not exist or error reading, default to 0
		val = "0"
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New(string(dotpip.ErrMsgValueNotInt))
	}

	num += increment
	newValStr := strconv.Itoa(num)
	_, setErr := f.internalSet(key, newValStr)
	if setErr != nil {
		return 0, setErr
	}

	return num, nil
}

// IncrByFloat increments the float value of a key by a number.
func (f *FileSystem) IncrByFloat(key dotpip.Key, increment float64) (ret float64, err error) {
	defer func() {
		if err == nil {
			f.emitKeyspaceEvent(key, "incrbyfloat", '$')
		}
	}()

	val, err := f.Get(key)
	if err != nil {
		val = "0"
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New(string(dotpip.ErrMsgValueNotFloat))
	}

	num += increment

	// Format the float consistently, e.g. using %g
	newValStr := strconv.FormatFloat(num, 'g', -1, 64)
	_, setErr := f.internalSet(key, newValStr)
	if setErr != nil {
		return 0, setErr
	}

	return num, nil
}

// Decr decrements the integer value of a key by one.
func (f *FileSystem) Decr(key dotpip.Key) (int, error) {
	return f.DecrBy(key, 1)
}

// DecrBy decrements the integer value of a key by a number.
func (f *FileSystem) DecrBy(key dotpip.Key, decrement int) (int, error) {
	return f.IncrBy(key, -decrement)
}


// GetDel gets the value of a key and deletes it.
func (f *FileSystem) GetDel(key dotpip.Key) (string, error) {
	val, err := f.Get(key)
	if err != nil {
		return "", err
	}
	f.Del(key)
	return val, nil
}


// GetRange gets a substring of the string stored at a key.
func (f *FileSystem) GetRange(key dotpip.Key, start int, end int) (string, error) {
	val, err := f.Get(key)
	if err != nil {
		return "", err
	}

	length := len(val)
	if length == 0 {
		return "", nil
	}

	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}

	if start >= length {
		return "", nil
	}
	if end >= length {
		end = length - 1
	}

	if start > end {
		return "", nil
	}

	return val[start : end+1], nil
}

// SetRange overwrites part of a string.
// Set sets the string value of a key.
func (f *FileSystem) SetRange(key dotpip.Key, offset int, value string) (ret int, err error) {
	defer func() {
		if err == nil {
			f.emitKeyspaceEvent(key, "setrange", '$')
		}
	}()

	if offset < 0 {
		return 0, errors.New(string(dotpip.ErrMsgOffsetOutOfRange))
	}

	val, err := f.Get(key)
	if err != nil {
		val = ""
	}

	// Pad with zero bytes if offset is greater than length
	if offset > len(val) {
		padLen := offset - len(val)
		padBytes := make([]byte, padLen)
		val += string(padBytes)
	}

	// Calculate new length
	newLen := offset + len(value)
	if newLen < len(val) {
		newLen = len(val)
	}

	newBytes := make([]byte, newLen)
	copy(newBytes, val)
	copy(newBytes[offset:], value)

	newVal := string(newBytes)
	_, setErr := f.internalSet(key, newVal)
	if setErr != nil {
		return 0, setErr
	}

	return len(newVal), nil
}

// MGet returns the values of all specified keys.
func (f *FileSystem) MGet(keys ...dotpip.Key) ([]string, error) {
	res := make([]string, len(keys))
	for i, key := range keys {
		val, err := f.Get(key)
		if err != nil {
			// In Redis MGET returns nil for non-existing keys.
			// Since our Go type is []string, returning an empty string for missing.
			res[i] = ""
		} else {
			res[i] = val
		}
	}
	return res, nil
}

// MSet sets the given keys to their respective values.
func (f *FileSystem) MSet(kvs ...dotpip.KV) error {
	for _, kv := range kvs {
		_, err := f.Set(kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// MSetNX sets the given keys to their respective values, only if none of the keys exist.
func (f *FileSystem) MSetNX(kvs ...dotpip.KV) (bool, error) {
	// First check if any keys exist
	keys := make([]dotpip.Key, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}

	exists, err := f.Exists(keys...)
	if err != nil {
		return false, err
	}

	for _, exist := range exists {
		if exist {
			return false, nil // At least one key exists, MSETNX fails
		}
	}

	// Now set all keys
	for _, kv := range kvs {
		_, err := f.Set(kv.Key, kv.Value)
		if err != nil {
			// Rollback would be ideal, but for simplicity we return err
			return false, err
		}
	}

	return true, nil
}

// Set sets the string value of a key.
func (f *FileSystem) Set(key dotpip.Key, value string, options ...dotpip.SetOption) (result string, err error) {
	res, err := f.internalSet(key, value, options...)
	if err == nil {
		f.emitKeyspaceEvent(key, "set", '$')
	}
	return res, err
}
