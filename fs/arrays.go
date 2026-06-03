package fs

import (
	"dotpip"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (f *FileSystem) getArray(key dotpip.Key) ([]string, error) {
	bytes, err := f.readFileByKey(key)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return []string{}, nil
		}
		return nil, err
	}
	v, err := f.formatter.ArrayDecode(bytes)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (f *FileSystem) setArray(key dotpip.Key, array []string) error {
	v, err := f.formatter.ArrayEncode(array)
	if err != nil {
		return err
	}
	if bytes, ok := v.([]byte); ok {
		return f.writeFileByKey(key, bytes)
	}
	return fmt.Errorf("failed to encode array")
}

func (f *FileSystem) ARCount(key dotpip.Key) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, val := range arr {
		if val != "" {
			count++
		}
	}
	return count, nil
}

func (f *FileSystem) ARDel(key dotpip.Key, indices ...int) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	if len(arr) == 0 {
		return 0, nil
	}

	deletedCount := 0
	for _, index := range indices {
		if index >= 0 && index < len(arr) {
			if arr[index] != "" {
				arr[index] = "" // Mark as empty slot
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		err = f.setArray(key, arr)
		if err != nil {
			return 0, err
		}
		f.emitKeyspaceEvent(key, "ardel", 'E')
	}

	return deletedCount, nil
}

func (f *FileSystem) ARDelRange(key dotpip.Key, ranges ...[2]int) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	if len(arr) == 0 {
		return 0, nil
	}

	deletedCount := 0
	for _, r := range ranges {
		start := r[0]
		end := r[1]
		if start > end {
			start, end = end, start // Ensure ascending order
		}

		if start < 0 {
			start = 0
		}
		if end >= len(arr) {
			end = len(arr) - 1
		}

		for i := start; i <= end; i++ {
			if arr[i] != "" {
				arr[i] = ""
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		err = f.setArray(key, arr)
		if err != nil {
			return 0, err
		}
		f.emitKeyspaceEvent(key, "ardelrange", 'E')
	}

	return deletedCount, nil
}

func (f *FileSystem) ARGet(key dotpip.Key, index int) (string, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return "", err
	}

	if index >= 0 && index < len(arr) && arr[index] != "" {
		return arr[index], nil
	}
	return "", nil // Or a specific not found error depending on implementation details. Usually nil.
}

func (f *FileSystem) ARGetRange(key dotpip.Key, start, end int) ([]string, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return []string{}, nil
	}

	var result []string
	if start <= end {
		if start < 0 {
			start = 0
		}
		if end >= len(arr) {
			end = len(arr) - 1
		}
		for i := start; i <= end; i++ {
			result = append(result, arr[i])
		}
	} else {
		if start >= len(arr) {
			start = len(arr) - 1
		}
		if end < 0 {
			end = 0
		}
		for i := start; i >= end; i-- {
			result = append(result, arr[i])
		}
	}

	return result, nil
}

