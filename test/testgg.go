package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/sbromberger/gographs/persistence"

	"github.com/sbromberger/gographs"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var fn = flag.String("f", "", "filename to open")

func sum(a []int) int {
	s := 0
	for _, r := range a {
		s += r
	}
	return s
}

func main() {
	ggg := 94
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		fmt.Println("created cpu profile: ", *cpuprofile)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	//      0  3  6
	// 0  [ T  F
	// 1    T  F
	// 2    T  T
	// 3    F  T
	// 4    F  T     ]

	// HouseGraph
	// rind := []uint32{1, 2, 0, 3, 0, 3, 4, 1, 2, 4, 2, 3}
	// cptr := []uint32{0, 2, 4, 7, 10, 12}
	// g := gographs.MakeGraph(rind, cptr)

	// fmt.Println("g created")
	// fmt.Println("has_edge (g, 0, 1) = ", g.HasEdge(0, 1))
	// fmt.Println("has_edge (g, 1, 0) = ", g.HasEdge(1, 0))
	// fmt.Println("has_edge (g, 1, 2) = ", g.HasEdge(1, 2))

	// fmt.Println("edges(g) = ", g.Edges())
	// fmt.Println("outneighbors(g, 2) = ", g.OutNeighborsInt(2))

	h := persistence.ReadStaticGraph(*fn)
	// fmt.Println("edges(h) = ", h.Edges())
	// fmt.Println("Fadj(h) = ", h.Fadj())
	fmt.Println("Order(h) = ", h.Order())
	fmt.Println("Size(h) = ", h.Size())
	start := time.Now()
	gographs.Dijkstra(&h, 0)
	elapsed := time.Since(start)
	fmt.Println("dijkstra done:")
	fmt.Println(elapsed)

	// fmt.Println("doing all dijkstras")
	// start = time.Now()
	// ds := gographs.AllDijkstraShortestPaths(&h)
	// elapsed = time.Since(start)
	// fmt.Println("elapsed = ", elapsed)
	// fmt.Println("len(ds) = ", len(ds))
	// // fmt.Println("len(ds[5999].Dists) = ", len(ds[5999].Dists))
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		fmt.Println("made mem profile", *memprofile)
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

}

// fmt.Println("d.dists[10:20] = ", d.Dists)
// fmt.Println("d.parents[20:30] = ", d.Parents)
// fmt.Println("d.preds[30:40] = ", d.Predecessors)
// fmt.Println("d.pathcounts[50:60] = ", d.Pathcounts)
// fmt.Println("sum d.dists = ", sum(d.Dists))
// fmt.Println("sum d.pathcounts = ", sum(d.Pathcounts))
// z := sparsevecs.UInt32SparseVec{v1, v2}
// fmt.Println("z.rowidx = ", z.Rowidx)
// fmt.Println("test1: ", z.GetIndexInt(3, 0)) // F
// fmt.Println("test2: ", z.GetIndexInt(2, 1)) // T
// fmt.Println("test3: ", z.GetIndexInt(1, 0)) // T
// fmt.Println("-----------------------------------------")
// fmt.Println("ok, testing insert")
// fmt.Println("test4:")
// z.Insert(2, 1)
// //      0  3  6
// // 0  [ T  F
// // 1    T  F
// // 2    T  T
// // 3    F  T
// // 4    F  T     ]
// fmt.Println("-----------------------------------------")
// fmt.Println("test5: ", z)
// z.Insert(2, 2)
// //      0  3  6  7
// // 0  [ T  F  F
// // 1    T  F  F
// // 2    T  T  T
// // 3    F  T  F
// // 4    F  T  F    ]
// // 0 1 2 / 2 3 4 / 2
// fmt.Println("test6: ", z)
