package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"path/filepath"
	"testing"
)

func TestGraphCommands(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer dotfs.FlushAll()

	key := dotpip.NewKey("mygraph")

	res, err := dotfs.GraphQuery(key, "CREATE (:Person {name: 'Alice'})-[:KNOWS]->(:Person)")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("Expected response")
	}

	resRO, err := dotfs.GraphROQuery(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	if len(resRO) == 0 {
		t.Fatal("Expected response")
	}

	_, err = dotfs.GraphExplain(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphExplain(key, "CREATE (n)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphExplain(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphExplain(key, "MATCH (n) SET n.name='Bob'")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphList()
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphProfile(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphProfile(key, "CREATE (n)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphProfile(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphProfile(key, "MATCH (n) SET n.name='Bob'")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphSlowlog(key)
	if err != nil {
		t.Fatal(err)
	}

	resDelete, err := dotfs.GraphQuery(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resDelete) == 0 {
		t.Fatal("Expected response")
	}

	resSet, err := dotfs.GraphQuery(key, "MATCH (n) SET n.name = 'Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resSet) == 0 {
		t.Fatal("Expected response")
	}

	_, err = dotfs.GraphDelete(key)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGraphROQueryError(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer dotfs.FlushAll()

	key := dotpip.NewKey("mygraph")

	_, err := dotfs.GraphROQuery(key, "CREATE (:Person {name: 'Alice'})")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGraphQueryParseError(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer dotfs.FlushAll()

	key := dotpip.NewKey("mygraph")

	_, err := dotfs.GraphQuery(key, "INVALID QUERY")
	if err == nil {
		t.Fatal("Expected error")
	}

	_, err = dotfs.GraphROQuery(key, "INVALID QUERY")
	if err == nil {
		t.Fatal("Expected error")
	}

	_, err = dotfs.GraphExplain(key, "INVALID QUERY")
	if err == nil {
		t.Fatal("Expected error")
	}

	_, err = dotfs.GraphProfile(key, "INVALID QUERY")
	if err == nil {
		t.Fatal("Expected error")
	}
}
