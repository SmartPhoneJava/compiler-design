package internal

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

	// неявная индексация, когда указываются значения в конструкторе
	// без индексов
	implicitIndex map[int]int
	currentLvl    int
	reserveName   string
}

func (tables *Tables) GetCallStackTop() *Table {
	if len(tables.callStack) == 0 {
		mainTable, _ := tables.GetTable(MainFunc)
		tables.pushToStack(mainTable)
		return mainTable
	}
	return tables.callStack[len(tables.callStack)-1]
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

func (f *Tables) GetTable(name string) (*Table, bool) {
	name = FuncName(name)
	funcObj, ok := f.Tables[name]
	if ok {
		return funcObj, true
	}
	newTable := NewTable(name)
	f.Tables[name] = newTable
	return newTable, false
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

	table += "<TR>" + cell(colorMapping[name+":naming"],
		NewPair("Name"),
		NewPair("num of fields"),
		NewPair("num of methods"),
	) + "</TR>"

	for _, tableObj := range tables {
		table += "<TR>" + cell(
			colorMapping[name+":body"],
			NewPair(tableObj.NormalizedName()),
			NewPair(strconv.Itoa(len(tableObj.LocalVars)+len(tableObj.LocalTables))),
			NewPair(strconv.Itoa(len(tableObj.LocalFuncs))),
		) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func (tables tableMap) Len() int {
	return len(tables)
}

func (table Table) Tables() tableMap {
	return table.LocalTables
}

func (table Table) Funcs() funcMap {
	return table.LocalFuncs
}

func (table Table) Vars() varMap {
	return table.LocalVars
}

func (table Table) NormalizedName() string {
	names := strings.Split(table.Name, " ")
	return names[len(names)-1]
}

func (table Table) Path() string {
	return table.Name
}

func (s *InfoCollector) createTable(content string) {
	var headTable = s.Tables.GetCallStackTop()
	var name = headTable.Name + " "
	var namedTable *Table

	if s.Tables.currentLvl > 0 {
		if s.Tables.reserveName != "" {
			name += " " + s.Tables.reserveName
		} else {
			index := s.Tables.implicitIndex[s.Tables.currentLvl]
			name += " anonymous " + strconv.Itoa(index+1)
		}
		namedTable, _ = s.Tables.GetTable(name)
		headTable.LocalTables[namedTable.NormalizedName()] = namedTable
		s.Tables.pushToStack(namedTable)
	} else {
		err := s.pickTable(content)
		if err != nil {
			log.Println("errr", err)
		}
	}
	//s.Tables.implicitIndex[s.Tables.currentLvl]++
	s.Tables.currentLvl++

}

// 168
