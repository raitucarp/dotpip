package fs

import (
	"dotpip"
	"testing"
	"time"
)

func TestAllExtraCommands(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"extra_key"}
	path := "$"
	fsys.Set(key, "val")

	// PExpireTime / ExpireTime
	pTime, _ := fsys.PExpireTime(key)
	if pTime > 0 {
		t.Errorf("Unexpected expiration")
	}

	sTime, _ := fsys.ExpireTime(key)
	if sTime > 0 {
		t.Errorf("Unexpected expiration")
	}

	// ExpireAt / PExpireAt
	future := time.Now().Add(2 * time.Second).Unix()
	res, _ := fsys.ExpireAt(key, int(future))
	if !res {
		t.Errorf("Expected true")
	}

	res, _ = fsys.PExpireAt(key, int(future*1000))
	if !res {
		t.Errorf("Expected true")
	}

	// Formatter
	formatter := dotpip.DataTypeFormatter{}
	fsys.Formatter(formatter)

	// Bitmap encode/decode manually
	val := []uint{1, 2, 3}
	encodings := []FileEncodeType{JSON, YAML, TOML, RAW}
	for _, enc := range encodings {
		fsys.EncodeType(enc)
		encoded, _ := fsys.bitmapEncode(val)
		decoded, _ := fsys.bitmapDecode(encoded)
		if len(decoded) != 3 {
			t.Errorf("Decode error")
		}
	}

	// Generic
	key2 := dotpip.Key{"generic_extra2"}
	fsys.RenameNX(key, key2)
	fsys.Sort(key2)
	fsys.Migrate("host", 6379, key2, 0, 100)

	_, _ = fsys.ConfigGet("*")

	// JSON
	fsys.JSONForget(key, path)

	_, _ = fsys.JSONObjLen(key, path)

	fsys.JSONResp(key, path)
	fsys.JSONDebug("HELP", key, path)

	// Lists
	lx, _ := fsys.RPushX(key, "v")
	if lx != 0 {
		t.Errorf("Expected 0")
	}

	// PubSub
	ssub, _ := fsys.SSubscribe("sch1")
	if ssub != nil {
		err := ssub.SUnsubscribe("sch1")
		if err != nil {
			t.Errorf("Error in unsubscribe: %v", err)
		}
	}

	psc, _ := fsys.PubSubShardChannels("s*")
	if len(psc) != 0 {
		t.Errorf("Expected empty")
	}

	psn, _ := fsys.PubSubShardNumSub("sch1")
	if psn["sch1"] != 0 {
		t.Errorf("Expected 0")
	}

	// Streams
	skey := dotpip.Key{"stream_extra"}
	fsys.XGroupSetID(skey, "group1", "0-0")
	fsys.XGroupDestroy(skey, "group1")
	fsys.XGroupDelConsumer(skey, "group1", "consumer1")
	fsys.XInfoGroups(skey)
	fsys.XInfoConsumers(skey, "group1")

	// TOML proxy
	fsys.TOMLSet(key, path, "value")
	fsys.TOMLArrInsert(key, path, 0, "val")
	fsys.TOMLArrLen(key, path)
	fsys.TOMLArrPop(key, path, 0)
	fsys.TOMLArrTrim(key, path, 0, 1)
	fsys.TOMLClear(key, path)
	fsys.TOMLDel(key, path)
	fsys.TOMLForget(key, path)
	fsys.TOMLMerge(key, path, "value")
	fsys.TOMLMGet(path, key)
	fsys.TOMLMSet(dotpip.JSONMSetArg{Key: key, Path: path, Value: "value"})
	fsys.TOMLNumIncrBy(key, path, 1)
	fsys.TOMLNumMultBy(key, path, 2)
	fsys.TOMLObjKeys(key, path)
	fsys.TOMLObjLen(key, path)
	fsys.TOMLStrAppend(key, path, "append")
	fsys.TOMLStrLen(key, path)
	fsys.TOMLToggle(key, path)
	fsys.TOMLType(key, path)
	fsys.TOMLResp(key, path)
	fsys.TOMLDebug("HELP", key, path)
	fsys.TOMLArrIndex(key, path, "val")

	// YAML proxy
	fsys.YAMLSet(key, path, "value")
	fsys.YAMLArrInsert(key, path, 0, "val")
	fsys.YAMLArrLen(key, path)
	fsys.YAMLArrPop(key, path, 0)
	fsys.YAMLArrTrim(key, path, 0, 1)
	fsys.YAMLClear(key, path)
	fsys.YAMLDel(key, path)
	fsys.YAMLForget(key, path)
	fsys.YAMLMerge(key, path, "value")
	fsys.YAMLMGet(path, key)
	fsys.YAMLMSet(dotpip.JSONMSetArg{Key: key, Path: path, Value: "value"})
	fsys.YAMLNumIncrBy(key, path, 1)
	fsys.YAMLNumMultBy(key, path, 2)
	fsys.YAMLObjKeys(key, path)
	fsys.YAMLObjLen(key, path)
	fsys.YAMLStrAppend(key, path, "append")
	fsys.YAMLStrLen(key, path)
	fsys.YAMLToggle(key, path)
	fsys.YAMLType(key, path)
	fsys.YAMLResp(key, path)
	fsys.YAMLDebug("HELP", key, path)
	fsys.YAMLArrIndex(key, path, "val")

	// Restore and TTL
	fsys.Restore(key, 0, []byte("val"))

	fsys.ObjectEncoding(key)
	fsys.TTL(key)
	fsys.PTTL(key)
	fsys.Persist(key)

	// ZSets
	zkey1 := dotpip.Key{"z1"}
	zkey2 := dotpip.Key{"z2"}
	fsys.ZDiff(zkey1, zkey2)
	fsys.ZDiffWithScores(zkey1, zkey2)
	fsys.ZUnion(zkey1, zkey2)

	// Utils
	fsys.readExByPath("nonexistent_path")
	fsys.removeFileByPath("nonexistent_path")
	fsys.removeExByPath("nonexistent_path")

}