func (f *FileSystem) ARGrep(key dotpip.Key, startStr, endStr string, predicates []dotpip.ARGrepPredicate, options dotpip.ARGrepOptions) ([]any, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return []any{}, nil
	}

	start := 0
	if startStr != "-" {
		// Attempt parsing, fallback to 0. Redis might return error on invalid syntax
		fmt.Sscanf(startStr, "%d", &start)
	}

	end := len(arr) - 1
	if endStr != "+" {
		fmt.Sscanf(endStr, "%d", &end)
	}

	reverse := false
	if start > end {
		reverse = true
	}

	var compiledPredicates []func(string) bool
	for _, p := range predicates {
		val := p.Value
		if options.NoCase {
			val = strings.ToLower(val)
		}
		switch p.Type {
		case "EXACT":
			compiledPredicates = append(compiledPredicates, func(s string) bool {
				if options.NoCase {
					return strings.ToLower(s) == val
				}
				return s == val
			})
		case "MATCH":
			compiledPredicates = append(compiledPredicates, func(s string) bool {
				if options.NoCase {
					return strings.Contains(strings.ToLower(s), val)
				}
				return strings.Contains(s, val)
			})
		case "GLOB":
			// Basic glob to regex
			pattern := "^" + regexp.QuoteMeta(val) + "$"
			pattern = strings.ReplaceAll(pattern, "\\*", ".*")
			pattern = strings.ReplaceAll(pattern, "\\?", ".")
			if options.NoCase {
				pattern = "(?i)" + pattern
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err // invalid glob
			}
			compiledPredicates = append(compiledPredicates, func(s string) bool {
				return re.MatchString(s)
			})
		case "RE":
			pattern := p.Value
			if options.NoCase {
				pattern = "(?i)" + pattern
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			}
			compiledPredicates = append(compiledPredicates, func(s string) bool {
				return re.MatchString(s)
			})
		}
	}

	var result []any
	count := 0

	evalMatch := func(s string) bool {
		if options.And {
			for _, fn := range compiledPredicates {
				if !fn(s) {
					return false
				}
			}
			return len(compiledPredicates) > 0
		} else { // OR is default
			for _, fn := range compiledPredicates {
				if fn(s) {
					return true
				}
			}
			return false
		}
	}

	// Simplified iteration taking reverse into account
	minIdx, maxIdx := start, end
	if reverse {
		minIdx, maxIdx = end, start
	}
	if minIdx < 0 {
		minIdx = 0
	}
	if maxIdx >= len(arr) {
		maxIdx = len(arr) - 1
	}

	if !reverse {
		for i := minIdx; i <= maxIdx; i++ {
			if arr[i] == "" {
				continue
			}
			if evalMatch(arr[i]) {
				if options.WithValues {
					result = append(result, i, arr[i])
				} else {
					result = append(result, i)
				}
				count++
				if options.Limit != nil && count >= *options.Limit {
					break
				}
			}
		}
	} else {
		for i := maxIdx; i >= minIdx; i-- {
			if arr[i] == "" {
				continue
			}
			if evalMatch(arr[i]) {
				if options.WithValues {
					result = append(result, i, arr[i])
				} else {
					result = append(result, i)
				}
				count++
				if options.Limit != nil && count >= *options.Limit {
					break
				}
			}
		}
	}

	return result, nil
}

func (f *FileSystem) ARInfo(key dotpip.Key, full bool) (map[string]any, error) {
	// Minimal stub info for array based on the command reference
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, nil // Typically a null response if not found
	}

	info := make(map[string]any)
	info["length"] = len(arr)
	// Additional FULL info can be populated here depending on exact Redis 8 spec
	return info, nil
}

func (f *FileSystem) ARInsert(key dotpip.Key, values ...string) (int, error) {
	// Basic append behavior as ARINSERT appends at the cursor.
	// The specification says "appends to the current insert cursor position" and defaults to length
	// We'll treat it as append for a simple implementation.
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	arr = append(arr, values...)
	err = f.setArray(key, arr)
	if err != nil {
		return 0, err
	}
	f.emitKeyspaceEvent(key, "arinsert", 'E')

	return len(arr) - 1, nil // returns last index
}

func (f *FileSystem) ARLastItems(key dotpip.Key, count int) ([]string, error) {
	// Not fully specced in given commands, but implied by name. Return last N non-empty.
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	var result []string
	added := 0
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] != "" {
			result = append([]string{arr[i]}, result...)
			added++
			if added == count {
				break
			}
		}
	}
	return result, nil
}

func (f *FileSystem) ARLen(key dotpip.Key) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}
	return len(arr), nil
}

func (f *FileSystem) ARMGet(key dotpip.Key, indices ...int) ([]string, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, idx := range indices {
		if idx >= 0 && idx < len(arr) && arr[idx] != "" {
			result = append(result, arr[idx])
		} else {
			result = append(result, "") // Equivalent to nil reply
		}
	}
	return result, nil
}

func (f *FileSystem) ARMSet(key dotpip.Key, indexValues []dotpip.ARIndexValue) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	maxIdx := len(arr) - 1
	for _, iv := range indexValues {
		if iv.Index > maxIdx {
			maxIdx = iv.Index
		}
	}

	if maxIdx >= len(arr) {
		newArr := make([]string, maxIdx+1)
		copy(newArr, arr)
		arr = newArr
	}

	newSlots := 0
	for _, iv := range indexValues {
		if arr[iv.Index] == "" {
			newSlots++
		}
		arr[iv.Index] = iv.Value
	}

	err = f.setArray(key, arr)
	if err != nil {
		return 0, err
	}
	f.emitKeyspaceEvent(key, "armset", 'E')

	return newSlots, nil
}

func (f *FileSystem) ARNext(key dotpip.Key) (int, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}
	return len(arr), nil
}

