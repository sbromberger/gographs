// compiled with -gcflags="-B"

package gographs

import (
	"fmt"
	"sync"

	"github.com/sbromberger/gographs/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

func processOneV(ch chan<- uint32, g *Graph, v uint32, visited bitvec.ABitVec) {
	for _, neighbor := range g.OutNeighbors(v) {
		if !visited.AGet(neighbor) {
			visited.ASet(neighbor)
			ch <- neighbor
		}
	}
}

// BFSpar computes a vector of levels from src.
func BFSpar(g *Graph, src uint32) {
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
		var wg sync.WaitGroup
		ch := make(chan uint32)
		for _, v := range curLevel {
			wg.Add(1)
			go func() {
				processOneV(ch, g, v, visited)
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		for n := range ch {
			nextLevel = append(nextLevel, n)
			vertLevel[n] = nLevel
		}
		zuint32.SortBYOB(nextLevel, curLevel[:nv])
		// for _, v := range nextLevel {
		// 	vertLevel[v] = nLevel
		// }

		fmt.Printf("completed level %d, size = %d\n", nLevel-1, len(nextLevel))

		nLevel++
		curLevel, nextLevel = nextLevel, curLevel[:0]
	}
}
