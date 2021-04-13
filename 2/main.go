package main

import (
	"fmt"
	"lab2/internal"
	"lab2/parsing"
	"log"
	"os"
)

const (
	EliminateLeftRecursion = iota + 1
	LeftFactorization
	LeftFactorizationWithoutRec
	RemoveUnreachable
	RemoveNongenerating
	RemoveUseless
	RemoveChains
	RemoveLambda
	Bring
)

func main() {

	var inputPath string
	log.Println("Введите название файла, откуда будет загружена грамматика:")
	fmt.Scanf("%s\n", &inputPath)

	// var outputPath string
	// log.Println("Введите название файла, куда будет записана новая грамматика:")
	// fmt.Scanf("%s\n", &outputPath)

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	var cfr = c.ToInternal()
	var newCfr internal.CFR
	cfr.Print("Грамматика загружена:")

	var action = 1
	for action > 0 {
		log.Println("Какое действие применить?(Впишите цифру от 1 до 9)")
		log.Println("1. Устранить левую рекурсию")
		log.Println("2. Провести левую факторизацию")
		log.Println("3. Устранить левую рекурсию с факторизацией")
		log.Println("4. Удалить недостижимые символы")
		log.Println("5. Удалить непорождающие символы")
		log.Println("6. Удалить бесполезные символы")
		log.Println("7. Удалить циклы")
		log.Println("8. Удалить лямбда-переходы")
		log.Println("9. Привести грамматику")
		log.Println("Любой другой код расценивается как команда 'завершить работу программы'")
		fmt.Scanf("%d\n", &action)
		if action < 1 || action > 9 {
			os.Exit(0)
		}

		switch action {
		case EliminateLeftRecursion:
			newCfr = cfr.ElrWithE2(true)
		case LeftFactorization:
			newCfr = cfr.LeftFactorization()
		case LeftFactorizationWithoutRec:
			newCfr = cfr.ElrWithE2(true).LeftFactorization()
		case RemoveUnreachable:
			newCfr = cfr.RemoveUnreachableNonterminal()
		case RemoveNongenerating:
			newCfr = cfr.RemoveNongeneratingNonterminal()
		case RemoveUseless:
			newCfr = cfr.RemoveUselessNonterms()
		case RemoveChains:
			newCfr = cfr.RemoveChains()
		case RemoveLambda:
			newCfr = cfr.RemoveLambda()
		case Bring:
			newCfr = cfr.Bring()
		}

		newCfr.Print("Операция выполнена. Полученная грамматика:")
	}

	// var grammar = internal.CFR{}
	// var str string
	// log.Println("Введите набор нетермов через пробел")
	// fmt.Scanf("%s\n", &str)

	// grammar.N = strings.Split(str, " ")
	// if len(grammar.N) == 0 {
	// 	log.Fatal("Нетермы не указаны")
	// }

	// log.Println("Введите набор термов через пробел")
	// fmt.Scanf("%s\n", &str)

	// log.Println("Какой терм стартовый?")

	// log.Println("Введите правила")

	// log.Println("Добро пожаловать. Введите номер действия, которое хотите выполнить:")
	// log.Println("1. Построить НКА по регулярному выражению")
	// log.Println("2. Построить ДКА по НКА")
	// log.Println("3. Минимизировать КА")
	// log.Println("4. Моделировать КА")
	// log.Println("0. Выйти")
	// var mainCode = -1

	// for mainCode != 0 {
	// 	repository.GetInt(&mainCode, `Ваш выбор?`)
	// 	switch mainCode {
	// 	case constants.Action2NKA:
	// 		var text string
	// 		repository.GetString(&text, `Введите регулярное выражение`)
	// 		var rw = expressions.NewRW(text)
	// 		kda := rw.ToENKA()
	// 		err := visualizer.VisualizeFSM(kda, "./assets", "main.dot")
	// 		if err != nil {
	// 			log.Printf("Не удалось визуализировтаь граф: %s", err)
	// 			break
	// 		}
	// 		cache.Put(mainCode, kda)
	// 	case constants.Action2DKA:
	// 		repository.LoadGraf(
	// 			constants.Action2NKA,
	// 			mainCode,
	// 			*cache,
	// 			func(graf *fsm.FSM) error {
	// 				*graf = *graf.ToDFA()
	// 				return nil
	// 			})
	// 	case constants.ActionModel:
	// 		var text, example string
	// 		repository.GetString(&text, `Введите регулярку:`)
	// 		repository.GetString(&example, `Введите пример:`)
	// 		var rw = expressions.NewRW(text)
	// 		kda := rw.ToENKA().ToDFA()
	// 		found := kda.ContainString(example)
	// 		if found {
	// 			log.Println("Да")
	// 		} else {
	// 			log.Println("Нет")
	// 		}
	// 	}

	// }
}

// 464
