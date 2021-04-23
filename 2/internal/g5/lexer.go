package g5

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

type Node struct {
	Value       string
	Parent      *Node
	ParentValue string
	Type        string
}

type ResT struct {
	Resolver Resolver
	Rule     Rule
	Symbols  []string

	Symbol Symbol
	Parent *ResT
}
type ResTs []ResT

func (lexer Lexer) ValidateDebug(text string, speed time.Duration) error {
	var (
		rows         = strings.Split(strings.ReplaceAll(text, "\n", ""), " ")
		comprassions int
		success      bool
	)
	var rules ResTs
	_, err := lexer.Start.GoTo(rows, 0, true, &rules, &success, &comprassions, speed)

	goterm.Println("Comprassions: ", comprassions)
	goterm.Flush()

	var nodes = []*Node{}
	m := make(map[string]*Node, 0)
	for _, r := range rules {
		var right string
		for _, s := range r.Rule.Symbols {
			right += " " + s.Value
		}
		log.Printf("%s->%s", r.Resolver.Symbol, right)
	}
	// var nodes = []*Node{}
	// m := make(map[string]*Node, 0)
	//m[lexer.Start.Symbol] = startNode

	//var goTextI = 0

	for i := len(rules) - 1; i >= 0; i-- {
		parent, ok := m[rules[i].Resolver.Symbol]
		if !ok {
			parent = &Node{
				Value: fmt.Sprintf("%d.", i) + rules[i].Resolver.Symbol,
				Type:  NonTerm,
			}
			nodes = append(nodes, parent)
		}
		for j, s := range rules[i].Symbols {
			var node = &Node{
				Value:       fmt.Sprintf("%d.%d... %s", i, j, s),
				Parent:      parent,
				ParentValue: parent.Value,
				//Type:        rules[i].Rule.Symbols[j].Type,
			}
			log.Println("node.Value", node.Value)
			// if rules[i].Rule.Symbols[j].Type == Term {
			// 	node.Value = rows[goTextI]
			// } else {
			// 	m[rules[i].Rule.Symbols[j].Value] = node
			// }
			if len(rules[i].Rule.Symbols) > j {
				if rules[i].Rule.Symbols[j].Type == NonTerm {
					m[rules[i].Rule.Symbols[j].Value] = node
					node.Value = fmt.Sprintf("%d.%d... %s", i, j, rules[i].Rule.Symbols[j].Value)
				}

				node.Type = rules[i].Rule.Symbols[j].Type
				if node.Type == Reserved {
					node.Type = Term
				}
			}

			//node.Value = s
			nodes = append(nodes, node)
		}
		delete(m, rules[i].Resolver.Symbol)
	}

	for _, node := range nodes {
		if node.Parent == nil {
			v := strings.Split(node.Value, "...")
			if len(v) == 2 {
				node.Parent = m[v[1]]
			}
		}

	}

	MustVisualize(nodes, "assets", "hello.dot")
	return err
}

func (lexer Lexer) Validate(text string, isDebug bool) error {
	var t time.Duration
	if isDebug {
		t = time.Second
	}
	return lexer.ValidateDebug(text, t)
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
