package fs

import (
	"dotpip"
	"math/rand"
	"sort"
	"strconv"
)

func (f *FileSystem) readSortedSet(key dotpip.Key) (map[string]float64, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return make(map[string]float64), nil
	}

	zset, err := f.formatter.SortedSetDecode(content)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return make(map[string]float64), nil
	}

	return zset, nil
}

func (f *FileSystem) writeSortedSet(key dotpip.Key, zset map[string]float64) error {
	if len(zset) == 0 {
		f.Del(key)
		return nil
	}
	content, err := f.formatter.SortedSetEncode(zset)
	if err != nil {
		return err
	}
	return f.writeFileByKey(key, content.([]byte))
}

// ZAdd adds one or more members to a sorted set.
func (f *FileSystem) ZAdd(key dotpip.Key, members []dotpip.Z, options ...dotpip.ZAddOption) (int, error) {
	cmd := &dotpip.ZAddCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	added := 0
	changed := 0

	for _, member := range members {
		score, exists := zset[member.Member]

		if exists {
			if cmd.NX {
				continue
			}
			if cmd.LT && member.Score >= score {
				continue
			}
			if cmd.GT && member.Score <= score {
				continue
			}

			if cmd.INCR {
				newScore := score + member.Score
				if newScore != score {
					zset[member.Member] = newScore
					changed++
				}
			} else if score != member.Score {
				zset[member.Member] = member.Score
				changed++
			}
		} else if !cmd.XX {
			zset[member.Member] = member.Score
			added++
		}
	}

	if added > 0 || changed > 0 {
		err = f.writeSortedSet(key, zset)
		if err != nil {
			return 0, err
		}
	}

	if cmd.CH {
		return added + changed, nil
	}
	return added, nil
}

// ZCard returns the number of members in a sorted set.
func (f *FileSystem) ZCard(key dotpip.Key) (int, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}
	return len(zset), nil
}

// ZCount returns the number of members in a sorted set with scores within the given values.
func (f *FileSystem) ZCount(key dotpip.Key, minVal float64, maxVal float64) (int, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, score := range zset {
		if score >= minVal && score <= maxVal {
			count++
		}
	}
	return count, nil
}

// ZDiff returns the difference between multiple sorted sets.
func (f *FileSystem) ZDiff(keys ...dotpip.Key) ([]string, error) {
	diff, err := f.ZDiffWithScores(keys...)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, z := range diff {
		res = append(res, z.Member)
	}
	return res, nil
}

// ZDiffWithScores returns the difference between multiple sorted sets with scores.
func (f *FileSystem) ZDiffWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
	if len(keys) == 0 {
		return []dotpip.Z{}, nil
	}

	baseZset, err := f.readSortedSet(keys[0])
	if err != nil {
		return nil, err
	}

	diff := make(map[string]float64)
	for k, v := range baseZset {
		diff[k] = v
	}

	for _, key := range keys[1:] {
		zset, err := f.readSortedSet(key)
		if err != nil {
			return nil, err
		}
		for member := range zset {
			delete(diff, member)
		}
	}

	var result []dotpip.Z
	for member, score := range diff {
		result = append(result, dotpip.Z{Member: member, Score: score})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Member < result[j].Member
		}
		return result[i].Score < result[j].Score
	})

	return result, nil
}

// ZIncrBy increments the score of a member in a sorted set.
func (f *FileSystem) ZIncrBy(key dotpip.Key, increment float64, member string) (float64, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	score := zset[member]
	newScore := score + increment
	zset[member] = newScore

	err = f.writeSortedSet(key, zset)
	if err != nil {
		return 0, err
	}

	return newScore, nil
}

// ZInter returns the intersection of multiple sorted sets.
func (f *FileSystem) ZInter(keys ...dotpip.Key) ([]string, error) {
	inter, err := f.ZInterWithScores(keys...)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, z := range inter {
		res = append(res, z.Member)
	}
	return res, nil
}

// ZInterWithScores returns the intersection of multiple sorted sets with scores.
func (f *FileSystem) ZInterWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
	if len(keys) == 0 {
		return []dotpip.Z{}, nil
	}

	baseZset, err := f.readSortedSet(keys[0])
	if err != nil {
		return nil, err
	}

	inter := make(map[string]float64)
	for k, v := range baseZset {
		inter[k] = v
	}

	for _, key := range keys[1:] {
		zset, err := f.readSortedSet(key)
		if err != nil {
			return nil, err
		}
		for member := range inter {
			if score, exists := zset[member]; !exists {
				delete(inter, member)
			} else {
				// SUM is the default for REDIS ZINTER
				inter[member] += score
			}
		}
	}

	var result []dotpip.Z
	for member, score := range inter {
		result = append(result, dotpip.Z{Member: member, Score: score})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Member < result[j].Member
		}
		return result[i].Score < result[j].Score
	})

	return result, nil
}

