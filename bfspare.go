// compiled with -gcflags="-B"

package gographs

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/egonelbre/async"
	"github.com/sbromberger/gographs/bitvec"
	"github.com/shawnsmithdev/zermelo/zuint32"
)

const (
	ReadBlockSize  = 256
	WriteBlockSize = 256
	MaxBlockSize   = 256 // max(ReadBlockSize, WriteBlockSize)
	EmptySentinel  = ^uint32(0)
)

type Frontier struct {
	Data []uint32
	Head uint32
}

func (front *Frontier) NextRead() (low, high uint32) {
	high = atomic.AddUint32(&front.Head, ReadBlockSize)
	low = high - ReadBlockSize
	if high > uint32(len(front.Data)) {
		high = uint32(len(front.Data))
	}
	return
}

func (front *Frontier) NextWrite() (low, high uint32) {
	high = atomic.AddUint32(&front.Head, ReadBlockSize)
	low = high - ReadBlockSize
	return
}

func (front *Frontier) Write(low, high *uint32, v uint32) {
	if *low >= *high {
		*low, *high = front.NextWrite()
	}
	front.Data[*low] = v
	*low += 1
}

func processLevel(g *Graph, currLevel, nextLevel *Frontier, visited *bitvec.BBitVec) {
	runtime.LockOSThread()

	writeLow, writeHigh := uint32(0), uint32(0)
	for {
		readLow, readHigh := currLevel.NextRead()
		if readLow >= readHigh {
			break
		}

		for _, v := range currLevel.Data[readLow:readHigh] {
			if v == EmptySentinel {
				continue
			}

			neighbors := g.OutNeighbors(v)
			i := 0
			for ; i < len(neighbors)-3; i += 4 {
				n1, n2, n3, n4 := neighbors[i], neighbors[i+1], neighbors[i+2], neighbors[i+3]
				if visited.TrySet(n1) {
					nextLevel.Write(&writeLow, &writeHigh, n1)
				}
				if visited.TrySet(n2) {
					nextLevel.Write(&writeLow, &writeHigh, n2)
				}
				if visited.TrySet(n3) {
					nextLevel.Write(&writeLow, &writeHigh, n3)
				}
				if visited.TrySet(n4) {
					nextLevel.Write(&writeLow, &writeHigh, n4)
				}
			}
			for _, n := range neighbors[i:] {
				if visited.TrySet(n) {
					nextLevel.Write(&writeLow, &writeHigh, n)
				}
			}
		}
	}

	for i := writeLow; i < writeHigh; i += 1 {
		nextLevel.Data[i] = EmptySentinel
	}
}

// BFSxp computes a vector of levels from src.
func BFSpare(g *Graph, src uint32, procs int) {
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
	for len(currLevel.Data) > 0 {

		async.Spawn(procs, func(i int) {
			processLevel(g, currLevel, nextLevel, &visited)
		}, func() { wait <- struct{}{} })

		<-wait

		nextLevel.Data = nextLevel.Data[:nextLevel.Head]
		nextLevel.Head = 0

		sentinelCount := uint32(0)
		async.BlockIter(len(nextLevel.Data), procs, func(low, high int) {
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
		nextLevel.Data = nextLevel.Data[:maxSize:maxSize]
		nextLevel.Head = 0
	}
}

// BFSpare2 computes a vector of levels from src.
func BFSpare2(g *Graph, src uint32, procs int) {
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

	var waitForLast1, waitForLast2 sync.WaitGroup
	doneProcessingCounter := int32(procs)
	// busyWait := int32(procs)
	waitForLast1.Add(procs)

	allDone := uint32(0)
	sentinelCount := uint32(0)

	worker := func(gid int) {
		runtime.LockOSThread()

		for atomic.LoadUint32(&allDone) == 0 {
			processLevel(g, currLevel, nextLevel, &visited)

			if atomic.AddInt32(&doneProcessingCounter, -1) == 0 {
				nextLevel.Data = nextLevel.Data[:nextLevel.Head]
				nextLevel.Head = 0
				atomic.StoreUint32(&sentinelCount, 0)

				atomic.StoreInt32(&doneProcessingCounter, int32(procs))
				waitForLast2.Add(procs)
				waitForLast1.Done()
			} else {
				waitForLast1.Done()
				waitForLast1.Wait()
			}

			// wait for sorting signal
			blockSize := (len(nextLevel.Data) + procs - 1) / procs
			low := blockSize * gid
			high := low + blockSize
			if high > len(nextLevel.Data) {
				high = len(nextLevel.Data)
			}

			zuint32.SortBYOB(nextLevel.Data[low:high], currLevel.Data[low:high])
			for index, v := range nextLevel.Data[low:high] {
				if v == EmptySentinel {
					atomic.AddUint32(&sentinelCount, uint32(high-low-index))
					break
				}
				vertLevel[v] = currentLevel
			}

			if atomic.AddInt32(&doneProcessingCounter, -1) == 0 {
				fmt.Printf("completed level %d, size = %d\n", currentLevel-1, len(nextLevel.Data))

				currentLevel++
				currLevel, nextLevel = nextLevel, currLevel
				nextLevel.Data = nextLevel.Data[:maxSize:maxSize]
				nextLevel.Head = 0

				if len(currLevel.Data) == 0 {
					atomic.StoreUint32(&allDone, 1)
				}

				atomic.StoreInt32(&doneProcessingCounter, int32(procs))
				waitForLast1.Add(procs)
				waitForLast2.Done()
			} else {
				waitForLast2.Done()
				waitForLast2.Wait()
			}
		}
	}

	for gid := 1; gid < procs; gid++ {
		go worker(gid)
	}
	worker(0)
}
