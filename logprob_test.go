package ngrams

import (
	"math"
	"testing"
)

func TestLogProbValidate(t *testing.T) {
	tests := []struct {
		name  string
		lp    logprob
		valid bool
	}{
		{"zero (log of 1)", 0.0, true},
		{"negative (log of fraction)", -1.0, true},
		{"negative infinity (log of 0)", logprob(math.Inf(-1)), true},
		{"positive (invalid)", 0.1, false},
		{"NaN (invalid)", logprob(math.NaN()), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.lp.validate() != tt.valid {
				t.Errorf("logprob(%v).validate() = %v, want %v", tt.lp, !tt.valid, tt.valid)
			}
		})
	}
}

func TestLogProbMinMax(t *testing.T) {
	lp := logprob(-1.0)
	if !math.IsInf(float64(lp.min()), -1) {
		t.Errorf("logprob.min() = %v, want -Inf", lp.min())
	}
	if lp.max() != 0 {
		t.Errorf("logprob.max() = %v, want 0", lp.max())
	}
}

func TestNewLogProb(t *testing.T) {
	tests := []struct {
		name     string
		v        uint64
		total    uint64
		expected float64
	}{
		{"half", 50, 100, math.Log(0.5)},
		{"quarter", 25, 100, math.Log(0.25)},
		{"all", 100, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := newLogProb(tt.v, tt.total)
			got := float64(lp)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("newLogProb(%d, %d) = %v, want %v", tt.v, tt.total, got, tt.expected)
			}
		})
	}
}

func TestNewLogProbTable(t *testing.T) {
	a := NewLowercaseAlphabet()
	lpt := NewLogProbTable(2, 1000, a)

	if lpt.n != 2 {
		t.Errorf("NewLogProbTable(2, ...).n = %d, want 2", lpt.n)
	}
	if lpt.total != 1000 {
		t.Errorf("NewLogProbTable(2, 1000, ...).total = %d, want 1000", lpt.total)
	}
	if len(lpt.freqs) != 26*26 {
		t.Errorf("NewLogProbTable(2, ...).freqs len = %d, want %d", len(lpt.freqs), 26*26)
	}

	// All entries should be initialized to -Inf
	for i, f := range lpt.freqs {
		if !math.IsInf(float64(f), -1) {
			t.Errorf("NewLogProbTable freqs[%d] = %v, want -Inf", i, f)
			break
		}
	}
}

func TestSetLogProbFromCount(t *testing.T) {
	a := NewLowercaseAlphabet()
	lpt := NewLogProbTable(1, 100, a)

	t.Run("valid count", func(t *testing.T) {
		err := lpt.SetLogProbFromCount(50, 0)
		if err != nil {
			t.Errorf("SetLogProbFromCount(50, 0) error = %v, want nil", err)
		}
		got := float64(*lpt.MustAt(0))
		expected := math.Log(0.5)
		if math.Abs(got-expected) > 1e-10 {
			t.Errorf("After SetLogProbFromCount(50, 0), At(0) = %v, want %v", got, expected)
		}
	})

	t.Run("wrong arity", func(t *testing.T) {
		err := lpt.SetLogProbFromCount(50, 0, 0)
		if err == nil {
			t.Error("SetLogProbFromCount(50, 0, 0) error = nil, want error")
		}
	})
}

func TestMustSetLogProbFromCount(t *testing.T) {
	a := NewLowercaseAlphabet()
	lpt := NewLogProbTable(1, 100, a)

	t.Run("valid count", func(t *testing.T) {
		lpt.MustSetLogProbFromCount(30, 1)
		got := float64(*lpt.MustAt(1))
		expected := math.Log(0.3)
		if math.Abs(got-expected) > 1e-10 {
			t.Errorf("After MustSetLogProbFromCount(30, 1), At(1) = %v, want %v", got, expected)
		}
	})

	t.Run("wrong arity panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSetLogProbFromCount with wrong arity did not panic")
			}
		}()
		lpt.MustSetLogProbFromCount(30, 0, 0)
	})
}
