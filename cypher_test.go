package dotpip

import (
	"testing"
)

func TestCypherParser(t *testing.T) {
	queries := []string{
		"MATCH (n:Person {name: 'Alice'})-[:KNOWS]->(m:Person) RETURN m",
		"CREATE (:Person {name: 'Bob'})",
		"MATCH (n) DELETE n",
		"MATCH (n) SET n.name = 'Charlie'",
	}

	for _, query := range queries {
		_, err := CypherParser.ParseString("", query)
		if err != nil {
			t.Fatalf("Failed to parse query %q: %v", query, err)
		}
	}
}
