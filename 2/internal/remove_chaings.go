package internal

// GroupByChains - Сгруппировать нетермы по наличию цепных правил
/*
первый возвращаемый параметр - мапа, где ключ - нетерм
значения - цепные правила(правила вида A->B)

второй возвращаемый параметр - мапа нетермов без цепн. правил

O(P)
*/
func (cfr CFR) groupByChains() (map[string][]string, map[string][]string) {
	var withChains = make(map[string][]string, 0)
	var noChains = make(map[string][]string, 0)
	for _, rule := range cfr.P {
		if cfr.isChainRule(rule) {
			withChains[rule.From] = append(withChains[rule.From], rule.To)
		} else {
			noChains[rule.From] = append(noChains[rule.From], rule.To)
		}
	}

	return withChains, noChains
}

// IsChainRule является ли правило цепным
func (cfr CFR) isChainRule(r Rule) bool {
	nt := cfr.toNoneTerminals(r.To)
	return len(nt) == 1 && len(r.To) == len(nt[0])
}
