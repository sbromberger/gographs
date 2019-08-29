package graph

import (
	"math"

	"github.com/sbromberger/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

const unvisited = math.MaxUint32

// BFS computes a vector of levels from src and returns a vector
// of vertices visited in order along with a vector of distances
// indexed by vertex.
func BFS(g Graph, src uint32) (vertexList, vertLevel []uint32) {
	nv := g.NumVertices()
	vertLevel = make([]uint32, nv)
	for i := uint32(0); i < nv; i++ {
		vertLevel[i] = unvisited
	}

	visited := bitvec.NewBitVec(nv)
	curLevel := make([]uint32, 0, nv)
	nextLevel := make([]uint32, 0, nv)
	nLevel := uint32(1)
	vertLevel[src] = 0
	visited.TrySet(src)
	curLevel = append(curLevel, src)
	vertexList = make([]uint32, 0, nv)
	vertexList = append(vertexList, src)
	for len(curLevel) > 0 {
		for _, v := range curLevel {
			for _, neighbor := range g.OutNeighbors(v) {
				if visited.TrySet(neighbor) {
					nextLevel = append(nextLevel, neighbor)
					vertLevel[neighbor] = nLevel
					vertexList = append(vertexList, neighbor)
				}
			}
		}
		nLevel++
		curLevel = curLevel[:0]
		curLevel, nextLevel = nextLevel, curLevel
		zuint32.SortBYOB(curLevel, nextLevel[:nv])
	}
	return vertexList, vertLevel
}
