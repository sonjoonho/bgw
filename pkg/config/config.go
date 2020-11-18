// Package config contains hardcoded circuit configurations.
package config

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/circuit"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/gate"
	"log"
	"math/rand"
	"os"
	"time"
)

var logger = log.New(os.Stderr, "Config: ", 0)

// Config is a configuration for the protocol.
type Config struct {
	// Secrets are the private values for each party.
	Secrets []int
	// Circuit is the circuit to be evaluated. A *copy* of this should be passed to each party to ensure that they do
	// not share memory.
	Circuit *circuit.Circuit
	// Field is the "finite field" that modular arithmetic is performed in.
	Field field.Field
	// Degree, also referred to as T, is the degree of polynomials used in Shamir Secret Sharing.
	Degree int
}

// New selects a configuration and performs validation on user inputs.
func New(prime int, seed, defaultSeed int64, degree, defaultDegree, circuit int) (*Config, error) {
	if seed == defaultSeed {
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)

	fld := field.New(prime)

	var cfg *Config
	switch circuit {
	case 1:
		cfg = config1(fld)
	case 2:
		cfg = config2(fld)
	case 3:
		cfg = config3(fld)
	case 4:
		cfg = config4(fld)
	case 5:
		cfg = config5(fld)
	case 6:
		cfg = config6(fld)
	case 7:
		cfg = config7(fld)
	default:
		logger.Fatalf("Unrecognised circuit number: %d", circuit)
	}

	if degree == defaultDegree {
		degree = (cfg.Circuit.NParties - 1) / 2
	}

	if degree < -1 {
		return nil, fmt.Errorf("degree=%d cannot be negative", degree)
	}

	if !(2*degree < cfg.Circuit.NParties) {
		return nil, fmt.Errorf("degree=%d does not satisfy 2T < N", degree)
	}

	cfg.Degree = degree

	if len(cfg.Secrets) != cfg.Circuit.NParties {
		return nil, fmt.Errorf("length mismatch between number of secrets and number of parties")
	}

	return cfg, nil
}
func config1(field field.Field) *Config {
	return &Config{
		Secrets: []int{5, 28},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				&gate.Input{Party: 0},
				&gate.Input{Party: 1},
				field,
			),
			NParties: 2,
		},
	}
}

func config2(field field.Field) *Config {
	return &Config{
		Secrets: []int{5, 28, 6},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				&gate.Input{Party: 0},
				gate.NewAdd(
					&gate.Input{Party: 1},
					&gate.Input{Party: 2},
					field,
				),
				field,
			),
			NParties: 3,
		},
	}
}

func config3(field field.Field) *Config {
	return &Config{
		Secrets: []int{10, 20, 30},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewAdd(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
					field,
				),
				&gate.Input{Party: 2},
				field,
			),
			NParties: 3,
		},
	}
}

func config4(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2, 3},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewMul(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
					field,
				),
				&gate.Input{Party: 2},
				field,
			),
			NParties: 3,
		},
	}
}

// This is the example from Smart (p. 445).
func config5(field field.Field) *Config {
	return &Config{
		Secrets: []int{20, 40, 21, 31, 1, 71},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				gate.NewAdd(
					gate.NewMul(
						&gate.Input{Party: 0},
						&gate.Input{Party: 1},
						field,
					),
					gate.NewMul(
						&gate.Input{Party: 2},
						&gate.Input{Party: 3},
						field,
					),
					field),
				gate.NewMul(
					&gate.Input{Party: 4},
					&gate.Input{Party: 5},
					field),
				field,
			),
			NParties: 6,
		},
	}
}

// Lots of additions.
func config6(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2, 3, 4, 5, 6},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				gate.NewAdd(
					gate.NewAdd(
						&gate.Input{Party: 0},
						&gate.Input{Party: 1},
						field,
					),
					gate.NewAdd(
						&gate.Input{Party: 2},
						&gate.Input{Party: 3},
						field,
					),
					field),
				gate.NewAdd(
					&gate.Input{Party: 4},
					&gate.Input{Party: 5},
					field),
				field,
			),
			NParties: 6,
		},
	}
}

func config7(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewMul(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
					field,
				),
				&gate.Input{Party: 0},
				field,
			),
			NParties: 2,
		},
	}
}
