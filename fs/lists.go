package fs

import (
	"dotpip"
	"errors"
	"os"
)

func (f *fileSystem) readList(key dotpip.Key) ([]string, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	anyList, err := f.formatter.ListDecode(content)
	if err != nil {
		return nil, err
	}

	strList := make([]string, len(anyList))
	for i, v := range anyList {
		switch val := v.(type) {
		case string:
			strList[i] = val
		default:
			// If not string, might be float64 from JSON unmarshal if numeric
			// but Redis lists should contain strings.
			// Coerce using fmt or similar if needed. For now, assume string.
			if valStr, ok := v.(string); ok {
				strList[i] = valStr
			}
		}
	}
	return strList, nil
}

func (f *fileSystem) writeList(key dotpip.Key, list []string) error {
	if len(list) == 0 {
		return f.removeFileByKey(key)
	}

	anyList := make([]any, len(list))
	for i, v := range list {
		anyList[i] = v
	}

	encoded, err := f.formatter.ListEncode(anyList)
	if err != nil {
		return err
	}

	return f.writeFileByKey(key, encoded.([]byte))
}

func (f *fileSystem) LIndex(key dotpip.Key, index int) (string, error) {
	list, err := f.readList(key)
	if err != nil {
		return "", err
	}

	length := len(list)
	if length == 0 {
		return "", nil // or error depending on if we return error on empty
	}

	if index < 0 {
		index = length + index
	}

	if index < 0 || index >= length {
		return "", nil // Redis returns nil
	}

	return list[index], nil
}

func (f *fileSystem) LInsert(key dotpip.Key, option dotpip.LInsertOption, pivot string, element string) (int, error) {
	list, err := f.readList(key)
	if err != nil {
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil // Key does not exist
	}

	pivotIndex := -1
	for i, v := range list {
		if v == pivot {
			pivotIndex = i
			break
		}
	}

	if pivotIndex == -1 {
		return -1, nil // Pivot not found
	}

	insertIndex := pivotIndex
	if option == dotpip.After {
		insertIndex++
	}

	list = append(list[:insertIndex], append([]string{element}, list[insertIndex:]...)...)

	err = f.writeList(key, list)
	if err != nil {
		return 0, err
	}

	return len(list), nil
}

func (f *fileSystem) LLen(key dotpip.Key) (int, error) {
	list, err := f.readList(key)
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (f *fileSystem) LMove(source dotpip.Key, destination dotpip.Key, srcDir dotpip.LMoveDir, destDir dotpip.LMoveDir) (string, error) {
	srcList, err := f.readList(source)
	if err != nil {
		return "", err
	}
	if len(srcList) == 0 {
		return "", nil
	}

	var val string
	if srcDir == dotpip.Left {
		val = srcList[0]
		srcList = srcList[1:]
	} else {
		val = srcList[len(srcList)-1]
		srcList = srcList[:len(srcList)-1]
	}

	// Handle case where source and destination are the same
	sameKey := false
	if len(source) == len(destination) {
		sameKey = true
		for i := range source {
			if source[i] != destination[i] {
				sameKey = false
				break
			}
		}
	}

	var destList []string
	if sameKey {
		destList = srcList
	} else {
		destList, err = f.readList(destination)
		if err != nil {
			return "", err
		}
	}

	if destDir == dotpip.Left {
		destList = append([]string{val}, destList...)
	} else {
		destList = append(destList, val)
	}

	if sameKey {
		err = f.writeList(source, destList)
		if err != nil {
			return "", err
		}
	} else {
		err = f.writeList(source, srcList)
		if err != nil {
			return "", err
		}

		err = f.writeList(destination, destList)
		if err != nil {
			return "", err
		}
	}

	return val, nil
}

func (f *fileSystem) LPop(key dotpip.Key, count int) ([]string, error) {
	list, err := f.readList(key)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	if count <= 0 {
		count = 1
	}

	if count > len(list) {
		count = len(list)
	}

	popped := list[:count]
	list = list[count:]

	err = f.writeList(key, list)
	if err != nil {
		return nil, err
	}

	return popped, nil
}

func (f *fileSystem) LPos(key dotpip.Key, element string, options ...dotpip.LPosOption) ([]int, error) {
	cmd := &dotpip.LPosCommand{
		Rank: 1,
		Count: 1,
		MaxLen: 0,
	}
	for _, opt := range options {
		opt(cmd)
	}

	list, err := f.readList(key)
	if err != nil {
		return nil, err
	}

	if cmd.Rank == 0 {
		return nil, errors.New("ERR RANK can't be zero")
	}

	startIndex := 0
	if cmd.MaxLen > 0 && cmd.MaxLen < len(list) {
		if cmd.Rank > 0 {
			list = list[:cmd.MaxLen]
		} else {
			startIndex = len(list) - cmd.MaxLen
			list = list[startIndex:]
		}
	}

	var results []int
	matches := 0

	if cmd.Rank > 0 {
		for i := 0; i < len(list); i++ {
			if list[i] == element {
				matches++
				if matches >= cmd.Rank {
					results = append(results, startIndex + i)
					if cmd.Count > 0 && len(results) >= cmd.Count {
						break
					}
				}
			}
		}
	} else {
		targetRank := -cmd.Rank
		for i := len(list) - 1; i >= 0; i-- {
			if list[i] == element {
				matches++
				if matches >= targetRank {
					results = append(results, startIndex + i)
					if cmd.Count > 0 && len(results) >= cmd.Count {
						break
					}
				}
			}
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func (f *fileSystem) LPush(key dotpip.Key, elements ...string) (int, error) {
	if len(elements) == 0 {
		return f.LLen(key)
	}

	list, err := f.readList(key)
	if err != nil {
		return 0, err
	}

	// Prepend elements, note that LPush reverses order when pushing multiple?
	// Redis: LPUSH mylist a b c results in c, b, a
	for _, el := range elements {
		list = append([]string{el}, list...)
	}

	err = f.writeList(key, list)
	if err != nil {
		return 0, err
	}

	return len(list), nil
}

func (f *fileSystem) LPushX(key dotpip.Key, elements ...string) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, nil
	}
	return f.LPush(key, elements...)
}

func (f *fileSystem) LRange(key dotpip.Key, start int, stop int) ([]string, error) {
	list, err := f.readList(key)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return []string{}, nil
	}

	length := len(list)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}
	if start >= length {
		return []string{}, nil
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return []string{}, nil
	}

	return list[start : stop+1], nil
}

