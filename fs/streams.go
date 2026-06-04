package fs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dotpip"
)

func (f *FileSystem) parseStreamID(id string) (int64, int64, error) {
	if id == "*" {
		return time.Now().UnixMilli(), 0, nil
	}
	parts := strings.Split(id, "-")
	if len(parts) > 2 {
		return 0, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidStreamID))
	}

	ms, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidStreamID))
	}

	var seq int64
	if len(parts) == 2 {
		if parts[1] == "*" {
			// handled dynamically based on last ID
			seq = -1
		} else {
			seq, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidStreamID))
			}
		}
	} else {
		// if no dash provided, assuming sequence is 0? Actually redis requires seq unless using * for it
		// wait, if "1526919030474", seq is implicitly 0?
		// No, usually auto generated. But for strictness let's assume it means 0 if no star.
		// Wait, redis says if only time is given, seq is 0.
		seq = 0
	}

	return ms, seq, nil
}

func compareIDs(ms1, seq1, ms2, seq2 int64) int {
	if ms1 > ms2 {
		return 1
	}
	if ms1 < ms2 {
		return -1
	}
	if seq1 > seq2 {
		return 1
	}
	if seq1 < seq2 {
		return -1
	}
	return 0
}

func parseIDString(id string) (int64, int64) {
	parts := strings.Split(id, "-")
	ms, _ := strconv.ParseInt(parts[0], 10, 64)
	seq, _ := strconv.ParseInt(parts[1], 10, 64)
	return ms, seq
}

func (f *FileSystem) getStream(key dotpip.Key) (dotpip.Stream, error) {
	var stream dotpip.Stream
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return dotpip.Stream{
				Entries: make([]dotpip.StreamEntry, 0),
				Groups:  make(map[string]dotpip.StreamGroup),
			}, nil
		}
		return stream, err
	}
	if len(content) == 0 {
		return dotpip.Stream{
			Entries: make([]dotpip.StreamEntry, 0),
			Groups:  make(map[string]dotpip.StreamGroup),
		}, nil
	}
	return f.formatter.StreamDecode(content)
}

func (f *FileSystem) setStream(key dotpip.Key, stream dotpip.Stream) error {
	encoded, err := f.formatter.StreamEncode(stream)
	if err != nil {
		return err
	}
	return f.writeFileByKey(key, encoded.([]byte))
}

func (f *FileSystem) XAck(key dotpip.Key, group string, ids ...string) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return 0, nil // group doesn't exist
	}

	acked := 0
	for _, id := range ids {
		if pendingEntry, exists := g.Pending[id]; exists {
			delete(g.Pending, id)

			// remove from consumer's pending list as well
			if cons, consOk := g.Consumers[pendingEntry.Consumer]; consOk {
				delete(cons.Pending, id)
				g.Consumers[pendingEntry.Consumer] = cons
			}

			acked++
		}
	}

	if acked > 0 {
		stream.Groups[group] = g
		err = f.setStream(key, stream)
	}

	return acked, err
}

