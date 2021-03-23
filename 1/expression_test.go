package main

// import (
// 	"gocompiler/internal/expressions"
// 	"gocompiler/internal/fsm"
// 	"gocompiler/internal/graph"
// 	"gocompiler/internal/visualizer"
// 	"testing"
// )

// // Пример из https://habr.com/ru/post/166777/
// func TestExpression1(t *testing.T) {
// 	var (
// 		rw     = expressions.NewRW("(xy* | ab | (x | a*)) (x | y*)")
// 		kda    = rw.ToENKA()
// 		folder = "assets/test/expressions/2"
// 	)

// 	visualizer.MustVisualizeFSM(kda, folder, "v1.dot")

// 	kda.ToDka()
// 	visualizer.MustVisualizeFSM(kda, folder, "v2.dot")

// 	var expected = fsm.NewDRFromEdges([]graph.Edge{
// 		{
// 			From:   "p0",
// 			To:     "p1",
// 			Weight: "x",
// 		},
// 		{
// 			From:   "p1",
// 			To:     "p1",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p1",
// 			To:     "p4",
// 			Weight: "x",
// 		},
// 		{
// 			From:   "p0",
// 			To:     "p2",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p2",
// 			To:     "p2",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p0",
// 			To:     "p3",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "p3",
// 			To:     "p2",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p3",
// 			To:     "p5",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "p3",
// 			To:     "p6",
// 			Weight: "b",
// 		},
// 		{
// 			From:   "p5",
// 			To:     "p5",
// 			Weight: "a",
// 		},
// 		{
// 			From:   "p5",
// 			To:     "p2",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p5",
// 			To:     "p4",
// 			Weight: "x",
// 		},
// 		{
// 			From:   "p6",
// 			To:     "p2",
// 			Weight: "y",
// 		},
// 		{
// 			From:   "p6",
// 			To:     "p4",
// 			Weight: "x",
// 		},
// 		{
// 			From:   "p3",
// 			To:     "p4",
// 			Weight: "x",
// 		},
// 	}, []string{"p0"}, []string{"p4"})

// 	origin := fsm.NewDRFromFS(*kda)

// 	visualizer.MustVisualizeDR(origin.CompareMode(), folder, "real.dot")
// 	visualizer.MustVisualizeDR(expected.CompareMode(), folder, "expected.dot")

// 	if !expected.IsSame(*origin) {
// 		t.Fatalf("Графы не сошлись, см. картинки в /assets/test")
// 	}
// }
