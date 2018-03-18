// compiled with -gcflags="-B"

package gographs

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/egonelbre/async"
	"github.com/sbromberger/gographs/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

const (
	ReadBlockSize  = 256 // number of neighbors to process per block
	WriteBlockSize = 256 // number of empty cells to allocate initially for nextLevel
	MaxBlockSize   = 256 // max(ReadBlockSize, WriteBlockSize)
	EmptySentinel  = ^uint32(0)
)

type Frontier struct {
	Data []uint32 // the vector of data
	Head uint32   // the index to the next unused element in the vector
}

// NextRead returns the low and high offsets into Frontier for reading ReadBlockSize blocks.
// It increases Head by the ReadBlockSize.
// Note: we only read from currLevel.
func (front *Frontier) NextRead() (low, high uint32) {
	high = atomic.AddUint32(&front.Head, ReadBlockSize)
	low = high - ReadBlockSize
	if high > uint32(len(front.Data)) {
		high = uint32(len(front.Data))
	}
	return
}

// NextWrite returns the low and high offsets into Frontier for writing WriteBlockSize blocks.
// It increases Head by WriteBlockSize.
// Note: we only write to nextLevel.
func (front *Frontier) NextWrite() (low, high uint32) {
	high = atomic.AddUint32(&front.Head, WriteBlockSize)
	low = high - WriteBlockSize
	return
}

// Write inserts `v` into the next available position in the Frontier, allocating as necessary.
// Note: we only write to nextLevel.
func (front *Frontier) Write(low, high *uint32, v uint32) {
	if *low >= *high {
		*low, *high = front.NextWrite()
	}
	front.Data[*low] = v
	*low++
}

// processLevel uses Frontiers to dequeue work from currLevel in ReadBlockSize increments.
func processLevel(g GoGraph, currLevel, nextLevel *Frontier, visited *bitvec.BBitVec) {
	writeLow, writeHigh := uint32(0), uint32(0)
	for {
		readLow, readHigh := currLevel.NextRead() // if currLevel still has vertices to process, get the indices of a ReadBlockSize block of them
		if readLow >= readHigh {                  // otherwise exit
			break
		}

		for _, v := range currLevel.Data[readLow:readHigh] { // get and loop through a slice of ReadBlockSize vertices from currLevel
			if v == EmptySentinel { // if we hit a sentinel within the block, skip it
				continue
			}

			neighbors := g.OutNeighbors(v) // get the outNeighbors of the vertex under inspection
			i := 0
			for ; i < len(neighbors)-3; i += 4 { // unroll loop for visited
				n1, n2, n3, n4 := neighbors[i], neighbors[i+1], neighbors[i+2], neighbors[i+3]
				x1, x2, x3, x4 := visited.GetBuckets4(n1, n2, n3, n4)
				if visited.TrySetWith(x1, n1) { // if not visited, add to the list of vertices for nextLevel
					nextLevel.Write(&writeLow, &writeHigh, n1)
				}
				if visited.TrySetWith(x2, n2) {
					nextLevel.Write(&writeLow, &writeHigh, n2)
				}
				if visited.TrySetWith(x3, n3) {
					nextLevel.Write(&writeLow, &writeHigh, n3)
				}
				if visited.TrySetWith(x4, n4) {
					nextLevel.Write(&writeLow, &writeHigh, n4)
				}
			}
			for _, n := range neighbors[i:] { // process any remaining (< 4) neighbors for this vertex
				if visited.TrySet(n) {
					nextLevel.Write(&writeLow, &writeHigh, n)
				}
			}
		}
	}

	for i := writeLow; i < writeHigh; i++ {
		nextLevel.Data[i] = EmptySentinel // ensure the rest of the nextLevel block is "empty" using the sentinel
	}
}

// BFSpare computes a vector of levels from src in parallel.
func BFSpare(g GoGraph, src uint32, procs int) {
	N := g.Order()
	vertLevel := make([]uint32, N)
	visited := bitvec.NewBBitVec(N)

	maxSize := N + MaxBlockSize*procs
	currLevel := &Frontier{make([]uint32, 0, maxSize), 0}
	nextLevel := &Frontier{make([]uint32, maxSize, maxSize), 0}

	currentLevel := uint32(2)
	vertLevel[src] = 0
	visited.TrySet(src)

	currLevel.Data = append(currLevel.Data, src)

	wait := make(chan struct{})
	for len(currLevel.Data) > 0 { // while we have vertices in currentLevel

		async.Spawn(procs, func(i int) { // spawn `procs` goroutines to process vertices in this level,
			runtime.LockOSThread() // using currLevel as the work queue. Make sure only one goroutine per thread.
			processLevel(g, currLevel, nextLevel, &visited)
		}, func() { wait <- struct{}{} })

		<-wait // this is equivalent to using a WaitGroup but uses a single channel message instead.

		nextLevel.Data = nextLevel.Data[:nextLevel.Head] // "truncate" nextLevel.Data to just the valid data...
		// ... we need to do this because Frontier.ReadNext uses `len`.

		sentinelCount := uint32(0)
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
}
