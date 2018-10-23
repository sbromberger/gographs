package gographs

import (
	"fmt"

	"github.com/sbromberger/gographs/sparsevecs"
)

// Edge is a graph edge.
type Edge struct {
	Src uint32
	Dst uint32
}

// Reverse an edge
func Reverse(e Edge) Edge {
	return Edge{e.Dst, e.Src}
}

// EdgeList is a slice of edges
type EdgeList []Edge

func (e EdgeList) Len() int {
	return len(e)
}

func (e EdgeList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e EdgeList) Less(i, j int) bool {
	if e[i].Src < e[j].Src {
		return true
	}
	if e[i].Src > e[j].Src {
		return false
	}
	return e[i].Dst < e[j].Dst
}

// GoGraph is an interface to Graphs and DiGraphs.
type GoGraph interface {
	IsDirected() bool
	HasEdge(u, v uint32) bool
	Fadj() *sparsevecs.UInt32SparseVec
	Badj() *sparsevecs.UInt32SparseVec
	// AddEdge(u, v int) err
	// InNeighbors(u uint32) []uint32
	OutNeighbors(u uint32) []uint32
	Order() int
	Size() int
}

// Graph is a graph structure
type Graph struct {
	mx sparsevecs.UInt32SparseVec
}

// DiGraph is a directed Graph.
type DiGraph struct {
	fmx, bmx sparsevecs.UInt32SparseVec
}

// Fadj returns a pointer to the forward adjacencies.
func (g *Graph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

// Badj returns a ptr to the backward adjacencies.
func (g *Graph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

// Fadj returns a pointer to the forward adjacencies.
func (g *DiGraph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.fmx
}

// Badj returns a ptr to the backward adjacencies.
func (g *DiGraph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.bmx
}

// IsDirected is true if the graph is directed.
func (g *Graph) IsDirected() bool {
	return false
}

// IsDirected is true if the graph is directed.
func (g DiGraph) IsDirected() bool {
	return true
}

// MakeGraph creates an undirected graph from rowidx and colptr vectors.
func MakeGraph(r []uint32, c []uint64) Graph {
	return Graph{sparsevecs.UInt32SparseVec{r, c}}
}

// HasEdge returns true if an edge exists between u and v.
func (g *Graph) HasEdge(u, v uint32) bool {
	f := g.Fadj()
	return f.GetIndex(u, v) || f.GetIndex(v, u)
}

// HasEdge returns true if an edge exists between u and v.
func (g *DiGraph) HasEdge(u, v int) bool {
	f := g.Fadj()
	b := g.Badj()
	return f.GetIndexInt(u, v) || b.GetIndexInt(v, u)
}

// AddEdge Adds an edge to graph g
// func (g *Graph) AddEdge(u, v uint32) error {
// 	return g.mx.Insert(u, v)
// }

// Size returns the number of edges
func (g *Graph) Size() int {
	return len(g.mx.Rowidx) / 2
}

// Size returns the number of edges
func (g DiGraph) Size() int {
	return len(g.fmx.Rowidx)
}

// Order returns the number of vertices
func (g *Graph) Order() int {
	return len(g.mx.Colptr) - 1
}

// Order returns the number of vertices
func (g *DiGraph) Order() int {
	return len(g.fmx.Colptr) - 1
}

// Edges returns a list of edges
func (g *Graph) Edges() []Edge {
	ne := g.Size()
	fmt.Println("ne = ", ne)
	edgelist := make([]Edge, ne)
	edgesAdded := 0
	for c := range g.mx.Colptr[:len(g.mx.Colptr)-1] {
		uc := uint32(c)
		rrange := g.mx.GetRange(uc)
		for _, r := range rrange {
			if r <= uc {
				fmt.Printf("Adding edge %d, %d\n", r, uc)
				edgelist[edgesAdded] = Edge{r, uc}
				edgesAdded++
			}
		}
	}
	return edgelist
}

// AddEdgeInt takes ints and calls AddEdge
// func (g *Graph) AddEdgeInt(u, v int) error {
// 	return g.AddEdge(uint32(u), uint32(v))

// }

// OutNeighbors returns the neighbors of vertex v.
func (g *Graph) OutNeighbors(v uint32) []uint32 {
	return g.mx.GetRange(v)
}

// OutNeighborsInt returns the neighbors of vertex v.
func (g *Graph) OutNeighborsInt(v int) []uint32 {
	return g.mx.GetRangeInt(v)
}