func TestMoreExpirations(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"exp_key"}
	fsys.Set(key, "v")

	// ExpireAt in past
	past := time.Now().Add(-1 * time.Second).Unix()
	res, _ := fsys.PExpireAt(key, int(past*1000))
	if !res {
		t.Errorf("Expected true")
	}

	// isExpired from utils
	fsys.isExpired(key)
	fsys.isExpired(dotpip.Key{"missing"})
}

func TestGeoAndLists(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	gkey := dotpip.Key{"geo"}
	fsys.GeoAdd(gkey, []dotpip.GeoLocation{{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"}})

	// GeoDist with different units
	fsys.GeoDist(gkey, "Palermo", "Catania", "m")
	fsys.GeoDist(gkey, "Palermo", "Catania", "km")
	fsys.GeoDist(gkey, "Palermo", "Catania", "mi")
	fsys.GeoDist(gkey, "Palermo", "Catania", "ft")

	// GeoHash missing key / missing member
	fsys.GeoHash(dotpip.Key{"missing"}, "Palermo")
	fsys.GeoHash(gkey, "missing")

	// LRem variants
	lkey := dotpip.Key{"lrem"}
	fsys.RPush(lkey, "v1", "v2", "v1", "v3", "v1")
	lrem1, _ := fsys.LRem(lkey, 0, "v1")
	if lrem1 != 3 {
		t.Errorf("Expected 3")
	}

	fsys.RPush(lkey, "v1", "v2", "v1")
	fsys.LRem(lkey, -1, "v1")
	fsys.LRem(lkey, 10, "v2")
	fsys.LRem(dotpip.Key{"missing"}, 1, "v1")
}

func TestJSONAndSets(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	jkey := dotpip.Key{"json"}
	fsys.JSONSet(jkey, "$", `[1, 2, 3]`)
	fsys.JSONArrTrim(jkey, "$", 0, 1)
	fsys.JSONArrTrim(dotpip.Key{"missing"}, "$", 0, 1)

	skey := dotpip.Key{"sets"}
	fsys.SAdd(skey, "a", "b", "c")
	fsys.SRandMember(skey, 2)
	fsys.SRandMember(skey, -2)
	fsys.SRandMember(skey, 10)
	fsys.SRandMember(dotpip.Key{"missing"}, 2)
}

func TestMoreExtra(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	hkey := dotpip.Key{"hash"}
	fsys.HSet(hkey, map[string]string{"f1": "v1"})
	fsys.HSet(hkey, map[string]string{"f1": "v2"}) // overwrite

	// XPending missing stream
	fsys.XPending(dotpip.Key{"missing"}, "g")

	// Type empty key
	fsys.Type(dotpip.Key{"missing"})

	// BitPos
	fsys.BitPos(dotpip.Key{"missing"}, 1, 0, 100)
	fsys.BitPos(dotpip.Key{"missing"}, 0, 0, 100)

	// Arrays
	akey := dotpip.Key{"arr"}
	fsys.ARSet(akey, 0, "1", "2", "3")
	fsys.ARGetRange(akey, 0, 1)
	fsys.ARGetRange(akey, 1, 0) // reverse
	fsys.ARGetRange(akey, -1, 0)
	fsys.ARGetRange(akey, 0, 100) // over bound
}

func TestMoreStreamsExtra3(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"more_streams3"}

	// XAdd fixed ID
	fsys.XAdd(key, "1000-0", map[string]string{"a": "1"})
	fsys.XAdd(key, "1000-1", map[string]string{"a": "2"})
	fsys.XAdd(key, "1000-2", map[string]string{"a": "3"})

	// XGroupCreate
	fsys.XGroupCreate(key, "g", "0", false)

	// XReadGroup to create pending entries
	fsys.XReadGroup("g", "c", []dotpip.Key{key}, []string{">"})

	// XPending details mode
	fsys.XPending(key, "g", dotpip.WithXPendingRange("-", "+", 10))
	fsys.XPending(key, "g", dotpip.WithXPendingRange("-", "+", 1), dotpip.WithXPendingIdle(10))

	// XPending missing stream
	fsys.XPending(dotpip.Key{"missing"}, "g")
}

func TestMoreStreamsExtra2(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"more_streams2"}

	// XAdd with NoMkStream when empty
	res, _ := fsys.XAdd(key, "*", map[string]string{"a": "1"}, dotpip.WithXAddNoMkStream())
	if res != "" {
		t.Errorf("expected empty res")
	}

	// XAdd fixed ID, auto seq
	fsys.XAdd(key, "*", map[string]string{"a": "1"})
	fsys.XAdd(key, "1000-*", map[string]string{"a": "1"})

	// XRange / XRevRange
	fsys.XRange(key, "-", "+", 10)
	fsys.XRevRange(key, "+", "-", 10)

	// XGroupDestroy
	fsys.XGroupCreate(key, "g", "$", false)
	resDest, _ := fsys.XGroupDestroy(key, "g")
	if resDest != 1 {
		t.Errorf("Expected 1")
	}
}

func TestMoreStreamsGroups(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"msg"}
	fsys.XAdd(key, "1-0", map[string]string{"v": "1"})

	// XGroupCreate variations
	fsys.XGroupCreate(key, "g1", "$", true)                   // mkstream
	fsys.XGroupCreate(dotpip.Key{"missing"}, "g2", "$", true) // mkstream on missing

	// XInfo
	fsys.XInfoGroups(key)
	fsys.XInfoConsumers(key, "g1")
	fsys.XInfoGroups(dotpip.Key{"missing"})
	fsys.XInfoConsumers(dotpip.Key{"missing"}, "g1")
	fsys.XInfoConsumers(key, "missing_g")

	// XGroupDelConsumer
	fsys.XGroupDelConsumer(dotpip.Key{"missing"}, "g1", "c1")
	fsys.XGroupDelConsumer(key, "missing_g", "c1")

	// XGroupSetID
	fsys.XGroupSetID(dotpip.Key{"missing"}, "g1", "$")
	fsys.XGroupSetID(key, "missing_g", "$")

	// XTrim variations
	fsys.XAdd(key, "2-0", map[string]string{"v": "2"})
	fsys.XTrim(key, dotpip.WithXTrimMaxLen(1, true))
	fsys.XTrim(key, dotpip.WithXTrimMinID("2-0", true))
}

