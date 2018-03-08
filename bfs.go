package gographs

import "fmt"

// BFS returns a vector of levels from src.
func BFS(g *Graph, src uint32) []uint32 {
	nv := g.Order()
	vertLevel := make([]uint32, nv)
	visited := make([]bool, nv)
	curLevel := make([]uint32, 0, nv)
	nextLevel := make([]uint32, 0, nv)
	nLevel := uint32(2)
	vertLevel[src] = 0
	visited[src] = true
	curLevel = append(curLevel, src)
	for len(curLevel) > 0 {
		for _, v := range curLevel {
			for _, neighbor := range g.OutNeighbors(v) {
				if !visited[neighbor] {
					nextLevel = append(nextLevel, neighbor)
					vertLevel[neighbor] = nLevel
					visited[neighbor] = true
				}
			}
		}
		fmt.Printf("completed level %d, size = %d\n", nLevel-1, len(nextLevel))
		nLevel++
		curLevel = curLevel[:0]
		curLevel, nextLevel = nextLevel, curLevel
		// sort.Slice(curLevel, func(i, j int) bool { return curLevel[i] < curLevel[j] })
	}
	return vertLevel
}
