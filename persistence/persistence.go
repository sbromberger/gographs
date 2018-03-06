package persistence

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sbromberger/gographs"
	"gonum.org/v1/hdf5"
)

func ReadText(fn string) gographs.Graph {
	var fInd, fVec []uint32
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal("Cannot open file: ", err)
	}
	scanner := bufio.NewScanner(f)
	l := ""
	inFInd := true
	for scanner.Scan() {
		l = scanner.Text()
		if l == "-----" {
			inFInd = false
			continue
		}
		v64, err := strconv.ParseUint(l, 10, 32)

		if err != nil {
			log.Fatal("Parsing err: ", err)
		}
		v := uint32(v64 - 1)
		if inFInd {
			fInd = append(fInd, v)
		} else {
			fVec = append(fVec, v)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Other error: ", err)
	}

	return gographs.MakeGraph(fVec, fInd)
}

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
