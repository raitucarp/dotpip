package dotpip

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	cypherLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Keyword", Pattern: `(?i)\b(MATCH|OPTIONAL|WHERE|WITH|RETURN|ORDER|BY|SKIP|LIMIT|ASC|DESC|ASCENDING|DESCENDING|CREATE|MERGE|SET|DELETE|REMOVE|DETACH|AND|OR|NOT|XOR|TRUE|FALSE|NULL|AS)\b`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "String", Pattern: `'[^']*'|"[^"]*"`},
		{Name: "Number", Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{Name: "Punct", Pattern: `[-+*/<>!=(){}\[\],.:;|]`},
		{Name: "Whitespace", Pattern: `\s+`},
	})

	// CypherParser parses Cypher queries.
	CypherParser = participle.MustBuild[CypherQuery](
		participle.Lexer(cypherLexer),
		participle.Elide("Whitespace"),
		participle.UseLookahead(2),
		participle.CaseInsensitive("Keyword"),
	)
)

// CypherQuery represents a parsed Cypher query.
type CypherQuery struct {
	Clauses []*Clause `parser:"@@*"`
}

// Clause represents a cypher clause.
type Clause struct {
	Match  *MatchClause  `parser:"  @@"`
	Create *CreateClause `parser:"| @@"`
	Return *ReturnClause `parser:"| @@"`
	Delete *DeleteClause `parser:"| @@"`
	Set    *SetClause    `parser:"| @@"`
}

// NodePattern represents a cypher node pattern.
type NodePattern struct {
	Variable   *string     `parser:"\"(\" @Ident?"`
	Labels     []string    `parser:"(\":\" @Ident)*"`
	Properties *Properties `parser:"@@? \")\""`
}

// Property represents a node property.
type Property struct {
	Key   string `parser:"@Ident \":\""`
	Value *Value `parser:"@@"`
}

// Properties represents a set of node properties.
type Properties struct {
	Props []*Property `parser:"\"{\" @@ (\",\" @@)* \"}\""`
}

// Value represents a property value.
type Value struct {
	String *string  `parser:"@String"`
	Number *float64 `parser:"| @Number"`
	Bool   *bool    `parser:"| @(\"TRUE\" | \"FALSE\")"`
}

// RelationshipPattern represents a relationship between nodes.
type RelationshipPattern struct {
	LeftArrow  bool                 `parser:"@\"<\"? \"-\""`
	Details    *RelationshipDetails `parser:"(\"[\" @@ \"]\")? \"-\""`
	RightArrow bool                 `parser:"@\">\"?"`
}

// RelationshipDetails contains relationship info.
type RelationshipDetails struct {
	Variable   *string     `parser:"@Ident?"`
	Types      []string    `parser:"(\":\" @Ident (\"|\" @Ident)*)?"`
	Properties *Properties `parser:"@@?"`
}

// PatternElement represents an element in a pattern.
type PatternElement struct {
	Node  *NodePattern    `parser:"@@"`
	Chain []*PatternChain `parser:"@@*"`
}

// PatternChain represents a chain of patterns.
type PatternChain struct {
	Relationship *RelationshipPattern `parser:"@@"`
	Node         *NodePattern         `parser:"@@"`
}

// MatchClause represents a MATCH clause.
type MatchClause struct {
	Optional bool            `parser:"@\"OPTIONAL\"? \"MATCH\""`
	Pattern  *PatternElement `parser:"@@"`
}

// CreateClause represents a CREATE clause.
type CreateClause struct {
	Pattern *PatternElement `parser:"\"CREATE\" @@"`
}

// ReturnClause represents a RETURN clause.
type ReturnClause struct {
	Items []string `parser:"\"RETURN\" @Ident (\",\" @Ident)*"`
}

// DeleteClause represents a DELETE clause.
type DeleteClause struct {
	Detach bool     `parser:"@\"DETACH\"? \"DELETE\""`
	Items  []string `parser:"@Ident (\",\" @Ident)*"`
}

// SetItem represents an item in a SET clause.
type SetItem struct {
	Variable string `parser:"@Ident"`
	Property string `parser:"(\".\" @Ident)?"`
	Operator string `parser:"\"=\""`
	Value    *Value `parser:"@@"`
}

// SetClause represents a SET clause.
type SetClause struct {
	Items []*SetItem `parser:"\"SET\" @@ (\",\" @@)*"`
}
