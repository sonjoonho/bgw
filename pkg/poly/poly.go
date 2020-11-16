// Package poly provides a simple implementation of polynomials.
package poly

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"strings"
)

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

func (p *Poly) String() string {
	ss := []string{fmt.Sprint(p.coeffs[0])}
	for i, c := range p.coeffs {
		if i == 0 {
			continue
		}
		ss = append(ss, fmt.Sprintf("%dx^%d", c, i))
	}
	return strings.Join(ss, " + ")
}

func Recombination(party, deg int, field field.Field) int {
	i := party + 1
	terms := make([]int, deg+1, deg+1)
	for j := 1; j <= deg+1; j++ {
		if j == i {
			terms[j-1] = 1
		} else {
			terms[j-1] = field.Div(j, field.Sub(j, i))
		}
	}
	return field.Product(terms)
}
