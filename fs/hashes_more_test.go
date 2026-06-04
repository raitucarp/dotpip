package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashesHKeysHVals(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_hashes_keysvals_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("myhash")
	_, _ = dotfs.HSet(key, map[string]string{"f1": "v1", "f2": "v2"})

	keys, _ := dotfs.HKeys(key)
	assert.Len(t, keys, 2)

	vals, _ := dotfs.HVals(key)
	assert.Len(t, vals, 2)
}

func TestHashesHStrLenAndHRandField(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_hashes_strlen_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("myhashscan")
	_, _ = dotfs.HSet(key, map[string]string{"f": "val"})

	// HStrLen
	resLen, err := dotfs.HStrLen(key, "f")
	assert.NoError(t, err)
	assert.Equal(t, 3, resLen)

	// HRandField
	resRand, err := dotfs.HRandField(key, 1)
	assert.NoError(t, err)
	assert.Len(t, resRand, 1)

	// HRandField with values
	resRandVal, err := dotfs.HRandField(key, 1, dotpip.WithHRandFieldWithValues())
	assert.NoError(t, err)
	assert.Len(t, resRandVal, 2) // [field, value]
}

func TestSetsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_sets_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	k1 := dotpip.NewKey("set1")
	k2 := dotpip.NewKey("set2")
	k3 := dotpip.NewKey("set3")

	_, _ = dotfs.SAdd(k1, "a", "b", "c")
	_, _ = dotfs.SAdd(k2, "c", "d", "e")

	// SDiff
	resDiff, _ := dotfs.SDiff(k1, k2)
	assert.Len(t, resDiff, 2) // a, b

	// SDiffStore
	resDiffStore, _ := dotfs.SDiffStore(k3, k1, k2)
	assert.Equal(t, 2, resDiffStore)

	// SInter
	resInter, _ := dotfs.SInter(k1, k2)
	assert.Len(t, resInter, 1) // c

	// SInterStore
	resInterStore, _ := dotfs.SInterStore(k3, k1, k2)
	assert.Equal(t, 1, resInterStore)

	// SUnion
	resUnion, _ := dotfs.SUnion(k1, k2)
	assert.Len(t, resUnion, 5)

	// SUnionStore
	resUnionStore, _ := dotfs.SUnionStore(k3, k1, k2)
	assert.Equal(t, 5, resUnionStore)

	// SMove
	resMove, _ := dotfs.SMove(k1, k2, "a")
	assert.Equal(t, true, resMove)

	// SPop
	resPop, _ := dotfs.SPop(k1, 1)
	assert.Len(t, resPop, 1)

	// SRandMember
	resRand, _ := dotfs.SRandMember(k1, 1)
	assert.Len(t, resRand, 1)
}

func TestZSetsMoreCoverage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_zsets_morecov_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("zset3")
	_, _ = dotfs.ZAdd(key, []dotpip.Z{{Score: 10, Member: "x"}, {Score: 20, Member: "y"}, {Score: 30, Member: "z"}})

	// ZRank
	rank, _ := dotfs.ZRank(key, "y")
	assert.Equal(t, 1, rank)

	// ZRevRank
	revrank, _ := dotfs.ZRevRank(key, "y")
	assert.Equal(t, 1, revrank)

	// ZPopMin
	popmin, _ := dotfs.ZPopMin(key, 1)
	assert.Len(t, popmin, 1)
	assert.Equal(t, "x", popmin[0].Member)

	// ZPopMax
	popmax, _ := dotfs.ZPopMax(key, 1)
	assert.Len(t, popmax, 1)
	assert.Equal(t, "z", popmax[0].Member)

	// ZScore
	score, _ := dotfs.ZScore(key, "y")
	assert.Equal(t, 20.0, score)

	// Missing score
	scoreMiss, err := dotfs.ZScore(key, "missing")
	assert.NoError(t, err)
	assert.Equal(t, float64(0), scoreMiss)

	// ZCard
	card, _ := dotfs.ZCard(key)
	assert.Equal(t, 1, card) // only y is left

	// ZRem
	rem, _ := dotfs.ZRem(key, "y")
	assert.Equal(t, 1, rem)

	card2, _ := dotfs.ZCard(key)
	assert.Equal(t, 0, card2)
}

