package g5

import (
	"fmt"
	"log"
	"strings"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

type Lexer struct {
	// ключ - нетерм, значение - распознаватель
	NonTerms map[string]*Resolver
	// Термы в мапе для быстрого поиска
	Terms map[string]bool
	Start *Resolver
}

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
	Value  string
	IsTerm bool
}

type Symbols []Symbol

// func (l Lexer) ToCFR() internal.CFR {
// 	// сделать приведение туда и обратно через хэши
// 	var (
// 		// N — конечный алфавит нетерминальных символов
// 		n []string
// 		// T —  конечный алфавит терминальных символов
// 		t []string
// 		// P — конечное множество правил порождения
// 		p internal.Rules
// 		// S — начальный нетерминал грамматики G
// 		s []string
// 	)
// 	for _, nonterm := range cg.Nonterms {
// 		n = append(n, nonterm.Name)
// 	}
// 	for _, term := range cg.Terms {
// 		t = append(t, term.Name)
// 	}
// 	for _, start := range cg.Start {
// 		s = append(s, start.Name)
// 	}
// 	for _, rule := range cg.Rules {
// 		var rightPart string
// 		for _, right := range rule.RightSide {
// 			rightPart += right.Name
// 		}
// 		p = append(p, internal.Rule{
// 			From: rule.LeftSide.Name,
// 			To:   rightPart,
// 		})
// 	}
// 	println("P", len(p), len(cg.Rules))
// 	return internal.CFR{
// 		N: n,
// 		S: s,
// 		P: p,
// 		T: t,
// 	}
// }

func (lexer Lexer) Validate(text string, isDebug bool) error {
	var (
		rows         = strings.Split(strings.Replace(text, "\n", "", 0), " ")
		comprassions int
		success      bool
	)
	_, err := lexer.Start.GoTo(rows, 0, isDebug, &success, &comprassions)
	if isDebug {
		goterm.Println("Comprassions: ", comprassions)
		goterm.Flush()
	}
	return err

}

// !! вынести в отдельный тип __IDENT и иже с ними
// Print - распечатать грамматику
func (lexer Lexer) Print(text string) {
	log.Println(text)

	color.Cyan("Список нетермов: \n")
	for nt := range lexer.NonTerms {
		fmt.Printf("%s\n", nt)
	}

	color.Cyan("Список термов: \n")
	for nt := range lexer.Terms {
		fmt.Printf("%s\n", nt)
	}

	color.Cyan("Стартовый нетерм: \n")
	fmt.Println(lexer.Start.Symbol)

	//sort.Sort(cfr.P)

	color.Cyan("Набор правил: \n")
	for nt, res := range lexer.NonTerms {
		for _, rule := range res.Rules {
			var right string
			for _, s := range rule.Symbols {
				if s.IsTerm {
					right += color.GreenString(s.Value)
				} else {
					right += " <" + color.YellowString(s.Value) + "> "
				}
			}
			if lexer.Start.Symbol == nt {
				nt = color.RedString(nt)
			} else {
				nt = color.YellowString(nt)
			}
			fmt.Printf("%s → %s\n", nt, right)
		}
	}
}

const (
	IndexStatusFound = iota + 1
	IndexStatusNotFound
)

func (lexer Lexer) PrintState(
	text []string,
	indexes map[int]int,
	currentResolverSymbol string,
	currentTextIndex, currentRuleI, currentSymbolI int,
) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)

	goterm.Println("Рассматриваемые индексы:")
	for index, status := range indexes {
		var indexS = fmt.Sprintf("%d", index)
		if index == currentTextIndex {
			indexS = color.CyanString(indexS)
		} else {
			if status == IndexStatusFound {
				indexS = color.GreenString(indexS)
			} else if status == IndexStatusNotFound {
				indexS = color.RedString(indexS)
			}
		}
		goterm.Println(indexS + " ")
	}

	goterm.Println("Рассматриваемый код:")
	for i, row := range text {
		//goterm.Println("currentTextIndex", i, row, currentTextIndex)
		if i == currentTextIndex {
			row = color.CyanString(row)
		}
		goterm.Print(row + " ")
	}

	//sort.Sort(cfr.P)

	goterm.Println("Разбор:", currentResolverSymbol)
	for nt, res := range lexer.NonTerms {
		for i, rule := range res.Rules {
			var right string
			for j, s := range rule.Symbols {
				if i == currentRuleI && j == currentSymbolI && res.Symbol == currentResolverSymbol {
					right += color.RedString(s.Value)
				} else if s.IsTerm {
					right += color.GreenString(s.Value)
				} else {
					right += " <" + color.YellowString(s.Value) + "> "
				}
			}
			if i == currentRuleI {
				if res.Symbol == currentResolverSymbol {
					nt = color.RedString(nt)
					goterm.Printf("%s → %s\n", nt, right)
				} else {
					nt = color.YellowString(nt)
				}

			} else {
				nt = color.YellowString(nt)
			}

		}
	}

	goterm.Flush()
	//time.Sleep(time.Second / 5)
}