func (f *FileSystem) XAdd(key dotpip.Key, id string, values map[string]string, options ...dotpip.XAddOption) (string, error) {
	cmd := &dotpip.XAddCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	stream, err := f.getStream(key)
	if err != nil {
		return "", err
	}

	if cmd.NoMkStream && len(stream.Entries) == 0 {
		return "", nil // don't add
	}

	ms, seq, err := f.parseStreamID(id)
	if err != nil {
		return "", err
	}

	var lastMs, lastSeq int64 = 0, 0
	if len(stream.Entries) > 0 {
		lastEntry := stream.Entries[len(stream.Entries)-1]
		lastMs, lastSeq = parseIDString(lastEntry.ID)
	}

	switch {
	case id == "*":
		ms = time.Now().UnixMilli()
		switch {
		case ms == lastMs:
			seq = lastSeq + 1
		case ms < lastMs:
			ms = lastMs
			seq = lastSeq + 1
		default:
			seq = 0
		}
	case strings.HasSuffix(id, "-*"):
		switch {
		case ms == lastMs:
			seq = lastSeq + 1
		case ms < lastMs:
			return "", fmt.Errorf(string(dotpip.ErrMsgXAddIDEqualSmaller))
		default:
			seq = 0
		}
	default:
		if compareIDs(ms, seq, lastMs, lastSeq) <= 0 {
			if ms == 0 && seq == 0 {
				return "", fmt.Errorf(string(dotpip.ErrMsgXAddIDGreaterZero))
			}
			return "", fmt.Errorf(string(dotpip.ErrMsgXAddIDEqualSmaller))
		}
	}

	newID := fmt.Sprintf("%d-%d", ms, seq)

	stream.Entries = append(stream.Entries, dotpip.StreamEntry{
		ID:     newID,
		Values: values,
	})

	if cmd.MaxLen > 0 {
		if len(stream.Entries) > cmd.MaxLen {
			// strict trim for now
			stream.Entries = stream.Entries[len(stream.Entries)-cmd.MaxLen:]
		}
	}

	err = f.setStream(key, stream)
	if err != nil {
		return "", err
	}

	return newID, nil
}

func (f *FileSystem) XDel(key dotpip.Key, ids ...string) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	deleted := 0
	newEntries := make([]dotpip.StreamEntry, 0, len(stream.Entries))
	for _, entry := range stream.Entries {
		if idMap[entry.ID] {
			deleted++
		} else {
			newEntries = append(newEntries, entry)
		}
	}

	if deleted > 0 {
		stream.Entries = newEntries
		err = f.setStream(key, stream)
	}

	return deleted, err
}

func (f *FileSystem) XGroupCreate(key dotpip.Key, group string, id string, mkStream bool) (string, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return "", err
	}

	if len(stream.Entries) == 0 && !mkStream {
		return "", fmt.Errorf(string(dotpip.ErrMsgXGroupKeyExists))
	}

	if stream.Groups == nil {
		stream.Groups = make(map[string]dotpip.StreamGroup)
	}

	if _, exists := stream.Groups[group]; exists {
		return "", fmt.Errorf(string(dotpip.ErrMsgBusyGroup))
	}

	var lastDeliveredID string
	switch {
	case id == "$":
		if len(stream.Entries) > 0 {
			lastDeliveredID = stream.Entries[len(stream.Entries)-1].ID
		} else {
			lastDeliveredID = "0-0"
		}
	case id == "0":
		lastDeliveredID = "0-0"
	default:
		// Validating format
		_, _, err := f.parseStreamID(id)
		if err != nil {
			return "", err
		}
		// In case they pass just "123", we should convert to "123-0"
		if !strings.Contains(id, "-") {
			id += "-0"
		}
		lastDeliveredID = id
	}

	stream.Groups[group] = dotpip.StreamGroup{
		LastDeliveredID: lastDeliveredID,
		Pending:         make(map[string]dotpip.StreamPendingEntry),
		Consumers:       make(map[string]dotpip.StreamConsumer),
	}

	err = f.setStream(key, stream)
	if err != nil {
		return "", err
	}

	return string(dotpip.StatusOK), nil
}

func (f *FileSystem) XGroupCreateConsumer(key dotpip.Key, group string, consumer string) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return 0, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	if g.Consumers == nil {
		g.Consumers = make(map[string]dotpip.StreamConsumer)
	}

	if _, ok := g.Consumers[consumer]; ok {
		return 0, nil // Consumer already exists
	}

	g.Consumers[consumer] = dotpip.StreamConsumer{
		Pending: make(map[string]int64),
	}
	stream.Groups[group] = g

	err = f.setStream(key, stream)
	return 1, err
}

