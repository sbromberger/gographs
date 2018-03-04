package gographs

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"

	"github.com/sbromberger/gographs/priorityqueue"
	"github.com/sbromberger/gographs/sparsevecs"
)

type Edge struct {
	src uint32
	dst uint32
}

type GoGraph interface {
	IsDirected() bool
	HasEdge(u, v int32) bool
	Fadj() *sparsevecs.UInt32SparseVec
	Badj() *sparsevecs.UInt32SparseVec
	// AddEdge(u, v int) err
	InNeighbors(u int32) []int32
	OutNeighbors(u int32) []int32
}

// Graph is a graph structure
type Graph struct {
	mx sparsevecs.UInt32SparseVec
}

type DiGraph struct {
	fmx, bmx sparsevecs.UInt32SparseVec
}

func (g *Graph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

func (g *Graph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.mx
}

func (g *DiGraph) Fadj() *sparsevecs.UInt32SparseVec {
	return &g.fmx
}

func (g *DiGraph) Badj() *sparsevecs.UInt32SparseVec {
	return &g.bmx
}

func (g *Graph) IsDirected() bool {
	return false
}

func (g DiGraph) IsDirected() bool {
	return true
}

// MakeGraph creates an undirected graph from rowidx and colptr vectors.
func MakeGraph(r, c []uint32) Graph {
	return Graph{sparsevecs.UInt32SparseVec{r, c}}
}

// HasEdge returns true if an edge exists between u and v.
func (g *Graph) HasEdge(u, v int) bool {
	f := g.Fadj()
	return f.GetIndexInt(u, v) || f.GetIndexInt(v, u)
}

func (g *DiGraph) HasEdge(u, v int) bool {
	f := g.Fadj()
	b := g.Badj()
	return f.GetIndexInt(u, v) || b.GetIndexInt(v, u)
}

// AddEdge Adds an edge to graph g
func (g *Graph) AddEdge(u, v uint32) error {
	return g.mx.Insert(u, v)
}

// Size returns the number of edges
func (g *Graph) Size() int {
	return len(g.mx.Rowidx) / 2
}

func (g DiGraph) Size() int {
	return len(g.fmx.Rowidx) / 2
}

// Order returns the number of vertices
func (g *Graph) Order() int {
	return len(g.mx.Colptr) - 1
}

// Order returns the number of vertices
func (g *DiGraph) Order() int {
	return len(g.fmx.Colptr) - 1
}

// Edges returns a list of edges
func (g *Graph) Edges() []Edge {
	ne := g.Size()
	fmt.Println("ne = ", ne)
	edgelist := make([]Edge, ne)
	edgesAdded := 0
	for c := range g.mx.Colptr[:len(g.mx.Colptr)-1] {
		uc := uint32(c)
		rrange := g.mx.GetRange(uc)
		for _, r := range rrange {
			if r <= uc {
				fmt.Printf("Adding edge %d, %d\n", r, uc)
				edgelist[edgesAdded] = Edge{r, uc}
				edgesAdded++
			}
		}
	}
	return edgelist
}

// AddEdgeInt takes ints and calls AddEdge
func (g *Graph) AddEdgeInt(u, v int) error {
	return g.AddEdge(uint32(u), uint32(v))

}

// OutNeighbors returns the neighbors of vertex v.
func (g *Graph) OutNeighbors(v uint32) []uint32 {
	return g.mx.GetRange(v)
}

// OutNeighborsInt returns the neighbors of vertex v.
func (g *Graph) OutNeighborsInt(v int) []uint32 {
	return g.mx.GetRangeInt(v)
}

// DijkstraState is a state holding dijkstra SP info
type DijkstraState struct {
	Parents      []uint32
	Dists        []int
	Predecessors [][]uint32
	Pathcounts   []int
}

const MaxDist = math.MaxInt64

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DijkstraShortestPaths returns dijkstra shortest paths
func dijkstraShortestPaths(g *Graph, srcs []uint32, withPreds bool) DijkstraState {
	nv := g.Order()
	pq := priorityqueue.New(nv + 1)
	visited := make([]bool, nv)
	parents := make([]uint32, nv)
	dists := make([]int, nv)
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

	for _, src := range srcs {
		dists[src] = 0
		pathcounts[src] = 1
		visited[src] = true
		pq.Push(uint32(src), 0)
	}

	for !pq.IsEmpty() {
		up, err := pq.Pop()
		if err != nil {
			log.Fatal("error in dequeue: ", err)
		}
		u := up.Value
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
				pq.Push(v, min(nv, alt))
			} else {
				if alt < dists[v] {
					dists[v] = alt
					parents[v] = u
					pathcounts[v] = 0
					if withPreds {
						preds[v] = make([]uint32, maxNPreds)
					}
					pq.Push(v, min(nv, alt))
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
