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
	local val = Get(KEYS[1])
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
	local err = UNKNOWN_CMD()
	return err
	`
	_, err := fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for unknown command")
	}

	// Since pcall catches the error, it might return a string or table.
	// Actually pcall is designed to return the error to Lua.
	// Let's test a simple pcall:
	scriptPCallRedis := `
	local status, err = pcall(UNKNOWN_CMD)
	if not status then
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

	_, _ = fsys.Set(dotpip.NewKey("foo"), "bar")

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

func TestLuaScriptsEmptyArgs(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// should be able to append empty string
	res, err := fsys.Eval(`return Append("testkey", "")`, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if num, ok := res.(float64); !ok || num != 0 {
		t.Fatalf("Expected 0, got %v", res)
	}
}

func TestLuaScriptsErrorScenarios2(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// EXISTS takes at least one parameter
	script := `return Exists()`
	_, err := fsys.Eval(script, 0, nil, nil)
	// It actually doesn't error out because Exists is variadic! `Exists(keys ...Key)` can take 0 keys and returns an empty slice!
	// It is converted to an empty Lua table, which becomes a map[] in convertLuaToGo if empty
	if err != nil {
		t.Fatalf("Exists without args should not error, got %v", err)
	}

	// Missing args for non-variadic command: GET takes 1 arg
	script = `return Get()`
	_, err = fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for missing args on GET")
	}
}

func TestLuaScriptsErrorScenarios3(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// Missing args for non-variadic command in pcall
	script := `local status, err = pcall(Get); if not status then return {err=err} else return err end`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	resMap, ok := res.(map[string]any)
	if !ok || resMap["err"] == nil {
		t.Fatalf("Expected error table from pcall")
	}
}

func TestLuaScriptsGoMethodCallError(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// Should trigger method error (Get on missing key)
	script := `return Get("missing")`
	_, err := fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for missing key in GET")
	}

	// Should trigger method error in pcall (Get on missing key)
	script = `local status, err = pcall(Get, "missing"); if not status then return {err=err} else return err end`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	resMap, ok := res.(map[string]any)
	if !ok || resMap["err"] == nil {
		t.Fatalf("Expected error table from pcall")
	}
}

func TestLuaScriptsVariousReturnTypes(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// bool
	res, _ := fsys.Eval(`return true`, 0, nil, nil)
	if v, ok := res.(bool); !ok || !v {
		t.Fatalf("Expected true")
	}

	// number
	res, _ = fsys.Eval(`return 123.45`, 0, nil, nil)
	if v, ok := res.(float64); !ok || v != 123.45 {
		t.Fatalf("Expected 123.45")
	}
}

func TestLuaScriptsGoMethodArgs(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// passing number
	res, _ := fsys.Eval(`return IncrBy("numkey", 5)`, 0, nil, nil)
	if v, ok := res.(float64); !ok || v != 5 {
		t.Fatalf("Expected 5")
	}
	res, _ = fsys.Eval(`return IncrByFloat("numkey", 1.5)`, 0, nil, nil)
	if v, ok := res.(float64); !ok || v != 6.5 {
		t.Fatalf("Expected 6.5")
	}

	// passing boolean (just to test boolean argument logic, we might need a dummy command)
	_, _ = fsys.Eval(`return ARInfo("numkey", true)`, 0, nil, nil)
}

func TestLuaScriptsExcessArgs(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// Call non-variadic with too many arguments
	_, _ = fsys.Set(dotpip.NewKey("a"), "1")
	script := `return Get("a", "b", "c")`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Expected Get to ignore excess args or handle it, got error: %v", err)
	}
	if v, ok := res.(string); !ok || v != "1" {
		t.Fatalf("Expected 1, got %v", res)
	}
}

func TestLuaScriptsErrorScenarios4(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()


	// Passing table to string
	script := `return Get({a=1})`
	_, err := fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for passing table to string arg")
	}

	// Passing string to struct slice variadic
	script = `return ZAdd("key", "not_a_slice")`
	_, err = fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for passing string to ZAdd elements")
	}

	// Invalid array structure (containing mixed unsupported types) for Slice
	script = `return MGet(1)`
	_, err = fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Expected MGet(1) to be valid (it wraps string automatically), got %v", err)
	}

	// Unimplemented method natively
	script = `return UNKNOWN_CMD()`
	_, err = fsys.Eval(script, 0, nil, nil)
	if err == nil {
		t.Fatalf("Expected error for UNKNOWN_CMD")
	}

}