func (f *FileSystem) XGroupDelConsumer(key dotpip.Key, group string, consumer string) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return 0, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	cons, ok := g.Consumers[consumer]
	if !ok {
		return 0, nil // Consumer doesn't exist
	}

	pendingCount := len(cons.Pending)

	// Remove consumer's pending entries from the group's pending list
	for id := range cons.Pending {
		delete(g.Pending, id)
	}

	delete(g.Consumers, consumer)
	stream.Groups[group] = g

	err = f.setStream(key, stream)
	return pendingCount, err
}

func (f *FileSystem) XGroupDestroy(key dotpip.Key, group string) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	if _, ok := stream.Groups[group]; !ok {
		return 0, nil
	}

	delete(stream.Groups, group)
	err = f.setStream(key, stream)
	return 1, err
}

func (f *FileSystem) XGroupSetID(key dotpip.Key, group string, id string) (string, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return "", err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return "", fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	var lastDeliveredID string
	switch {
	case id == "$":
		if len(stream.Entries) > 0 {
			lastDeliveredID = stream.Entries[len(stream.Entries)-1].ID
		} else {
			lastDeliveredID = "0-0"
		}
	case id == "0":
		lastDeliveredID = "0-0"
	default:
		_, _, err := f.parseStreamID(id)
		if err != nil {
			return "", err
		}
		if !strings.Contains(id, "-") {
			id += "-0"
		}
		lastDeliveredID = id
	}

	g.LastDeliveredID = lastDeliveredID
	stream.Groups[group] = g

	err = f.setStream(key, stream)
	if err != nil {
		return "", err
	}

	return string(dotpip.StatusOK), nil
}

func (f *FileSystem) XLen(key dotpip.Key) (int, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}
	return len(stream.Entries), nil
}

func (f *FileSystem) XRange(key dotpip.Key, start string, end string, count int) ([]dotpip.StreamEntry, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	var startMs, startSeq int64
	if start != "-" {
		startMs, startSeq, err = f.parseStreamID(start)
		if err != nil {
			return nil, err
		}
	}

	var endMs, endSeq int64 = -1, -1
	if end != "+" {
		endMs, endSeq, err = f.parseStreamID(end)
		if err != nil {
			return nil, err
		}
		// if "123" is given as end, it means "123-max" which technically we can represent by not filtering seq if it's not provided, but parseStreamID makes it 0. Let's fix this for end string.
		if !strings.Contains(end, "-") {
			endSeq = 9223372036854775807 // max int64
		}
	}

	var results []dotpip.StreamEntry
	for _, entry := range stream.Entries {
		ms, seq := parseIDString(entry.ID)

		if start != "-" {
			if compareIDs(ms, seq, startMs, startSeq) < 0 {
				continue
			}
		}

		if end != "+" {
			if compareIDs(ms, seq, endMs, endSeq) > 0 {
				continue
			}
		}

		results = append(results, entry)
		if count > 0 && len(results) >= count {
			break
		}
	}

	return results, nil
}

func (f *FileSystem) XRevRange(key dotpip.Key, end string, start string, count int) ([]dotpip.StreamEntry, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	var endMs, endSeq int64 = -1, -1
	if end != "+" {
		endMs, endSeq, err = f.parseStreamID(end)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(end, "-") {
			endSeq = 9223372036854775807 // max int64
		}
	}

	var startMs, startSeq int64
	if start != "-" {
		startMs, startSeq, err = f.parseStreamID(start)
		if err != nil {
			return nil, err
		}
	}

	var results []dotpip.StreamEntry
	for i := len(stream.Entries) - 1; i >= 0; i-- {
		entry := stream.Entries[i]
		ms, seq := parseIDString(entry.ID)

		if end != "+" {
			if compareIDs(ms, seq, endMs, endSeq) > 0 {
				continue
			}
		}

		if start != "-" {
			if compareIDs(ms, seq, startMs, startSeq) < 0 {
				continue
			}
		}

		results = append(results, entry)
		if count > 0 && len(results) >= count {
			break
		}
	}

	return results, nil
}

