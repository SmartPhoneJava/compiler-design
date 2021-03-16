package main

import (
	"errors"
	"fmt"
	"log"
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
	NoSort
)

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

func (rules Rules) HasLeftRecursion(
	noneTerminal string,
) bool {
	for _, r := range rules {
		if r.From == noneTerminal && r.IsLeftRecursive() {
			log.Println("recurse", r.From, r.To)
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
	for _, to := range to {
		log.Printf("\n Want add %s -> %s", from, to)
	}
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

func (r *Rules) AppendRules(b Rules) {
	var unique = make(map[string]bool)
	for _, r := range *r {
		unique[r.From+r.To] = true
	}

	for _, b := range b {
		_, ok := unique[b.From+b.To]
		if ok { // боремся с дублями
			continue
		}
		unique[b.From+b.To] = true
		*r = append(*r, Rule{From: b.From, To: b.To})
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

// RemoveRules удалить правила содержащие from
func (r Rules) RemoveRules(from string) Rules {
	var newRules Rules
	for i := range r {
		if r[i].From == from {
			continue
		}
		newRules = append(newRules, r[i])
	}
	return newRules
}

func (r Rule) RemoveSymbol(s string) string {
	return r.RemoveFirst(len(s))
}

func (r Rule) RemoveFirst(n int) string {
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
	return len(rule.To) >= len(symbol) && rule.To[:len(symbol)] == symbol
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

// Устранить левую рекурсию
// Алгоритм 2.13
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
				alpha, beta = cfr.P.Filter(r.From, LeftRecursion).AlphaBeta(r.From)
				marked      = r.NewMarked()
			)
			alpha = alpha.RemoveFirst(len(r.From))

			var (
				betas          = beta.GetRPart()
				betasWithNewA  = beta.Add(marked).GetRPart()
				alphas         = alpha.GetRPart()
				alphasWithNewA = alpha.Add(marked).GetRPart()
			)

			(&newRules).Append(r.From, append(betas, betasWithNewA...)...)
			(&newRules).Append(marked, append(alphas, alphasWithNewA...)...)
			markedPoints[r.From] = false
			newSymbols = append(newSymbols, marked)
		} else {
			(&newRules).Append(r.From, r.To)
		}
	}

	return CFR{
		N: newSymbols,
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}
}

// ElrWithE - Устранить левую рекурсию, оставив e-продукцию
/*
Алгоритм 4.8 из "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание", стр 277
Гарантированно работает с грамматиками, не имеющими:
- циклов(порождений A -> A)
- e-продукций(продукций вида A -> e)
*/
//  4.8 и 4.10.
func (cfr CFR) ElrWithE() CFR {
	var (
		newRules   Rules
		newSymbols = make([]string, len(cfr.N))
	)

	copy(newSymbols, cfr.N)

	for i := 0; i < len(cfr.N); i++ {
		var (
			Aᵢ    = cfr.N[i]
			fromA = cfr.P.Filter(Aᵢ, NoSort)
		)
		for j := 0; j < i; j++ {
			var (
				A𝚥 = cfr.N[j]
				β  = cfr.P.Filter(A𝚥, NoSort)
			)
			(&newRules).Append(A𝚥, β.GetRPart()...)
			for _, ruleA := range fromA {
				if ruleA.RightBeginFrom(A𝚥) {
					var (
						α  = ruleA.RemoveSymbol(A𝚥)
						αβ = β.Add(α).GetRPart()
					)
					(&newRules).Append(Aᵢ, αβ...)
				} else {
					(&newRules).Append(Aᵢ, ruleA.To)
				}

			}
		}

		if newRules.HasLeftRecursion(Aᵢ) {
			var (
				alpha, beta = newRules.Filter(Aᵢ, LeftRecursion).AlphaBeta(Aᵢ)
				marked      = Aᵢ + "'"
			)
			if len(alpha)+len(beta) > 0 {
				alpha = alpha.RemoveFirst(len(Aᵢ))

				var (
					betasWithNewA  = beta.Add(marked).GetRPart()
					alphasWithNewA = alpha.Add(marked).GetRPart()
				)

				newRules = newRules.RemoveRules(Aᵢ)
				(&newRules).Append(Aᵢ, betasWithNewA...)
				(&newRules).Append(marked, append(alphasWithNewA, Epsilon)...)
				newSymbols = append(newSymbols, marked)
			}
		}

	}
	return CFR{
		N: newSymbols,
		T: cfr.T,
		P: newRules.DeleteE(),
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

// 464
