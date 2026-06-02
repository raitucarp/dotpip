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
	return Key(s)
}

func NewKeyWithDelimiter(keyString string, delimiter string) Key {
	return Key(strings.Split(keyString, delimiter))
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

	Formatter(fmap DataTypeFormatter)
}

func New(dotface DotPip) DotPip {
	return dotface
}
