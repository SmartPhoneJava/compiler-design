package internal

import "strings"

type varMap map[string]*Var

type Var struct {
	Name, Value string
	CreatedIn   *Func
	UsedIn      map[string]*Func
}

func handleLocalVar(content string, headFunc *Func) {
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
			Value:     varValues[i],
			CreatedIn: headFunc,
			UsedIn: map[string]*Func{
				headFunc.Name: headFunc,
			},
		}

	}
}

func handleGlobalVar(varName, varValue string, currentFunc, mainFunc *Func) {
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
			varNames[i], varValues[i],
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
