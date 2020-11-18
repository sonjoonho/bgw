package gate

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"

// Mul is an arithmetic multiplication gate.
type Mul struct {
	// first is the first input to this gate.
	first Gate
	// second is the second input to this gate.
	second Gate
	// field is the Field that we perform arithmetic over.
	field field.Field
	// output is the output value of this gate.
	output int
}

func NewMul(first Gate, second Gate, field field.Field) Gate {
	return &Mul{
		first:  first,
		second: second,
		field:  field,
	}
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
	return g.output
}

func (g *Mul) Type() string {
	return "MUL"
}

func (g *Mul) Copy() Gate {
	return &Mul{
		first:  g.first.Copy(),
		second: g.second.Copy(),
		// field is a value struct so assignment makes a copy.
		field: g.field,
	}
}
