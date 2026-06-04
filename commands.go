package dotpip

import "strings"

type Key []string

type BitOp string

const (
	BitOpAnd BitOp = "AND"
	BitOpOr  BitOp = "OR"
	BitOpXor BitOp = "XOR"
	BitOpNot BitOp = "NOT"
)

func NewKey(s ...string) Key {
	var k Key
	for _, part := range s {
		k = append(k, strings.Split(part, ":")...)
	}
	return k
}

func NewKeyWithDelimiter(args ...string) Key {
	if len(args) == 0 {
		return Key{}
	}
	if len(args) == 1 {
		return Key{}
	}
	delimiter := args[len(args)-1]
	var k Key
	for i := 0; i < len(args)-1; i++ {
		k = append(k, strings.Split(args[i], delimiter)...)
	}
	return k
}

type DataTypeFormatter struct {
	StringEncode func(value string) (any, error)
	StringDecode func(value any) (string, error)

	HashEncode func(value map[string]string) (any, error)
	HashDecode func(value any) (map[string]string, error)

	ListEncode func(value []any) (any, error)
	ListDecode func(value any) ([]any, error)

	SetEncode func(value map[string]any) (any, error)
	SetDecode func(value any) (map[string]any, error)

	SortedSetEncode func(value map[string]float64) (any, error)
	SortedSetDecode func(value any) (map[string]float64, error)

	VectorSetDecode func(value any) ([]float64, error)
	VectorSetEncode func(value []float64) (any, error)

	StreamEncode func(value Stream) (any, error)
	StreamDecode func(value any) (Stream, error)

	BitmapEncode func(value []uint) (any, error)
	BitmapDecode func(value any) ([]uint, error)

	BitfieldEncode func(value any) ([]any, error)
	BitfieldDecode func(value []any) (any, error)

	ArrayEncode func(value []string) (any, error)
	ArrayDecode func(value any) ([]string, error)

	GeospatialEncode func(value map[string]GeoLocation) (any, error)
	GeospatialDecode func(value any) (map[string]GeoLocation, error)

	JSONEncode func(value any) (any, error)
	JSONDecode func(value any) (any, error)

	HyperLogLogEncode func(value []byte) (any, error)
	HyperLogLogDecode func(value any) ([]byte, error)
}

type KV struct {
	Key   Key
	Value string
}

