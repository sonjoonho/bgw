// Package circuit contains the representation of a circuit.
package circuit

import (
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"
)

// Circuit represents an arithmetic circuit to be computed by parties. It is not thread safe -- each Goroutine should be
// given a copy.
type Circuit struct {
	// Root is the output gate of the circuit.
	Root gate.Gate
	// NParties is the number of parties that have inputs in this circuit.
	NParties int
	// gates are the gates of this circuit, ordered linearly in the order which they should be processed. It is lazily
	// initialised.
	gates []gate.Gate
}

// Copy makes a deep copy of this Circuit.
func (c *Circuit) Copy() *Circuit {
	return &Circuit{
		Root:     c.Root.Copy(),
		NParties: c.NParties,
	}
}

// Gate returns the gate for index i. If the order of of gates has not been determined, we run Traverse.
func (c *Circuit) Gate(i int) gate.Gate {
	if c.gates == nil {
		c.gates = c.Traverse()
	}
	return c.gates[i]
}

// Traverse traverses the circuit and returns the gates in order. This must be deterministic so that each party receives
//gates with matching indexes.
func (c *Circuit) Traverse() []gate.Gate {
	if c.gates != nil {
		return c.gates
	}

	stack := []gate.Gate{c.Root}
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

	c.gates = res

	return res
}

// ComputeExpected recursively evaluates the circuit using secrets as input and returns the expected value.
func (c *Circuit) ComputeExpected(secrets []int) int {
	return eval(c.Root, secrets)
}

func eval(root gate.Gate, s []int) int {
	fst := root.First()
	snd := root.Second()

	switch v := root.(type) {
	case *gate.Input:
		return s[v.Party]
	case *gate.Add:
		return eval(fst, s) + eval(snd, s)
	case *gate.Mul:
		return eval(fst, s) * eval(snd, s)
	default:
		panic("Unrecognised gate type in circuit")
	}
}

func peek(stack []gate.Gate) gate.Gate {
	return stack[len(stack)-1]
}

func pop(stack []gate.Gate) ([]gate.Gate, gate.Gate) {
	g := peek(stack)
	return stack[0 : len(stack)-1], g
}
