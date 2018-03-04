package persistence

import (
	"fmt"
	"log"

	"github.com/sbromberger/gographs"
	"gonum.org/v1/hdf5"
)

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

	fmt.Println("dset1 = ", dset1)
	dset2, err := f.OpenDataset(ds2name)
	if err != nil {
		log.Fatal("Cannot open dataset ds2: ", err)
	}

	s1 := dset1.Space()
	s2 := dset2.Space()

	fVec := make([]uint32, s1.SimpleExtentNPoints())
	fInd := make([]uint32, s2.SimpleExtentNPoints())

	if err := dset1.Read(&fVec); err != nil {
		log.Fatal("can't read f_vec: ", err)
	}
	fmt.Println("dset1.datatypw = ")
	if err := dset2.Read(&fInd); err != nil {
		log.Fatal("can't read f_ind: ", err)
	}

	for i := range fVec {
		fVec[i]--
	}

	for i := range fInd {
		fInd[i]--
	}

	return gographs.MakeGraph(fVec, fInd)
}
