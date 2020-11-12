// Package poly provides a simple implementation of polynomials.
package poly

import "bgw/pkg/field"

// Poly is a polynomial.
type Poly struct {
	coeffs []int
	field  field.Field
}

// New returns a new polynomial.
func New(coeffs []int, field field.Field) *Poly {
	return &Poly{coeffs: coeffs, field: field}
}

// Eval evaluates this polynomial at this value of x in the field.
func (p *Poly) Eval(x int) int {
	r := 0
	for i, c := range p.coeffs {
		r += c * p.field.Pow(x, i)
	}
	return p.field.Mod(r)
}
