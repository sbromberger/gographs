package graph

import "fmt"

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

func (e SimpleEdge) String() string {
	return fmt.Sprintf("SimpleEdge %d -> %d", e.Src(), e.Dst())
}

// Reverse reverses a SimpleEdge.
func Reverse(e SimpleEdge) Edge {
	return SimpleEdge{e.dst, e.src}
}
