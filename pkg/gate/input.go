package gate

// Input is an implicit party gate.
type Input struct {
	Party int
	Share int
}

func (g *Input) First() Gate {
	return nil
}

func (g *Input) Second() Gate {
	return nil
}

func (g *Input) Output() int {
	return g.Share
}

func (g *Input) Copy() Gate {
	return &Input{
		Party: g.Party,
		Share: g.Share,
	}
}
