package ast

import (
	"fmt"
	"lab2/internal/g5"
	"log"
	"os"

	"github.com/awalterschulze/gographviz"
)

func ToNumOperator(r g5.Rule) (NumOperator, error) {
	var (
		terms, nonTerms []string
	)
	for _, s := range r.Symbols {
		if s.Type == g5.Term {
			terms = append(terms, s.Value)
		} else {
			nonTerms = append(nonTerms, s.Value)
		}
	}
	switch {
	case len(terms) == 1 && len(nonTerms) == 2:
		return OneTwoOperatored{
			Main: terms[0],
		}, nil
	case len(terms) == 1 && len(nonTerms) == 0:
		return NoOperatored{
			Main: terms[0],
		}, nil
	case len(r.Symbols) == 6:
		if IsIfThenElseOperator(r) {
			return IfThenElseOperatored{}, nil
		}
	case len(r.Symbols) == 4:
		if IsIfThenOperator(r) {
			return IfThenOperatored{}, nil
		}
	}
	return nil, fmt.Errorf("нет модели для правила с %d термами и %d нетермами: %v", len(terms), len(nonTerms), r)
}

func IsIfThenElseOperator(r g5.Rule) bool {
	if len(r.Symbols) != 6 {
		return false
	}
	return r.Symbols[0].Type == g5.Term && r.Symbols[0].Value == "if" &&
		r.Symbols[1].Type == g5.NonTerm &&
		r.Symbols[2].Type == g5.Term && r.Symbols[2].Value == "then" &&
		r.Symbols[3].Type == g5.NonTerm &&
		r.Symbols[4].Type == g5.Term && r.Symbols[4].Value == "else" &&
		r.Symbols[5].Type == g5.NonTerm
}

func IsIfThenOperator(r g5.Rule) bool {
	if len(r.Symbols) != 4 {
		return false
	}
	return r.Symbols[0].Type == g5.Term && r.Symbols[0].Value == "if" &&
		r.Symbols[1].Type == g5.NonTerm &&
		r.Symbols[2].Type == g5.Term && r.Symbols[2].Value == "then" &&
		r.Symbols[3].Type == g5.NonTerm
}

type Node struct {
	ID          string
	Value       string
	Parent      *Node
	ParentValue string
	Type        string
}

type NumOperator interface {
	/*
		 Возврщает 2 слайса узлов:
		 	- те, которые надо добавить в массив вакантных
			 узлов, для дальнейшнего назначения им значения
			 и помещения в АСТ
			- те, которые сразу добавляются в АСТ
	*/
	ToNodes(node *Node, counter *int) ([]*Node, []*Node)
}

// +, *...
type OneTwoOperatored struct {
	Main string
}

func (two OneTwoOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = two.Main

	var leftNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: two.Main,
	}
	*counter++
	var rightNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: two.Main,
	}
	*counter++
	return []*Node{leftNode, rightNode}, nil
}

// if A then B
type IfThenOperatored struct{}

func (two IfThenOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = "branch"
	node.Type = g5.NonTerm

	var conditionNodeParent = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Value:       "if",
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++
	var trueNodeParent = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Value:       "then",
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++

	var conditionNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      conditionNodeParent,
		ParentValue: conditionNodeParent.Value,
	}
	*counter++
	var trueNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      trueNodeParent,
		ParentValue: trueNodeParent.Value,
	}
	*counter++
	return []*Node{
			conditionNode,
			trueNode,
		}, []*Node{
			conditionNodeParent,
			trueNodeParent,
		}
}

// if A then B else C
type IfThenElseOperatored struct{}

func (two IfThenElseOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = "branch"
	node.Type = g5.NonTerm

	var conditionNodeParent = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Value:       "if",
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++
	var trueNodeParent = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Value:       "then",
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++
	var falseNodeParent = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Value:       "else",
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++

	var conditionNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      conditionNodeParent,
		ParentValue: conditionNodeParent.Value,
	}
	*counter++
	var trueNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      trueNodeParent,
		ParentValue: trueNodeParent.Value,
	}
	*counter++
	var falseNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      falseNodeParent,
		ParentValue: falseNodeParent.Value,
	}
	*counter++
	return []*Node{
			conditionNode,
			trueNode,
			falseNode,
		}, []*Node{
			conditionNodeParent,
			trueNodeParent,
			falseNodeParent,
		}
}

// A -> a
type NoOperatored struct {
	Main string
}

func (no NoOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = no.Main
	return nil, nil
}

// VisualizeFSM - визуализировать граф
func Visualize(nodes []*Node, path, name string) error {
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

		// if v.Type == g5.Term {
		// 	vattr["label"] = fmt.Sprintf(`<<font color="green">%s</font>>`, v.Value)
		// } else {
		vattr["label"] = fmt.Sprintf(`<<font color="red">%s</font>>`, v.Value)
		//}
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
