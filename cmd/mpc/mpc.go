package main

import (
	"flag"
	"log"
)

var circuit int
var prime int

func init() {
	flag.IntVar(&circuit, "circuit", 1, "Circuit to run")
	flag.IntVar(&prime, "prime", 101, "Prime number to use for modular arithmetic")
}

func main() {
	flag.Parse()
	log.Printf("Running circuit %d", circuit)
}
