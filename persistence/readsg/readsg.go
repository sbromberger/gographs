package staticgraph

import (
	"log"

	"github.com/sbromberger/gographs"
	"gonum.org/v1/hdf5"
)

// ReadStaticGraph returns a graph read from a LightGraphs.StaticGraph.
func ReadStaticGraph(fn string) gographs.Graph {
	ds1name := "f_vec"
	ds2name := "f_ind"
	f, err := hdf5.OpenFile(fn, hdf5.F_ACC_RDONLY)
	defer f.Close()
	if err != nil {
		log.Fatal("Cannot open file: ", err)
	}
	dset1, err := f.OpenDataset(ds1name)
	if err != nil {
		log.Fatal("Cannot open dataset ds1: ", err)
	}

	dset2, err := f.OpenDataset(ds2name)
	if err != nil {
		log.Fatal("Cannot open dataset ds2: ", err)
	}

	s1 := dset1.Space()
	s2 := dset2.Space()

	fVec := make([]uint32, s1.SimpleExtentNPoints())
	fInd32 := make([]uint32, s2.SimpleExtentNPoints())
	fInd := make([]uint64, s2.SimpleExtentNPoints())

	if err := dset1.Read(&fVec); err != nil {
		log.Fatal("can't read f_vec: ", err)
	}
	if err := dset2.Read(&fInd32); err != nil {
		log.Fatal("can't read f_ind: ", err)
	}

	for i := range fVec {
		fVec[i]--
	}

	for i := range fInd32 {
		fInd[i] = uint64(fInd32[i]) - 1
	}

	return gographs.MakeGraph(fVec, fInd)
}
