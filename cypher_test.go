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
		"MATCH (n:Person {age: 30})-[:FRIEND]->(m:Person)-[:WORKS_AT]->(c:Company {name: 'Tech'}) RETURN n, m, c",
		"CREATE (p:User {id: 1, active: TRUE})-[:FOLLOWS {since: '2023'}]->(q:User {id: 2})",
		"OPTIONAL MATCH (x)-[r:LIKES|LOVES]->(y) RETURN r",
		"MATCH (a)<-[:PARENT]-(b)-[:SIBLING]->(c) SET a.age = 50, c.status = 'single'",
	}

	for _, query := range queries {
		_, err := CypherParser.ParseString("", query)
		if err != nil {
			t.Fatalf("Failed to parse query %q: %v", query, err)
		}
	}
}

func TestCypherParserInvalid(t *testing.T) {
	queries := []string{
		"MATCH n DELETE n",         // Missing parenthesis for node
		"CREATE (n SET n.val = 1",  // Invalid grammar mixed
		"MATCH ()--()--()--() SET", // Trailing SET
	}

	for _, query := range queries {
		_, err := CypherParser.ParseString("", query)
		if err == nil {
			t.Fatalf("Expected error for query %q, but got none", query)
		}
	}
}
