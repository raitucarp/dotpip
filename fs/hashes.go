package fs

import (
	"dotpip"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

func (f *FileSystem) readHash(key dotpip.Key) (map[string]string, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	return f.formatter.HashDecode(content)
}

func (f *FileSystem) writeHash(key dotpip.Key, hash map[string]string) error {
	if len(hash) == 0 {
		err := f.removeFileByKey(key)
		if err == nil {
			f.emitKeyspaceEvent(key, "del", 'g')
		}
		return err
	}

	encoded, err := f.formatter.HashEncode(hash)
	if err != nil {
		return err
	}

	return f.writeFileByKey(key, encoded.([]byte))
}

// HDel deletes fields from a hash.
func (f *FileSystem) HDel(key dotpip.Key, fields ...string) (int, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}
	if len(hash) == 0 {
		return 0, nil
	}

	deleted := 0
	var deletedFields []string
	for _, field := range fields {
		if _, exists := hash[field]; exists {
			delete(hash, field)
			deleted++
			deletedFields = append(deletedFields, field)
		}
	}

	if deleted > 0 {
		err = f.writeHash(key, hash)
		if err == nil {
			f.emitSubkeyEvent(key, "hdel", 'h', deletedFields)
		}
		if err != nil {
			return 0, err
		}
	}

	return deleted, nil
}

// HExists checks if a hash field exists.
func (f *FileSystem) HExists(key dotpip.Key, field string) (bool, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return false, err
	}
	if len(hash) == 0 {
		return false, nil
	}
	_, exists := hash[field]
	return exists, nil
}

// HGet gets a hash field value.
func (f *FileSystem) HGet(key dotpip.Key, field string) (string, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return "", err
	}
	if len(hash) == 0 {
		return "", nil // Redis returns nil for non-existing field
	}
	return hash[field], nil
}

// HGetAll gets all hash fields and values.
// HGet gets a hash field value.
// HGetAll gets all hash fields and values.
func (f *FileSystem) HGetAll(key dotpip.Key) (map[string]string, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// HIncrBy increments a hash field.
func (f *FileSystem) HIncrBy(key dotpip.Key, field string, increment int) (int, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}

	valStr, exists := hash[field]
	var current int
	if exists {
		current, err = strconv.Atoi(valStr)
		if err != nil {
			return 0, err // ERR hash value is not an integer
		}
	}

	current += increment
	hash[field] = strconv.Itoa(current)

	err = f.writeHash(key, hash)
	if err == nil {
		f.emitKeyspaceEvent(key, "hincrby", 'h')
		f.emitSubkeyEvent(key, "hincrby", 'h', []string{field})
	}
	if err != nil {
		return 0, err
	}

	return current, nil
}

// HIncrByFloat increments a hash field by float.
func (f *FileSystem) HIncrByFloat(key dotpip.Key, field string, increment float64) (float64, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}

	valStr, exists := hash[field]
	var current float64
	if exists {
		current, err = strconv.ParseFloat(valStr, 64)
		if err != nil {
			return 0, err // ERR hash value is not a float
		}
	}

	current += increment
	// Redis standard format
	hash[field] = strconv.FormatFloat(current, 'f', -1, 64)

	err = f.writeHash(key, hash)
	if err == nil {
		f.emitKeyspaceEvent(key, "hincrbyfloat", 'h')
		f.emitSubkeyEvent(key, "hincrbyfloat", 'h', []string{field})
	}
	if err != nil {
		return 0, err
	}

	return current, nil
}

// HKeys gets all hash fields.
func (f *FileSystem) HKeys(key dotpip.Key) ([]string, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}
	return keys, nil
}

// HLen gets hash length.
func (f *FileSystem) HLen(key dotpip.Key) (int, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}
	return len(hash), nil
}

// HMGet gets multiple hash fields.
func (f *FileSystem) HMGet(key dotpip.Key, fields ...string) ([]string, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return nil, err
	}

	results := make([]string, len(fields))
	for i, field := range fields {
		if val, exists := hash[field]; exists {
			results[i] = val
		} else {
			results[i] = "" // Redis returns nil for missing fields
		}
	}
	return results, nil
}

