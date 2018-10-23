package main

import (
	"fmt"
	"os"

	"github.com/sbromberger/gographs"
	"github.com/sbromberger/gographs/persistence/edgelist"
)

func main() {
	fn := os.Args[1]
	edgelist.GraphFromEdgeList(fn)
	h := gographs.HouseGraph()
	fmt.Println("h.Fadj = ", h.Fadj())

}
