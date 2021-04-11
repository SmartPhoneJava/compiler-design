package main

import (
	"lab2/internal"
	"log"
	"testing"
)

// http://mathhelpplanet.com/static.php?p=privedennaya-forma-ks-grammatiki
func Test_Bring(t *testing.T) {
	var G = internal.CFR{
		N: []string{"S"},
		T: []string{"b", "a"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aSa"},
			{From: "S", To: "bSb"},
			{From: "S", To: "a"},
			{From: "S", To: "b"},
			{From: "S", To: "e"},
		},
	}

	var expected = internal.CFR{
		N: []string{"S"},
		T: []string{"b", "a"},
		S: []string{"S"},
		P: internal.Rules{
			{From: "S", To: "aSa"},
			{From: "S", To: "bSb"},
			{From: "S", To: "aa"},
			{From: "S", To: "bb"},
			{From: "S", To: "a"},
			{From: "S", To: "b"},
		},
	}

	var real = G.Bring()

	if err := expected.IsSame(real); err != nil {
		log.Println("Ожидалось:", expected.N)
		log.Println("Получено:", real.N)
		log.Println("Ожидалось:", expected.P)
		log.Println("Получено:", real.P)
		t.Fatalf("Ожидание и реальность не сошлись: %s", err)
	}
}
