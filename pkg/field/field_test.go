package field

import (
	"fmt"
	"testing"
)

func TestField_Mod(t *testing.T) {
	tests := []struct {
		a    int
		p    int
		want int
	}{{
		a:    602,
		p:    101,
		want: 97,
	}, {
		a:    -42,
		p:    101,
		want: 59,
	}, {
		a:    -100,
		p:    101,
		want: 1,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d p=%d", tc.a, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Mod(tc.a), tc.want; got != want {
				t.Errorf("%v.Mod(%d) = %d, want %d", f, tc.a, got, want)
			}
		})
	}
}

func TestField_Add(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		p    int
		want int
	}{{
		a:    602,
		b:    103,
		p:    101,
		want: 99,
	}, {
		a:    -42,
		b:    103,
		p:    101,
		want: 61,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d b=%d p=%d", tc.a, tc.b, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Add(tc.a, tc.b), tc.want; got != want {
				t.Errorf("%v.Add(%d, %d) = %d, want %d", f, tc.a, tc.b, got, want)
			}
		})
	}
}
func TestField_Sub(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		p    int
		want int
	}{{
		a:    602,
		b:    103,
		p:    101,
		want: 95,
	}, {
		a:    -42,
		b:    103,
		p:    101,
		want: 57,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d b=%d p=%d", tc.a, tc.b, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Sub(tc.a, tc.b), tc.want; got != want {
				t.Errorf("%v.Sub(%d, %d) = %d, want %d", f, tc.a, tc.b, got, want)
			}
		})
	}
}

func TestField_Mul(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		p    int
		want int
	}{{
		a:    21,
		b:    1032,
		p:    101,
		want: 58,
	}, {
		a:    -100,
		b:    103,
		p:    101,
		want: 2,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d b=%d p=%d", tc.a, tc.b, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Mul(tc.a, tc.b), tc.want; got != want {
				t.Errorf("%v.Mul(%d, %d) = %d, want %d", f, tc.a, tc.b, got, want)
			}
		})
	}
}
func TestField_Pow(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		p    int
		want int
	}{{
		a:    21,
		b:    3,
		p:    101,
		want: 70,
	}, {
		a:    -3,
		b:    3,
		p:    11,
		want: 6,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d b=%d p=%d", tc.a, tc.b, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Pow(tc.a, tc.b), tc.want; got != want {
				t.Errorf("%v.Pow(%d, %d) = %d, want %d", f, tc.a, tc.b, got, want)
			}
		})
	}
}
func TestField_Inv(t *testing.T) {
	tests := []struct {
		a    int
		p    int
		want int
	}{{
		a:    29,
		p:    11,
		want: 8,
	}, {
		a:    -3,
		p:    11,
		want: 7,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d p=%d", tc.a, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Inv(tc.a), tc.want; got != want {
				t.Errorf("%v.Inv(%d) = %d, want %d", f, tc.a, got, want)
			}
		})
	}
}

func TestField_Div(t *testing.T) {
	tests := []struct {
		a    int
		b    int
		p    int
		want int
	}{{
		a:    29,
		b:    3,
		p:    11,
		want: 6,
	}, {
		a:    -3,
		b:    20,
		p:    11,
		want: 7,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("a=%d b=%d p=%d", tc.a, tc.b, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Div(tc.a, tc.b), tc.want; got != want {
				t.Errorf("%v.Div(%d, %d) = %d, want %d", f, tc.a, tc.b, got, want)
			}
		})
	}
}
func TestField_Rand(t *testing.T) {
	p := 3
	f := Field{Prime: p}

	for i := 0; i < 10; i++ {
		got := f.Rand()
		if got < 0 || got > 2 {
			t.Errorf("%v.Rand() = %v which is out of range for p = %d", f, got, p)
		}
	}
}

func TestField_Summation(t *testing.T) {
	tests := []struct {
		s    []int
		p    int
		want int
	}{{
		s:    []int{3, 1, 2, 6},
		p:    11,
		want: 1,
	}, {
		s:    []int{-19, 8},
		p:    11,
		want: 0,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("s=%d p=%d", tc.s, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Summation(tc.s), tc.want; got != want {
				t.Errorf("%v.Summation(%v) = %d, want %d", f, tc.s, got, want)
			}
		})
	}
}

func TestField_Product(t *testing.T) {
	tests := []struct {
		s    []int
		p    int
		want int
	}{{
		s:    []int{3, 1, 2, 6},
		p:    11,
		want: 3,
	}, {
		s:    []int{-19, 8},
		p:    11,
		want: 2,
	}}

	for _, tc := range tests {
		name := fmt.Sprintf("s=%d p=%d", tc.s, tc.p)
		t.Run(name, func(t *testing.T) {
			f := Field{Prime: tc.p}

			if got, want := f.Product(tc.s), tc.want; got != want {
				t.Errorf("%v.Product(%v) = %d, want %d", f, tc.s, got, want)
			}
		})
	}
}
