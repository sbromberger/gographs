package converter

import (
	"bufio"
	"fmt"
	"os"

	graph "github.com/sbromberger/gographs"
)

// GraphFromLG reads a graph from lightgraphs format. Ignores header.
func GraphFromLG(fn string) (graph.SimpleGraph, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if err != nil {
		return graph.SimpleGraph{}, fmt.Errorf("Open failed: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan() // read header

	ss, ds, err := readEdgeList(scanner, 1) // offset is 1 because Julia is 1-indexed
	if err != nil {
		return graph.SimpleGraph{}, err
	}
	return graph.New(ss, ds)
}
