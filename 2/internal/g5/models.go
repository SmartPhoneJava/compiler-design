package g5

type Lexer struct {
	// ключ - нетерм, значение - распознаватель
	NonTerms map[string]*Resolver
	// Термы в мапе для быстрого поиска
	Terms map[string]bool
	Start *Resolver
}

// Узел дерева AST
type AST struct {
	ID     uint32
	Parent *AST
	Value  string
	Type   int
}

var ID_AST = 0

type Resolver struct {
	Rules
	Symbol string
	Lexer  *Lexer
}

type Rule struct {
	Symbols Symbols
}

type Rules []Rule

type Symbol struct {
	Value string
	Type  string
}

type Symbols []Symbol
