package main

import (
	"fmt"
	"sort"
)

type Pair struct {
	src int
	dst int
}

type Pairs []Pair

func (p Pairs) Len() int {
	return len(p)
}

func (p Pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Pairs) Less(i, j int) bool {
	if p[i].src < p[j].src {
		return true
	}
	if p[i].src > p[j].src {
		return false
	}
	return p[i].dst < p[j].dst
}

func main() {
	a := make([]Pair, 0)
	a = append(a, Pair{1, 2})
	a = append(a, Pair{2, 20})
	a = append(a, Pair{2, 10})
	a = append(a, Pair{1, 60})

	sort.Sort(Pairs(a))
	fmt.Println(a)
}
