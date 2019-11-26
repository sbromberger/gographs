package graph

// This file contains interface definitions for Graph and Edge types.

import (
	"github.com/sbromberger/graphmatrix"
)

type Edge interface {
	// Src returns the source index of the edge
	Src() uint32
	// Dst returns the destination index of the edge
	Dst() uint32
	// Weight returns the weight of the edge
	Weight() float64
}

// EdgeList is a slice of edges
type EdgeList []Edge

// Len returns the length of an edgelist.
func (e EdgeList) Len() int {
	return len(e)
}

func (e EdgeList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e EdgeList) Less(i, j int) bool {
	if e[i].Src() < e[j].Src() {
		return true
	}
	if e[i].Src() > e[j].Src() {
		return false
	}
	return e[i].Dst() < e[j].Dst()
}

type EdgeIter struct {
	mxiter *graphmatrix.NZIter
}

type Graph interface {
	// OutDegree returns the outdegree of vertex u.
	OutDegree(u uint32) uint32
	// InDegree returns the indegree of vertex u.
	InDegree(u uint32) uint32
	// OutNeighbors returns the out neighbors of vertex u.
	OutNeighbors(u uint32) []uint32
	// InNeighbors returns the in neighbors of vertex u.
	InNeighbors(u uint32) []uint32
	// HasEdge returns true if an edge exists between u and v.
	HasEdge(u, v uint32) bool
	// IsDirected is true if the graph is directed
	IsDirected() bool
	// NumEdges returns the number of edges in a graph.
	NumEdges() uint64
	// NumVertices returns the number of vertices in a graph.
	NumVertices() uint32
}
