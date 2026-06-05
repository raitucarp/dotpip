package fs

import (
	"bufio"
	"dotpip"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileSystem struct {
	pathRoot   string
	formatter  *dotpip.DataTypeFormatter
	encodeType FileEncodeType

	expirations map[string]int64
	expMutex    sync.RWMutex
	expStop     chan struct{}

	subscriptions map[*PubSubSubscription]struct{}
	subMutex      sync.RWMutex
	watcher       *fsnotify.Watcher

	pubsubOffsets map[string]int64
	config        map[string]string
}

func NewFileSystem(pathRoot string) *FileSystem {
	watcher, _ := fsnotify.NewWatcher()
	if watcher != nil {
		_ = watcher.Add(pathRoot)
	}

	f := FileSystem{
		pathRoot:      pathRoot,
		formatter:     &dotpip.DataTypeFormatter{},
		encodeType:    JSON,
		expirations:   make(map[string]int64),
		expStop:       make(chan struct{}),
		subscriptions: make(map[*PubSubSubscription]struct{}),
		watcher:       watcher,
		pubsubOffsets: make(map[string]int64),
		config:        make(map[string]string),
	}

	f.formatter.StringEncode = f.stringEncode
	f.formatter.StringDecode = f.stringDecode
	f.formatter.HashEncode = f.hashEncode
	f.formatter.HashDecode = f.hashDecode
	f.formatter.ListEncode = f.listEncode
	f.formatter.ListDecode = f.listDecode
	f.formatter.SetEncode = f.setEncode
	f.formatter.SetDecode = f.setDecode
	f.formatter.SortedSetEncode = f.sortedSetEncode
	f.formatter.SortedSetDecode = f.sortedSetDecode
	f.formatter.BitmapEncode = f.bitmapEncode
	f.formatter.BitmapDecode = f.bitmapDecode
	f.formatter.GeospatialEncode = f.geospatialEncode
	f.formatter.GeospatialDecode = f.geospatialDecode
	f.formatter.HyperLogLogEncode = f.hyperLogLogEncode
	f.formatter.HyperLogLogDecode = f.hyperLogLogDecode
	f.formatter.StreamEncode = f.streamEncode
	f.formatter.StreamDecode = f.streamDecode
	f.formatter.ArrayEncode = f.arrayEncode
	f.formatter.ArrayDecode = f.arrayDecode
	f.formatter.VectorSetEncode = f.vectorSetEncode
	f.formatter.VectorSetDecode = f.vectorSetDecode

	f.loadExpirations()
	go f.scanExpirations()
	if f.watcher != nil {
		go f.watchFS()
	}

	return &f
}

func (f *FileSystem) Close() {
	if f.expStop != nil {
		close(f.expStop)
	}
	if f.watcher != nil {
		_ = f.watcher.Close()
	}
}

func (f *FileSystem) addSubscription(sub *PubSubSubscription) {
	f.subMutex.Lock()
	f.subscriptions[sub] = struct{}{}
	f.subMutex.Unlock()
}

func (f *FileSystem) removeSubscription(sub *PubSubSubscription) {
	f.subMutex.Lock()
	delete(f.subscriptions, sub)
	f.subMutex.Unlock()
}

func (f *FileSystem) notifySubscribers(channel string, message string) int {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	count := 0
	for sub := range f.subscriptions {
		sub.mu.RLock()
		matched := false
		if _, ok := sub.channels[channel]; ok {
			matched = true
		}
		if !matched {
			for pat := range sub.patterns {
				if matchPattern(pat, channel) {
					matched = true
					break
				}
			}
		}
		if !matched {
			if _, ok := sub.shardChannels[channel]; ok {
				matched = true
			}
		}
		sub.mu.RUnlock()

		if matched {
			select {
			case sub.ch <- dotpip.PubSubMessage{Channel: channel, Payload: message}:
				count++
			default:
				// channel full, drop message
			}
		}
	}
	return count
}

func (f *FileSystem) loadExpirations() {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()

	_ = filepath.WalkDir(f.pathRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(path, ".ex") {
			content, err := f.readExByPath(path)
			if err == nil {
				expireAt, err := strconv.ParseInt(string(content), 10, 64)
				if err == nil {
					// We use the absolute .ex file path or the relative key as map key.
					// Let's use the key path relative to pathRoot.
					// For simplicity, we can use the absolute path of the data file as the map key.
					dataPath := strings.TrimSuffix(path, ".ex")
					f.expirations[dataPath] = expireAt
				}
			}
		}
		return nil
	})
}

func (f *FileSystem) scanExpirations() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-f.expStop:
			return
		case <-ticker.C:
			f.processExpirations()
		}
	}
}

func (f *FileSystem) processExpirations() {
	f.expMutex.Lock()
	defer f.expMutex.Unlock()

	now := time.Now().UnixMilli()
	for dataPath, expireAt := range f.expirations {
		if now >= expireAt {
			// Expired!
			keyStr := dataPath[len(f.pathRoot):]
			if keyStr != "" && keyStr[0] == filepath.Separator {
				keyStr = keyStr[1:]
			}
			ext := filepath.Ext(keyStr)
			if ext == ".json" || ext == ".yaml" || ext == ".toml" {
				keyStr = keyStr[:len(keyStr)-len(ext)]
			}
			keyParts := strings.Split(keyStr, string(filepath.Separator))
			f.emitKeyspaceEvent(keyParts, "expired", 'x')
			_ = f.removeFileByPath(dataPath)
			_ = f.removeExByPath(dataPath + ".ex")
			delete(f.expirations, dataPath)
		}
	}
}