func (f *FileSystem) XRead(keys []dotpip.Key, ids []string, options ...dotpip.XReadOption) (map[string][]dotpip.StreamEntry, error) {
	cmd := &dotpip.XReadCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	if len(keys) != len(ids) {
		return nil, fmt.Errorf(string(dotpip.ErrMsgUnbalancedXRead))
	}

	// For simplicity, blocking is not fully implemented with actual wait loops here.
	// Real redis blocks if data is empty. We will just return what we have right away.

	results := make(map[string][]dotpip.StreamEntry)

	for i, key := range keys {
		idStr := ids[i]

		stream, err := f.getStream(key)
		if err != nil {
			continue // skip or return err? Redis might just return empty for non-existent.
		}

		var startMs, startSeq int64
		if idStr == "$" {
			if len(stream.Entries) > 0 {
				startMs, startSeq = parseIDString(stream.Entries[len(stream.Entries)-1].ID)
			} else {
				startMs, startSeq = 0, 0
			}
		} else {
			startMs, startSeq, err = f.parseStreamID(idStr)
			if err != nil {
				return nil, err
			}
		}

		var entries []dotpip.StreamEntry
		for _, entry := range stream.Entries {
			ms, seq := parseIDString(entry.ID)

			// XRead strictly strictly greater than ID
			if compareIDs(ms, seq, startMs, startSeq) > 0 {
				entries = append(entries, entry)
				if cmd.Count > 0 && len(entries) >= cmd.Count {
					break
				}
			}
		}

		if len(entries) > 0 {
			keyStr := strings.Join(key, ".")
			results[keyStr] = entries
		}
	}

	return results, nil
}

func (f *FileSystem) XReadGroup(group string, consumer string, keys []dotpip.Key, ids []string, options ...dotpip.XReadGroupOption) (map[string][]dotpip.StreamEntry, error) {
	cmd := &dotpip.XReadGroupCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	if len(keys) != len(ids) {
		return nil, fmt.Errorf(string(dotpip.ErrMsgUnbalancedXRead))
	}

	results := make(map[string][]dotpip.StreamEntry)

	for i, key := range keys {
		idStr := ids[i]

		stream, err := f.getStream(key)
		if err != nil {
			continue
		}

		g, ok := stream.Groups[group]
		if !ok {
			return nil, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
		}

		// Ensure consumer exists
		if g.Consumers == nil {
			g.Consumers = make(map[string]dotpip.StreamConsumer)
		}
		cons, consOk := g.Consumers[consumer]
		if !consOk {
			cons = dotpip.StreamConsumer{Pending: make(map[string]int64)}
			g.Consumers[consumer] = cons
		}
		if cons.Pending == nil {
			cons.Pending = make(map[string]int64)
		}

		if g.Pending == nil {
			g.Pending = make(map[string]dotpip.StreamPendingEntry)
		}

		var entries []dotpip.StreamEntry

		if idStr == ">" {
			// Read new messages never delivered to other consumers
			lastMs, lastSeq := parseIDString(g.LastDeliveredID)

			for _, entry := range stream.Entries {
				ms, seq := parseIDString(entry.ID)

				if compareIDs(ms, seq, lastMs, lastSeq) > 0 {
					entries = append(entries, entry)

					// Update LastDeliveredID
					g.LastDeliveredID = entry.ID

					if !cmd.NoAck {
						// Add to pending
						now := time.Now().UnixMilli()
						g.Pending[entry.ID] = dotpip.StreamPendingEntry{
							Consumer:      consumer,
							DeliveryTime:  now,
							DeliveryCount: 1,
						}
						cons.Pending[entry.ID] = now
					}

					if cmd.Count > 0 && len(entries) >= cmd.Count {
						break
					}
				}
			}
		} else {
			// Read from pending messages for this consumer
			// "0" usually means read from the beginning of consumer's pending list
			var startMs, startSeq int64
			if idStr != "0" && idStr != "0-0" {
				startMs, startSeq, err = f.parseStreamID(idStr)
				if err != nil {
					return nil, err
				}
			}

			// We need to fetch entries that are in this consumer's pending list
			for _, entry := range stream.Entries {
				ms, seq := parseIDString(entry.ID)
				if compareIDs(ms, seq, startMs, startSeq) > 0 {
					if _, isPending := cons.Pending[entry.ID]; isPending {
						entries = append(entries, entry)

						// update delivery count and time
						if pEntry, pOk := g.Pending[entry.ID]; pOk {
							pEntry.DeliveryTime = time.Now().UnixMilli()
							pEntry.DeliveryCount++
							g.Pending[entry.ID] = pEntry
						}

						if cmd.Count > 0 && len(entries) >= cmd.Count {
							break
						}
					}
				}
			}
		}

		if len(entries) > 0 {
			keyStr := strings.Join(key, ".")
			results[keyStr] = entries

			g.Consumers[consumer] = cons
			stream.Groups[group] = g
			_ = f.setStream(key, stream)
		}
	}

	return results, nil
}

