package main

import (
	"os"

	"github.com/sbromberger/gographs/persistence"
)

func main() {
	fn := os.Args[1]
	fn2 := os.Args[2]
	h := persistence.ReadText(fn)
	f := h.Fadj()
	raw.SaveRaw(fn2, f.Rowidx, f.Colptr)
}
