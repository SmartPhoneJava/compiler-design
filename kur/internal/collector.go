package internal

import (
	"fmt"
	"kurs/parser"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const (
	NotLocal = iota
	MaybeLocalVar
	LocalFunc
	LocalVar
)

type InfoCollector struct {
	origin string
	rows   []string

	*parser.BaseLuaListener

	Funcs  *Funcs
	Tables *Tables

	candidateVar    string
	candidate       string
	candidateFindIn string

	expression string

	tableFunc string

	localStatus uint

	isReturning bool
}

func NewInfoCollector(origin string) *InfoCollector {
	var (
		mainFunc  = NewFunc(MainFunc)
		mainTable = NewTable(MainFunc)
	)
	return &InfoCollector{
		BaseLuaListener: &parser.BaseLuaListener{},
		Funcs: &Funcs{
			AllFuncs: map[string]*Func{
				mainFunc.Name: mainFunc,
			},
			callStack: []*Func{
				mainFunc,
			},
		},
		Tables: &Tables{
			Tables: map[string]*Table{
				mainFunc.Name: mainTable,
			},
			callStack: []*Table{
				mainTable,
			},
			implicitIndex: make(map[int]int),
		},
		origin: origin,
		rows:   strings.Split(origin, "\n"),
	}
}

type AntlrNode interface {
	GetText() string
	GetStart() antlr.Token
	GetStop() antlr.Token
}

func (s InfoCollector) GetText(ctx AntlrNode) string {
	var str string
	start := ctx.GetStart()
	end := ctx.GetStop()
	for i := start.GetLine(); i < end.GetLine()+1; i++ {
		if i == start.GetLine() {
			str += s.rows[i-1][start.GetColumn():]
		} else if i == end.GetLine() {
			str += s.rows[i-1][:end.GetColumn()]
		} else {
			str += s.rows[i-1]
		}
	}

	var text = ctx.GetText()
	words := strings.Split(str, " ")
	var wordsNum = 0
	var summSymbols = 0
	for summSymbols < len(text) && wordsNum < len(words) {
		summSymbols += len(words[wordsNum])
		if summSymbols > len(text) {
			words[wordsNum] = words[wordsNum][:len(words[wordsNum])-(summSymbols-len(text))]
		}
		wordsNum += 1
	}
	return strings.Join(words[:wordsNum], " ")
}

// Example
func (s *InfoCollector) VisitTerminal(node antlr.TerminalNode) {
	//fmt.Printf("term %v\n", node.GetText())
}

func (s *InfoCollector) EnterEveryRule(c antlr.ParserRuleContext) {
	fmt.Printf("rule %v %v\n", c.GetText(), c.GetRuleIndex())
}

func (s *InfoCollector) ExitEveryRule(c antlr.ParserRuleContext) {
	//fmt.Printf("aaa %v\n", c.GetText())
}