func (f *FileSystem) XTrim(key dotpip.Key, options ...dotpip.XTrimOption) (int, error) {
	cmd := &dotpip.XTrimCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	stream, err := f.getStream(key)
	if err != nil {
		return 0, err
	}

	initialLen := len(stream.Entries)

	if cmd.MaxLen > 0 {
		if len(stream.Entries) > cmd.MaxLen {
			// If approx is true, redis might not strictly trim to maxlen, but for simplicity we'll just trim exactly or slightly off.
			// exact trim for now since approx depends on internal structures like radix tree nodes.
			stream.Entries = stream.Entries[len(stream.Entries)-cmd.MaxLen:]
		}
	}

	if cmd.MinID != "" {
		minMs, minSeq, err := f.parseStreamID(cmd.MinID)
		if err != nil {
			return 0, err
		}

		var newEntries []dotpip.StreamEntry
		for _, entry := range stream.Entries {
			ms, seq := parseIDString(entry.ID)
			if compareIDs(ms, seq, minMs, minSeq) >= 0 {
				newEntries = append(newEntries, entry)
			}
		}
		stream.Entries = newEntries
	}

	if cmd.Limit > 0 {
		// Limit indicates the maximum number of entries to delete.
		deleted := initialLen - len(stream.Entries)
		if deleted > cmd.Limit {
			// Re-fetch stream and just trim `cmd.Limit` entries
			stream, _ = f.getStream(key)
			stream.Entries = stream.Entries[cmd.Limit:]
		}
	}

	deleted := initialLen - len(stream.Entries)
	if deleted > 0 {
		err = f.setStream(key, stream)
	}

	return deleted, err
}

func (f *FileSystem) XPending(key dotpip.Key, group string, options ...dotpip.XPendingOption) ([]any, error) {
	cmd := &dotpip.XPendingCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return nil, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	if cmd.Start == "" && cmd.End == "" && cmd.Count == 0 {
		// Summary mode
		totalPending := len(g.Pending)
		var minID, maxID string
		consumerStats := make(map[string]int)

		for id, pending := range g.Pending {
			if minID == "" || strings.Compare(id, minID) < 0 {
				minID = id
			}
			if maxID == "" || strings.Compare(id, maxID) > 0 {
				maxID = id
			}
			consumerStats[pending.Consumer]++
		}

		var stats []any
		for cons, count := range consumerStats {
			stats = append(stats, []any{cons, count})
		}

		return []any{totalPending, minID, maxID, stats}, nil
	}

	// Extended mode
	var startMs, startSeq int64
	if cmd.Start != "-" {
		startMs, startSeq, err = f.parseStreamID(cmd.Start)
		if err != nil {
			return nil, err
		}
	}

	var endMs, endSeq int64 = -1, -1
	if cmd.End != "+" {
		endMs, endSeq, err = f.parseStreamID(cmd.End)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(cmd.End, "-") {
			endSeq = 9223372036854775807
		}
	}

	var results []any
	now := time.Now().UnixMilli()

	// Needs to sort keys to maintain order, but for simplicity here we just iterate stream entries.
	for _, entry := range stream.Entries {
		ms, seq := parseIDString(entry.ID)

		if cmd.Start != "-" {
			if compareIDs(ms, seq, startMs, startSeq) < 0 {
				continue
			}
		}

		if cmd.End != "+" {
			if compareIDs(ms, seq, endMs, endSeq) > 0 {
				continue
			}
		}

		if pendingEntry, isPending := g.Pending[entry.ID]; isPending {
			idleTime := now - pendingEntry.DeliveryTime

			// Optional: filter by consumer... but wait, XPendingCommand doesn't have a consumer filter yet, redis allows it

			if cmd.Idle > 0 && idleTime < int64(cmd.Idle) {
				continue
			}

			results = append(results, []any{
				entry.ID,
				pendingEntry.Consumer,
				idleTime,
				pendingEntry.DeliveryCount,
			})

			if cmd.Count > 0 && len(results) >= cmd.Count {
				break
			}
		}
	}

	return results, nil
}

