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
	visited := bitvec.NewBitVec(nv)
	curLevel := make([]uint32, 0, nv)
	nextLevel := make([]uint32, 0, nv)

	nLevel := uint32(2)
	vertLevel[src] = 0
	visited.Set(src)
	curLevel = append(curLevel, src)

	for len(curLevel) > 0 {
		for _, v := range curLevel {
			neighbors := g.OutNeighbors(v)
			i := 0

			for ; i < len(neighbors)-3; i += 4 {
				n1, n2, n3, n4 := neighbors[i], neighbors[i+1], neighbors[i+2], neighbors[i+3]
				if !visited.IsSet(n1) {
					visited.Set(n1)
					nextLevel = append(nextLevel, n1)
				}
				if !visited.IsSet(n2) {
					visited.Set(n2)
					nextLevel = append(nextLevel, n2)
				}
				if !visited.IsSet(n3) {
					visited.Set(n3)
					nextLevel = append(nextLevel, n3)
				}
				if !visited.IsSet(n4) {
					visited.Set(n4)
					nextLevel = append(nextLevel, n4)
				}
			}

			for _, n := range neighbors[i:] {
				if !visited.IsSet(n) {
					visited.Set(n)
					nextLevel = append(nextLevel, n)
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
