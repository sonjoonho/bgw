package circuit

import (
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"
	"testing"
)

func TestCircuit_ComputeExpected(t *testing.T) {
	secrets := []int{5, 28, 6}
	// fld is not important for this function.
	fld := field.New(0)
	circuit := &Circuit{
		Root: gate.NewAdd(
			&gate.Input{Party: 0},
			gate.NewAdd(
				&gate.Input{Party: 1},
				&gate.Input{Party: 2},
				fld,
			),
			fld,
		),
	}

	if got, want := circuit.ComputeExpected(secrets), 39; got != want {
		t.Errorf("circuit.ComputeExpected(%v) = %d, want %d", secrets, got, want)
	}
}
