package main

import (
	"fmt"

	"github.com/sbromberger/gographs/persistence"
)

func main() {
	fmt.Println("Starting persistence")
	g := persistence.ReadStaticGraph("uint32.sg")
	fmt.Printf("g has order %d and size %d\n", g.Order(), g.Size())
	fmt.Println("edges(g) = ", g.Edges())
}
