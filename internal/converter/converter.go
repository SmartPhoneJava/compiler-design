package converter

import (
	"gocompiler/internal/expressions"
	"gocompiler/internal/fsm"
)

// ExpressionToNKA привести регулярное выражение к
// 	непрерывному конечному автомату
func ExpressionToNKA(str *expressions.RW) *fsm.FSM {
	gr := str.ToGraph()
	return &fsm.FSM{gr}
}