func (f *fileSystem) LRem(key dotpip.Key, count int, element string) (int, error) {
	list, err := f.readList(key)
	if err != nil {
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil
	}

	removed := 0
	var newList []string

	if count > 0 {
		for _, v := range list {
			if v == element && removed < count {
				removed++
			} else {
				newList = append(newList, v)
			}
		}
	} else if count < 0 {
		targetCount := -count
		// Iterate backwards to remove from end
		for i := len(list) - 1; i >= 0; i-- {
			if list[i] == element && removed < targetCount {
				removed++
			} else {
				newList = append([]string{list[i]}, newList...)
			}
		}
	} else {
		for _, v := range list {
			if v == element {
				removed++
			} else {
				newList = append(newList, v)
			}
		}
	}

	if removed > 0 {
		err = f.writeList(key, newList)
		if err != nil {
			return 0, err
		}
	}

	return removed, nil
}

func (f *fileSystem) LSet(key dotpip.Key, index int, element string) error {
	list, err := f.readList(key)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return errors.New("ERR no such key")
	}

	length := len(list)
	if index < 0 {
		index = length + index
	}

	if index < 0 || index >= length {
		return errors.New("ERR index out of range")
	}

	list[index] = element
	return f.writeList(key, list)
}

func (f *fileSystem) LTrim(key dotpip.Key, start int, stop int) error {
	list, err := f.readList(key)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}

	length := len(list)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}

	if start >= length || start > stop {
		return f.writeList(key, []string{}) // Empty list removes key
	}
	if stop >= length {
		stop = length - 1
	}

	return f.writeList(key, list[start:stop+1])
}

func (f *fileSystem) RPop(key dotpip.Key, count int) ([]string, error) {
	list, err := f.readList(key)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	if count <= 0 {
		count = 1
	}

	if count > len(list) {
		count = len(list)
	}

	// RPOP returns items starting from the rightmost
	var popped []string
	for i := 0; i < count; i++ {
		idx := len(list) - 1 - i
		popped = append(popped, list[idx])
	}

	list = list[:len(list)-count]

	err = f.writeList(key, list)
	if err != nil {
		return nil, err
	}

	return popped, nil
}

func (f *fileSystem) RPush(key dotpip.Key, elements ...string) (int, error) {
	if len(elements) == 0 {
		return f.LLen(key)
	}

	list, err := f.readList(key)
	if err != nil {
		return 0, err
	}

	list = append(list, elements...)

	err = f.writeList(key, list)
	if err != nil {
		return 0, err
	}

	return len(list), nil
}

func (f *fileSystem) RPushX(key dotpip.Key, elements ...string) (int, error) {
	exist, err := f.checkExistByKey(key)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, nil
	}
	return f.RPush(key, elements...)
}
