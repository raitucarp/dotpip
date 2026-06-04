package fs

import (
	"testing"
)

func TestLuaScriptSandbox(t *testing.T) {
	tempDir := t.TempDir()
	fsys := NewFileSystem(tempDir)
	defer fsys.Close()

	// Should not have io
	script := `return type(io)`
	res, err := fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if res != "nil" {
		t.Fatalf("Expected io to be nil, got %v", res)
	}

	// Should not have os
	script = `return type(os)`
	res, err = fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if res != "nil" {
		t.Fatalf("Expected os to be nil, got %v", res)
	}

	// Should not have package or require
	script = `return type(package)`
	res, err = fsys.Eval(script, 0, nil, nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if res != "nil" {
		t.Fatalf("Expected package to be nil, got %v", res)
	}
}
