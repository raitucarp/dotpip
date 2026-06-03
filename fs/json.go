package fs

import (
	"encoding/json"
	"errors"
	"os"

	"dotpip"
	"github.com/ohler55/ojg/jp"
)

func jsonCompare(a, b any) bool {
	ba, err1 := json.Marshal(a)
	bb, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(ba) == string(bb)
}

func (f *FileSystem) readJSON(key dotpip.Key) (any, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return f.JSONDecode(content)
}

func (f *FileSystem) writeJSON(key dotpip.Key, value any) error {
	if value == nil {
		err := f.removeFileByKey(key)
		if err == nil {
			f.emitKeyspaceEvent(key, "del", 'g')
		}
		return err
	}

	encoded, err := f.JSONEncode(value)
	if err != nil {
		return err
	}

	err = f.writeFileByKey(key, encoded.([]byte))
	if err == nil {
		f.emitKeyspaceEvent(key, "set", 'g')
	}
	return err
}

func (f *FileSystem) JSONSet(key dotpip.Key, path string, value any) (string, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return "", err
	}

	doc, err := f.readJSON(key)
	if err != nil {
		return "", err
	}

	if doc == nil {
		if path == "$" || path == "." {
			doc = value
		} else {
			return "", errors.New("ERR new objects must be created at the root")
		}
	} else {
		err = expr.Set(doc, value)
		if err != nil {
			return "", err
		}
	}

	err = f.writeJSON(key, doc)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func (f *FileSystem) JSONGet(key dotpip.Key, paths ...string) (any, error) {
	doc, err := f.readJSON(key)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, nil
	}

	if len(paths) == 0 {
		return doc, nil
	}

	if len(paths) == 1 {
		expr, err := jp.ParseString(paths[0])
		if err != nil {
			return nil, err
		}
		res := expr.Get(doc)
		if len(res) == 0 {
			return nil, nil
		}
		if len(res) == 1 {
			return res[0], nil
		}
		return res, nil
	}

	result := make(map[string]any)
	for _, path := range paths {
		expr, err := jp.ParseString(path)
		if err != nil {
			return nil, err
		}
		res := expr.Get(doc)
		switch len(res) {
		case 0:
			result[path] = nil
		case 1:
			result[path] = res[0]
		default:
			result[path] = res
		}
	}

	return result, nil
}

func (f *FileSystem) JSONDel(key dotpip.Key, path string) (int, error) {
	if path == "" {
		path = "$"
	}
	expr, err := jp.ParseString(path)
	if err != nil {
		return 0, err
	}

	if path == "$" {
		doc, err := f.readJSON(key)
		if err != nil || doc == nil {
			return 0, err
		}
		err = f.writeJSON(key, nil)
		return 1, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return 0, err
	}

	matches := expr.Get(doc)
	deletedCount := len(matches)

	err = expr.Del(doc)
	if err != nil {
		return 0, nil
	}

	if deletedCount > 0 {
		err = f.writeJSON(key, doc)
	}
	return deletedCount, err
}

func (f *FileSystem) JSONArrAppend(key dotpip.Key, path string, values ...any) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if arr, ok := v.([]any); ok {
			arr = append(arr, values...)
			results = append(results, len(arr))
			changedAny = true
			return arr, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONArrIndex(key dotpip.Key, path string, value any, startAndStop ...int) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	start := 0
	stop := 0
	hasStart := len(startAndStop) > 0
	hasStop := len(startAndStop) > 1

	if hasStart {
		start = startAndStop[0]
	}
	if hasStop {
		stop = startAndStop[1]
	}

	var results []any
	res := expr.Get(doc)
	for _, match := range res {
		if arr, ok := match.([]any); ok {
			arrLen := len(arr)

			s := start
			if s < 0 {
				s = arrLen + s
				if s < 0 {
					s = 0
				}
			}

			e := stop
			if hasStop {
				if e < 0 {
					e = arrLen + e
					if e < 0 {
						e = 0
					}
				}
				if e > arrLen {
					e = arrLen
				}
			} else {
				e = arrLen
			}

			foundIdx := -1
			if s <= e {
				for i := s; i < e; i++ {
					if jsonCompare(arr[i], value) {
						foundIdx = i
						break
					}
				}
			}
			results = append(results, foundIdx)
		} else {
			results = append(results, nil)
		}
	}

	return results, nil
}

