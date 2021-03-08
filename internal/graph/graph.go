package graph

import (
	"fmt"
	"log"
)

// Vertex - структура вершин
type Vertex struct {
	ID  string
	Out map[string]*Edge
	In  map[string]*Edge
}

// Edge - структура ребра
type Edge struct {
	From   string
	To     string
	Weight string
}

// ID - получить уникальный айди, по которому можно
// сравнивать ребра
func (e *Edge) ID() string {
	return e.From + " " + e.To + " " + e.Weight
}

// Graph - структура графа
type Graph struct {
	Vertexes map[string]*Vertex
	Edges    map[string]*Edge

	// счётчик для автоматического нумерования узлов
	counter int
	// первый узел графа
	First []string
	// первый узел графа
	Last []string
}

// NewGraph - создать новый граф
func NewGraph() *Graph {
	return &Graph{
		Vertexes: make(map[string]*Vertex, 0),
		Edges:    make(map[string]*Edge, 0),
	}
}

type VertexOpts struct {
	changeName func(string) string
	id         string
}

type VertexOpt func(*VertexOpts)

func VertexOptChangeName(changeName func(string) string) VertexOpt {
	return func(o *VertexOpts) {
		o.changeName = changeName
	}
}

func VertexOptID(id string) VertexOpt {
	return func(o *VertexOpts) {
		o.id = id
	}
}

// AddVertex - добавить вершину
func (g *Graph) AddVertex(opts ...VertexOpt) string {
	options := &VertexOpts{}
	for _, opt := range opts {
		opt(options)
	}
	var id string
	if options.changeName != nil {
		id = options.changeName(options.id)
	}
	if id == "" {
		if options.id != "" {
			id = options.id
		} else {
			id = fmt.Sprintf("s%d", g.counter)
		}
	}
	_, ok := g.Vertexes[id]
	if !ok {
		v := &Vertex{
			ID:  id,
			Out: make(map[string]*Edge, 0),
			In:  make(map[string]*Edge, 0),
		}
		g.Vertexes[id] = v

		g.counter++
	}
	return id
}

// SetFirstLast - установить первый и последний узлы
func (g *Graph) SetFirstLast(f, l []string) {
	g.First = f
	g.Last = l

}

// SetFirstLast - установить первый и последний узлы
func (g *Graph) SwitchFirstLast() {
	g.First, g.Last = g.Last, g.First
}

// // GetFirst - получить исток
// func (g *Graph) GetFirst() *Vertex {
// 	// log.Println("g.firstVertex", g.firstVertex, g.firstVertex.countIn())
// 	// if g.firstVertex != nil && g.firstVertex.countIn() == 0 {
// 	// 	return g.firstVertex
// 	// }
// 	// //var vertex = &Vertex{}
// 	// for _, v := range g.Vertexes {
// 	// 	vertex = v // Нельзя оставить вершину не проинициализированной
// 	// 	if v.countIn() == 0 {
// 	// 		return v
// 	// 	}
// 	// }
// 	return g.First
// }

// GetFirst - получить исток
func (v Vertex) countIn() int {
	var c int
	for _, edge := range v.In {
		if edge.From == edge.To {
			continue
		}
		c++
	}
	return c
}

// VertexesArr - получить список вершин
func (g *Graph) VertexesArr() []*Vertex {
	var (
		vertexes = make([]*Vertex, len(g.Vertexes))
		i        int
	)
	for _, v := range g.Vertexes {
		vertexes[i] = v
		i++
	}
	return vertexes
}

// AddEdge - добавить ребро
func (g *Graph) AddEdge(ce *Edge, opts ...VertexOpt) {
	var e = &Edge{
		From:   ce.From,
		To:     ce.To,
		Weight: ce.Weight,
	}
	e.From = g.AddVertex(append(opts, VertexOptID(e.From))...)
	e.To = g.AddVertex(append(opts, VertexOptID(e.To))...)
	id := e.ID()

	g.Edges[id] = e
	g.Vertexes[e.From].Out[id] = e
	g.Vertexes[e.To].In[id] = e
}

// RemoveVertex - убрать вершину
func (g *Graph) RemoveVertex(id string) {
	for _, e := range g.Edges {
		if e.From == id || e.To == id {
			g.RemoveEdge(e)
		}
	}
	delete(g.Vertexes, id)
}

// RemoveEdge - убрать ребро
func (g *Graph) RemoveEdge(e *Edge) {
	id := e.ID()
	delete(g.Vertexes[e.From].Out, id)
	delete(g.Vertexes[e.To].In, id)
	delete(g.Edges, id)
}

// SplitEdge - разбить ребро на несколько
func (g *Graph) SplitEdge(e *Edge, newWeights ...string) {
	if len(newWeights) == 0 {
		return
	}
	g.RemoveEdge(e)
	var prevEdge = &Edge{
		To: e.From,
	}
	for i, weight := range newWeights {
		var newEdge = &Edge{
			From:   prevEdge.To,
			Weight: weight,
		}
		if i == len(newWeights)-1 {
			newEdge.To = e.To
		}
		g.AddEdge(newEdge)
		prevEdge = newEdge
	}
}

// MultiplyEdge создать несколько ребер вместо одного
func (g *Graph) MultiplyEdge(e *Edge, newWeights ...string) {
	if len(newWeights) < 2 {
		return
	}
	g.RemoveEdge(e)

	for _, weight := range newWeights {
		g.AddEdge(&Edge{
			From:   e.From,
			To:     e.To,
			Weight: weight,
		})
	}
}

func (g *Graph) FindInString(find string, ids []string) bool {
	for _, id := range ids {
		if id == find {
			return true
		}
	}
	return false
}

// Beautify установить красивые названия
func (g *Graph) Beautify() *Graph {
	//return g
	// last и first тоже заменять
	log.Println("before", len(g.First), len(g.Last))
	var (
		namesReplacer = make(map[string]string, 0)
		newFSM        = NewGraph()
		i             int
	)
	for _, v := range g.Vertexes {
		namesReplacer[v.ID] = fmt.Sprintf("%d", i+1)
		i++
	}
	for _, e := range g.Edges {
		newFSM.AddEdge(&Edge{
			From:   namesReplacer[e.From],
			To:     namesReplacer[e.To],
			Weight: e.Weight,
		})
	}
	for _, str := range g.First {
		newFSM.First = append(newFSM.First, namesReplacer[str])
	}
	for _, str := range g.Last {
		newFSM.Last = append(newFSM.Last, namesReplacer[str])
	}
	//newFSM.SetFirstLast(namesReplacer[g.First], namesReplacer[g.Last])
	*g = *newFSM
	log.Println("after", len(g.First), len(g.Last))
	return newFSM
}
