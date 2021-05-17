package opa

// operator_precedence_analyzer.go
import (
	"fmt"
	"lab2/internal/g5"
)

const (
	NO byte = iota
	EQ
	MORE
	LESS
)

// Matrix of operator precedence relations
type OperatorsMatrix map[string]map[string]byte

/*
func MakeMatrix(lexer g5.Lexer) OperatorsMatrix {
	var (
		AllTermsInNoneTerms   = MakeNT2T(lexer)
		leftNT, rightNT, some = MakeT2NT(lexer)
		matrix                = make(OperatorsMatrix)
	)
	log.Println("leftNT", leftNT)
	log.Println("rightNT", rightNT)
	log.Println("some", some)

	for row, rules := range rightNT {
		_, ok := matrix[row]
		if !ok {
			matrix[row] = make(map[string]byte)
		}
		for noneTerm := range rules {
			for column := range AllTermsInNoneTerms[noneTerm] {
				matrix[row][column] = MORE
				log.Println("row more", row, column)
			}
		}
	}
	for row, rules := range leftNT {
		_, ok := matrix[row]
		if !ok {
			matrix[row] = make(map[string]byte)
		}
		for noneTerm := range rules {
			for column := range AllTermsInNoneTerms[noneTerm] {
				matrix[row][column] = MORE
				log.Println("row less", row, column)
			}
		}
	}

	for row, rules := range some {
		_, ok := matrix[row]
		if !ok {
			matrix[row] = make(map[string]byte)
		}
		for column := range rules {
			matrix[row][column] = MORE
			log.Println("row equal", row, column)
		}
	}
	log.Println("matrix", matrix)
	return matrix
}

// Matrix of operator precedence relations
type OperatorsSet map[string]map[string]interface{}

type NoneTermsSet map[string]interface{}

// Из каких нетермов в какие термы можно попасть
func MakeNT2T(lexer g5.Lexer) OperatorsSet {
	var nt2nt = make(OperatorsSet)
	var symbolsSeen = make(NoneTermsSet)
	makeNT2NT(lexer, "", nt2nt, symbolsSeen, nil)
	return nt2nt
}

// какие термы в каких нетермах лежат
func MakeT2NT(lexer g5.Lexer) (OperatorsSet, OperatorsSet, OperatorsSet) {
	var (
		leftNT  = make(OperatorsSet)
		rightNT = make(OperatorsSet)
		rightT  = make(OperatorsSet)
	)
	for _, r := range lexer.NonTerms {
		for _, r := range r.Rules {
			var (
				savedTerms    []string
				savedNonTerms []string
			)
			for _, s := range r.Symbols {
				//log.Println("savedSymbols", s, savedSymbols)
				for _, savedS := range savedNonTerms {
					if s.Type != g5.NonTerm {
						_, ok := leftNT[savedS]
						if !ok {
							leftNT[savedS] = make(NoneTermsSet)
						}
						leftNT[savedS][s.Value] = nil
					}
				}
				for _, savedS := range savedTerms {
					if s.Type == g5.NonTerm {
						_, ok := rightNT[savedS]
						if !ok {
							rightNT[savedS] = make(NoneTermsSet)
						}
						rightNT[savedS][s.Value] = nil
					} else {
						_, ok := rightT[savedS]
						if !ok {
							rightT[savedS] = make(NoneTermsSet)
						}
						rightT[savedS][s.Value] = nil
					}
				}
				if s.Type != g5.NonTerm {
					savedTerms = append(savedTerms, s.Value)
				} else {
					savedNonTerms = append(savedNonTerms, s.Value)
				}
			}
		}
	}
	return leftNT, rightNT, rightT
}

func makeNT2NT(
	lexer g5.Lexer,
	currSymb string,
	nt2nt OperatorsSet,
	symbolsSeen NoneTermsSet,
	alreadySearching map[string]interface{},
) {
	for row, r := range lexer.NonTerms {
		if currSymb != "" && row != currSymb {
			continue
		}
		if currSymb == "" {
			alreadySearching = make(map[string]interface{})
		}
		_, ok := nt2nt[row]
		if !ok {
			nt2nt[row] = make(NoneTermsSet)
		}
		for _, rule := range r.Rules {
			for _, s := range rule.Symbols {
				if s.Type != g5.NonTerm {
					//log.Println("ssss", s.Type, s.Value, row, s.Value)
					nt2nt[row][s.Value] = nil
				} else {
					if s.Value == row {
						continue
					}
					_, seen := symbolsSeen[s.Value]

					if !seen {
						_, already := alreadySearching[s.Value]
						if !already {
							log.Println(s.Value)
							log.Println(alreadySearching[s.Value])
							alreadySearching[s.Value] = nil
							log.Println("s.Value", s.Value, row)
							makeNT2NT(lexer, s.Value, nt2nt, symbolsSeen, alreadySearching)
						}
					}

					for ts := range nt2nt[s.Value] {
						nt2nt[row][ts] = nil
					}
				}
			}
		}
		symbolsSeen[row] = true
	}
}
*/
/*

Для каждого нетерминального символа А ищем все правила, содержащие А в левой части.
Во множество L(А) включаем самый левый терминальный символ из правой части правил,
игнорируя нетерминальные символы, а во множество R(A) — самый крайний правый
терминальный символ из правой части правил. Переходим к шагу 2

*/

