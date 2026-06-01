package fs

import (
	"dotpip"
	"errors"
	"fmt"
	"strconv"
	"time"

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

	if cmd.NX && exists {
		return returnOldValueIfGet()
	} else if cmd.XX && !exists {
		return returnOldValueIfGet()
	} else if cmd.IfEq != "" && (!exists || oldValue != cmd.IfEq) {
		return returnOldValueIfGet()
	} else if cmd.IfNe != "" && (!exists || oldValue == cmd.IfNe) {
		return returnOldValueIfGet()
	} else if cmd.IfDeq != "" && (exists || oldValue != cmd.IfDeq) {
		return returnOldValueIfGet()
	} else if cmd.IfDne != "" && (exists || oldValue == cmd.IfDne) {
		return returnOldValueIfGet()
	}

	err = f.writeFileByKey(key, finalValue.([]byte))
	if err != nil {
		return "", err
	}

	if !cmd.KeepTTL {
		f.removeExByKey(key) // Removing existing ttl first

		var ttlMs int64
		if cmd.Ex > 0 {
			ttlMs = int64(cmd.Ex) * 1000
		} else if cmd.Px > 0 {
			ttlMs = int64(cmd.Px)
		} else if cmd.ExAt > 0 {
			nowMs := time.Now().UnixMilli()
			ttlMs = int64(cmd.ExAt)*1000 - nowMs
		} else if cmd.PxAt > 0 {
			nowMs := time.Now().UnixMilli()
			ttlMs = int64(cmd.PxAt) - nowMs
		}

		if ttlMs > 0 {
			expireAt := time.Now().UnixMilli() + ttlMs
			expireContent := strconv.FormatInt(expireAt, 10)
			f.writeExByKey(key, []byte(expireContent))
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

func (f *fileSystem) Digest(key dotpip.Key) (hexHash string, err error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return "", err
	}

	xxhash := xxh3.Hash(content)
	hexHash = fmt.Sprintf("%x", xxhash)

	return hexHash, nil
}

func (f *fileSystem) StrLen(key dotpip.Key) int {
	val, err := f.Get(key)
	if err != nil {
		return 0
	}
	return len(val)
}

func (f *fileSystem) Incr(key dotpip.Key) (int, error) {
	return f.IncrBy(key, 1)
}

func (f *fileSystem) IncrBy(key dotpip.Key, increment int) (int, error) {
	val, err := f.Get(key)
	if err != nil {
		// If key does not exist or error reading, default to 0
		val = "0"
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("ERR value is not an integer or out of range")
	}

	num += increment
	newValStr := strconv.Itoa(num)
	_, setErr := f.Set(key, newValStr)
	if setErr != nil {
		return 0, setErr
	}

	return num, nil
}

func (f *fileSystem) IncrByFloat(key dotpip.Key, increment float64) (float64, error) {
	val, err := f.Get(key)
	if err != nil {
		val = "0"
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New("ERR value is not a valid float")
	}

	num += increment

	// Format the float consistently, e.g. using %g
	newValStr := strconv.FormatFloat(num, 'g', -1, 64)
	_, setErr := f.Set(key, newValStr)
	if setErr != nil {
		return 0, setErr
	}

	return num, nil
}

func (f *fileSystem) Decr(key dotpip.Key) (int, error) {
	return f.DecrBy(key, 1)
}

func (f *fileSystem) DecrBy(key dotpip.Key, decrement int) (int, error) {
	return f.IncrBy(key, -decrement)
}

func (f *fileSystem) GetDel(key dotpip.Key) (string, error) {
	val, err := f.Get(key)
	if err != nil {
		return "", err
	}
	f.Del(key)
	return val, nil
}

func (f *fileSystem) GetRange(key dotpip.Key, start int, end int) (string, error) {
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

func (f *fileSystem) SetRange(key dotpip.Key, offset int, value string) (int, error) {
	if offset < 0 {
		return 0, errors.New("ERR offset is out of range")
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
	_, setErr := f.Set(key, newVal)
	if setErr != nil {
		return 0, setErr
	}

	return len(newVal), nil
}

func (f *fileSystem) MGet(keys ...dotpip.Key) ([]string, error) {
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

func (f *fileSystem) MSet(kvs ...dotpip.KV) error {
	for _, kv := range kvs {
		_, err := f.Set(kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *fileSystem) MSetNX(kvs ...dotpip.KV) (bool, error) {
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
