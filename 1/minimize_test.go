package main

import (
	"gocompiler/internal/fsm"
	"gocompiler/internal/graph"
	"gocompiler/internal/visualizer"
	"testing"
)

// Пример из http://neerc.ifmo.ru/wiki/index.php?title=Алгоритм_Бржозовского
func TestMinimize1(t *testing.T) {
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
	}, []string{"0"}, []string{"3"})
	visualizer.MustVisualizeFSM(&origin.FSM, "assets/test/min/1", "origin.dot")
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
	}, []string{"0"}, []string{"2"})
	origin.R().D().R().D()

	visualizer.MustVisualizeFSM(&origin.FSM, "assets/test/min/1", "real.dot")

	if !expected.IsSame(*origin) {
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/1", "expected.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}

// Пример из http://neerc.ifmo.ru/wiki/index.php?title=Минимизация_ДКА,_алгоритм_за_O(n%5E2)_с_построением_пар_различимых_состояний
func TestMinimize2(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
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
	}, []string{"A"}, []string{"G", "F"})
	visualizer.MustVisualizeFSM(&origin.FSM, "assets/test/min/2", "origin.dot")
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
	}, []string{"A"}, []string{"F"})
	var real = *origin

	real.R().D().R().D()

	visualizer.MustVisualizeDR(real.CompareMode(), "assets/test/min/2", "real.dot")
	visualizer.MustVisualizeDR(expected.CompareMode(), "assets/test/min/2", "expected.dot")

	if !expected.IsSame(real) {
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}

// https://lektsii.org/6-91118.html
// второй пример
/*
func TestMinimize3(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "B",
			Weight: "0",
		},
		{
			From:   "B",
			To:     "D",
			Weight: "1",
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
		{
			From:   "G",
			To:     "F",
			Weight: "0",
		},
		{
			From:   "F",
			To:     "G",
			Weight: "0",
		},
		{
			From:   "F",
			To:     "E",
			Weight: "1",
		},
		{
			From:   "G",
			To:     "D",
			Weight: "1",
		},
	}, []string{"A"}, []string{"D"})
	log.Println("test3")
	visualizer.MustVisualizeFSM(&origin.FSM, "assets/test/min/3", "origin.dot")
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
	}, []string{"A"}, []string{"D"})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real0.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real1.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real2.dot")
	real.R()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real3.dot")
	real.D()
	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real4.dot")

	//real.Chain(real.D, real.R, real.D, real.R)

	if !expected.IsSame(real) {
		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real.dot")
		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/3", "expected.dot")
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}
*/

// https://github.com/navin-mohan/dfa-minimization/blob/master/DFA%20minimization.ipynb
/*
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
	}, []string{"1"}, []string{"5"})
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
	}, []string{"3"}, []string{"2"})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "origin.dot")
	real.R().D().R().D()

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/3", "real.dot")
	visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/3", "expected.dot")

	if !expected.IsSame(real) {
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}
*/
/*
// // http://esyr.org/wiki/Конструирование_Компиляторов%2C_Алгоритмы_решения_задач#.D0.9F.D0.BE.D1.81.D1.82.D1.80.D0.BE.D0.B5.D0.BD.D0.B8.D0.B5_.D0.94.D0.9A.D0.90_.D0.BF.D0.BE_.D0.9D.D0.9A.D0.90
func TestMinimize5(t *testing.T) {
	var origin = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "A",
			To:     "C",
			Weight: "c",
		},
		{
			From:   "C",
			To:     "E",
			Weight: "c",
		},
		{
			From:   "E",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "B",
			To:     "E",
			Weight: "c",
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
			To:     "C",
			Weight: "c",
		},
	}, []string{"A"}, []string{"C"})
	var expected = fsm.NewDRFromEdges([]graph.Edge{
		{
			From:   "A",
			To:     "B",
			Weight: "a",
		},
		{
			From:   "B",
			To:     "A",
			Weight: "b",
		},
		{
			From:   "B",
			To:     "A",
			Weight: "c",
		},
		{
			From:   "A",
			To:     "C",
			Weight: "c",
		},
	}, []string{"B"}, []string{"С"})
	var real = *origin

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/5", "origin.dot")
	real.R().D().R().D()

	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/5", "real.dot")
	visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/5", "expected.dot")

	if !expected.IsSame(real) {
		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
	}
}
*/

// // https://intuit.ru/studies/courses/26/26/lecture/801?page=4
// func TestMinimize6(t *testing.T) {
// 	return
// 	var origin = fsm.NewDRFromEdges([]graph.Edge{
// 		{
// 			From:   "A",
// 			To:     "C",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "C",
// 			To:     "C",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "C",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "A",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "B",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "B",
// 			To:     "D",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "D",
// 			To:     "B",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "D",
// 			To:     "E",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "E",
// 			To:     "C",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "E",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		// {
// 		// 	From:   "F",
// 		// 	To:     "E",
// 		// 	Weight: "a",
// 		// },
// 	})
// 	origin.SetFirstLast([]string{"A"}, []string{"E"})
// 	var expected = fsm.NewDRFromEdges([]graph.Edge{
// 		{
// 			From:   "C",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "C",
// 			To:     "C",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "B",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "B",
// 			To:     "D",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "D",
// 			To:     "B",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "D",
// 			To:     "E",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "E",
// 			To:     "B",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "E",
// 			To:     "C",
// 			Weight: "b",
// 		},
// 	})
// 	var real = *origin

// 	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real0.dot")
// 	real.R()
// 	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real1.dot")
// 	real.D()
// 	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real2.dot")
// 	real.R()
// 	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real3.dot")
// 	real.D()
// 	visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real4.dot")

// 	//real.Chain(real.D, real.R, real.D, real.R)

// 	if !expected.IsSame(real) {
// 		visualizer.MustVisualizeFSM(&real.FSM, "assets/test/min/real.dot")
// 		visualizer.MustVisualizeFSM(&expected.FSM, "assets/test/min/expected.dot")
// 		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
// 	}
// }
