package internal

import (
	"fmt"
	"strings"
)

const Epsilon = "e"

// Rule - правило перехода
type Rule struct {
	// From - всегда одиночный нетерминал в случае КС грамматики
	From, To string
}

// Правила
type Rules []Rule

const (
	// LeftRecursion - сортировка вида
	// A -> Aα1 | Aα2 | ...| Aαn | β1 | β2 | ... | βn
	LeftRecursion = iota
	NoSort
)

func (rule Rule) ID() string {
	return rule.To + rule.From
}

// rules - получить правила, связанные с нетерминалом
// 	в порядке, заданным кодом sortType
func (rules Rules) Filter(
	noneTerminal string,
	sortType uint,
) Rules {
	var alpha, beta Rules

	for _, rule := range rules {
		if rule.From == noneTerminal {
			if sortType == LeftRecursion {
				if rule.IsLeftRecursive() {
					alpha = append(alpha, rule)
				} else {
					beta = append(beta, rule)
				}
			} else if sortType == NoSort {
				alpha = append(alpha, rule)
			}
		}
	}
	return append(alpha, beta...)
}

// FilterByTwo - отфильтровать правила так, чтобы в левой части был
// left, а правая начиналась с right
func (rules Rules) FilterByTwo(
	left, right string,
	noneTerms []string,
) Rules {
	var filtered Rules
	for _, rule := range rules {
		if rule.From == left && rule.IsItFirstNoneTerminal(right, noneTerms) {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

// HasDirectLeftRecursion - проверить, есть ли
//  прямая левая рекурсия в наборе правил
func (rules Rules) HasDirectLeftRecursion(
	noneTerminal string,
) bool {
	for _, r := range rules {
		if r.From == noneTerminal && r.IsLeftRecursive() {
			return true
		}
	}
	return false
}

// ConnectedPair - вернуть правила, у которых левая часть
//  равна `a`, а правая начинается с `b`
func (rules Rules) ConnectedPair(
	a, b string,
) Rules {
	var ab Rules
	for _, rule := range rules {
		if rule.From == a && rule.RightBeginFrom(b) {
			ab = append(ab, rule)
		}
	}
	return ab
}

// AlphaBeta разбить набор правил на следующие:
// Aα1 | Aα2 | ...| Aαn
// β1 | β2 | ... | βn, где β не начинается на A
// в роли A выступает нетерминал symbol
func (rules Rules) AlphaBeta(symbol string) (Rules, Rules) {
	var (
		alpha, beta Rules
	)
	for _, rule := range rules {
		if rule.RightBeginFrom(symbol) {

			alpha = append(alpha, rule)
		} else {
			beta = append(beta, rule)
		}
	}
	return alpha.RemoveFirst(len(symbol)), beta
}

// MarkLeftRecursives - получить нетерминалы с левой
//  рекурсией
func (rules Rules) MarkLeftRecursives() map[string]bool {
	var (
		noneTerminalsMap = make(map[string]bool, 0)
	)
	for _, rule := range rules {
		if rule.IsLeftRecursive() {
			_, ok := noneTerminalsMap[rule.From]
			if !ok {
				noneTerminalsMap[rule.From] = true
			}
		}
	}
	return noneTerminalsMap
}

// Пример 1: A' -> AV
// первое условие позволяет убедиться, что A' != A, недостаточно проверять
// первый символ обоих частей, надо смотреть то же число символов
// Пример 2: A -> A'V
// Недостаточно смотреть только на 1-ый символ, надо убедиться
// что нет символа ', поэтому проверяем что помеченный A не будет равен
// обнаруженному A'
func (r Rule) IsLeftRecursive() bool {
	return r.RightBeginFrom(r.From) && !r.RightBeginFrom(r.NewMarked())
}

// NewMarked вернуть помеченный нетерминал
func (r Rule) NewMarked() string {
	return r.From + "'"
}

// Append добавить правила переходов из from в каждый из to
func (r *Rules) Append(from string, to ...string) {
	var unique = make(map[string]bool)
	for _, r := range *r {
		unique[r.From+r.To] = true
	}

	for _, to := range to {
		// не добавляем пустые
		if len(strings.TrimSpace(to)) == 0 {
			continue
		}
		_, ok := unique[from+to]
		if ok { // боремся с дублями
			continue
		}
		unique[from+to] = true
		*r = append(*r, Rule{From: from, To: to})
	}
}

// DeleteE удалить пустые порождения
func (r Rules) DeleteE() Rules {
	for i, rule := range r {
		if r[i].To == "e" {
			continue
		}
		r[i].To = strings.ReplaceAll(rule.To, "e", "")
	}
	return r
}

// Add добавить цепочку к правой части всех правил
func (r Rules) Add(addMe string) Rules {
	// var newR = make(Rules, len(r))
	// copy(newR, r)
	// for i := range r {
	// 	newR[i].To = newR[i].To + addMe
	// }
	// return newR

	for i := range r {
		r[i].To = r[i].To + addMe
	}
	return r
}

// RemoveFirst удалить c левой части n символов
func (r Rules) RemoveFirst(n int) Rules {
	for i := range r {
		if len(r[i].To) >= n {
			r[i].To = r[i].To[n:]
		}
	}
	return r
}

// RemoveRules удалить правила где в левой части стоит A
// Например A->a, A->Be, A->e и т.д.
func (r Rules) RemoveRules(A string) Rules {
	var newRules Rules
	for i := range r {
		if r[i].From == A {
			continue
		}
		newRules = append(newRules, r[i])
	}
	return newRules
}

// RemoveRules удалить правила вида A -> Ba
func (r Rules) RemoveRulesFT(A, B string) Rules {
	var newRules Rules
	for i := range r {
		if r[i].From == A && r[i].RightBeginFrom(B) {
			continue
		}
		newRules = append(newRules, r[i])
	}
	return newRules
}

// RemoveSymbol - удалить из правой части правила
//  столько символов, сколько занимает строка s
func (r Rule) RemoveSymbol(s string) string {
	return r.RemoveFirstN(len(s))
}

// RemoveFirst - удалить из правой части правила
//  первые n символов
func (r Rule) RemoveFirstN(n int) string {
	if len(r.To) >= n {
		r.To = r.To[n:]
	}
	return r.To
}

// GetRPart получить правую часть правил
func (rules Rules) GetRPart() []string {
	var rights = make([]string, len(rules))
	for i, rule := range rules {
		rights[i] = rule.To
	}
	return rights
}

// BeginFrom - проверить, что правая часть начинается с symbol
func (rule Rule) RightBeginFrom(symbol string) bool {
	if len(rule.To) > len(symbol) {
		if rule.To[len(symbol)] == '\'' {
			return false
		}
	}
	return len(rule.To) >= len(symbol) &&
		rule.To[:len(symbol)] == symbol
}

// BeginFrom - проверить, что правая часть начинается с symbol
func (rule Rule) IsItFirstNoneTerminal(
	symbol string,
	noneTerms []string,
) bool {
	return symbol == rule.FirstNoneTerminal(noneTerms)
}

// FirstNoneTerminal получить первый нетерминал правой части
func (rule Rule) FirstNoneTerminal(
	noneTerms []string,
) string {
	var (
		found     string
		searchStr string
	)
	for _, r := range rule.To {
		searchStr += string(r)
		if len(found) != 0 {
			if r == '\'' {
				found += "'"
				continue
			}
			return found
		} else {
			for _, v := range noneTerms {
				if searchStr == v {
					found = v
					break
				}
			}
		}
	}
	return found
}

// сравнить два набора правил
func (a Rules) IsSame(b Rules) error {
	if len(a) != len(b) {
		return fmt.Errorf("Размерности правил не сходятся. Ожидалось: %d, получено %d", len(a), len(b))
	}
	var (
		fromA = make([]string, 0)
		toA   = make([]string, 0)
		fromB = make([]string, 0)
		toB   = make([]string, 0)
	)
	for i := range a {
		fromA = append(fromA, a[i].From)
		toA = append(toA, a[i].To)
		fromB = append(fromB, b[i].From)
		toB = append(toB, b[i].To)
	}
	err1 := areStringSlicesSame(fromA, fromB)
	err2 := areStringSlicesSame(toA, toB)
	var err error
	if err1 != nil {
		err = fmt.Errorf("Вершины из: %s", err1)
		if err2 != nil {
			err = fmt.Errorf("%s Вершины из: %s", err1, err2)
		}
	} else {
		if err2 != nil {
			err = fmt.Errorf("Вершины в: %s", err2)
		}
	}
	return err
}

// ApplyEpsilon Получить все правила, которые могут получиться
//  из данного если вместо каждого нетерма nt подставить епсилон
//  Возвращает всевозможные правые части и флаг было ли совершено
//  преобразование
func (a Rule) ApplyEpsilon(cfr CFR, nt string) []string {
	m := cfr.ToNoneTerminalsMap(a.To)
	_, ok := m[nt]
	if !ok {
		return nil
	}
	return applyEpsilon(0, "", a.To, nt, cfr.CountOfNoneE(nt) == 0)
}

func applyEpsilon(
	index int,
	start, word, nt string,
	noOriginNT bool) []string {
	if len(word) == index {
		return []string{start}
	}
	var strs []string
	s := word[index]
	if string(s) == nt {
		strs = append(strs, applyEpsilon(
			index+1, start, word,
			nt, noOriginNT,
		)...)
		if !noOriginNT {
			strs = append(strs, applyEpsilon(
				index+1, start+string(word[index]),
				word, nt, noOriginNT,
			)...)

		}
	} else {
		strs = append(strs, applyEpsilon(
			index+1, start+string(word[index]),
			word, nt, noOriginNT,
		)...)
	}

	return strs
}

func (cfr *CFR) CountOfNoneE(nt string) int {
	var count int
	for _, r := range cfr.P {
		if r.From == nt && r.To != "e" {
			count++
		}
	}
	return count
}

// Убрать повторяющиеся e, примеры
// AAeeAeA -> AAeAeA
// eeee -> e
// eeeeAeee -> eAe
func deleteEduplicates(str string) string {
	var newStr []rune
	var foundE bool
	for _, r := range str {
		if r == 'e' {
			if foundE {
				continue
			}
			foundE = true
		} else {
			foundE = false
		}
		newStr = append(newStr, r)
	}
	return string(newStr)
}
