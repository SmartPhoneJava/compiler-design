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
	visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/example/1", "expected.dot")

	var real = fsm.DR{*origin.ToENKA()}

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/example/1", "real.dot")

	if !expected.IsSame(real) {
		t.Fatalf("Графы не сошлись, см. картинки в /assets/example/1")
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
	visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/example/2", "expected.dot")
	visualizer.MustVisualizeFSM(&origin.FSM, "assets/test/example/2", "origin.dot")
	real := fsm.DR{*origin.ToDFA()}

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/example/2", "real.dot")

	if !expected.IsSame(real) {
		t.Fatalf("Графы не сошлись, см. картинки в /assets/example/2")
	}
}
