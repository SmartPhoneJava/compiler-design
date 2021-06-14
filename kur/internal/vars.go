package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type varMap map[string]*Var

type Var struct {
	Name, Value string
	CreatedIn   *Func
	UsedIn      map[string]*Func
}

func (c *InfoCollector) handleLocalVar(content string, headFunc *Func) {
	vars := strings.Split(content, "=")
	if len(vars) != 2 {
		return
	}

	vars[0] = strings.TrimLeft(vars[0], "local")
	var (
		varNames  = strings.Split(vars[0], ",")
		varValues = strings.Split(vars[1], ",")
	)
	if len(varNames) != len(varValues) {
		return
	}

	for i := 0; i < len(varNames); i++ {

		headFunc.LocalVars[varNames[i]] = &Var{
			Name:      varNames[i],
			Value:     c.getVarName(varValues[i]),
			CreatedIn: headFunc,
			UsedIn: map[string]*Func{
				headFunc.Name: headFunc,
			},
		}

	}
}

func (c *InfoCollector) handleGlobalVar(varName, varValue string, currentFunc, mainFunc *Func) {
	if varValue == "" || varName == "" {
		return
	}

	var (
		varNames  = strings.Split(varName, ",")
		varValues = strings.Split(varValue, ",")
	)
	if len(varNames) != len(varValues) {
		return
	}

	for i := 0; i < len(varNames); i++ {
		updateVar(
			mainFunc.LocalVars,
			varNames[i],
			c.getVarName(varValues[i]),
			currentFunc,
		)
	}
}

func updateVar(
	varsMap map[string]*Var,
	varName, varValue string,
	currFunc *Func,
) {
	_, alreadyCreated := varsMap[varName]
	if !alreadyCreated {
		varsMap[varName] = &Var{
			Name:      varName,
			Value:     varValue,
			CreatedIn: currFunc,
			UsedIn:    map[string]*Func{},
		}
	}

	if currFunc != nil {
		varsMap[varName].UsedIn[currFunc.Name] = currFunc
	}
}

func (vars varMap) Node(isGlobal bool) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " variables"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" color="black">
  <TR><TD PORT="0" COLSPAN="3" BGCOLOR="%s">table: %s</TD></TR>
	`, colorMapping[name+":header"], name)

	table += "<TR>" + cell(colorMapping[name+":naming"],
		NewPair("Name"),
		NewPair("Value"),
		NewPair("first initialized by"),
	) + "</TR>"

	for _, varObj := range vars {
		createdIn := varObj.CreatedIn
		var createdInName string

		if createdIn != nil {
			createdInName = createdIn.Name
		}

		table += "<TR>" + cell(colorMapping[name+":body"],
			NewPair(varObj.Name),
			NewPair(varObj.Value),
			NewPair(createdInName),
		) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func (vars varMap) Len() int {
	return len(vars)
}

func (s *InfoCollector) getVarName(varName string) string {
	varName = strings.TrimLeft(varName, ".")
	varName = strings.TrimPrefix(varName, "[")
	varName = strings.TrimSuffix(varName, "]")
	if !(strings.HasPrefix(varName, `"`) &&
		strings.HasSuffix(varName, `"`)) {
		_, err := strconv.Atoi(varName)
		if err != nil {
			varName = `"` + varName + `"`
		}
	}
	return varName
}
