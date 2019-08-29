package main

import (
	"log"
	"os"

	"github.com/sbromberger/graph/persistence"
)

func main() {
	fn1 := os.Args[1]
	fn2 := os.Args[2]
	g, err := persistence.GraphFromLG(fn1)
	if err != nil {
		log.Fatal(err)
	}
	if err := persistence.SaveRaw(fn2, g); err != nil {
		log.Fatal(err)
	}
}
