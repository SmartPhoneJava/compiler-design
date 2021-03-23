package main

import (
	"gocompiler/internal/expressions"
	"gocompiler/internal/fsm"
	"gocompiler/internal/graph"
	"gocompiler/internal/visualizer"
	"testing"
)

// Пример 3.24
func TestExample3_24(t *testing.T) {
	var origin expressions.RW = "(a|b)*a"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "0",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "0",
			To:     "1",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "2",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "4",
			Weight: "e",
		},
		{
			From:   "2",
			To:     "3",
			Weight: "a",
		},
		{
			From:   "4",
			To:     "5",
			Weight: "b",
		},
		{
			From:   "3",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "5",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "1",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "7",
			To:     "8",
			Weight: "a",
		},
	}, []string{"0"}, []string{"7"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/1"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

// Пример 3.24
func TestExample3_21(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "0",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "0",
			To:     "1",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "2",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "4",
			Weight: "e",
		},
		{
			From:   "2",
			To:     "3",
			Weight: "a",
		},
		{
			From:   "4",
			To:     "5",
			Weight: "b",
		},
		{
			From:   "3",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "5",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "1",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "7",
			To:     "8",
			Weight: "a",
		},
		{
			From:   "8",
			To:     "9",
			Weight: "b",
		},
		{
			From:   "9",
			To:     "10",
			Weight: "b",
		},
	}, []string{"0"}, []string{"10"})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "C",
			Weight: "b",
		},
		{
			From:   "C",
			To:     "C",
			Weight: "b",
		},
		{
			From:   "A",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "B",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "C",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "E",
			To:     "C",
			Weight: "b",
		},
		{
			From:   "B",
			To:     "D",
			Weight: "b",
		},
		{
			From:   "D",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "D",
			To:     "E",
			Weight: "b",
		},
		{
			From:   "E",
			To:     "B",
			Weight: "a",
		},
	}, []string{"A"}, []string{"E"})
	var (
		real   = fsm.DR{*origin.ToDFA()}
		folder = "assets/test/example/2"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

// Пример мой
func TestExampleMy1(t *testing.T) {
	var origin expressions.RW = "a*|b"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "e",
		},

		{
			From:   "2",
			To:     "3",
			Weight: "b",
		},
		{
			From:   "3",
			To:     "4",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "5",
			Weight: "e",
		},
		{
			From:   "5",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "5",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "7",
			To:     "8",
			Weight: "a",
		},
		{
			From:   "8",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "8",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "4",
			Weight: "e",
		},
	}, []string{"1"}, []string{"4"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/3"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy2(t *testing.T) {
	var origin expressions.RW = "((((ab))))"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "a",
		},

		{
			From:   "2",
			To:     "3",
			Weight: "b",
		},
	}, []string{"1"}, []string{"3"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/4"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy3(t *testing.T) {
	var origin expressions.RW = " ((a|b)*b)"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "e",
		},
		{
			From:   "1",
			To:     "3",
			Weight: "e",
		},
		{
			From:   "3",
			To:     "4",
			Weight: "e",
		},
		{
			From:   "3",
			To:     "5",
			Weight: "e",
		},
		{
			From:   "4",
			To:     "6",
			Weight: "a",
		},
		{
			From:   "5",
			To:     "7",
			Weight: "b",
		},
		{
			From:   "6",
			To:     "8",
			Weight: "e",
		},
		{
			From:   "7",
			To:     "8",
			Weight: "e",
		},
		{
			From:   "8",
			To:     "2",
			Weight: "e",
		},
		{
			From:   "8",
			To:     "3",
			Weight: "e",
		},
		{
			From:   "2",
			To:     "9",
			Weight: "b",
		},
	}, []string{"1"}, []string{"9"})

	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/5"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy4(t *testing.T) {
	var origin expressions.RW = "((a|b)b*)"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "e",
		},

		{
			From:   "1",
			To:     "3",
			Weight: "e",
		},
		{
			From:   "2",
			To:     "4",
			Weight: "a",
		},
		{
			From:   "3",
			To:     "5",
			Weight: "b",
		},
		{
			From:   "4",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "5",
			To:     "6",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "7",
			To:     "8",
			Weight: "b",
		},
		{
			From:   "8",
			To:     "7",
			Weight: "e",
		},
		{
			From:   "6",
			To:     "9",
			Weight: "e",
		},
		{
			From:   "8",
			To:     "9",
			Weight: "e",
		},
	}, []string{"1"}, []string{"9"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/6"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy5(t *testing.T) {
	var origin expressions.RW = "((a|b)(b))"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "b",
		},

		{
			From:   "1",
			To:     "1",
			Weight: "a",
		},
	}, []string{"1"}, []string{"2"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/7"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy6(t *testing.T) {
	var origin expressions.RW = " ((a|b)(b)*)"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "b",
		},

		{
			From:   "1",
			To:     "1",
			Weight: "a",
		},
	}, []string{"1"}, []string{"2"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/8"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}

func TestExampleMy7(t *testing.T) {
	var origin expressions.RW = "((((a))*))"
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "b",
		},

		{
			From:   "1",
			To:     "1",
			Weight: "a",
		},
	}, []string{"1"}, []string{"2"})
	var (
		real   = fsm.DR{*origin.ToENKA()}
		folder = "assets/test/example/9"
	)
	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&expected.FSM, folder, "expected.dot")
		visualizer.MustVisualizeFSM(&real.FSM, folder, "real.dot")
		t.Fatalf("Графы не сошлись, см. картинки в %s", folder)
	}
}
