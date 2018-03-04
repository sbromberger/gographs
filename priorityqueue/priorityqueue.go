package priorityqueue

import (
	"errors"
	"log"
)

// Item stores value and priority.
type Item struct {
	Value    uint32
	Priority int
}

// PriorityQueue implements a priority queue.
type PriorityQueue struct {
	q                       [][]uint32
	numpri, lowind, highind int
}

// IsEmpty returns true if the queue is empty.
func (pq *PriorityQueue) IsEmpty() bool {
	return (pq.lowind > pq.highind) ||
		(pq.lowind == pq.highind && len(pq.q[pq.lowind]) == 0)
}

// lenPriInRange will return the length of the pq at priority n if
// lowind <= n <= highind, or -1.
func (pq *PriorityQueue) lenPriInRange(n int) int {
	if pq.lowind <= n && n <= pq.highind {
		return len(pq.q[n])
	}
	return -1
}

// Push adds a value/priority to a PriorityQueue
func (pq *PriorityQueue) Push(val uint32, pri int) {
	atpri := append((*pq).q[pri], val)
	pq.q[pri] = atpri
	if pri < pq.lowind {
		pq.lowind = pri
	}
	if pri > pq.highind {
		pq.highind = pri
	}
}

// Pop retrieves the highest-priority item from the PriorityQueue
func (pq *PriorityQueue) Pop() (Item, error) {

	for pq.lenPriInRange(pq.lowind) == 0 && pq.lowind < pq.numpri {
		pq.lowind++
	}

	if pq.lowind > pq.highind {
		return Item{}, errors.New("queue is empty")
	}

	pri := pq.lowind
	atpri := pq.q[pri]
	lenAtpri := len(atpri)

	if pq.lowind == pq.numpri && lenAtpri == 0 {
		return Item{}, errors.New("queue is empty")
	}

	val := (atpri)[lenAtpri-1]
	atpri = (atpri)[:lenAtpri-1]
	pq.q[pri] = atpri

	return Item{val, pri}, nil
}

// New creates a PriorityQueue with `n` priorities.
func New(n int) PriorityQueue {

	if n == 0 {
		log.Fatal("priority queue must have more than zero priority levels")
	}
	q := make([][]uint32, n)
	return PriorityQueue{q, n, n, 0}
}
