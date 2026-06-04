package fs

import (
	"dotpip"
	"testing"
)

func TestLuaScriptsBasic(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	_, _ = fsys.Set(dotpip.NewKey("mykey"), "myvalue")

	script := `
	local val = redis.call("get", KEYS[1])
	return val
	`
	res, err := fsys.Eval(script, 1, []string{"mykey"}, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	if res != "myvalue" {
		t.Fatalf("Expected 'myvalue', got '%v'", res)
	}
}

func TestLuaScriptsDirectGlobal(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `
	Set("test", "hello")
	return Get("test")
	`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	if res != "hello" {
		t.Fatalf("Expected 'hello', got '%v'", res)
	}
}

func TestLuaScriptExistsLoadFlush(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `return 1`
	hash, err := fsys.ScriptLoad(script)
	if err != nil {
		t.Fatalf("ScriptLoad error: %v", err)
	}

	exists, err := fsys.ScriptExists(hash)
	if err != nil {
		t.Fatalf("ScriptExists error: %v", err)
	}
	if len(exists) != 1 || !exists[0] {
		t.Fatalf("Expected script to exist")
	}

	err = fsys.ScriptFlush()
	if err != nil {
		t.Fatalf("ScriptFlush error: %v", err)
	}

	exists, err = fsys.ScriptExists(hash)
	if err != nil {
		t.Fatalf("ScriptExists error: %v", err)
	}
	if len(exists) != 1 || exists[0] {
		t.Fatalf("Expected script to be flushed")
	}
}

func TestLuaScriptsErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `
	local err = redis.call("UNKNOWN_CMD")
	return err
	`
	_, err := fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for unknown command")
	}

	// Since pcall catches the error, it might return a string or table.
	// Actually redis.pcall is designed to return the error to Lua.
	// Let's test a simple pcall:
	scriptPCallRedis := `
	local err = redis.pcall("UNKNOWN_CMD")
	if type(err) == "table" and err.err then
		return "caught error"
	end
	return "no error"
	`
	res, err := fsys.Eval(scriptPCallRedis, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	if res != "caught error" {
		t.Fatalf("Expected 'caught error', got '%v'", res)
	}
}

func TestLuaScriptsArrayReturnType(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	_, _ = fsys.Set(dotpip.NewKey("mykey"), "myvalue")

	script := `
	return {1, 2, "hello"}
	`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	resSlice, ok := res.([]any)
	if !ok {
		t.Fatalf("Expected array return type")
	}
	if len(resSlice) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(resSlice))
	}
	if resSlice[0] != float64(1) {
		t.Fatalf("Expected first element to be 1")
	}
	if resSlice[2] != "hello" {
		t.Fatalf("Expected third element to be 'hello'")
	}
}

func TestLuaScriptsMapReturnType(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `
	return {a=1, b="hello"}
	`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	resMap, ok := res.(map[string]any)
	if !ok {
		t.Fatalf("Expected map return type, got %T", res)
	}
	if len(resMap) != 2 {
		t.Fatalf("Expected map length 2")
	}
	if resMap["a"] != float64(1) {
		t.Fatalf("Expected key 'a' to be 1")
	}
	if resMap["b"] != "hello" {
		t.Fatalf("Expected key 'b' to be 'hello'")
	}
}

func TestLuaScriptsEvalSha(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `return "evalsha"`
	hash, err := fsys.ScriptLoad(script)
	if err != nil {
		t.Fatalf("ScriptLoad error: %v", err)
	}

	res, err := fsys.EvalSha(hash, 0, nil, nil)
	if err != nil {
		t.Fatalf("EvalSha error: %v", err)
	}

	if res != "evalsha" {
		t.Fatalf("Expected 'evalsha', got '%v'", res)
	}

	_, err = fsys.EvalSha("invalidhash", 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for non-existent script hash")
	}
}

func TestLuaScriptsStructArg(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `
	local z = {score=1.5, member="mem1"}
	return ZAdd("myzset", {z})
	`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	count, ok := res.(float64)
	if !ok || count != 1 {
		t.Fatalf("Expected ZAdd to add 1 member, got %v", res)
	}
}

func TestLuaScriptsRO(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	script := `return 1`
	hash, err := fsys.ScriptLoad(script)
	if err != nil {
		t.Fatalf("ScriptLoad error: %v", err)
	}

	res, err := fsys.EvalRO(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("EvalRO error: %v", err)
	}

	if num, ok := res.(float64); !ok || num != 1 {
		t.Fatalf("Expected 1, got %v", res)
	}

	res, err = fsys.EvalShaRO(hash, 0, nil, nil)
	if err != nil {
		t.Fatalf("EvalShaRO error: %v", err)
	}

	if num, ok := res.(float64); !ok || num != 1 {
		t.Fatalf("Expected 1, got %v", res)
	}
}

func TestLuaScriptKill(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	err := fsys.ScriptKill()
	if err != nil {
		t.Fatalf("Expected no error from ScriptKill, got %v", err)
	}
}


func TestLuaScriptsReturnTypes(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	fsys.Set(dotpip.NewKey("foo"), "bar")

	// Array of primitives
	res, err := fsys.Eval(`return {"a", 1, true}`, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	arr, ok := res.([]any)
	if !ok || len(arr) != 3 {
		t.Fatalf("Expected slice of length 3, got %T %v", res, res)
	}

	// Nil
	res, err = fsys.Eval(`return nil`, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if res != nil {
		t.Fatalf("Expected nil, got %v", res)
	}
}
