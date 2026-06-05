package fs_test

import (
	"testing"
	"dotpip"
	"dotpip/fs"
	"path/filepath"
)

func TestGraphCommands(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	res, err := dotfs.GraphQuery(key, "CREATE (:Person {name: 'Alice'})-[:KNOWS]->(:Person)")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("Expected response")
	}

	// Make sure it really persists
	resRO, err := dotfs.GraphROQuery(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	if len(resRO) == 0 {
		t.Fatal("Expected response")
	}

	// Test the actual value logic returned from MATCH based on Gonum calculations
	nodesCalculated := resRO[0]["NodesCalculated"].(int)
	edgesCalculated := resRO[0]["EdgesCalculated"].(int)

	if nodesCalculated != 2 {
		t.Fatalf("Expected 2 nodes calculated, got %d", nodesCalculated)
	}

	if edgesCalculated != 1 {
		t.Fatalf("Expected 1 edge calculated, got %d", edgesCalculated)
	}

	resExp, err := dotfs.GraphExplain(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) == 0 || resExp[0] != "MATCH" {
		t.Fatalf("Expected MATCH explanation, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "CREATE (n)")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) == 0 || resExp[0] != "CREATE" {
		t.Fatalf("Expected CREATE explanation, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) < 2 || resExp[1] != "DELETE" {
		t.Fatalf("Expected DELETE explanation, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "MATCH (n) SET n.name='Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) < 2 || resExp[1] != "SET" {
		t.Fatalf("Expected SET explanation, got %v", resExp)
	}

	_, err = dotfs.GraphList()
	if err != nil {
		t.Fatal(err)
	}

	resProf, err := dotfs.GraphProfile(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) == 0 || resProf[0] != "MATCH" {
		t.Fatalf("Expected MATCH profile, got %v", resProf)
	}

	resProf, err = dotfs.GraphProfile(key, "CREATE (n)")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) == 0 || resProf[0] != "CREATE" {
		t.Fatalf("Expected CREATE profile, got %v", resProf)
	}

	resProf, err = dotfs.GraphProfile(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) < 2 || resProf[1] != "DELETE" {
		t.Fatalf("Expected DELETE profile, got %v", resProf)
	}

	resProf, err = dotfs.GraphProfile(key, "MATCH (n) SET n.name='Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) < 2 || resProf[1] != "SET" {
		t.Fatalf("Expected SET profile, got %v", resProf)
	}

	_, err = dotfs.GraphSlowlog(key)
	if err != nil {
		t.Fatal(err)
	}

	resDelete, err := dotfs.GraphQuery(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resDelete) < 2 || resDelete[1]["NodesDeleted"] == nil {
		t.Fatal("Expected nodes deleted response")
	}

	resSet, err := dotfs.GraphQuery(key, "MATCH (n) SET n.name = 'Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resSet) < 2 || resSet[1]["PropertiesSet"] == nil {
		t.Fatal("Expected properties set response")
	}

	_, err = dotfs.GraphDelete(key)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGraphROQueryError(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	_, err := dotfs.GraphROQuery(key, "CREATE (:Person {name: 'Alice'})")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGraphQueryParseError(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

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
