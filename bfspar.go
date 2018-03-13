// compiled with -gcflags="-B"

package gographs

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/sbromberger/gographs/bitvec"
)

type VisitedStruct struct {
	vec   bitvec.BitVec
	mutex sync.RWMutex
}

type NextLevelStruct struct {
	vec   []uint32
	mutex sync.RWMutex
}

func NewVisitedStruct(n int) VisitedStruct {
	v := bitvec.NewBitVec(n)
	m := new(sync.RWMutex)
	return VisitedStruct{v, *m}
}

func NewNextLevelStruct(n int) NextLevelStruct {
	v := make([]uint32, 0, n)
	m := new(sync.RWMutex)
	return NextLevelStruct{v, *m}
}

func processOneBlock(wg *sync.WaitGroup, g *Graph, vs []uint32, visited *VisitedStruct, vertLevel *[]uint32, nextLevel *NextLevelStruct, nLevel uint32) {
	// fmt.Printf("    processing %d vertices\n", len(vs))
	defer wg.Done()
	tlQueue := make([]uint32, 0, len(vs))
	for _, v := range vs {
		for _, neighbor := range g.OutNeighbors(v) {
			// fmt.Println("   processing neighbor ", neighbor)
			visited.mutex.RLock()
			if !visited.vec.IsSet(neighbor) {
				visited.mutex.RUnlock()
				visited.mutex.Lock()
				visited.vec.Set(neighbor)
				visited.mutex.Unlock()
				(*vertLevel)[neighbor] = nLevel
				// fmt.Printf("  appending %d to tlQueue (len %d)\n", neighbor, len(tlQueue))
				tlQueue = append(tlQueue, neighbor)
			} else {
				visited.mutex.RUnlock()
			}
		}
	}
	if len(tlQueue) > 0 {
		nextLevel.mutex.Lock()
		// fmt.Println("len(tlQueue) = ", len(tlQueue))
		nextLevel.vec = append(nextLevel.vec, tlQueue...)
		// fmt.Println("len(nextLevel.vec) = %d", len(nextLevel.vec))
		nextLevel.mutex.Unlock()
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// BFSpar computes a vector of levels from src in parallel.
func BFSpar(g *Graph, src uint32) {
	np := runtime.NumCPU()
	fmt.Println("GOMAXPROCS = ", runtime.GOMAXPROCS(-1))
	nBlocks := np
	fmt.Printf("Using %d cores\n", np)
	nv := g.Order()
	vertLevel := make([]uint32, nv)
	visited := NewVisitedStruct(nv)
	curLevel := make([]uint32, 0, nv)
	nextLevel := NewNextLevelStruct(nv)
	nLevel := uint32(2)
	vertLevel[src] = 0
	visited.vec.Set(src)
	curLevel = append(curLevel, src)
	var wg sync.WaitGroup

	for len(curLevel) > 0 {
		curLen := len(curLevel)
		blockSize := curLen / nBlocks
		fmt.Printf("spawning %d goroutines of length %d\n", nBlocks, blockSize)
		// fmt.Println("blockSize = ", blockSize)
		// fmt.Println("curLen = ", curLen)
		for i := 0; i < nBlocks-1; i++ {
			blockStart := i * blockSize
			blockEnd := minInt((i+1)*(blockSize), curLen)
			vs := curLevel[blockStart:blockEnd]
			// fmt.Printf("spawning process %d with %d elems starting at offset %d\n", i, len(vs), blockStart)
			wg.Add(1)
			go processOneBlock(&wg, g, vs, &visited, &vertLevel, &nextLevel, nLevel)
		}

		blockStart := blockSize * (nBlocks)
		blockEnd := curLen
		fmt.Println("adding last block of size ", blockEnd-blockStart)
		wg.Add(1)
		go processOneBlock(&wg, g, curLevel[blockStart:blockEnd], &visited, &vertLevel, &nextLevel, nLevel)
		wg.Wait()
		fmt.Printf("completed level %d, size = %d\n", nLevel-1, len(nextLevel.vec))
		nLevel++
		curLevel, nextLevel.vec = nextLevel.vec, curLevel[:0]
		// fmt.Println("curLevel at bottom = ", curLevel)
	}
}
