package graph

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/egonelbre/async"
	"github.com/sbromberger/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

// processDijkstraLevel uses Frontiers to dequeue work from currLevel in ReadBlockSize increments.
func processDijkstraLevel(g Graph, currLevel, nextLevel *Frontier, visited *bitvec.ABitVec, dists []float32, parents []uint32, pathcounts []uint32) {
	writeLow, writeHigh := u0, u0
	for {
		readLow, readHigh := currLevel.NextRead() // if currLevel still has vertices to process, get the indices of a ReadBlockSize block of them
		if readLow >= readHigh {                  // otherwise exit
			break
		}

		for _, u := range currLevel.Data[readLow:readHigh] { // get and loop through a slice of ReadBlockSize vertices from currLevel
			if u == EmptySentinel { // if we hit a sentinel within the block, skip it
				continue
			}
			alt := min(maxDist, dists[u]+1)
			vs := g.OutNeighbors(u) // get the outNeighbors of the vertex under inspection
			// i := 0
			for _, v := range vs {
				// for ; i < len(vs)-3; i += 4 { // unroll loop for visited
				// 	v1, v2, v3, v4 := vs[i], vs[i+1], vs[i+2], vs[i+3]
				// 	x1, x2, x3, x4 := visited.GetBuckets4(v1, v2, v3, v4)
				// 	if visited.TrySetWith(x1, v1) { // if not visited, add to the list of vertices for nextLevel
				// 		nextLevel.Write(&writeLow, &writeHigh, v1)
				// 		dists[v1] = alt
				// 		parents[v1] = u
				// 		pathcounts[v1] += pathcounts[u]
				// 	} else {
				// 		if alt < dists[v1] {
				// 			dists[v1] = alt
				// 			parents[v1] = u
				// 			pathcounts[v1] = 0
				// 		}
				// 		if alt == dists[v1] {
				// 			pathcounts[v1] += pathcounts[u]
				// 		}
				// 	}
				// 	if visited.TrySetWith(x2, v2) {
				// 		nextLevel.Write(&writeLow, &writeHigh, v2)
				// 		dists[v2] = alt
				// 		dists[v2] = alt
				// 		parents[v2] = u
				// 		pathcounts[v2] += pathcounts[u]
				// 	} else {
				// 		if alt < dists[v2] {
				// 			dists[v2] = alt
				// 			parents[v2] = u
				// 			pathcounts[v2] = 0
				// 		}
				// 		if alt == dists[v2] {
				// 			pathcounts[v2] += pathcounts[u]
				// 		}
				// 	}
				// 	if visited.TrySetWith(x3, v3) {
				// 		nextLevel.Write(&writeLow, &writeHigh, v3)
				// 		dists[v3] = alt
				// 		dists[v3] = alt
				// 		parents[v3] = u
				// 		pathcounts[v3] += pathcounts[u]
				// 	} else {
				// 		if alt < dists[v3] {
				// 			dists[v3] = alt
				// 			parents[v3] = u
				// 			pathcounts[v3] = 0
				// 		}
				// 		if alt == dists[v3] {
				// 			pathcounts[v3] += pathcounts[u]
				// 		}
				// 	}
				// 	if visited.TrySetWith(x4, v4) {
				// 		nextLevel.Write(&writeLow, &writeHigh, v4)
				// 		dists[v4] = alt
				// 		dists[v4] = alt
				// 		parents[v4] = u
				// 		pathcounts[v4] += pathcounts[u]
				// 	} else {
				// 		if alt < dists[v4] {
				// 			dists[v4] = alt
				// 			parents[v4] = u
				// 			pathcounts[v4] = 0
				// 		}
				// 		if alt == dists[v4] {
				// 			pathcounts[v4] += pathcounts[u]
				// 		}
				// 	}
				// }
				// for _, v := range vs[i:] { // process any remaining (< 4) neighbors for this vertex
				if visited.TrySet(v) {
					nextLevel.Write(&writeLow, &writeHigh, v)
					dists[v] = alt
					parents[v] = u
					atomic.AddUint32(&pathcounts[v], pathcounts[u])
				} else {
					if alt < dists[v] {
						dists[v] = alt
						parents[v] = u
						pathcounts[v] = 0
					}
					if alt == dists[v] {
						atomic.AddUint32(&pathcounts[v], pathcounts[u])
					}
				}
			}
		}
	}

	for i := writeLow; i < writeHigh; i++ {
		nextLevel.Data[i] = EmptySentinel // ensure the rest of the nextLevel block is "empty" using the sentinel
	}
}

// ParallelDijkstra computes a vector of levels from src in parallel.
func ParallelDijkstra(g Graph, src uint32, procs int) DijkstraState {
	N := g.NumVertices()
	vertLevel := make([]uint32, N)
	visited := bitvec.NewABitVec(N)

	maxSize := N + MaxBlockSize*uint32(procs)
	currLevel := &Frontier{make([]uint32, 0, maxSize), 0}
	nextLevel := &Frontier{make([]uint32, maxSize), 0}

	currentLevel := uint32(2)
	parents := make([]uint32, N)
	pathcounts := make([]uint32, N)
	dists := make([]float32, N)

	for i := range dists {
		dists[i] = maxDist
	}

	vertLevel[src] = 0
	dists[src] = 0
	parents[src] = src
	pathcounts[src] = 1
	visited.TrySet(src)

	currLevel.Data = append(currLevel.Data, src)

	wait := make(chan struct{})
	for len(currLevel.Data) > 0 { // while we have vertices in currentLevel

		async.Spawn(procs, func(i int) { // spawn `procs` goroutines to process vertices in this level,
			runtime.LockOSThread() // using currLevel as the work queue. Make sure only one goroutine per thread.
			processDijkstraLevel(g, currLevel, nextLevel, &visited, dists, parents, pathcounts)
		}, func() { wait <- struct{}{} })

		<-wait // this is equivalent to using a WaitGroup but uses a single channel message instead.

		nextLevel.Data = nextLevel.Data[:nextLevel.Head] // "truncate" nextLevel.Data to just the valid data...
		// ... we need to do this because Frontier.ReadNext uses `len`.

		sentinelCount := u0
		// now sort nextLevel by block. After this, all data within a given block will be sorted. This ensures that
		// "most" data are ordered, which preserves some linearity in cache access, but this might not be significant.
		// More testing is needed.
		async.BlockIter(int(nextLevel.Head), procs, func(low, high int) {
			zuint32.SortBYOB(nextLevel.Data[low:high], currLevel.Data[low:high])
			for index, v := range nextLevel.Data[low:high] {
				if v == EmptySentinel {
					atomic.AddUint32(&sentinelCount, uint32(high-low-index))
					break
				}
				vertLevel[v] = currentLevel
			}
		})

		fmt.Printf("completed level %d, size = %d\n", currentLevel-1, len(nextLevel.Data)-int(sentinelCount))

		currentLevel++
		currLevel, nextLevel = nextLevel, currLevel
		currLevel.Head = 0 // start reading from 0
		// reset buffer for next level
		nextLevel.Data = nextLevel.Data[:maxSize:maxSize] // resize the buffer to `maxSize` elements. We don't care what's in it, because...
		nextLevel.Head = 0                                // ... we start writing to index 0.
	}
	ds := DijkstraState{
		Parents:    parents,
		Dists:      dists,
		Pathcounts: pathcounts,
	}
	return ds
}
