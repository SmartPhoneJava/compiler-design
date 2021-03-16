package main

import (
	"log"
	"testing"
)

// https://studizba.com/files/show/djvu/3050-1-tom-1.html
// страница 179 справа внизу
func TestRemoveLeftReverse(t *testing.T) {
	var G = CFR{
		N: []string{"E", "T", "F"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{},
		P: Rules{
			{From: "E", To: "E+T"},
			{From: "E", To: "T"},
			{From: "T", To: "T*F"},
			{From: "T", To: "F"},
			{From: "F", To: "(E)"},
			{From: "F", To: "a"},
		},
	}

	var expected = CFR{
		N: []string{"E", "T", "F", "E'", "T'"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{},
		P: Rules{
			{From: "E", To: "TE'"},
			{From: "E", To: "T"},
			{From: "E'", To: "+T"},
			{From: "E'", To: "+TE'"},
			{From: "T", To: "FT'"},
			{From: "T", To: "F"},
			{From: "T'", To: "*F"},
			{From: "T'", To: "*FT'"},
			{From: "F", To: "(E)"},
			{From: "F", To: "a"},
		},
	}

	var real = G.EliminateLeftRecursion()
	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}

}
