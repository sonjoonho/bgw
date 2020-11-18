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
	"strings"
)

// message represents a message used for inter-party communication.
type message struct {
	// party is the *source* party.
	party int
	gate  int
	share int
}

// Party is a party which can communicate with other parties.
type Party struct {
	// id is the identifier of this Party. It starts from 0.
	id int
	// secret is this party's secret.
	secret int
	// ch is a channel through which this Party receives messages.
	ch chan *message
	// done is...
	done chan bool
	// subs is a slice of send-only channels that this party uses to send message to subscribers. Its capacity is equal
	// to the number of parties, specified during initialisation.
	subs []chan<- *message
	// shares is a buffer for received shares. It maps from Party id to gate.Gate to share.
	shares []map[gate.Gate]*int
	// field is the field that we perform arithmetic over.
	field field.Field
	// circuit is the circuit that this party evaluates.
	circuit *circuit.Circuit
	// degree is the degree of the polynomial in Shamir Secret Sharing.
	degree int
	logger *log.Logger
}

// New initialises and returns a new Party. nParties specifies the number of parties participating in the protocol,
// and nGates specifies the number of gates in the circuit.
func New(id int, secret int, circuit *circuit.Circuit, field field.Field, degree int) *Party {
	nParties := circuit.NParties

	p := &Party{
		id:      id,
		secret:  secret,
		circuit: circuit,
		field:   field,
		ch:      make(chan *message, nParties+1),
		subs:    make([]chan<- *message, nParties, nParties),
		shares:  make([]map[gate.Gate]*int, nParties, nParties),
		degree:  degree,
		logger:  log.New(os.Stdout, fmt.Sprintf("Party %d: ", id), log.Lmicroseconds),
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
	ch := p.subs[to]
	msg := &message{party: p.id, gate: gate, share: share}
	ch <- msg
}

// RecvShare receives a share
func (p *Party) RecvShare() *message {
	msg := <-p.ch
	return msg
}

// Run runs the BGW protocol for this party.
func (p *Party) Run() int {
	p.logger.Printf("Running with secret %d\n", p.secret)

	// 2. Run circuit.
	//    We iterate over all the gates in the circuit and wait for the inputs from the source gate, and then run the
	//    gate computation on the shares that we have. Then broadcast those shares.
	gates := p.circuit.Traverse()
	// After a each iteration of this loop, the value of g.Output() must be set.
	for gIdx, g := range gates {
		switch v := g.(type) {
		case *gate.Input:
			p.processInput(gIdx, v)
		case *gate.Add:
			p.processAdd(gIdx, v)
		case *gate.Mul:
			p.processMul(gIdx, v)
		}
	}

	// 3. Create final result. The final gate will always be the output gate.
	outputGate := gates[len(gates)-1]
	output := p.processOutput(len(gates)-1, outputGate)

	return output
}

func (p *Party) processInput(gateIdx int, gate *gate.Input) {
	gatePrefix := fmt.Sprintf("[%d | %s]", gateIdx, gate.Type())

	nParties := p.circuit.NParties

	// 1. Split secret and share.
	//	  Each party generates a random polynomial with the 0th initialDeg term as the secret. This polynomial needs to
	//	  be evaluated for each party to generate shares party.e. P(party) for party in 1..number of parties. Then
	//    broadcast shares to every other party.

	// If this gate is an input corresponding to this party, we want to send shares to all parties (including itself).
	// Otherwise, we want to receive shares from all other parties.
	if gate.Party == p.id {
		po := poly.Random(p.secret, p.degree, p.field)

		// sentShares are the shares sent from this party. This variable is used for logging only.
		sentShares := make([]int, nParties, nParties)
		p.logger.Printf("%s using polynomial %s", gatePrefix, po)

		for party := 0; party < nParties; party++ {
			i := party + 1
			share := po.Eval(i)
			if party != p.id {
				p.SendShare(party, share, gateIdx)
			} else {
				p.shares[party][gate] = &share
			}

			sentShares[party] = share
		}

		p.logger.Printf("%s sent shares %v", gatePrefix, sentShares)
	} else {
		for p.shares[gate.Party][gate] == nil {
			msg := p.RecvShare()
			p.shares[msg.party][p.circuit.Gate(msg.gate)] = &msg.share

			p.logger.Printf("%s received share %d from party %d", gatePrefix, msg.share, msg.party)
		}
	}
	gate.SetOutput(*p.shares[gate.Party][gate])
}

func (p *Party) processAdd(gateIdx int, gate *gate.Add) {
	gatePrefix := gatePrefix(gateIdx, gate)
	fst := gate.First().Output()
	snd := gate.Second().Output()
	prime := p.field.Prime

	out := p.field.Add(gate.First().Output(), gate.Second().Output())

	p.logger.Printf("%s %d + %d mod %d = %d", gatePrefix, fst, snd, prime, out)

	gate.SetOutput(out)
}

func (p *Party) processMul(gateIdx int, gate *gate.Mul) {
	// gatePrefix marks this gate in the logging output for readability.
	gatePrefix := gatePrefix(gateIdx, gate)
	fst := gate.First().Output()
	snd := gate.Second().Output()
	prime := p.field.Prime

	// 1. Each party locally computes d = a * b.
	out := p.field.Mul(gate.First().Output(), gate.Second().Output())

	p.logger.Printf("%s %d × %d mod %d = %d", gatePrefix, fst, snd, prime, out)

	// 2. Each party produces a polynomial delta of degree at most degree such delta_i(0) = d^i.
	nParties := p.circuit.NParties
	po := poly.Random(out, p.degree, p.field)

	p.logger.Printf("%s using polynomial %s", gatePrefix, po)

	// 3. Each party i distributes to party j the value d_{i, j} = delta_i(j).

	// sentShares are the shares sent from this party. This variable is used for logging only.
	sentShares := make([]int, nParties, nParties)
	for party := 0; party < nParties; party++ {
		i := party + 1
		share := po.Eval(i)
		p.SendShare(party, share, gateIdx)

		sentShares[party] = share
	}

	p.logger.Printf("%s sent shares %v", gatePrefix, sentShares)

	// recvShares are the shares for all parties for this gate. This variable is used for logging only.
	recvShares := make([]int, nParties, nParties)
	for party := 0; party < nParties; party++ {
		for p.shares[party][gate] == nil {
			msg := p.RecvShare()
			p.shares[msg.party][p.circuit.Gate(msg.gate)] = &msg.share
		}

		recvShares[party] = *p.shares[party][gate]
	}
	// At this point, all shares for this gate will have been received.
	// i.e. p.shares[parties][gate] != nil for all parties.

	p.logger.Printf("%s received shares %v", gatePrefix, recvShares)

	// Each party j computes c^j.
	terms := make([]int, nParties, nParties)

	// termStrings are the terms of the summation formatted as a string for debugging.
	termsStrings := make([]string, nParties, nParties)
	for party := 0; party < nParties; party++ {
		share := *p.shares[party][gate]
		basis := poly.Recombination(party, nParties)
		terms[party] = p.field.Mul(share, basis)

		termsStrings[party] = fmt.Sprintf("(%d × %d)", share, basis)
	}
	output := p.field.Summation(terms)

	summationString := strings.Join(termsStrings, " + ")
	p.logger.Printf("%s %s mod %d = %d", gatePrefix, summationString, prime, output)

	gate.SetOutput(output)
}

func (p *Party) processOutput(gateIdx int, gate gate.Gate) int {
	gatePrefix := fmt.Sprintf("[%d | OUTPUT]", gateIdx)

	nParties := p.circuit.NParties

	// We broadcast our share to all other parties, and receive shares from all other parties.
	sentShares := make([]int, nParties, nParties)
	for party := 0; party < nParties; party++ {
		share := gate.Output()
		// gateIdx + 1 identifies the implicit "output gate".
		p.SendShare(party, share, gateIdx+1)

		sentShares[party] = share
	}

	p.logger.Printf("%s sent shares %v", gatePrefix, sentShares)

	outputShares := make([]*int, nParties, nParties)
	for party := 0; party < nParties; party++ {
		for outputShares[party] == nil {
			msg := p.RecvShare()
			outputShares[msg.party] = &msg.share
		}
	}
	// All elements of outputShares will be populated.
	outputSharesStrings := make([]string, nParties, nParties)
	for party := 0; party < nParties; party++ {
		outputSharesStrings[party] = fmt.Sprint(*outputShares[party])
	}
	outputSharesString := "[" + strings.Join(outputSharesStrings, " ") + "]"
	p.logger.Printf("%s received shares %v", gatePrefix, outputSharesString)

	terms := make([]int, nParties, nParties)
	termsStrings := make([]string, nParties, nParties)
	for party := 0; party < nParties; party++ {
		share := *outputShares[party]
		basis := poly.Recombination(party, nParties)
		terms[party] = p.field.Mul(basis, share)

		termsStrings[party] = fmt.Sprintf("(%d × %d)", share, basis)
	}
	output := p.field.Summation(terms)

	prime := p.field.Prime
	summationString := strings.Join(termsStrings, " + ")
	p.logger.Printf("%s %s mod %d = %d\n", gatePrefix, summationString, prime, output)
	p.logger.Printf("%s output: %d\n", gatePrefix, output)

	return output
}

// gatePrefix returns a formatted tag representing a gate e.g. [3 | MUL].
func gatePrefix(gateIdx int, gate gate.Gate) string {
	return fmt.Sprintf("[%d | %s]", gateIdx, gate.Type())
}
