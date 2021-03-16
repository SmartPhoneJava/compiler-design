package main

import (
	"log"
	"testing"
)

// https://studizba.com/files/show/djvu/3050-1-tom-1.html
// страница 179 справа внизу
func TestExample2_27(t *testing.T) {
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

// Пример 4.9 "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание", стр 277
func TestExample4_9(t *testing.T) {
	log.Println("!!!!!!!!!!!!!!!!!!")
	var G = CFR{
		N: []string{"S", "A"},
		T: []string{"a", "b", "c", "d"},
		S: []string{},
		P: Rules{
			{From: "S", To: "Aa"},
			{From: "S", To: "b"},
			{From: "A", To: "Ac"},
			{From: "A", To: "Sd"},
			{From: "A", To: "e"},
		},
	}

	var expected = CFR{
		N: []string{"S", "A", "A'"},
		T: []string{"a", "b", "c", "d"},
		S: []string{},
		P: Rules{
			{From: "S", To: "Aa"},
			{From: "S", To: "b"},
			{From: "A", To: "bdA'"},
			{From: "A", To: "A'"},
			{From: "A'", To: "cA'"},
			{From: "A'", To: "adA'"},
			{From: "A'", To: "e"},
		},
	}

	var real = G.ElrWithE()
	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}

}
