package ast

import (
	"fmt"
	"lab2/internal/g5"
	"os"

	"github.com/awalterschulze/gographviz"
)

// Типы генератора узлов для АСТ
const (
	Type1to2 = iota + 1
	Type1to1
	Type2to1
	TypeHard
	TypeTerm
)

/*
	ToNumOperator - определить модель, по которой строить узлы для АСТ
	на основе правила

	symbols - слайс символов, которые будут подставлены в
		зарезервированные термы
	anyCounter - счётчик неиспользованных зарезервированных
		символов - индекс для получения элемента массива symbols
	left - левая часть правила
	right - правая часть правила
	start - стартовый нетерм грамматики
*/
func ToNumOperator(
	symbols []string,
	left, start string,
	right g5.Rule,
	anyCounter *int,
) (NumOperator, error) {
	var (
		terms, nonTerms []string
	)
	for _, s := range right.Symbols {
		switch s.Type {
		case g5.Term:
			terms = append(terms, s.Value)
		case g5.Reserved:
			if *anyCounter >= 0 {
				terms = append(terms, symbols[*anyCounter])
				*anyCounter--
			}
		default:
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
	case len(terms) == 1 && len(nonTerms) == 1:
		if terms[0] == ";" {
			return IgnoreOperatored{}, nil
		}
		return OneOneOperatored{
			Main: terms[0],
		}, nil
	case len(right.Symbols) == 6:
		if IsIfThenElseOperator(right) {
			return IfThenElseOperatored{}, nil
		}
	case len(right.Symbols) == 4:
		if IsIfThenOperator(right) {
			return IfThenOperatored{}, nil
		}
	case len(terms) == 2 && len(nonTerms) == 1:
		if terms[1] == "=" {
			return LeftTemEqOperatored{
				Left: terms[0],
			}, nil
		}
		return TwoOneOperatored{
			Left:  terms[0],
			Right: terms[1],
		}, nil
	}
	return nil, fmt.Errorf("нет модели для правила с %d термами и %d нетермами: %s",
		len(terms), len(nonTerms), right.String(left, start),
	)
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
	Type        uint
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
	node.Type = Type1to2

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

// :+, :-
type OneOneOperatored struct {
	Main string
}

func (two OneOneOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = two.Main
	node.Type = Type1to1

	var downNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: two.Main,
	}
	*counter++
	return []*Node{downNode}, nil
}

// { E }, ( E ) ...
type TwoOneOperatored struct {
	Left, Right string
}

func (two TwoOneOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = two.Left + " " + two.Right
	node.Type = Type2to1

	var downNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++
	return []*Node{downNode}, nil
}

// { E }, ( E ) ...
type LeftTemEqOperatored struct {
	Left string
}

func (two LeftTemEqOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = "="
	node.Type = Type1to2

	var rightNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: node.Value,
	}
	*counter++
	var leftNode = &Node{
		ID:          fmt.Sprintf("%d.", *counter),
		Parent:      node,
		ParentValue: node.Value,
		Value:       two.Left,
		Type:        TypeTerm,
	}
	*counter++
	return []*Node{rightNode}, []*Node{leftNode}
}

// if A then B
type IfThenOperatored struct{}

func (two IfThenOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = "branch"
	node.Type = TypeHard

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
	node.Type = TypeHard

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

// Для ;
type IgnoreOperatored struct {
	Main string
}

func (no IgnoreOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = no.Main
	return []*Node{node}, nil
}

// A -> a
type NoOperatored struct {
	Main string
}

func (no NoOperatored) ToNodes(
	node *Node, counter *int,
) ([]*Node, []*Node) {
	node.Value = no.Main
	node.Type = TypeTerm
	return nil, nil
}

/*
 visualize - визуализировать граф
	nodes - узлы АСТ
	path - путь до файла формата .dot
	name - имя выходного файла
*/
func visualize(nodes []*Node, path, name string) error {
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
		var vattr = make(map[string]string, 0)
		var color = "black"
		switch v.Type {
		case Type1to2:
			color = "blue"
		case Type1to1:
			color = "cyan"
		case Type2to1:
			color = "orange"
		case TypeHard:
			color = "brown"
		case TypeTerm:
			color = "red"
		}
		vattr["label"] = fmt.Sprintf(`<<font color="%s">%s</font>>`, color, v.Value)
		graph.AddNode("G", toString(v.ID), vattr)
	}
	for _, e := range nodes {
		if e.Parent != nil {
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

// ToAst привести правила к АСТ
func Visualize(
	nt, start string,
	symbols []string,
	rules g5.Rules,
	path, name string,
) error {
	var (
		counter = 2
		root    = &Node{
			ID: fmt.Sprintf("%d.", 1),
		}
		freeNodes       = []*Node{root}
		nodes           = []*Node{root}
		alreadyAddNodes = make(map[string]interface{})
	)
	var anyCounter = len(symbols) - 1
	for i := len(rules) - 1; i >= 0; i-- {
		var r = rules[i]
		var model, err = ToNumOperator(symbols, nt, start, r, &anyCounter)
		if err != nil {
			return err
		}

		node := freeNodes[len(freeNodes)-1]
		newNodes, toAst := model.ToNodes(node, &counter)

		_, ok := alreadyAddNodes[node.ID]
		if !ok {
			nodes = append(nodes, node)
			alreadyAddNodes[node.ID] = nil
		}

		nodes = append(nodes, toAst...)
		freeNodes = append(freeNodes[:len(freeNodes)-1], newNodes...)
	}

	return visualize(nodes, path, name)
}
