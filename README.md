# BGW
Implementation of the BGW MPC protocol in Go.

## Getting Started

On a CSG machines, Go 1.15 is available at `/vol/linux/apps/go/bin/go`. If, for some reason, this doesn't work then Go
 1.14.3 should be installed at `/usr/lib/go-1.14/bin/go`. Below we will assume the `go` binary is on your path already. 

Run with 
```sh
go run cmd/mpc/mpc.go 
```

We recommend sorting the output for readability. To do this, run:

```sh
go run cmd/mpc/mpc.go | sort
```

If the protocol finishes successfully, you should see that the program finishes with something like:

```sh
[...]
MPC: 17:59:37.362092 Expected output: 7
MPC: 17:59:37.362104 Actual output:   7
MPC: 17:59:37.362112 Protocol succeeded (:
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

### Log Format

Log lines for each party's computation follows the format:
```sh
00<PARTY_NUMBER>: <TIMESTAMP>  [<GATE_NUMBER> | <GATE_TYPE>] <LOG MESSAGE>
``` 

Gates are indented according to their level in the tree.

### Party Communication

The main protocol is implemented in `party.Run`. It first traverses the circuit "tree" and processes each gate in turn. 
After computing the output of a gate, it's `Output` value is set so that other gates that depend on it can access its value. 
The traversal is done in such a way that a gate's dependencies are always available. 

Parties communicate via their `msg` channel, through which *all* inter-party communication 
is done. Each party is initialised with a copy of the circuit, to prevent any accidental shared memory.

Parties are indexed from 0, although they are indexed from 1 for the purpose of calculations (e.g. computing the 
recombination vector).

### Circuit Definition

Circuits are represented using the struct `circuit.Circuit`. They are defined using a tree-like structure, with
 `gate.Gate`s as nodes. 
 
```go
&circuit.Circuit{
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
```

Circuits that have multiple inputs into the circuit are supported. For example:

```go
&circuit.Circuit{
    NParties: 2,
    Root: gate.NewMul(
        gate.NewMul(
            &gate.Input{Party: 0},
            &gate.Input{Party: 1},
        ),
        &gate.Input{Party: 0},
    ),
}
``` 

See `pkg/config/config.go` for the full list of hardcoded circuit configurations.

### Finite Field

We chose not to use `big.Int` for simplicity, opting for the standard `int` type. Instead, all modular arithmetic functions are implemented in package 
`field`. However, this does limit the size of input/prime that can be used.

## Authors
* Joon-Ho Son `<js6317>`
* William George Burr `<wb2117>`
