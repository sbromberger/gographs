package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/sbromberger/gographs"
	"github.com/sbromberger/gographs/persistence/raw"
	"github.com/sbromberger/gographs/persistence/readsg"
	"github.com/sbromberger/gographs/persistence/readtext"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var fn = flag.String("f", "", "filename to open")
var src = flag.Int("v", 0, "source vertex")
var procs = flag.Int("procs", 0, "number of procs to use")
var ff = flag.String("fmt", "t", "file format: (t)ext, (m)map, (s)taticgraph")
var parallel = flag.Bool("p", false, "run in parallel")

func sum(a []int) int {
	s := 0
	for _, r := range a {
		s += r
	}
	return s
}

func main() {
	flag.Parse()

	fmt.Println("reading graph")
	if *procs == 0 {
		*procs = runtime.NumCPU()
	} else if *procs > runtime.NumCPU() {
		*procs = runtime.NumCPU()
	}
	fmt.Println("Procs = ", *procs)

	// h := readtext.ReadText(*fn)
	var h gographs.Graph
	switch *ff {
	case "t":
		h = readtext.ReadText(*fn)
	case "m":
		h = raw.GraphFromRaw(*fn)
	case "s":
		h = staticgraph.ReadStaticGraph(*fn)
	default:
		log.Fatal("Unknown file format (use one of t, m, or s)")
	}
	fmt.Println("Order(h) = ", h.Order())
	fmt.Println("Size(h) = ", h.Size())
	fmt.Println("Neighbors(1) = ", h.OutNeighbors(1))
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
	runtime.GC()
	debug.SetGCPercent(-1)
	// runtime.LockOSThread()
	times := make([]time.Duration, 10)
	for i := range times {
		start := time.Now()
		if *parallel {
			gographs.ParallelBFS(&h, uint32(*src), *procs)
		} else {
			gographs.BFS(&h, uint32(*src))
		}
		elapsed := time.Since(start)
		fmt.Print("BFS done: ")
		fmt.Println(elapsed)
		times[i] = elapsed
	}
	sumTime := time.Duration(0)
	for i := range times {
		sumTime += times[i]
	}
	avgintns := int64(sumTime/time.Nanosecond) / int64(len(times))
	avg := time.Duration(time.Nanosecond * time.Duration(avgintns))
	fmt.Printf("Average for %d runs: %s\n", len(times), avg)
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
