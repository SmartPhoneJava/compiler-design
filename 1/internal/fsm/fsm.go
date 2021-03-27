package fsm

import (
	"gocompiler/internal/graph"
	"sort"
	"strings"
)

// FSM НКА
type FSM struct {
	*graph.Graph
}

// NewVertex - переход к новой вершине для построения ДКА из НКА
type NewVertex struct {
	From string
	To   []string
}

// IsInLast проверить, что данная вершина является терминальной
func (fsm FSM) IsInLast(v string) bool {
	for _, l := range fsm.Last {
		if v == l {
			return true
		}
	}
	return false
}

// EClosure - возвращает все вершины, в которые можно попасть
// по епсилон переходам стартовав из  vertexes
// Если переход ведет в терминальную вершину, она добавляется
// в массив терминальных вершин
func (fsm FSM) EClosure(
	addToLast func(),
	vertexes ...string,
) []string {
	var (
		stack, visited []string
		unique         = make(map[string]bool)
	)
	for _, v := range vertexes {
		_, ok := unique[v]
		if ok {
			continue
		}
		unique[v] = true
		visited = append(visited, v)
		stack = append(stack, v)
	}
	for len(stack) > 0 {
		var head = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, edge := range fsm.Edges {
			if edge.From == head && edge.Weight == "e" {
				_, ok := unique[edge.To]
				if ok {
					continue
				}
				if fsm.IsInLast(edge.To) {
					addToLast()
				}

				unique[edge.To] = true
				stack = append(stack, edge.To)
				visited = append(visited, edge.To)
			}
		}
	}
	return visited
}

// ContainString - содержит ли строку
func (fsm *FSM) ContainString(str string) bool {
	str = strings.TrimSpace(str)
	var vertex = fsm.First[0]

	for len(str) > 0 {
		var found bool
		for _, edge := range fsm.Edges {
			if edge.From == vertex && edge.Weight == string(str[0]) {
				vertex = edge.To
				found = true
				break
			}
		}
		str = str[1:]
		if !found {
			return false
		}
	}
	for _, l := range fsm.Last {
		if l == vertex {
			return true
		}
	}
	return false
}

// ToDFA - построение эквивалентного ДКА к НКА
// Алгоритм 3.20
func (fsm *FSM) ToDFA() *FSM {
	if len(fsm.Vertexes) == 0 {
		return fsm
	}
	var (
		visitedCombinations = make(map[string]bool, 0)
		newFSM              = &FSM{graph.NewGraph()}
		queue               = []NewVertex{{
			From: fsm.First[0],
			To:   fsm.First,
		}}
		lastVertexes []string
	)
	// Пока очередь не пуста
	for len(queue) != 0 {
		// Итерируемся по очереди
		head := queue[0]
		head.To = fsm.EClosure(func() { lastVertexes = append(lastVertexes, head.From) }, head.To...)
		queue = queue[1:]

		// Определяем куда и по каким путям можно прийти отсюда
		var paths = fsm.MoveTo(head)

		if len(paths) == 0 {
			lastVertexes = append(lastVertexes, head.From)
		}

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
			ids = fsm.EClosure(func() { withLast = true }, ids...)

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

	delete(paths, "e")
	return paths
}
