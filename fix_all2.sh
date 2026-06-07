#!/bin/bash
sed -i 's/	CypherParser = participle\.MustBuild\[CypherQuery\](/\t\/\/ CypherParser parses Cypher queries.\n\tCypherParser = participle.MustBuild[CypherQuery](/' cypher.go
sed -i 's/type CypherQuery struct {/\/\/ CypherQuery represents a parsed Cypher query.\ntype CypherQuery struct {/' cypher.go
sed -i 's/type Clause struct {/\/\/ Clause represents a cypher clause.\ntype Clause struct {/' cypher.go
sed -i 's/type NodePattern struct {/\/\/ NodePattern represents a cypher node pattern.\ntype NodePattern struct {/' cypher.go
sed -i 's/type Property struct {/\/\/ Property represents a node property.\ntype Property struct {/' cypher.go
sed -i 's/type Properties struct {/\/\/ Properties represents a set of node properties.\ntype Properties struct {/' cypher.go
sed -i 's/type Value struct {/\/\/ Value represents a property value.\ntype Value struct {/' cypher.go
sed -i 's/type RelationshipPattern struct {/\/\/ RelationshipPattern represents a relationship between nodes.\ntype RelationshipPattern struct {/' cypher.go
sed -i 's/type RelationshipDetails struct {/\/\/ RelationshipDetails contains relationship info.\ntype RelationshipDetails struct {/' cypher.go
sed -i 's/type PatternElement struct {/\/\/ PatternElement represents an element in a pattern.\ntype PatternElement struct {/' cypher.go
sed -i 's/type PatternChain struct {/\/\/ PatternChain represents a chain of patterns.\ntype PatternChain struct {/' cypher.go
sed -i 's/type MatchClause struct {/\/\/ MatchClause represents a MATCH clause.\ntype MatchClause struct {/' cypher.go
sed -i 's/type CreateClause struct {/\/\/ CreateClause represents a CREATE clause.\ntype CreateClause struct {/' cypher.go
sed -i 's/type ReturnClause struct {/\/\/ ReturnClause represents a RETURN clause.\ntype ReturnClause struct {/' cypher.go
sed -i 's/type DeleteClause struct {/\/\/ DeleteClause represents a DELETE clause.\ntype DeleteClause struct {/' cypher.go
sed -i 's/type SetItem struct {/\/\/ SetItem represents an item in a SET clause.\ntype SetItem struct {/' cypher.go
sed -i 's/type SetClause struct {/\/\/ SetClause represents a SET clause.\ntype SetClause struct {/' cypher.go
