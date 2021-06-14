package internal

import (
	"fmt"
	"strconv"
)

type Table struct {
	Name string

	CreatedIn *Func
	UsedIn    funcMap

	LocalVars   varMap
	LocalTables tableMap
	LocalFuncs  funcMap
}

type tableMap map[string]*Table

type Tables struct {
	Tables    tableMap
	callStack []*Table

	CandidateTable string
	CandidateVar   string
	AddByDots      bool
}

func (f *Tables) GetCallStackTop() *Table {
	if len(f.callStack) == 0 {
		mainTable := f.GetTable(MainFunc)
		f.pushToStack(mainTable)
		return mainTable
	}
	return f.callStack[len(f.callStack)-1]
}

func (f *Tables) pushToStack(newTable *Table) {
	f.callStack = append(f.callStack, newTable)
}

func (f *Tables) popFromStack() {
	if len(f.callStack) == 0 {
		return
	}
	f.callStack = f.callStack[:len(f.callStack)-1]
}

func (f *Tables) GetTable(name string) *Table {
	name = FuncName(name)
	funcObj, ok := f.Tables[name]
	if ok {
		return funcObj
	}
	newTable := NewTable(name)
	f.Tables[name] = newTable
	return newTable
}

func NewTable(name string) *Table {
	name = FuncName(name)
	return &Table{
		Name:        name,
		LocalFuncs:  map[string]*Func{},
		LocalVars:   map[string]*Var{},
		LocalTables: map[string]*Table{},
	}
}

func (tables tableMap) Node(isGlobal bool) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " tables"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" color="black">
  <TR><TD PORT="0" COLSPAN="3" BGCOLOR="%s">table: %s</TD></TR>
	`, colorMapping[name+":header"], name)

	table += "<TR>" + cell(colorMapping[name+":naming"], "Name", "num of fields", "num of methods") + "</TR>"

	for _, tableObj := range tables {
		createdIn := tableObj.CreatedIn
		var createdInName string

		if createdIn != nil {
			createdInName = createdIn.Name
		}
		table += "<TR>" + cell(
			colorMapping[name+":body"],
			tableObj.Name,
			strconv.Itoa(len(tableObj.LocalVars)+len(tableObj.LocalTables)),
			createdInName,
		) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func (tables tableMap) Len() int {
	return len(tables)
}
