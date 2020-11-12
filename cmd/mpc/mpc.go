package main

import (
	"bgw/pkg/circuit"
	"bgw/pkg/field"
	"bgw/pkg/gate"
	"bgw/pkg/party"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

var logger = log.New(os.Stderr, "MPC: ", 0)

const (
	defaultCircuitNumber = 1
	defaultPrime         = 101
	defaultNParties      = 2
	defaultSeed          = 0
)

var (
	circuitNumber int
	prime         int
	nParties      int
	seed          int64
)

func init() {
	flag.IntVar(&circuitNumber, "circuit", defaultCircuitNumber, "Circuit to run")
	flag.IntVar(&prime, "prime", defaultPrime, "Prime number to use for modular arithmetic")
	flag.IntVar(&nParties, "parties", defaultNParties, "Number of parties")
	flag.Int64Var(&seed, "seed", defaultSeed, "Seed for pseudorandom number generation")
}

func main() {
	flag.Parse()
	logger.Println("Starting BGW protocol...")

	if seed == defaultSeed {
		seed = time.Now().UnixNano()
	}
	fld := field.New(prime, seed)

	// TODO(sonjoonho): Read circuit definition from configuration file.
	// TODO(sonjoonho): Assert len(inputs) == len(secrets) == nParties.
	crct := circuit.New(
		[]int{0, 0},
		[]gate.Gate{
			gate.NewAdd(1, fld),
		},
	)
	secrets := []int{5, 28}

	logger.Printf("Running circuit %d with %d parties", circuitNumber, nParties)

	// Initialise each party.
	parties := make([]*party.Party, nParties, nParties)
	for i := 0; i < nParties; i++ {
		p := party.New(i, secrets[i], crct, fld)
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