func (f *FileSystem) XClaim(key dotpip.Key, group string, consumer string, minIdleTime int, ids []string, options ...dotpip.XClaimOption) ([]dotpip.StreamEntry, error) {
	cmd := &dotpip.XClaimCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return nil, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	if g.Consumers == nil {
		g.Consumers = make(map[string]dotpip.StreamConsumer)
	}
	cons, consOk := g.Consumers[consumer]
	if !consOk {
		cons = dotpip.StreamConsumer{Pending: make(map[string]int64)}
	}
	if cons.Pending == nil {
		cons.Pending = make(map[string]int64)
	}

	now := time.Now().UnixMilli()
	if cmd.Time > 0 {
		now = cmd.Time
	}

	var entries []dotpip.StreamEntry
	changed := false

	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	for _, entry := range stream.Entries {
		if !idMap[entry.ID] {
			continue
		}

		pendingEntry, isPending := g.Pending[entry.ID]
		if !isPending && !cmd.Force {
			continue
		}

		idleTime := now - pendingEntry.DeliveryTime
		if isPending && idleTime < int64(minIdleTime) {
			continue
		}

		// Claim it!
		if isPending {
			// Remove from old consumer
			if oldCons, oldOk := g.Consumers[pendingEntry.Consumer]; oldOk {
				delete(oldCons.Pending, entry.ID)
				g.Consumers[pendingEntry.Consumer] = oldCons
			}
		}

		pendingEntry.Consumer = consumer
		if cmd.Idle > 0 {
			pendingEntry.DeliveryTime = now - int64(cmd.Idle)
		} else {
			pendingEntry.DeliveryTime = now
		}

		if cmd.RetryCount > 0 { // Or JustID? Actually we just reset or increment count.
			pendingEntry.DeliveryCount++
		}

		g.Pending[entry.ID] = pendingEntry
		cons.Pending[entry.ID] = pendingEntry.DeliveryTime
		entries = append(entries, entry)
		changed = true
	}

	if changed {
		g.Consumers[consumer] = cons
		stream.Groups[group] = g
		err = f.setStream(key, stream)
	}

	return entries, err
}