func TestZSetsMoreCoverage2(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_zsets_morecov2_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("zset_incr")
	_, _ = dotfs.ZAdd(key, []dotpip.Z{{Score: 10, Member: "m"}})

	resIncr, err := dotfs.ZIncrBy(key, 5.5, "m")
	assert.NoError(t, err)
	assert.Equal(t, 15.5, resIncr)

	// LexCount
	_, _ = dotfs.ZAdd(key, []dotpip.Z{{Score: 15.5, Member: "a"}, {Score: 15.5, Member: "b"}, {Score: 15.5, Member: "c"}})
	lexCount, err := dotfs.ZLexCount(key, "[a", "[c")
	assert.NoError(t, err)
	assert.Equal(t, 0, lexCount)

	// ZCount
	count, err := dotfs.ZCount(key, 10, 20)
	assert.NoError(t, err)
	assert.Equal(t, 4, count) // a, b, c, m

	// ZRangeWithScores
	resRange, err := dotfs.ZRangeWithScores(key, "0", "-1")
	assert.NoError(t, err)
	assert.Len(t, resRange, 4)

	// ZInterWithScores
	key2 := dotpip.NewKey("zset_inter2")
	_, _ = dotfs.ZAdd(key2, []dotpip.Z{{Score: 10, Member: "a"}, {Score: 20, Member: "m"}})

	resInterScores, err := dotfs.ZInterWithScores(key, key2)
	assert.NoError(t, err)
	assert.Len(t, resInterScores, 2) // a, m

	// ZUnionWithScores
	resUnionScores, err := dotfs.ZUnionWithScores(key, key2)
	assert.NoError(t, err)
	assert.Len(t, resUnionScores, 4)
}

func TestStreamsMoreCoverage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_streams_morecov_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystream_more")
	id1, _ := dotfs.XAdd(key, "*", map[string]string{"k1": "v1"})
	id2, _ := dotfs.XAdd(key, "*", map[string]string{"k2": "v2"})

	// XDel
	resDel, err := dotfs.XDel(key, id1)
	assert.NoError(t, err)
	assert.Equal(t, 1, resDel)

	// XAck
	_, _ = dotfs.XGroupCreate(key, "myg", "0-0", false)
	_, _ = dotfs.XReadGroup("myg", "alice", []dotpip.Key{key}, []string{">"}, dotpip.WithXReadGroupCount(1))

	resAck, err := dotfs.XAck(key, "myg", id2)
	assert.NoError(t, err)
	assert.Equal(t, 1, resAck)
}

func TestListLPos(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_list_lpos_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("lpos")
	_, _ = dotfs.RPush(key, "a", "b", "c", "b")

	// Basic LPos
	pos, err := dotfs.LPos(key, "b")
	assert.NoError(t, err)
	assert.Equal(t, []int{1}, pos)

	// LPos Rank
	posRank, err := dotfs.LPos(key, "b", dotpip.WithLPosRank(2))
	assert.NoError(t, err)
	assert.Equal(t, []int{3}, posRank)

	// LPos Count
	posCount, err := dotfs.LPos(key, "b", dotpip.WithLPosCount(2))
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3}, posCount)
}

func TestStringsMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_strings_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystr")

	// Set NX
	res, _ := dotfs.Set(key, "hello", dotpip.WithNX())
	assert.Equal(t, "hello", res)

	// Set XX
	res2, _ := dotfs.Set(key, "world", dotpip.WithXX())
	assert.Equal(t, "world", res2)

	// Set Get
	resGetSet, _ := dotfs.Set(key, "new", dotpip.WithGet())
	assert.Equal(t, "world", resGetSet)

	// GetDel
	resGetDel, err := dotfs.GetDel(key)
	assert.NoError(t, err)
	assert.Equal(t, "new", resGetDel)

	// MSetNX
	k1 := dotpip.NewKey("k1")
	k2 := dotpip.NewKey("k2")
	resMSetNX, err := dotfs.MSetNX(
		dotpip.KV{Key: k1, Value: "v1"},
		dotpip.KV{Key: k2, Value: "v2"},
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resMSetNX)

	resMSetNX2, err := dotfs.MSetNX(
		dotpip.KV{Key: k1, Value: "v1_new"},
		dotpip.KV{Key: dotpip.NewKey("k3"), Value: "v3"},
	)
	assert.NoError(t, err)
	assert.Equal(t, false, resMSetNX2)

	// IncrByFloat
	knum := dotpip.NewKey("knum")
	_, _ = dotfs.Set(knum, "10.5")
	resIncrF, err := dotfs.IncrByFloat(knum, 2.5)
	assert.NoError(t, err)
	assert.Equal(t, 13.0, resIncrF)

	// GetRange
	_, _ = dotfs.Set(key, "helloworld")
	resGetRange, err := dotfs.GetRange(key, 0, 4)
	assert.NoError(t, err)
	assert.Equal(t, "hello", resGetRange)

	// SetRange
	resSetRange, err := dotfs.SetRange(key, 5, "there")
	assert.NoError(t, err)
	assert.Equal(t, 10, resSetRange)
}

