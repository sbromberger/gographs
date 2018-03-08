package priorityqueue

import (
	"github.com/sbromberger/gographs/heap"
)

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*heap.Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) IsEmpty() bool { return pq.Len() == 0 }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(item *heap.Item) {
	n := len(*pq)
	// item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() *heap.Item {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the Priority and Value of an Item in the queue.
func (pq *PriorityQueue) update(item *heap.Item, Value uint32, Priority float32) {
	item.Value = Value
	item.Priority = Priority
	heap.Fix(pq, item.Index)
}
