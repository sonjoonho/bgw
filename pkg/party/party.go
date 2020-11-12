// Package party implements the representation of a party in the BGW protocol.
package party

import (
	"bgw/pkg/circuit"
	"bgw/pkg/field"
	"bgw/pkg/poly"
	"fmt"
	"log"
	"os"
	"sync"
)

type message struct {
	party int
	gate  int
	share int
}

// Party is a party which can communicate with other parties.
type Party struct {
	// id is the index of this Party.
	id int
	// secret is...
	secret int
	// ch is a channel through which this Party receives messages.
	ch chan *message
	// done is...
	done chan bool
	// subs is a slice of send-only channels that this party uses to send message to subscribers. Its capacity is equal
	// to the number of parties, specified during initialisation.
	subs []chan<- *message
	// shares is a buffer for received shares. It maps from Party id to gate.Gate id to share. The first column contains
	// the initial share for each party.
	shares  [][]int
	field   field.Field
	circuit *circuit.Circuit
	logger  *log.Logger
}

// New initialises and returns a new Party. nParties specifies the number of parties participating in the protocol,
// and nGates specifies the number of gates in the circuit.
func New(id int, secret int, circuit *circuit.Circuit, field field.Field) *Party {
	nParties := circuit.NParties()
	nGates := circuit.NGates()

	p := &Party{
		id:      id,
		secret:  secret,
		circuit: circuit,
		field:   field,
		ch:      make(chan *message, nParties),
		subs:    make([]chan<- *message, nParties, nParties),
		shares:  make([][]int, nParties, nParties),
		logger:  log.New(os.Stderr, fmt.Sprintf("Party %d: ", id), 0),
	}

	// Initialise nested slices of shares.
	for i := 0; i < nParties; i++ {
		p.shares[i] = make([]int, nGates+1)
	}

	return p
}

// Id returns the id for this Party.
func (p *Party) Id() int {
	return p.id
}

// SubscribeAll subscribes this Party to all parties.
func (p *Party) SubscribeAll(parties []*Party) {
	for _, pty := range parties {
		p.subs[pty.id] = pty.ch
	}
}

// SendShare sends the specified share to another Party.
func (p *Party) SendShare(to int, share int, srcGate int) {
	p.logger.Printf("Sending share %d at gate %d to party %d\n", share, to, srcGate)
	ch := p.subs[to]
	msg := &message{party: p.id, gate: srcGate, share: share}
	ch <- msg
}

// RecvShare receives a share
func (p *Party) RecvShare() *message {
	msg := <-p.ch
	p.logger.Printf("Received message %+v", msg)
	return msg
}

// Run runs the BGW protocol for this party.
func (p *Party) Run(wg *sync.WaitGroup) {
	p.logger.Println("Running...")

	nParties := p.circuit.NParties()
	nGates := p.circuit.NGates()

	// 1. Split secret and share.
	//	  Each party generates a random polynomial with the 0th degree term as the secret. This polynomial needs to be
	//    evaluated for each party to generate shares i.e. P(i) for i in 1..number of parties. Then broadcast shares to
	//    every other party.
	coeffs := make([]int, nParties, nParties)
	coeffs[0] = p.secret
	for d := 1; d < nParties; d++ {
		coeffs[d] = p.field.Rand()
	}
	po := poly.New(coeffs, p.field)

	for i := 1; i <= nParties; i++ {
		share := po.Eval(i)
		p.SendShare(i-1, share, 0)
	}

	// Ensure we have received all initial shares.
	gateCounts := make([]int, nGates+1, nGates+1)
	for gateCounts[0] < nParties {
		msg := p.RecvShare()
		p.shares[msg.party][msg.gate] = msg.share
		gateCounts[msg.gate]++
	}

	// 2. Run circuit.
	//    We iterate over all the gates in the circuit and wait for the inputs from the source gate, and then run the
	//    gate computation on the shares that we have. Then broadcast those shares.

	// 3. Create final result.

	wg.Done()
}
