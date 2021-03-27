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
	log.Println("4. Моделировать КА")
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
					*graf = *graf.ToDFA()
					return nil
				})
		case constants.ActionModel:
			var text, example string
			repository.GetString(&text, `Введите регулярку:`)
			repository.GetString(&example, `Введите пример:`)
			var rw = expressions.NewRW(text)
			kda := rw.ToENKA().ToDFA()
			found := kda.ContainString(example)
			if found {
				log.Println("Да")
			} else {
				log.Println("Нет")
			}
		}

	}
}

func buildTheory() {
	var rw = expressions.NewRW("(a|b)*abb")
	kda := rw.ToENKA()
	visualizer.MustVisualizeFSM(kda, "./assets/theory", "nka.dot")
}

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
// (00|01|10|11)*
