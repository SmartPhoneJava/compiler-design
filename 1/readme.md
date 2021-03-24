## **Распознавание цепочек регулярного языка**

*Выполнено* ***Доктором А. А.*** _ИУ7-22м_

В рамках лабораторной работы №1 по курсу _Конструирование компиляторов_

## 1. Цель и задачи работы

**Цель работы**: приобретение практических навыков реализации важнейших элементов лексических анализаторов на примере распознавания цепочек регулярного языка.

**Задачи работы**:

1) Ознакомиться с основными понятиями и определениями, лежащими в основе построения лексических анализаторов.
2) Прояснить связь между регулярным множеством, регулярным выражением, праволинейным языком, конечно- автоматным языком и недетерминированным конечно-автоматным языком.
3) Разработать, тестировать и отладить программу распознавания цепочек регулярного или праволинейного языка в соответствии с предложенным вариантом грамматики.

## 2. Вариант задания

Напишите программу, которая в качестве входа принимает произвольное регулярное выражение, и выполняет

следующие преобразования:

1) По регулярному выражению строит НКА.
2) По НКА строит эквивалентный ему ДКА.
3) По ДКА строит эквивалентный ему КА, имеющий наименьшее возможное количество состояний. Указание. Воспользоваться алгоритмом, приведенным по адресу http://neerc.ifmo.ru/wiki/index.php?title=Алгоритм_Бржозовского
4) Моделирует минимальный КА для входной цепочки из терминалов исходной грамматики.

## 3. Реализация

### Структуры данных

