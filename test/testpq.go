package main

import (
	"fmt"

	"github.com/sbromberger/gographs/priorityqueue"
)

func main() {
	pq := priorityqueue.NewFloatPQ(10)
	pq.Push(5, 2.0)
	// fmt.Println("pq is ", pq)
	pq.Push(2, 1.1)
	// fmt.Println("pq is ", pq)
	pq.Push(8, 4.3)
	// fmt.Println("pq is ", pq)
	pq.Push(9, 2.2)
	// fmt.Println("pq is ", pq)
	pq.Push(15, 2.0)

	for i := 0; i < 6; i++ {
		// fmt.Println("pq is ", pq)
		if item, err := pq.Pop(); err != nil {
			fmt.Println("error received: ", err)
		} else {
			fmt.Printf("popped val = %d with pri = %f\n", item.Val, item.Pri)
		}
	}

}
