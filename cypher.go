package dotpip

import (
	participle "github.com/alecthomas/participle/v2"
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

	CypherParser = participle.MustBuild[CypherQuery](
		participle.Lexer(cypherLexer),
		participle.Elide("Whitespace"),
		participle.UseLookahead(2),
		participle.CaseInsensitive("Keyword"),
	)
)

type CypherQuery struct {
	Clauses []*Clause `parser:"@@*"`
}

type Clause struct {
	Match  *MatchClause  `parser:"  @@"`
	Create *CreateClause `parser:"| @@"`
	Return *ReturnClause `parser:"| @@"`
	Delete *DeleteClause `parser:"| @@"`
	Set    *SetClause    `parser:"| @@"`
}

type NodePattern struct {
	Variable   *string     `parser:"\"(\" @Ident?"`
	Labels     []string    `parser:"(\":\" @Ident)*"`
	Properties *Properties `parser:"@@? \")\""`
}

type Property struct {
	Key   string `parser:"@Ident \":\""`
	Value *Value `parser:"@@"`
}

type Properties struct {
	Props []*Property `parser:"\"{\" @@ (\",\" @@)* \"}\""`
}

type Value struct {
	String *string  `parser:"@String"`
	Number *float64 `parser:"| @Number"`
	Bool   *bool    `parser:"| @(\"TRUE\" | \"FALSE\")"`
}

type RelationshipPattern struct {
	LeftArrow  bool                 `parser:"@\"<\"? \"-\""`
	Details    *RelationshipDetails `parser:"(\"[\" @@ \"]\")? \"-\""`
	RightArrow bool                 `parser:"@\">\"?"`
}

type RelationshipDetails struct {
	Variable   *string     `parser:"@Ident?"`
	Types      []string    `parser:"(\":\" @Ident (\"|\" @Ident)*)?"`
	Properties *Properties `parser:"@@?"`
}

type PatternElement struct {
	Node  *NodePattern    `parser:"@@"`
	Chain []*PatternChain `parser:"@@*"`
}

type PatternChain struct {
	Relationship *RelationshipPattern `parser:"@@"`
	Node         *NodePattern         `parser:"@@"`
}

type MatchClause struct {
	Optional bool            `parser:"@\"OPTIONAL\"? \"MATCH\""`
	Pattern  *PatternElement `parser:"@@"`
}

type CreateClause struct {
	Pattern *PatternElement `parser:"\"CREATE\" @@"`
}

type ReturnClause struct {
	Items []string `parser:"\"RETURN\" @Ident (\",\" @Ident)*"`
}

type DeleteClause struct {
	Detach bool     `parser:"@\"DETACH\"? \"DELETE\""`
	Items  []string `parser:"@Ident (\",\" @Ident)*"`
}

type SetItem struct {
	Variable string `parser:"@Ident"`
	Property string `parser:"(\".\" @Ident)?"`
	Operator string `parser:"\"=\""`
	Value    *Value `parser:"@@"`
}

type SetClause struct {
	Items []*SetItem `parser:"\"SET\" @@ (\",\" @@)*"`
}
