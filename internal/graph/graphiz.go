package graph

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

func NewGraphiz(g Graph) string {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	var attrs = make(map[string]string, 0)
	for _, v := range g.Vertexes {
		graph.AddNode("G", v.ID, nil)
	}
	for _, e := range g.Edges {
		attrs["label"] = fmt.Sprintf(`<<font color="blue">%s</font>>`, e.Weight)
		graph.AddEdge(e.From, e.To, true, attrs)
	}
	return graph.String()
}
