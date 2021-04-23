package main

import (
	"lab2/internal/g5"
	"lab2/parsing"
	"log"
	"testing"
)

func TestFirst(t *testing.T) {
	var inputPath = "test1.json"

	c, err := parsing.MakeGrammar(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	lexer, err := c.ToLexer()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	var expected = g5.NoneTermsSet{
		"E": g5.StringSet{
			"n": nil,
			"(": nil,
		},
		"E'": g5.StringSet{
			"+": nil,
			"e": nil,
		},
		"T": g5.StringSet{
			"n": nil,
			"(": nil,
		},
		"T'": g5.StringSet{
			"*": nil,
			"e": nil,
		},
		"F": g5.StringSet{
			"n": nil,
			"(": nil,
		},
	}
	realOne := lexer.ConstructFirst()
	for nt, ts := range expected {
		for term := range ts {
			_, ok := realOne[nt]
			if !ok {
				t.Fatalf("Не сошлось")
			}
			_, ok = realOne[nt][term]
			if !ok {
				expected.Println("Ожидалось:")
				realOne.Println("\nПолучено:")
				t.Fatalf("Не сошлось")
			}
		}
	}
	if err != nil {
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// func TestFollow(t *testing.T) {
// 	var inputPath = "test1.json"

// 	c, err := parsing.MakeGrammar(inputPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	lexer, err := c.ToLexer()
// 	if err != nil {
// 		log.Fatalf("%s", err.Error())
// 	}
// 	var expected = g5.NoneTermsSet{
// 		"E": g5.StringSet{
// 			"$": nil,
// 			")": nil,
// 		},
// 		"E'": g5.StringSet{
// 			"$": nil,
// 			")": nil,
// 		},
// 		"T": g5.StringSet{
// 			"+": nil,
// 			"$": nil,
// 			")": nil,
// 		},
// 		"T'": g5.StringSet{
// 			"+": nil,
// 			"$": nil,
// 			")": nil,
// 		},
// 		"F": g5.StringSet{
// 			"*": nil,
// 			"+": nil,
// 			"$": nil,
// 			")": nil,
// 		},
// 	}
// 	first := lexer.ConstructFirst()
// 	realOne := lexer.ConstructFollow(first)
// 	for nt, ts := range expected {
// 		for term := range ts {
// 			_, ok := realOne[nt]
// 			if !ok {
// 				t.Fatalf("Не сошлось")
// 			}
// 			_, ok = realOne[nt][term]
// 			if !ok {
// 				expected.Println("Ожидалось:")
// 				realOne.Println("\nПолучено:")
// 				t.Fatalf("Не сошлось")
// 			}
// 		}
// 	}
// 	if err != nil {
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}
// }