func (f *FileSystem) JSONArrInsert(key dotpip.Key, path string, index int, values ...any) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if arr, ok := v.([]any); ok {
			arrLen := len(arr)
			idx := index
			if idx < 0 {
				idx = arrLen + idx
			}
			if idx < 0 {
				idx = 0
			}
			if idx > arrLen {
				idx = arrLen
			}

			newArr := make([]any, 0, arrLen+len(values))
			newArr = append(newArr, arr[:idx]...)
			newArr = append(newArr, values...)
			newArr = append(newArr, arr[idx:]...)

			results = append(results, len(newArr))
			changedAny = true
			return newArr, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONArrLen(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	res := expr.Get(doc)
	for _, match := range res {
		if arr, ok := match.([]any); ok {
			results = append(results, len(arr))
		} else {
			results = append(results, nil)
		}
	}
	return results, nil
}

func (f *FileSystem) JSONArrPop(key dotpip.Key, path string, index ...int) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	idxToPop := -1
	if len(index) > 0 {
		idxToPop = index[0]
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if arr, ok := v.([]any); ok {
			arrLen := len(arr)
			if arrLen == 0 {
				results = append(results, nil)
				return v, false
			}

			idx := idxToPop
			if idx < 0 {
				idx = arrLen + idx
			}
			if idx < 0 {
				idx = 0
			}
			if idx >= arrLen {
				idx = arrLen - 1
			}

			poppedVal := arr[idx]
			newArr := make([]any, 0, arrLen-1)
			newArr = append(newArr, arr[:idx]...)
			newArr = append(newArr, arr[idx+1:]...)

			results = append(results, poppedVal)
			changedAny = true
			return newArr, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONArrTrim(key dotpip.Key, path string, start int, stop int) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if arr, ok := v.([]any); ok {
			arrLen := len(arr)
			s := start
			if s < 0 {
				s = arrLen + s
				if s < 0 {
					s = 0
				}
			}
			e := stop
			if e < 0 {
				e = arrLen + e
				if e < 0 {
					e = 0
				}
			}

			if s > arrLen {
				s = arrLen
			}
			if e >= arrLen {
				e = arrLen - 1
			}

			if s > e {
				newArr := make([]any, 0)
				results = append(results, 0)
				changedAny = true
				return newArr, true
			}

			newArr := make([]any, e-s+1)
			copy(newArr, arr[s:e+1])

			results = append(results, len(newArr))
			changedAny = true
			return newArr, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONClear(key dotpip.Key, path string) (int, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return 0, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return 0, err
	}

	clearedCount := 0
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		switch val := v.(type) {
		case []any:
			if len(val) > 0 {
				clearedCount++
				changedAny = true
				return make([]any, 0), true
			}
			return v, false
		case map[string]any:
			if len(val) > 0 {
				clearedCount++
				changedAny = true
				return make(map[string]any), true
			}
			return v, false
		case float64, int, int64, uint64:
			if val != float64(0) {
				clearedCount++
				changedAny = true
				return float64(0), true
			}
			return v, false
		}
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return clearedCount, err
}

func (f *FileSystem) JSONForget(key dotpip.Key, path string) (int, error) {
	return f.JSONDel(key, path)
}

func (f *FileSystem) JSONMerge(key dotpip.Key, path string, value any) (string, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return "", err
	}

	doc, err := f.readJSON(key)
	if err != nil {
		return "", err
	}
	if doc == nil {
		doc = make(map[string]any)
	}

	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		// Merge logic:
		// If both are objects, merge keys (value keys overwrite or add, null values delete keys)
		// Otherwise, replace the whole value.

		if mapVal, isMap := v.(map[string]any); isMap {
			if mergeVal, isMergeMap := value.(map[string]any); isMergeMap {
				for k, v2 := range mergeVal {
					if v2 == nil {
						delete(mapVal, k)
					} else {
						mapVal[k] = v2
					}
				}
				changedAny = true
				return mapVal, true
			}
		}
		changedAny = true
		return value, true
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return "OK", err
}

func (f *FileSystem) JSONMGet(path string, keys ...dotpip.Key) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	results := make([]any, len(keys))
	for i, key := range keys {
		doc, err := f.readJSON(key)
		if err != nil || doc == nil {
			results[i] = nil
			continue
		}

		res := expr.Get(doc)
		switch len(res) {
		case 0:
			results[i] = nil
		case 1:
			results[i] = res[0]
		default:
			results[i] = res
		}
	}
	return results, nil
}

func (f *FileSystem) JSONMSet(args ...dotpip.JSONMSetArg) (string, error) {
	for _, arg := range args {
		_, err := f.JSONSet(arg.Key, arg.Path, arg.Value)
		if err != nil {
			return "", err
		}
	}
	return "OK", nil
}

func (f *FileSystem) JSONNumIncrBy(key dotpip.Key, path string, value float64) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		switch num := v.(type) {
		case float64:
			newVal := num + value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case int:
			newVal := float64(num) + value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case int64:
			newVal := float64(num) + value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case uint64:
			newVal := float64(num) + value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONNumMultBy(key dotpip.Key, path string, value float64) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		switch num := v.(type) {
		case float64:
			newVal := num * value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case int:
			newVal := float64(num) * value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case int64:
			newVal := float64(num) * value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		case uint64:
			newVal := float64(num) * value
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONObjKeys(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	res := expr.Get(doc)
	var results []any
	for _, match := range res {
		if mapVal, ok := match.(map[string]any); ok {
			var keys []string
			for k := range mapVal {
				keys = append(keys, k)
			}
			results = append(results, keys)
		} else {
			results = append(results, nil)
		}
	}
	return results, nil
}

func (f *FileSystem) JSONObjLen(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	res := expr.Get(doc)
	var results []any
	for _, match := range res {
		if mapVal, ok := match.(map[string]any); ok {
			results = append(results, len(mapVal))
		} else {
			results = append(results, nil)
		}
	}
	return results, nil
}

func (f *FileSystem) JSONStrAppend(key dotpip.Key, path string, value string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if str, ok := v.(string); ok {
			newStr := str + value
			results = append(results, len(newStr))
			changedAny = true
			return newStr, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONStrLen(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	res := expr.Get(doc)
	var results []any
	for _, match := range res {
		if str, ok := match.(string); ok {
			results = append(results, len(str))
		} else {
			results = append(results, nil)
		}
	}
	return results, nil
}

func (f *FileSystem) JSONToggle(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	var results []any
	changedAny := false
	doc, err = expr.Modify(doc, func(v any) (any, bool) {
		if b, ok := v.(bool); ok {
			newVal := !b
			results = append(results, newVal)
			changedAny = true
			return newVal, true
		}
		results = append(results, nil)
		return v, false
	})

	if changedAny {
		err = f.writeJSON(key, doc)
	}
	return results, err
}

func (f *FileSystem) JSONType(key dotpip.Key, path string) ([]any, error) {
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	res := expr.Get(doc)
	var results []any
	for _, match := range res {
		switch match.(type) {
		case string:
			results = append(results, "string")
		case float64, int, int64, uint64:
			results = append(results, "number")
		case bool:
			results = append(results, "boolean")
		case []any:
			results = append(results, "array")
		case map[string]any:
			results = append(results, "object")
		case nil:
			results = append(results, "null")
		default:
			results = append(results, "unknown")
		}
	}
	return results, nil
}

func (f *FileSystem) JSONResp(key dotpip.Key, path string) (any, error) {
	if path == "" {
		path = "$"
	}
	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	doc, err := f.readJSON(key)
	if err != nil || doc == nil {
		return nil, err
	}

	res := expr.Get(doc)

	// This command converts the JSON representation to RESP
	// In the file system implementation, returning the value itself is enough since RESP mapping is done at networking level.
	// We'll return an array of values for each match.
	return res, nil
}

func (f *FileSystem) JSONDebug(subcommand string, key dotpip.Key, path string) (any, error) {
	// For MEMORY subcommand, return an estimated size
	if subcommand == "MEMORY" {
		expr, err := jp.ParseString(path)
		if err != nil {
			return nil, err
		}

		doc, err := f.readJSON(key)
		if err != nil || doc == nil {
			return nil, err
		}

		res := expr.Get(doc)
		var results []any
		for _, match := range res {
			b, err := json.Marshal(match)
			if err == nil {
				results = append(results, len(b))
			} else {
				results = append(results, 0)
			}
		}
		return results, nil
	}

	// Other subcommands can just return basic OK for now
	return "OK", nil
}
