package fs

import "dotpip"

func (f *FileSystem) YAMLArrAppend(key dotpip.Key, path string, values ...any) ([]any, error) {
	return f.JSONArrAppend(key, path, values...)
}

func (f *FileSystem) YAMLArrIndex(key dotpip.Key, path string, value any, startAndStop ...int) ([]any, error) {
	return f.JSONArrIndex(key, path, value, startAndStop...)
}

func (f *FileSystem) YAMLArrInsert(key dotpip.Key, path string, index int, values ...any) ([]any, error) {
	return f.JSONArrInsert(key, path, index, values...)
}

func (f *FileSystem) YAMLArrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONArrLen(key, path)
}

func (f *FileSystem) YAMLArrPop(key dotpip.Key, path string, index ...int) ([]any, error) {
	return f.JSONArrPop(key, path, index...)
}

func (f *FileSystem) YAMLArrTrim(key dotpip.Key, path string, start int, stop int) ([]any, error) {
	return f.JSONArrTrim(key, path, start, stop)
}

func (f *FileSystem) YAMLClear(key dotpip.Key, path string) (int, error) {
	return f.JSONClear(key, path)
}

func (f *FileSystem) YAMLDebug(subcommand string, key dotpip.Key, path string) (any, error) {
	return f.JSONDebug(subcommand, key, path)
}

func (f *FileSystem) YAMLDel(key dotpip.Key, path string) (int, error) {
	return f.JSONDel(key, path)
}

func (f *FileSystem) YAMLForget(key dotpip.Key, path string) (int, error) {
	return f.JSONForget(key, path)
}

func (f *FileSystem) YAMLGet(key dotpip.Key, paths ...string) (any, error) {
	return f.JSONGet(key, paths...)
}

func (f *FileSystem) YAMLMerge(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONMerge(key, path, value)
}

func (f *FileSystem) YAMLMGet(path string, keys ...dotpip.Key) ([]any, error) {
	return f.JSONMGet(path, keys...)
}

func (f *FileSystem) YAMLMSet(args ...dotpip.JSONMSetArg) (string, error) {
	return f.JSONMSet(args...)
}

func (f *FileSystem) YAMLNumIncrBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumIncrBy(key, path, value)
}

func (f *FileSystem) YAMLNumMultBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumMultBy(key, path, value)
}

func (f *FileSystem) YAMLObjKeys(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjKeys(key, path)
}

func (f *FileSystem) YAMLObjLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjLen(key, path)
}

func (f *FileSystem) YAMLResp(key dotpip.Key, path string) (any, error) {
	return f.JSONResp(key, path)
}

func (f *FileSystem) YAMLSet(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONSet(key, path, value)
}

func (f *FileSystem) YAMLStrAppend(key dotpip.Key, path string, value string) ([]any, error) {
	return f.JSONStrAppend(key, path, value)
}

func (f *FileSystem) YAMLStrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONStrLen(key, path)
}

func (f *FileSystem) YAMLToggle(key dotpip.Key, path string) ([]any, error) {
	return f.JSONToggle(key, path)
}

func (f *FileSystem) YAMLType(key dotpip.Key, path string) ([]any, error) {
	return f.JSONType(key, path)
}
