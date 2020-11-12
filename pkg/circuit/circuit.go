// Package circuit contains the representation of a circuit.
package circuit

import "bgw/pkg/gate"

// Circuit represents an arithmetic circuit to be computed by parties.
type Circuit struct {
	// inputs is a mapping from party to input gate.
	inputs []int
	gates  []gate.Gate
}

func New(inputs []int, gates []gate.Gate) *Circuit {
	return &Circuit{inputs: inputs, gates: gates}
}

func (c *Circuit) InitialGate(party int) int {
	return c.inputs[party]
}

// NParties returns the number of parties in the protocol.
func (c *Circuit) NParties() int {
	return len(c.inputs)
}

// NParties returns the number of parties in the circuit.
func (c *Circuit) NGates() int {
	return len(c.gates)
}
