package main

import (
	"gocompiler/internal/constants"
	"gocompiler/internal/expressions"
	"gocompiler/internal/fsm"
	"gocompiler/internal/repository"
	"gocompiler/internal/visualizer"
	"log"
)

func main() {

	buildTheory()

	var cache = repository.NewCache()

	log.Println("Добро пожаловать. Введите номер действия, которое хотите выполнить:")
	log.Println("1. Построить НКА по регулярному выражению")
	log.Println("2. Построить ДКА по НКА")
	log.Println("3. Минимизировать КА")
	log.Println("0. Выйти")
	var mainCode = -1

	for mainCode != 0 {
		repository.GetInt(&mainCode, `Ваш выбор?`)
		switch mainCode {
		case constants.Action2NKA:
			var text string
			repository.GetString(&text, `Введите регулярное выражение`)
			var rw = expressions.NewRW(text)
			kda := rw.ToENKA()
			kda.RemoveShortCircuits()
			kda.ReplaceEpsilons()
			kda.ReplaceEqualEdges()
			kda.AutoDetectFirstLast()
			err := visualizer.VisualizeFSM(kda, "./assets", "main.dot")
			if err != nil {
				log.Printf("Не удалось визуализировтаь граф: %s", err)
				break
			}
			cache.Put(mainCode, kda)
		case constants.Action2DKA:
			repository.LoadGraf(
				constants.Action2NKA,
				mainCode,
				*cache,
				func(graf *fsm.FSM) error {
					graf = graf.ToDFA().ReplaceEqualEdges()
					return nil
				})
		case constants.ActionMinimize:
			repository.LoadGraf(
				constants.Action2DKA,
				mainCode,
				*cache,
				func(graf *fsm.FSM) error {
					f := fsm.NewDRFromFS(*graf).R().D().R().D()
					*graf = f.FSM
					return nil
				})
		}

	}
}

func buildTheory() {
	var rw = expressions.NewRW("(a|b)*abb")
	//var rw = expressions.NewRW("(00|01|10|11)*")
	kda := rw.ToENKA()
	visualizer.MustVisualizeFSM(kda, "./assets/theory", "nka_0.dot")
	//kda.RemoveShortCircuits()
	//kda.ReplaceEpsilons()
	//kda.ReplaceEqualEdges()
	kda.AutoDetectFirstLast()
	visualizer.MustVisualizeFSM(kda, "./assets/theory", "nka_1.dot")
}

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
