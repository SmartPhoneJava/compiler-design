package internal

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

// Посчитать не епсилон символы
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
