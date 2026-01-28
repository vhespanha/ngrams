package ngrams

import (
	"math"
	"testing"
)

func TestProbValidate(t *testing.T) {
	tests := []struct {
		name  string
		p     prob
		valid bool
	}{
		{"zero", 0.0, true},
		{"one", 1.0, true},
		{"middle", 0.5, true},
		{"negative", -0.1, false},
		{"greater than one", 1.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.validate() != tt.valid {
				t.Errorf("prob(%v).validate() = %v, want %v", tt.p, !tt.valid, tt.valid)
			}
		})
	}
}

func TestProbMinMax(t *testing.T) {
	p := prob(0.5)
	if p.min() != 0 {
		t.Errorf("prob.min() = %v, want 0", p.min())
	}
	if p.max() != 1 {
		t.Errorf("prob.max() = %v, want 1", p.max())
	}
}

func TestNewProbTable(t *testing.T) {
	a := NewLowercaseAlphabet()
	pt := NewProbTable(2, 1000, a)

	if pt.n != 2 {
		t.Errorf("NewProbTable(2, ...).n = %d, want 2", pt.n)
	}
	if pt.total != 1000 {
		t.Errorf("NewProbTable(2, 1000, ...).total = %d, want 1000", pt.total)
	}
	if len(pt.freqs) != 26*26 {
		t.Errorf("NewProbTable(2, ...).freqs len = %d, want %d", len(pt.freqs), 26*26)
	}
}

func TestSetProbFromCount(t *testing.T) {
	a := NewLowercaseAlphabet()
	pt := NewProbTable(1, 100, a)

	t.Run("valid count", func(t *testing.T) {
		err := pt.SetProbFromCount(50, 0)
		if err != nil {
			t.Errorf("SetProbFromCount(50, 0) error = %v, want nil", err)
		}
		got := *pt.MustAt(0)
		if math.Abs(float64(got-0.5)) > 1e-10 {
			t.Errorf("After SetProbFromCount(50, 0), At(0) = %v, want 0.5", got)
		}
	})

	t.Run("wrong arity", func(t *testing.T) {
		err := pt.SetProbFromCount(50, 0, 0)
		if err == nil {
			t.Error("SetProbFromCount(50, 0, 0) error = nil, want error")
		}
	})
}

func TestMustSetProbFromCount(t *testing.T) {
	a := NewLowercaseAlphabet()
	pt := NewProbTable(1, 100, a)

	t.Run("valid count", func(t *testing.T) {
		pt.MustSetProbFromCount(30, 1)
		got := *pt.MustAt(1)
		if math.Abs(float64(got-0.3)) > 1e-10 {
			t.Errorf("After MustSetProbFromCount(30, 1), At(1) = %v, want 0.3", got)
		}
	})

	t.Run("wrong arity panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSetProbFromCount with wrong arity did not panic")
			}
		}()
		pt.MustSetProbFromCount(30, 0, 0)
	})
}
