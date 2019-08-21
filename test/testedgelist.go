package main

import (
	"fmt"
	"os"

	"github.com/sbromberger/gographs"
)

func main() {
	fn := os.Args[1]
	persistene.GraphFromEdgeList(fn)
	h := gographs.HouseGraph()
	fmt.Println("h.Fadj = ", h.Fadj())

}
