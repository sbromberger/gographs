package graph

// HouseGraph returns an undirected house graph.
func HouseGraph() SimpleGraph {
	ss := []uint32{1, 1, 2, 3, 3, 4}
	sd := []uint32{2, 3, 4, 4, 5, 5}
	g, _ := New(ss, sd)
	return g
}

func PathGraph(n int) (SimpleGraph, error) {
	ss := make([]uint32, n-1)
	sd := make([]uint32, n-1)
	for i := uint32(0); i < uint32(n-1); i++ {
		ss[i] = i
		sd[i] = i + 1
	}
	return New(ss, sd)
}

func CycleGraph(n int) (SimpleGraph, error) {
	ss := make([]uint32, n)
	sd := make([]uint32, n)
	for i := uint32(0); i < uint32(n-1); i++ {
		ss[i] = i
		sd[i] = i + 1
	}
	ss[n-1] = uint32(n - 1)
	sd[n-1] = 0

	return New(ss, sd)
}

func WheelGraph(n int) (SimpleGraph, error) {
	ss := make([]uint32, 2*(n-1))
	sd := make([]uint32, 2*(n-1))
	n32 := uint32(n)
	for i := uint32(0); i < uint32(n-1); i++ {
		ss[i] = 0
		sd[i] = i + 1
		if i == 0 {
			continue
		}
		ss[n32-1+i] = i
		sd[n32-1+i] = i + 1
	}
	ss[n-1] = uint32(n - 1)
	sd[n-1] = 1

	return New(ss, sd)
}