// ZLexCount returns the number of members in a sorted set within the given lexicographical range.
func (f *FileSystem) ZLexCount(key dotpip.Key, minVal string, maxVal string) (int, error) {
	// A simplified implementation. Redis supports [min, (min, +, -
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	count := 0
	for member := range zset {
		if (minVal == "-" || member >= minVal) && (maxVal == "+" || member <= maxVal) {
			count++
		}
	}
	return count, nil
}

func (f *FileSystem) getSortedZSet(key dotpip.Key) ([]dotpip.Z, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return nil, err
	}

	var result []dotpip.Z
	for member, score := range zset {
		result = append(result, dotpip.Z{Member: member, Score: score})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Member < result[j].Member
		}
		return result[i].Score < result[j].Score
	})

	return result, nil
}

// ZPopMax removes and returns the members with the highest scores in a sorted set.
func (f *FileSystem) ZPopMax(key dotpip.Key, count int) ([]dotpip.Z, error) {
	sorted, err := f.getSortedZSet(key)
	if err != nil {
		return nil, err
	}

	if count <= 0 || len(sorted) == 0 {
		return []dotpip.Z{}, nil
	}

	limit := count
	if limit > len(sorted) {
		limit = len(sorted)
	}

	popped := []dotpip.Z{}
	zset, _ := f.readSortedSet(key)

	// Pop from the end
	for i := 0; i < limit; i++ {
		z := sorted[len(sorted)-1-i]
		popped = append(popped, z)
		delete(zset, z.Member)
	}

	err = f.writeSortedSet(key, zset)
	if err != nil {
		return nil, err
	}

	return popped, nil
}

// ZPopMin removes and returns the members with the lowest scores in a sorted set.
func (f *FileSystem) ZPopMin(key dotpip.Key, count int) ([]dotpip.Z, error) {
	sorted, err := f.getSortedZSet(key)
	if err != nil {
		return nil, err
	}

	if count <= 0 || len(sorted) == 0 {
		return []dotpip.Z{}, nil
	}

	limit := count
	if limit > len(sorted) {
		limit = len(sorted)
	}

	popped := []dotpip.Z{}
	zset, _ := f.readSortedSet(key)

	// Pop from the beginning
	for i := 0; i < limit; i++ {
		z := sorted[i]
		popped = append(popped, z)
		delete(zset, z.Member)
	}

	err = f.writeSortedSet(key, zset)
	if err != nil {
		return nil, err
	}

	return popped, nil
}

// ZRandMember returns random members from a sorted set.
func (f *FileSystem) ZRandMember(key dotpip.Key, count int) ([]string, error) {
	randZ, err := f.ZRandMemberWithScores(key, count)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, z := range randZ {
		res = append(res, z.Member)
	}
	return res, nil
}

// ZRandMemberWithScores returns random members with scores from a sorted set.
func (f *FileSystem) ZRandMemberWithScores(key dotpip.Key, count int) ([]dotpip.Z, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return nil, err
	}

	if len(zset) == 0 {
		return []dotpip.Z{}, nil
	}

	var members []dotpip.Z
	for member, score := range zset {
		members = append(members, dotpip.Z{Member: member, Score: score})
	}

	if count == 0 {
		return []dotpip.Z{}, nil
	}

	if count < 0 {
		// allow duplicates
		absCount := -count
		result := make([]dotpip.Z, absCount)
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

// ZRange returns members in a sorted set within a range.
func (f *FileSystem) ZRange(key dotpip.Key, start string, stop string, options ...dotpip.ZRangeOption) ([]string, error) {
	resZ, err := f.ZRangeWithScores(key, start, stop, options...)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, z := range resZ {
		res = append(res, z.Member)
	}
	return res, nil
}

// ZRangeWithScores returns members with scores in a sorted set within a range.
func (f *FileSystem) ZRangeWithScores(key dotpip.Key, start string, stop string, options ...dotpip.ZRangeOption) ([]dotpip.Z, error) {
	sorted, err := f.getSortedZSet(key)
	if err != nil {
		return nil, err
	}

	if len(sorted) == 0 {
		return []dotpip.Z{}, nil
	}

	cmd := &dotpip.ZRangeCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	if cmd.Rev {
		// Redis reverses the entire list first if REV is provided
		for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
			sorted[i], sorted[j] = sorted[j], sorted[i]
		}
	}

	switch {
	case cmd.ByScore:
		// Start and stop are scores
		startScore, _ := strconv.ParseFloat(start, 64)
		stopScore, _ := strconv.ParseFloat(stop, 64)

		var filtered []dotpip.Z
		for _, z := range sorted {
			// When REV is used, Redis checks if it's within start/stop properly,
			// but we will do a simple bounds check here.
			if cmd.Rev {
				if z.Score >= stopScore && z.Score <= startScore {
					filtered = append(filtered, z)
				}
			} else {
				if z.Score >= startScore && z.Score <= stopScore {
					filtered = append(filtered, z)
				}
			}
		}
		sorted = filtered
	case cmd.ByLex:
		// Start and stop are strings
		var filtered []dotpip.Z
		for _, z := range sorted {
			// Simplified lex matching
			if cmd.Rev {
				if (stop == "-" || z.Member >= stop) && (start == "+" || z.Member <= start) {
					filtered = append(filtered, z)
				}
			} else {
				if (start == "-" || z.Member >= start) && (stop == "+" || z.Member <= stop) {
					filtered = append(filtered, z)
				}
			}
		}
		sorted = filtered
	default:
		// Start and stop are indices
		startIdx, _ := strconv.Atoi(start)
		stopIdx, _ := strconv.Atoi(stop)

		if startIdx < 0 {
			startIdx = len(sorted) + startIdx
		}
		if stopIdx < 0 {
			stopIdx = len(sorted) + stopIdx
		}

		if startIdx < 0 {
			startIdx = 0
		}
		if stopIdx >= len(sorted) {
			stopIdx = len(sorted) - 1
		}

		if startIdx > stopIdx || startIdx >= len(sorted) {
			sorted = []dotpip.Z{}
		} else {
			sorted = sorted[startIdx : stopIdx+1]
		}
	}

	if cmd.Limit {
		if cmd.Offset >= len(sorted) {
			return []dotpip.Z{}, nil
		}

		end := cmd.Offset + cmd.Count
		if end > len(sorted) || cmd.Count < 0 {
			end = len(sorted)
		}

		return sorted[cmd.Offset:end], nil
	}

	return sorted, nil
}

