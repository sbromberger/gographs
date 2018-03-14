// compiled with -gcflags="-B"

package gographs

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/sbromberger/gographs/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

func processOneBlock(ch chan<- []uint32, g *Graph, vs []uint32, visited bitvec.ABitVec) {
	for _, v := range vs {
		neighbors := make([]uint32, 0)
		for _, neighbor := range g.OutNeighbors(v) {
			if !visited.AGet(neighbor) {
				visited.ASet(neighbor)
				neighbors = append(neighbors, neighbor)
			}
		}
		ch <- neighbors
	}
}

// BFSpar computes a vector of levels from src.
func BFSpar(g *Graph, src uint32) {
	np := runtime.GOMAXPROCS(-1) / 4
	fmt.Println("nprocs: ", np)
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

		chunkSize := (len(curLevel) + np - 1) / np
		var workblocks [][]uint32
		for i := 0; i < len(curLevel); i += chunkSize {
			end := i + chunkSize

			if end > len(curLevel) {
				end = len(curLevel)
			}
			workblocks = append(workblocks, curLevel[i:end])
		}
		fmt.Println("  len(workblocks) = ", len(workblocks))
		ch := make(chan []uint32, len(workblocks))
		wg.Add(len(workblocks))
		for _, vs := range workblocks {
			go func(vs []uint32) {
				processOneBlock(ch, g, vs, visited)
				wg.Done()
			}(vs)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		for ns := range ch {
			nextLevel = append(nextLevel, ns...)
			for _, n := range ns {
				vertLevel[n] = nLevel
			}
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