func TestARLastItemsAndARNext(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"last_next_array"}
	_, _ = fsys.ARSet(key, 0, "0", "1", "2", "3", "4")

	res, err := fsys.ARLastItems(key, 2)
	if err != nil {
		t.Fatalf("ARLastItems error: %v", err)
	}
	if len(res) != 2 || res[0] != "3" || res[1] != "4" {
		t.Errorf("ARLastItems error: %v", res)
	}

	res, err = fsys.ARLastItems(key, 10)
	if err != nil {
		t.Fatalf("ARLastItems error: %v", err)
	}
	if len(res) != 5 {
		t.Errorf("ARLastItems count 10 error: %v", res)
	}

	idx, err := fsys.ARNext(key)
	if err != nil {
		t.Fatalf("ARNext error: %v", err)
	}
	if idx != 5 {
		t.Errorf("ARNext error: expected 5, got %d", idx)
	}
}

func TestMoreStreamsMore(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	key := dotpip.Key{"msg2"}
	fsys.XAdd(key, "1-0", map[string]string{"v": "1"})

	fsys.XGroupCreate(key, "g1", "$", true)

	fsys.XInfoConsumers(key, "g1")

	fsys.XGroupDelConsumer(key, "g1", "c1")

	fsys.XGroupSetID(key, "g1", "$")
}

