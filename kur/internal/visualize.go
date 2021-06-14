package internal

import (
	"fmt"
	"kurs/internal/visualizer"
	"log"
)

func VisualInternal() {

}

func FuncTableMustVisualize(f *Funcs, path, name string) {
	nodes := []*visualizer.Node{
		{
			Name: "Statistics",
			Style: func() map[string]string {
				return map[string]string{
					"shape": `"box"`,
					"style": `"rounded,filled"`,
				}
			},
		},
	}
	color := chooseColor(0)

	edges := []*visualizer.Edge{}

	for _, funcObj := range f.AllFuncs {
		if len(funcObj.LocalVars) == 0 && len(funcObj.LocalFuncs) == 0 {
			continue
		}

		var (
			funcName     = funcObj.NormalizedName()
			funcNodeName = funcObj.Name
		)

		fromNode := "Statistics"

		if funcName != MainFunc {
			fromNode = funcNodeName
			// add funcObj
			nodes = append(nodes, &visualizer.Node{
				Name: funcNodeName,
				Style: func() map[string]string {
					return map[string]string{
						"label": funcName,
						"shape": `"plaintext"`,
						"style": `"rounded,filled"`,
						"color": "white",
					}
				},
			})

			fromForEdge := funcObj.ParentName() //"Statistics"
			if fromForEdge == MainFunc {
				fromForEdge = MainFunc + " local functions"
			}
			edges = append(edges, &visualizer.Edge{
				From: fromForEdge,
				To:   funcNodeName,
				Style: func() map[string]string {
					return map[string]string{
						"color": color,
					}
				},
			})
		}

		log.Println("local var", funcObj.LocalVars)
		// add func's vars
		if len(funcObj.LocalVars) > 0 {
			funcsVarsName := funcNodeName + " local variables"
			funcsVars := VarsTable(funcName == MainFunc, funcObj.LocalVars)
			nodes = append(nodes, &visualizer.Node{
				Name: funcsVarsName,
				Style: func() map[string]string {
					return map[string]string{
						"label": funcsVars,
						"shape": `"plaintext"`,
						"style": `"rounded,filled"`,
						"color": "white",
					}
				},
			})
			edges = append(edges, &visualizer.Edge{
				From: fromNode,
				To:   funcsVarsName,
				Style: func() map[string]string {
					return map[string]string{
						"color": color,
					}
				},
			})
		}

		// // add func's funcs
		// if len(funcObj.LocalFuncs) > 0 {
		// 	funcsInnerFuncsName := funcNodeName + " local functions"
		// 	funcsInnerFuncs := FuncsTable(funcName == MainFunc, funcObj.LocalFuncs)
		// 	nodes = append(nodes, &visualizer.Node{
		// 		Name: funcsInnerFuncsName,
		// 		Style: func() map[string]string {
		// 			return map[string]string{
		// 				"label": funcsInnerFuncs,
		// 				"shape": `"plaintext"`,
		// 				"style": `"rounded,filled"`,
		// 				"color": "white",
		// 			}
		// 		},
		// 	})
		// 	edges = append(edges, &visualizer.Edge{
		// 		From: fromNode,
		// 		To:   funcsInnerFuncsName,
		// 		Style: func() map[string]string {
		// 			return map[string]string{
		// 				"color": color,
		// 			}
		// 		},
		// 	})
		// }

		// add funcs
		addNodeEdge(
			funcObj.LocalFuncs, funcObj,
			funcNodeName, "local functions",
			&nodes, &edges,
		)

		// add table
		addNodeEdge(
			funcObj.LocalTables, funcObj,
			funcNodeName, "local tables",
			&nodes, &edges,
		)
	}

	visualizer.Visualize(nodes, edges, path, name)
}

type ToNode interface {
	Node(isGlobal bool) string
	Len() int
}

func addNodeEdge(
	toNode ToNode,
	currFunc *Func,
	fromNode, extra string,
	nodes *[]*visualizer.Node,
	edges *[]*visualizer.Edge,
) {
	if toNode.Len() == 0 {
		return
	}
	var (
		funcNodeName = currFunc.Name
		funcName     = currFunc.NormalizedName()
		nodeName     = funcNodeName + " " + extra
		nodeContent  = toNode.Node(funcName == MainFunc)
	)
	*nodes = append(*nodes, &visualizer.Node{
		Name: nodeName,
		Style: func() map[string]string {
			return map[string]string{
				"label": nodeContent,
				"shape": `"plaintext"`,
				"style": `"rounded,filled"`,
				"color": "white",
			}
		},
	})

	fromForEdge := fromNode //currFunc.ParentName() //"Statistics"
	if fromForEdge == MainFunc {
		fromForEdge = "Statistics"
	}
	log.Println("nodeName", fromForEdge)
	log.Println("nodeName", nodeName)
	*edges = append(*edges, &visualizer.Edge{
		From: fromForEdge,
		To:   nodeName,
		Style: func() map[string]string {
			return map[string]string{
				"color": "grey",
			}
		},
	})
}

func FuncsTable(isGlobal bool, funcs map[string]*Func) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " functions"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" color="black" CELLSPACING="0">
  <TR><TD COLSPAN="3" BGCOLOR="%s">%s</TD></TR>
	`, colorMapping[name+":header"], name)

	table += "<TR>" + cell(colorMapping[name+":naming"], "Name", "Input", "Output") + "</TR>"

	for _, funcObj := range funcs {
		normalized := funcObj.NormalizedName()
		if isGlobal && normalized == MainFunc {
			continue
		}
		table += "<TR>" + cell(colorMapping[name+":body"], normalized, funcObj.Args, funcObj.Return) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func VarsTable(isGlobal bool, vars map[string]*Var) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " variables"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" color="black">
  <TR><TD PORT="0" COLSPAN="3" BGCOLOR="%s">table: %s</TD></TR>
	`, colorMapping[name+":header"], name)

	table += "<TR>" + cell(colorMapping[name+":naming"], "Name", "Value", "first initialized by") + "</TR>"

	for _, varObj := range vars {
		createdIn := varObj.CreatedIn
		var createdInName string

		if createdIn != nil {
			createdInName = createdIn.Name
		}
		table += "<TR>" + cell(colorMapping[name+":body"], varObj.Name, varObj.Value, createdInName) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func cell(color string, values ...string) string {
	var newString string

	mapColor, ok := colorMapping[color]
	if ok {
		color = mapColor
	}

	for _, val := range values {
		newString += fmt.Sprintf(`<TD BGCOLOR="%s" PORT="k1">%s</TD>\n`, color, val)
	}
	return newString
}

// https://colorscheme.ru/html-colors.html
var colorMapping = map[string]string{
	"Global functions:header": "#FFD700",
	"Global functions:naming": "#FFFF00",
	"Global functions:body":   "#FFFFE0",

	"Global variables:header": "#FFDAB9",
	"Global variables:naming": "#FFE4B5",
	"Global variables:body":   "#FFEFD5",

	"Global tables:header": "#BDB76B",
	"Global tables:naming": "#F0E68C",
	"Global tables:body":   "#EEE8AA",

	"Local functions:header": "#00FA9A",
	"Local functions:naming": "#90EE90",
	"Local functions:body":   "#98FB98",

	"Local variables:header": "#7FFFD4",
	"Local variables:naming": "#AFEEEE",
	"Local variables:body":   "#E0FFFF",

	"Local tables:header": "#00BFFF",
	"Local tables:naming": "#87CEFA",
	"Local tables:body":   "#B0E0E6",
}
