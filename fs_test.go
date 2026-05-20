package dotpip

import "testing"

var fs = FileSystem("./data")
var dotface = New(fs)

var firstTestValue = "set testa"
var testKey = NewKey("test")

func TestFlushAll(t *testing.T) {
	err := dotface.FlushAll()
	if err != nil {
		t.Errorf("Should not have returned an error when flushing all data")
	}
}

func TestFirstGet(t *testing.T) {

	_, err := dotface.Get(testKey)

	if err == nil {
		t.Errorf("Should have returned an error for key %s", testKey)
	}
}

func TestFirstSet(t *testing.T) {
	_, err := dotface.Set(testKey, firstTestValue)

	if err != nil {
		t.Errorf("Should not have returned an error for key %s", testKey)
	}
}

func TestGet(t *testing.T) {
	testValue, err := dotface.Get(testKey)

	if err != nil {
		t.Errorf("Should not have returned an error for key %s", testKey)
	}

	if testValue != firstTestValue {
		t.Errorf("Expected value '%s' but got '%s'", firstTestValue, testValue)
	}
}
