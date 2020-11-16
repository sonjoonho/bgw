# Privacy Engineering Coursework
Implementation of the BGW MPC protocol in Go.

## Getting Started

On a CSG machines, Go 1.15 is available at `/vol/linux/apps/go/bin/go`. If, for some reason, this doesn't work then Go
 1.14.3 should be installed at `/usr/lib/go-1.14/bin/go`. Below we will assume the `go` binary is on your path already. 

Run with 
```sh
go run cmd/mpc/mpc.go 
```

## Details

### Circuit Definition

TODO(sonjoonho)

### Party Communication

Parties are indexed from 0.

TODO(sonjoonho)

### Finite Field

We chose not to use `big.Int` for simplicity. Instead, all modular arithmetic functions are implemented in package 
`field`. 

## Authors
* Joon-Ho Son `<js6317>`
* William George Burr `<wb2117>`