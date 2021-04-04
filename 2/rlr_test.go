package main

import (
	"lab2/internal"
	"log"
	"testing"
)

// https://studizba.com/files/show/djvu/3050-1-tom-1.html
// страница 179 справа внизу
func TestExample2_27(t *testing.T) {
	var G = internal.CFR{
		N: []string{"E", "T", "F"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{},
		P: internal.Rules{
			{From: "E", To: "E+T"},
			{From: "E", To: "T"},
			{From: "T", To: "T*F"},
			{From: "T", To: "F"},
			{From: "F", To: "(E)"},
			{From: "F", To: "a"},
		},
	}

	var expected = internal.CFR{
		N: []string{"E", "T", "F", "E'", "T'"},
		T: []string{"+", "*", "(", ")", "a"},
		S: []string{},
		P: internal.Rules{
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
	var G = internal.CFR{
		N: []string{"S", "A"},
		T: []string{"a", "b", "c", "d"},
		S: []string{},
		P: internal.Rules{
			{From: "S", To: "Aa"},
			{From: "S", To: "b"},
			{From: "A", To: "Ac"},
			{From: "A", To: "Sd"},
			{From: "A", To: "e"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S", "A", "A'"},
		T: []string{"a", "b", "c", "d"},
		S: []string{},
		P: internal.Rules{
			{From: "S", To: "Aa"},
			{From: "S", To: "b"},
			{From: "A", To: "bdA'"},
			{From: "A", To: "A'"},
			{From: "A'", To: "cA'"},
			{From: "A'", To: "adA'"},
			{From: "A'", To: "e"},
		},
	}

	var real = G.ElrWithE2(true)
	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// Пример 4.7, стр 276, "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание"
// Выражение 4.2, стр 254
func TestExample4_7(t *testing.T) {
	var G = internal.CFR{
		N: []string{"E", "T", "F"},
		T: []string{"+", "(", ")", "i", "d"},
		S: []string{},
		P: internal.Rules{
			{From: "E", To: "E+T"},
			{From: "E", To: "T"},
			{From: "T", To: "T*F"},
			{From: "T", To: "F"},
			{From: "F", To: "(E)"},
			{From: "F", To: "id"},
		},
	}

	var expected = internal.CFR{
		N: []string{"E", "T", "F", "E'", "T'"},
		T: []string{"+", "(", ")", "i", "d"},
		S: []string{},
		P: internal.Rules{
			{From: "E", To: "TE'"},
			{From: "E'", To: "+TE'"},
			{From: "E'", To: "e"},
			{From: "T", To: "FT'"},
			{From: "T", To: "*FT'"},
			{From: "T'", To: "e"},
			{From: "F", To: "(E)"},
			{From: "F", To: "id"},
		},
	}

	var real = G.ElrWithE2(true)
	if err := real.IsSame(expected); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// // Пример 4.7, стр 276, "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание"
// // Выражение 4.2, стр 254
// func TestExample4_7_f(t *testing.T) {
// 	var G = internal.CFR{
// 		N: []string{"E", "T", "F"},
// 		T: []string{"+", "(", ")", "i", "d"},
// 		S: []string{},
// 		P: internal.Rules{
// 			{From: "E", To: "E+T"},
// 			{From: "E", To: "T"},
// 			{From: "T", To: "T*F"},
// 			{From: "T", To: "F"},
// 			{From: "F", To: "(E)"},
// 			{From: "F", To: "id"},
// 		},
// 	}

// 	var expected = internal.CFR{
// 		N: []string{"E", "T", "F", "E'", "T'"},
// 		T: []string{"+", "(", ")", "i", "d"},
// 		S: []string{},
// 		P: internal.Rules{
// 			{From: "E", To: "TE'"},
// 			{From: "E'", To: "+TE'"},
// 			{From: "E'", To: "e"},
// 			{From: "T", To: "FT'"},
// 			{From: "T", To: "*FT'"},
// 			{From: "T'", To: "e"},
// 			{From: "F", To: "(E)"},
// 			{From: "F", To: "id"},
// 		},
// 	}

// 	var real = G.LeftFactorization()
// 	if err := real.IsSame(expected); err != nil {
// 		log.Println("Ожидалось:", expected.P)
// 		log.Println("Получено:", real.P)
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}
// }

// Левая факторизация
// Пример 4.11, стр 279, "Ахо, Сети, Ульман. Компиляторы. Принципы, технологии, инструменты, 2008, 2-ое издание"
func TestExample4_11(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S", "E"},
		T: []string{"i", "t", "e", "a", "b"},
		S: []string{},
		P: internal.Rules{
			{From: "S", To: "iEtS"},
			{From: "S", To: "iEtSeS"},
			{From: "S", To: "a"},
			{From: "E", To: "b"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S", "E", "S'"},
		T: []string{"i", "t", "e", "a", "b"},
		S: []string{},
		P: internal.Rules{
			{From: "S", To: "iEtSS'"},
			{From: "S", To: "a"},
			{From: "S'", To: "eS"},
			{From: "S'", To: "e"},
			{From: "E", To: "b"},
		},
	}

	var real = G.LeftFactorization()
	if err := expected.IsSame(real); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}

// // Удаление левой рекурсии
// // https://neerc.ifmo.ru/wiki/index.php?title=Устранение_левой_рекурсии
// func TestExampleMy1(t *testing.T) {
// 	var G = internal.CFR{
// 		N: []string{"A", "S"},
// 		T: []string{"a", "b", "c"},
// 		S: []string{},
// 		P: internal.Rules{
// 			{From: "A", To: "Sa"},
// 			{From: "S", To: "Sb"},
// 			{From: "S", To: "SAc"},
// 			{From: "S", To: "b"},
// 		},
// 	}

// 	var expected = internal.CFR{
// 		N: []string{"A", "S", "S'"},
// 		T: []string{"a", "b", "c"},
// 		S: []string{"A"},
// 		P: internal.Rules{
// 			{From: "A", To: "Sa"},
// 			{From: "S", To: "bS'"},
// 			{From: "S'", To: "bS'"},
// 			{From: "S'", To: "acS'"},
// 			{From: "S'", To: "b"},
// 			{From: "S'", To: "ac"},
// 		},
// 	}

// 	log.Println("++++++++++++++++++++++++++++++++++")
// 	log.Println("Оригинал:", G.P)
// 	var real = G.ElrWithE2()
// 	if err := expected.IsSame(real); err != nil {
// 		log.Println("Ожидалось:", expected.P)
// 		log.Println("Получено:", real.P)
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}

// }

// // Удаление левой рекурсии
// // https://neerc.ifmo.ru/wiki/index.php?title=Устранение_левой_рекурсии
// func TestExampleMy2(t *testing.T) {
// 	var G = internal.CFR{
// 		N: []string{"A", "S"},
// 		T: []string{"a", "b"},
// 		S: []string{"S"},
// 		P: internal.Rules{
// 			{From: "A", To: "Sa"},
// 			{From: "A", To: "Aa"},
// 			{From: "S", To: "Ab"},
// 		},
// 	}

// 	var expected = internal.CFR{
// 		N: []string{"A", "S"},
// 		T: []string{"a", "b"},
// 		S: []string{"S"},
// 		P: internal.Rules{
// 			{From: "A", To: "SaA'"},
// 			{From: "A", To: "Sa"},
// 			{From: "A'", To: "a"},
// 			{From: "A'", To: "aA'"},
// 			{From: "S", To: "Ab"},
// 		},
// 	}

// 	log.Println("========================")
// 	log.Println("Оригинал:", G.P)
// 	var real = G.ElrWithE()
// 	if err := expected.IsSame(real); err != nil {
// 		log.Println("Ожидалось:", expected.P)
// 		log.Println("Получено:", real.P)
// 		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
// 	}

// }

// Удаление левой рекурсии
// http://espressocode.top/removing-direct-and-indirect-left-recursion-in-a-grammar/
func TestExampleMy3(t *testing.T) {
	var G = internal.CFR{
		N: []string{"A1", "A2", "A3"},
		T: []string{"a", "b"},
		S: []string{"A1"},
		P: internal.Rules{
			{From: "A1", To: "A2A3"},
			{From: "A2", To: "A3A1"},
			{From: "A2", To: "b"},
			{From: "A3", To: "A1A1"},
			{From: "A3", To: "a"},
		},
	}

	var expected = internal.CFR{
		N: []string{"A1", "A2", "A3", "A3'"},
		T: []string{"a", "b"},
		S: []string{"A1"},
		P: internal.Rules{
			{From: "A1", To: "A2A3"},
			{From: "A2", To: "A3A1"},
			{From: "A2", To: "b"},
			{From: "A3", To: "a"},
			{From: "A3", To: "bA3A1"},
			{From: "A3", To: "aA3'"},
			{From: "A3", To: "bA3A1A3'"},
			{From: "A3'", To: "A1A3A1"},
			{From: "A3'", To: "A1A3A1A3'"},
		},
	}

	log.Println("++++++++++++++++++++++++++++++++++")
	log.Println("Оригинал:", G.P)
	var real = G.ElrWithE2(false)
	if err := expected.IsSame(real); err != nil {
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}

}
