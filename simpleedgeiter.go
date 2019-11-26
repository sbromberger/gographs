package graph

func (it EdgeIter) Next() (e SimpleEdge, done bool) {
	r, c, d := it.mxiter.Next()
	return SimpleEdge{src: r, dst: c}, d
}
