// Package main contains the main driver code for the protocol.
package main

import (
	"flag"
	"fmt"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/config"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/party"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "MPC: ", log.Lmicroseconds)

const (
	defaultCircuitNumber = 1
	defaultPrime         = 101
	defaultSeed          = 0
	defaultDegree        = -1
)

var (
	circuitNumber int
	degree        int
	prime         int
	seed          int64
)

func init() {
	flag.IntVar(&circuitNumber, "circuit", defaultCircuitNumber, "Circuit to run.")
	flag.IntVar(&degree, "degree", defaultDegree, "Degree of polynomial. If unset, it is set to N-1/2")
	flag.IntVar(&prime, "prime", defaultPrime, "Prime number to use for modular arithmetic.")
	flag.Int64Var(&seed, "seed", defaultSeed, "Seed for pseudorandom number generation. If unset, the current time is used.")
}

func main() {
	flag.Parse()
	logger.Println("Starting BGW protocol...")

	cfg, err := config.New(prime, seed, defaultSeed, degree, defaultDegree, circuitNumber)
	if err != nil {
		logger.Fatalf("Configuration failed: %v", err)
	}

	nParties := cfg.Circuit.NParties

	logger.Println("")
	logger.Printf("Circuit Configuration")
	logger.Println("===================================")
	logger.Printf("  Circuit number:    %d", circuitNumber)
	logger.Printf("  Number of parties: %d", nParties)
	logger.Printf("  Secrets:           %v", cfg.Secrets)
	logger.Printf("  Polynomial degree: %d", cfg.Degree)
	logger.Println("")

	actual, err := RunProtocol(cfg)
	if err != nil {
		logger.Fatalf("Protocol failed: %v", err)
	}

	expected := cfg.Field.Mod(cfg.Circuit.ComputeExpected(cfg.Secrets))
	logger.Printf("Expected output: %d", expected)
	logger.Printf("Actual output:   %d", actual)

	if expected == actual {
		logger.Println("Protocol succeeded (:")
	} else {
		logger.Fatal("Protocol failed ):")
	}
}

// RunProtocol runs the BGW protocol using the provided configuration.
func RunProtocol(cfg *config.Config) (int, error) {
	nParties := cfg.Circuit.NParties

	// Initialise each party.
	parties := make([]*party.Party, nParties, nParties)
	for i := 0; i < nParties; i++ {
		// Note that cfg.Circuit is copied, and the rest of the parameters are values so parties do not share memory.
		p := party.New(i, cfg.Secrets[i], cfg.Circuit.Copy(), cfg.Field, cfg.Degree)
		parties[i] = p
	}

	// results stores the final output values of each party. These are then checked for consistency.
	results := make([]int, nParties, nParties)
	// Go!
	var wg sync.WaitGroup
	for i, p := range parties {
		p.SubscribeAll(parties)
		wg.Add(1)
		// These parameters are important!
		// Reference: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables.
		go func(i int, p *party.Party) {
			defer wg.Done()
			results[i] = p.Run()
		}(i, p)
	}

	// Block until all parties have finished.
	wg.Wait()

	// Check results for consistency.
	for _, r := range results {
		if r != results[0] {
			return 0, fmt.Errorf("protocol failed: return values do not match")
		}
	}

	return results[0], nil
}
