package main

import (
	"fmt"
	"lab2/internal"
	"lab2/parsing"
	"log"
	"os"
	"strings"
	"time"

	"github.com/buger/goterm"
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
	goterm.Clear() // Clear current screen

	// for {
	// 	// By moving cursor to top-left position we ensure that console output
	// 	// will be overwritten each time, instead of adding new.
	// 	goterm.MoveCursor(1, 1)

	// 	goterm.Println("Current Time:", time.Now().Format(time.RFC1123))

	// 	goterm.Flush() // Call it every time at the end of rendering

	// 	time.Sleep(time.Second)
	// }
	// goterm.Clear()
	mainLab3()
}

func mainLab3() {
	var inputPath = "g5.json"
	//log.Println("Введите название файла, откуда будет загружена грамматика:")
	//fmt.Scanf("%s\n", &inputPath)

	// var outputPath string
	// log.Println("Введите название файла, куда будет записана новая грамматика:")
	// fmt.Scanf("%s\n", &outputPath)

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	lexer.Print("Грамматика загружена:")
	time.Sleep(time.Second)
	var text = "{ }"
	_, err = lexer.Start.GoTo(strings.Split(text, " "), 0)
	if err != nil {
		goterm.Println("Произошла ошибка", err)
	} else {
		goterm.Println("Успех")
	}
	goterm.Flush()
}

func mainLab2() {

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
}
