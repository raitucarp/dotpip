package fs

import (
	"dotpip"
	"math/rand"
	"sort"
	"strconv"
)

func (f *fileSystem) readSortedSet(key dotpip.Key) (map[string]float64, error) {
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

func (f *fileSystem) writeSortedSet(key dotpip.Key, zset map[string]float64) error {
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

func (f *fileSystem) ZAdd(key dotpip.Key, members []dotpip.Z, options ...dotpip.ZAddOption) (int, error) {
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
			} else {
				if score != member.Score {
					zset[member.Member] = member.Score
					changed++
				}
			}
		} else {
			if cmd.XX {
				continue
			}
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

func (f *fileSystem) ZCard(key dotpip.Key) (int, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}
	return len(zset), nil
}

func (f *fileSystem) ZCount(key dotpip.Key, min float64, max float64) (int, error) {
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, score := range zset {
		if score >= min && score <= max {
			count++
		}
	}
	return count, nil
}

func (f *fileSystem) ZDiff(keys ...dotpip.Key) ([]string, error) {
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

func (f *fileSystem) ZDiffWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZIncrBy(key dotpip.Key, increment float64, member string) (float64, error) {
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

func (f *fileSystem) ZInter(keys ...dotpip.Key) ([]string, error) {
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

func (f *fileSystem) ZInterWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZLexCount(key dotpip.Key, min string, max string) (int, error) {
	// A simplified implementation. Redis supports [min, (min, +, -
	zset, err := f.readSortedSet(key)
	if err != nil {
		return 0, err
	}

	count := 0
	for member := range zset {
		if (min == "-" || member >= min) && (max == "+" || member <= max) {
			count++
		}
	}
	return count, nil
}

func (f *fileSystem) getSortedZSet(key dotpip.Key) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZPopMax(key dotpip.Key, count int) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZPopMin(key dotpip.Key, count int) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZRandMember(key dotpip.Key, count int) ([]string, error) {
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

func (f *fileSystem) ZRandMemberWithScores(key dotpip.Key, count int) ([]dotpip.Z, error) {
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

func (f *fileSystem) ZRange(key dotpip.Key, start string, stop string, options ...dotpip.ZRangeOption) ([]string, error) {
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

func (f *fileSystem) ZRangeWithScores(key dotpip.Key, start string, stop string, options ...dotpip.ZRangeOption) ([]dotpip.Z, error) {
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

	if cmd.ByScore {
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
	} else if cmd.ByLex {
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
	} else {
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

func (f *fileSystem) ZRank(key dotpip.Key, member string) (int, error) {
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

func (f *fileSystem) ZRem(key dotpip.Key, members ...string) (int, error) {
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

func (f *fileSystem) ZRevRank(key dotpip.Key, member string) (int, error) {
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

func (f *fileSystem) ZScore(key dotpip.Key, member string) (float64, error) {
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

func (f *fileSystem) ZUnion(keys ...dotpip.Key) ([]string, error) {
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

func (f *fileSystem) ZUnionWithScores(keys ...dotpip.Key) ([]dotpip.Z, error) {
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