// ZRank returns the rank of a member in a sorted set.
func (f *FileSystem) ZRank(key dotpip.Key, member string) (int, error) {
	sorted, err := f.getSortedZSet(key)
	if err != nil {
		return -1, err
	}

	for i, z := range sorted {
		if z.Member == member {
			return i, nil
		}
	}
	return -1, nil // not found
}

// ZRem removes one or more members from a sorted set.
func (f *FileSystem) ZRem(key dotpip.Key, members ...string) (int, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	removed := 0
	for _, member := range members {
		if _, exists := zset[member]; exists {
			delete(zset, member)
			removed++
		}
	}

	if removed > 0 {
		err = f.writeSortedSet(key, zset)
		if err != nil {
			return 0, err
		}
	}

	return removed, nil
}

// ZRevRank returns the reverse rank of a member in a sorted set.
func (f *FileSystem) ZRevRank(key dotpip.Key, member string) (int, error) {
	sorted, err := f.getSortedZSet(key)
	if err != nil {
		return -1, err
	}

	for i, z := range sorted {
		if z.Member == member {
			return len(sorted) - 1 - i, nil
		}
	}
	return -1, nil
}

// ZScore returns the score of a member in a sorted set.
func (f *FileSystem) ZScore(key dotpip.Key, member string) (float64, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	if score, exists := zset[member]; exists {
		return score, nil
	}
	// Return 0 if not exists? In real Redis it returns nil
	return 0, nil
}

// ZUnion returns the union of multiple sorted sets.
func (f *FileSystem) ZUnion(keys ...dotpip.Key) ([]string, error) {
	union, err := f.ZUnionWithScores(keys...)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, z := range union {
		res = append(res, z.Member)
	}
	return res, nil
}

// ZUnionWithScores returns the union of multiple sorted sets with scores.
func (f *FileSystem) ZUnionWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
	union := make(map[string]float64)

	for _, key := range keys {
		zset, err := f.readSortedSet(key)
		if err != nil {
			return nil, err
		}
		for member, score := range zset {
			union[member] += score // SUM default
		}
	}

	var result []dotpip.Z
	for member, score := range union {
		result = append(result, dotpip.Z{Member: member, Score: score})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Member < result[j].Member
		}
		return result[i].Score < result[j].Score
	})

	return result, nil
}

// ZScan incrementally iterates over a sorted set.
func (f *FileSystem) ZScan(key dotpip.Key, cursor uint64, options ...dotpip.ScanOption) (uint64, []dotpip.Z, error) {
	cmd := &dotpip.ScanCommand{Count: 10}
	for _, option := range options {
		option(cmd)
	}

	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, nil, err
	}

	var allMembers []string
	for member := range zset {
		if cmd.Match != "" && !matchPattern(cmd.Match, member) {
			continue
		}
		allMembers = append(allMembers, member)
	}

	sort.Strings(allMembers)

	if cursor >= uint64(len(allMembers)) {
		return 0, []dotpip.Z{}, nil
	}

	end := cursor + uint64(cmd.Count)
	nextCursor := end
	if end >= uint64(len(allMembers)) {
		end = uint64(len(allMembers))
		nextCursor = 0
	}

	var result []dotpip.Z
	for _, member := range allMembers[cursor:end] {
		result = append(result, dotpip.Z{
			Member: member,
			Score:  zset[member],
		})
	}

	return nextCursor, result, nil
}
