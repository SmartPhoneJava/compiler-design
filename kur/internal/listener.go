package internal

import (
	"fmt"
	"kurs/parser"
	"log"
	"strings"
)

// EnterChunk is called when production chunk is entered.
func (s *InfoCollector) EnterChunk(ctx *parser.ChunkContext) {}

// ExitChunk is called when production chunk is exited.
func (s *InfoCollector) ExitChunk(ctx *parser.ChunkContext) {}

// EnterBlock is called when production block is entered.
func (s *InfoCollector) EnterBlock(ctx *parser.BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *InfoCollector) ExitBlock(ctx *parser.BlockContext) {}

// EnterStat is called when production stat is entered.
func (s *InfoCollector) EnterStat(ctx *parser.StatContext) {
	statValue := ctx.GetText()
	s.expression = ""
	log.Println("ENTER STAT")

	switch {
	case strings.HasPrefix(statValue, "localfunc"):
		s.localStatus = LocalFunc
		s.candidate = statValue
	case strings.HasPrefix(statValue, "local"):
		s.localStatus = MaybeLocalVar
		s.candidateVar = statValue

		tableName := statValue[len("local"):]
		equalI := strings.Index(tableName, "=")
		if equalI != -1 {
			tableName = tableName[:equalI]
		}

		s.Tables.CandidateVar = tableName
	default:
		s.Tables.CandidateTable = statValue
		s.localStatus = NotLocal
	}
}

// ExitStat is called when production stat is exited.
func (s *InfoCollector) ExitStat(ctx *parser.StatContext) {}

// EnterAttnamelist is called when production attnamelist is entered.
func (s *InfoCollector) EnterAttnamelist(ctx *parser.AttnamelistContext) {
	if s.localStatus == MaybeLocalVar {
		s.localStatus = LocalVar
		s.candidateVar = ctx.GetText()

	}
}

// ExitAttnamelist is called when production attnamelist is exited.
func (s *InfoCollector) ExitAttnamelist(ctx *parser.AttnamelistContext) {}

// EnterAttrib is called when production attrib is entered.
func (s *InfoCollector) EnterAttrib(ctx *parser.AttribContext) {}

// ExitAttrib is called when production attrib is exited.
func (s *InfoCollector) ExitAttrib(ctx *parser.AttribContext) {}

// EnterRetstat is called when production retstat is entered.
func (s *InfoCollector) EnterRetstat(ctx *parser.RetstatContext) {
	retValue := ctx.GetText()
	retValue = strings.TrimPrefix(retValue, "return")
	namedFunc := s.Funcs.GetCallStackTop()
	log.Println("reeeeeturn", retValue, namedFunc.Name)
	namedFunc.Return = retValue
	namedFunc.ReturnRaw = strings.TrimLeft(s.GetText(ctx), "return")
	s.isReturning = true
	s.candidate = retValue
}

// #6 функции
// ExitRetstat is called when production retstat is exited.
func (s *InfoCollector) ExitRetstat(ctx *parser.RetstatContext) {
	s.isReturning = false
}

// EnterLabel is called when production label is entered.
func (s *InfoCollector) EnterLabel(ctx *parser.LabelContext) {}

// ExitLabel is called when production label is exited.
func (s *InfoCollector) ExitLabel(ctx *parser.LabelContext) {}

// #1 граф вызовов
// EnterFuncname is called when production funcname is entered.
func (s *InfoCollector) EnterFuncname(ctx *parser.FuncnameContext) {
	name := ctx.GetText()
	var namedFunc *Func
	blocks, ok := s.fields(name)
	s.tableFunc = name

	if ok {
		name, err := s.pickField(blocks, "", TypeFunc)
		if err != nil {
			log.Println("EEEEEEEEEEEEEEEERRR")
			return
		}
		namedFunc = s.Funcs.GetFunc(name)
		s.tableFunc = name
		log.Println("ADDDDDDDDDD", name)
	} else {
		namedFunc = s.Funcs.GetFunc(s.Funcs.GetCallStackTop().Name + " " + name)
		if s.localStatus != LocalFunc {
			s.candidate = name
		}
	}

	s.Funcs.pushToStack(namedFunc)
	log.Println("EnterFuncname", namedFunc.Name)
	return

}

