package opa

// operator_precedence_analyzer.go
import (
	"fmt"
	"lab2/internal/g5"
	"log"

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
func (analyser *Analyzer) Build(lexer g5.Lexer) {
	var (
		ruleMap = make(map[string]*g5.Rule)
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
			for i, symbol := range rule.Symbols {
				if symbol.Value != (*row)[len(*row)-1-i] {
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

func (analyser Analyzer) ToAst(rules g5.Rules) {
	var nodes = []*g5.Node{}
	savedArr := make([]*g5.Node, 0)
	var counter = 0
	for i := 0; i < len(rules); i++ {
		var r = rules[i]
		var ok bool
		var parent *g5.Node
		for _, el := range savedArr {
			parent = el
			break
		}
		if len(savedArr) > 0 {
			log.Println("mmmmmm", i)
		}
		for _, m := range savedArr {
			log.Println("m", m)
		}
		if !ok {
			parent = &g5.Node{
				ID:    fmt.Sprintf("%d.", counter),
				Value: AnalyserNonTerm,
				Type:  g5.NonTerm,
			}
			nodes = append(nodes, parent)
			counter++
			savedArr = append(savedArr, parent)
		}
		var right string
		for _, s := range r.Symbols {
			var node = &g5.Node{
				ID:          fmt.Sprintf("%d.", counter),
				Value:       s.Value,
				Parent:      parent,
				ParentValue: parent.Value,
				//Type:        rules[i].Rule.Symbols[j].Type,
			}
			log.Printf("\nconnect %s -> %s", s, parent.Value)
			counter++
			// if len(r.Symbols) > j {
			// 	if r.Symbols[j].Type == g5.NonTerm {
			// 		m[r.Symbols[j].Value] = node
			// 		node.Value = r.Symbols[j].Value
			// 	}

			// 	node.Type = r.Symbols[j].Type
			// 	if node.Type == g5.Reserved {
			// 		node.Type = g5.Term
			// 	}
			// }
			nodes = append(nodes, node)

		}
		for _, s := range r.Symbols {
			right += " " + s.Value
		}
		log.Printf("%s->%s", AnalyserNonTerm, right)
		if ok {
			savedArr = savedArr[1:]
		}
	}

	g5.MustVisualize(nodes, "assets", "hello2.dot")
}

func (analyser Analyzer) ToAstV2(rules g5.Rules) error {
	var nodes = []*g5.Node{}

	var counter = 2
	var root = &g5.Node{
		ID:    fmt.Sprintf("%d.", 1),
		Value: AnalyserNonTerm,
		Type:  g5.NonTerm,
	}
	freeNodes := []*g5.Node{root}
	nodes = append(nodes, root)

	for i := len(rules) - 1; i >= 0; i-- {
		var r = rules[i]
		var model, err = ToNumOperator(r)
		if err != nil {
			return err
		}

		node := freeNodes[0]
		newNodes := model.ToNodes(node, &counter)
		nodes = append(nodes, node)
		freeNodes = append(freeNodes[1:], newNodes...)
	}

	return g5.VisualizeFSM(nodes, "assets", "hello2.dot")
}

// ast переделка

func ToNumOperator(r g5.Rule) (NumOperator, error) {
	var (
		terms, nonTerms []string
	)
	for _, s := range r.Symbols {
		if s.Type == g5.Term {
			terms = append(terms, s.Value)
		} else {
			nonTerms = append(nonTerms, s.Value)
		}
	}
	switch {
	case len(terms) == 1 && len(nonTerms) == 2:
		return OneTwoOperatored{
			Main: terms[0],
		}, nil
	case len(terms) == 1 && len(nonTerms) == 0:
		return NoOperatored{
			Main: terms[0],
		}, nil
	default:
		return nil, fmt.Errorf("Нет модели для правила с %d термами и %d нетермами", len(terms), len(nonTerms))
	}
}

type NumOperator interface {
	ToNodes(node *g5.Node, counter *int) []*g5.Node
}

type OneTwoOperatored struct {
	Main string
}

func (two OneTwoOperatored) ToNodes(node *g5.Node, counter *int) []*g5.Node {
	node.Value = two.Main

	var leftNode = &g5.Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: two.Main,
		Type:        g5.Term,
	}
	*counter++
	var rightNode = &g5.Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: two.Main,
		Type:        g5.Term,
	}
	*counter++
	return []*g5.Node{leftNode, rightNode}
}

type NoOperatored struct {
	Main string
}

func (no NoOperatored) ToNodes(node *g5.Node, counter *int) []*g5.Node {
	node.Value = no.Main
	return nil
}
