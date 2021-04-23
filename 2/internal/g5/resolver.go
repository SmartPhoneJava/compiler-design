package g5

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

func (r Resolver) handleReserved(
	input []string,
	value string,
	index int,
	current, next map[int]int,
	isDebug bool,
) {

	var ok bool
	switch value {
	case TermAny:
		ok = regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(input[index])
	case TermNumber:
		ok = regexp.MustCompile(`^[0-9]+$`).MatchString(input[index])
	case TermIDENT:
		ok = regexp.MustCompile(`^[A-Za-z]+$`).MatchString(input[index])
	}
	if !ok {
		return
	}
	current[index] = IndexStatusFound
	next[index+1] = 0
	if isDebug {
		goterm.Println(color.GreenString("Терм найден"))
		goterm.Flush()
	}
}

func (r Resolver) handleTerm(
	input []string,
	value string,
	index int,
	current, next map[int]int,
	isDebug bool,
) {
	if input[index] == value {
		current[index] = IndexStatusFound
		next[index+1] = 0
		if isDebug {
			goterm.Println(color.GreenString("Терм найден"))
			goterm.Flush()
		}
	}
}

func (r Resolver) handleNonTerm(
	input []string,
	value string,
	index int,
	current, next map[int]int,
	isDebug bool,
	getItRules *ResTs,
	finishedSuccess *bool,
	comprassions *int,
	speed time.Duration,
) {
	newResolver, ok := r.Lexer.NonTerms[value]
	if !ok {
		current[index] = IndexStatusNotFound
		return
	}

	if isDebug {
		goterm.Println(color.GreenString("Переходим в новый нетерм: " + value))
		goterm.Flush()
	}
	lm, err := newResolver.GoTo(
		input, index, isDebug, getItRules,
		finishedSuccess, comprassions, speed,
	)
	if *finishedSuccess {
		return
	}
	if err != nil {
		current[index] = IndexStatusNotFound
		return
	}
	if len(lm) > 0 {
		current[index] = IndexStatusFound
		for newI := range lm {
			next[newI] = 0
		}
	} else {
		current[index] = IndexStatusNotFound
	}
}

// принимает на вход массив всех символов сканируемого текста
// и указатель на текущий символ
// возвращает мапка новых индексов и ошибку
// вовзращается именно мапка, чтобы не было повторений
func (r Resolver) GoTo(
	input []string,
	inputI int,
	isDebug bool,
	getItRules *ResTs,
	finishedSuccess *bool,
	comprassions *int,
	speed time.Duration,
) (map[int]interface{}, error) {
	if len(input) <= inputI {
		return nil, errors.New("указатель за пределы массива")
	}

	// Мапка новых индексов
	var newIndexes = make(map[int]interface{}, 0)
	for ri, rule := range r.Rules {
		var ruleI = make(map[int]int, 0)
		ruleI[inputI] = 0

		var initial = make(ResTs, len(*getItRules))
		copy(initial, *getItRules)

		var symbols = []string{}
		for si, s := range rule.Symbols {
			var symbolI = make(map[int]int, 0)
			for index := range ruleI {
				if isDebug {
					r.Lexer.PrintState(
						input, ruleI, r.Symbol,
						index, ri, si, speed,
					)
				}
				if index >= len(input) {
					ruleI[index] = IndexStatusNotFound
					continue
				}
				*comprassions++
				switch s.Type {
				case Reserved:
					symbols = append(symbols, input[index])
					r.handleReserved(
						input, s.Value, index, ruleI,
						symbolI, isDebug,
					)
				case Term:
					symbols = append(symbols, input[index])
					r.handleTerm(
						input, s.Value, index, ruleI,
						symbolI, isDebug,
					)
				case NonTerm:
					symbols = append(symbols, input[index])
					r.handleNonTerm(
						input, s.Value, index, ruleI,
						symbolI, isDebug, getItRules,
						finishedSuccess, comprassions, speed,
					)
					if *finishedSuccess {
						*getItRules = append(*getItRules, ResT{
							Resolver: r,
							Rule:     rule,
						})
						return nil, nil
					}
				}
			}
			ruleI = map[int]int{}
			for i := range symbolI {
				ruleI[i] = 0
			}
		}
		if isDebug {
			var text string
			if len(ruleI) > 0 {
				*getItRules = append(*getItRules, ResT{
					Resolver: r,
					Rule:     rule,
					Symbols:  symbols,
				})
				text = color.GreenString("Правило подошло")
			} else {
				text = color.RedString("Правило не подошло")
				*getItRules = initial
			}
			goterm.Println(text)
			goterm.Flush()
		}

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

	if len(newIndexes) == 0 {
		return newIndexes, errors.New("Грамматика покрывает не весь код")
	}

	return newIndexes, nil
}

// !! написать тесты

// Множество строк
type StringSet map[string]interface{}

// Множество нетермов, где каждому соответсвует множество строк
type NoneTermsSet map[string]StringSet

func (set NoneTermsSet) Println(text string) {
	log.Println(text)
	for nt, ts := range set {
		var right string
		for term := range ts {
			right += term + " "
		}
		log.Printf("%s -> { %s}\n", nt, right)
	}
}

// https://neerc.ifmo.ru/wiki/index.php?title=Построение_FIRST_и_FOLLOW
func (lexer Lexer) ConstructFirst() NoneTermsSet {
	var first = make(NoneTermsSet, 0)
	for nt := range lexer.NonTerms {
		first[nt] = make(StringSet, 0)
	}
	var changed = true
	for changed {
		changed = false
		for noneTerm, r := range lexer.NonTerms {
			var beforeLen = len(first[noneTerm])
			for _, rule := range r.Rules {
				for _, s := range rule.Symbols {
					if s.Type == Term {
						first[noneTerm][s.Value] = nil
					} else if s.Type == NonTerm {
						for term := range first[s.Value] {
							first[noneTerm][term] = nil
						}
					}
				}
			}
			var afterLen = len(first[noneTerm])
			if beforeLen != afterLen {
				changed = true
			}
		}
	}
	return first
}

// https://neerc.ifmo.ru/wiki/index.php?title=Построение_FIRST_и_FOLLOW
func (lexer Lexer) ConstructFollow(first NoneTermsSet) NoneTermsSet {
	var follow = make(NoneTermsSet, 0)
	for nt := range lexer.NonTerms {
		follow[nt] = make(StringSet, 0)
	}
	// в стартовый нетерминал помещается символ конца строки
	follow[lexer.Start.Symbol]["$"] = nil

	var changed = true
	for changed {
		changed = false
		// for (A→α∈P)
		for A, r := range lexer.NonTerms {
			for _, rule := range r.Rules {
				// for (B:α=βBγ)
				for i, s := range rule.Symbols {
					if s.Type == NonTerm {
						B := s.Value
						var beforeLen = len(follow[s.Value])

						// FOLLOW[B] ∪= FIRST(γ)
						for j := i; j < len(rule.Symbols); j++ {
							if rule.Symbols[j].Type == Term {
								γ := rule.Symbols[j].Value
								follow[B][γ] = nil
							}
						}

						//  FOLLOW[B] ∪= FOLLOW[A]
						for term := range follow[A] {
							follow[B][term] = nil
						}

						// changed = true if FOLLOW[B] изменился
						var afterLen = len(follow[s.Value])
						if beforeLen != afterLen {
							changed = true
						}
					}
				}
			}
		}
	}
	return follow
}
