// Package gate contains implementations of circuit gates.
package gate

// Gate represents a gate in the circuit. All gates should implement this interface.
type Gate interface {
	// First is the first input into this two-input gate.
	First() Gate
	// Second is the second input into this two-input gate.
	Second() Gate
	// SetOutput sets the output value of this gate (which can be retrieved with Output).
	SetOutput(int)
	// Output is the output value of this gate.
	Output() int
	// Type returns a human-readable representation of the type of this Gate, primarily for debugging purposes.
	Type() string
	// Copy returns a deep copy of this gate.
	Copy() Gate
}