// HRandField gets random hash fields.
func (f *FileSystem) HRandField(key dotpip.Key, count int, options ...dotpip.HRandFieldOption) ([]string, error) {
	cmd := &dotpip.HRandFieldCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	hash, err := f.readHash(key)
	if err != nil {
		return nil, err
	}
	if len(hash) == 0 {
		return nil, nil
	}

	// Extract keys
	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}

	// Implement logic for count
	// If count is not specified (default count is 1 for simplicity here, but Redis has HRANDFIELD key [count [WITHVALUES]])
	// Note: We're taking count as an argument. Let's say count == 0 means 1 if no options, but if count is passed...
	// Standard Redis: HRANDFIELD key -> 1 element, HRANDFIELD key count -> count elements

	actualCount := count
	allowDuplicates := false
	if count < 0 {
		actualCount = -count
		allowDuplicates = true
	}

	var selectedKeys []string
	if allowDuplicates {
		for i := 0; i < actualCount; i++ {
			idx := rand.Intn(len(keys))
			selectedKeys = append(selectedKeys, keys[idx])
		}
	} else {
		if actualCount > len(keys) {
			actualCount = len(keys)
		}
		// Shuffle keys for unique random elements
		rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
		selectedKeys = keys[:actualCount]
	}

	if cmd.WithValues {
		results := make([]string, 0, len(selectedKeys)*2)
		for _, k := range selectedKeys {
			results = append(results, k, hash[k])
		}
		return results, nil
	}

	return selectedKeys, nil
}

// HSet sets hash fields.
func (f *FileSystem) HSet(key dotpip.Key, values map[string]string) (int, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}

	added := 0
	for k, v := range values {
		if _, exists := hash[k]; !exists {
			added++
		}
		hash[k] = v
	}

	err = f.writeHash(key, hash)
	if err == nil {
		f.emitKeyspaceEvent(key, "hset", 'h')

		var fields []string
		for k := range values {
			fields = append(fields, k)
		}
		f.emitSubkeyEvent(key, "hset", 'h', fields)
	}
	if err != nil {
		return 0, err
	}

	return added, nil
}

// HSetNX sets a hash field if it does not exist.
func (f *FileSystem) HSetNX(key dotpip.Key, field string, value string) (bool, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return false, err
	}

	if _, exists := hash[field]; exists {
		return false, nil
	}

	hash[field] = value
	err = f.writeHash(key, hash)
	if err == nil {
		f.emitKeyspaceEvent(key, "hset", 'h')
		f.emitSubkeyEvent(key, "hset", 'h', []string{field})
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// HStrLen gets hash field string length.
func (f *FileSystem) HStrLen(key dotpip.Key, field string) (int, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return 0, err
	}

	val, exists := hash[field]
	if !exists {
		return 0, nil
	}

	return len(val), nil
}

// HVals gets all hash values.
func (f *FileSystem) HVals(key dotpip.Key) ([]string, error) {
	hash, err := f.readHash(key)
	if err != nil {
		return nil, err
	}

	vals := make([]string, 0, len(hash))
	for _, v := range hash {
		vals = append(vals, v)
	}

	return vals, nil
}

// HScan iterates hash fields.
func (f *FileSystem) HScan(key dotpip.Key, cursor uint64, options ...dotpip.ScanOption) (uint64, map[string]string, error) {
	cmd := &dotpip.ScanCommand{Count: 10}
	for _, option := range options {
		option(cmd)
	}

	hash, err := f.readHash(key)
	if err != nil {
		return 0, nil, err
	}

	var allKeys []string
	for k := range hash {
		if cmd.Match != "" && !matchPattern(cmd.Match, k) {
			continue
		}
		allKeys = append(allKeys, k)
	}

	sort.Strings(allKeys)

	if cursor >= uint64(len(allKeys)) {
		return 0, map[string]string{}, nil
	}

	end := cursor + uint64(cmd.Count)
	nextCursor := end
	if end >= uint64(len(allKeys)) {
		end = uint64(len(allKeys))
		nextCursor = 0
	}

	result := make(map[string]string)
	for _, k := range allKeys[cursor:end] {
		result[k] = hash[k]
	}

	return nextCursor, result, nil
}
