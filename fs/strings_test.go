package fs

import (
	"dotpip"
	"testing"
)

var fssa = FileSystem("../data")
var dotfs = dotpip.New(fssa)

var firstTestValue = "set testa"
var testKey = dotpip.NewKey("test")

func TestFlushAll(t *testing.T) {
	err := dotfs.FlushAll()
	if err != nil {
		t.Errorf("Should not have returned an error when flushing all data")
	}
}

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