func TestStringsMoreCoverage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_strings_more2_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystr2")

	// Append
	_, _ = dotfs.Set(key, "hello")
	resApp := dotfs.Append(key, " world")
	assert.Equal(t, 11, resApp)

	// SetRange out of bounds
	_, _ = dotfs.Set(key, "hello")
	resSetR, _ := dotfs.SetRange(key, 10, "world") // expands to length 15
	assert.Equal(t, 15, resSetR)

	// Digest
	_, _ = dotfs.Set(dotpip.NewKey("k1"), "hello")
	digest, err := dotfs.Digest(dotpip.NewKey("k1"))
	assert.NoError(t, err)
	assert.NotEmpty(t, digest)

	// StrLen
	strLen := dotfs.StrLen(key)
	assert.NoError(t, err)
	assert.Equal(t, 15, strLen)

	// MGet
	k1 := dotpip.NewKey("k1")
	k2 := dotpip.NewKey("k2")
	k3 := dotpip.NewKey("k3")
	_, _ = dotfs.Set(k1, "1")
	_, _ = dotfs.Set(k2, "2")
	resMGet, err := dotfs.MGet(k1, k2, k3)
	assert.NoError(t, err)
	assert.Len(t, resMGet, 3)
	assert.Equal(t, "1", resMGet[0])
	assert.Equal(t, "2", resMGet[1])
	assert.Equal(t, "", resMGet[2])
}

func TestStringsMoreCoverage3(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_strings_more3_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystr3")

	// Incr
	_, _ = dotfs.Set(key, "10")
	resIncr, err := dotfs.Incr(key)
	assert.NoError(t, err)
	assert.Equal(t, 11, resIncr)

	// Decr
	resDecr, err := dotfs.Decr(key)
	assert.NoError(t, err)
	assert.Equal(t, 10, resDecr)

	// IncrBy
	resIncrBy, err := dotfs.IncrBy(key, 5)
	assert.NoError(t, err)
	assert.Equal(t, 15, resIncrBy)

	// DecrBy
	resDecrBy, err := dotfs.DecrBy(key, 2)
	assert.NoError(t, err)
	assert.Equal(t, 13, resDecrBy)

	// MSet
	k1 := dotpip.NewKey("k1_mset")
	k2 := dotpip.NewKey("k2_mset")
	err = dotfs.MSet(dotpip.KV{Key: k1, Value: "v1"}, dotpip.KV{Key: k2, Value: "v2"})
	assert.NoError(t, err)

}

func TestStringsMoreCoverage4(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_strings_more4_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mystr4")

	// Set with Ex options
	res1, _ := dotfs.Set(key, "10", dotpip.WithEx(10))
	assert.Equal(t, "10", res1)

	res2, _ := dotfs.Set(key, "10", dotpip.WithPx(10000))
	assert.Equal(t, "10", res2)

	res3, _ := dotfs.Set(key, "10", dotpip.WithExAt(1999999999))
	assert.Equal(t, "10", res3)

	res4, _ := dotfs.Set(key, "10", dotpip.WithPxAt(1999999999000))
	assert.Equal(t, "10", res4)

	res5, _ := dotfs.Set(key, "10", dotpip.WithKeepTTL())
	assert.Equal(t, "10", res5)
}
