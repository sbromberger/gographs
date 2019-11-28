package graph

import "github.com/sbromberger/graphmatrix"

type SimpleEdgeIter struct {
	mxiter graphmatrix.NZIter
}

func (it *SimpleEdgeIter) Next() Edge {
	r, c, _ := it.mxiter.Next()
	return SimpleEdge{src: r, dst: c}
}

func (it *SimpleEdgeIter) Done() bool {
	return it.mxiter.Done()
}
