package fs

import (
	"dotpip"
	"strconv"
	"time"
)

func (f *fileSystem) applyExpireOptions(key dotpip.Key, expireAt int64, cmd *dotpip.ExpireCommand) (bool, error) {
	dataPath := f.keyToAbsoluteFilePath(key)
	currentExpireAt, hasTTL := f.getExpiration(dataPath)

	if cmd.NX && hasTTL {
		return false, nil
	}
	if cmd.XX && !hasTTL {
		return false, nil
	}
	if cmd.GT && (!hasTTL || expireAt <= currentExpireAt) {
		return false, nil
	}
	if cmd.LT && hasTTL && expireAt >= currentExpireAt {
		return false, nil
	}

	f.setExpiration(dataPath, expireAt)
	expireContent := strconv.FormatInt(expireAt, 10)
	f.writeExByKey(key, []byte(expireContent))

	return true, nil
}

func (f *fileSystem) Expire(key dotpip.Key, seconds int, options ...dotpip.ExpireOption) (bool, error) {
	return f.PExpire(key, seconds*1000, options...)
}

func (f *fileSystem) PExpire(key dotpip.Key, milliseconds int, options ...dotpip.ExpireOption) (bool, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return false, err
	}

	cmd := &dotpip.ExpireCommand{}
	for _, option := range options {
		option(cmd)
	}

	expireAt := time.Now().UnixMilli() + int64(milliseconds)
	return f.applyExpireOptions(key, expireAt, cmd)
}

func (f *fileSystem) ExpireAt(key dotpip.Key, timestamp int, options ...dotpip.ExpireOption) (bool, error) {
	return f.PExpireAt(key, timestamp*1000, options...)
}

func (f *fileSystem) PExpireAt(key dotpip.Key, timestamp int, options ...dotpip.ExpireOption) (bool, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return false, err
	}

	cmd := &dotpip.ExpireCommand{}
	for _, option := range options {
		option(cmd)
	}

	expireAt := int64(timestamp)

	// If timestamp is in the past, delete key immediately
	if expireAt <= time.Now().UnixMilli() {
	    // If options prevent expiration update, do nothing
	    dataPath := f.keyToAbsoluteFilePath(key)
	    currentExpireAt, hasTTL := f.getExpiration(dataPath)
	    if cmd.NX && hasTTL { return false, nil }
	    if cmd.XX && !hasTTL { return false, nil }
	    if cmd.GT && (!hasTTL || expireAt <= currentExpireAt) { return false, nil }
	    if cmd.LT && hasTTL && expireAt >= currentExpireAt { return false, nil }

	    f.Del(key)
	    return true, nil
	}

	return f.applyExpireOptions(key, expireAt, cmd)
}

func (f *fileSystem) ExpireTime(key dotpip.Key) (int64, error) {
    res, err := f.PExpireTime(key)
    if res > 0 {
        return res / 1000, err
    }
    return res, err
}

func (f *fileSystem) PExpireTime(key dotpip.Key) (int64, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return -2, err
	}
	if !exist {
		return -2, nil
	}

	dataPath := f.keyToAbsoluteFilePath(key)
	expireAt, hasTTL := f.getExpiration(dataPath)
	if !hasTTL {
		return -1, nil
	}

	return expireAt, nil
}

func (f *fileSystem) TTL(key dotpip.Key) (int64, error) {
    res, err := f.PTTL(key)
    if res > 0 {
        return (res + 500) / 1000, err
    }
    return res, err
}

func (f *fileSystem) PTTL(key dotpip.Key) (int64, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return -2, err
	}
	if !exist {
		return -2, nil
	}

	dataPath := f.keyToAbsoluteFilePath(key)
	expireAt, hasTTL := f.getExpiration(dataPath)
	if !hasTTL {
		return -1, nil
	}

	ttl := expireAt - time.Now().UnixMilli()
	if ttl < 0 {
	    // It's expired but not yet cleaned up
		return -2, nil
	}

	return ttl, nil
}

func (f *fileSystem) Persist(key dotpip.Key) (bool, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil || !exist {
		return false, err
	}

	dataPath := f.keyToAbsoluteFilePath(key)
	_, hasTTL := f.getExpiration(dataPath)
	if !hasTTL {
		return false, nil
	}

	f.unsetExpiration(dataPath)
	f.removeExByKey(key)
	return true, nil
}
