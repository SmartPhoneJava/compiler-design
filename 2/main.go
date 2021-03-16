package main

import (
	"errors"
	"fmt"
	"strings"
)

// CFR(Context-free grammar) - Контекстно-свободная грамматика
// https://habr.com/ru/post/177701/
type CFR struct {
	// N — конечный алфавит нетерминальных символов
	N []string
	// T —  конечный алфавит терминальных символов (совпадает с алфавитом языка, задаваемого грамматикой)
	T []string
	// P — конечное множество правил порождения
	P Rules
	// S — начальный нетерминал грамматики G
	S []string
}

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
)

// rules - получить правила, связанные с нетерминалом
// 	в порядке, заданным кодом sortType
func (rules Rules) Filter(
	noneTerminal string,
	sortType uint,
) Rules {
	if sortType != LeftRecursion {
		return rules
	}
	var (
		alpha, beta Rules
	)
	for _, rule := range rules {
		if rule.From == noneTerminal {
			if rule.IsLeftRecursive() {
				alpha = append(alpha, rule)
			} else {
				beta = append(beta, rule)
			}
		}
	}
	return append(alpha, beta...)
}

// AlphaBeta разбить набор правил на следующие:
// Aα1 | Aα2 | ...| Aαn
// β1 | β2 | ... | βn, где β не начинается на A
func (rules Rules) AlphaBeta() (Rules, Rules) {
	var (
		alpha, beta Rules
	)
	for _, rule := range rules {
		if rule.IsLeftRecursive() {
			r := rule
			r.To = r.To[1:]
			alpha = append(alpha, r)
		} else {
			beta = append(beta, rule)
		}
	}
	return alpha, beta
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

func (r Rule) IsLeftRecursive() bool {
	return r.From == r.To[:1]
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

// Add добавить цепочку к правой части всех правил
func (r Rules) Add(addMe string) Rules {
	for i := range r {
		r[i].To = r[i].To + addMe
	}
	return r
}

// GetRPart получить правую часть правил
func (rules Rules) GetRPart() []string {
	var rights = make([]string, len(rules))
	for _, rule := range rules {
		rights = append(rights, rule.To)
	}
	return rights
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
	err := areStringSlicesSame(fromA, fromB)
	if err != nil {
		return err
	}

	err = areStringSlicesSame(toA, toB)
	if err != nil {
		return err
	}
	return nil
}

// Устранить левую рекурсию
// https://studizba.com/files/show/djvu/3050-1-tom-1.html стр. 180
// https://intuit.ru/studies/courses/1157/173/lecture/4697?page=7
func (cfr CFR) EliminateLeftRecursion() CFR {
	var (
		newRules     Rules
		newSymbols   = make([]string, len(cfr.N))
		markedPoints = cfr.P.MarkLeftRecursives()
	)

	/*
		отмечаем, что посетили нетерминал с левой рекурсией
		ставя false в markedPoints
	*/

	copy(newSymbols, cfr.N)
	for _, r := range cfr.P {
		if markedPoints[r.From] {
			var (
				alpha, beta = cfr.P.Filter(r.From, LeftRecursion).AlphaBeta()
				marked      = r.NewMarked()
			)

			var (
				betas          = beta.GetRPart()
				betasWithNewA  = beta.Add(marked).GetRPart()
				alphas         = alpha.GetRPart()
				alphasWithNewA = alpha.Add(marked).GetRPart()
			)

			newRules.Append(r.From, append(betas, betasWithNewA...)...)
			newRules.Append(marked, append(alphas, alphasWithNewA...)...)
			markedPoints[r.From] = false
			newSymbols = append(newSymbols, marked)
		} else {
			newRules.Append(r.From, r.To)
		}
	}

	return CFR{
		N: newSymbols,
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}
}

func (a CFR) IsSame(b CFR) error {
	if err := areStringSlicesSame(a.N, b.N); err != nil {
		return fmt.Errorf("Нетерминальные алфавиты не сходятся: %s", err)
	}

	if err := a.P.IsSame(b.P); err != nil {
		return fmt.Errorf("Правила не сходятся: %s", err)
	}

	if err := areStringSlicesSame(a.N, b.N); err != nil {
		return fmt.Errorf("Терминальные алфавиты не сходятся: %s", err)
	}

	if err := areStringSlicesSame(a.S, b.S); err != nil {
		return fmt.Errorf("Стартовые нетерминалы не сходятся: %s", err)
	}
	return nil
}

// сравнить два массива со строками
func areStringSlicesSame(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("Размерности не сходятся. Ожидалось: %d, получено %d", len(a), len(b))
	}
	type Pair struct{ E, R bool }
	var seen = make(map[string]Pair, 0)
	for i := range a {
		s := seen[a[i]]
		s.E = true
		seen[a[i]] = s

		s = seen[b[i]]
		s.R = true
		seen[b[i]] = s
	}
	var (
		errMsg string
	)
	for id, s := range seen {
		if s.E && s.R {
			continue
		}
		if !s.E {
			errMsg += fmt.Sprintf("\nНе хватает элемента: %s", id)
		} else {
			errMsg += fmt.Sprintf("\nЛишний элемент: %s", id)
		}
	}
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}

func main() {

}