type LR map[string]map[g5.Symbol]interface{}

// Сделать копию L/R
func (lr LR) copy() LR {
	var copied = make(LR)
	for noneTerm, symbols := range lr {
		copied[noneTerm] = make(map[g5.Symbol]interface{})
		for symbol := range symbols {
			copied[noneTerm][symbol] = nil
		}
	}
	return copied
}

// обновить L/R - возвращает true, если итерация изменила содержимое L/R
func (lr LR) iterate() bool {
	// чтобы не итерироваться по новым элементам на той же итерации
	// итерируемся по копии, а добавляем в оригинал
	copied := lr.copy()

	// флаг было ли изменение на этой итерации
	var changed bool
	for noneTerm, symbols := range copied {
		for symbol := range symbols {
			if symbol.Type == g5.NonTerm {
				var before = len(lr[noneTerm])
				for innerSymbol := range copied[symbol.Value] {
					lr[noneTerm][innerSymbol] = nil
				}
				var after = len(lr[noneTerm])
				if before != after {
					changed = true
				}
			}
		}
	}
	return changed
}

// L
func makeMostLeft(lexer g5.Lexer) LR {
	// шаг 1
	var mostL = make(LR)
	for noneTerm, resolver := range lexer.NonTerms {
		mostL[noneTerm] = make(map[g5.Symbol]interface{})
		for _, r := range resolver.Rules {
			for _, symbol := range r.Symbols {
				mostL[noneTerm][symbol] = nil
				break
			}
		}
	}

	// шаг 2 + шаг 3
	// Итерируемся, пока что то меняется
	for mostL.iterate() {
	}
	return mostL
}

// R
func makeMostRight(lexer g5.Lexer) LR {
	// шаг 1
	var mostR = make(LR)
	for nt, resolver := range lexer.NonTerms {
		mostR[nt] = make(map[g5.Symbol]interface{})
		for _, r := range resolver.Rules {
			for i := len(r.Symbols) - 1; i >= 0; i-- {
				var s = r.Symbols[i]
				mostR[nt][s] = nil
				break
			}
		}
	}

	// шаг 2 + шаг 3
	// Итерируемся, пока что то меняется
	for mostR.iterate() {
	}
	return mostR
}

type LRt map[string]map[string]interface{}

func (lrt LRt) Println(text string) {
	for nt, symbols := range lrt {
		fmt.Printf("\n%s(%s):", text, nt)
		for symbol := range symbols {
			fmt.Printf("%s ", symbol)

		}
	}
}

// обновить Lt/Rt - возвращает true, если итерация изменила содержимое Lt/Rt
func (lrt LRt) iterate(lr LR) bool {
	// флаг было ли изменение на этой итерации
	var changed bool
	for noneTerm, symbols := range lr {
		for symbol := range symbols {
			if symbol.Type == g5.NonTerm {
				var before = len(lrt[noneTerm])
				for innerSymbol := range lrt[symbol.Value] {
					lrt[noneTerm][innerSymbol] = nil
				}
				var after = len(lrt[noneTerm])
				if before != after {
					changed = true
				}
			}
		}
	}
	return changed
}

// Lt
func MakeMostLeftTerm(lexer g5.Lexer) LRt {
	var mostLt = make(LRt)
	for nt, resolver := range lexer.NonTerms {
		mostLt[nt] = make(map[string]interface{})
		for _, r := range resolver.Rules {
			for _, s := range r.Symbols {
				if s.Type != g5.NonTerm {
					mostLt[nt][s.Value] = nil
					break
				}
			}
		}
	}

	var mostL = makeMostLeft(lexer)
	// шаг 2 + шаг 3
	// Итерируемся, пока что то меняется
	for mostLt.iterate(mostL) {
	}
	return mostLt
}

