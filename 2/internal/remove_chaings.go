package internal

// GroupByChains - Сгруппировать нетермы по наличию цепных правил
/*
первый возвращаемый параметр - мапа, где ключ - нетерм
значения - цепные правила(правила вида A->B)

второй возвращаемый параметр - мапа нетермов без цепн. правил

O(P)
*/
func (cfr CFR) groupByChains() (
	map[string]map[string]interface{},
	map[string]map[string]interface{},
) {
	var withChains = make(map[string]map[string]interface{}, 0)
	var noChains = make(map[string]map[string]interface{}, 0)
	for _, rule := range cfr.P {
		if cfr.isChainRule(rule) {
			_, ok := withChains[rule.From]
			if !ok {
				withChains[rule.From] = make(map[string]interface{}, 0)
			}
			withChains[rule.From][rule.To] = nil
		} else {
			_, ok := noChains[rule.From]
			if !ok {
				noChains[rule.From] = make(map[string]interface{}, 0)
			}
			noChains[rule.From][rule.To] = nil
		}
	}

	return withChains, noChains
}

// IsChainRule является ли правило цепным
func (cfr CFR) isChainRule(r Rule) bool {
	nt := cfr.toNoneTerminals(r.To)
	return len(nt) == 1 && len(r.To) == len(nt[0])
}
