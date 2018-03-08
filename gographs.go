package gographs

import (
	"fmt"

	"github.com/sbromberger/gographs/sparsevecs"
)

type Edge struct {
	src uint32
	dst uint32
}

type GoGraph interface {
	IsDirected() bool
	HasEdge(u, v int32) bool
	Fadj() *sparsevecs.UInt32SparseVec
	Badj() *sparsevecs.UInt32SparseVec
	// AddEdge(u, v int) err
	InNeighbors(u int32) []int32
	OutNeighbors(u int32) []int32
}

// Graph is a graph structure
type Graph struct {
	mx sparsevecs.UInt32SparseVec
}

type DiGraph struct {
	fmx, bmx sparsevecs.UInt32SparseVec
}

func (g *Graph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

func (g *Graph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

func (g *DiGraph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.fmx
}

func (g *DiGraph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.bmx
}

func (g *Graph) IsDirected() bool {
	return false
}

func (g DiGraph) IsDirected() bool {
	return true
}

// MakeGraph creates an undirected graph from rowidx and colptr vectors.
func MakeGraph(r []uint32, c []uint64) Graph {
	return Graph{sparsevecs.UInt32SparseVec{r, c}}
}

// HasEdge returns true if an edge exists between u and v.
func (g *Graph) HasEdge(u, v int) bool {
	f := g.Fadj()
	return f.GetIndexInt(u, v) || f.GetIndexInt(v, u)
}

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

func (g DiGraph) Size() int {
	return len(g.fmx.Rowidx) / 2
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