// ExitFuncname is called when production funcname is exited.
func (s *InfoCollector) ExitFuncname(ctx *parser.FuncnameContext) {}

// #5 поля
// EnterVarlist is called when production varlist is entered.
func (s *InfoCollector) EnterVarlist(ctx *parser.VarlistContext) {
	s.candidateVar = ctx.GetText()
}

// ExitVarlist is called when production varlist is exited.
func (s *InfoCollector) ExitVarlist(ctx *parser.VarlistContext) {}

// EnterNamelist is called when production namelist is entered.
func (s *InfoCollector) EnterNamelist(ctx *parser.NamelistContext) {
}

// ExitNamelist is called when production namelist is exited.
func (s *InfoCollector) ExitNamelist(ctx *parser.NamelistContext) {}

// #1 граф вызовов
// EnterExplist is called when production explist is entered.
func (s *InfoCollector) EnterExplist(ctx *parser.ExplistContext) {
	if s.expression == "" {
		s.expression = ctx.GetText()
	}
}

// ExitExplist is called when production explist is exited.
func (s *InfoCollector) ExitExplist(ctx *parser.ExplistContext) {}

// EnterExp is called when production exp is entered.
func (s *InfoCollector) EnterExp(ctx *parser.ExpContext) {
	content := ctx.GetText() // s.GetText(ctx.BaseParserRuleContext)

	headTable := s.Tables.GetCallStackTop()
	if headTable.NormalizedName() == "main" {
		if !strings.HasPrefix(content, "{") {
			log.Println("aaaaaaaaaaaaaaaaaa")
			s.pickVar(ctx)
		}
	}
}

// ExitExp is called when production exp is exited.
func (s *InfoCollector) ExitExp(ctx *parser.ExpContext) {}

// EnterPrefixexp is called when production prefixexp is entered.
func (s *InfoCollector) EnterPrefixexp(ctx *parser.PrefixexpContext) {
	content := ctx.GetText()
	s.pickVar(ctx)

	tI := strings.Index(s.expression, content)
	if tI != -1 {
		s.expression = s.expression[tI+len(content):]
	}
}

// ExitPrefixexp is called when production prefixexp is exited.
func (s *InfoCollector) ExitPrefixexp(ctx *parser.PrefixexpContext) {}

// #1 граф вызовов
// EnterFunctioncall is callefd when production functioncall is entered.
func (s *InfoCollector) EnterFunctioncall(ctx *parser.FunctioncallContext) {
	s.candidate = ctx.GetText()
}

// ExitFunctioncall is called when production functioncall is exited.
func (s *InfoCollector) ExitFunctioncall(ctx *parser.FunctioncallContext) {}

// EnterVarOrExp is called when production varOrExp is entered.
func (s *InfoCollector) EnterVarOrExp(ctx *parser.VarOrExpContext) {}

// ExitVarOrExp is called when production varOrExp is exited.
func (s *InfoCollector) ExitVarOrExp(ctx *parser.VarOrExpContext) {}

// EnterVar_ is called when production var_ is entered.
func (s *InfoCollector) EnterVar_(ctx *parser.Var_Context) {}

// ExitVar_ is called when production var_ is exited.
func (s *InfoCollector) ExitVar_(ctx *parser.Var_Context) {}

// EnterVarSuffix is called when production varSuffix is entered.
func (s *InfoCollector) EnterVarSuffix(ctx *parser.VarSuffixContext) {}

// ExitVarSuffix is called when production varSuffix is exited.
func (s *InfoCollector) ExitVarSuffix(ctx *parser.VarSuffixContext) {}

