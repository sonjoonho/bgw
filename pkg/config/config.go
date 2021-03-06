// Package config contains hardcoded circuit configurations.
package config

import (
	"fmt"
	"github.com/sonjoonho/bgw/pkg/circuit"
	"github.com/sonjoonho/bgw/pkg/field"
	"github.com/sonjoonho/bgw/pkg/gate"
	"log"
	"math"
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
	case 8:
		cfg = config8(fld)
	case 9:
		cfg = config9(fld)
	case 10:
		cfg = config10(fld)
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

	if nSecrets, nParties := len(cfg.Secrets), cfg.Circuit.NParties; nSecrets != nParties {
		return nil, fmt.Errorf("length mismatch between number of secrets (%d) and number of parties (%d)", nSecrets, nParties)
	}

	return cfg, nil
}

// This is the example from Smart (p. 445).
func config1(field field.Field) *Config {
	return &Config{
		Secrets: []int{20, 40, 21, 31, 1, 71},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				gate.NewAdd(
					gate.NewMul(
						&gate.Input{Party: 0},
						&gate.Input{Party: 1},
					),
					gate.NewMul(
						&gate.Input{Party: 2},
						&gate.Input{Party: 3},
					),
				), gate.NewMul(
					&gate.Input{Party: 4},
					&gate.Input{Party: 5},
				),
			),
			NParties: 6,
		},
	}
}

// multree creates the gates for a multiplication of every party's input and returns the root gate.
func multree(nParties int, partyIdx int, field field.Field) gate.Gate {
	if nParties%2 == 1 {
		return gate.NewMul(
			&gate.Input{Party: partyIdx},
			multree(nParties-1, partyIdx+1, field),
		)
	}
	left := gate.NewMul(
		&gate.Input{Party: partyIdx},
		&gate.Input{Party: partyIdx + 1},
	)
	if nParties == 2 {
		return left
	}
	return gate.NewMul(
		left,
		multree(nParties-2, partyIdx+2, field),
	)
}

func config2(field field.Field) *Config {
	nParties := int(math.Pow(2, 3))
	root := multree(nParties, 0, field)

	secrets := make([]int, nParties)
	for i := 0; i < nParties; i++ {
		secrets[i] = i + 1
	}
	return &Config{
		Secrets: secrets,
		Field:   field,
		Circuit: &circuit.Circuit{
			Root:     root,
			NParties: nParties,
		},
	}
}

// Fibonacci sequence.
func config3(field field.Field) *Config {
	n := 10
	return &Config{
		Secrets: []int{0, 1},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root:     fibonacci(n),
			NParties: 2,
		},
	}
}

// fibonacci creates the gates for computing the nth fibonacci number and returns the root gate.
func fibonacci(n int) gate.Gate {
	fibGates := []gate.Gate{&gate.Input{Party: 0}, &gate.Input{Party: 1}}
	for k := 2; k <= n; k++ {
		fibGates = append(fibGates, gate.NewAdd(fibGates[len(fibGates)-1], fibGates[len(fibGates)-2]))
	}
	return fibGates[len(fibGates)-1]
}

// A single add gate.
func config4(field field.Field) *Config {
	return &Config{
		Secrets: []int{5, 28},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				&gate.Input{Party: 0},
				&gate.Input{Party: 1},
			),
			NParties: 2,
		},
	}
}

// Two add gates.
func config5(field field.Field) *Config {
	return &Config{
		Secrets: []int{5, 28, 6},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(
				&gate.Input{Party: 0},
				gate.NewAdd(
					&gate.Input{Party: 1},
					&gate.Input{Party: 2},
				),
			),
			NParties: 3,
		},
	}
}

// An add gate and multiplication gate.
func config6(field field.Field) *Config {
	return &Config{
		Secrets: []int{10, 20, 30},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewAdd(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
				),
				&gate.Input{Party: 2},
			),
			NParties: 3,
		},
	}
}

// Two multiplication gates.
func config7(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2, 3},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewMul(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
				),
				&gate.Input{Party: 2}),
			NParties: 3,
		},
	}
}

// Many addition gates.
func config8(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewAdd(&gate.Input{Party: 0}, gate.NewAdd(
				&gate.Input{Party: 1},
				gate.NewAdd(
					&gate.Input{Party: 2},
					gate.NewAdd(
						&gate.Input{Party: 3},
						gate.NewAdd(
							&gate.Input{Party: 4},
							gate.NewAdd(
								&gate.Input{Party: 5},
								gate.NewAdd(
									&gate.Input{Party: 6},
									gate.NewAdd(
										&gate.Input{Party: 7},
										gate.NewAdd(
											&gate.Input{Party: 8},
											gate.NewAdd(
												&gate.Input{Party: 9},
												gate.NewAdd(
													&gate.Input{Party: 10},
													gate.NewAdd(
														&gate.Input{Party: 11},
														gate.NewAdd(
															&gate.Input{Party: 12},
															gate.NewAdd(
																&gate.Input{Party: 13},
																gate.NewAdd(
																	&gate.Input{Party: 14},
																	gate.NewAdd(
																		&gate.Input{Party: 15},
																		&gate.Input{Party: 16},
																	),
																),
															),
														),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			)),
			NParties: 17,
		},
	}
}

// Party 0 has two inputs into the circuit.
func config9(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(
				gate.NewMul(
					&gate.Input{Party: 0},
					&gate.Input{Party: 1},
				),
				&gate.Input{Party: 0}),
			NParties: 2,
		},
	}
}

// Many multiplication gates.
func config10(field field.Field) *Config {
	return &Config{
		Secrets: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		Field:   field,
		Circuit: &circuit.Circuit{
			Root: gate.NewMul(&gate.Input{Party: 0}, gate.NewMul(
				&gate.Input{Party: 1},
				gate.NewMul(
					&gate.Input{Party: 2},
					gate.NewMul(
						&gate.Input{Party: 3},
						gate.NewMul(
							&gate.Input{Party: 4},
							gate.NewMul(
								&gate.Input{Party: 5},
								gate.NewMul(
									&gate.Input{Party: 6},
									gate.NewMul(
										&gate.Input{Party: 7},
										gate.NewMul(
											&gate.Input{Party: 8},
											gate.NewMul(
												&gate.Input{Party: 9},
												gate.NewMul(
													&gate.Input{Party: 10},
													gate.NewMul(
														&gate.Input{Party: 11},
														gate.NewMul(
															&gate.Input{Party: 12},
															gate.NewMul(
																&gate.Input{Party: 13},
																gate.NewMul(
																	&gate.Input{Party: 14},
																	gate.NewMul(
																		&gate.Input{Party: 15},
																		&gate.Input{Party: 16},
																	),
																),
															),
														),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			)),
			NParties: 17,
		},
	}
}
