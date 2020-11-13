package main

import (
	"flag"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/config"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/party"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "MPC: ", 0)

const (
	defaultCircuitNumber = 1
	defaultPrime         = 101
	defaultSeed          = 0
)

var (
	circuitNumber int
	prime         int
	seed          int64
)

func init() {
	flag.IntVar(&circuitNumber, "circuit", defaultCircuitNumber, "Circuit to run")
	flag.IntVar(&prime, "prime", defaultPrime, "Prime number to use for modular arithmetic")
	flag.Int64Var(&seed, "seed", defaultSeed, "Seed for pseudorandom number generation")
}

func main() {
	flag.Parse()
	logger.Println("Starting BGW protocol...")

	cfg := config.New(prime, seed, defaultSeed, circuitNumber)

	logger.Printf("Running circuit %d: %+v", circuitNumber, cfg)

	nParties := cfg.Circuit.NParties()

	// Initialise each party.
	parties := make([]*party.Party, nParties, nParties)
	for i := 0; i < nParties; i++ {
		// Note that cfg.Circuit is copied, and the rest of the parameters are values so parties do not share memory.
		p := party.New(i, cfg.Secrets[i], cfg.Circuit.Copy(), cfg.Field)
		parties[i] = p
	}

	// Start protocol.
	var wg sync.WaitGroup
	for _, p := range parties {
		p.SubscribeAll(parties)
		wg.Add(1)
		go p.Run(&wg)
	}

	// Block until all parties have finished.
	wg.Wait()
}
