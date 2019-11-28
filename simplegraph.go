package graph

import (
	"fmt"

	"github.com/sbromberger/graphmatrix"
)

// SimpleGraph is a graph structure representing an undirected graph.
type SimpleGraph struct {
	fmx, bmx graphmatrix.GraphMatrix
}

func (g SimpleGraph) String() string {
	return fmt.Sprintf("(%d, %d) graph", g.NumVertices(), g.NumEdges())
}

// MakeSimpleGraph creates an undirected graph from vectors of source and dest vertices.
func New(ss, ds []uint32) (SimpleGraph, error) {
	rSs := make([]uint32, len(ss))
	rDs := make([]uint32, len(ds))
	copy(rSs, ss)
	copy(rDs, ds)

	err := graphmatrix.SortIJ(&ss, &ds)
	if err != nil {
		return SimpleGraph{}, err
	}
	err = graphmatrix.SortIJ(&rDs, &rSs)
	if err != nil {
		return SimpleGraph{}, err
	}
	fmx, err := graphmatrix.NewFromSortedIJ(ss, ds)
	if err != nil {
		return SimpleGraph{}, err
	}
	bmx, err := graphmatrix.NewFromSortedIJ(rDs, rSs)
	if err != nil {
		return SimpleGraph{}, err
	}
	return SimpleGraph{fmx: fmx, bmx: bmx}, nil
}

func NewFromEdgeList(l EdgeList) (SimpleGraph, error) {
	ss := make([]uint32, len(l))
	ds := make([]uint32, len(l))
	for i, e := range l {
		ss[i] = e.Src()
		ds[i] = e.Dst()
	}
	return New(ss, ds)
}

// OutDegree returns the out degree of vertex u.
func (g SimpleGraph) OutDegree(u uint32) uint32 { return uint32(len(g.OutNeighbors(u))) }

// InDegree returns the indegree of vertex u.
func (g SimpleGraph) InDegree(u uint32) uint32 { return uint32(len(g.InNeighbors(u))) }

// OutNeighbors returns the out neighbors of vertex u.
func (g SimpleGraph) OutNeighbors(u uint32) []uint32 {
	r, _ := g.fmx.GetRow(u)
	return r
}

// InNeighbors returns the in neighbors of vertex u.
func (g SimpleGraph) InNeighbors(u uint32) []uint32 {
	r, _ := g.bmx.GetRow(u)
	return r
}

// HasEdge returns true if an edge exists between u and v.
func (g SimpleGraph) HasEdge(u, v uint32) bool {
	un := g.OutNeighbors(u)
	vn := g.InNeighbors(v)
	lenun := uint64(len(un))
	lenvn := uint64(len(vn))
	var found bool
	if lenvn > lenun {
		_, found = graphmatrix.SearchSorted32(un, v, 0, lenun)
	} else {
		_, found = graphmatrix.SearchSorted32(vn, u, 0, lenvn)
	}
	return found
}

// AddEdge adds an edge to graph g
func (g *SimpleGraph) AddEdge(u, v uint32) error {
	if err := g.fmx.SetIndex(u, v); err != nil {
		return err
	}
	if err := g.bmx.SetIndex(v, u); err != nil {
		return err
	}
	return nil
}

// NumEdges returns the number of edges
func (g SimpleGraph) NumEdges() uint64 {
	return g.fmx.N()
}

// NumVertices returns the number of vertices
func (g SimpleGraph) NumVertices() uint32 {
	return g.fmx.Dim()
}

// Edges returns an iterator of edges
func (g SimpleGraph) Edges() EdgeIter {
	return &SimpleEdgeIter{mxiter: g.fmx.NewNZIter()}
}

// FMat returns the forward matrix of the graph.
func (g SimpleGraph) FMat() graphmatrix.GraphMatrix {
	return g.fmx
}

// BMat returns the backward matrix of the graph.
func (g SimpleGraph) BMat() graphmatrix.GraphMatrix {
	return g.bmx
}

func (g SimpleGraph) IsDirected() bool { return false }

func FromRaw(findptr []uint64, find []uint32, bindptr []uint64, bind []uint32) (SimpleGraph, error) {
	fmx := graphmatrix.GraphMatrix{IndPtr: findptr, Indices: find}
	bmx := graphmatrix.GraphMatrix{IndPtr: bindptr, Indices: bind}
	return SimpleGraph{fmx, bmx}, nil
}
