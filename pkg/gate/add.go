package gate

// Add is an arithmetic addition gate.
type Add struct {
	// first is the first input to this gate.
	first Gate
	// second is the second input to this gate.
	second Gate
	// output is the output value of this gate.
	output int
}

func NewAdd(first Gate, second Gate) Gate {
	return &Add{
		first:  first,
		second: second,
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
	return g.output
}

func (g *Add) Type() string {
	return "ADD"
}

func (g *Add) Copy() Gate {
	return &Add{
		first:  g.first.Copy(),
		second: g.second.Copy(),
		output: g.output,
	}
}