type DotPip interface {
	Append(key Key, value string) int
	Get(key Key) (string, error)
	Set(key Key, value string, options ...SetOption) (string, error)
	Digest(key Key) (string, error)

	StrLen(key Key) int
	Incr(key Key) (int, error)
	IncrBy(key Key, increment int) (int, error)
	IncrByFloat(key Key, increment float64) (float64, error)
	Decr(key Key) (int, error)
	DecrBy(key Key, decrement int) (int, error)
	GetDel(key Key) (string, error)
	GetRange(key Key, start int, end int) (string, error)
	SetRange(key Key, offset int, value string) (int, error)
	MGet(keys ...Key) ([]string, error)
	MSet(kvs ...KV) error
	MSetNX(kvs ...KV) (bool, error)

	Del(keys ...Key) int
	Copy(source Key, destination Key, options ...CopyOption) int
	Exists(keys ...Key) ([]bool, error)
	FlushAll() error
	Rename(key Key, newKey Key) error
	RenameNX(key Key, newKey Key) (bool, error)
	Keys(pattern string) ([]Key, error)
	Type(key Key) (string, error)
	RandomKey() (Key, error)
	Touch(keys ...Key) (int, error)
	Unlink(keys ...Key) int
	Dump(key Key) ([]byte, error)
	Restore(key Key, ttl int, serializedValue []byte, options ...RestoreOption) error
	Sort(key Key) ([]string, error)
	Scan(cursor uint64, options ...ScanOption) (uint64, []Key, error)
	Move(key Key, db DotPip) (int, error)
	Migrate(host string, port int, key Key, destinationDB DotPip, timeout int, options ...MigrateOption) error
	Wait(numReplicas int, timeout int) (int, error)
	WaitAOF(numLocal int, numReplicas int, timeout int) (int, int, error)
	DBSize() (int, error)
	ObjectEncoding(key Key) (string, error)
	ObjectFreq(key Key) (int, error)
	ObjectIdletime(key Key) (int, error)
	ObjectRefcount(key Key) (int, error)

	Expire(key Key, seconds int, options ...ExpireOption) (bool, error)
	ExpireAt(key Key, timestamp int, options ...ExpireOption) (bool, error)
	ExpireTime(key Key) (int64, error)
	PExpire(key Key, milliseconds int, options ...ExpireOption) (bool, error)
	PExpireAt(key Key, timestamp int, options ...ExpireOption) (bool, error)
	PExpireTime(key Key) (int64, error)
	TTL(key Key) (int64, error)
	PTTL(key Key) (int64, error)
	Persist(key Key) (bool, error)

	LIndex(key Key, index int) (string, error)
	LInsert(key Key, option LInsertOption, pivot string, element string) (int, error)
	LLen(key Key) (int, error)
	LMove(source Key, destination Key, srcDir LMoveDir, destDir LMoveDir) (string, error)
	LPop(key Key, count int) ([]string, error)
	LPos(key Key, element string, options ...LPosOption) ([]int, error)
	LPush(key Key, elements ...string) (int, error)
	LPushX(key Key, elements ...string) (int, error)
	LRange(key Key, start int, stop int) ([]string, error)
	LRem(key Key, count int, element string) (int, error)
	LSet(key Key, index int, element string) error
	LTrim(key Key, start int, stop int) error
	RPop(key Key, count int) ([]string, error)
	RPush(key Key, elements ...string) (int, error)
	RPushX(key Key, elements ...string) (int, error)

	HDel(key Key, fields ...string) (int, error)
	HExists(key Key, field string) (bool, error)
	HGet(key Key, field string) (string, error)
	HGetAll(key Key) (map[string]string, error)
	HIncrBy(key Key, field string, increment int) (int, error)
	HIncrByFloat(key Key, field string, increment float64) (float64, error)
	HKeys(key Key) ([]string, error)
	HLen(key Key) (int, error)
	HMGet(key Key, fields ...string) ([]string, error)
	HRandField(key Key, count int, options ...HRandFieldOption) ([]string, error)
	HSet(key Key, values map[string]string) (int, error)
	HSetNX(key Key, field string, value string) (bool, error)
	HStrLen(key Key, field string) (int, error)
	HVals(key Key) ([]string, error)
	HScan(key Key, cursor uint64, options ...ScanOption) (uint64, map[string]string, error)

	SAdd(key Key, members ...string) (int, error)
	SCard(key Key) (int, error)
	SDiff(keys ...Key) ([]string, error)
	SDiffStore(destination Key, keys ...Key) (int, error)
	SInter(keys ...Key) ([]string, error)
	SInterCard(limit int, keys ...Key) (int, error)
	SInterStore(destination Key, keys ...Key) (int, error)
	SIsMember(key Key, member string) (bool, error)
	SMembers(key Key) ([]string, error)
	SMIsMember(key Key, members ...string) ([]bool, error)
	SMove(source Key, destination Key, member string) (bool, error)
	SPop(key Key, count int) ([]string, error)
	SRandMember(key Key, count int) ([]string, error)
	SRem(key Key, members ...string) (int, error)
	SUnion(keys ...Key) ([]string, error)
	SUnionStore(destination Key, keys ...Key) (int, error)
	SScan(key Key, cursor uint64, options ...ScanOption) (uint64, []string, error)

	ZAdd(key Key, members []Z, options ...ZAddOption) (int, error)
	ZCard(key Key) (int, error)
	ZCount(key Key, minVal float64, maxVal float64) (int, error)
	ZDiff(keys ...Key) ([]string, error)
	ZDiffWithScores(keys ...Key) ([]Z, error)
	ZIncrBy(key Key, increment float64, member string) (float64, error)
	ZInter(keys ...Key) ([]string, error)
	ZInterWithScores(keys ...Key) ([]Z, error)
	ZLexCount(key Key, minVal string, maxVal string) (int, error)
	ZPopMax(key Key, count int) ([]Z, error)
	ZPopMin(key Key, count int) ([]Z, error)
	ZRandMember(key Key, count int) ([]string, error)
	ZRandMemberWithScores(key Key, count int) ([]Z, error)
	ZRange(key Key, start string, stop string, options ...ZRangeOption) ([]string, error)
	ZRangeWithScores(key Key, start string, stop string, options ...ZRangeOption) ([]Z, error)
	ZRank(key Key, member string) (int, error)
	ZRem(key Key, members ...string) (int, error)
	ZRevRank(key Key, member string) (int, error)
	ZScore(key Key, member string) (float64, error)
	ZUnion(keys ...Key) ([]string, error)
	ZUnionWithScores(keys ...Key) ([]Z, error)
	ZScan(key Key, cursor uint64, options ...ScanOption) (uint64, []Z, error)

	BitCount(key Key, start int, end int) (int, error)
	BitField(key Key, args ...any) ([]any, error)
	BitOp(operation BitOp, destKey Key, keys ...Key) (int, error)
	BitPos(key Key, bit int, start int, end int) (int, error)
	GetBit(key Key, offset int) (int, error)
	SetBit(key Key, offset int, value int) (int, error)

	GeoAdd(key Key, members []GeoLocation, options ...GeoAddOption) (int, error)
	GeoDist(key Key, member1 string, member2 string, unit GeoUnit) (float64, error)
	GeoHash(key Key, members ...string) ([]string, error)
	GeoPos(key Key, members ...string) ([]*GeoLocation, error)
	GeoSearch(key Key, options ...GeoSearchOption) ([]GeoSearchResult, error)
	GeoSearchStore(destination Key, source Key, searchOptions []GeoSearchOption, storeOptions ...GeoSearchStoreOption) (int, error)

	PFAdd(key Key, elements ...string) (int, error)
	PFCount(keys ...Key) (int, error)
	PFMerge(destKey Key, sourceKeys ...Key) error

	XAck(key Key, group string, ids ...string) (int, error)
	XAdd(key Key, id string, values map[string]string, options ...XAddOption) (string, error)
	XDel(key Key, ids ...string) (int, error)
	XGroupCreate(key Key, group string, id string, mkStream bool) (string, error)
	XGroupCreateConsumer(key Key, group string, consumer string) (int, error)
	XGroupDelConsumer(key Key, group string, consumer string) (int, error)
	XGroupDestroy(key Key, group string) (int, error)
	XGroupSetID(key Key, group string, id string) (string, error)
	XLen(key Key) (int, error)
	XRange(key Key, start string, end string, count int) ([]StreamEntry, error)
	XRevRange(key Key, end string, start string, count int) ([]StreamEntry, error)
	XRead(keys []Key, ids []string, options ...XReadOption) (map[string][]StreamEntry, error)
	XReadGroup(group string, consumer string, keys []Key, ids []string, options ...XReadGroupOption) (map[string][]StreamEntry, error)
	XTrim(key Key, options ...XTrimOption) (int, error)
	XPending(key Key, group string, options ...XPendingOption) ([]any, error)
	XClaim(key Key, group string, consumer string, minIdleTime int, ids []string, options ...XClaimOption) ([]StreamEntry, error)
	XAutoClaim(key Key, group string, consumer string, minIdleTime int, start string, options ...XAutoClaimOption) (string, []StreamEntry, error)
	XInfoStream(key Key) (map[string]any, error)
	XInfoGroups(key Key) ([]map[string]any, error)
	XInfoConsumers(key Key, group string) ([]map[string]any, error)

	Publish(channel string, message string) (int, error)
	Subscribe(channels ...string) (PubSubSubscription, error)
	PSubscribe(patterns ...string) (PubSubSubscription, error)
	SSubscribe(shardChannels ...string) (PubSubSubscription, error)
	PubSubChannels(pattern string) ([]string, error)
	PubSubNumPat() (int, error)
	PubSubNumSub(channels ...string) (map[string]int, error)
	PubSubShardChannels(pattern string) ([]string, error)
	PubSubShardNumSub(shardChannels ...string) (map[string]int, error)

	ConfigGet(parameter string) (map[string]string, error)
	ConfigSet(parameter string, value string) error

	JSONArrAppend(key Key, path string, values ...any) ([]any, error)
	JSONArrIndex(key Key, path string, value any, startAndStop ...int) ([]any, error)
	JSONArrInsert(key Key, path string, index int, values ...any) ([]any, error)
	JSONArrLen(key Key, path string) ([]any, error)
	JSONArrPop(key Key, path string, index ...int) ([]any, error)
	JSONArrTrim(key Key, path string, start int, stop int) ([]any, error)
	JSONClear(key Key, path string) (int, error)
	JSONDebug(subcommand string, key Key, path string) (any, error)
	JSONDel(key Key, path string) (int, error)
	JSONForget(key Key, path string) (int, error)
	JSONGet(key Key, paths ...string) (any, error)
	JSONMerge(key Key, path string, value any) (string, error)
	JSONMGet(path string, keys ...Key) ([]any, error)
	JSONMSet(args ...JSONMSetArg) (string, error)
	JSONNumIncrBy(key Key, path string, value float64) ([]any, error)
	JSONNumMultBy(key Key, path string, value float64) ([]any, error)
	JSONObjKeys(key Key, path string) ([]any, error)
	JSONObjLen(key Key, path string) ([]any, error)
	JSONResp(key Key, path string) (any, error)
	JSONSet(key Key, path string, value any) (string, error)
	JSONStrAppend(key Key, path string, value string) ([]any, error)
	JSONStrLen(key Key, path string) ([]any, error)
	JSONToggle(key Key, path string) ([]any, error)
	JSONType(key Key, path string) ([]any, error)

	YAMLArrAppend(key Key, path string, values ...any) ([]any, error)
	YAMLArrIndex(key Key, path string, value any, startAndStop ...int) ([]any, error)
	YAMLArrInsert(key Key, path string, index int, values ...any) ([]any, error)
	YAMLArrLen(key Key, path string) ([]any, error)
	YAMLArrPop(key Key, path string, index ...int) ([]any, error)
	YAMLArrTrim(key Key, path string, start int, stop int) ([]any, error)
	YAMLClear(key Key, path string) (int, error)
	YAMLDebug(subcommand string, key Key, path string) (any, error)
	YAMLDel(key Key, path string) (int, error)
	YAMLForget(key Key, path string) (int, error)
	YAMLGet(key Key, paths ...string) (any, error)
	YAMLMerge(key Key, path string, value any) (string, error)
	YAMLMGet(path string, keys ...Key) ([]any, error)
	YAMLMSet(args ...JSONMSetArg) (string, error)
	YAMLNumIncrBy(key Key, path string, value float64) ([]any, error)
	YAMLNumMultBy(key Key, path string, value float64) ([]any, error)
	YAMLObjKeys(key Key, path string) ([]any, error)
	YAMLObjLen(key Key, path string) ([]any, error)
	YAMLResp(key Key, path string) (any, error)
	YAMLSet(key Key, path string, value any) (string, error)
	YAMLStrAppend(key Key, path string, value string) ([]any, error)
	YAMLStrLen(key Key, path string) ([]any, error)
	YAMLToggle(key Key, path string) ([]any, error)
	YAMLType(key Key, path string) ([]any, error)

	TOMLArrAppend(key Key, path string, values ...any) ([]any, error)
	TOMLArrIndex(key Key, path string, value any, startAndStop ...int) ([]any, error)
	TOMLArrInsert(key Key, path string, index int, values ...any) ([]any, error)
	TOMLArrLen(key Key, path string) ([]any, error)
	TOMLArrPop(key Key, path string, index ...int) ([]any, error)
	TOMLArrTrim(key Key, path string, start int, stop int) ([]any, error)
	TOMLClear(key Key, path string) (int, error)
	TOMLDebug(subcommand string, key Key, path string) (any, error)
	TOMLDel(key Key, path string) (int, error)
	TOMLForget(key Key, path string) (int, error)
	TOMLGet(key Key, paths ...string) (any, error)
	TOMLMerge(key Key, path string, value any) (string, error)
	TOMLMGet(path string, keys ...Key) ([]any, error)
	TOMLMSet(args ...JSONMSetArg) (string, error)
	TOMLNumIncrBy(key Key, path string, value float64) ([]any, error)
	TOMLNumMultBy(key Key, path string, value float64) ([]any, error)
	TOMLObjKeys(key Key, path string) ([]any, error)
	TOMLObjLen(key Key, path string) ([]any, error)
	TOMLResp(key Key, path string) (any, error)
	TOMLSet(key Key, path string, value any) (string, error)
	TOMLStrAppend(key Key, path string, value string) ([]any, error)
	TOMLStrLen(key Key, path string) ([]any, error)
	TOMLToggle(key Key, path string) ([]any, error)
	TOMLType(key Key, path string) ([]any, error)

	ARCount(key Key) (int, error)
	ARDel(key Key, indices ...int) (int, error)
	ARDelRange(key Key, ranges ...[2]int) (int, error)
	ARGet(key Key, index int) (string, error)
	ARGetRange(key Key, start, end int) ([]string, error)
	ARGrep(key Key, start, end string, predicates []ARGrepPredicate, options ARGrepOptions) ([]any, error)
	ARInfo(key Key, full bool) (map[string]any, error)
	ARInsert(key Key, values ...string) (int, error)
	ARLastItems(key Key, count int) ([]string, error)
	ARLen(key Key) (int, error)
	ARMGet(key Key, indices ...int) ([]string, error)
	ARMSet(key Key, indexValues []ARIndexValue) (int, error)
	ARNext(key Key) (int, error)
	AROp(key Key, start, end int, operation string, matchValue *string) (any, error)
	ARRing(key Key, size int, values ...string) (int, error)
	ARScan(key Key, start, end int, limit *int) ([]any, error)
	ARSeek(key Key, index int) (int, error)
	ARSet(key Key, index int, values ...string) (int, error)

	Formatter(fmap DataTypeFormatter)
}

type ARIndexValue struct {
	Index int
	Value string
}

type ARGrepPredicate struct {
	Type  string
	Value string
}

type ARGrepOptions struct {
	And        bool
	Or         bool
	Limit      *int
	WithValues bool
	NoCase     bool
}

type JSONMSetArg struct {
	Key   Key
	Path  string
	Value any
}

type PubSubMessage struct {
	Type    string
	Pattern string
	Channel string
	Payload string
}

type PubSubSubscription interface {
	Channel() <-chan PubSubMessage
	Unsubscribe(channels ...string) error
	PUnsubscribe(patterns ...string) error
	SUnsubscribe(shardChannels ...string) error
	Close() error
}

func New(dotface DotPip) DotPip {
	return dotface
}

func (f *DataTypeFormatter) JSONSetEncode(value any) (any, error) {
	if f.JSONEncode != nil {
		return f.JSONEncode(value)
	}
	return value, nil
}
