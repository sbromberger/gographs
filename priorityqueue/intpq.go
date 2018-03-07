package priorityqueue

import (
	"errors"
	"log"
)

// IntItem stores value and priority.
type IntItem struct {
	Value    uint32
	Priority int
}

// IntPQ implements a priority queue.
type IntPQ struct {
	q                       [][]uint32
	numpri, lowind, highind int
}

// IsEmpty returns true if the queue is empty.
func (pq *IntPQ) IsEmpty() bool {
	return (pq.lowind > pq.highind) ||
		(pq.lowind == pq.highind && len(pq.q[pq.lowind]) == 0)
}

// lenPriInRange will return the length of the pq at priority n if
// lowind <= n <= highind, or -1.
func (pq *IntPQ) lenPriInRange(n int) int {
	if pq.lowind <= n && n <= pq.highind {
		return len(pq.q[n])
	}
	return -1
}

// Push adds a value/priority to a PQ
func (pq *IntPQ) Push(val uint32, pri int) {
	atpri := append((*pq).q[pri], val)
	pq.q[pri] = atpri
	if pri < pq.lowind {
		pq.lowind = pri
	}
	if pri > pq.highind {
		pq.highind = pri
	}
}

// Pop retrieves the highest-priority IntIntItem from the PQ
func (pq *IntPQ) Pop() (IntItem, error) {

	for pq.lenPriInRange(pq.lowind) == 0 && pq.lowind < pq.numpri {
		pq.lowind++
	}

	if pq.lowind > pq.highind {
		return IntItem{}, errors.New("queue is empty")
	}

	pri := pq.lowind
	atpri := pq.q[pri]
	lenAtpri := len(atpri)

	if pq.lowind == pq.numpri && lenAtpri == 0 {
		return IntItem{}, errors.New("queue is empty")
	}

	val := (atpri)[lenAtpri-1]
	atpri = (atpri)[:lenAtpri-1]
	pq.q[pri] = atpri

	return IntItem{val, pri}, nil
}

// NewIntPQ creates a PQ with `n` priorities.
func NewIntPQ(n int) IntPQ {

	if n == 0 {
		log.Fatal("priority queue must have more than zero priority levels")
	}
	q := make([][]uint32, n)
	return IntPQ{q, n, n, 0}
}
