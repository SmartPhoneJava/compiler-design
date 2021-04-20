package g5

import (
	"errors"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

// принимает на вход массив всех символов сканируемого текста
// и указатель на текущий символ
// возвращает мапка новых индексов и ошибку
// вовзращается именно мапка, чтобы не было повторений
func (r Resolver) GoTo(
	input []string,
	inputI int,
	isDebug bool,
	finishedSuccess *bool,
	comprassions *int,
) (map[int]interface{}, error) {
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

		for si, s := range rule.Symbols {
			var symbolI = make(map[int]int, 0)
			//log.Println("ruleI", si, s.Value, len(rule.Symbols), ruleI)
			for index := range ruleI {
				//	log.Println("symbolI", symbolI)
				if isDebug {
					r.Lexer.PrintState(input, ruleI, r.Symbol, index, ri, si)
				}
				if index >= len(input) {
					ruleI[index] = IndexStatusNotFound
					continue
				}
				*comprassions++
				if s.IsTerm {
					if s.Value == "__EMPTY" {
						//	log.Println("found", index)
						ruleI[index] = IndexStatusFound
						symbolI[index] = 0
						//	log.Println("want symbolI", symbolI)
						goterm.Println(color.GreenString("Терм найден"))
						goterm.Flush()
					}
					if input[index] == s.Value || s.Value == "__ANY" || s.Value == "__NUMBER" || s.Value == "__IDENT" {
						//	log.Println("found", index)
						ruleI[index] = IndexStatusFound
						symbolI[index+1] = 0
						//	log.Println("want symbolI", symbolI)
						goterm.Println(color.GreenString("Терм найден"))
						goterm.Flush()
					}
				} else {
					newResolver, ok := r.Lexer.NonTerms[s.Value]
					if !ok {
						ruleI[index] = IndexStatusNotFound
						continue
					}

					goterm.Println(color.GreenString("Переходим в новый нетерм: " + s.Value))
					goterm.Flush()
					lm, err := newResolver.GoTo(input, index, isDebug, finishedSuccess, comprassions)
					if *finishedSuccess {
						return nil, nil
					}
					//log.Println("returned", len(lm), err)
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
			ruleI = map[int]int{}
			//log.Println("symbolIsymbolI", symbolI)
			for i := range symbolI {
				ruleI[i] = 0
			}
		}
		if len(ruleI) > 0 {
			goterm.Println(color.GreenString("Правило подошло"))
		} else {
			goterm.Println(color.RedString("Правило не подошло"))
		}
		goterm.Flush()
		//log.Println("xit", len(ruleI))
		for i := range ruleI {
			newIndexes[i] = true
		}
	}

	for index := range newIndexes {
		if index > len(input) {
			delete(newIndexes, index)
		}
		if index == len(input) {
			*finishedSuccess = true
		}
	}

	//log.Println("return newIndexes", newIndexes)
	//goterm.Println("newIndexes", len(newIndexes))
	if len(newIndexes) == 0 {
		return newIndexes, errors.New("Грамматика покрывает не весь код")
	}

	return newIndexes, nil
}
