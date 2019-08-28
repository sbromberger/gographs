package graph

// HouseGraph returns an undirected house graph.
func HouseGraph() Graph {
	ss := []uint32{1, 1, 2, 3, 3, 4}
	sd := []uint32{2, 3, 4, 4, 5, 5}
	g, _ := New(ss, sd)
	return g
}
