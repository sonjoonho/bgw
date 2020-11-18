# Privacy Engineering Coursework
Implementation of the BGW MPC protocol in Go.

[![Pipeline status](https://gitlab.doc.ic.ac.uk/js6317/bgw/badges/master/pipeline.svg)](https://gitlab.doc.ic.ac.uk/js6317/bgw/)

## Getting Started

On a CSG machines, Go 1.15 is available at `/vol/linux/apps/go/bin/go`. If, for some reason, this doesn't work then Go
 1.14.3 should be installed at `/usr/lib/go-1.14/bin/go`. Below we will assume the `go` binary is on your path already. 

Run with 
```sh
go run cmd/mpc/mpc.go 
```

Output can be sorted by party by running

```sh
go run cmd/mpc/mpc.go | sort
```

The protocol can be configured using command line arguments. For example:

```sh
go run cmd/mpc/mpc.go -circuit 5 -degree 2 -prime 1003 -seed 4
```

The various circuit definitions can be found in `pkg/config/config.go`. The full usage is detailed below:

```
Usage of mpc:
  -circuit int
    	Circuit to run. (default 1)
  -degree int
    	Degree of polynomial. If unset, it is set to N-1/2 (default -1)
  -prime int
    	Prime number to use for modular arithmetic. (default 101)
  -seed int
    	Seed for pseudorandom number generation. If unset, the current time is used.
```

## Details
### Party Communication

The main protocol is implemented in `party.Run`. It first traversed the circuit "tree" and processes each gate in turn. 
After computing the output of a gate, it's `Output` value is set so that other gates that depend on it can access its value. 
The traversal is done in such a way that a gate's dependencies are always available. 

Parties communicate with channels. Each party has a single `msg` channel, through which **all** inter-party communication 
is done. Each party is initialised with a copy of the circuit, to prevent any accidental shared memory.

Parties are indexed from 0, although they are indexed from 1 for the purposed of calculations (e.g. computing the 
recombination vector).

### Circuit Definition

Circuits are represented using the struct `circuit.Circuit`. They are defined using a tree-like structure, with
 `gate.Gate`s as nodes. See `pkg/config/config.go` for examples.

### Finite Field

We chose not to use `big.Int` for simplicity. Instead, all modular arithmetic functions are implemented in package 
`field`. 

## Authors
* Joon-Ho Son `<js6317>`
* William George Burr `<wb2117>`