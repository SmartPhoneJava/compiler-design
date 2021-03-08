package fsm

import (
	"gocompiler/internal/graph"
	"log"
)

// DR - КА с поддержкой операций d,r
/*
	A — конечный автомат,
	d(A) — детерминизированный автомат для A,
	r(A) — обратный автомат для A,
	dr(A) — результат d(r(A)). Аналогично для rdr(A) и drdr(A).
*/
type DR struct {
	FSM
}

// NewDR - создать новый экземпляр DR
func NewDR() *DR {
	return &DR{FSM{graph.NewGraph()}}
}

// NewDRFromFS - как NewDR, но на основе другого КА
func NewDRFromFS(fsm FSM) *DR {
	return &DR{fsm}
}

// NewDRFromEdges - как NewDR, но на основе списка ребёр
func NewDRFromEdges(
	edges []graph.Edge,
	first, last []string,
) *DR {
	var dr = NewDR()
	for _, edge := range edges {
		dr.AddEdge(&edge)
	}
	dr.SetFirstLast(first, last)
	return dr
}

// D Детерминизировать КА
func (A *DR) D() *DR {
	v := &DR{FSM{A.ToDka().Beautify()}}
	*A = *v
	return A
}

// R Построить обратный КА
func (A *DR) R() *DR {
	var reverseMe = NewDR()
	for _, edge := range A.Edges {
		reverseMe.AddEdge(&graph.Edge{
			From:   edge.To,
			To:     edge.From,
			Weight: edge.Weight,
		})
	}

	log.Println("reverse last", A.Last)
	log.Println("reverse first", A.First)

	reverseMe.SetFirstLast(A.Last, A.First)
	reverseMe = &DR{FSM{reverseMe.Beautify()}}
	*A = *reverseMe
	return A
}

// IsSame - сранить два КА
func (A DR) IsSame(B DR) bool {
	var try = 1
	for try > 0 {
		try--
		right := A.isSame(B)
		if right {
			return true
		}
	}
	return false
}

// IsSame - сранить два КА
func (A DR) isSame(B DR) bool {
	var (
		// Если не бьютифаить графы, то могут не
		// совпасть названия узлов
		copyA = A.Beautify()
		copyB = B.Beautify()
	)
	if len(copyA.Vertexes) != len(copyB.Vertexes) ||
		len(copyA.Edges) != len(copyB.Edges) {
		return false
	}
	// Проверяем равенство вершин и ребёр
	for _, vA := range copyA.Vertexes {
		vB, ok := copyB.Vertexes[vA.ID]
		if !ok || len(vA.In) != len(vB.In) || len(vA.Out) != len(vB.Out) {
			return false
		}
		for k := range vA.In {
			delete(vB.In, k)
		}
		if len(vB.In) != 0 {
			return false
		}
		for k := range vA.Out {
			delete(vB.Out, k)
		}
		if len(vB.Out) != 0 {
			return false
		}
		delete(copyB.Vertexes, vA.ID)
	}
	if len(copyB.Vertexes) != 0 {
		return false
	}

	// На самом деле дальнейшая проверка ребер не нужна,
	//  поскольку ребра итак проверены ниже, но раз уж
	//  это проверка на полное равенство структур, то пусть
	//  будет
	for _, vA := range copyA.Edges {
		_, ok := copyB.Edges[vA.ID()]
		if !ok {
			return false
		}
		delete(copyB.Edges, vA.ID())
	}
	if len(copyB.Edges) != 0 {
		return false
	}
	return true
}
