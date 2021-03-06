// Package poly provides a simple implementation of polynomials.
package poly

import (
	"fmt"
	"github.com/sonjoonho/bgw/pkg/field"
	"math"
	"strings"
)

// Poly is a polynomial.
type Poly struct {
	Coeffs []int
	field  field.Field
}

// New returns a new polynomial with the specified coefficients.
func New(coeffs []int, field field.Field) *Poly {
	return &Poly{Coeffs: coeffs, field: field}
}

// Random returns a polynomial with constant c, and random coefficient for all other terms.
func Random(c int, deg int, field field.Field) *Poly {
	coeffs := make([]int, deg+1, deg+1)
	coeffs[0] = c
	for d := 1; d <= deg; d++ {
		coeffs[d] = field.Rand()
	}

	return New(coeffs, field)
}

// Eval evaluates this polynomial at this value of x in the field.
func (p *Poly) Eval(x int) int {
	r := 0
	for i, c := range p.Coeffs {
		r += c * p.field.Pow(x, i)
	}
	return p.field.Mod(r)
}

// String returns the string representation of this polynomial.
func (p *Poly) String() string {
	ss := []string{fmt.Sprint(p.Coeffs[0])}
	for i, c := range p.Coeffs {
		if i == 0 {
			continue
		}
		ss = append(ss, fmt.Sprintf("%dx^%d", c, i))
	}
	return strings.Join(ss, " + ")
}

// Recombination returns an element, r_i (or delta_i(0)), of the recombination vector of the specified length.
func Recombination(party, length int) int {
	i := party + 1
	terms := make([]float64, length, length)
	for j := 1; j <= length; j++ {
		if j == i {
			terms[j-1] = 1
		} else {
			terms[j-1] = float64(j) / float64(j-i)
		}
	}
	return product(terms)
}

// product returns the product of the elements of a slice, rounded to the nearest integer.
func product(s []float64) int {
	prod := 1.0
	for _, n := range s {
		prod *= n
	}
	return int(math.Round(prod))
}
