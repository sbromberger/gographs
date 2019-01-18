package gographs

import "github.com/sbromberger/gographs/sparsevecs"

// HouseGraph returns an undirected house graph.
func HouseGraph() Graph {
	fVec := []uint32{1, 2, 0, 3, 0, 3, 4, 1, 2, 4, 2, 3}
	fInt := []uint64{0, 2, 4, 7, 10, 12}
	mx := sparsevecs.UInt32SparseVec{Rowidx: fVec, Colptr: fInt}
	return Graph{mx}
}
