package visualizer

import (
	"fmt"
	"gocompiler/internal/fsm"
	"log"
	"os"

	"github.com/awalterschulze/gographviz"
)

// MustVisualizeFSM - попробовать визуализировать граф
//  если не выйдет, то бросаем фатальную ошибку
func MustVisualizeFSM(g *fsm.FSM, path string) {
	err := VisualizeFSM(g, path)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

// VisualizeFSM - визуализировать граф
func VisualizeFSM(g *fsm.FSM, path string) error {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	var attrs = make(map[string]string, 0)
	for _, v := range g.Vertexes {
		graph.AddNode("G", toString(v.ID), nil)
	}
	for _, e := range g.Edges {
		attrs["label"] = fmt.Sprintf(`<<font color="blue">%s</font>>`, e.Weight)
		graph.AddEdge(toString(e.From), toString(e.To), true, attrs)
	}
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(graph.String())

	return nil
}

func toString(s string) string {
	return `"` + s + `"`
}
