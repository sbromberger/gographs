package readtext

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"github.com/sbromberger/gographs"
)

func ReadText(fn string) gographs.Graph {
	var fInd []uint64
	var fVec []uint32
	f, err := os.Open(fn)
	defer f.Close()
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
		v := v64 - 1
		if inFInd {
			fInd = append(fInd, v)
		} else {
			fVec = append(fVec, uint32(v))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Other error: ", err)
	}

	return gographs.MakeGraph(fVec, fInd)
}
