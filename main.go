package main

import (
	"gocompiler/internal/converter"
	"gocompiler/internal/expressions"
	"gocompiler/internal/visualizer"
	"time"
)

const (
	Addr         = ":2997"
	ReadTimeout  = time.Hour * 3600
	WriteTimeout = time.Hour * 3600
	IdleTimeout  = time.Hour * 3600
)

func main() {
	var rw = expressions.NewRW("(xy* | ab | (x | a*)) (x | y*)")
	kda := converter.ExpressionToNKA(&rw)
	visualizer.MustVisualizeFSM(kda, "assets/v1.dot")

	kda.RemoveShortCircuits()
	visualizer.MustVisualizeFSM(kda, "assets/v2.dot")

	kda.ReplaceEpsilons()
	visualizer.MustVisualizeFSM(kda, "assets/v3.dot")

	kda.ReplaceEqualEdges()
	visualizer.MustVisualizeFSM(kda, "assets/v4.dot")

	kda.NkaToDka()
	visualizer.MustVisualizeFSM(kda, "assets/v5.dot")

	kda.ReplaceEqualEdges()
	visualizer.MustVisualizeFSM(kda, "assets/v6.dot")

	// for _, v := range kda.Vertexes {
	// 	log.Println("vertex is", v.ID)
	// 	log.Println("IN:", v.In)
	// 	log.Println("Out:", v.Out)
	// }

	// r := mux.NewRouter()
	// r.HandleFunc("/",
	// 	func(w http.ResponseWriter, r *http.Request) {

	// 		graph, err := graph.NewChartsGraph(*kda)
	// 		if err != nil {
	// 			log.Fatal("/", err)
	// 		}

	// 		if err = graph.Render(w); err != nil {
	// 			log.Println(err)
	// 		}
	// 	})

	// server := &http.Server{
	// 	Addr:           Addr,
	// 	Handler:        r,
	// 	ReadTimeout:    ReadTimeout,
	// 	WriteTimeout:   WriteTimeout,
	// 	IdleTimeout:    IdleTimeout,
	// 	MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	// }

	//log.Println("Server is on localhost" + Addr + "/")
	//server.ListenAndServe()
}

// func (kda *KDA) Create(rw RW) {
// 	var (
// 		specSymbols = "()* "
// 		operation = No
// 	)
// 	for _, r := range rw {
// 		if strings.ContainsRune(specSymbols, r) {
// 			continue
// 		}
// 		if operation {}
// 	}

// }

// xy* (x | y*) | ab (x | y*) | (x | a*) (x | y*)
// xy* | ab (x | y*) | (x | a*) (x | y*)
// (xy* | ab | (x | a*)) (x | y*)
