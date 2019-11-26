package graph

func (e SimpleEdge) Src() uint32 {
	return e.src
}

func (e SimpleEdge) Dst() uint32 {
	return e.dst
}

func (e SimpleEdge) Weight() float64 {
	return 1.0
}

// Reverse an edge
func Reverse(e SimpleEdge) Edge {
	return SimpleEdge{e.dst, e.src}
}