// EnterNameAndArgs is called when production nameAndArgs is entered.
func (s *InfoCollector) EnterNameAndArgs(ctx *parser.NameAndArgsContext) {}

// ExitNameAndArgs is called when production nameAndArgs is exited.
func (s *InfoCollector) ExitNameAndArgs(ctx *parser.NameAndArgsContext) {}

// EnterArgs is called when production args is entered.
func (s *InfoCollector) EnterArgs(ctx *parser.ArgsContext) {
	var (
		name     = ctx.GetText()
		argsI    = strings.Index(s.candidate, name)
		funcName = s.candidate
		funcArgs = s.candidate
	)
	if argsI >= 0 {
		funcName = s.candidate[:argsI]
		funcArgs = s.candidate[argsI:]

		args := strings.Split(funcArgs, ",")
		funcArgs = ""
		for i := range args {
			funcArgs += fmt.Sprintf("arg%d", i+1)
			if i != len(args)-1 {
				funcArgs += ", "
			}
		}
	}

	namedFunc := s.Funcs.GetFunc(s.Funcs.GetCallStackTop().Name + " " + funcName)
	namedFunc.Args = funcArgs

	headFunc := s.Funcs.GetCallStackTop()
	headFunc.Calls = append(headFunc.Calls, namedFunc)
}

// ExitArgs is called when production args is exited.
func (s *InfoCollector) ExitArgs(ctx *parser.ArgsContext) {}

// EnterFunctiondef is called when production functiondef is entered.
func (s *InfoCollector) EnterFunctiondef(ctx *parser.FunctiondefContext) {}

// ExitFunctiondef is called when production functiondef is exited.
func (s *InfoCollector) ExitFunctiondef(ctx *parser.FunctiondefContext) {}

// EnterFuncbody is called when production funcbody is entered.
func (s *InfoCollector) EnterFuncbody(ctx *parser.FuncbodyContext) {
	if s.tableFunc != "" {
		tables := strings.Split(s.tableFunc, ".")
		if len(tables) < 2 {
			return
		}
		mainFunc := s.Funcs.GetFunc(MainFunc)
		currentTable := &Table{}
		for i, t := range tables {
			switch {
			case i == 0:
				newTable, ok := mainFunc.LocalTables[t]
				if !ok {
					return
				}
				currentTable = newTable
			case i == len(tables)-1:
				var (
					funcName    = currentTable.Name + " " + t
					content     = ctx.GetText()
					createdFunc = s.createFunc(content, funcName, content)
				)
				currentTable.LocalFuncs[createdFunc.NormalizedName()] = createdFunc
			default:
				newTable, ok := currentTable.LocalTables[t]
				if !ok {
					return
				}
				currentTable = newTable
			}
		}
	}

	var newFunc, addTo *Func

	// Добавление локальной функции в мапу всех функций
	if s.localStatus == LocalFunc || s.isReturning {
		addTo = s.Funcs.GetCallStackTop()

		newFunc = s.createFunc(s.candidate, "", ctx.GetText())
		s.Funcs.GetCallStackTop().LocalFuncs[newFunc.NormalizedName()] = newFunc
		s.Funcs.pushToStack(newFunc)
	} else {
		newFunc = s.Funcs.GetCallStackTop()
		s.initFuncInfo(ctx.GetText(), newFunc)
		addTo = s.Funcs.GetFunc(MainFunc)
	}

	addTo.LocalFuncs[newFunc.NormalizedName()] = newFunc
}

// #1 функции
// ExitFuncbody is called when production funcbody is exited.
func (s *InfoCollector) ExitFuncbody(ctx *parser.FuncbodyContext) {
	s.Funcs.popFromStack()
}

// EnterParlist is called when production parlist is entered.
func (s *InfoCollector) EnterParlist(ctx *parser.ParlistContext) {
	head := s.Funcs.GetCallStackTop()
	if head.Name == MainFunc {
		return
	}
	head.ArgsRaw = s.GetText(ctx)
}

