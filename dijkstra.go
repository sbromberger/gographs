package graph

import (
	"fmt"
	"math"

	"github.com/shawnsmithdev/zermelo/zuint32"
)

// DijkstraState is a state holding dijkstra SP info
type DijkstraState struct {
	Parents      []uint32
	Dists        []float32
	Pathcounts   []uint32
	Predecessors [][]uint32
}

func (d DijkstraState) String() string {
	s := fmt.Sprintln("DijkstraState with")
	s += fmt.Sprintln("  Parents:      ", d.Parents)
	s += fmt.Sprintln("  Dists:        ", d.Dists)
	s += fmt.Sprintln("  PathCounts:   ", d.Pathcounts)
	s += fmt.Sprintln("  Predecessors: ", d.Predecessors)
	return s
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

const maxDist = float32(math.MaxFloat32 - 1)

// Dijkstra returns a DijkstraState.
func Dijkstra(g Graph, src uint32, weightFn func(uint32, uint32) float32, withPreds bool) DijkstraState {
	nv := g.NumVertices()
	vertLevel := make([]uint32, nv)
	for i := uint32(0); i < nv; i++ {
		vertLevel[i] = unvisited
	}
	curLevel := make([]uint32, 0, nv)
	nextLevel := make([]uint32, 0, nv)
	nLevel := uint32(2)
	parents := make([]uint32, nv)
	pathcounts := make([]uint32, nv)
	dists := make([]float32, nv)

	preds := make([][]uint32, 0)
	if withPreds {
		preds = make([][]uint32, nv)
	}

	for i := range dists {
		dists[i] = maxDist
	}

	vertLevel[src] = 0
	dists[src] = 0
	parents[src] = src
	pathcounts[src] = 1
	curLevel = append(curLevel, src)
	for len(curLevel) > 0 {
		for _, u := range curLevel {
			for _, v := range g.OutNeighbors(u) {
				alt := min(maxDist, dists[u]+weightFn(u, v))
				if vertLevel[v] == unvisited { // if not visited
					dists[v] = alt
					parents[v] = u
					pathcounts[v] += pathcounts[u]
					if withPreds {
						preds[v] = append(preds[v], u)
					}
					nextLevel = append(nextLevel, v)
					vertLevel[v] = nLevel
				} else {
					if alt < dists[v] {
						dists[v] = alt
						parents[v] = u
						pathcounts[v] = 0
						if withPreds {
							preds[v] = preds[v][:0]
						}
					}
					if alt == dists[v] {
						pathcounts[v] += pathcounts[u]
						if withPreds {
							preds[v] = append(preds[v], u)
						}
					}
				}
			}
		}
		fmt.Printf("completed level %d, size = %d\n", nLevel-1, len(nextLevel))
		nLevel++
		curLevel = curLevel[:0]
		curLevel, nextLevel = nextLevel, curLevel
		zuint32.SortBYOB(curLevel, nextLevel[:nv])
	}
	pathcounts[src] = 1
	parents[src] = 0
	if withPreds {
		preds[src] = preds[src][:0]
	}
	ds := DijkstraState{
		Parents:      parents,
		Dists:        dists,
		Pathcounts:   pathcounts,
		Predecessors: preds,
	}
	return ds
}
