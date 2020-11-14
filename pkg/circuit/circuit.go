// Package circuit contains the representation of a circuit.
package circuit

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"

// Circuit represents an arithmetic circuit to be computed by parties.
type Circuit struct {
	Root     gate.Gate
	NParties int
	NGates   int
}

// Copy makes a deep copy of this Circuit.
func (c *Circuit) Copy() *Circuit {
	return &Circuit{
		Root:     c.Root.Copy(),
		NParties: c.NParties,
		NGates:   c.NGates,
	}
}
