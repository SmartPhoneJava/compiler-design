package internal

import (
	"fmt"
	"kurs/internal/visualizer"
	"regexp"
	"strings"
)

type funcMap map[string]*Func

type Funcs struct {
	AllFuncs  map[string]*Func
	callStack []*Func
}

type Func struct {
	Name string

	Args    string
	ArgsRaw string

	Return    string
	ReturnRaw string

	LocalVars   varMap
	LocalTables tableMap
	LocalFuncs  funcMap
	Calls       []*Func
}

func NewFunc(name string) *Func {
	name = FuncName(name)
	return &Func{
		Name:        name,
		LocalFuncs:  map[string]*Func{},
		LocalVars:   map[string]*Var{},
		LocalTables: map[string]*Table{},
	}
}

const MainFunc = "main"

func (f *Funcs) GetFunc(name string) *Func {
	name = FuncName(name)
	funcObj, ok := f.AllFuncs[name]
	if ok {
		return funcObj
	}
	newFunc := NewFunc(name)
	f.AllFuncs[name] = newFunc
	return newFunc
}

func (f *Funcs) GetCallStackTop() *Func {
	if len(f.callStack) == 0 {
		mainFunc := f.GetFunc(MainFunc)
		f.pushToStack(mainFunc)
		return mainFunc
	}
	return f.callStack[len(f.callStack)-1]
}

func (f *Funcs) pushToStack(newFunc *Func) {
	f.callStack = append(f.callStack, newFunc)
}

func (f *Funcs) popFromStack() {
	if len(f.callStack) == 0 {
		return
	}
	f.callStack = f.callStack[:len(f.callStack)-1]
}

func chooseColor(n int) string {
	switch n {
	case 0:
		return `"#ffa000"`
	case 1:
		return `"#03a9f4"`
	case 2:
		return `"#0288d1"`
	default:
		return `"#ffc107"`
	}
}

func (f *Funcs) MustVisualize(path, name string) {
	var (
		mainFunc = f.GetFunc(MainFunc)
		nodes    = []*visualizer.Node{{
			Name: mainFunc.Path(),
			Style: func() map[string]string {
				return map[string]string{
					"color": chooseColor(len(mainFunc.Calls)),
					"shape": `"box"`,
					"style": `"rounded,filled"`,
				}
			},
		}}
		edges []*visualizer.Edge
	)

	f.visualizeInternal(mainFunc, &nodes, &edges, map[string]interface{}{})
	visualizer.Visualize(nodes, edges, path, name)
}

func (f *Funcs) visualizeInternal(
	from *Func,
	nodes *[]*visualizer.Node,
	edges *[]*visualizer.Edge,
	visited map[string]interface{},
) {
	_, ok := visited[from.Name]
	if ok {
		return
	}
	visited[from.Name] = nil
	for _, to := range from.Calls {
		color := chooseColor(len(to.Calls))
		*nodes = append(*nodes, &visualizer.Node{
			Name: to.Path(),
			Style: func() map[string]string {
				return map[string]string{
					"color": color,
					"label": to.NormalizedName(),
					"shape": `"box"`,
					"style": `"rounded,filled"`,
				}
			},
		})
		*edges = append(*edges, &visualizer.Edge{
			From: from.Path(),
			To:   to.Path(),
			Style: func() map[string]string {
				return map[string]string{
					"color": color,
				}
			},
		})
		f.visualizeInternal(to, nodes, edges, visited)
	}
}

func FuncName(name string) string {
	bracketIndex := strings.Index(name, "(")

	var funcName = name
	if bracketIndex >= 0 {
		funcName = name[:bracketIndex]
	}
	var start = len(funcName) - 1
	rex := regexp.MustCompile("^[ a-z0-9A-Z_]+$")
	for start >= 0 {
		if !rex.Match([]byte(funcName[start:])) {
			return funcName[start+1:]
		}
		start--
	}
	return funcName
}

func (funcs funcMap) Node(isGlobal bool) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " functions"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" color="black" CELLSPACING="0">
	<TR><TD COLSPAN="3" BGCOLOR="%s">%s</TD></TR>
	  `, colorMapping[name+":header"], name)

	table += "<TR>" + cell(
		colorMapping[name+":naming"],
		NewPair("Name"),
		NewPair("Input"),
		NewPair("Output"),
	) + "</TR>"

	for _, funcObj := range funcs {
		normalized := funcObj.NormalizedName()
		if isGlobal && normalized == MainFunc {
			continue
		}
		table += "<TR>" + cell(
			colorMapping[name+":body"],
			NewPair(normalized),
			stringPair{"", funcObj.ArgsRaw},
			stringPair{"", funcObj.ReturnRaw},
		) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func (funcs funcMap) Len() int {
	return len(funcs)
}

func (f Func) NormalizedName() string {
	names := strings.Split(f.Name, " ")
	return names[len(names)-1]
}

func (f Func) ParentName() string {
	last := strings.LastIndexByte(f.Name, ' ')
	if last == -1 {
		return f.Name
	}
	return f.Name[:last]
}

func (f Func) Tables() tableMap {
	return f.LocalTables
}

func (f Func) Funcs() funcMap {
	return f.LocalFuncs
}

func (f Func) Vars() varMap {
	return f.LocalVars
}

func (f Func) Path() string {
	return f.Name
}

func (s *InfoCollector) createFunc(content, funcName, bodyContent string) *Func {

	var (
		startI = strings.Index(content, "function")
		endI   = strings.Index(content, bodyContent)
		head   = s.Funcs.GetCallStackTop()
	)

	if content != bodyContent && funcName == "" {
		funcName = content[startI+len("function") : endI]
	}

	if funcName == "" {
		funcName = "anonymous"
	}

	namedFunc := s.Funcs.GetFunc(head.Name + " " + funcName)
	if content[endI] == '(' {
		leftBracket := endI
		params := content[leftBracket+1:]
		rightBracket := strings.Index(params, ")")
		if rightBracket != -1 {
			params = params[:rightBracket]
			namedFunc.Args = params
		}
	}

	// head.Calls = append(head.Calls, namedFunc)
	//s.Funcs.pushToStack(namedFunc)
	return namedFunc
}

func (s *InfoCollector) initFuncInfo(bodyContent string, funcObj *Func) {
	if bodyContent == "" {
		return
	}

	if bodyContent[0] == '(' {
		rightBracket := strings.Index(bodyContent, ")")
		if rightBracket != -1 {
			funcObj.Args = bodyContent[:rightBracket]
		}
	}
}
