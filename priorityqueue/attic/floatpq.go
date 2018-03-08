package priorityqueue

import (
	"errors"
	"math"
	"sort"
)

// FloatItem holds value and priority
type FloatItem struct {
	Val uint32
	Pri float64
}

func newFloatItem() FloatItem {
	return FloatItem{Val: 0, Pri: math.MaxFloat64}
}

// FloatPQ implements a priority queue.
type FloatPQ struct {
	Q    []FloatItem
	Size int
}

// IsEmpty returns true if PQ is empty.
func (pq *FloatPQ) IsEmpty() bool {
	return pq.Size == 0
}

// NewFloatPQ creates a PQ with `n` priorities.
func NewFloatPQ(n int) FloatPQ {
	q := make([]FloatItem, n)
	for i := range q {
		q[i] = newFloatItem()
	}
	return FloatPQ{Q: q, Size: 0}
}

// Push pushes a value and priority to the PQ
func (pq *FloatPQ) Push(v uint32, p float64) int {
	n := sort.Search(pq.Size, func(n int) bool {
		return pq.Q[n].Pri >= p
	})

	item := FloatItem{Val: v, Pri: p}
	if n == 0 { // at beginning
		pq.Q = append([]FloatItem{item}, pq.Q...)
		pq.Size++
		return pq.Size
	}
	if n == len(pq.Q) { // at end
		pq.Q = append(pq.Q, item)
		pq.Size++
		return pq.Size
	}
	if n == pq.Size { // at the end of the values, but before end of slice
		pq.Q[n] = item
		pq.Size++
		return pq.Size
	}

	// insert at position n
	pq.Q = append(pq.Q, newFloatItem())
	copy(pq.Q[n+1:pq.Size+1], pq.Q[n:pq.Size])
	pq.Q[n] = item
	pq.Size++
	return pq.Size
}

// Pop returns the Item with the highest priorty
func (pq *FloatPQ) Pop() (FloatItem, error) {
	if pq.IsEmpty() || pq.Q[0] == newFloatItem() {
		return newFloatItem(), errors.New("queue is empty")
	}
	p := pq.Q[0]
	pq.Q = pq.Q[1:]
	pq.Size--
	return FloatItem{Val: p.Val, Pri: p.Pri}, nil
}
