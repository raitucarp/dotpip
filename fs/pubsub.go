package fs

import (
	"dotpip"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ensure FileSystem implements PubSub interface.

// PubSubSubscription represents a file system backed pubsub subscription.
// PubSubSubscription represents a file system backed pubsub subscription.
type PubSubSubscription struct {
	fs            *FileSystem
	channels      map[string]struct{}
	patterns      map[string]struct{}
	shardChannels map[string]struct{}
	ch            chan dotpip.PubSubMessage
	closeCh       chan struct{}
	mu            sync.RWMutex
}

// Channel returns a channel to receive messages.
// Channel returns a channel to receive messages.
func (s *PubSubSubscription) Channel() <-chan dotpip.PubSubMessage {
	return s.ch
}

// Unsubscribe unsubscribes from channels.
// Unsubscribe unsubscribes from channels.
func (s *PubSubSubscription) Unsubscribe(channels ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(channels) == 0 {
		s.channels = make(map[string]struct{})
		return nil
	}

	for _, ch := range channels {
		delete(s.channels, ch)
	}
	return nil
}

// PUnsubscribe unsubscribes from patterns.
// PUnsubscribe unsubscribes from patterns.
func (s *PubSubSubscription) PUnsubscribe(patterns ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(patterns) == 0 {
		s.patterns = make(map[string]struct{})
		return nil
	}

	for _, pat := range patterns {
		delete(s.patterns, pat)
	}
	return nil
}

// SUnsubscribe unsubscribes from shard channels.
// SUnsubscribe unsubscribes from shard channels.
func (s *PubSubSubscription) SUnsubscribe(shardChannels ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(shardChannels) == 0 {
		s.shardChannels = make(map[string]struct{})
		return nil
	}

	for _, ch := range shardChannels {
		delete(s.shardChannels, ch)
	}
	return nil
}

// Close closes the subscription.
// Close closes the subscription.
func (s *PubSubSubscription) Close() error {
	select {
	case <-s.closeCh:
	default:
		close(s.closeCh)
	}
	s.fs.removeSubscription(s)
	return nil
}

func matchPattern(pattern, str string) bool {
	// simplified matching for tests, ideally use filepath.Match or custom
	matched, _ := filepath.Match(pattern, str)
	return matched
}

// Subscribe subscribes to channels.
// Subscribe subscribes to channels.
func (f *FileSystem) Subscribe(channels ...string) (dotpip.PubSubSubscription, error) {
	sub := &PubSubSubscription{
		fs:            f,
		channels:      make(map[string]struct{}),
		patterns:      make(map[string]struct{}),
		shardChannels: make(map[string]struct{}),
		ch:            make(chan dotpip.PubSubMessage, 100),
		closeCh:       make(chan struct{}),
	}

	for _, ch := range channels {
		sub.channels[ch] = struct{}{}
	}

	f.addSubscription(sub)
	return sub, nil
}

// PSubscribe subscribes to patterns.
// PSubscribe subscribes to patterns.
func (f *FileSystem) PSubscribe(patterns ...string) (dotpip.PubSubSubscription, error) {
	sub := &PubSubSubscription{
		fs:            f,
		channels:      make(map[string]struct{}),
		patterns:      make(map[string]struct{}),
		shardChannels: make(map[string]struct{}),
		ch:            make(chan dotpip.PubSubMessage, 100),
		closeCh:       make(chan struct{}),
	}

	for _, pat := range patterns {
		sub.patterns[pat] = struct{}{}
	}

	f.addSubscription(sub)
	return sub, nil
}

// SSubscribe subscribes to shard channels.
// SSubscribe subscribes to shard channels.
func (f *FileSystem) SSubscribe(shardChannels ...string) (dotpip.PubSubSubscription, error) {
	sub := &PubSubSubscription{
		fs:            f,
		channels:      make(map[string]struct{}),
		patterns:      make(map[string]struct{}),
		shardChannels: make(map[string]struct{}),
		ch:            make(chan dotpip.PubSubMessage, 100),
		closeCh:       make(chan struct{}),
	}

	for _, ch := range shardChannels {
		sub.shardChannels[ch] = struct{}{}
	}

	f.addSubscription(sub)
	return sub, nil
}

// Publish publishes a message to a channel.
// Publish publishes a message to a channel.
func (f *FileSystem) Publish(channel string, message string) (int, error) {
	// Write to file for fsnotify if it's not a keyspace notification
	// We use a specific directory or pattern for pubsub.
	if !strings.HasPrefix(channel, "__keyspace@") && !strings.HasPrefix(channel, "__keyevent@") {
		pubsubFile := filepath.Join(f.pathRoot, ".pubsub_"+channel)

		// Encode using the formatter
		encodedMsg, err := f.formatter.StringEncode(message)
		var writeMsg string
		if err == nil {
			switch v := encodedMsg.(type) {
			case []byte:
				writeMsg = string(v)
			case string:
				writeMsg = v
			default:
				writeMsg = message
			}
		} else {
			writeMsg = message
		}

		file, err := os.OpenFile(pubsubFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			_, _ = file.WriteString(writeMsg + "\n")
			_ = file.Close()

			// Update offset so we don't double notify
			f.subMutex.Lock()
			stat, _ := os.Stat(pubsubFile)
			if stat != nil {
				f.pubsubOffsets[pubsubFile] = stat.Size()
			}
			f.subMutex.Unlock()
		}
	}

	receivers := f.notifySubscribers(channel, message)

	return receivers, nil
}

// PubSubChannels lists active channels.
// PubSubChannels lists active channels.
func (f *FileSystem) PubSubChannels(pattern string) ([]string, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	channelsMap := make(map[string]struct{})
	for sub := range f.subscriptions {
		sub.mu.RLock()
		for ch := range sub.channels {
			if pattern == "" || matchPattern(pattern, ch) {
				channelsMap[ch] = struct{}{}
			}
		}
		sub.mu.RUnlock()
	}

	var res []string
	for ch := range channelsMap {
		res = append(res, ch)
	}
	return res, nil
}

// PubSubNumPat returns the number of active patterns.
// PubSubNumPat returns the number of active patterns.
func (f *FileSystem) PubSubNumPat() (int, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	patternsMap := make(map[string]struct{})
	for sub := range f.subscriptions {
		sub.mu.RLock()
		for pat := range sub.patterns {
			patternsMap[pat] = struct{}{}
		}
		sub.mu.RUnlock()
	}

	return len(patternsMap), nil
}

// PubSubNumSub returns the number of subscribers for channels.
// PubSubNumSub returns the number of subscribers for channels.
func (f *FileSystem) PubSubNumSub(channels ...string) (map[string]int, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	res := make(map[string]int)
	for _, ch := range channels {
		res[ch] = 0
	}

	for sub := range f.subscriptions {
		sub.mu.RLock()
		for _, ch := range channels {
			if _, ok := sub.channels[ch]; ok {
				res[ch]++
			}
		}
		sub.mu.RUnlock()
	}

	return res, nil
}

// PubSubShardChannels lists shard channels.
// PubSubShardChannels lists shard channels.
func (f *FileSystem) PubSubShardChannels(pattern string) ([]string, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	channelsMap := make(map[string]struct{})
	for sub := range f.subscriptions {
		sub.mu.RLock()
		for ch := range sub.shardChannels {
			if pattern == "" || matchPattern(pattern, ch) {
				channelsMap[ch] = struct{}{}
			}
		}
		sub.mu.RUnlock()
	}

	var res []string
	for ch := range channelsMap {
		res = append(res, ch)
	}
	return res, nil
}

// PubSubShardNumSub gets subscription count for shard channels.
// PubSubShardNumSub gets subscription count for shard channels.
func (f *FileSystem) PubSubShardNumSub(shardChannels ...string) (map[string]int, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	res := make(map[string]int)
	for _, ch := range shardChannels {
		res[ch] = 0
	}

	for sub := range f.subscriptions {
		sub.mu.RLock()
		for _, ch := range shardChannels {
			if _, ok := sub.shardChannels[ch]; ok {
				res[ch]++
			}
		}
		sub.mu.RUnlock()
	}

	return res, nil
}
