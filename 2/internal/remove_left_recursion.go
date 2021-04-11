package internal

import (
	"strings"
	"unicode"
)

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

// replaceRules - заменить правило на множество
func (cfr CFR) replaceRule(
	rules, updated Rules,
	Ai string,
) Rules {
	if len(rules) == 0 {
		return rules
	}
	var symbolsNum = make(map[string]int)
	for _, n := range cfr.N {
		_, ok := symbolsNum[n]
		if ok {
			continue
		}
		symbolsNum[n] = len(symbolsNum)
	}

	var (
		newRules    = make(Rules, len(rules))
		changed     = true
		returnRules Rules
		groups      = updated.Group()
	)

	copy(newRules, rules)
	for changed {
		changed = false
		var newRulesAgain Rules
		for _, a := range newRules {
			arr := cfr.toNoneTerminals(a.To)
			// Если правило теперь ведет только в терминальное состояние
			if len(arr) == 0 {
				returnRules.Append(a.From, a.To)
				continue
			}
			var isLower bool
			for _, r := range a.To {
				isLower = unicode.IsLower(r)
				break
			}
			// Нет левой рекурсии
			if isLower {
				returnRules.Append(a.From, a.To)
				continue
			}
			Aj := arr[0]
			// Если нетерм дальше или равен по нумерации
			if symbolsNum[Ai] <= symbolsNum[Aj] {
				returnRules.Append(a.From, a.To)
				continue
			}
			changed = true
			a.To = a.RemoveSymbol(Aj)
			fromAj := groups[Aj]

			rpart := fromAj.Add(a.To).GetRPart()

			// Подставляем новые правила
			newRulesAgain.Append(a.From, rpart...)
		}
		if len(newRulesAgain) == 0 {
			break
		}
		newRules = make(Rules, len(newRulesAgain))
		copy(newRules, newRulesAgain)
	}
	return returnRules
}
