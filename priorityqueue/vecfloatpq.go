package priorityqueue

import (
	"errors"
	"math"
	"sort"
)

type Float2Item struct {
	Val uint32
	Pri float64
}

type Float2PQ struct {
	Vals []uint32
	Size uint32
	Pris []float64
}

func newFloat2Item() Float2Item {
	return Float2Item{Val: uint32(math.MaxUint32 - 1), Pri: math.MaxFloat64}
}

// IsEmpty returns true if PQ is empty.
func (pq *Float2PQ) IsEmpty() bool {
	return pq.Size == 0
}

func NewFloat2PQ(n int) Float2PQ {
	vs := make([]uint32, n)
	ps := make([]float64, n)
	for i := 0; i < n; i++ {
		vs[i] = uint32(math.MaxUint32 - 1)
		ps[i] = math.MaxFloat64
	}
	return Float2PQ{Vals: vs, Size: 0, Pris: ps}
}

// Push adds a value/priority to a PQ
func (pq *Float2PQ) Push(val uint32, pri float64) uint32 {
	n := sort.SearchFloat64s(pq.Pris[:pq.Size], pri)

	pq.Pris = append(pq.Pris, math.MaxFloat64)
	pq.Vals = append(pq.Vals, uint32(math.MaxUint32-1))
	copy(pq.Pris[n+1:pq.Size+1], pq.Pris[n:pq.Size])
	copy(pq.Vals[n+1:pq.Size+1], pq.Vals[n:pq.Size])
	pq.Pris[n] = pri
	pq.Vals[n] = val
	pq.Size++
	return pq.Size
}

// Pop retrieves the highest-priority IntIntItem from the PQ
func (pq *Float2PQ) Pop() (Float2Item, error) {
	if pq.IsEmpty() || pq.Vals[0] == uint32(math.MaxUint32-1) {
		return newFloat2Item(), errors.New("queue is empty")
	}
	p := pq.Pris[0]
	v := pq.Vals[0]
	pq.Pris = pq.Pris[1:]
	pq.Vals = pq.Vals[1:]
	pq.Size--
	return Float2Item{Val: v, Pri: p}, nil
}
