package persistence

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sbromberger/gographs"
)

// ReadEdgeList returns a sorted list of edges given a scanner.
func ReadEdgeList(scanner *bufio.Scanner, offset int) []gographs.Edge {
	var l string
	edges := make([]gographs.Edge, 0)
	for scanner.Scan() {
		l = scanner.Text()
		pieces := strings.Split(l, ",")
		if len(pieces) > 2 {
			log.Fatal("Parsing err: got ", l)
		}
		u64, err := strconv.ParseUint(pieces[0], 10, 32)
		if err != nil {
			log.Fatal("Parsing err: ", err)
		}
		v64, err := strconv.ParseUint(pieces[1], 10, 32)
		if err != nil {
			log.Fatal("Parsing err: ", err)
		}

		u32 := uint32(u64) - uint32(offset)
		v32 := uint32(v64) - uint32(offset)
		edges = append(edges, gographs.Edge{Src: u32, Dst: v32})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Other error: ", err)
	}
	revedges := make([]gographs.Edge, len(edges))
	for i, e := range edges {
		revedges[i] = gographs.Reverse(e)
	}
	edges = append(edges, revedges...)
	sort.Sort(gographs.EdgeList(edges))

	// fmt.Println("edges = ", edges)

	return edges
}

// GraphFromEdgeList returns a graph from an edgelist.
func GraphFromEdgeList(fn string) gographs.Graph {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal("Open failed: ", err)
	}

	scanner := bufio.NewScanner(f)

	edges := ReadEdgeList(scanner, 0)
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
