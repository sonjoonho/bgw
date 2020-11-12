// Package gate contains implementations of circuit gates.
package gate

// Gate represents an arithmetic gate. All gates should implement this interface.
type Gate interface {
	// NextGate returns the next gate in the circuit.
	NextGate() int
	// Output computes the output of this gate.
	Output(x, y int) int
}
