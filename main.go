package main

import (
	"log"
	"rangetree"
)

func main() {
	rt := rangetree.NewRangeTree()

	rt.AddRange(5, 15)
	rt.AddRange(14, 20)
	rt.AddRange(16, 20)
	rt.AddRange(30, 50)

	log.Printf("addrange complete\n")
	rt.Dump(false)
}
