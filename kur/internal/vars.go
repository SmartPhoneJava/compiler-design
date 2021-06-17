package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type varMap map[string]*Var

type Var struct {
	Name, Value, RealText string
	CreatedIn             *Func
	UsedIn                map[string]*Func
}

func (vars varMap) Node(isGlobal bool) string {
	var name = "Local"
	if isGlobal {
		name = "Global"
	}
	name += " variables"

	var table = fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" color="black">
  <TR><TD PORT="0" COLSPAN="2" BGCOLOR="%s">table: %s</TD></TR>
	`, colorMapping[name+":header"], name)

	table += "<TR>" + cell(colorMapping[name+":naming"],
		NewPair("Name"),
		NewPair("Value"),
	) + "</TR>"

	for _, varObj := range vars {
		value := varObj.RealText
		if value == "" {
			value = varObj.Value
		}
		table += "<TR>" + cell(colorMapping[name+":body"],
			stringPair{"", varObj.Name},
			stringPair{"", decodeString(value)},
		) + "</TR>"
	}

	table += `\n</TABLE>>`
	return table
}

func (vars varMap) Len() int {
	return len(vars)
}

func (s *InfoCollector) getVarName(varName string) string {
	varName = strings.TrimPrefix(varName, ".")
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

func (s *InfoCollector) pickVar(node AntlrNode) error {
	value := node.GetText()
	var (
		vars    = strings.Split(s.candidateVar, ",")
		varName = vars[0]
		head    = s.Funcs.GetCallStackTop()
		varObj  = &Var{
			Name:      varName,
			Value:     value,
			CreatedIn: head,
			RealText:  s.GetText(node),
		}
	)
	if varName == "" {
		return nil
	}

	s.expression = strings.TrimPrefix(s.expression, ",")
	pickI := strings.Index(s.expression, value)

	if pickI != 0 {
		return nil
	}

	var err error
	var withPoints = strings.Split(varName, ".")

	if len(withPoints) > 1 {
		_, err = s.pickField(withPoints, value, false)
	}
	if err != nil {
		return err
	}
	s.expression = s.expression[pickI+len(value):]

	if len(vars) > 0 {
		s.candidateVar = strings.Join(vars[1:], ",")
	}
	if len(withPoints) > 1 {
		return nil
	}

	funcObj := &Func{}
	if s.localStatus == LocalVar {
		funcObj = head
	} else {
		funcObj = s.Funcs.GetFunc(MainFunc)
	}

	funcObj.LocalVars[varName] = varObj
	return nil
}

func (s *InfoCollector) pickTable(value string) error {
	var (
		vars      = strings.Split(s.candidateVar, ",")
		tableName = vars[0]
	)

	var withPoints = strings.Split(tableName, ".")
	var err error
	var name string
	if len(withPoints) > 1 {
		name, err = s.pickField(withPoints, value, true)
	}

	if err != nil {
		return err
	}
	if len(vars) > 0 {
		s.candidateVar = strings.Join(vars[1:], ",")
	}

	if len(withPoints) > 1 {
		table, _ := s.Tables.GetTable(name)
		s.Tables.pushToStack(table)
		return nil
	}

	if tableName == "" {
		tableName = "0" // !! ???????
	}

	s.expression = strings.TrimPrefix(s.expression, ",")
	pickI := strings.Index(s.expression, value)
	if pickI != 0 {
		return errors.New("cant pick")
	}

	s.expression = s.expression[pickI+len(value):]

	table := &Table{}
	if s.Tables.currentLvl == 0 {
		funcObj := &Func{}
		if s.localStatus == LocalVar {
			funcObj = s.Funcs.GetCallStackTop()
		} else {
			funcObj = s.Funcs.GetFunc(MainFunc)
		}
		table = NewTable(funcObj.Name + " " + tableName)
		funcObj.LocalTables[table.NormalizedName()] = table
	} else {
		parent := s.Tables.GetCallStackTop()
		table = NewTable(parent.Name + " " + tableName)
		parent.LocalTables[table.NormalizedName()] = table
	}
	s.Tables.pushToStack(table)
	return nil
}

func (s *InfoCollector) pickField(
	withPoints []string,
	value string,
	isTable bool,
) (string, error) {

	if len(withPoints) > 1 {
		table, ok := s.Funcs.GetFunc(MainFunc).LocalTables[withPoints[0]]
		if !ok {
			return "", errors.New("no such table")
		}
		for i, tableName := range withPoints {
			if i == 0 {
				continue
			}
			if i == len(withPoints)-1 {
				var name string
				if isTable {
					name = table.Name + " " + tableName
					newTable, _ := s.Tables.GetTable(name)
					s.Tables.pushToStack(newTable)
					table.LocalTables[newTable.NormalizedName()] = newTable
				} else {
					name = tableName
					table.LocalVars[tableName] = &Var{
						Name:  tableName,
						Value: value,
					}
				}
				return name, nil
			}
			table, ok = table.LocalTables[tableName]
			if !ok {
				return "", errors.New("no such table")
			}
		}
	}
	return "", nil
}

func decodeString(str string) string {
	for oldStr, newStr := range mapSymbols {
		str = strings.ReplaceAll(str, oldStr, newStr)
	}
	return str
}

var mapSymbols = map[string]string{
	"<": "_less_",
	">": "_more_",
}