func (f *FileSystem) XAutoClaim(key dotpip.Key, group string, consumer string, minIdleTime int, start string, options ...dotpip.XAutoClaimOption) (string, []dotpip.StreamEntry, error) {
	cmd := &dotpip.XAutoClaimCommand{Count: 100}
	for _, opt := range options {
		opt(cmd)
	}

	stream, err := f.getStream(key)
	if err != nil {
		return "", nil, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return "", nil, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	if g.Consumers == nil {
		g.Consumers = make(map[string]dotpip.StreamConsumer)
	}
	cons, consOk := g.Consumers[consumer]
	if !consOk {
		cons = dotpip.StreamConsumer{Pending: make(map[string]int64)}
	}
	if cons.Pending == nil {
		cons.Pending = make(map[string]int64)
	}

	var startMs, startSeq int64
	if start != "-" && start != "0-0" {
		startMs, startSeq, err = f.parseStreamID(start)
		if err != nil {
			return "", nil, err
		}
	}

	now := time.Now().UnixMilli()
	var entries []dotpip.StreamEntry
	nextID := "0-0"
	changed := false

	count := 0
	for _, entry := range stream.Entries {
		ms, seq := parseIDString(entry.ID)

		if compareIDs(ms, seq, startMs, startSeq) < 0 {
			continue
		}

		pendingEntry, isPending := g.Pending[entry.ID]
		if !isPending {
			continue
		}

		idleTime := now - pendingEntry.DeliveryTime
		if idleTime < int64(minIdleTime) {
			continue
		}

		// Claim it!
		if oldCons, oldOk := g.Consumers[pendingEntry.Consumer]; oldOk {
			delete(oldCons.Pending, entry.ID)
			g.Consumers[pendingEntry.Consumer] = oldCons
		}

		pendingEntry.Consumer = consumer
		pendingEntry.DeliveryTime = now
		pendingEntry.DeliveryCount++

		g.Pending[entry.ID] = pendingEntry
		cons.Pending[entry.ID] = pendingEntry.DeliveryTime

		entries = append(entries, entry)
		changed = true

		count++
		if cmd.Count > 0 && count >= cmd.Count {
			// What is the next ID?
			nextID = entry.ID
			break
		}
	}

	if len(entries) > 0 && nextID == entries[len(entries)-1].ID {
		// Redis usually returns the ID just greater than the last one returned, or something similar
		// We'll approximate by finding the next entry, or appending 0
		ms, seq := parseIDString(nextID)
		nextID = fmt.Sprintf("%d-%d", ms, seq+1)
	} else if len(entries) == 0 {
		nextID = "0-0"
	}

	if changed {
		g.Consumers[consumer] = cons
		stream.Groups[group] = g
		err = f.setStream(key, stream)
	}

	return nextID, entries, err
}

func (f *FileSystem) XInfoStream(key dotpip.Key) (map[string]any, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	info := make(map[string]any)
	info["length"] = len(stream.Entries)
	info["groups"] = len(stream.Groups)

	if len(stream.Entries) > 0 {
		info["first-entry"] = stream.Entries[0]
		info["last-entry"] = stream.Entries[len(stream.Entries)-1]
	}

	return info, nil
}

func (f *FileSystem) XInfoGroups(key dotpip.Key) ([]map[string]any, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	var groups []map[string]any
	for name, g := range stream.Groups {
		groupInfo := make(map[string]any)
		groupInfo["name"] = name
		groupInfo["consumers"] = len(g.Consumers)
		groupInfo["pending"] = len(g.Pending)
		groupInfo["last-delivered-id"] = g.LastDeliveredID
		groups = append(groups, groupInfo)
	}

	return groups, nil
}

func (f *FileSystem) XInfoConsumers(key dotpip.Key, group string) ([]map[string]any, error) {
	stream, err := f.getStream(key)
	if err != nil {
		return nil, err
	}

	g, ok := stream.Groups[group]
	if !ok {
		return nil, fmt.Errorf(string(dotpip.ErrMsgNoGroup), strings.Join(key, "."), group)
	}

	var consumers []map[string]any
	for name, cons := range g.Consumers {
		consInfo := make(map[string]any)
		consInfo["name"] = name
		consInfo["pending"] = len(cons.Pending)

		// Find max idle
		var maxIdle int64
		now := time.Now().UnixMilli()
		for id := range cons.Pending {
			if pEntry, pOk := g.Pending[id]; pOk {
				idle := now - pEntry.DeliveryTime
				if idle > maxIdle {
					maxIdle = idle
				}
			}
		}
		consInfo["idle"] = maxIdle

		consumers = append(consumers, consInfo)
	}

	return consumers, nil
}
