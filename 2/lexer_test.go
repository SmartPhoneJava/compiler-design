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
	err = lexer.Validate(`{ a = 5 ; }`, false)

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
	err = lexer.Validate(`{ a = a + a ; }`, false)

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
	err = lexer.Validate(`{ a = ( a - b ) + a ; }`, false) // true

	if err != nil {
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexer4(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ b = 1 + 1 / 5 > 3 + 7 ; }`, false)

	if err != nil {
		lexer.Print("")
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexerNegative1(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ b = 123 232 ; }`, false)

	if err == nil {
		lexer.Print("")
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexerNegative2(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ b = + + + ; }`, false)

	if err == nil {
		lexer.Print("")
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

func TestLexer5(t *testing.T) {
	var inputPath = "g5.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = lexer.Validate(`{ b = 2 + ( 1 + ( a + a ) ) ; }`, false)

	if err != nil {
		lexer.Print("")
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// func TestLexerNegative2(t *testing.T) {
// 	var inputPath = "g5.json"

// 	c, err := parsing.MakeGrammar(inputPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	lexer, err := c.ToLexer()
// 	if err != nil {
// 		log.Fatalf("%s", err.Error())
// 	}
// 	err = lexer.Validate(`{ a = 1 > < ; }`, false)

// 	if err == nil {
// 		lexer.Print("")
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}
// }

// func BenchmarkLexer1(b *testing.B) {
// 	var inputPath = "g5.json"

// 	c, err := parsing.MakeGrammar(inputPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	lexer, err := c.ToLexer()
// 	if err != nil {
// 		log.Fatalf("%s", err.Error())
// 	}
// 	err = lexer.Validate(`{ b = 1 + 1 / 5 > 3 + 7 ; }`, false)

// }
