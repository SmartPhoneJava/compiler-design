package fsm

import (
	"gocompiler/internal/graph"
	"log"
	"sort"
	"strings"
)

// FSM НКА
type FSM struct {
	*graph.Graph
}

// RemoveShortCircuits убрать замыкания
func (fsm *FSM) RemoveShortCircuits() *FSM {
	for _, e := range fsm.Edges {
		if e.Weight == "*" {
			vIn := fsm.Vertexes[e.From].In
			for _, ei := range vIn {
				var newEdge = graph.Edge{
					From:   ei.From,
					To:     ei.To,
					Weight: "e",
				}
				fsm.RemoveEdge(ei)
				fsm.AddEdge(&newEdge)
				fsm.AddEdge(&graph.Edge{
					From:   ei.To,
					To:     ei.To,
					Weight: ei.Weight,
				})
				break
			}
			e.Weight = "e"
		}
	}
	for _, v := range fsm.Vertexes {
		if len(v.In)+len(v.Out) == 0 {
			fsm.RemoveVertex(v.ID)
		}
	}
	return fsm
}

// ReplaceEpsilons заменить епсилон-переходы
func (fsm *FSM) ReplaceEpsilons() *FSM {
	log.Println("Заменить епсилон-переходы")
	var newFSM = &FSM{graph.NewGraph()}
	for _, v := range fsm.Vertexes {
		m := fsm.replaceEpsilons(v, "e", false)

		for k, val := range m {
			if k == v.ID {
				continue
			}
			newFSM.AddEdge(&graph.Edge{
				From:   v.ID,
				To:     k,
				Weight: val,
			})
		}

		// Добавляем в новый граф завихренные ребра(на самого себя)
		for _, edge := range v.Out {
			if edge.From == edge.To {
				newFSM.AddEdge(&graph.Edge{
					From:   edge.From,
					To:     edge.To,
					Weight: edge.Weight,
				})
			}
		}
	}

	*fsm = *newFSM
	return fsm
}

// ReplaceEqualEdges - убрать ребра-дубли
func (fsm *FSM) ReplaceEqualEdges() *FSM {
	log.Println("Убрать ребра дубли")
	var (
		removeVertexes = make([]*graph.Vertex, 0)
		vertexCount    = len(fsm.Vertexes)
		vertexes       = fsm.VertexesArr()
	)
	for i := 0; i < vertexCount; i++ {
		for j := i + 1; j < vertexCount; j++ {
			var (
				v1 = vertexes[i]
				v2 = vertexes[j]
			)
			if len(v1.Out) != len(v2.Out) {
				continue
			}
			var m = make(map[string]bool, 0)

			for _, edge := range v1.Out {
				m[edge.To] = true
			}
			for _, edge := range v2.Out {
				delete(m, edge.To)
			}
			// Полное совпадение исходящих дуг!
			if len(m) == 0 {
				// Добавляем к одной из вершин входящие дуги второй

				var m = make(map[string]*graph.Edge, 0)

				for _, edge := range v1.In {
					m[edge.From] = edge
				}
				for _, edge := range v2.In {
					delete(m, edge.From)
				}
				for k, v := range m {
					fsm.AddEdge(&graph.Edge{
						From:   k,
						To:     v2.ID,
						Weight: v.Weight,
					})
				}

				log.Println("Следующие состояния эквивалентны", v1.ID, v2.ID)

				//  Первая вершина нам больше не нужна, но мы
				// не можем ее сразу удалить поскольку итерируемся
				// по списку вершин, поэтому удалим ее попозже
				removeVertexes = append(removeVertexes, v1)
			}
		}
	}

	for _, v := range removeVertexes {
		fsm.RemoveVertex(v.ID)
	}
	return fsm
}

type DKAVertex struct {
	From string
	Olds []string
}

// ToDka - построение эквивалентного ДКА к НКА
// http://esyr.org/wiki/Конструирование_Компиляторов%2C_Алгоритмы_решения_задач#.D0.9F.D0.BE.D1.81.D1.82.D1.80.D0.BE.D0.B5.D0.BD.D0.B8.D0.B5_.D0.94.D0.9A.D0.90_.D0.BF.D0.BE_.D0.9D.D0.9A.D0.90
func (fsm *FSM) ToDka() *FSM {
	log.Println("Построить ДКА, эквивалентное указанному НКА")
	if len(fsm.Vertexes) == 0 {
		return fsm
	}
	var (
		visitedCombinations = make(map[string]bool, 0)
		newFSM              = &FSM{graph.NewGraph()}
		queue               = []DKAVertex{
			{
				Olds: fsm.First,
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
		for _, old := range head.Olds {
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
				Olds: ids,
				From: newVertex,
			})
		}

		queue = queue[1:]
	}
	newFSM.SetFirstLast(fsm.First, lastVertexes)
	*fsm = *newFSM
	return fsm
}

func (fsm *FSM) replaceEpsilons(
	v *graph.Vertex,
	path string,
	fromE bool,
) map[string]string {
	var (
		m    = make(map[string]string, 0)
		vOut = v.Out
	)
	for _, ei := range vOut {
		var kr map[string]string
		if path == "e" {
			if ei.Weight != "e" {
				// если раньше были только епсилоны и тут появилс
				// обычный символ
				m[ei.To] = ei.Weight
				// при вызове рекурсии необходимо пометить, что путь
				// создан из епсилона
			}
			kr = fsm.replaceEpsilons(fsm.Vertexes[ei.To], ei.Weight, ei.Weight != path)
			for k, r := range kr {
				m[k] = r
			}
		} else if ei.Weight == "e" || ei.Weight == path {
			// Если путь не состоит из епсилонов, то не стоит
			// заходить в самого себя, ведь это вызовет бесконечный
			// цикл
			// Если путь создан из епсилона, то мы больше не можем
			// ходить по ребрам с епсилонами
			if ei.From == ei.To || (fromE && ei.Weight == "e") {
				continue
			}
			m[ei.To] = path
			kr = fsm.replaceEpsilons(fsm.Vertexes[ei.To], path, fromE)
			for k, r := range kr {
				m[k] = r
			}
		}
	}
	// Раз мы добрались до этой вершины, пометим как мы это сделали
	if path != "e" {
		m[v.ID] = path
	}
	return m
}
