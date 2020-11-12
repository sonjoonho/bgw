// Package field provides functions for mathematically correct modular arithmetic.
package field

import (
	"math/rand"
)

// Field is supposed to almost approximately represent something akin to the finite field defined by integers mod p (but
// not really). It provides arithmetic operations module Prime. For simplicity it uses the standard Go int type so will
// not support numbers bigger than 2^63 - 1. It is not robust, please be gentle.
// See: https://en.wikipedia.org/wiki/Finite_field.
type Field struct {
	// Prime is a prime number. The primeness of Prime is not enforced.
	Prime int
}

func New(prime int, seed int64) Field {
	rand.Seed(seed)
	return Field{Prime: prime}
}

// Mod implements the modulus function for Prime. Note that unlike some other languages the % operator implements
// remainder, which can return a negative value.
func (f Field) Mod(a int) int {
	m := a % f.Prime
	if a < 0 && f.Prime < 0 {
		m -= f.Prime
	}
	if a < 0 && f.Prime > 0 {
		m += f.Prime
	}
	return m
}

// Add adds two integers modulo Prime.
func (f Field) Add(a, b int) int {
	return f.Mod(a + b)
}

// Sub subtracts two integers modulo Prime.
func (f Field) Sub(a, b int) int {
	return f.Mod(a - b)
}

// Mul multiplies two integers modulo Prime.
func (f Field) Mul(a, b int) int {
	return f.Mod(a * b)
}

// Pow performs integer exponentiation modulo Prime.
func (f Field) Pow(a, b int) int {
	r := 1
	for b > 0 {
		if b&1 != 0 {
			r *= a
		}
		b >>= 1
		a *= a
	}
	return f.Mod(r)
}

// Inv computes the multiplicative inverse modulo Prime using Fermat's little theorem.
// See https://en.wikipedia.org/wiki/Fermat%27s_little_theorem.
func (f Field) Inv(a int) int {
	return f.Pow(a, f.Prime-2)
}

// Div performs integer division modulo Prime.
func (f Field) Div(a, b int) int {
	return f.Mul(a, f.Inv(b))
}

// Rand returns a random integer between n, 0 <= n < Prime.
func (f Field) Rand() int {
	return rand.Intn(f.Prime)
}

// Summation returns the sum of slice modulo Prime.
func (f Field) Summation(s []int) int {
	sum := 0
	for _, n := range s {
		sum = f.Add(sum, n)
	}
	return sum
}

// Product returns the product of a slice modulo Prime.
func (f Field) Product(s []int) int {
	prod := 1
	for _, n := range s {
		prod = f.Mul(prod, n)
	}
	return prod
}
