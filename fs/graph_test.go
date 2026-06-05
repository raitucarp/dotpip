package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"path/filepath"
	"testing"
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

	resRO, err := dotfs.GraphROQuery(key, "MATCH (n:Person)-[:KNOWS]->(m)-[:WORKS_AT]->(c) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	if len(resRO) < 2 {
		t.Fatal("Expected response")
	}

	nodesCalculated := resRO[0][string(dotpip.GraphKeywordNodesFound)].(int)
	pathsMatched := resRO[0][string(dotpip.GraphKeywordPathsMatched)].(int)

	if nodesCalculated != 8 {
		t.Fatalf("Expected 8 nodes found across executions, got %d", nodesCalculated)
	}

	if pathsMatched != 1 {
		t.Fatalf("Expected 1 deep path matched corresponding to length 2 across graph edges, got %d", pathsMatched)
	}

	// Test actual RETURN parameters mapping variable 'n'
	if nData, ok := resRO[1]["n"]; ok {
		nList := nData.([]any)
		if len(nList) == 0 {
			t.Fatal("Expected at least one returned node variable for n")
		}
	} else {
		t.Fatal("Expected returned structure to contain key 'n'")
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
		t.Fatalf("Expected slowlog")
	}

	_, err = dotfs.GraphSlowlog(dotpip.NewKey("not_exist"))
	if err != nil {
		t.Fatal(err)
	}

	resDelete, err := dotfs.GraphQuery(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resDelete) < 2 || resDelete[1][string(dotpip.GraphKeywordNodesDeleted)] == nil {
		t.Fatal("Expected nodes deleted response")
	}

	_, err = dotfs.GraphQuery(key, "CREATE (n:Person)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphQuery(key, "MATCH (n) DELETE n")
	if err != nil {
		t.Fatal(err)
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

	_, err = dotfs.GraphROQuery(key, "MATCH (n) DELETE n")
	if err == nil {
		t.Fatal("Expected error")
	}

	_, err = dotfs.GraphROQuery(key, "MATCH (n) SET n.name = 'Bob'")
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

func TestGraphPropertiesParsing(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	// Verify property setting works
	_, err := dotfs.GraphQuery(key, "CREATE (a:Person {name: 'Alice', age: 30, active: TRUE})")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dotfs.GraphQuery(key, "MATCH (a) SET a.age = 31, a.active = FALSE")
	if err != nil {
		t.Fatal(err)
	}

	resProf, err := dotfs.GraphProfile(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resProf) == 0 {
		t.Fatal("Expected response")
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
}

func TestGraphMissingMatchClause(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	resRO, err := dotfs.GraphROQuery(key, "RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resRO) == 0 {
		t.Fatal("Expected response handled safely")
	}
}

func TestGraphMatchNoVariables(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	_, _ = dotfs.GraphQuery(key, "CREATE (:Person)")

	// test no variable mapping
	resRO, err := dotfs.GraphROQuery(key, "MATCH (:Person) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	// length handles mapping dynamically; should correctly parse safe returning empty for unspecified references without failing test.
	if len(resRO) < 2 {
		t.Fatal("Expected parsed execution map safe response")
	}
}

func TestGraphQueryEmptyGraph(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	res, err := dotfs.GraphQuery(key, "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) < 2 {
		t.Fatal("Expected empty match execution returned properly")
	}
}

func TestGraphMissingMatchClauseQuery(t *testing.T) {
	dotfs := fs.NewFileSystem(filepath.Join(t.TempDir(), "db"))
	defer func() { _ = dotfs.FlushAll() }()

	key := dotpip.NewKey("mygraph")

	resRO, err := dotfs.GraphQuery(key, "RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resRO) == 0 {
		t.Fatal("Expected response handled safely")
	}
}
