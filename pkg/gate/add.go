package gate

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"

// Add is an arithmetic addition gate.
type Add struct {
	// first is the first input to this gate.
	first Gate
	// second is the second input to this gate.
	second Gate
	// field is the field that we perform arithmetic over.
	field  field.Field
	output int
}

func NewAdd(first Gate, second Gate, field field.Field) Gate {
	return &Add{
		first:  first,
		second: second,
		field:  field,
	}
}

func (g *Add) First() Gate {
	return g.first
}

func (g *Add) Second() Gate {
	return g.second
}

func (g *Add) SetOutput(output int) {
	g.output = output
}

func (g *Add) Output() int {
	return g.field.Add(g.first.Output(), g.second.Output())
}

func (g *Add) Copy() Gate {
	return &Add{
		first:  g.first.Copy(),
		second: g.second.Copy(),
		// field is a value struct so assignment makes a copy.
		field:  g.field,
		output: g.output,
	}
}
