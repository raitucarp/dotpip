package fs

import "dotpip"

// TOMLArrAppend appends values to a TOML array.
func (f *FileSystem) TOMLArrAppend(key dotpip.Key, path string, values ...any) ([]any, error) {
	return f.JSONArrAppend(key, path, values...)
}

// TOMLArrIndex returns the index of a value in a TOML array.
func (f *FileSystem) TOMLArrIndex(key dotpip.Key, path string, value any, startAndStop ...int) ([]any, error) {
	return f.JSONArrIndex(key, path, value, startAndStop...)
}

// TOMLArrInsert inserts values into a TOML array at the specified index.
func (f *FileSystem) TOMLArrInsert(key dotpip.Key, path string, index int, values ...any) ([]any, error) {
	return f.JSONArrInsert(key, path, index, values...)
}

// TOMLArrLen returns the length of a TOML array.
func (f *FileSystem) TOMLArrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONArrLen(key, path)
}

// TOMLArrPop removes and returns an element from the index in the array.
func (f *FileSystem) TOMLArrPop(key dotpip.Key, path string, index ...int) ([]any, error) {
	return f.JSONArrPop(key, path, index...)
}

// TOMLArrTrim trims an array so that it contains only the specified inclusive range of elements.
func (f *FileSystem) TOMLArrTrim(key dotpip.Key, path string, start int, stop int) ([]any, error) {
	return f.JSONArrTrim(key, path, start, stop)
}

// TOMLClear clear container values.
func (f *FileSystem) TOMLClear(key dotpip.Key, path string) (int, error) {
	return f.JSONClear(key, path)
}

// TOMLDebug reports memory usage.
func (f *FileSystem) TOMLDebug(subcommand string, key dotpip.Key, path string) (any, error) {
	return f.JSONDebug(subcommand, key, path)
}

// TOMLDel deletes a value.
func (f *FileSystem) TOMLDel(key dotpip.Key, path string) (int, error) {
	return f.JSONDel(key, path)
}

// TOMLForget is an alias for TOMLDel.
func (f *FileSystem) TOMLForget(key dotpip.Key, path string) (int, error) {
	return f.JSONForget(key, path)
}

// TOMLGet returns the value at path.
func (f *FileSystem) TOMLGet(key dotpip.Key, paths ...string) (any, error) {
	return f.JSONGet(key, paths...)
}

// TOMLMerge merges a given TOML value into matching paths.
func (f *FileSystem) TOMLMerge(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONMerge(key, path, value)
}

// TOMLMGet returns the values at path from multiple key arguments.
func (f *FileSystem) TOMLMGet(path string, keys ...dotpip.Key) ([]any, error) {
	return f.JSONMGet(path, keys...)
}

// TOMLMSet sets or updates one or more TOML values.
func (f *FileSystem) TOMLMSet(args ...dotpip.JSONMSetArg) (string, error) {
	return f.JSONMSet(args...)
}

// TOMLNumIncrBy increments the number value stored at path by number.
func (f *FileSystem) TOMLNumIncrBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumIncrBy(key, path, value)
}

// TOMLNumMultBy multiplies the number value stored at path by number.
func (f *FileSystem) TOMLNumMultBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumMultBy(key, path, value)
}

// TOMLObjKeys returns the keys in the object that s referenced by path.
func (f *FileSystem) TOMLObjKeys(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjKeys(key, path)
}

// TOMLObjLen returns TOML object length.
func (f *FileSystem) TOMLObjLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjLen(key, path)
}

// TOMLResp returns RESP formatted TOML.
func (f *FileSystem) TOMLResp(key dotpip.Key, path string) (any, error) {
	return f.JSONResp(key, path)
}

// TOMLSet sets a TOML value.
func (f *FileSystem) TOMLSet(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONSet(key, path, value)
}

// TOMLStrAppend appends to a TOML string.
func (f *FileSystem) TOMLStrAppend(key dotpip.Key, path string, value string) ([]any, error) {
	return f.JSONStrAppend(key, path, value)
}

// TOMLStrLen returns TOML string length.
func (f *FileSystem) TOMLStrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONStrLen(key, path)
}

// TOMLToggle toggles a TOML boolean.
func (f *FileSystem) TOMLToggle(key dotpip.Key, path string) ([]any, error) {
	return f.JSONToggle(key, path)
}

// TOMLType returns the type of a TOML value.
func (f *FileSystem) TOMLType(key dotpip.Key, path string) ([]any, error) {
	return f.JSONType(key, path)
}
