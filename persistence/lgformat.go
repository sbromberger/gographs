package persistence

import (
	"bufio"
	"log"
	"os"

	"github.com/sbromberger/gographs"
)

// GraphFromLG reads a graph from lightgraphs format. Ignores header.
func GraphFromLG(fn string) gographs.Graph {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal("Open failed: ", err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan() // read header

	edges := ReadEdgeList(scanner, 1)
	rowval := make([]uint32, len(edges))
	colptr := make([]uint64, 0)
	currsrc := uint32(0)

	for i, e := range edges {
		src := e.Src
		for currsrc <= src {
			colptr = append(colptr, uint64(i))
			currsrc++
		}

		rowval[i] = e.Dst
	}
	colptr = append(colptr, uint64(len(rowval)))
	// fmt.Println("rowval = ", rowval)
	// fmt.Println("colptr = ", colptr)
	return gographs.MakeGraph(rowval, colptr)
}
