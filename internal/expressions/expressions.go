package expressions

import (
	"gocompiler/internal/graph"
	"log"
)

// RW регулярное выражение
type RW string

type RWS []RW

func (rws RWS) toString() []string {
	var arr = make([]string, 0)
	for _, rw := range rws {
		arr = append(arr, string(rw))
	}
	return arr
}

func NewRW(str string) RW {
	return RW(str)
}

// !! simplify - упростить регулярку
// !! пока не реализовано
func (*RW) simplify() {

}

// Concatenations Разбить на список конкатенаций
func (str RW) Concatenations() []RW {
	var lexemes = make([]RW, 0)
	var (
		bracketOn, begin int
	)
	for i, r := range str {
		switch r {
		case ' ':
			continue
		case '(':
			if bracketOn == 0 {
				begin = i + 1
			}
			bracketOn++
		case ')':
			bracketOn--
			if bracketOn == 0 {
				lexemes = append(lexemes, str[begin:i])
				begin = i + 1
			}
		case '|':
			if bracketOn == 0 {
				return nil
			}
		default:
			if bracketOn == 0 {
				lexemes = append(lexemes, RW(r))
			}
		}
	}
	return lexemes
}

func (str *RW) trumBrackets() {
	if (*str)[0] != '(' {
		return
	}
	var (
		bracketOn int
	)
	for i, r := range *str {
		switch r {
		case '(':
			bracketOn++
		case ')':
			bracketOn--
			if bracketOn == 0 {
				if i == len(*str)-1 {
					*str = (*str)[1 : len(*str)-1]
					return
				}
			}
		}
	}
}

// Unions разбить на множество объединений
func (str *RW) Unions() RWS {
	var lexemes = make(RWS, 0)
	var (
		bracketOn, begin int
	)
	for i, r := range *str {
		switch r {
		case '(':
			bracketOn++
		case ')':
			bracketOn--
		case '|':
			if bracketOn == 0 {
				lexemes = append(lexemes, (*str)[begin:i])
				begin = i + 1
			}
		}
	}
	if begin > 0 && begin < len(*str) {
		lexemes = append(lexemes, (*str)[begin:])
	}
	return lexemes
}

// ToGraph привести к виду графа
func (str *RW) ToGraph() *graph.Graph {
	var kda = graph.NewGraph()
	kda.AddEdge(&graph.Edge{
		From:   "q0",
		To:     "q1",
		Weight: string(*str),
	})
	var maxTry = 1000
	for maxTry > 0 {
		arr := kda.Edges
		var changes int

		for _, edge := range arr {
			if len(edge.Weight) == 1 {
				continue
			}

			ew := RW(edge.Weight)
			rws := (&ew).Unions()
			kda.MultiplyEdge(edge, rws.toString()...)
			changes += len(rws)

			rws = RW(edge.Weight).Concatenations()
			kda.SplitEdge(edge, rws.toString()...)
			changes += len(rws)
		}
		if changes == 0 {
			break
		}
		maxTry--
	}
	if maxTry == 0 {
		log.Println("error happened")
	}
	return kda
}

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