// ExitParlist is called when production parlist is exited.
func (s *InfoCollector) ExitParlist(ctx *parser.ParlistContext) {}

// #4 таблицы
// EnterTableconstructor is called when production tableconstructor is entered.
func (s *InfoCollector) EnterTableconstructor(ctx *parser.TableconstructorContext) {
	s.createTable(ctx.GetText())
}

// #4 таблицы
// ExitTableconstructor is called when production tableconstructor is exited.
func (s *InfoCollector) ExitTableconstructor(ctx *parser.TableconstructorContext) {
	s.Tables.currentLvl--
	if s.Tables.currentLvl == 0 {
		s.Tables.implicitIndex = make(map[int]int)
	}
	s.Tables.popFromStack()
}

// EnterFieldlist is called when production fieldlist is entered.
func (s *InfoCollector) EnterFieldlist(ctx *parser.FieldlistContext) {
	content := ctx.GetText()
	log.Println("EnterFieldlist", content)
	if strings.HasPrefix(content, "{") {
		return
	} else {
		//varName := strings.Split(content, "=")[0]
		//varName = s.getVarName(varName)
		//s.Tables.reserveName = varName
		s.expression = content
	}
}

// ExitFieldlist is called when production fieldlist is exited.
func (s *InfoCollector) ExitFieldlist(ctx *parser.FieldlistContext) {}

// #4 таблицы
// EnterField is called when production field is entered.
func (s *InfoCollector) EnterField(ctx *parser.FieldContext) {

	var (
		field             = ctx.GetText()
		fieldParts        = strings.Split(field, "=")
		varName, varValue string
		namedTable        = s.Tables.GetCallStackTop()
	)
	log.Println("EnterFieldEnterField", namedTable.Path())
	if strings.HasPrefix(field, "{") {
		return
	}

	if len(fieldParts) > 1 {
		varName = fieldParts[0]
		varValue = fieldParts[1]
	} else {
		varName = fmt.Sprintf("%d", s.Tables.implicitIndex[s.Tables.currentLvl])
		s.Tables.implicitIndex[s.Tables.currentLvl]++

		varValue = field
	}

	log.Println("EnterField", varName, "!", varValue)

	varName = s.getVarName(varName)
	// if s.Tables.currentLvl < 1 {
	// 	s.Funcs.GetCallStackTop().LocalTables[namedTable.NormalizedName()] = namedTable
	// }
	varValue = strings.Trim(varValue, `"`)
	if !strings.HasPrefix(varValue, "{") {

		realTexts := strings.Split(s.GetText(ctx), "=")
		var realText = s.GetText(ctx)
		if len(realTexts) > 1 {
			realText = realTexts[1]
		}
		log.Println("----------------", varName, "!", varValue)

		namedTable.LocalVars[varName] = &Var{
			Name:     varName,
			Value:    varValue,
			RealText: realText,
		}
	} else {
		varName = strings.Trim(varName, `"`)
		log.Println("create table11", varName, "!", namedTable.Name)
		newTable, _ := s.Tables.GetTable(namedTable.Name + " " + varName)

		log.Println("heeeead", len(s.Tables.callStack), namedTable.Name)
		namedTable.LocalTables[newTable.NormalizedName()] = newTable
		s.Tables.pushToStack(newTable)
		log.Println("create table", newTable.Path())
		//s.Tables.reserveName = ""
		//s.Tables.reserveName = strings.Trim(s.getVarName(varName), `"`)
	}
}

// ExitField is called when production field is exited.
func (s *InfoCollector) ExitField(ctx *parser.FieldContext) {}

// EnterFieldsep is called when production fieldsep is entered.
func (s *InfoCollector) EnterFieldsep(ctx *parser.FieldsepContext) {}

// ExitFieldsep is called when production fieldsep is exited.
func (s *InfoCollector) ExitFieldsep(ctx *parser.FieldsepContext) {}

