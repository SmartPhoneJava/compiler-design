package parsing

import (
	"lab2/internal"
)

//грамматика это четверка
type ContextGrammar struct {
	Name string `json:"-name,omitempty"`
	//нетерминалы
	Nonterms []Symbol
	//терминалы
	Terms []Symbol
	//правила
	Rules []Rule `json:"productions,omitempty"`
	//стартовые символы
	Start []Symbol `json:"startsymbol,omitempty"`
}

type SymbolStruct struct {
	sm Symbol
}
type Symbol struct {
	//сама строчка с названием символа
	Name string `json:"-name,omitempty"`
	//строчка с написанием символа
	Spell string `json:"-spell,omitempty"`
	Type  string
}

type Rule struct {
	LeftSide  Symbol   `json:"lhs,omitempty"`
	RightSide []Symbol `json:"rhs,omitempty"`
}

type GrammarStruct struct {
	Gm ContextGrammar `json:"grammar"`
}

func (cg ContextGrammar) ToInternal() internal.CFR {
	var (
		// N — конечный алфавит нетерминальных символов
		n []string
		// T —  конечный алфавит терминальных символов
		t []string
		// P — конечное множество правил порождения
		p internal.Rules
		// S — начальный нетерминал грамматики G
		s []string
	)
	for _, nonterm := range cg.Nonterms {
		n = append(n, nonterm.Name)
	}
	for _, term := range cg.Terms {
		t = append(t, term.Name)
	}
	for _, start := range cg.Start {
		s = append(s, start.Name)
	}
	for _, rule := range cg.Rules {
		var rightPart string
		for _, right := range rule.RightSide {
			rightPart += right.Name
		}
		p = append(p, internal.Rule{
			From: rule.LeftSide.Name,
			To:   rightPart,
		})
	}
	println("P", len(p), len(cg.Rules))
	return internal.CFR{
		N: n,
		S: s,
		P: p,
		T: t,
	}
}

func (cg *ContextGrammar) FromInternal(cfr internal.CFR) {
	for _, nonterm := range cfr.N {
		cg.Nonterms = append(cg.Nonterms, Symbol{
			Name:  nonterm,
			Spell: nonterm,
			Type:  "nonterm",
		})
	}
	for _, term := range cfr.T {
		cg.Terms = append(cg.Terms, Symbol{
			Name:  term,
			Spell: term,
			Type:  "term",
		})
	}
	for _, start := range cfr.S {
		cg.Start = append(cg.Start, Symbol{
			Name:  start,
			Spell: start,
			Type:  "nonterm",
		})
	}
	for _, rule := range cfr.P {
		symbols := cfr.TermsAndNonTerms(rule.To)
		var corSymbols []Symbol
		for _, symbol := range symbols {
			corSymbols = append(corSymbols, Symbol{
				Spell: symbol.Spell,
				Name:  symbol.Spell,
				Type:  symbol.Type,
			})
		}
		cg.Rules = append(cg.Rules, Rule{
			LeftSide: Symbol{
				Name:  rule.From,
				Spell: rule.From,
				Type:  "nonterm",
			},
			RightSide: corSymbols,
		})
	}
}
