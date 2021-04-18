package g5

import (
	"errors"

	"github.com/buger/goterm"
)

// принимает на вход массив всех символов сканируемого текста
// и указатель на текущий символ
// возвращает мапка новых индексов и ошибку
// вовзращается именно мапка, чтобы не было повторений
func (r Resolver) GoTo(input []string, inputI int) (map[int]interface{}, error) {
	if len(input) <= inputI {
		return nil, errors.New("указатель за пределы массива")
	}

	// var currInput = input[inputI]
	// var foundRules = make(Rules, 0)
	// for _, rule := range r.Rules {
	// 	// Если первый символ правой части совпал, значит правило подходит
	// 	for _, s := range rule.Symbols {
	// 		if currInput == s.Value {
	// 			foundRules = append(foundRules, rule)
	// 		}
	// 		break
	// 	}
	// }

	// Мапка новых индексов
	var newIndexes = make(map[int]interface{}, 0)
	for ri, rule := range r.Rules {
		//var allSymbolsFound = true

		var ruleI = make(map[int]int, 0)
		ruleI[inputI] = 0
		var symbolI = make(map[int]int, 0)
		for si, s := range rule.Symbols {
			for index := range ruleI {
				r.Lexer.PrintState(input, ruleI, r.Symbol, index, ri, si)
				if index >= len(input) {
					ruleI[index] = IndexStatusNotFound
					continue
				}
				if s.IsTerm {
					if input[index] == s.Value {
						ruleI[index] = IndexStatusFound
						symbolI[index+1] = 0
					}
				} else {
					newResolver, ok := r.Lexer.NonTerms[s.Value]
					if !ok {
						ruleI[index] = IndexStatusNotFound
						continue
					}

					lm, err := newResolver.GoTo(input, index)
					if err != nil {
						ruleI[index] = IndexStatusNotFound
						continue
					}
					if len(lm) > 0 {
						ruleI[index] = IndexStatusFound
						for newI := range lm {
							symbolI[newI] = 0
						}
					} else {
						ruleI[index] = IndexStatusNotFound
					}
				}
			}
			ruleI = make(map[int]int, 0)
			for i := range symbolI {
				ruleI[i] = 0
			}
		}
		for i := range ruleI {
			newIndexes[i] = true
		}
	}
	goterm.Println("newIndexes", len(newIndexes))
	if len(newIndexes) == 0 {
		return newIndexes, errors.New("Грамматика покрывает не весь код")
	}
	for index := range newIndexes {
		if index == len(input) {
			return newIndexes, nil
		}
	}
	return newIndexes, errors.New("Найдены лишние символы")
}
