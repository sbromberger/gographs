package gographs

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/sbromberger/gographs/heap"
	"github.com/sbromberger/gographs/priorityqueue"
)

// DijkstraState is a state holding dijkstra SP info
type oDijkstraState struct {
	Parents      []uint32
	Dists        []float32
	Predecessors [][]uint32
	Pathcounts   []int
}

const oMaxDist = float32(math.MaxFloat32 - 1)

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// DijkstraShortestPaths returns dijkstra shortest paths
func dijkstraShortestPaths(g *Graph, srcs []uint32, withPreds bool) DijkstraState {
	nv := g.Order()
	pq := make(priorityqueue.PriorityQueue, len(srcs))
	visited := make([]bool, nv)
	parents := make([]uint32, nv)
	dists := make([]float32, nv)
	pathcounts := make([]int, nv)
	nPreds := 0
	maxNPreds := 0
	if withPreds {
		nPreds = nv
		maxNPreds = int(0.0015 * float64(nv))
	}
	preds := make([][]uint32, nPreds)

	for i := range preds {
		preds[i] = make([]uint32, maxNPreds)
	}
	for i := range dists {
		dists[i] = MaxDist
	}

	for i, src := range srcs {
		dists[src] = 0
		pathcounts[src] = 1
		visited[src] = true
		pq[i] = &heap.Item{Value: uint32(src), Priority: 0, Index: i}
	}

	heap.Init(&pq)
	for pq.Len() > 0 {
		item := heap.Pop(&pq)
		u := item.Value
		for _, v := range g.OutNeighbors(u) {
			alt := MaxDist
			if dists[u] < MaxDist {
				alt = dists[u] + 1
			}
			if !visited[v] {
				dists[v] = alt
				parents[v] = u
				pathcounts[v] += pathcounts[u]
				visited[v] = true
				if withPreds {
					preds[v] = append(preds[v], u)
				}
				heap.Push(&pq, &heap.Item{Value: v, Priority: min(float32(nv), alt)})
			} else {
				if alt < dists[v] {
					dists[v] = alt
					parents[v] = u
					pathcounts[v] = 0
					if withPreds {
						preds[v] = make([]uint32, maxNPreds)
					}
					heap.Push(&pq, &heap.Item{Value: v, Priority: min(float32(nv), alt)})
				}
				if alt == dists[v] {
					pathcounts[v]++
					if withPreds {
						preds[v] = append(preds[v], u)
					}
				}
			}
		}
	}
	for _, src := range srcs {
		pathcounts[src] = 1
		parents[src] = uint32(src)
		if withPreds {
			preds[src] = make([]uint32, 0)
		}
	}
	return DijkstraState{parents, dists, preds, pathcounts}
}

func DijkstraShortestPaths(g *Graph, s uint32) DijkstraState {
	return dijkstraShortestPaths(g, []uint32{s}, false)
}

func DijkstraShortestPathsWithPreds(g *Graph, s uint32) DijkstraState {
	return dijkstraShortestPaths(g, []uint32{s}, true)
}

type DijkstraStateOneSrc struct {
	src   uint32
	state DijkstraState
}

func multiDijkstra(g *Graph, sVertex, eVertex uint32, ch chan<- DijkstraStateOneSrc) {
	fmt.Printf("in multidijkstra from %d to %d\n", sVertex, eVertex)
	for i := sVertex; i < eVertex; i++ {
		ch <- DijkstraStateOneSrc{i, DijkstraShortestPaths(g, i)}
	}
}

func AllDijkstraShortestPaths(g *Graph) []DijkstraState {
	var wg sync.WaitGroup
	ncpu := runtime.NumCPU() + 2
	nv := g.Order()
	blocksize := nv / ncpu
	fmt.Println("ncpu = ", ncpu)
	fmt.Println("blocksize = ", blocksize)
	blockbounds := make([]uint32, ncpu+1)
	for i := range blockbounds {
		fmt.Println("i = ", i)
		blockbounds[i] = uint32(i * blocksize)
	}
	blockbounds[ncpu] = uint32(nv)
	ds := make([]DijkstraState, nv)
	ch := make(chan DijkstraStateOneSrc, nv)
	fmt.Println("starting goroutines")
	for i := 0; i < ncpu; i++ {
		fmt.Println("starting goroutine ", i)
		wg.Add(1)
		go func(block int) {
			defer wg.Done()
			multiDijkstra(g, blockbounds[block], blockbounds[block+1], ch)
			fmt.Printf("block %d done\n", block)
		}(i)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for d := range ch {
		if d.src == 5999 {
			fmt.Println("got 5999")
			fmt.Println(d.state)
		}
		ds[d.src] = d.state
	}
	// wg.Wait()
	return ds
}
