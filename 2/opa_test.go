package main

import (
	"errors"
	"lab2/internal/ast"
	"lab2/internal/opa"
	"lab2/parsing"
	"strings"
	"testing"
)

func Test_Opa1(t *testing.T) {
	var inputPath = "assets/grammar/bigexample.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var (
		analyzer = opa.NewAnalyzer(lexer)
		inputRow = "if a or a and a then a = a xor a ;"
	)
	_, outR, err := analyzer.Exec(strings.Split(inputRow, " "))
	if err != nil {
		t.Fatal(err)
	}
	err = ast.Visualize(outR, "assets", "Test_Opa1.dot")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Opa2(t *testing.T) {
	var inputPath = "assets/grammar/bigexample.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var (
		analyzer    = opa.NewAnalyzer(lexer)
		inputRow    = "if a a or a and a then a = a xor a ;"
		expectedErr = errors.New("Введенный код содержит ошибку: ключевое слово `a` не может находиться слева от `a`")
	)
	_, _, err = analyzer.Exec(strings.Split(inputRow, " "))

	if err.Error() != expectedErr.Error() {
		t.Fatalf("Ожидалось '%s', а получено '%s'", expectedErr, err)
	}
}

func Test_Opa3(t *testing.T) {
	var inputPath = "assets/grammar/smallexample.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var (
		analyzer    = opa.NewAnalyzer(lexer)
		inputRow    = "a + + a * a"
		expectedErr = errors.New("Введенный код содержит ошибку: Правила не найдено для &[⏊ E +]")
	)
	_, _, err = analyzer.Exec(strings.Split(inputRow, " "))

	if err.Error() != expectedErr.Error() {
		t.Fatalf("Ожидалось '%s', а получено '%s'", expectedErr, err)
	}
}

func Test_Opa4(t *testing.T) {
	var inputPath = "assets/grammar/smallexample.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var (
		analyzer = opa.NewAnalyzer(lexer)
		inputRow = "a + a * a"
	)
	_, outR, err := analyzer.Exec(strings.Split(inputRow, " "))
	if err != nil {
		t.Fatal(err)
	}
	err = ast.Visualize(outR, "assets", "Test_Opa4.dot")
	if err != nil {
		t.Fatal(err)
	}
}