// Rt
func MakeMostRightTerm(lexer g5.Lexer) LRt {
	var mostRt = make(LRt)
	for nt, resolver := range lexer.NonTerms {
		mostRt[nt] = make(map[string]interface{})
		for _, r := range resolver.Rules {
			for i := len(r.Symbols) - 1; i >= 0; i-- {
				var s = r.Symbols[i]
				if s.Type != g5.NonTerm {
					mostRt[nt][s.Value] = nil
					break
				}
			}
		}
	}

	var mostR = makeMostRight(lexer)
	// шаг 2 + шаг 3
	// Итерируемся, пока что то меняется
	for mostRt.iterate(mostR) {
	}
	return mostRt
}

func MakeMatrixV2(lexer g5.Lexer) OperatorsMatrix {
	var (
		left  = MakeMostLeftTerm(lexer)
		right = MakeMostRightTerm(lexer)
	)

	var matrix = make(OperatorsMatrix)
	for _, res := range lexer.NonTerms {
		for _, rule := range res.Rules {
			for i, s := range rule.Symbols {
				if s.Type != g5.NonTerm {
					_, ok := matrix[s.Value]
					if !ok {
						matrix[s.Value] = make(map[string]byte)
					}
					if i != 0 {
						prevValue := rule.Symbols[i-1].Value
						if rule.Symbols[i-1].Type == g5.NonTerm {
							for symbol := range right[prevValue] {
								_, ok = matrix[symbol]
								if !ok {
									matrix[symbol] = make(map[string]byte)
								}
								matrix[symbol][s.Value] = MORE
							}
						}
					}
					if i != len(rule.Symbols)-1 {
						nextValue := rule.Symbols[i+1].Value

						if rule.Symbols[i+1].Type != g5.NonTerm {
							matrix[s.Value][nextValue] = EQ
						} else {
							for symbol := range left[nextValue] {
								matrix[s.Value][symbol] = LESS
							}
							if i+2 < len(rule.Symbols) && rule.Symbols[i+2].Type != g5.NonTerm {
								matrix[s.Value][rule.Symbols[i+2].Value] = EQ
							}
						}
					}
				}
			}
		}
	}

	var start = lexer.Start.Symbol
	for symbol := range right[start] {
		_, ok := matrix[symbol]
		if !ok {
			matrix[symbol] = make(map[string]byte)
		}
		matrix[symbol]["⏊"] = MORE
	}
	for symbol := range left[start] {
		_, ok := matrix["⏊"]
		if !ok {
			matrix["⏊"] = make(map[string]byte)
		}
		matrix["⏊"][symbol] = LESS
	}

	return matrix
}

func (matrix OperatorsMatrix) Println() {
	//var arr = make([]string, 0, len(matrix))
	var arr = []string{"(", "a", "*", "+", ")", "⏊"}
	var elemIndex = map[string]int{
		"(": 0,
		"a": 1,
		"*": 2,
		"+": 3,
		")": 4,
		"⏊": 5,
	}

	// for v := range matrix {
	// 	elemIndex[v] = len(arr)
	// 	arr = append(arr, v)
	// }

	var realMatrix = make([][]string, len(arr))
	for i := range realMatrix {
		realMatrix[i] = make([]string, len(arr))
	}

	for c, terms := range matrix {
		for r, symbol := range terms {
			var text = ""
			switch symbol {
			case EQ:
				text = "="
			case MORE:
				text = "▶"
			case LESS:
				text = "◀"
			}
			realMatrix[elemIndex[c]][elemIndex[r]] = text
		}
	}

	var printHorLine = func() {

		fmt.Printf("\n")
		for range arr {
			fmt.Printf("------|")
		}
		fmt.Printf("------|")
		fmt.Printf("\n")
	}

	fmt.Printf("\n\n%6s|", "")
	for _, s := range arr {
		fmt.Printf("%6s|", center(s, 6))
	}
	printHorLine()
	for i, row := range realMatrix {
		fmt.Printf("%6s|", center(arr[i], 6))
		for _, v := range row {
			fmt.Printf("%6s|", center(v, 6))
		}
		printHorLine()
	}
}

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}
