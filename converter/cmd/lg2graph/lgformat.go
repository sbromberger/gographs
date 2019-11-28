package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sbromberger/gographs/converter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], " infile outfile")
		os.Exit(1)
	}
	in := os.Args[1]
	out := os.Args[2]

	g, err := converter.GraphFromLG(in)
	if err != nil {
		log.Fatalf("Can't read %s: %v", in, err)
	}
	err = g.Save(out)
	if err != nil {
		log.Fatalf("Can't write %s: %v", out, err)
	}
}