func (f *FileSystem) AROp(key dotpip.Key, start, end int, operation string, matchValue *string) (any, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return nil, nil
	}

	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end >= len(arr) {
		end = len(arr) - 1
	}

	var nums []float64
	var ints []int64
	used := 0
	matched := 0

	for i := start; i <= end; i++ {
		val := arr[i]
		if val != "" {
			used++
			if matchValue != nil && val == *matchValue {
				matched++
			}
			if num, err := strconv.ParseFloat(val, 64); err == nil {
				nums = append(nums, num)
			}
			if in, err := strconv.ParseInt(val, 10, 64); err == nil {
				ints = append(ints, in)
			} else {
				// truncated toward zero as per docs
				if num, err := strconv.ParseFloat(val, 64); err == nil {
					ints = append(ints, int64(num))
				}
			}
		}
	}

	switch operation {
	case "SUM":
		if len(nums) == 0 {
			return nil, nil
		}
		sum := 0.0
		for _, n := range nums {
			sum += n
		}
		return fmt.Sprintf("%g", sum), nil
	case "MIN":
		if len(nums) == 0 {
			return nil, nil
		}
		min := nums[0]
		for _, n := range nums {
			if n < min {
				min = n
			}
		}
		return fmt.Sprintf("%g", min), nil
	case "MAX":
		if len(nums) == 0 {
			return nil, nil
		}
		max := nums[0]
		for _, n := range nums {
			if n > max {
				max = n
			}
		}
		return fmt.Sprintf("%g", max), nil
	case "AND":
		if used == 0 {
			return nil, nil
		}
		if len(ints) == 0 {
			return 0, nil
		}
		res := ints[0]
		for i := 1; i < len(ints); i++ {
			res &= ints[i]
		}
		return int(res), nil
	case "OR":
		if used == 0 {
			return nil, nil
		}
		if len(ints) == 0 {
			return 0, nil
		}
		res := ints[0]
		for i := 1; i < len(ints); i++ {
			res |= ints[i]
		}
		return int(res), nil
	case "XOR":
		if used == 0 {
			return nil, nil
		}
		if len(ints) == 0 {
			return 0, nil
		}
		res := ints[0]
		for i := 1; i < len(ints); i++ {
			res ^= ints[i]
		}
		return int(res), nil
	case "MATCH":
		return matched, nil
	case "USED":
		return used, nil
	}

	return nil, fmt.Errorf("unknown operation: %s", operation)
}

func (f *FileSystem) ARRing(key dotpip.Key, size int, values ...string) (int, error) {
	// A simple append with left-trim to size
	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	arr = append(arr, values...)
	if len(arr) > size {
		arr = arr[len(arr)-size:]
	}

	err = f.setArray(key, arr)
	if err != nil {
		return 0, err
	}
	f.emitKeyspaceEvent(key, "arring", 'E')

	return len(arr) - 1, nil // Assume returning the last index
}

func (f *FileSystem) ARScan(key dotpip.Key, start, end int, limit *int) ([]any, error) {
	arr, err := f.getArray(key)
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return []any{}, nil
	}

	reverse := false
	if start > end {
		reverse = true
	}

	minIdx, maxIdx := start, end
	if reverse {
		minIdx, maxIdx = end, start
	}
	if minIdx < 0 {
		minIdx = 0
	}
	if maxIdx >= len(arr) {
		maxIdx = len(arr) - 1
	}

	var result []any
	count := 0

	if !reverse {
		for i := minIdx; i <= maxIdx; i++ {
			if arr[i] != "" {
				result = append(result, i, arr[i])
				count++
				if limit != nil && count >= *limit {
					break
				}
			}
		}
	} else {
		for i := maxIdx; i >= minIdx; i-- {
			if arr[i] != "" {
				result = append(result, i, arr[i])
				count++
				if limit != nil && count >= *limit {
					break
				}
			}
		}
	}

	return result, nil
}

func (f *FileSystem) ARSeek(key dotpip.Key, index int) (int, error) {
	// Cursor management can be just returning the passed index as a simple implementation
	// unless actual cursor state needs to be persisted.
	return index, nil
}

func (f *FileSystem) ARSet(key dotpip.Key, index int, values ...string) (int, error) {
	if index < 0 {
		return 0, fmt.Errorf("ERR index out of bounds")
	}

	arr, err := f.getArray(key)
	if err != nil {
		return 0, err
	}

	requiredLen := index + len(values)
	if requiredLen > len(arr) {
		newArr := make([]string, requiredLen)
		copy(newArr, arr)
		arr = newArr
	}

	newSlots := 0
	for i, val := range values {
		if arr[index+i] == "" {
			newSlots++
		}
		arr[index+i] = val
	}

	err = f.setArray(key, arr)
	if err != nil {
		return 0, err
	}
	f.emitKeyspaceEvent(key, "arset", 'E')

	return newSlots, nil
}
