package main

import (
	"gocompiler/internal/converter"
	"gocompiler/internal/expressions"
	"gocompiler/internal/visualizer"
)

func main() {
	var rw = expressions.NewRW("(xy* | ab | (x | a*)) (x | y*)")
	kda := converter.ExpressionToNKA(&rw)
	visualizer.MustVisualizeFSM(kda, "assets", "v1.dot")

	kda.RemoveShortCircuits()
	visualizer.MustVisualizeFSM(kda, "assets", "v2.dot")

	kda.ReplaceEpsilons()
	visualizer.MustVisualizeFSM(kda, "assets", "v3.dot")

	kda.ReplaceEqualEdges()
	visualizer.MustVisualizeFSM(kda, "assets", "v4.dot")

	kda.ToDka()
	visualizer.MustVisualizeFSM(kda, "assets", "v5.dot")

	kda.ReplaceEqualEdges()
	visualizer.MustVisualizeFSM(kda, "assets", "v6.dot")

}

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