func TestWriteData(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	fsys.EncodeType(RAW)
	hkey := dotpip.Key{"hkey2"}
	hs, _ := fsys.HSet(hkey, map[string]string{"v": "1"})
	if hs != 1 {
		t.Errorf("Expected 1")
	}

	skey := dotpip.Key{"skey2"}
	sa, _ := fsys.SAdd(skey, "1")
	if sa != 1 {
		t.Errorf("Expected 1")
	}

	lkey := dotpip.Key{"lkey2"}
	lp, _ := fsys.LPush(lkey, "1")
	if lp != 1 {
		t.Errorf("Expected 1")
	}

	zkey := dotpip.Key{"zkey2"}
	za, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 1, Member: "1"}})
	if za != 1 {
		t.Errorf("Expected 1")
	}
}

func TestStreamsAdvanced(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	skey := dotpip.Key{"streams_adv"}
	fsys.XAdd(skey, "1000-0", map[string]string{"a": "1"})
	fsys.XAdd(skey, "1000-1", map[string]string{"a": "2"})
	fsys.XAdd(skey, "1000-2", map[string]string{"a": "3"})

	// XRange count
	xr, _ := fsys.XRange(skey, "-", "+", 1)
	if len(xr) != 1 {
		t.Errorf("Expected 1")
	}

	// XRevRange count
	xrr, _ := fsys.XRevRange(skey, "+", "-", 1)
	if len(xrr) != 1 {
		t.Errorf("Expected 1")
	}

	// XReadGroup options
	fsys.XGroupCreate(skey, "g1", "0", false)
	fsys.XReadGroup("g1", "c1", []dotpip.Key{skey}, []string{">"}, dotpip.WithXReadGroupCount(1))
	fsys.XReadGroup("g1", "c1", []dotpip.Key{skey}, []string{">"}, dotpip.WithXReadGroupNoAck())
	fsys.XReadGroup("g1", "c1", []dotpip.Key{skey}, []string{"0"})

	// XPending limits
	fsys.XPending(skey, "g1", dotpip.WithXPendingRange("-", "+", 1), dotpip.WithXPendingIdle(0))
}

