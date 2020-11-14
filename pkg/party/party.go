// Package party implements the representation of a party in the BGW protocol.
package party

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/circuit"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/poly"
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
	// id is the identifier of this Party. It starts from 0.
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
	nParties := circuit.NParties
	nGates := circuit.NGates

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
		pShares := make([]int, nGates+1)
		for j := 0; j < nGates+1; j++ {
			pShares[j] = -1
		}
		p.shares[i] = pShares
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

	nParties := p.circuit.NParties
	nGates := p.circuit.NGates

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
	gates := traverse(p.circuit)
	for _, g := range gates {
		p.logger.Printf("Processing gate %#v...\n", g)
		switch v := g.(type) {
		case *gate.Input:
			v.Share = p.shares[v.Party][0]
			p.logger.Printf("Input: %d\n", v.Output())
		case *gate.Add:
			p.logger.Printf("Add: %d\n", v.Output())
		}
	}

	// 3. Create final result.

	wg.Done()
}

// traverse performs an iterative post-order traversal of the circuit's gates.
func traverse(circuit *circuit.Circuit) []gate.Gate {
	stack := []gate.Gate{circuit.Root}
	var res []gate.Gate

	for len(stack) > 0 {
		var next gate.Gate
		stack, next = pop(stack)
		res = append([]gate.Gate{next}, res...)
		if next.First() != nil {
			stack = append(stack, next.First())
		}
		if next.Second() != nil {
			stack = append(stack, next.Second())
		}
	}

	return res
}

func peek(stack []gate.Gate) gate.Gate {
	return stack[len(stack)-1]
}

func pop(stack []gate.Gate) ([]gate.Gate, gate.Gate) {
	g := peek(stack)
	return stack[0 : len(stack)-1], g
}
