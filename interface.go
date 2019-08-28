package graph

import (
	"github.com/sbromberger/graphmatrix"
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

type EdgeIter struct {
	mxiter *graphmatrix.NZIter
}

func (it EdgeIter) Next() (e Edge, done bool) {
	r, c, d := it.mxiter.Next()
	return Edge{Src: r, Dst: c}, d
}

// Graph is a graph structure representing an undirected graph.
type Graph struct {
	fmx, bmx graphmatrix.GraphMatrix
}

// MakeGraph creates an undirected graph from vectors of source and dest vertices.
func New(ss, ds []uint32) (Graph, error) {
	r_ss := make([]uint32, len(ss))
	r_ds := make([]uint32, len(ds))
	copy(r_ss, ss)
	copy(r_ds, ds)

	graphmatrix.SortIJ(&ss, &ds)
	graphmatrix.SortIJ(&r_ds, &r_ss)
	fmx, err := graphmatrix.NewFromSortedIJ(ss, ds)
	if err != nil {
		return Graph{}, err
	}
	bmx, err := graphmatrix.NewFromSortedIJ(r_ds, r_ss)
	if err != nil {
		return Graph{}, err
	}
	return Graph{fmx: fmx, bmx: bmx}, nil
}

func NewFromEdgeList(l EdgeList) (Graph, error) {
	ss := make([]uint32, len(l))
	ds := make([]uint32, len(l))
	for i, e := range l {
		ss[i] = e.Src
		ds[i] = e.Dst
	}
	return New(ss, ds)
}

// OutDegree returns the out degree of vertex u.
func (g Graph) OutDegree(u uint32) uint32 { return uint32(len(g.OutNeighbors(u))) }

// InDegree returns the indegree of vertex u.
func (g Graph) InDegree(u uint32) uint32 { return uint32(len(g.InNeighbors(u))) }

// OutNeighbors returns the out neighbors of vertex u.
func (g Graph) OutNeighbors(u uint32) []uint32 { return g.fmx.GetRow(u) }

// InNeighbors returns the in neighbors of vertex u.
func (g Graph) InNeighbors(u uint32) []uint32 { return g.bmx.GetRow(u) }

// HasEdge returns true if an edge exists between u and v.
func (g Graph) HasEdge(u, v uint32) bool {
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
func (g *Graph) AddEdge(u, v uint32) error {
	if err := g.fmx.SetIndex(u, v); err != nil {
		return err
	}
	if err := g.bmx.SetIndex(v, u); err != nil {
		return err
	}
	return nil
}

// Size returns the number of edges
func (g Graph) NumEdges() uint64 {
	return g.fmx.N()
}

// Order returns the number of vertices
func (g Graph) NumVertices() uint32 {
	return g.fmx.Dim()
}

// Edges returns an iterator of edges
func (g Graph) Edges() EdgeIter {
	return EdgeIter{mxiter: g.fmx.NewNZIter()}
}

// FMat returns the forward matrix of the graph.
func (g Graph) FMat() graphmatrix.GraphMatrix {
	return g.fmx
}

// BMat returns the backward matrix of the graph.
func (g Graph) BMat() graphmatrix.GraphMatrix {
	return g.bmx
}

func FromRaw(findptr []uint64, find []uint32, bindptr []uint64, bind []uint32) (Graph, error) {
	fmx := graphmatrix.GraphMatrix{IndPtr: findptr, Indices: find}
	bmx := graphmatrix.GraphMatrix{IndPtr: bindptr, Indices: bind}
	return Graph{fmx, bmx}, nil
}
