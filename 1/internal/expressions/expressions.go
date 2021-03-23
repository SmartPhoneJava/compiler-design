package expressions

import (
	"gocompiler/internal/fsm"
	"gocompiler/internal/graph"
	"log"
	"strings"
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
				begin = i //+ 1
			}
			bracketOn++
		case ')':
			bracketOn--
			if bracketOn == 0 {
				lexemes = append(lexemes, str[begin:i+1])
				begin = i + 1
			}
		case '*':
			if i > 0 && str[i-1] == ')' && bracketOn == 0 {
				begin = i + 1
			}
			lexemes[len(lexemes)-1] += "*"
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

func removeBrackets(s string) string {
	if len(s) == 0 || s[0] != '(' {
		return s
	}
	var countBrackets int
	for i, symbol := range s {
		if symbol == '(' {
			countBrackets++
		} else if symbol == ')' {
			countBrackets--
			if countBrackets == 0 {
				if i == len(s)-1 { // последняя скобка первого уровня
					return s[1 : len(s)-1]
				}
				return s
			}
		}
	}

	return s
}

const (
	NoStar = iota
	StarNoBrackets
	StarBrackets
)

// RemoveZw убрать звездочку и вернуть флаг произошедшего события
func removeZw(s string) (string, int) {
	var oldS = ""
	for len(s) != len(oldS) {
		oldS = s
		s = removeBrackets(s)
	}
	if len(s) == 0 || s[len(s)-1] != '*' {
		return s, NoStar
	}
	var countBrackets int
	for _, s := range s {
		if s == '(' || s == ')' {
			countBrackets++
		}
	}
	// скобок нет
	if countBrackets == 0 {
		if len(s) > 0 && s[len(s)-1] == '*' {
			return s[:len(s)-1], StarNoBrackets // убираем *
		}
	} else if len(s)-len(removeBrackets(s[:len(s)-1])) == 3 {
		// есть внешняя скобка
		return s[1 : len(s)-2], StarBrackets // убираем скобки и *

	}
	return s, NoStar
}

// ToGraph привести к НКА
// АХО А.В, ЛАМ М.С., СЕТИ Р., УЛЬМАН Дж.Д. Компиляторы: принципы, технологии и инструменты. – М.:Вильямс, 2008.
// Алгоритм 3.23 Алгоритм Мак-Нотона-Ямады-Томпсона, стр. 213
func (str *RW) ToENKA() *fsm.FSM {
	var kda = &fsm.FSM{Graph: graph.NewGraph()}
	kda.AddEdge(&graph.Edge{
		From:   "q0",
		To:     "q1",
		Weight: string(*str),
	})
	var maxTry = 100
	var changes = 1
	for changes > 0 && maxTry > 0 {
		arr := kda.Edges
		changes = 0

		for _, edge := range arr {
			weight := strings.TrimSpace(edge.Weight)
			if len(weight) == 1 {
				continue
			}
			weight, changed := removeZw(weight)
			if changed != 0 {
				edge = kda.EpsilonEdge(edge, weight)
			}

			ew := RW(weight)
			rws := (&ew).Unions()
			kda.MultiplyEdge(edge, rws.toString()...)
			changes += len(rws)

			rws = (&ew).Concatenations()
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
	kda.SetFirstLast([]string{"q0"}, []string{"q1"})
	return kda
}

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
