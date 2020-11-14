// Package gate contains implementations of circuit gates.
package gate

// Gate represents a gate in the circuit. All gates should implement this interface.
type Gate interface {
	First() Gate
	Second() Gate
	Output() int
	Copy() Gate
}