func (f *FileSystem) watchFS() {
	for {
		select {
		case <-f.expStop:
			return
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}
			f.handleFSEvent(event)
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			_ = err
		}
	}
}

func (f *FileSystem) handleFSEvent(event fsnotify.Event) {
	// Emit keyspace and keyevent notifications
	// e.g. for SET we might get a CREATE or WRITE event.
	baseName := filepath.Base(event.Name)

	// Ignore internal files
	if strings.HasPrefix(baseName, ".") || strings.HasSuffix(baseName, ".ex") {
		if strings.HasPrefix(baseName, ".pubsub_") && (event.Op&fsnotify.Write == fsnotify.Write) {
			channel := strings.TrimPrefix(baseName, ".pubsub_")
			lines, _ := f.readTail(event.Name)
			for _, line := range lines {
				decodedMsg, err := f.formatter.StringDecode([]byte(line))
				if err == nil {
					f.notifySubscribers(channel, decodedMsg)
				} else {
					f.notifySubscribers(channel, line)
				}
			}
		}
		return
	}

}

func (f *FileSystem) readTail(path string) ([]string, error) {
	f.subMutex.Lock()
	defer f.subMutex.Unlock()

	if f.pubsubOffsets == nil {
		f.pubsubOffsets = make(map[string]int64)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	offset := f.pubsubOffsets[path]
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	var newLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newLines = append(newLines, scanner.Text())
	}

	newOffset, err := file.Seek(0, 1) // current position
	if err == nil {
		f.pubsubOffsets[path] = newOffset
	}

	return newLines, scanner.Err()
}

func (f *FileSystem) ConfigSet(parameter string, value string) error {
	f.subMutex.Lock()
	defer f.subMutex.Unlock()
	f.config[parameter] = value
	return nil
}

func (f *FileSystem) ConfigGet(parameter string) (map[string]string, error) {
	f.subMutex.RLock()
	defer f.subMutex.RUnlock()

	res := make(map[string]string)
	if parameter == "*" {
		for k, v := range f.config {
			res[k] = v
		}
	} else if val, ok := f.config[parameter]; ok {
		res[parameter] = val
	}

	return res, nil
}

func (f *FileSystem) emitKeyspaceEvent(key []string, event string, typeChar rune) {
	if key == nil {
		return
	}
	f.subMutex.RLock()
	notifyEvents := f.config["notify-keyspace-events"]
	f.subMutex.RUnlock()

	if notifyEvents == "" {
		return
	}

	keyspaceEnabled := strings.Contains(notifyEvents, "K")
	keyeventEnabled := strings.Contains(notifyEvents, "E")

	if !keyspaceEnabled && !keyeventEnabled {
		return
	}

	matchAll := strings.Contains(notifyEvents, "A")
	isExcludedFromA := typeChar == 'm' || typeChar == 'n' || typeChar == 'o' || typeChar == 'c'

	if !(matchAll && !isExcludedFromA) && !strings.ContainsRune(notifyEvents, typeChar) {
		return // Not enabled
	}

	keyString := strings.Join(key, string(filepath.Separator))

	if keyspaceEnabled {
		f.notifySubscribers("__keyspace@0__:"+keyString, event)
	}
	if keyeventEnabled {
		f.notifySubscribers("__keyevent@0__:"+event, keyString)
	}
}

func (f *FileSystem) emitSubkeyEvent(key []string, event string, typeChar rune, subkeys []string) {
	if key == nil || len(subkeys) == 0 {
		return
	}

	f.subMutex.RLock()
	notifyEvents := f.config["notify-keyspace-events"]
	f.subMutex.RUnlock()

	if notifyEvents == "" {
		return
	}

	subkeyspaceEnabled := strings.Contains(notifyEvents, "S")
	subkeyeventEnabled := strings.Contains(notifyEvents, "T")
	subkeyspaceitemEnabled := strings.Contains(notifyEvents, "I")
	subkeyspaceeventEnabled := strings.Contains(notifyEvents, "V")

	if !subkeyspaceEnabled && !subkeyeventEnabled && !subkeyspaceitemEnabled && !subkeyspaceeventEnabled {
		return
	}

	// Check if the typeChar is enabled or if all ('A') are enabled
	matchAll := strings.Contains(notifyEvents, "A")
	if !matchAll && !strings.ContainsRune(notifyEvents, typeChar) {
		return // Not enabled
	}

	keyString := strings.Join(key, string(filepath.Separator))

	var subkeyParts []string
	for _, sk := range subkeys {
		subkeyParts = append(subkeyParts, strconv.Itoa(len(sk))+":"+sk)
	}
	subkeysJoined := strings.Join(subkeyParts, ",")

	if subkeyspaceEnabled {
		if !strings.Contains(event, "|") {
			payload := event + "|" + subkeysJoined
			f.notifySubscribers("__subkeyspace@0__:"+keyString, payload)
		}
	}

	if subkeyeventEnabled {
		payload := strconv.Itoa(len(keyString)) + ":" + keyString + "|" + subkeysJoined
		f.notifySubscribers("__subkeyevent@0__:"+event, payload)
	}

	if subkeyspaceitemEnabled {
		if !strings.Contains(keyString, "\n") {
			for _, sk := range subkeys {
				channel := "__subkeyspaceitem@0__:" + keyString + "\n" + sk
				f.notifySubscribers(channel, event)
			}
		}
	}

	if subkeyspaceeventEnabled {
		if !strings.Contains(event, "|") {
			channel := "__subkeyspaceevent@0__:" + event + "|" + keyString
			f.notifySubscribers(channel, subkeysJoined)
		}
	}
}
