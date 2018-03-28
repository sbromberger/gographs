package main

import (
	"fmt"
	"time"

	"github.com/sbromberger/gographs"
)

func main() {
	g := gographs.HouseGraph()
	fmt.Println("Order(g) = ", g.Order())
	fmt.Println("Size(g) = ", g.Size())
	for i := 0; i < g.Order(); i++ {
		fmt.Printf("Neighbors(%d) = %v\n", i, g.OutNeighbors(uint32(i)))
	}
	start := time.Now()
	d := gographs.Dijkstra(&g, 0, true)
	elapsed := time.Since(start)
	fmt.Print("BFS done: ")
	fmt.Println(elapsed)
	fmt.Println(d)

}
