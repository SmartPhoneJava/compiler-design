package parsing

import (
	"errors"
	"strings"

	"github.com/thedevsaddam/gojsonq"
)

func MakeGrammar(filename string) (ContextGrammar, error) {
	cg := ContextGrammar{}
	//здесь термы
	name := gojsonq.New().File(filename).
		From("grammar.-name").Get()
	if name == nil {
		return ContextGrammar{}, errors.New("Файл не найден или не содержит данных")
	}
	cg.Name = name.(string)

	terms := gojsonq.New().File(filename).
		From("grammar.terminalsymbols.term").Get()
	term, _ := terms.([]interface{})
	for i := 0; i < len(term); i++ {
		list := term[i].(map[string]interface{})
		s := Symbol{}
		s.Spell = list["-name"].(string)
		s.Name = list["-spell"].(string)
		s.Type = "term"
		cg.Terms = append(cg.Terms, s)
	}
	//здесь не термы
	nonterms := gojsonq.New().File(filename).
		From("grammar.nonterminalsymbols.nonterm").Get()
	nonterm, ok := nonterms.([]interface{})
	if !ok {
		return cg, errors.New("nonterm, ok := nonterms.([]interface{})")
	}
	for i := 0; i < len(nonterm); i++ {
		list := nonterm[i].(map[string]interface{})
		s := Symbol{}
		s.Spell = list["-name"].(string)
		s.Name = list["-name"].(string)
		s.Type = "nonterm"
		cg.Nonterms = append(cg.Nonterms, s)
	}

	productions := gojsonq.New().File(filename).
		From("grammar.productions.production").Get()
	production, _ := productions.([]interface{})

	for i := 0; i < len(production); i++ {
		r := Rule{}
		list := production[i].(map[string]interface{})
		lhs := list["lhs"].(map[string]interface{})
		r.LeftSide = Symbol{
			Name:  lhs["-name"].(string),
			Spell: lhs["-name"].(string),
		}
		for i := 0; i < len(cg.Terms); i++ {
			if r.LeftSide.Name == cg.Terms[i].Name {
				r.LeftSide.Type = "term"
			}
		}
		if r.LeftSide.Type == "" {
			r.LeftSide.Type = "nonterm"
		}

		rhs := list["rhs"].(map[string]interface{})["symbol"]

		if _, b := rhs.(interface{}); b {
			//это всегда
			a, _ := rhs.(interface{})
			//это иногда
			inside, b := a.([]interface{})
			if b {
				//в правиле несколько символов
				for i := 0; i < len(inside); i++ {
					wow := inside[i].(map[string]interface{})
					s := Symbol{
						Type:  wow["-type"].(string),
						Spell: wow["-name"].(string),
						Name:  wow["-name"].(string),
					}

					for _, t := range cg.Terms {
						s.Name = strings.ReplaceAll(s.Name, t.Spell, t.Name)
					}
					r.RightSide = append(r.RightSide, s)
				}
			} else {
				inside2, ok := a.(map[string]interface{})
				if !ok {
					return cg, errors.New("	inside2, ok := a.(map[string]interface{})")
				}
				s := Symbol{
					Type:  inside2["-type"].(string),
					Spell: inside2["-name"].(string),
					Name:  inside2["-name"].(string),
				}
				for _, t := range cg.Terms {
					s.Name = strings.ReplaceAll(s.Name, t.Spell, t.Name)
				}
				r.RightSide = append(r.RightSide, s)
			}
		}
		cg.Rules = append(cg.Rules, r)
	}
	//стартовый символ
	startInterface := gojsonq.New().File(filename).
		From("grammar.startsymbol").Get()
	start, flag := startInterface.(map[string]interface{})
	//это один символ
	if flag {
		//символ один
		startName := start["-name"].(string)
		for i := 0; i < len(cg.Nonterms); i++ {
			if cg.Nonterms[i].Name == startName {
				cg.Start = append(cg.Start, cg.Nonterms[i])
			}
		}

		for i := 0; i < len(cg.Terms); i++ {
			if cg.Terms[i].Name == startName {
				cg.Start = append(cg.Start, cg.Terms[i])
			}
		}
	}

	return cg, nil
}
