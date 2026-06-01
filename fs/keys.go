package fs

import (
	"dotpip"
	"os"
)

func (f *fileSystem) Copy(source dotpip.Key, destination dotpip.Key, options ...dotpip.CopyOption) int {
	cmd := &dotpip.CopyCommand{}
	for _, option := range options {
		option(cmd)
	}

	sourceExist, err := f.checkExistByKey(source)
	if err != nil {
		return 0
	}

	if !sourceExist {
		return 0
	}

	if cmd.Destination != nil {
		content, err := f.Get(source)
		if err != nil {
			return 0
		}

		_, err = cmd.Destination.Set(destination, content)
		if err != nil {
			return 0
		}

		return 1
	}

	if cmd.Replace {
		err = f.removeFileByKey(destination)
		if err != nil {
			return 0
		}
	}

	content, err := f.readFileByKey(source)
	if err != nil {
		return 0
	}

	err = f.writeFileByKey(destination, content)
	if err != nil {
		return 0
	}

	return 1
}

func (f *fileSystem) Del(keys ...dotpip.Key) int {
	count := 0
	for _, key := range keys {
		err := f.removeFileByKey(key)
		if err != nil {
			return count
		}
		count++
	}

	return count
}

func (f *fileSystem) Exists(keys ...dotpip.Key) ([]bool, error) {
	results := make([]bool, len(keys))
	for i, key := range keys {
		exist, err := f.checkExistByKey(key)
		results[i] = exist && err == nil
	}

	return results, nil
}

func (f *fileSystem) FlushAll() (err error) {
	err = os.RemoveAll(f.pathRoot)
	if err != nil {
		return err
	}

	return os.MkdirAll(f.pathRoot, 0755)
}
