package gate

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"

// Add is an arithmetic addition gate.
type Add struct {
	// next is the next gate in the circuit.
	next int
	// field is the field that we perform arithmetic over.
	field field.Field
}

// NewAdd returns a new Add gate.
func NewAdd(next int, field field.Field) *Add {
	return &Add{next: next, field: field}
}

// NextGate returns the next gate in the circuit.
func (a *Add) NextGate() int {
	return a.next
}

// Output computes the output of this Add gate.
func (a *Add) Output(x, y int) int {
	return a.field.Add(x, y)
}

//
func (a *Add) Copy() Gate {
	return &Add{
		next: a.next,
		// field is a value struct, so no deep copy function is needed.
		field: a.field,
	}
}
