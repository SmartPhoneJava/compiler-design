package internal

import (
	"fmt"
	"kurs/internal/visualizer"
	"strings"
)

func VisualInternal() {

}

func FuncTableMustVisualize(f *Funcs, path, name string) {
	var (
		nodes []*visualizer.Node
		edges []*visualizer.Edge
	)

	visualizeInternal("", f.GetFunc(MainFunc), &nodes, &edges, map[string]interface{}{})
	visualizer.Visualize(nodes, edges, path, name)
}

type Object interface {
	Vars() varMap
	Tables() tableMap
	Funcs() funcMap

	Path() string
	NormalizedName() string
}

var repeat = 1000

func visualizeInternal(
	parent string,
	object Object,
	nodes *[]*visualizer.Node,
	edges *[]*visualizer.Edge,
	visited map[string]interface{},
) {
	_, ok := visited[object.Path()]
	if ok {
		return
	}
	visited[object.Path()] = nil

	repeat--
	*nodes = append(*nodes, &visualizer.Node{
		Name: object.Path(),
		Style: func() map[string]string {
			return map[string]string{
				"label": object.NormalizedName(),
				"shape": `"plaintext"`,
				"style": `"rounded,filled"`,
				"color": "white",
			}
		},
	})

	if parent != "" {
		*edges = append(*edges, &visualizer.Edge{
			From:     parent,
			To:       object.Path(),
			FromPort: object.NormalizedName(),
			Style: func() map[string]string {
				return map[string]string{
					"color": "lightgrey",
				}
			},
		})
	}

	var localGlobal = "local"
	if object.NormalizedName() == MainFunc {
		localGlobal = "global"
	}

	// add vars
	addNodeEdge(
		object.Vars(), object,
		object.Path(), localGlobal+" variables",
		nodes, edges,
	)

	// add funcs
	addNodeEdge(
		object.Funcs(), object,
		object.Path(), localGlobal+" functions",
		nodes, edges,
	)

	// add table
	addNodeEdge(
		object.Tables(), object,
		object.Path(), localGlobal+" tables",
		nodes, edges,
	)

	for _, funcObj := range object.Funcs() {
		visualizeInternal(
			object.Path()+" "+localGlobal+" functions",
			funcObj,
			nodes,
			edges,
			visited,
		)
	}

	for _, tableObj := range object.Tables() {
		visualizeInternal(
			object.Path()+" "+localGlobal+" tables",
			tableObj,
			nodes,
			edges,
			visited,
		)
	}
}

type ToNode interface {
	Node(isGlobal bool) string
	Len() int
}

func addNodeEdge(
	toNode ToNode,
	obj Object,
	parent, extra string,
	nodes *[]*visualizer.Node,
	edges *[]*visualizer.Edge,
) {
	if toNode.Len() == 0 {
		return
	}
	var (
		funcNodeName = obj.Path()
		funcName     = obj.NormalizedName()
		child        = funcNodeName + " " + extra
		nodeContent  = toNode.Node(funcName == MainFunc)
	)
	*nodes = append(*nodes, &visualizer.Node{
		Name: child,
		Style: func() map[string]string {
			return map[string]string{
				"label": nodeContent,
				"shape": `"plaintext"`,
				"style": `"rounded,filled"`,
				"color": "white",
			}
		},
	})

	color := "black"
	mapColor, ok := colorMapping["Arrow "+extra]
	if ok {
		color = mapColor
	}
	*edges = append(*edges, &visualizer.Edge{
		From: parent,
		To:   child,
		Style: func() map[string]string {
			return map[string]string{
				"color": `"` + color + `"`,
			}
		},
	})
}

type stringPair struct {
	One, Two string
}

func NewPair(str string) stringPair {
	return stringPair{One: str, Two: str}
}

func cell(color string, values ...stringPair) string {
	var newString string

	mapColor, ok := colorMapping[color]
	if ok {
		color = mapColor
	}

	for _, pair := range values {
		pair.One = strings.Trim(pair.One, `"`)
		newString += fmt.Sprintf(`<TD BGCOLOR="%s" PORT="%s">%s</TD>\n`, color, pair.One, pair.Two)
	}
	return newString
}

// https://colorscheme.ru/html-colors.html
var colorMapping = map[string]string{
	"Global functions:header": "#FFD700",
	"Global functions:naming": "#FFFF00",
	"Global functions:body":   "#FFFFE0",
	"Arrow global functions":  "#FF8C00",

	"Global variables:header": "#FFDAB9",
	"Global variables:naming": "#FFE4B5",
	"Global variables:body":   "#FFEFD5",
	"Arrow global variables":  "brown",

	"Global tables:header": "#BDB76B",
	"Global tables:naming": "#F0E68C",
	"Global tables:body":   "#EEE8AA",
	"Arrow global tables":  "#BDB76B",

	"Local functions:header": "#00FA9A",
	"Local functions:naming": "#BDECB6",
	"Local functions:body":   "#D0F0C0",
	"Arrow local functions":  "#00FA9A",

	"Local variables:header": "#7FFFD4",
	"Local variables:naming": "#AFEEEE",
	"Local variables:body":   "#E0FFFF",
	"Arrow local variables":  "#7FFFD4",

	"Local tables:header": "#00BFFF",
	"Local tables:naming": "#87CEFA",
	"Local tables:body":   "#B0E0E6",
	"Arrow local tables":  "#00BFFF",
}
