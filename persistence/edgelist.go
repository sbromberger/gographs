package persistence

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sbromberger/graph"
)

func readEdgeList(scanner *bufio.Scanner, offset uint32) ([]uint32, []uint32, error) {
	var l string
	ss := make([]uint32, 0, 100)
	ds := make([]uint32, 0, 100)
	for scanner.Scan() {
		l = scanner.Text()
		pieces := strings.Split(l, ",")
		if len(pieces) > 2 {
			return []uint32{}, []uint32{}, fmt.Errorf("Parsing error: got %s", l)
		}
		u64, err := strconv.ParseUint(pieces[0], 10, 32)
		if err != nil {
			return []uint32{}, []uint32{}, fmt.Errorf("Parsing error: got %s", l)
		}
		v64, err := strconv.ParseUint(pieces[1], 10, 32)
		if err != nil {
			return []uint32{}, []uint32{}, fmt.Errorf("Parsing error: got %s", l)
		}
		u := uint32(u64) - offset
		v := uint32(v64) - offset
		ss = append(ss, u)
		ds = append(ds, v)
	}
	if err := scanner.Err(); err != nil {
		return []uint32{}, []uint32{}, fmt.Errorf("Other error: %v", err)
	}
	return ss, ds, nil
}

// ReadEdgeList returns a graph from an edgelist.
func ReadEdgeList(fn string) (graph.Graph, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if err != nil {
		return graph.Graph{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ss, ds, err := readEdgeList(scanner, 0)
	if err != nil {
		return graph.Graph{}, err
	}

	return graph.New(ss, ds)
}
