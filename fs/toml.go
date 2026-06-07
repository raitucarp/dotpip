package fs

import "dotpip"

func (f *FileSystem) TOMLArrAppend(key dotpip.Key, path string, values ...any) ([]any, error) {
	return f.JSONArrAppend(key, path, values...)
}

func (f *FileSystem) TOMLArrIndex(key dotpip.Key, path string, value any, startAndStop ...int) ([]any, error) {
	return f.JSONArrIndex(key, path, value, startAndStop...)
}

func (f *FileSystem) TOMLArrInsert(key dotpip.Key, path string, index int, values ...any) ([]any, error) {
	return f.JSONArrInsert(key, path, index, values...)
}

func (f *FileSystem) TOMLArrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONArrLen(key, path)
}

func (f *FileSystem) TOMLArrPop(key dotpip.Key, path string, index ...int) ([]any, error) {
	return f.JSONArrPop(key, path, index...)
}

func (f *FileSystem) TOMLArrTrim(key dotpip.Key, path string, start int, stop int) ([]any, error) {
	return f.JSONArrTrim(key, path, start, stop)
}

func (f *FileSystem) TOMLClear(key dotpip.Key, path string) (int, error) {
	return f.JSONClear(key, path)
}

func (f *FileSystem) TOMLDebug(subcommand string, key dotpip.Key, path string) (any, error) {
	return f.JSONDebug(subcommand, key, path)
}

func (f *FileSystem) TOMLDel(key dotpip.Key, path string) (int, error) {
	return f.JSONDel(key, path)
}

func (f *FileSystem) TOMLForget(key dotpip.Key, path string) (int, error) {
	return f.JSONForget(key, path)
}

func (f *FileSystem) TOMLGet(key dotpip.Key, paths ...string) (any, error) {
	return f.JSONGet(key, paths...)
}

func (f *FileSystem) TOMLMerge(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONMerge(key, path, value)
}

func (f *FileSystem) TOMLMGet(path string, keys ...dotpip.Key) ([]any, error) {
	return f.JSONMGet(path, keys...)
}

func (f *FileSystem) TOMLMSet(args ...dotpip.JSONMSetArg) (string, error) {
	return f.JSONMSet(args...)
}

func (f *FileSystem) TOMLNumIncrBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumIncrBy(key, path, value)
}

func (f *FileSystem) TOMLNumMultBy(key dotpip.Key, path string, value float64) ([]any, error) {
	return f.JSONNumMultBy(key, path, value)
}

func (f *FileSystem) TOMLObjKeys(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjKeys(key, path)
}

func (f *FileSystem) TOMLObjLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONObjLen(key, path)
}

func (f *FileSystem) TOMLResp(key dotpip.Key, path string) (any, error) {
	return f.JSONResp(key, path)
}

func (f *FileSystem) TOMLSet(key dotpip.Key, path string, value any) (string, error) {
	return f.JSONSet(key, path, value)
}

func (f *FileSystem) TOMLStrAppend(key dotpip.Key, path string, value string) ([]any, error) {
	return f.JSONStrAppend(key, path, value)
}

func (f *FileSystem) TOMLStrLen(key dotpip.Key, path string) ([]any, error) {
	return f.JSONStrLen(key, path)
}

func (f *FileSystem) TOMLToggle(key dotpip.Key, path string) ([]any, error) {
	return f.JSONToggle(key, path)
}

func (f *FileSystem) TOMLType(key dotpip.Key, path string) ([]any, error) {
	return f.JSONType(key, path)
}
