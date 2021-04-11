package main

import (
	"lab2/internal"
	"log"
	"testing"
)

// file:///home/artyom/labs/bauman/10/compiler-design/2/formal.languages.theory.3.pdf
func Test_RemoveUnreachableNonterminal(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "A", "B", "C", "D"},
		T: []string{"a", "b", "c", "d", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aAB"},
			{From: "S", To: "C"},
			{From: "D", To: "cDc"},
			{From: "D", To: "d"},
			{From: "C", To: "aCD"},
			{From: "A", To: "bA"},
			{From: "A", To: "a"},
			{From: "A", To: "e"},
			{From: "B", To: "b"},
		},
	}

	var expectedNoneGenerating = internal.CFR{
		N: []string{"S", "A", "B", "D"},
		T: []string{"a", "b", "c", "d", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aAB"},
			{From: "D", To: "cDc"},
			{From: "D", To: "d"},
			{From: "A", To: "bA"},
			{From: "A", To: "a"},
			{From: "A", To: "e"},
			{From: "B", To: "b"},
		},
	}

	var expectedUnreachable = internal.CFR{
		N: []string{"S", "A", "B"},
		T: []string{"a", "b", "c", "d", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aAB"},
			{From: "A", To: "bA"},
			{From: "A", To: "a"},
			{From: "A", To: "e"},
			{From: "B", To: "b"},
		},
	}

	var realNoneGenerating = G.RemoveNongeneratingNonterminal()

	if err := realNoneGenerating.IsSame(expectedNoneGenerating); err != nil {
		log.Println("Ожидалось:", expectedNoneGenerating.N)
		log.Println("Получено:", realNoneGenerating.N)
		log.Println("Ожидалось:", expectedNoneGenerating.P)
		log.Println("Получено:", realNoneGenerating.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}

	var realUnreachable = realNoneGenerating.RemoveUnreachableNonterminal()

	if err := realUnreachable.IsSame(expectedUnreachable); err != nil {
		log.Println("Ожидалось:", expectedUnreachable.N)
		log.Println("Получено:", realUnreachable.N)
		log.Println("Ожидалось:", expectedUnreachable.P)
		log.Println("Получено:", realUnreachable.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// http://mathhelpplanet.com/static.php?p=privedennaya-forma-ks-grammatiki
func Test_RemoveUnreachableNonterminal2(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "A", "B", "C"},
		T: []string{"a", "b"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aA"},
			{From: "S", To: "bB"},
			{From: "A", To: "bAa"},
			{From: "B", To: "aB"},
			{From: "B", To: "bS"},
			{From: "B", To: "a"},
			{From: "B", To: "b"},
			{From: "C", To: "BaA"},
		},
	}

	var expectedNoneGenerating = internal.CFR{
		N: []string{"S", "B"},
		T: []string{"a", "b"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "bB"},
			{From: "B", To: "aB"},
			{From: "B", To: "bS"},
			{From: "B", To: "a"},
			{From: "B", To: "b"},
		},
	}

	var realNoneGenerating = G.RemoveNongeneratingNonterminal()

	if err := expectedNoneGenerating.IsSame(realNoneGenerating); err != nil {
		log.Println("1.Ожидалось:", expectedNoneGenerating.N)
		log.Println("1. Получено:", realNoneGenerating.N)
		log.Println("1. Ожидалось:", expectedNoneGenerating.P)
		log.Println("1. Получено:", realNoneGenerating.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
