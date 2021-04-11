package main

import (
	"lab2/internal"
	"log"
	"testing"
)

// https://studizba.com/files/show/djvu/3050-1-tom-1.html
// пример 2.24, 174 страница
func Test_RemoveChains1(t *testing.T) {
	var G = internal.CFR{
		N: []string{"E", "T", "F"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{"E"},
		P: internal.Rules{
			{From: "E", To: "E+T"},
			{From: "E", To: "T"},
			{From: "T", To: "T*F"},
			{From: "T", To: "(E)"},
			{From: "T", To: "a"},
			{From: "F", To: "(E)"},
			{From: "F", To: "a"},
		},
	}

	var expected = internal.CFR{
		N: []string{"E", "T", "F"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{"E"},
		P: internal.Rules{
			{From: "E", To: "E+T"},
			{From: "E", To: "T*F"},
			{From: "E", To: "(E)"},
			{From: "E", To: "a"},
			{From: "T", To: "T*F"},
			{From: "T", To: "(E)"},
			{From: "T", To: "a"},
			{From: "F", To: "(E)"},
			{From: "F", To: "a"},
		},
	}

	var real = G.RemoveChains()

	if err := expected.IsSame(real); err != nil {
		log.Println("Ожидалось:", expected.N)
		log.Println("Получено:", real.N)
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
