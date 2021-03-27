package main

import (
	"gocompiler/internal/expressions"
	"testing"
)

// Пример 3.19
func TestContain1(t *testing.T) {
	var (
		rwStr   = "(a|b)*abb"
		example = "ababb"
	)
	if !expressions.NewRW(rwStr).ToENKA().ToDFA().ContainString(example) {
		t.Fatalf("Получено Нет, а должно Да")
	}
}

func TestContain2(t *testing.T) {
	var (
		rwStr   = "(a|b)*abb"
		example = "abba"
		dfa     = expressions.NewRW(rwStr).ToENKA().ToDFA()
		//folder  = "assets/test/contain/1"
	)
	//visualizer.MustVisualizeFSM(dfa, folder, "real.dot")
	if dfa.ContainString(example) {
		t.Fatalf("Получено Да, а должно Нет")
	}
}

func TestContai3(t *testing.T) {
	var (
		rwStr   = "(a|b)*abb(a|b)*"
		example = "aaabbaaaaa"
		dfa     = expressions.NewRW(rwStr).ToENKA().ToDFA()
	)
	if !dfa.ContainString(example) {
		t.Fatalf("Получено Нет, а должно Да")
	}
}

func TestContai4(t *testing.T) {
	var (
		rwStr   = "a|b"
		example = "ab"
		//folder  = "assets/test/contain/2"
	)
	var nka = expressions.NewRW(rwStr).ToENKA()
	//visualizer.MustVisualizeFSM(nka, folder, "v1.dot")
	var dfa = nka.ToDFA()
	//visualizer.MustVisualizeFSM(dfa, folder, "v2.dot")
	if dfa.ContainString(example) {
		t.Fatalf("Получено Да, а должно Нет")
	}
}

func TestContai5(t *testing.T) {
	var (
		rwStr   = "a*|b"
		example = "aaaaaa"
		//folder  = "assets/test/contain/3"
	)
	var nka = expressions.NewRW(rwStr).ToENKA()
	// visualizer.MustVisualizeFSM(nka, folder, "v1.dot")
	var dfa = nka.ToDFA()
	// visualizer.MustVisualizeFSM(&dr.FSM, folder, "v2.dot")
	if !dfa.ContainString(example) {
		t.Fatalf("Получено Нет, а должно Да")
	}
}

func TestContai6(t *testing.T) {
	var (
		rwStr   = "a*|b"
		example = "bba"
		//folder  = "assets/test/contain/2"
	)
	var nka = expressions.NewRW(rwStr).ToENKA()
	//visualizer.MustVisualizeFSM(nka, folder, "v1.dot")
	var dfa = nka.ToDFA()
	//visualizer.MustVisualizeFSM(dfa, folder, "v2.dot")
	if dfa.ContainString(example) {
		t.Fatalf("Получено Да, а должно Нет")
	}
}
func TestContai7(t *testing.T) {
	var (
		rwStr   = "a*|b"
		example = "ba"
		//folder  = "assets/test/contain/2"
	)
	var nka = expressions.NewRW(rwStr).ToENKA()
	//visualizer.MustVisualizeFSM(nka, folder, "v1.dot")
	var dfa = nka.ToDFA()
	//visualizer.MustVisualizeFSM(dfa, folder, "v2.dot")
	if dfa.ContainString(example) {
		t.Fatalf("Получено Да, а должно Нет")
	}
}

func TestContai8(t *testing.T) {
	var (
		rwStr   = "a*|b"
		example = " "
		// folder  = "assets/test/contain/2"
	)
	var nka = expressions.NewRW(rwStr).ToENKA()
	// visualizer.MustVisualizeFSM(nka, folder, "v1.dot")
	var dfa = nka.ToDFA()
	// visualizer.MustVisualizeFSM(dfa, folder, "v2.dot")
	if !dfa.ContainString(example) {
		t.Fatalf("Получено Нет, а должно Да")
	}
}
