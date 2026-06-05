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

	createQueries := []string{
		"CREATE (:Person {name: 'Alice'})-[:KNOWS]->(:Person {name: 'Bob'})-[:WORKS_AT]->(:Company {name: 'Tech'})",
		"CREATE (:Person {name: 'Charlie'})-[:KNOWS]->(:Person {name: 'Alice'})",
		"CREATE (:Company {name: 'Startup'})<-[:FOUNDED]-(:Person {name: 'Dave'})-[:KNOWS]->(:Person {name: 'Charlie'})",
	}

	for _, query := range createQueries {
		res, err := dotfs.GraphQuery(key, query)
		if err != nil {
			t.Fatalf("Failed to execute query %q: %v", query, err)
		}
		if len(res) == 0 {
			t.Fatal("Expected response")
		}
	}

	resRO, err := dotfs.GraphROQuery(key, "MATCH (n)-[:KNOWS]->(m)-[:WORKS_AT]->(c)")
	if err != nil {
		t.Fatal(err)
	}

	if len(resRO) == 0 {
		t.Fatal("Expected response")
	}

	nodesCalculated := resRO[0][string(dotpip.GraphKeywordNodesFound)].(int)
	pathsMatched := resRO[0][string(dotpip.GraphKeywordPathsMatched)].(int)

	if nodesCalculated != 8 {
		t.Fatalf("Expected 8 nodes found across executions, got %d", nodesCalculated)
	}

	if pathsMatched != 2 {
		t.Fatalf("Expected 2 deep path matched corresponding to length 2 across graph edges, got %d", pathsMatched)
	}

	var dpip dotpip.DotPip = dotfs
	dpRes, dpErr := dpip.GraphQuery(key, "MATCH (n) RETURN n")
	if dpErr != nil {
		t.Fatal(dpErr)
	}
	if len(dpRes) == 0 {
		t.Fatal("Expected response")
	}

	dpResRO, dpErrRO := dpip.GraphROQuery(key, "MATCH (n) RETURN n")
	if dpErrRO != nil {
		t.Fatal(dpErrRO)
	}
	if len(dpResRO) == 0 {
		t.Fatal("Expected response")
	}

	resProf, err := dotfs.GraphProfile(key, "MATCH (n)-[:KNOWS]->(m)-[:WORKS_AT]->(c)")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) == 0 || resProf[0] != string(dotpip.GraphKeywordMatch) {
		t.Fatalf("Expected MATCH profile, got %v", resProf)
	}

	// TEST EXPLAIN
	resExp, err := dotfs.GraphExplain(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) == 0 || resExp[0] != string(dotpip.GraphKeywordMatch) {
		t.Fatalf("Expected MATCH profile, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "CREATE (n)")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) == 0 || resExp[0] != string(dotpip.GraphKeywordCreate) {
		t.Fatalf("Expected CREATE explanation, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) < 2 || resExp[1] != string(dotpip.GraphKeywordDelete) {
		t.Fatalf("Expected DELETE explanation, got %v", resExp)
	}

	resExp, err = dotfs.GraphExplain(key, "MATCH (n) SET n.name='Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resExp) < 2 || resExp[1] != string(dotpip.GraphKeywordSet) {
		t.Fatalf("Expected SET explanation, got %v", resExp)
	}


	resSet, err := dotfs.GraphQuery(key, "MATCH (n) SET n.name = 'Bob'")
	if err != nil {
		t.Fatal(err)
	}
	if len(resSet) < 2 || resSet[1][string(dotpip.GraphKeywordPropertiesSet)] == nil {
		t.Fatal("Expected properties set response")
	}

	// Test List and Slowlog functionality properly BEFORE DELETE
	keysList, err := dotfs.GraphList()
	if err != nil {
		t.Fatal(err)
	}
	if len(keysList) == 0 {
		t.Fatalf("Expected at least one graph in list, got none")
	}

	slowlogRes, err := dotfs.GraphSlowlog(key)
	if err != nil {
		t.Fatal(err)
	}
	if len(slowlogRes) == 0 {
		t.Fatalf("Expected mocked slowlog, got none")
	}

	resDelete, err := dotfs.GraphQuery(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resDelete) < 2 || resDelete[1][string(dotpip.GraphKeywordNodesDeleted)] == nil {
		t.Fatal("Expected nodes deleted response")
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
