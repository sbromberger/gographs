package main

import (
	"fmt"

	"github.com/sbromberger/gographs/priorityqueue"
)

func main() {
	pq := priorityqueue.New(10)
	pq.Push(5, 2)
	pq.Push(2, 1)
	pq.Push(8, 4)
	pq.Push(9, 2)

	for i := 0; i < 5; i++ {
		fmt.Println("pq is ", pq)
		if item, err := pq.Pop(); err != nil {
			fmt.Println("error received: ", err)
		} else {
			fmt.Printf("popped val = %d with pri = %d\n", item.Value, item.Priority)
		}
	}

}
