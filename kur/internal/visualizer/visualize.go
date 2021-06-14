package visualizer

import (
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
)

type Node struct {
	Name  string
	Style func() map[string]string
}

type Edge struct {
	From, To, FromPort string
	Style              func() map[string]string
}

// VisualizeFSM - визуализировать таблицу
func Visualize(nodes []*Node, edges []*Edge, path, name string) error {
	graphAst, err := gographviz.ParseString(`digraph G {}`)
	if err != nil {
		return err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return err
	}

	var nodesMap = make(map[string]*Node)
	for _, v := range nodes {
		nodesMap[v.Name] = v
	}

	for name, node := range nodesMap {
		graph.AddNode("G", toString(name), node.Style())
	}
	var edgesMap = make(map[string]map[string]interface{})
	for _, e := range edges {
		if e.To == "" {
			continue
		}
		_, ok := edgesMap[e.From]
		if !ok {
			edgesMap[e.From] = make(map[string]interface{})
		}
		_, ok = edgesMap[e.From][e.To]
		if ok {
			continue
		}
		edgesMap[e.From][e.To] = nil
		fromEdge := toString(e.From)

		graph.AddEdge(fromEdge, toString(e.To), true, e.Style())

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

	str := graph.String()
	for _, e := range edges {
		if e.FromPort != "" {
			str = strings.Replace(str,
				toString(e.From)+"->"+toString(e.To),
				toString(e.From)+":"+e.FromPort+"->"+toString(e.To),
				1)
		}
	}
	file.WriteString(str)

	return nil
}

func toString(s string) string {
	return `"` + s + `"`
}
