package internal

// LeftFactorizationRules - компоненты новых правил одного нетерма
type LeftFactorizationRules struct {
	// α - самый длинный префикс
	Alpha string
	// β - Остатки правых частей, содержищих α
	Beta []string
	// γ - все правые части не содержащие α
	Gamma []string
}

type rulesPrefixGroup struct {
	// Правила с одним префиксом
	rules  Rules
	prefix string
}

type rulesPrefixGroups []rulesPrefixGroup

// rules - конвертировать правила к группам с правилами
func (rules Rules) ToGroup() rulesPrefixGroup {
	var group rulesPrefixGroup

	group.rules = make(Rules, len(rules))
	copy(group.rules, rules)

	var currGroups = group.Expand()
	for len(currGroups) > 0 {
		group = currGroups[0]
		var nextGroups = rulesPrefixGroups{}
		for _, g := range currGroups {
			newGroups := g.Expand()
			for _, ng := range newGroups {
				if len(ng.rules) > 1 {
					nextGroups = append(nextGroups, ng)
				}
			}
		}
		currGroups = nextGroups
	}

	return group
}

// Получить группы правил, где префиксы на 1 больше
func (group rulesPrefixGroup) Expand() rulesPrefixGroups {
	var m = make(map[string]*rulesPrefixGroup)
	for _, r := range group.rules {
		if len(r.To) == len(group.prefix) {
			continue
		}
		newPrefix := r.To[:len(group.prefix)+1]
		prefGroup, ok := m[newPrefix]
		if !ok {
			m[newPrefix] = &rulesPrefixGroup{
				prefix: newPrefix,
				rules:  Rules{r},
			}
		} else {
			prefGroup.rules.Append(r.From, r.To)
			m[newPrefix] = prefGroup
		}
	}
	var groups = make(rulesPrefixGroups, 0, len(m))
	for _, r := range m {
		groups = append(groups, *r)
	}
	return groups
}

// ToLFS - получить LeftFactorizationRules для нетерма noneTerminal
func (rules Rules) toLFS(group rulesPrefixGroup) LeftFactorizationRules {
	var (
		prefix = group.prefix
		lfr    = LeftFactorizationRules{Alpha: group.prefix}
	)

	for _, rule := range rules {
		to := rule.To
		if len(to) >= len(prefix) && to[:len(prefix)] == prefix {
			rest := to[len(prefix):]
			if len(rest) == 0 {
				rest = "e"
			}
			lfr.Beta = append(lfr.Beta, rest)
		} else {
			lfr.Gamma = append(lfr.Gamma, to)
		}

	}
	return lfr
}

func (rules Rules) ToLFS(
	noneTerminal string,
) (Rules, bool) {
	var (
		termRules = rules.Filter(noneTerminal, NoSort)
		group     = termRules.ToGroup()
	)
	if len(group.rules) < 2 {
		return termRules, false
	}
	lfs := termRules.toLFS(group)
	return *NewRulesfromLFS(noneTerminal, lfs), true
}

func NewRulesfromLFS(
	noneTerminal string,
	lfr LeftFactorizationRules,
) *Rules {
	var rules = &Rules{{
		From: noneTerminal,
		To:   lfr.Alpha + noneTerminal + "'",
	}}

	for _, gamma := range lfr.Gamma {
		rules.Append(noneTerminal, gamma)
	}
	for _, beta := range lfr.Beta {
		rules.Append(noneTerminal+"'", beta)
	}
	return rules
}
