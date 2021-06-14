package internal

import (
	"kurs/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const (
	NotLocal = iota
	MaybeLocalVar
	LocalFunc
	LocalVar
)

type InfoCollector struct {
	*parser.BaseLuaListener

	Funcs  *Funcs
	Tables *Tables

	candidateVar    string
	candidate       string
	candidateFindIn string

	localStatus uint

	isReturning bool

	ExitStatCallback func()
}

func NewInfoCollector() *InfoCollector {
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
		},
	}
}

// Example
func (s *InfoCollector) VisitTerminal(node antlr.TerminalNode) {
	//fmt.Printf("term %v\n", node.GetText())
}

func (s *InfoCollector) EnterEveryRule(c antlr.ParserRuleContext) {
	//fmt.Printf("rule %v %v\n", c.GetText(), c.GetRuleIndex())
}

func (s *InfoCollector) ExitEveryRule(c antlr.ParserRuleContext) {
	//fmt.Printf("aaa %v\n", c.GetText())
}