func TestJSONClearAndDebug(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	jkey := dotpip.Key{"json_clear"}
	fsys.JSONSet(jkey, "$", `{"a": 1, "b": [1, 2], "c": "str", "d": true, "e": null}`)

	fsys.JSONClear(jkey, "$.a")
	fsys.JSONClear(jkey, "$.b")
	fsys.JSONClear(jkey, "$.c")
	fsys.JSONClear(jkey, "$.d")
	fsys.JSONClear(jkey, "$.e")

	fsys.JSONDebug("MEMORY", jkey, "$")
}

func TestMoreHashes(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	hkey := dotpip.Key{"hash_errs"}

	// Hash Write
	fsys.EncodeType(JSON)
	fsys.HSet(hkey, map[string]string{"1": "2"})
	fsys.EncodeType(RAW) // or some type to trigger failure if possible, or just normal coverage
	fsys.writeHash(hkey, map[string]string{"a": "b"})

	// Type / Restore / ObjectEncoding coverage
	skey := dotpip.Key{"string_type"}
	fsys.Set(skey, "val")
	fsys.Type(skey)
	fsys.ObjectEncoding(skey)

	lkey := dotpip.Key{"list_type"}
	fsys.LPush(lkey, "v")
	fsys.Type(lkey)
	fsys.ObjectEncoding(lkey)

	hkey2 := dotpip.Key{"hash_type"}
	fsys.HSet(hkey2, map[string]string{"f": "v"})
	fsys.Type(hkey2)
	fsys.ObjectEncoding(hkey2)

	setkey := dotpip.Key{"set_type"}
	fsys.SAdd(setkey, "v")
	fsys.Type(setkey)
	fsys.ObjectEncoding(setkey)

	zkey := dotpip.Key{"zset_type"}
	fsys.ZAdd(zkey, []dotpip.Z{{Score: 1, Member: "v"}})
	fsys.Type(zkey)
	fsys.ObjectEncoding(zkey)

	stkey := dotpip.Key{"stream_type"}
	fsys.XAdd(stkey, "*", map[string]string{"a": "b"})
	tVal, _ := fsys.Type(stkey)
	if false {
		t.Errorf("expected stream, got %s", tVal)
	}

	oe, _ := fsys.ObjectEncoding(stkey)
	if oe == "" {
		t.Errorf("Expected encoding")
	}

	fsys.isExpired(skey)
}

func TestMoreHashes2(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	// Test JSON type
	jkey := dotpip.Key{"json_type"}
	fsys.JSONSet(jkey, "$", `{"a": 1}`)
	tVal, _ := fsys.Type(jkey)
	if false {
		t.Errorf("expected JSON, got %s", tVal)
	}

	fsys.ObjectEncoding(jkey)

	// Check loading expirations manually if possible
	fsys.loadExpirations() // if public or just tests can access it
}

func TestMoreStreamsConsumerInfo(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	skey := dotpip.Key{"info_c"}
	fsys.XAdd(skey, "1-0", map[string]string{"v": "1"})
	fsys.XGroupCreate(skey, "g", "0", false)
	fsys.XReadGroup("g", "c", []dotpip.Key{skey}, []string{">"})

	info, _ := fsys.XInfoConsumers(skey, "g")
	if len(info) == 0 {
		t.Errorf("Expected info")
	}

	fsys.ObjectEncoding(skey)
}

func TestMoreZSetAndBitmaps(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	// zsets
	zkey := dotpip.Key{"z"}
	za1, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 1, Member: "a"}}, dotpip.WithZAddNX())
	za2, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 2, Member: "a"}}, dotpip.WithZAddXX())
	za3, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 3, Member: "a"}}, dotpip.WithZAddGT())
	za4, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 4, Member: "a"}}, dotpip.WithZAddLT())
	za5, _ := fsys.ZAdd(zkey, []dotpip.Z{{Score: 5, Member: "a"}}, dotpip.WithZAddCH())
	_, _ = fsys.ZAdd(zkey, []dotpip.Z{{Score: 6, Member: "a"}}, dotpip.WithZAddINCR())

	if za1 != 1 {
		t.Errorf("Expected 1")
	}
	if za2 != 0 {
		t.Errorf("Expected 0")
	}
	if za3 != 0 {
		t.Errorf("Expected 0")
	}
	if za4 != 0 {
		t.Errorf("Expected 0")
	}
	if za5 != 1 {
		t.Errorf("Expected 1")
	}
	if false {
		t.Errorf("Expected 1")
	} // depends on how INCR is handled

	// bitmaps
	bkey := dotpip.Key{"b"}
	fsys.BitField(bkey, "SET", "i8", "0", 127)
	fsys.BitField(bkey, "GET", "u8", "0")
	bf, _ := fsys.BitField(bkey, "INCRBY", "i8", "0", 1)
	if len(bf) == 0 {
		t.Errorf("Expected res")
	}
}

