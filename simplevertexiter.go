package graph

type SimpleVertexIter struct {
	curr uint32
	max  uint32
}

func (it *SimpleVertexIter) Next() uint32 {
	it.curr++
	if it.curr > it.max {
		it.curr = it.max
	}
	return it.curr
}

func (it *SimpleVertexIter) Done() bool {
	return it.curr > it.max
}
