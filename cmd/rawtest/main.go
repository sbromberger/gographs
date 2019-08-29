package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sbromberger/graph"
	"github.com/sbromberger/graph/persistence"
)

func main() {
	g, err := persistence.ReadRaw(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(g)

	fmt.Println("outneighbors(0) = ", g.OutNeighbors(0))
	start := time.Now()
	v, w := graph.BFS(g, 0)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("elapsed BFS = ", elapsed, ", len(v) = ", len(v), ", len(w) = ", len(w))
	fmt.Println(v[:10])
	fmt.Println(w[:10])
	s := int(0)
	for _, n := range w {
		s += int(n)
	}
	fmt.Println(s)
}