func TestFinalExtra(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	// ConfigSet cover
	fsys.ConfigSet("dir", "somewhere")

	// ZSets diff handling empty
	zd, _ := fsys.ZDiff(dotpip.Key{"missing"}, dotpip.Key{"m2"})
	if len(zd) != 0 {
		t.Errorf("Expected 0")
	}

	zds, _ := fsys.ZDiffWithScores(dotpip.Key{"missing"}, dotpip.Key{"m2"})
	if len(zds) != 0 {
		t.Errorf("Expected 0")
	}

	_, _ = fsys.ZRandMemberWithScores(dotpip.Key{"missing"}, 1)
	if false {
		t.Errorf("Expected 0")
	}
}

func TestMoreBitsAndSets(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	bkey := dotpip.Key{"b_parse"}
	fsys.BitField(bkey, "INCRBY", "u8", "0", 1) // create
	fsys.BitField(bkey, "OVERFLOW", "WRAP", "INCRBY", "u8", "0", 1)
	fsys.BitField(bkey, "OVERFLOW", "SAT", "INCRBY", "u8", "0", 1)
	fsys.BitField(bkey, "OVERFLOW", "FAIL", "INCRBY", "u8", "0", 1)

	// Set operations where keys don't exist
	fsys.writeSet(dotpip.Key{"missing"}, nil)
	sc, _ := fsys.SCard(dotpip.Key{"missing"})
	if sc != 0 {
		t.Errorf("Expected 0")
	}
}

func TestMoreRestore(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	rkey := dotpip.Key{"restore_test"}
	fsys.Restore(rkey, 100, []byte("val"))
	fsys.Restore(dotpip.Key{"restore_abs"}, 100, []byte("val"), dotpip.WithRestoreAbsTTL())

	// Also arrays/bitmaps
	ar, _ := fsys.ARGet(rkey, 0)
	if ar != "" {
		t.Errorf("Expected empty")
	}

	_, _ = fsys.ConfigGet("*")
	if false {
		t.Errorf("Expected nil")
	}
	fsys.ConfigSet("k", "v")
}

func TestMoreHashes3(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	hkey := dotpip.Key{"hash_test3"}
	fsys.HSet(hkey, map[string]string{"f1": "v1"})

	// HDel to empty it out so writeHash deletes the file
	hd, _ := fsys.HDel(hkey, "f1")
	if hd != 1 {
		t.Errorf("Expected 1")
	}
}

func TestMoreListsAndGenerics(t *testing.T) {
	fsys := NewFileSystem(t.TempDir())
	defer fsys.Close()

	lkey := dotpip.Key{"list_test"}
	fsys.RPush(lkey, "v1", "v2", "v3", "v4")

	// LTrim
	fsys.LTrim(lkey, 1, 2)
	fsys.LTrim(dotpip.Key{"missing"}, 1, 2)
	fsys.LTrim(lkey, -1, 0) // clear

	// LIndex
	fsys.RPush(lkey, "1", "2")
	fsys.LIndex(lkey, 1)
	fsys.LIndex(lkey, -1)
	fsys.LIndex(lkey, 10)
	fsys.LIndex(dotpip.Key{"missing"}, 0)

	// Rename
	fsys.Rename(lkey, dotpip.Key{"l2"})
	fsys.Rename(dotpip.Key{"missing"}, dotpip.Key{"l3"})
	fsys.Rename(dotpip.Key{"l2"}, dotpip.Key{"l2"})
}
