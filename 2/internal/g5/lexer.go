package g5

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

func (lexer Lexer) ValidateDebug(text string, speed time.Duration) error {
	var (
		rows         = strings.Split(strings.ReplaceAll(text, "\n", ""), " ")
		comprassions int
		success      bool
	)
	_, err := lexer.Start.GoTo(rows, 0, true, &success, &comprassions, speed)

	goterm.Println("Comprassions: ", comprassions)
	goterm.Flush()

	return err
}

func (lexer Lexer) Validate(text string, isDebug bool) error {
	var (
		rows         = strings.Split(strings.ReplaceAll(text, "\n", ""), " ")
		comprassions int
		success      bool
	)
	_, err := lexer.Start.GoTo(rows, 0, isDebug, &success, &comprassions, 0)
	if isDebug {
		goterm.Println("Comprassions: ", comprassions)
		goterm.Flush()
	}
	return err
}

func (lexer Lexer) ColorSymbol(s Symbol, right *string) {
	switch s.Type {
	case Term:
		*right += color.GreenString(s.Value)
	case NonTerm:
		*right += " <" + color.YellowString(s.Value) + "> "
	case Reserved:
		*right += color.HiMagentaString(s.Value)
	}
}

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

	color.Cyan("Набор правил: \n")
	for nt, res := range lexer.NonTerms {
		for _, rule := range res.Rules {
			var right string
			for _, s := range rule.Symbols {
				lexer.ColorSymbol(s, &right)
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
	speed time.Duration,
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
				} else {
					lexer.ColorSymbol(s, &right)
				}
			}
			if i == currentRuleI {
				if res.Symbol == currentResolverSymbol {
					nt = color.RedString(nt)
					goterm.Printf("%s → %s\n", nt, right)
				}
				// else {
				// 	nt = color.YellowString(nt)
				// }
			}
			// else {
			// 	nt = color.YellowString(nt)
			// }
		}
	}

	goterm.Flush()
	if speed != 0 {
		time.Sleep(speed)
	}
}
