// Package config contains hardcoded circuit configurations.
package config

import (
	"bgw/pkg/circuit"
	"bgw/pkg/field"
	"bgw/pkg/gate"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stderr, "Config: ", 0)

type Config struct {
	Secrets []int
	Circuit *circuit.Circuit
	Field   field.Field
}

func New(prime int, seed int64, defaultSeed int64, circuit int) *Config {
	if seed == defaultSeed {
		seed = time.Now().UnixNano()
	}

	fld := field.New(prime, seed)

	var cfg *Config
	switch circuit {
	case 1:
		cfg = config1(fld)
	default:
		logger.Fatalf("Unrecognised circuit number: %d", circuit)
	}

	// Validation
	if len(cfg.Secrets) != cfg.Circuit.NParties() {
		logger.Fatalf("Length mismatch between number of secrets and number of parties.")
	}

	return cfg
}

func config1(fld field.Field) *Config {
	return &Config{
		Secrets: []int{5, 28},
		Field:   fld,
		Circuit: circuit.New(
			[]int{0, 0},
			[]gate.Gate{
				gate.NewAdd(1, fld),
			},
		),
	}
}
