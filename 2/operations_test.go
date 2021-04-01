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

// // https://studizba.com/files/show/djvu/3050-1-tom-1.html
// // Пример 2.22
// func Test_RemoveUnreachableNonterminal2(t *testing.T) {
// 	var G = internal.CFR{
// 		N: []string{"I", "A", "B", "C", "R"},
// 		T: []string{"int", "char", ",", ";"},
// 		S: []string{"I"},
// 		P: internal.Rules{
// 			{From: "I", To: "ABC"},
// 			{From: "A", To: "int"},
// 			{From: "A", To: "char"},
// 			{From: "B", To: "IR"},
// 			{From: "R", To: ",IR"},
// 			{From: "R", To: "I"},
// 			{From: "C", To: ";"},
// 		},
// 	}

// 	var expected = internal.CFR{
// 		N: []string{"I", "A", "B", "C", "R"},
// 		T: []string{"int", "char", ",", ";"},
// 		S: []string{"I"},
// 		P: internal.Rules{
// 			{From: "I", To: "ABC"},
// 			{From: "A", To: "int"},
// 			{From: "A", To: "char"},
// 			{From: "B", To: "IR"},
// 			{From: "R", To: ",IR"},
// 			{From: "R", To: "I"},
// 			{From: "C", To: ";"},
// 		},
// 	}

// 	log.Println("aaaa:", G.P)
// 	G = G.ElrWithE()
// 	log.Println("dddddd:", G.P)
// 	var real = G.RemoveNongeneratingNonterminal()

// 	if err := real.IsSame(expected); err != nil {

// 		log.Println("Ожидалось:", expected.N)
// 		log.Println("Получено:", real.N)
// 		log.Println("Ожидалось:", expected.P)
// 		log.Println("Получено:", real.P)
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}
// }

// file:///home/artyom/labs/bauman/10/compiler-design/2/formal.languages.theory.3.pdf
func Test_RemoveLambda1(t *testing.T) {
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

// http://mathhelpplanet.com/static.php?p=privedennaya-forma-ks-grammatiki
func Test_RemoveLambda2(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "T"},
		T: []string{"a", "b", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aSbT"},
			{From: "S", To: "bTaT"},
			{From: "S", To: "ab"},
			{From: "T", To: "baaST"},
			{From: "T", To: "TT"},
			{From: "T", To: "e"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S", "T"},
		T: []string{"a", "b", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aSbT"},
			{From: "S", To: "aSb"},
			{From: "S", To: "bTaT"},
			{From: "S", To: "bTa"},
			{From: "S", To: "baT"},
			{From: "S", To: "ba"},
			{From: "S", To: "ab"},
			{From: "T", To: "baaST"},
			{From: "T", To: "TT"},
			{From: "T", To: "baaS"},
		},
	}

	var real = G.RemoveLambda()

	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.N)
		log.Println("Получено:", real.N)
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// http://mathhelpplanet.com/static.php?p=privedennaya-forma-ks-grammatiki
func Test_RemoveLambda3(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "B", "C", "A"},
		T: []string{"a", "c", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "ABCd"},
			{From: "A", To: "a"},
			{From: "A", To: "e"},
			{From: "B", To: "AC"},
			{From: "C", To: "c"},
			{From: "C", To: "e"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S", "B", "C", "A"},
		T: []string{"a", "c", "e"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "Ad"},
			{From: "S", To: "ABd"},
			{From: "S", To: "ACd"},
			{From: "S", To: "ABCd"},
			{From: "S", To: "Bd"},
			{From: "S", To: "BCd"},
			{From: "S", To: "Cd"},
			{From: "S", To: "d"},
			{From: "B", To: "A"},
			{From: "B", To: "C"},
			{From: "B", To: "AC"},
			{From: "A", To: "a"},
			{From: "C", To: "c"},
		},
	}

	log.Println("++++++++++++++++++++++++++++++++++")
	var real = G.RemoveLambda()

	if err := expected.IsSame(real); err != nil {
		log.Println("Ожидалось:", expected.N)
		log.Println("Получено:", real.N)
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
