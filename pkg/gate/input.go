package gate

import "fmt"

// Input is an implicit party gate.
type Input struct {
	Party  int
	output int
}

func (g *Input) First() Gate {
	return nil
}

func (g *Input) Second() Gate {
	return nil
}

func (g *Input) SetOutput(output int) {
	g.output = output
}

func (g *Input) Output() int {
	return g.output
}

func (g *Input) Type() string {
	return fmt.Sprintf("IN%d", g.Party)
}

func (g *Input) Copy() Gate {
	return &Input{
		Party:  g.Party,
		output: g.output,
	}
}