// EnterOperatorOr is called when production operatorOr is entered.
func (s *InfoCollector) EnterOperatorOr(ctx *parser.OperatorOrContext) {}

// ExitOperatorOr is called when production operatorOr is exited.
func (s *InfoCollector) ExitOperatorOr(ctx *parser.OperatorOrContext) {}

// EnterOperatorAnd is called when production operatorAnd is entered.
func (s *InfoCollector) EnterOperatorAnd(ctx *parser.OperatorAndContext) {}

// ExitOperatorAnd is called when production operatorAnd is exited.
func (s *InfoCollector) ExitOperatorAnd(ctx *parser.OperatorAndContext) {}

// EnterOperatorComparison is called when production operatorComparison is entered.
func (s *InfoCollector) EnterOperatorComparison(ctx *parser.OperatorComparisonContext) {}

// ExitOperatorComparison is called when production operatorComparison is exited.
func (s *InfoCollector) ExitOperatorComparison(ctx *parser.OperatorComparisonContext) {}

// EnterOperatorStrcat is called when production operatorStrcat is entered.
func (s *InfoCollector) EnterOperatorStrcat(ctx *parser.OperatorStrcatContext) {}

// ExitOperatorStrcat is called when production operatorStrcat is exited.
func (s *InfoCollector) ExitOperatorStrcat(ctx *parser.OperatorStrcatContext) {}

// EnterOperatorAddSub is called when production operatorAddSub is entered.
func (s *InfoCollector) EnterOperatorAddSub(ctx *parser.OperatorAddSubContext) {}

// ExitOperatorAddSub is called when production operatorAddSub is exited.
func (s *InfoCollector) ExitOperatorAddSub(ctx *parser.OperatorAddSubContext) {}

// EnterOperatorMulDivMod is called when production operatorMulDivMod is entered.
func (s *InfoCollector) EnterOperatorMulDivMod(ctx *parser.OperatorMulDivModContext) {}

// ExitOperatorMulDivMod is called when production operatorMulDivMod is exited.
func (s *InfoCollector) ExitOperatorMulDivMod(ctx *parser.OperatorMulDivModContext) {}

// EnterOperatorBitwise is called when production operatorBitwise is entered.
func (s *InfoCollector) EnterOperatorBitwise(ctx *parser.OperatorBitwiseContext) {}

// ExitOperatorBitwise is called when production operatorBitwise is exited.
func (s *InfoCollector) ExitOperatorBitwise(ctx *parser.OperatorBitwiseContext) {}

// EnterOperatorUnary is called when production operatorUnary is entered.
func (s *InfoCollector) EnterOperatorUnary(ctx *parser.OperatorUnaryContext) {}

// ExitOperatorUnary is called when production operatorUnary is exited.
func (s *InfoCollector) ExitOperatorUnary(ctx *parser.OperatorUnaryContext) {}

// EnterOperatorPower is called when production operatorPower is entered.
func (s *InfoCollector) EnterOperatorPower(ctx *parser.OperatorPowerContext) {}

// ExitOperatorPower is called when production operatorPower is exited.
func (s *InfoCollector) ExitOperatorPower(ctx *parser.OperatorPowerContext) {}

// #5 переменные
// EnterNumber is called when production number is entered.
func (s *InfoCollector) EnterNumber(ctx *parser.NumberContext) {
	s.pickVar(ctx)
}

// ExitNumber is called when production number is exited.
func (s *InfoCollector) ExitNumber(ctx *parser.NumberContext) {}

// #5 переменные
// EnterStringg is called when production stringg is entered.
func (s *InfoCollector) EnterStringg(ctx *parser.StringgContext) {
	log.Println("bbbbbbbbbbbbb")
	s.pickVar(ctx)
}

// ExitStringg is called when production stringg is exited.
func (s *InfoCollector) ExitStringg(ctx *parser.StringgContext) {}

// 534 -> 490 -> 462
