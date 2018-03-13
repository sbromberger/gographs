// compiled with -gcflags="-B"

package gographs

import (
	"fmt"

	"github.com/sbromberger/gographs/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

// BFSxp computes a vector of levels from src.
func BFSxp(g *Graph, src uint32) {
	nv := g.Order()
	vertLevel := make([]uint32, nv)
	visited := bitvec.NewABitVec(nv)
	curLevel := make([]uint32, 0, nv)
	nextLevel := make([]uint32, 0, nv)

	nLevel := uint32(2)
	vertLevel[src] = 0
	visited.Set(src)
	curLevel = append(curLevel, src)

	for len(curLevel) > 0 {
		for _, v := range curLevel {
			for _, neighbor := range g.OutNeighbors(v) {
				if !visited.Get(neighbor) {
					nextLevel = append(nextLevel, neighbor)
					visited.Set(neighbor)
				}
			}
		}

		zuint32.SortBYOB(nextLevel, curLevel[:nv])
		for _, v := range nextLevel {
			vertLevel[v] = nLevel
		}

		fmt.Printf("completed level %d, size = %d\n", nLevel-1, len(nextLevel))

		nLevel++
		curLevel, nextLevel = nextLevel, curLevel[:0]
	}
}
