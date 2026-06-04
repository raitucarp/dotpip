package fs

import (
	"dotpip"
	"math/rand"
	"sort"
)

func (f *FileSystem) readSet(key dotpip.Key) (map[string]any, error) {
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

func (f *FileSystem) writeSet(key dotpip.Key, set map[string]any) error {
	if len(set) == 0 {
		f.Del(key) // Del internally emits "del" event already!
		return nil
	}
	content, err := f.formatter.SetEncode(set)
	if err != nil {
		return err
	}
	return f.writeFileByKey(key, content.([]byte))
}

func (f *FileSystem) SAdd(key dotpip.Key, members ...string) (int, error) {
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
		if err == nil {
			f.emitKeyspaceEvent(key, "sadd", 's')
		}
		if err != nil {
			return 0, err
		}
	}

	return added, nil
}

func (f *FileSystem) SCard(key dotpip.Key) (int, error) {
	set, err := f.readSet(key)
	if err != nil {
		return 0, err
	}
	return len(set), nil
}

func (f *FileSystem) SDiff(keys ...dotpip.Key) ([]string, error) {
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

func (f *FileSystem) SDiffStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	diff, err := f.SDiff(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range diff {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err == nil {
		f.emitKeyspaceEvent(destination, "sdiffstore", 's')
	}
	if err != nil {
		return 0, err
	}

	return len(diff), nil
}

func (f *FileSystem) SInter(keys ...dotpip.Key) ([]string, error) {
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

func (f *FileSystem) SInterCard(limit int, keys ...dotpip.Key) (int, error) {
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

func (f *FileSystem) SInterStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	inter, err := f.SInter(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range inter {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err == nil {
		f.emitKeyspaceEvent(destination, "sinterstore", 's')
	}
	if err != nil {
		return 0, err
	}

	return len(inter), nil
}

func (f *FileSystem) SIsMember(key dotpip.Key, member string) (bool, error) {
	set, err := f.readSet(key)
	if err != nil {
		return false, err
	}
	_, exists := set[member]
	return exists, nil
}

func (f *FileSystem) SMembers(key dotpip.Key) ([]string, error) {
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

func (f *FileSystem) SMIsMember(key dotpip.Key, members ...string) ([]bool, error) {
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

func (f *FileSystem) SMove(source dotpip.Key, destination dotpip.Key, member string) (bool, error) {
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
	if err == nil {
		f.emitKeyspaceEvent(source, "srem", 's')
		f.emitKeyspaceEvent(destination, "sadd", 's')
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (f *FileSystem) SPop(key dotpip.Key, count int) ([]string, error) {
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
	if err == nil {
		f.emitKeyspaceEvent(key, "spop", 's')
	}
	if err != nil {
		return nil, err
	}

	return popped, nil
}

func (f *FileSystem) SRandMember(key dotpip.Key, count int) ([]string, error) {
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

func (f *FileSystem) SRem(key dotpip.Key, members ...string) (int, error) {
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
		if err == nil {
			f.emitKeyspaceEvent(key, "srem", 's')
		}
		if err != nil {
			return 0, err
		}
	}

	return removed, nil
}

func (f *FileSystem) SUnion(keys ...dotpip.Key) ([]string, error) {
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

func (f *FileSystem) SUnionStore(destination dotpip.Key, keys ...dotpip.Key) (int, error) {
	union, err := f.SUnion(keys...)
	if err != nil {
		return 0, err
	}

	set := make(map[string]any)
	for _, member := range union {
		set[member] = struct{}{}
	}

	err = f.writeSet(destination, set)
	if err == nil {
		f.emitKeyspaceEvent(destination, "sunionstore", 's')
	}
	if err != nil {
		return 0, err
	}

	return len(union), nil
}

func (f *FileSystem) SScan(key dotpip.Key, cursor uint64, options ...dotpip.ScanOption) (uint64, []string, error) {
	cmd := &dotpip.ScanCommand{Count: 10}
	for _, option := range options {
		option(cmd)
	}

	set, err := f.readSet(key)
	if err != nil {
		return 0, nil, err
	}

	var allMembers []string
	for member := range set {
		if cmd.Match != "" && !matchPattern(cmd.Match, member) {
			continue
		}
		allMembers = append(allMembers, member)
	}

	sort.Strings(allMembers)

	if cursor >= uint64(len(allMembers)) {
		return 0, []string{}, nil
	}

	end := cursor + uint64(cmd.Count)
	nextCursor := end
	if end >= uint64(len(allMembers)) {
		end = uint64(len(allMembers))
		nextCursor = 0
	}

	return nextCursor, allMembers[cursor:end], nil
}
