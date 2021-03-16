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
func MustVisualizeFSM(g *fsm.FSM, path, name string) {
	err := VisualizeFSM(g, path, name)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

// MustVisualizeDR - попробовать визуализировать граф
//  если не выйдет, то бросаем фатальную ошибку
func MustVisualizeDR(g *fsm.DR, path, name string) {
	err := VisualizeFSM(&g.FSM, path, name)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

// VisualizeFSM - визуализировать граф
func VisualizeFSM(g *fsm.FSM, path, name string) error {
	graphAst, err := gographviz.ParseString(`digraph G {}`)
	if err != nil {
		return err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return err
	}

	var attrs = make(map[string]string, 0)
	for _, v := range g.Vertexes {
		var vattr = make(map[string]string, 0)

		if g.FindInString(v.ID, g.First) {
			vattr["label"] = fmt.Sprintf(`<<font color="green">%s</font>>`, v.ID)
		} else if g.FindInString(v.ID, g.Last) {
			vattr["label"] = fmt.Sprintf(`<<font color="red">%s</font>>`, v.ID)
		}
		graph.AddNode("G", toString(v.ID), vattr)
	}
	for _, e := range g.Edges {
		attrs["label"] = fmt.Sprintf(`<<font color="blue">%s</font>>`, e.Weight)
		graph.AddEdge(toString(e.From), toString(e.To), true, attrs)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(path + "/" + name)
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
