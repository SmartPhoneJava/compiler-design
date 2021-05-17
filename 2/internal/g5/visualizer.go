package g5

import (
	"fmt"
	"log"
	"os"

	"github.com/awalterschulze/gographviz"
)

// MustVisualizeFSM - попробовать визуализировать граф
//  если не выйдет, то бросаем фатальную ошибку
func MustVisualize(g []*Node, path, name string) {
	err := VisualizeFSM(g, path, name)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

// VisualizeFSM - визуализировать граф
func VisualizeFSM(nodes []*Node, path, name string) error {
	graphAst, err := gographviz.ParseString(`digraph G {}`)
	if err != nil {
		return err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return err
	}

	var attrs = make(map[string]string, 0)
	for _, v := range nodes {
		log.Println("node", v.ID+v.Value, v.Type)
		var vattr = make(map[string]string, 0)

		if v.Type == Term {
			vattr["label"] = fmt.Sprintf(`<<font color="green">%s</font>>`, v.Value)
		} else {
			vattr["label"] = fmt.Sprintf(`<<font color="red">%s</font>>`, v.Value)
		}
		graph.AddNode("G", toString(v.ID), vattr)
	}
	for _, e := range nodes {
		//attrs["label"] = fmt.Sprintf(`<<font color="blue">%s</font>>`, "hello")
		if e.Parent != nil {
			log.Println("edge", e.ID+e.Value, e.Parent.ID+e.Parent.Value, e.ParentValue)
			graph.AddEdge(toString(e.Parent.ID), toString(e.ID), true, attrs)
		}
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
