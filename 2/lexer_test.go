package main

import (
	"lab2/parsing"
	"log"
	"testing"
)

func TestLexer1(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ a = 5 ; }`, true)

	if err != nil {
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexer2(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ a = a + a ; }`, true)

	if err != nil {
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexer3(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ a = ( a - b ) + a ; }`, true)

	if err != nil {
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
