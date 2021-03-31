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

func Test_RemoveLambda(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "A", "B", "C"},
		T: []string{"a", "b", "c", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "BC"},
			{From: "S", To: "Ab"},
			{From: "B", To: "e"},
			{From: "C", To: "c"},
			{From: "A", To: "Aa"},
			{From: "A", To: "e"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S", "A", "C"},
		T: []string{"a", "b", "c", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "C"},
			{From: "S", To: "b"},
			{From: "S", To: "Ab"},
			{From: "C", To: "c"},
			{From: "A", To: "Aa"},
			{From: "A", To: "a"},
		},
	}

	var real = G.RemoveLambda().RemoveNongeneratingNonterminal() //.RemoveNongeneratingNonterminal().RemoveUnreachableNonterminal()

	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.N)
		log.Println("Получено:", real.N)
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
