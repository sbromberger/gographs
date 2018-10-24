package main

import (
	"os"

	"github.com/sbromberger/gographs/persistence/lgformat"
	"github.com/sbromberger/gographs/persistence/raw"
)

func main() {
	fn1 := os.Args[1]
	fn2 := os.Args[2]
	g := lgformat.GraphFromLG(fn1)
	raw.SaveRaw(fn2, g)
}
