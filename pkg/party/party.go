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
	// shares is a buffer for received shares. It maps from Party id to gate.Gate to share.
	shares  []map[gate.Gate]*int
	field   field.Field
	circuit *circuit.Circuit
	logger  *log.Logger
}

// New initialises and returns a new Party. nParties specifies the number of parties participating in the protocol,
// and nGates specifies the number of gates in the circuit.
func New(id int, secret int, circuit *circuit.Circuit, field field.Field) *Party {
	nParties := circuit.NParties

	p := &Party{
		id:      id,
		secret:  secret,
		circuit: circuit,
		field:   field,
		ch:      make(chan *message, nParties+1),
		subs:    make([]chan<- *message, nParties, nParties),
		shares:  make([]map[gate.Gate]*int, nParties, nParties),
		logger:  log.New(os.Stderr, fmt.Sprintf("Party %d: ", id), log.Lmicroseconds),
	}

	// Initialise slices of shares.
	for i := 0; i < nParties; i++ {
		p.shares[i] = make(map[gate.Gate]*int)
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
func (p *Party) SendShare(to int, share int, gate int) {
	p.logger.Printf("Sending share %d at gate %d to party %d\n", share, gate, to)
	ch := p.subs[to]
	msg := &message{party: p.id, gate: gate, share: share}
	ch <- msg
}

// RecvShare receives a share
func (p *Party) RecvShare() *message {
	msg := <-p.ch
	p.logger.Printf("Received message %+v", msg)
	return msg
}

// Run runs the BGW protocol for this party.
func (p *Party) Run() int {
	p.logger.Println("Running...")

	nParties := p.circuit.NParties
	//nGates := p.circuit.NGates

	// 1. Split secret and share.
	//	  Each party generates a random polynomial with the 0th initialDeg term as the secret. This polynomial needs to be
	//    evaluated for each party to generate shares party.e. P(party) for party in 1..number of parties. Then
	//    broadcast shares to every other party.
	coeffs := make([]int, initialDeg(nParties)+1, initialDeg(nParties)+1)
	coeffs[0] = p.secret
	for d := 1; d <= initialDeg(nParties); d++ {
		coeffs[d] = p.field.Rand()
	}
	po := poly.New(coeffs, p.field)
	p.logger.Printf("Polynomial: %v", po)

	// 2. Run circuit.
	//    We iterate over all the gates in the circuit and wait for the inputs from the source gate, and then run the
	//    gate computation on the shares that we have. Then broadcast those shares.
	gates := p.circuit.Traverse()
	for gIdx, g := range gates {
		p.logger.Printf("Processing gate %d: %#v...\n", gIdx, g)
		switch v := g.(type) {
		case *gate.Input:
			p.processInput(gIdx, v, po, nParties)
		case *gate.Add:
			p.processAdd(v)
		case *gate.Mul:
			p.processMul()
		}
	}

	// 3. Create final result.
	outputGate := gates[len(gates)-1]

	// We broadcast our share to all other parties, and receive shares from all other parties.
	for party := 0; party < nParties; party++ {
		o := outputGate.Output()
		p.SendShare(party, o, len(gates)-1)
	}

	for party := 0; party < nParties; party++ {
		for p.shares[party][outputGate] == nil {
			msg := p.RecvShare()
			p.shares[msg.party][p.circuit.Gate(msg.gate)] = &msg.share
		}
	}

	terms := make([]int, nParties, nParties)
	for party := 0; party < nParties; party++ {
		share := *p.shares[party][outputGate]
		basis := poly.Recombination(party, initialDeg(nParties), p.field)
		terms[party] = p.field.Mul(basis, share)
	}
	output := p.field.Summation(terms)
	p.logger.Printf("Output: %d\n", output)

	return output
}

func (p *Party) processInput(gateIdx int, gate *gate.Input, po *poly.Poly, nParties int) {
	if gate.Party == p.id {
		for party := 0; party < nParties; party++ {
			i := party + 1
			share := po.Eval(i)
			p.shares[party][gate] = &share
			p.SendShare(party, share, gateIdx)
		}
	} else {
		for p.shares[gate.Party][gate] == nil {
			msg := p.RecvShare()
			p.shares[msg.party][p.circuit.Gate(msg.gate)] = &msg.share
		}
	}
	gate.SetOutput(*p.shares[gate.Party][gate])
}

func (p *Party) processAdd(gate *gate.Add) {
	out := gate.Output()
	gate.SetOutput(out)
}

func (p *Party) processMul() {
	// Check if we have the share.
	// If not, check for messages from other parties.
}

func initialDeg(nParties int) int {
	return nParties - 1
}
