package fs

import "dotpip"

// YAMLArrAppend appends values to a YAML array.
func (f *FileSystem) YAMLArrAppend(key dotpip.Key, path string, values ...any) ([]any, error) {
	return f.JSONArrAppend(key, path, values...)
}

// YAMLArrIndex returns the index of a value in a YAML array.
func (f *FileSystem) YAMLArrIndex(key dotpip.Key, path string, value any, startAndStop ...int) ([]any, error) {
	return f.JSONArrIndex(key, path, value, startAndStop...)
}

// YAMLArrInsert inserts values into a YAML array at the specified index.
func (f *FileSystem) YAMLArrInsert(key dotpip.Key, path string, index int, values ...any) ([]any, error) {
	return f.JSONArrInsert(key, path, index, values...)
}

// YAMLArrLen returns the length of a YAML array.
func (f *FileSystem) YAMLArrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONArrLen(key, path)
}

// YAMLArrPop removes and returns an element from the index in the array.
func (f *FileSystem) YAMLArrPop(key dotpip.Key, path string, index ...int) ([]any, error) {
	return f.JSONArrPop(key, path, index...)
}

// YAMLArrTrim trims an array so that it contains only the specified inclusive range of elements.
func (f *FileSystem) YAMLArrTrim(key dotpip.Key, path string, start int, stop int) ([]any, error) {
	return f.JSONArrTrim(key, path, start, stop)
}

// YAMLClear clear container values.
func (f *FileSystem) YAMLClear(key dotpip.Key, path string) (int, error) {
	return f.JSONClear(key, path)
}

// YAMLDebug reports memory usage.
func (f *FileSystem) YAMLDebug(subcommand string, key dotpip.Key, path string) (any, error) {
	return f.JSONDebug(subcommand, key, path)
}

// YAMLDel deletes a value.
func (f *FileSystem) YAMLDel(key dotpip.Key, path string) (int, error) {
	return f.JSONDel(key, path)
}

// YAMLForget is an alias for YAMLDel.
func (f *FileSystem) YAMLForget(key dotpip.Key, path string) (int, error) {
	return f.JSONForget(key, path)
}

// YAMLGet returns the value at path.
func (f *FileSystem) YAMLGet(key dotpip.Key, paths ...string) (any, error) {
	return f.JSONGet(key, paths...)
}

// YAMLMerge merges a given YAML value into matching paths.
func (f *FileSystem) YAMLMerge(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONMerge(key, path, value)
}

// YAMLMGet returns the values at path from multiple key arguments.
func (f *FileSystem) YAMLMGet(path string, keys ...dotpip.Key) ([]any, error) {
	return f.JSONMGet(path, keys...)
}

// YAMLMSet sets or updates one or more YAML values.
func (f *FileSystem) YAMLMSet(args ...dotpip.JSONMSetArg) (string, error) {
	return f.JSONMSet(args...)
}

// YAMLNumIncrBy increments the number value stored at path by number.
func (f *FileSystem) YAMLNumIncrBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumIncrBy(key, path, value)
}

// YAMLNumMultBy multiplies the number value stored at path by number.
func (f *FileSystem) YAMLNumMultBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumMultBy(key, path, value)
}

// YAMLObjKeys returns the keys in the object that s referenced by path.
func (f *FileSystem) YAMLObjKeys(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjKeys(key, path)
}

// YAMLObjLen report the number of keys in the YAML object at path in key.
func (f *FileSystem) YAMLObjLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjLen(key, path)
}

// YAMLResp returns the YAML in key in Redis Serialization Protocol (RESP).
func (f *FileSystem) YAMLResp(key dotpip.Key, path string) (any, error) {
	return f.JSONResp(key, path)
}

// YAMLSet sets the YAML value at path in key.
func (f *FileSystem) YAMLSet(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONSet(key, path, value)
}

// YAMLStrAppend appends the yaml-string values to the string at path.
func (f *FileSystem) YAMLStrAppend(key dotpip.Key, path string, value string) ([]any, error) {
	return f.JSONStrAppend(key, path, value)
}

// YAMLStrLen report the length of the YAML String at path in key.
func (f *FileSystem) YAMLStrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONStrLen(key, path)
}

// YAMLToggle toggles a boolean value stored at path.
func (f *FileSystem) YAMLToggle(key dotpip.Key, path string) ([]any, error) {
	return f.JSONToggle(key, path)
}

// YAMLType reports the type of YAML value at path.
func (f *FileSystem) YAMLType(key dotpip.Key, path string) ([]any, error) {
	return f.JSONType(key, path)
}
