package opa

// operator_precedence_analyzer.go
import (
	"fmt"
	"lab2/internal/g5"

	"github.com/fatih/color"
)

const (
	NO byte = iota
	EQ
	MORE
	LESS
	DONE
)

const (
	EqText   = "="
	MoreText = "▶"
	LessText = "◀"
)

const StartEnd = "⏊"
const AnalyserNonTerm = "E"

// Matrix of operator precedence relations
type OperatorsMatrix map[string]map[string]byte

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
	// var arr = make([]string, 0, len(matrix))
	// var elemIndex = map[string]int{}
	var arr = []string{";", "if", "then", "else", "a", "=", "or", "xor", "and", "(", ")", "⏊"}
	var elemIndex = map[string]int{
		";":    0,
		"if":   1,
		"then": 2,
		"else": 3,
		"a":    4,
		"=":    5,
		"or":   6,
		"xor":  7,
		"and":  8,
		"(":    9,
		")":    10,
		"⏊":    11,
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
				text = EqText
			case MORE:
				text = MoreText
			case LESS:
				text = LessText
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

/*

1. Составляем новый список правил, заменив нетермы на E
2. У нас 4 структуры:
	-	строка разбора
	- 	предыдущий символ,
	-	лента(остаток кода юзера),
	-	номера примененных правил
3. Пока код юзера не опустел, смотрим его первым символ
4. находим в матрице элемент [прошлый_символ_ленты][текущий_символ_ленты]
4* прошлый_символ_ленты в первой итерации это символ старта $
5. Если <, =:
	Добавляем текущий элемент ленты в строку разбора
	Меняем прошлый элемент ленты на текущий
	Двигаем ленту на 1 вправо
   Если >:
	Ищем правило, где строка разбора находится справа и
	заменяем эту часть в строке разбора на правую часть правила
	В массив номеров добавляем номер правила, которое смогло перевести
	Если ни одного правило не отработало возвращаем ошибку

*/

// Operator Precedence Analyzer
type Analyzer struct {
	Rules  g5.Rules
	Matrix OperatorsMatrix
	terms  map[string]bool
}

// Установим новые правила для анализатора
func NewAnalyzer(lexer g5.Lexer) *Analyzer {
	var (
		analyser = &Analyzer{}
		ruleMap  = make(map[string]*g5.Rule)
	)
	for _, res := range lexer.NonTerms {
		for _, rule := range res.Rules {
			var (
				newRuleString string
				newRule       = &g5.Rule{}
			)
			for _, symbol := range rule.Symbols {
				if symbol.Type == g5.NonTerm {
					newRuleString += " <" + AnalyserNonTerm + "> "
					newRule.Symbols = append(newRule.Symbols, g5.Symbol{
						Value: AnalyserNonTerm,
						Type:  symbol.Type,
					})
				} else {
					newRuleString += " " + symbol.Value + " "
					newRule.Symbols = append(newRule.Symbols, g5.Symbol{
						Value: symbol.Value,
						Type:  symbol.Type,
					})
				}
			}
			if newRuleString == " <"+AnalyserNonTerm+"> " {
				continue
			}
			ruleMap[newRuleString] = newRule
		}
	}

	analyser.Rules = make(g5.Rules, 0, len(ruleMap))
	for _, rule := range ruleMap {
		analyser.Rules = append(analyser.Rules, *rule)
	}

	analyser.terms = lexer.Terms
	analyser.Matrix = MakeMatrixV2(lexer)
	analyser.Matrix[StartEnd][StartEnd] = DONE // помечаем, что при
	// совпадении спец. символа успешно зканчиваем
	return analyser
}

func (analyser Analyzer) PrintRules() {
	var resolver = &g5.Resolver{
		Rules:  analyser.Rules,
		Symbol: AnalyserNonTerm,
	}
	var printingLexer = &g5.Lexer{
		NonTerms: map[string]*g5.Resolver{
			AnalyserNonTerm: resolver,
		},
		Terms: analyser.terms,
		Start: resolver,
	}

	resolver.Lexer = printingLexer

	printingLexer.Print("\n\nПравила внутри анализатора")
}

func (analyser Analyzer) findRule(row *[]string) (g5.Rule, error) {
	for rowC := 1; rowC < len(*row); rowC++ {
		for _, rule := range analyser.Rules {
			if len(rule.Symbols) != rowC {
				continue
			}
			var matched = true

			for i := range rule.Symbols {
				if rule.Symbols[len(rule.Symbols)-1-i].Value != (*row)[len(*row)-1-i] {
					matched = false
					break
				}
			}

			if matched {
				*row = (*row)[:len(*row)-rowC+1]
				(*row)[len(*row)-1] = AnalyserNonTerm
				return rule, nil
			}
		}
	}
	return g5.Rule{}, fmt.Errorf("Правила не найдено для %s", row)
}

// Exec - попробовать считать введенные символы
/*
	Возвращает:
		- слайс символов грамматиков для составления АСТ
		- слайс правил, которые необходимо применить для получения
			введенной строки
		- ошибка, если введены некорретные входные данные
*/
func (analyser Analyzer) Exec(input []string) ([]string, g5.Rules, error) {
	// Добавляем символ конца ввода
	input = append(input, StartEnd)
	var (
		outputSymbols []string
		outputRules   g5.Rules
		// В стек сразу помещаем символа начала ввода
		stack = []string{StartEnd}
	)
	for len(stack) > 0 {
		var (
			currStack      = getFromStack(stack)
			currInput      = input[0]
			matrixOperator = analyser.Matrix[currStack][currInput]
		)

		switch {
		case matrixOperator == EQ || matrixOperator == LESS:
			stack = append(stack, currInput)
			input = input[1:]
		case matrixOperator == MORE:
			foundRule, err := analyser.findRule(&stack)
			if err != nil {
				return nil, nil, wrapError(err)
			}
			outputRules = append(outputRules, foundRule)
			outputSymbols = append(outputSymbols, currStack)
		case matrixOperator == DONE:
			return outputSymbols, outputRules, nil
		case matrixOperator == NO:
			var err error
			if currStack == StartEnd {
				err = fmt.Errorf("код не может начинаться с `%s`", currInput)
			} else {
				err = fmt.Errorf("ключевое слово `%s` не может находиться слева от `%s`", currStack, currInput)
			}
			return nil, nil, wrapError(err)
		default:
			var err = fmt.Errorf("обнаружен неопознанный символ в матрице '%d'", matrixOperator)
			return nil, nil, wrapError(err)
		}
	}
	return outputSymbols, outputRules, nil
}

func (analyser Analyzer) PrintlnExecResult(
	text, input string,
	symbols []string,
	rules g5.Rules,
) {
	fmt.Println("\n" + text)
	color.Cyan("Прочитанные символы: \n")
	var right string
	for _, s := range symbols {
		right += color.GreenString(s) + " "
	}
	fmt.Printf("%s\n", right)

	fmt.Printf("%s %s\n", color.CyanString("Правила, применённые к "), color.GreenString(input))

	for _, rule := range rules {
		var right string
		for _, s := range rule.Symbols {
			g5.ColorSymbol(s, &right)
		}
		fmt.Printf("%s → %s\n", color.RedString(AnalyserNonTerm), right)
	}
}

func wrapError(err error) error {
	return fmt.Errorf("Введенный код содержит ошибку: %s", err)
}

func getFromStack(stack []string) string {
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i] == AnalyserNonTerm {
			continue
		}
		return stack[i]
	}
	return ""
}
