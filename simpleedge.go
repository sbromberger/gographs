package graph

// SimpleEdge is a graph edge.
type SimpleEdge struct {
	src uint32
	dst uint32
}

func (e SimpleEdge) Src() uint32 {
	return e.src
}

func (e SimpleEdge) Dst() uint32 {
	return e.dst
}

func (e SimpleEdge) Weight() float64 {
	return 1.0
}

// Reverse reverses a SimpleEdge.
func Reverse(e SimpleEdge) Edge {
	return SimpleEdge{e.dst, e.src}
}