1. [Граф](https://github.com/SmartPhoneJava/compiler-design/blob/56032407c7aa534d473c42d233e4548204a47100/1/internal/graph/graph.go#L29-L38) представлен наборами вершин и узлов, а так же множеством начальных и конечных узлов:

```go
type Graph struct {
	Vertexes map[string]*Vertex
	Edges    map[string]*Edge

	// Истоки и стоки
	First, Last []string
}
```

2. [Вершина](https://github.com/SmartPhoneJava/compiler-design/blob/56032407c7aa534d473c42d233e4548204a47100/1/internal/graph/graph.go#L9-L13) представляет собой связку из идентификатора и хэш-таблиц с входными и выходными дугами:

```go
type Vertex struct {
	ID  string
	Out map[string]*Edge
	In  map[string]*Edge
}
```

3. [Дуга](https://github.com/SmartPhoneJava/compiler-design/blob/56032407c7aa534d473c42d233e4548204a47100/1/internal/graph/graph.go#L16-L20) - переход из состояния `from` в `to` по символу `weight`:

```go
type Edge struct {
	From   string
	To     string
	Weight string
}
```

### Построение НКА по регулярному выражению

На вход поступает строка(`input`), на выходе конечный автомат, представленный графом. В начале граф содержит две вершины и соединяющую их дугу с весом `input`. Дальше происходит "разрастание" графа по `алгоритму Мак-Нотона-Ямады-Томпсона`[3]

Пока есть изменения (`changes > 0`) для каждой дуги графа, где длина строки веса превышает размер символа (`len(weight) > 1`) выполнить:

1. Поиск и удаление внешних скобок в строке
2. Поиск и замена замыкания. Если обнаружен символ `*`, убираем его и создаём дополнительные 2 узла и 4 дуги. 2 новых узла связываются дугой с весом, сооветствующей подстроке РВ без символа замыкания. В обратную сторону идёт дуга весом `e`. Такие же дуги связывают уже связанные вершины графа с новосозданными вершинами
3. Поиск и замена объединения. Если удалось найти `k`объединенй, где`k>1`, то удаляем текущую дугу между вершинами `qN`и`qM`и создаем`k`параллельных дуг от`qN`и`qM`до`qM`с весами - подстроками полученных из`weight`
4. Поиск и замена конкатенации. Если удалось найти `k`конкатенаций, где`k>1`, то удаляем текущую дугу между вершинами `qN`и`qM`и создаем последовательно`k`дуг с весами - подстроками полученных из`weight`

Исходный код функции `removeZ`, отвечающей за устранение примыканий и лишних скобок:

```go
func removeZw(s string) (string, int) {
	var oldS = ""
	for len(s) != len(oldS) {
		oldS = s
		s = internal.RemoveBrackets(s)
	}
	if len(s) == 0 || s[len(s)-1] != '*' {
		return s, NoStar
	}
	var countBrackets int
	for _, s := range s {
		if s == '(' || s == ')' {
			countBrackets++
		}
	}
	// скобок нет
	if countBrackets == 0 {
		if len(s) > 0 && s[len(s)-1] == '*' {
			return s[:len(s)-1], StarNoBrackets // убираем *
		}
	} else if len(s)-len(internal.RemoveBrackets(s[:len(s)-1])) == 3 {
		// есть внешняя скобка
		return s[1 : len(s)-2], StarBrackets // убираем скобки и *

	}
	return s, NoStar
}

```

Исходный код функции `Unions`, отвечающей за разбиение на множество объединений:

```go
func (str *RW) Unions() RWS {
	var lexemes = make(RWS, 0)
	var (
		bracketOn, begin int
	)
	for i, r := range *str {
		switch r {
		case '(':
			bracketOn++
		case ')':
			bracketOn--
		case '|':
			if bracketOn == 0 {
				lexemes = append(lexemes, (*str)[begin:i])
				begin = i + 1
			}
		}
	}
	if begin > 0 && begin < len(*str) {
		lexemes = append(lexemes, (*str)[begin:])
	}
	return lexemes
}
```

Исходный код функции `Concatenations`, отвечающей за разбиение на множество конкатенаций:

```go
func (str RW) Concatenations() []RW {
	var lexemes = make([]RW, 0)
	var (
		bracketOn, begin int
	)
	for i, r := range str {
		switch r {
		case ' ':
			continue
		case '(':
			if bracketOn == 0 {
				begin = i
			}
			bracketOn++
		case ')':
			bracketOn--
			if bracketOn == 0 {
				lexemes = append(lexemes, str[begin:i+1])
				begin = i + 1
			}
		case '*':
			if i > 0 && str[i-1] == ')' && bracketOn == 0 {
				begin = i + 1
			}
			lexemes[len(lexemes)-1] += "*"
		case '|':
			if bracketOn == 0 {
				return nil
			}
		default:
			if bracketOn == 0 {
				lexemes = append(lexemes, RW(r))
			}
		}
	}
	return lexemes
}
```

Исходный код функции `MultiplyEdge`, отвечающей за создание параллельных дуг:

```go
func (g *Graph) MultiplyEdge(e *Edge, newWeights ...string) []Edge {
	if len(newWeights) < 2 {
		return nil
	}
	var edges = make([]Edge, 0)
	for _, weight := range newWeights {
		e1 := g.AddEdge(&Edge{
			From:   e.From,
			Weight: "e",
		})
		e2 := g.AddEdge(&Edge{
			To:     e.To,
			Weight: "e",
		})
		newEdge := g.AddEdge(&Edge{
			From:   e1.To,
			To:     e2.From,
			Weight: weight,
		})
		edges = append(edges, newEdge)
	}
	g.RemoveEdge(e)

	return edges
}
```

Исходный код функции `SplitEdge`, отвечающей за добавление последовательности дуг:

```go
func (g *Graph) SplitEdge(e *Edge, newWeights ...string) []Edge {
	if len(newWeights) < 2 {
		return nil
	}
	var (
		edges    = make([]Edge, 0)
		prevEdge = Edge{
			To: e.From,
		}
	)
	for i, weight := range newWeights {
		var newEdge = &Edge{
			From:   prevEdge.To,
			Weight: weight,
		}
		if i == len(newWeights)-1 {
			newEdge.To = e.To
		}
		prevEdge = g.AddEdge(newEdge)
		edges = append(edges, prevEdge)
	}
	g.RemoveEdge(e)
	return edges
}

```

Исходный код функции `EpsilonEdge`, отвечающей за создание дуг в случае обнаружения замыкания:

```go
func (g *Graph) EpsilonEdge(e *Edge, weight string) *Edge {
	e1 := g.AddEdge(&Edge{
		From:   e.From,
		Weight: "e",
	})
	e2 := g.AddEdge(&Edge{
		To:     e.To,
		Weight: "e",
	})
	g.AddEdge(&Edge{
		To:     e1.To,
		From:   e2.From,
		Weight: "e",
	})
	g.AddEdge(&Edge{
		From:   e.From,
		To:     e.To,
		Weight: "e",
	})

	g.RemoveEdge(e)

	newEdge := g.AddEdge(&Edge{
		From:   e1.To,
		To:     e2.From,
		Weight: weight,
	})
	return &newEdge
}
```

Описанные выше методы являются частью алгоритма перевода РВ в НКА. Исходный код соответствующего метода приведен ниже.

```go

func (str *RW) ToENKA() *fsm.FSM {
	var kda = &fsm.FSM{Graph: graph.NewGraph()}
	firstEdge := kda.AddEdge(&graph.Edge{
		From:   "q0",
		To:     "q1",
		Weight: string(*str),
	})
	var queue = []graph.Edge{firstEdge}
	for len(queue) > 0 {

		head := queue[0]
		queue = queue[1:]

		weight := strings.TrimSpace(head.Weight)
		if len(weight) == 1 {
			continue
		}

		weight, changed := removeZw(weight)
		if changed != 0 {
			head = *kda.EpsilonEdge(&head, weight)
			queue = append(queue, head)
			continue
		}

		ew := RW(weight)
		rws := (&ew).Unions()
		edges := kda.MultiplyEdge(&head, rws.toString()...)
		queue = append(queue, edges...)
		if len(edges) != 0 {
			continue
		}

		rws = (&ew).Concatenations()
		edges = kda.SplitEdge(&head, rws.toString()...)
		queue = append(queue, edges...)
	}
	kda.SetFirstLast([]string{"q0"}, []string{"q1"})
	return kda
}
```

### Строительство эквивалентного ДКА по НКА

В основе алгоритма составления ДКА лежит алгоритм Томпсона[3]

Алгоритм Томпсона строит по [НКА](https://neerc.ifmo.ru/wiki/index.php?title=%D0%9D%D0%B5%D0%B4%D0%B5%D1%82%D0%B5%D1%80%D0%BC%D0%B8%D0%BD%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%BD%D1%8B%D0%B5_%D0%BA%D0%BE%D0%BD%D0%B5%D1%87%D0%BD%D1%8B%D0%B5_%D0%B0%D0%B2%D1%82%D0%BE%D0%BC%D0%B0%D1%82%D1%8B "Недетерминированные конечные автоматы") эквивалентный [ДКА](https://neerc.ifmo.ru/wiki/index.php?title=%D0%94%D0%B5%D1%82%D0%B5%D1%80%D0%BC%D0%B8%D0%BD%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%BD%D1%8B%D0%B5_%D0%BA%D0%BE%D0%BD%D0%B5%D1%87%D0%BD%D1%8B%D0%B5_%D0%B0%D0%B2%D1%82%D0%BE%D0%BC%D0%B0%D1%82%D1%8B "Детерминированные конечные автоматы") следующим образом:

* Помещаем в очередь Q множество, состоящее только из стартовой вершины.
* Затем, пока очередь не пуста выполняем следующие действия:
  * Достаем из очереди множество, назовем его q
  * Для всех c∈Σ посмотрим в какое состояние ведет переход по символу c из каждого состояния в q. Полученное множество состояний положим в очередь Q только если оно не лежало там раньше. Каждое такое множество в итоговом ДКА будет отдельной вершиной, в которую будут вести переходы по соответствующим символам.
  * Если в множестве qq хотя бы одна из вершин была терминальной в НКА, то соответствующая данному множеству вершина в ДКА также будет терминальной.

Реализация представлена в листинге ниже, функция `ToDka`

```go
func (fsm *FSM) ToDFA() *FSM {
	if len(fsm.Vertexes) == 0 {
		return fsm
	}
	var (
		visitedCombinations = make(map[string]bool, 0)
		newFSM              = &FSM{graph.NewGraph()}
		queue               = []NewVertex{{
			From: fsm.First[0],
			To:   fsm.EClosure(fsm.First...),
		}}
		lastVertexes []string
	)
	// Пока очередь не пуста
	for len(queue) != 0 {
		// Итерируемся по очереди
		head := queue[0]
		queue = queue[1:]

		// Определяем куда и по каким путям можно прийти отсюда
		var paths = fsm.MoveTo(head)

		// Проход по каждому направлению
		for path, vertexes := range paths {
			var (
				ids = make([]string, 0)
				// следим за тем, чтобы не потерять, какая
				// вершина была замыкающей
				withLast bool
			)
			for vertex := range vertexes {
				ids = append(ids, vertex)
				if !withLast {
					withLast = fsm.FindInString(vertex, fsm.Last)
				}
			}
			var id string

			ids = fsm.EClosure(ids...)

			// Сортируем, чтобы потом можно было пометить пройденный путь
			// "1,3,5" и "3,1,5" - одна и та же комбинация для нас
			sort.Strings(ids)
			id = strings.Join(ids, " ")
			newVertex := newFSM.AddVertex(graph.VertexOptID(id))
			newFSM.AddEdge(&graph.Edge{
				From:   head.From,
				To:     newVertex,
				Weight: path,
			})
			if withLast {
				lastVertexes = append(lastVertexes, newVertex)
			}

			// Помечаем пройденную комбинацию вершин
			_, ok := visitedCombinations[id]
			if ok {
				continue
			}

			visitedCombinations[id] = true

			queue = append(queue, NewVertex{
				To:   ids,
				From: newVertex,
			})
		}
	}
	// устаналиваем обновленные начало и конец
	newFSM.SetFirstLast(fsm.First, lastVertexes)
	*fsm = *newFSM
	return fsm
}

// Пути, куда можно попасть из текущей вершины
// Ключ - путь, значения - в каких узлы ведет
// вложенная мэпа гарантирует уникальность узлов
type Dtran map[string]map[string]bool

// Получить пути из вершины head
func (fsm *FSM) MoveTo(head NewVertex) Dtran {
	var paths = make(Dtran)
	for _, old := range head.To {
		toWhom := fsm.Vertexes[old].Out
		for _, e := range toWhom {
			_, ok := paths[e.Weight]
			if !ok {
				paths[e.Weight] = make(map[string]bool)
			}
			paths[e.Weight][e.To] = true
		}
	}
	for path, _ := range paths {
		if path == "e" {
			continue
		}
		for vertex := range paths["e"] {
			paths[path][vertex] = true
		}
	}
	delete(paths, "e")
	return paths
}
```

### Минимизация с помошью алгоритма Бржозовского

Введём следующие обозначения:

* `A`— конечный автомат,
* `d(A)` — детерминизированный автомат для `A`,
* `r(A)` — обратный автомат для AA,
* `dr(A)` — результат `d(r(A))`. Аналогично для `rdr(A)` и `drdr(A)`.

Пусть `A` — автомат (необязательно детерминированный), распознающий язык `L`. Минимальный   детерминированный автомат `A` может быть вычислен следующим образом: `A=drdr(A)`

Для детерминизации служит функция `D`, вызывающая  `ToDka`, описанную пунктом выше. А обратный автомат строит функцию `R`:

```go
// D Детерминизировать КА
func (A *DR) D() *DR {
	v := &DR{FSM{A.ToDFA().Beautify()}}
	*A = *v
	return A
}

// R Построить обратный КА
func (A *DR) R() *DR {
	var reverseMe = NewDR()
	for _, edge := range A.Edges {
		reverseMe.AddEdge(&graph.Edge{
			From:   edge.To,
			To:     edge.From,
			Weight: edge.Weight,
		})
	}

	reverseMe.SetFirstLast(A.Last, A.First)
	reverseMe = &DR{FSM{reverseMe.Beautify()}}
	*A = *reverseMe
	return A
}
```

С учетом этих функций алгоритм минимизации выглядит следующим образом: `A.R().D().R().D()`, где `A` - КА

## 4. Набор тестов

Для тестирования использовались примеры 3.21, 3.24 из [3] а также был написан ряд тестов, проверяющих корректное создание НКА по следующим регулярным выражениям:

- `a*|b`
- `((((ab))))`
- `((a|b)*b)`
- `((a|b)b*)`
- `((a|b)(b))`
- `((a|b)(b)*)`
- `((((a))*)*)`

Примеры для проверки минимизации были взяты из [4], [5], [6]

 Для запуска тестов введите команду `go test`. Результаты отобразятся в консоли:

![alt text](assets/reports/1.png)

В случае падения тестов, в папке `assets/test` будут созданы файлы формата `.dot` для визуального сравнения ожидаемых(`expected.dot`) и полученных результатов(`real.dot`).

## 5. Результаты выполнения программы

Интерфейс программы выглядит следующим образом:

![alt text](assets/reports/2.png)

## 6. Выводы

В ходе выполнения лабораторной работы были выполнены следующие задачи:

1) Были изучены основные понятия и определения, лежащие в основе построения лексических анализаторов.
2) Проведен анализ связи между регулярным множеством, регулярным выражением, праволинейным языком, конечно - автоматным языком и недетерминированным конечно-автоматным языком.
3) Разработана, протестирована и отлажена программа распознавания цепочек регулярного или праволинейного языка в соответствии с предложенным вариантом грамматики.

## Список дополнительной использованной литературы

1. БЕЛОУСОВ А.И., ТКАЧЕВ С.Б. Дискретная математика: Учеб. Для вузов / Под ред. В.С. Зарубина, А.П.
   Крищенко. – М.: Изд-во МГТУ им. Н.Э. Баумана, 2001.
2. АХО А., УЛЬМАН Дж. Теория синтаксического анализа, перевода и компиляции: В 2-х томах. Т.1.:
   Синтаксичечкий анализ. - М.: Мир, 1978.
3. АХО А.В, ЛАМ М.С., СЕТИ Р., УЛЬМАН Дж.Д. Компиляторы: принципы, технологии и инструменты. – М.:
   Вильямс, 2008.
4. "Алгоритм Бржзовского", университет ИТМО, URL: http://neerc.ifmo.ru/wiki/index.php?title=Алгоритм_Бржозовского
5. "Минимизация ДКА, алгоритм за O(n^2) с построением пар различимых состояний", университет ИТМО, URL:университет ИТМО, URL: http://neerc.ifmo.ru/wiki/index.php?title=Минимизация_ДКА,_алгоритм_за_O(n%5E2)_с_построением_пар_различимых_состояний
6. "Скорость переработки и передачи данных", стр.2, URL:университет ИТМО, URL: https://lektsii.org/6-91118.html
