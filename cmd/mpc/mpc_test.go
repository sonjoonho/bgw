package main

import (
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/circuit"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/config"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"
	"testing"
)

func TestRunProtocol(t *testing.T) {
	fld := field.New(101)
	tests := []struct {
		name string
		cfg  *config.Config
		want int
	}{{
		name: "Textbook example",
		cfg: &config.Config{
			Secrets: []int{20, 40, 21, 31, 1, 71},
			Field:   fld,
			Degree:  2,
			Circuit: &circuit.Circuit{
				NParties: 6,
				Root: gate.NewAdd(
					gate.NewAdd(
						gate.NewMul(
							&gate.Input{Party: 0},
							&gate.Input{Party: 1},
							fld,
						),
						gate.NewMul(
							&gate.Input{Party: 2},
							&gate.Input{Party: 3},
							fld,
						),
						fld),
					gate.NewMul(
						&gate.Input{Party: 4},
						&gate.Input{Party: 5},
						fld),
					fld,
				),
			},
		},
		want: 7,
	}, {
		name: "Many adds",
		want: 21,
		cfg: &config.Config{
			Secrets: []int{1, 2, 3, 4, 5, 6},
			Field:   fld,
			Circuit: &circuit.Circuit{
				NParties: 6,
				Root: gate.NewAdd(
					gate.NewAdd(
						gate.NewAdd(
							&gate.Input{Party: 0},
							&gate.Input{Party: 1},
							fld,
						),
						gate.NewAdd(
							&gate.Input{Party: 2},
							&gate.Input{Party: 3},
							fld,
						),
						fld),
					gate.NewAdd(
						&gate.Input{Party: 4},
						&gate.Input{Party: 5},
						fld),
					fld,
				),
			},
		},
	}, {
		name: "Multiple inputs for one party",
		want: 6,
		cfg: &config.Config{
			Secrets: []int{1, 2},
			Field:   fld,
			Circuit: &circuit.Circuit{
				NParties: 2,
				Root: gate.NewAdd(
					gate.NewAdd(
						gate.NewMul(
							&gate.Input{Party: 0},
							&gate.Input{Party: 1},
							fld,
						),
						gate.NewAdd(
							&gate.Input{Party: 1},
							&gate.Input{Party: 0},
							fld,
						),
						fld),
					gate.NewMul(
						&gate.Input{Party: 0},
						&gate.Input{Party: 0},
						fld),
					fld,
				),
			},
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RunProtocol(tc.cfg)
			if err != nil {
				t.Errorf("RunProtocol(%v) failed with %v", tc.cfg, err)
			} else if got != tc.want {
				t.Errorf("RunProtocol(%v) = %d, want %d", tc.cfg, got, tc.want)
			}
		})
	}
}
