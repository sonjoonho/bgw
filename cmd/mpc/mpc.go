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

	nParties := cfg.Circuit.NParties

	// Initialise each party.
	parties := make([]*party.Party, nParties, nParties)
	for i := 0; i < nParties; i++ {
		// Note that cfg.Circuit is copied, and the rest of the parameters are values so parties do not share memory.
		p := party.New(i, cfg.Secrets[i], cfg.Circuit.Copy(), cfg.Field)
		parties[i] = p
	}

	results := make([]int, nParties, nParties)
	// Start protocol.
	var wg sync.WaitGroup
	for i, p := range parties {
		p.SubscribeAll(parties)
		wg.Add(1)
		go runParty(i, results, p, &wg)
	}

	// Block until all parties have finished.
	wg.Wait()

	for _, r := range results {
		if r != results[0] {
			logger.Fatalf("Protcol failed: return values do not match.")
		}
	}
	logger.Printf("Output: %d", results[0])
}

func runParty(i int, results []int, p *party.Party, wg *sync.WaitGroup) {
	results[i] = p.Run(wg)
}
