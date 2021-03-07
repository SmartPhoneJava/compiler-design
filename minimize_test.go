package main

import (
	"gocompiler/internal/fsm"
	"gocompiler/internal/graph"
	"gocompiler/internal/visualizer"
	"testing"
)

// Пример из http://neerc.ifmo.ru/wiki/index.php?title=Алгоритм_Бржозовского
func TestMinimize1(t *testing.T) {
	return
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "0",
			To:     "1",
			Weight: "a",
		},
		{
			From:   "0",
			To:     "2",
			Weight: "a",
		},
		{
			From:   "0",
			To:     "2",
			Weight: "b",
		},
		{
			From:   "1",
			To:     "2",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "1",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "2",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "3",
			Weight: "b",
		},
		{
			From:   "1",
			To:     "3",
			Weight: "b",
		},
	})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "0",
			To:     "1",
			Weight: "a",
		},
		{
			From:   "0",
			To:     "1",
			Weight: "b",
		},
		{
			From:   "1",
			To:     "1",
			Weight: "a",
		},
		{
			From:   "1",
			To:     "2",
			Weight: "b",
		},
	})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real0.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real1.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real2.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real3.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real4.dot")

	//real.Chain(real.D, real.R, real.D, real.R)

	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real.dot")
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/expected.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}

// Пример из http://neerc.ifmo.ru/wiki/index.php?title=Минимизация_ДКА,_алгоритм_за_O(n%5E2)_с_построением_пар_различимых_состояний
func TestMinimize2(t *testing.T) {
	return
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "А",
			To:     "B",
			Weight: "1",
		},
		{
			From:   "B",
			To:     "A",
			Weight: "1",
		},
		{
			From:   "B",
			To:     "H",
			Weight: "0",
		},
		{
			From:   "A",
			To:     "H",
			Weight: "0",
		},
		{
			From:   "H",
			To:     "C",
			Weight: "0",
		},
		{
			From:   "H",
			To:     "C",
			Weight: "1",
		},
		{
			From:   "C",
			To:     "E",
			Weight: "0",
		},
		{
			From:   "C",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "E",
			To:     "F",
			Weight: "0",
		},
		{
			From:   "D",
			To:     "E",
			Weight: "0",
		},
		{
			From:   "D",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "E",
			To:     "G",
			Weight: "1",
		},
		{
			From:   "G",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "G",
			To:     "G",
			Weight: "0",
		},
		{
			From:   "F",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "F",
			To:     "F",
			Weight: "0",
		},
	})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "A",
			Weight: "1",
		},
		{
			From:   "A",
			To:     "H",
			Weight: "0",
		},
		{
			From:   "H",
			To:     "C",
			Weight: "0",
		},
		{
			From:   "H",
			To:     "C",
			Weight: "1",
		},
		{
			From:   "C",
			To:     "E",
			Weight: "0",
		},
		{
			From:   "C",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "E",
			To:     "F",
			Weight: "1",
		},
		{
			From:   "E",
			To:     "F",
			Weight: "0",
		},
		{
			From:   "F",
			To:     "F",
			Weight: "0",
		},
		{
			From:   "F",
			To:     "F",
			Weight: "1",
		},
	})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real0_2.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real1_2.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real2_2.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real3_2.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real4_2.dot")

	//real.Chain(real.D, real.R, real.D, real.R)

	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real_2.dot")
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/expected_2.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}

// https://lektsii.org/6-91118.html
func TestMinimize3(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "B",
			Weight: "0",
		},
		{
			From:   "A",
			To:     "C",
			Weight: "1",
		},
		{
			From:   "C",
			To:     "D",
			Weight: "0",
		},
		{
			From:   "E",
			To:     "B",
			Weight: "0",
		},
		{
			From:   "C",
			To:     "E",
			Weight: "1",
		},
		{
			From:   "E",
			To:     "D",
			Weight: "1",
		},
		{
			From:   "D",
			To:     "E",
			Weight: "1",
		},
	})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "B",
			Weight: "1",
		},
		{
			From:   "A",
			To:     "B",
			Weight: "0",
		},
		{
			From:   "D",
			To:     "B",
			Weight: "1",
		},
		{
			From:   "B",
			To:     "D",
			Weight: "1",
		},
		{
			From:   "D",
			To:     "D",
			Weight: "1",
		},
	})
	var real = *origin

	//visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real0.dot")
	real.R()
	//visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real1.dot")
	real.D()
	//visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real2.dot")
	real.R()
	//visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real3.dot")
	real.D()
	//visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real4.dot")

	//real.Chain(real.D, real.R, real.D, real.R)

	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real.dot")
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/expected.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}

// https://github.com/navin-mohan/dfa-minimization/blob/master/DFA%20minimization.ipynb
func TestMinimize4(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "3",
			Weight: "a",
		},
		{
			From:   "3",
			To:     "5",
			Weight: "a",
		},
		{
			From:   "5",
			To:     "3",
			Weight: "a",
		},
		{
			From:   "3",
			To:     "4",
			Weight: "b",
		},
		{
			From:   "4",
			To:     "4",
			Weight: "a",
		},
		{
			From:   "4",
			To:     "4",
			Weight: "b",
		},
		{
			From:   "5",
			To:     "2",
			Weight: "b",
		},
		{
			From:   "2",
			To:     "4",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "1",
			Weight: "b",
		},
		{
			From:   "1",
			To:     "2",
			Weight: "b",
		},
	})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "1",
			To:     "2",
			Weight: "a",
		},
		{
			From:   "1",
			To:     "3",
			Weight: "b",
		},
		{
			From:   "3",
			To:     "1",
			Weight: "b",
		},
		{
			From:   "3",
			To:     "4",
			Weight: "a",
		},
		{
			From:   "4",
			To:     "3",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "2",
			Weight: "a",
		},
		{
			From:   "2",
			To:     "2",
			Weight: "b",
		},
		{
			From:   "4",
			To:     "2",
			Weight: "b",
		},
	})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real0.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real1.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real2.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real3.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real4.dot")

	//real.Chain(real.D, real.R, real.D, real.R)

	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real.dot")
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/expected.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}
