package fs

import (
	"dotpip"
	"math/rand"
)

func (f *fileSystem) readSet(key dotpip.Key) (map[string]any, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return make(map[string]any), nil // Return empty set if not exists
	}

	set, err := f.formatter.SetDecode(content)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return make(map[string]any), nil
	}

	return set, nil
}

func (f *fileSystem) writeSet(key dotpip.Key, set map[string]any) error {
	if len(set) == 0 {
		f.Del(key)
		return nil
	}
	content, err := f.formatter.SetEncode(set)
	if err != nil {
		return err
	}
	return f.writeFileByKey(key, content.([]byte))
}

func (f *fileSystem) SAdd(key dotpip.Key, members ...string) (int, error) {
	set, err := f.readSet(key)
	if err != nil {
		return 0, err
	}

	added := 0
	for _, member := range members {
		if _, exists := set[member]; !exists {
			set[member] = struct{}{}
			added++
		}
	}

	if added > 0 {
		err = f.writeSet(key, set)
		if err != nil {
			return 0, err
		}
	}

	return added, nil
}

func (f *fileSystem) SCard(key dotpip.Key) (int, error) {
	set, err := f.readSet(key)
	if err != nil {
		return 0, err
	}
	return len(set), nil
}

func (f *fileSystem) SDiff(keys ...dotpip.Key) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}

	baseSet, err := f.readSet(keys[0])
	if err != nil {
		return nil, err
	}

	diff := make(map[string]any)
	for k, v := range baseSet {
		diff[k] = v
	}

	for _, key := range keys[1:] {
		set, err := f.readSet(key)
		if err != nil {
			return nil, err
		}
		for member := range set {
			delete(diff, member)
		}
	}

	var result []string
	for member := range diff {
		result = append(result, member)
	}

	return result, nil
}

func (f *fileSystem) SDiffStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	diff, err := f.SDiff(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range diff {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err != nil {
		return 0, err
	}

	return len(diff), nil
}

func (f *fileSystem) SInter(keys ...dotpip.Key) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}

	baseSet, err := f.readSet(keys[0])
	if err != nil {
		return nil, err
	}

	inter := make(map[string]any)
	for k, v := range baseSet {
		inter[k] = v
	}

	for _, key := range keys[1:] {
		set, err := f.readSet(key)
		if err != nil {
			return nil, err
		}
		for member := range inter {
			if _, exists := set[member]; !exists {
				delete(inter, member)
			}
		}
	}

	var result []string
	for member := range inter {
		result = append(result, member)
	}

	return result, nil
}

func (f *fileSystem) SInterCard(limit int, keys ...dotpip.Key) (int, error) {
	inter, err := f.SInter(keys...)
	if err != nil {
		return 0, err
	}

	count := len(inter)
	if limit > 0 && count > limit {
		count = limit
	}

	return count, nil
}

func (f *fileSystem) SInterStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	inter, err := f.SInter(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range inter {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err != nil {
		return 0, err
	}

	return len(inter), nil
}

func (f *fileSystem) SIsMember(key dotpip.Key, member string) (bool, error) {
	set, err := f.readSet(key)
	if err != nil {
		return false, err
	}
	_, exists := set[member]
	return exists, nil
}

func (f *fileSystem) SMembers(key dotpip.Key) ([]string, error) {
	set, err := f.readSet(key)
	if err != nil {
		return nil, err
	}

	var result []string
	for member := range set {
		result = append(result, member)
	}

	return result, nil
}

func (f *fileSystem) SMIsMember(key dotpip.Key, members ...string) ([]bool, error) {
	set, err := f.readSet(key)
	if err != nil {
		return nil, err
	}

	result := make([]bool, len(members))
	for i, member := range members {
		_, exists := set[member]
		result[i] = exists
	}

	return result, nil
}

func (f *fileSystem) SMove(source dotpip.Key, destination dotpip.Key, member string) (bool, error) {
	sourceSet, err := f.readSet(source)
	if err != nil {
		return false, err
	}

	if _, exists := sourceSet[member]; !exists {
		return false, nil
	}

	destSet, err := f.readSet(destination)
	if err != nil {
		return false, err
	}

	delete(sourceSet, member)
	destSet[member] = struct{}{}

	err = f.writeSet(source, sourceSet)
	if err != nil {
		return false, err
	}

	err = f.writeSet(destination, destSet)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (f *fileSystem) SPop(key dotpip.Key, count int) ([]string, error) {
	set, err := f.readSet(key)
	if err != nil {
		return nil, err
	}

	if len(set) == 0 {
		return []string{}, nil
	}

	if count <= 0 {
		return []string{}, nil
	}

	var members []string
	for member := range set {
		members = append(members, member)
	}

	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	limit := count
	if limit > len(members) {
		limit = len(members)
	}

	popped := members[:limit]
	for _, member := range popped {
		delete(set, member)
	}

	err = f.writeSet(key, set)
	if err != nil {
		return nil, err
	}

	return popped, nil
}

func (f *fileSystem) SRandMember(key dotpip.Key, count int) ([]string, error) {
	set, err := f.readSet(key)
	if err != nil {
		return nil, err
	}

	if len(set) == 0 {
		return []string{}, nil
	}

	var members []string
	for member := range set {
		members = append(members, member)
	}

	if count == 0 {
		return []string{}, nil
	}

	if count < 0 {
		// allow duplicates
		absCount := -count
		result := make([]string, absCount)
		for i := 0; i < absCount; i++ {
			result[i] = members[rand.Intn(len(members))]
		}
		return result, nil
	}

	// no duplicates
	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	limit := count
	if limit > len(members) {
		limit = len(members)
	}

	return members[:limit], nil
}

func (f *fileSystem) SRem(key dotpip.Key, members ...string) (int, error) {
	set, err := f.readSet(key)
	if err != nil {
		return 0, err
	}

	removed := 0
	for _, member := range members {
		if _, exists := set[member]; exists {
			delete(set, member)
			removed++
		}
	}

	if removed > 0 {
		err = f.writeSet(key, set)
		if err != nil {
			return 0, err
		}
	}

	return removed, nil
}

func (f *fileSystem) SUnion(keys ...dotpip.Key) ([]string, error) {
	union := make(map[string]any)

	for _, key := range keys {
		set, err := f.readSet(key)
		if err != nil {
			return nil, err
		}
		for member := range set {
			union[member] = struct{}{}
		}
	}

	var result []string
	for member := range union {
		result = append(result, member)
	}

	return result, nil
}

func (f *fileSystem) SUnionStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	union, err := f.SUnion(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range union {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err != nil {
		return 0, err
	}

	return len(union), nil
}
