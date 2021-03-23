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

type DKAVertex struct {
	From string
	To   []string
}

// ToDka - построение эквивалентного ДКА к НКА
// http://esyr.org/wiki/Конструирование_Компиляторов%2C_Алгоритмы_решения_задач#.D0.9F.D0.BE.D1.81.D1.82.D1.80.D0.BE.D0.B5.D0.BD.D0.B8.D0.B5_.D0.94.D0.9A.D0.90_.D0.BF.D0.BE_.D0.9D.D0.9A.D0.90
func (fsm *FSM) ToDka() *FSM {
	//log.Println("Построить ДКА, эквивалентное указанному НКА")
	if len(fsm.Vertexes) == 0 {
		return fsm
	}
	var (
		visitedCombinations = make(map[string]bool, 0)
		newFSM              = &FSM{graph.NewGraph()}
		queue               = []DKAVertex{
			{
				To:   fsm.First,
				From: fsm.First[0],
			},
		}
		lastVertexes []string
	)
	for len(queue) != 0 {
		head := queue[0]

		// ключ - путь, значения - в каких узлы ведет
		// вложенная мэпа, чтобы гарантировать уникальность узлов
		var paths = make(map[string]map[string]bool, 0)
		for _, old := range head.To {
			toWhom := fsm.Vertexes[old].Out
			for _, e := range toWhom {
				_, ok := paths[e.Weight]
				if !ok {
					paths[e.Weight] = make(map[string]bool, 0)
				}
				paths[e.Weight][e.To] = true
			}
		}
		for path, vertexes := range paths {
			var (
				ids      = make([]string, 0)
				withLast bool
			)
			for vertex := range vertexes {
				ids = append(ids, vertex)
				if !withLast {
					withLast = fsm.FindInString(vertex, fsm.Last)
				}
			}
			var id string
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

			_, ok := visitedCombinations[id]
			if ok {
				continue
			}

			visitedCombinations[id] = true

			queue = append(queue, DKAVertex{
				To:   ids,
				From: newVertex,
			})
		}

		queue = queue[1:]
	}
	newFSM.SetFirstLast(fsm.First, lastVertexes)
	*fsm = *newFSM
	return fsm
}

// EClosure - возвращает все вершины, в которые можно попасть
// по епсилон переходам стартовав из  vertexes
func (fsm FSM) EClosure(vertexes ...string) []string {
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
				unique[edge.To] = true
				stack = append(stack, edge.To)
				visited = append(visited, edge.To)
			}
		}
	}
	return visited
}

// ToDFA - построение эквивалентного ДКА к НКА
// Алгоритм 3.20
func (fsm *FSM) ToDFA() *FSM {
	//log.Println("Построить ДКА, эквивалентное указанному НКА")
	if len(fsm.Vertexes) == 0 {
		return fsm
	}
	var (
		visitedCombinations = make(map[string]bool, 0)
		newFSM              = &FSM{graph.NewGraph()}
		queue               = []DKAVertex{{
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

			queue = append(queue, DKAVertex{
				To:   ids,
				From: newVertex,
			})
		}
	}
	// устаналиваем начало и обновленный конец
	newFSM.SetFirstLast(fsm.First, lastVertexes)
	*fsm = *newFSM
	return fsm
}

// Пути, куда можно попасть из текущей вершины
// Ключ - путь, значения - в каких узлы ведет
// вложенная мэпа гарантирует уникальность узлов
type Dtran map[string]map[string]bool

// Получить пути из вершины head
func (fsm *FSM) MoveTo(head DKAVertex) Dtran {
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
