package fs

import (
	"dotpip"
	"testing"
)

var fssa = NewFileSystem("../data")
var dotfs = dotpip.New(fssa)

var firstTestValue = "set testa"
var testKey = dotpip.NewKey("test")

func TestFirstGet(t *testing.T) {
	_, err := dotfs.Get(testKey)
	if err == nil {
		t.Errorf("Should have returned an error for key %s", testKey)
	}

}

func TestFirstSet(t *testing.T) {
	_, err := dotfs.Set(testKey, firstTestValue)

	if err != nil {
		t.Errorf("Should not have returned an error for key %s", testKey)
	}
}

func TestGet(t *testing.T) {
	testValue, err := dotfs.Get(testKey)

	if err != nil {
		t.Errorf("Should not have returned an error for key %s", testKey)
	}

	if testValue != firstTestValue {
		t.Errorf("Expected value '%s' but got '%s'", firstTestValue, testValue)
	}
}

func TestDigest(t *testing.T) {
	testValue := "test value"

	_, err := dotfs.Set(testKey, testValue)
	if err != nil {
		t.Errorf("Should not have returned an error when setting value for key %s", testKey)
	}

	digest, err := dotfs.Digest(testKey)

	if err != nil {
		t.Errorf("Should not have returned an error when digesting key %s", testKey)
	}

	if digest == "" {
		t.Errorf("Digest should not be empty for key %s", testKey)
	}
}

func TestAppend(t *testing.T) {
	appendValue := "test"
	appendLen := dotfs.Append(testKey, appendValue)
	if appendLen <= 0 {
		t.Errorf("Should not have returned less than zero key %s", testKey)
	}

	modifiedValue, err := dotfs.Get(testKey)
	if err != nil {
		t.Errorf("Should not have returned an error when digesting key %s", testKey)
	}

	if appendLen < len(appendValue)+len(firstTestValue) {
		t.Errorf("Append length should %d", len(appendValue)+len(modifiedValue))
	}
}

func TestStrLen(t *testing.T) {
	_ = dotfs.FlushAll()

	strlen := dotfs.StrLen(testKey)
	if strlen != 0 {
		t.Errorf("StrLen of non-existent key should be 0")
	}

	_, _ = dotfs.Set(testKey, "hello")
	strlen = dotfs.StrLen(testKey)
	if strlen != 5 {
		t.Errorf("StrLen of 'hello' should be 5, got %d", strlen)
	}
}

func TestIncrDecr(t *testing.T) {
	_ = dotfs.FlushAll()

	val, err := dotfs.Incr(testKey)
	if err != nil || val != 1 {
		t.Errorf("Incr of new key should be 1")
	}

	val, err = dotfs.IncrBy(testKey, 10)
	if err != nil || val != 11 {
		t.Errorf("IncrBy 10 should be 11, got %d", val)
	}

	val, err = dotfs.Decr(testKey)
	if err != nil || val != 10 {
		t.Errorf("Decr should be 10")
	}

	val, err = dotfs.DecrBy(testKey, 5)
	if err != nil || val != 5 {
		t.Errorf("DecrBy 5 should be 5")
	}

	_, _ = dotfs.Set(testKey, "abc")
	_, err = dotfs.Incr(testKey)
	if err == nil {
		t.Errorf("Incr on non-integer string should return error")
	}
}

func TestIncrByFloat(t *testing.T) {
	_ = dotfs.FlushAll()

	val, err := dotfs.IncrByFloat(testKey, 10.5)
	if err != nil || val != 10.5 {
		t.Errorf("IncrByFloat of new key should be 10.5")
	}

	val, err = dotfs.IncrByFloat(testKey, 0.1)
	if err != nil || val != 10.6 {
		t.Errorf("IncrByFloat should handle precision")
	}
}

func TestGetDel(t *testing.T) {
	_ = dotfs.FlushAll()
	_, _ = dotfs.Set(testKey, "value")

	val, err := dotfs.GetDel(testKey)
	if err != nil || val != "value" {
		t.Errorf("GetDel should return value")
	}

	_, err = dotfs.Get(testKey)
	if err == nil {
		t.Errorf("Key should be deleted after GetDel")
	}
}

func TestGetRange(t *testing.T) {
	_ = dotfs.FlushAll()
	_, _ = dotfs.Set(testKey, "This is a string")

	val, _ := dotfs.GetRange(testKey, 0, 3)
	if val != "This" {
		t.Errorf("GetRange(0, 3) expected 'This', got '%s'", val)
	}

	val, _ = dotfs.GetRange(testKey, -3, -1)
	if val != "ing" {
		t.Errorf("GetRange(-3, -1) expected 'ing', got '%s'", val)
	}
}

func TestSetRange(t *testing.T) {
	_ = dotfs.FlushAll()
	_, _ = dotfs.Set(testKey, "Hello World")

	newLen, err := dotfs.SetRange(testKey, 6, "Redis")
	if err != nil || newLen != 11 {
		t.Errorf("SetRange failed, newLen: %d", newLen)
	}

	val, _ := dotfs.Get(testKey)
	if val != "Hello Redis" {
		t.Errorf("SetRange value expected 'Hello Redis', got '%s'", val)
	}
}

func TestMGetMSet(t *testing.T) {
	_ = dotfs.FlushAll()

	k1 := dotpip.NewKey("k1")
	k2 := dotpip.NewKey("k2")

	_ = dotfs.MSet(dotpip.KV{Key: k1, Value: "v1"}, dotpip.KV{Key: k2, Value: "v2"})

	vals, _ := dotfs.MGet(k1, k2, dotpip.NewKey("k3"))
	if len(vals) != 3 || vals[0] != "v1" || vals[1] != "v2" || vals[2] != "" {
		t.Errorf("MGet returned incorrect values: %v", vals)
	}
}

func TestMSetNX(t *testing.T) {
	_ = dotfs.FlushAll()

	k1 := dotpip.NewKey("k1")
	k2 := dotpip.NewKey("k2")

	success, err := dotfs.MSetNX(dotpip.KV{Key: k1, Value: "v1"}, dotpip.KV{Key: k2, Value: "v2"})
	if err != nil || !success {
		t.Errorf("MSetNX should succeed on empty DB")
	}

	k3 := dotpip.NewKey("k3")
	success, err = dotfs.MSetNX(dotpip.KV{Key: k2, Value: "v2_new"}, dotpip.KV{Key: k3, Value: "v3"})
	if err != nil || success {
		t.Errorf("MSetNX should fail if one key exists")
	}

	val, _ := dotfs.Get(k3)
	if val != "" {
		t.Errorf("k3 should not be set because MSetNX failed")
	}
}

func TestEncodings(t *testing.T) {
	encodings := []FileEncodeType{JSON, YAML, TOML}

	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			testFS := NewFileSystem("../data")
			testFS.EncodeType(enc)
			dfs := dotpip.New(testFS)
			_ = dfs.FlushAll()

			key := dotpip.NewKey("test_enc")
			val := "hello encoding"

			_, err := dfs.Set(key, val)
			if err != nil {
				t.Fatalf("Failed to set with %s encoding: %v", enc, err)
			}

			getVal, err := dfs.Get(key)
			if err != nil {
				t.Fatalf("Failed to get with %s encoding: %v", enc, err)
			}

			if getVal != val {
				t.Errorf("Expected %q, got %q with %s encoding", val, getVal, enc)
			}
		})
	}
}
