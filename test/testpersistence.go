package main

import (
	"fmt"

	"github.com/sbromberger/gographs/persistence"
)

func main() {
	fmt.Println("Starting persistence")
	g := persistence.ReadText("sg-10k-250k.txt")
	fmt.Printf("g has order %d and size %d\n", g.Order(), g.Size())
	// fmt.Println("edges(g) = ", g.Edges())
}
