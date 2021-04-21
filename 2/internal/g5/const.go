package g5

// Тип символ
const (
	Term     = "term"
	NonTerm  = "nonterm"
	Reserved = "reserved"
)

// служебные термы
const (
	// Любой набор символов
	TermAny = "__ANY"
	// Любое число
	TermNumber = "__NUMBER"
	// Любое идентификатор
	TermIDENT = "__IDENT"
)
