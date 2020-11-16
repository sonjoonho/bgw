package gate

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"

// Mul is an arithmetic multiplication gate.
type Mul struct {
	// first is the first input to this gate.
	first Gate
	// second is the second input to this gate.
	second Gate
	// field is the Field that we perform arithmetic over.
	field  field.Field
	output int
}

func (g *Mul) First() Gate {
	return g.first
}

func (g *Mul) Second() Gate {
	return g.second
}

func (g *Mul) SetOutput(output int) {
	g.output = output
}

func (g *Mul) Output() int {
	return g.field.Mul(g.first.Output(), g.second.Output())
}

func (g *Mul) Copy() Gate {
	return &Mul{
		first:  g.first.Copy(),
		second: g.second.Copy(),
		// field is a value struct so assignment makes a copy.
		field: g.field,
	}
}
