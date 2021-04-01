package internal

import (
	"fmt"
	"log"
	"unicode"
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

func (a *CFR) UpdateN() {
	var (
		mapVisited = make(map[string]interface{})
		newN       = make([]string, 0)
	)
	for _, r := range a.P {
		mapVisited[r.From] = nil
		noneTerms := ToNoneTerminals(r.To)
		for _, nt := range noneTerms {
			mapVisited[nt] = nil
		}
	}

	for k := range mapVisited {
		newN = append(newN, k)
	}
	a.N = newN
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
func (cfr CFR) ElrWithE() CFR {
	var (
		newRules   = make(Rules, len(cfr.P))
		newSymbols = make([]string, len(cfr.N))
	)

	copy(newSymbols, cfr.N)
	copy(newRules, cfr.P)

	for i := 0; i < len(cfr.N); i++ {
		var Aᵢ = cfr.N[i]

		for j := 0; j < i; j++ {
			var (
				A𝚥    = cfr.N[j]
				fromA = cfr.P.FilterByTwo(Aᵢ, A𝚥, newSymbols)
				β     = cfr.P.Filter(A𝚥, NoSort)
				// fromA = newRules.FilterByTwo(Aᵢ, A𝚥, newSymbols)
				// β     = newR
			)

			// Починить багу
			log.Println("add", Aᵢ, A𝚥, len(fromA))
			for _, ruleA := range fromA {
				var (
					α  = ruleA.RemoveSymbol(A𝚥)
					αβ = β.Add(α).GetRPart()
				)
				log.Println("add1", Aᵢ, len(fromA))
				(&newRules).Append(Aᵢ, αβ...)

			}
			log.Println("Remove1", len(newRules))
			newRules = newRules.RemoveRulesFT(Aᵢ, A𝚥)
			log.Println("Remove2", len(newRules))
		}

		if newRules.HasDirectLeftRecursion(Aᵢ) {
			var (
				alpha, beta    = newRules.Filter(Aᵢ, LeftRecursion).AlphaBeta(Aᵢ)
				marked         = Aᵢ + "'"
				betasWithNewA  = beta.Add(marked).GetRPart()
				alphasWithNewA = alpha.Add(marked).GetRPart()
			)

			newRules = newRules.RemoveRules(Aᵢ)
			(&newRules).Append(Aᵢ, betasWithNewA...)
			(&newRules).Append(marked, append(alphasWithNewA, Epsilon)...)
			newSymbols = append(newSymbols, marked)
		}
	}

	return CFR{
		N: newSymbols,
		T: cfr.T,
		P: newRules.DeleteE(),
		S: cfr.S,
	}
}

/*
Алгоритм 4.10 из "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание", стр 278

Oсновная идея левой факторизации в том, что в том случае, когда неясно, какую из двух альтернатив
надо использовать для развертки нетерминала A, нужно изменить A - правила так, чтобы отложить решение до
тех пор, пока не будет достаточно информации для принятия правильного решения.

Если A ->  αB | αC - два A-правила и входная цепочка начинается с непустой строки, выводимой из α,
мы не знаем, разворачивать ли по первому правилу или по второму. Можно отложить решение, развернув A -> αA'.
Тогда после анализа того, что выводимо из α, можно развернуть по A'->B или по A'->C.

Левофакторизованные правила принимают вид

A -> αA'

A'-> B|C

https://intuit.ru/studies/courses/1157/173/lecture/4697?page=7
*/
func (cfr CFR) LeftFactorization() CFR {
	var (
		newRules   Rules
		newSymbols = make([]string, len(cfr.N))
	)
	copy(newSymbols, cfr.N)

	for _, n := range cfr.N {
		rulesGet, done := cfr.P.ToLFS(n)
		if done {
			newSymbols = append(newSymbols, n+"'")
		}

		newRules = append(newRules, rulesGet...)
	}

	// log.Println("newSymbols", newSymbols)
	return CFR{
		N: newSymbols,
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}
}

/*

Воспользоваться определением на стр. 175, алгоритмом 2.9. и алгоритмом 2.10. [1]. При
тестировании воспользоваться упражнением 2.4.13. [1].

стр. 171 алгоритм 2.9

Для того чтобы преобразовать произвольную КС-грамматику к приведенному виду, необходимо выполнить следующие действия:

удалить все бесплодные символы;
удалить все недостижимые символы;
удалить ^.-правила;
удалить цепные правила.
https://uz.denemetr.com/docs/294/index-20812-1.html?page=7

*/

// Удалить недостижимые нетерминалы
/*
 Алгоритм 2.7. , стр. 169
 https://neerc.ifmo.ru/wiki/index.php?title=Удаление_бесполезных_символов_из_грамматики

Шаг 0. Множество достижимых нетерминалов состоит из единственного элемента: {S}.
Шаг 1. Если найдено правило, в левой части которого стоит нетерминал,
	содержащийся в множестве, добавим в множество все нетерминалы из правой части.
Шаг 2. Повторим предыдущий шаг, если множество порождающих нетерминалов изменилось.
	Получаем множество всех достижимых нетерминалов, а нетерминалы, не попавшие в него,
	являются недостижимыми.

	// Сложность O(P)
*/
func (cfr CFR) RemoveUnreachableNonterminal() CFR {
	if len(cfr.N) == 0 {
		return cfr
	}
	var (
		mapVisited = make(map[string]interface{})
		fromTo     = make(map[string][]string)
		queue      = make([]string, len(cfr.S))
	)
	copy(queue, cfr.S)

	// O(P)
	for _, q := range cfr.P {
		goTo, ok := fromTo[q.From]
		if !ok {
			goTo = ToNoneTerminals(q.To)
		} else {
			goTo = append(goTo, ToNoneTerminals(q.To)...)
		}
		fromTo[q.From] = goTo
	}

	// O(N)
	for len(queue) > 0 {
		head := queue[0]
		queue = queue[1:]
		// O(1)
		for _, to := range fromTo[head] {
			_, ok := mapVisited[to]
			if ok {
				continue
			}
			mapVisited[to] = nil
			queue = append(queue, to)
		}
	}

	for _, s := range cfr.S {
		mapVisited[s] = nil
	}

	var newRules Rules

	// Добавляем правила, где левая часть содержит достижимый нетерминал
	for _, v := range cfr.P {
		_, ok := mapVisited[v.From]
		if ok {
			newRules = append(newRules, v)
		}
	}

	newCfr := &CFR{
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}

	newCfr.UpdateN()
	return *newCfr
}

// RemoveNongeneratingNonterminal - Удалить непорождающие нетерминалы
/*
 Алгоритм 2.8. , стр. 169
 https://neerc.ifmo.ru/wiki/index.php?title=Удаление_бесполезных_символов_из_грамматики#.D0.94.D0.BE.D1.81.D1.82.D0.B8.D0.B6.D0.B8.D0.BC.D1.8B.D0.B5_.D0.B8_.D0.BD.D0.B5.D0.B4.D0.BE.D1.81.D1.82.D0.B8.D0.B6.D0.B8.D0.BC.D1.8B.D0.B5_.D0.BD.D0.B5.D1.82.D0.B5.D1.80.D0.BC.D0.B8.D0.BD.D0.B0.D0.BB.D1.8B

Шаг 0. Множество порождающих нетерминалов пустое.
Шаг 1. Находим правила, не содержащие нетерминалов в правых частях
	и добавляем нетерминалы, встречающихся в левых частях таких правил,
	в множество.
Шаг 2. Если найдено такое правило, что все нетерминалы, стоящие в его
	правой части, уже входят в множество, то добавим в множество нетерминалы,
	стоящие в его левой части.
Шаг 3. Повторим предыдущий шаг, если множество порождающих нетерминалов
	изменилось.
	В результате получаем множество всех порождающих нетерминалов грамматики,
	а все нетерминалы, не попавшие в него, являются непорождающими.

	// Сложность O(P)
*/
func (cfr CFR) RemoveNongeneratingNonterminal() CFR {
	if len(cfr.N) == 0 {
		return cfr
	}

	var (
		// Число непорождающих нетермов в правилах
		ruleCounter = make(map[*RuleWithTerms]interface{})
		// Посещенные вершины
		mapVisited = make(map[string]interface{})
		// Очередь нетермов, которые могут привести в терминальное состояние
		queue = []string{}
	)

	for i, q := range cfr.P {
		var (
			noneTerms = ToNoneTerminalsMap(q.To)
			rterms    = RuleWithTerms{
				r:         &cfr.P[i],
				noneTerms: noneTerms,
			}
		)

		ruleCounter[&rterms] = len(noneTerms)
		if len(noneTerms) == 0 {
			_, ok := mapVisited[q.From]
			if ok {
				continue
			}
			queue = append(queue, q.From)
			mapVisited[q.From] = nil
		}
	}
	for _, s := range cfr.S {
		mapVisited[s] = nil
	}

	for len(queue) > 0 {
		var localQueue = make([]string, len(queue))
		copy(localQueue, queue)
		queue = []string{}

		for k := range ruleCounter {
			for _, lq := range localQueue {
				// Если у правила не осталось неподходящих нетермов
				if len(k.noneTerms) == 0 {
					break
				}
				// Удаляем текущий нетерм
				delete(k.noneTerms, lq)

				// Теперь если справа не осталось неподходящих нетермов
				// то нетерм слева стал подходящим, пометим это
				if len(k.noneTerms) == 0 {
					_, ok := mapVisited[k.r.From]
					if !ok {
						break
					}
					mapVisited[k.r.From] = nil
					queue = append(queue, k.r.From)
				}
			}
		}
	}

	var newRules Rules

	// Добавляем правила, где не осталось недопустимых нетермов
	for rc := range ruleCounter {
		if len(rc.noneTerms) == 0 {
			newRules = append(newRules, *rc.r)
		}
	}

	newCfr := &CFR{
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}

	newCfr.UpdateN()
	return *newCfr
}

// RuleWithTerms - правило с неподходящими вершинами -
// кандидами в непорождающие нетерминалы
type RuleWithTerms struct {
	r         *Rule
	noneTerms map[string]interface{}
}

// ToNoneTerminals Получить массив нетерминалов из строки
func ToNoneTerminals(str string) []string {
	var noneTerminals = []string{}
	for _, r := range str {
		if unicode.IsUpper(r) {
			noneTerminals = append(noneTerminals, string(r))
		}
	}
	return noneTerminals
}

func ToNoneTerminalsMap(str string) map[string]interface{} {
	var noneTerminals = make(map[string]interface{})
	for _, r := range str {
		if unicode.IsUpper(r) {
			noneTerminals[string(r)] = nil
		}
	}
	return noneTerminals
}

// http://mathhelpplanet.com/static.php?p=privedennaya-forma-ks-grammatiki
// file:///home/artyom/Загрузки/formal.languages.theory.3.pdf
func (cfr CFR) RemoveLambda() CFR {
	if len(cfr.N) == 0 {
		return cfr
	}

	var (
		// Обновленные правила
		mapNewRules = make(map[string]*Rule)
		// Посещенные вершины
		mapVisited = make(map[string]interface{})
		// Очередь нетермов, которые имеют пустой переход
		queue = []string{}
	)
	// Определяем нетермы с пустыми переходами
	for _, q := range cfr.P {
		if q.To == "e" {
			_, ok := mapVisited[q.From]
			if ok {
				continue
			}
			queue = append(queue, q.From)
			mapVisited[q.From] = nil
		}
	}

	// Помещаем все правила в mapNewRules
	for i, q := range cfr.P {
		if q.To == "e" {
			continue
		}
		_, ok := mapNewRules[q.To+q.From]
		if ok {
			continue
		}

		mapNewRules[q.From+q.To] = &cfr.P[i]
	}

	for len(queue) > 0 {
		var localQueue = make([]string, len(queue))
		copy(localQueue, queue)
		queue = []string{}

		for _, lq := range localQueue {
			for _, r := range mapNewRules {
				strs := r.ApplyEpsilon(lq)
				for _, str := range strs {
					if str == r.From || str == "" {
						if str == "" {
							_, ok := mapVisited[r.From]
							if !ok {
								queue = append(queue, r.From)
								mapVisited[r.From] = nil
							}
						}
						continue
					}
					_, ok := mapNewRules[str+r.From]
					if ok {
						continue
					}

					mapNewRules[str+r.From] = &Rule{
						From: r.From,
						To:   str,
					}
				}
			}
		}
	}

	var newRules Rules

	// Добавляем обновленные правила
	for _, rc := range mapNewRules {
		newRules.Append(rc.From, rc.To)
	}

	newCfr := &CFR{
		T: cfr.T,
		P: newRules,
		S: cfr.S,
	}

	newCfr.UpdateN()
	return *newCfr
}

// 551
