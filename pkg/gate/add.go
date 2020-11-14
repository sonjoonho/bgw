package gate

import "gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"

// Add is an arithmetic addition gate.
type Add struct {
	// first is the first input to this gate.
	first Gate
	// second is the second input to this gate.
	second Gate
	// Field is the Field that we perform arithmetic over.
	field field.Field
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

// Add returns the sum of the inputs First and Second over the field Field.
func (g *Add) Output() int {
	return g.field.Add(g.First().Output(), g.Second().Output())
}

func (g *Add) Copy() Gate {
	return &Add{
		first:  g.first.Copy(),
		second: g.second.Copy(),
		// field is a value struct so assignment makes a copy.
		field: g.field,
	}
}
