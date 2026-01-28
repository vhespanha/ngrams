package ngrams

import (
	"math"
	"testing"
)

func TestNewTable(t *testing.T) {
	a := NewLowercaseAlphabet()

	tests := []struct {
		name          string
		n             int
		expectedFreqs int
	}{
		{
			name:          "unigram",
			n:             1,
			expectedFreqs: 26,
		},
		{
			name:          "bigram",
			n:             2,
			expectedFreqs: 26 * 26,
		},
		{
			name:          "trigram",
			n:             3,
			expectedFreqs: 26 * 26 * 26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := newTable[uint64](tt.n, 1000, a)
			if tbl.n != tt.n {
				t.Errorf("newTable(%d, ...).n = %d, want %d", tt.n, tbl.n, tt.n)
			}
			if len(tbl.freqs) != tt.expectedFreqs {
				t.Errorf("newTable(%d, ...).freqs len = %d, want %d",
					tt.n, len(tbl.freqs), tt.expectedFreqs)
			}
			if tbl.total != 1000 {
				t.Errorf("newTable(%d, ...).total = %d, want 1000", tt.n, tbl.total)
			}
		})
	}
}

func TestTableIdx(t *testing.T) {
	a := NewLowercaseAlphabet()
	tbl := newTable[uint64](2, 1000, a)

	// For bigrams with lowercase alphabet (26 chars):
	// idx(a,a) = 0*26 + 0 = 0
	// idx(a,b) = 0*26 + 1 = 1
	// idx(b,a) = 1*26 + 0 = 26
	// idx(z,z) = 25*26 + 25 = 675

	tests := []struct {
		name     string
		symbols  []symbol
		expected int
	}{
		{
			name:     "aa",
			symbols:  []symbol{0, 0},
			expected: 0,
		},
		{
			name:     "ab",
			symbols:  []symbol{0, 1},
			expected: 1,
		},
		{
			name:     "ba",
			symbols:  []symbol{1, 0},
			expected: 26,
		},
		{
			name:     "zz",
			symbols:  []symbol{25, 25},
			expected: 675,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tbl.idx(tt.symbols)
			if idx != tt.expected {
				t.Errorf("idx(%v) = %d, want %d", tt.symbols, idx, tt.expected)
			}
		})
	}
}

func TestTableAt(t *testing.T) {
	a := NewLowercaseAlphabet()
	tbl := newTable[uint64](2, 1000, a)

	t.Run("valid access", func(t *testing.T) {
		ptr, err := tbl.At(0, 0)
		if err != nil {
			t.Errorf("At(0, 0) error = %v, want nil", err)
		}
		if ptr == nil {
			t.Error("At(0, 0) returned nil pointer")
		}
	})

	t.Run("wrong arity", func(t *testing.T) {
		_, err := tbl.At(0)
		if err == nil {
			t.Error("At(0) error = nil, want error for wrong arity")
		}
	})

	t.Run("wrong arity too many", func(t *testing.T) {
		_, err := tbl.At(0, 0, 0)
		if err == nil {
			t.Error("At(0, 0, 0) error = nil, want error for wrong arity")
		}
	})
}

func TestTableMustAt(t *testing.T) {
	a := NewLowercaseAlphabet()
	tbl := newTable[uint64](2, 1000, a)

	t.Run("valid access", func(t *testing.T) {
		ptr := tbl.MustAt(0, 0)
		if ptr == nil {
			t.Error("MustAt(0, 0) returned nil pointer")
		}
	})

	t.Run("wrong arity panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustAt(0) did not panic with wrong arity")
			}
		}()
		tbl.MustAt(0)
	})
}

func TestTableSet(t *testing.T) {
	a := NewLowercaseAlphabet()
	tbl := newTable[uint64](2, 1000, a)

	t.Run("valid set", func(t *testing.T) {
		err := tbl.Set(42, 0, 0)
		if err != nil {
			t.Errorf("Set(42, 0, 0) error = %v, want nil", err)
		}
		ptr, _ := tbl.At(0, 0)
		if *ptr != 42 {
			t.Errorf("After Set(42, 0, 0), At(0, 0) = %d, want 42", *ptr)
		}
	})

	t.Run("wrong arity", func(t *testing.T) {
		err := tbl.Set(42, 0)
		if err == nil {
			t.Error("Set(42, 0) error = nil, want error for wrong arity")
		}
	})
}

func TestTableMustSet(t *testing.T) {
	a := NewLowercaseAlphabet()
	tbl := newTable[uint64](2, 1000, a)

	t.Run("valid set", func(t *testing.T) {
		tbl.MustSet(99, 1, 1)
		ptr := tbl.MustAt(1, 1)
		if *ptr != 99 {
			t.Errorf("After MustSet(99, 1, 1), MustAt(1, 1) = %d, want 99", *ptr)
		}
	})

	t.Run("wrong arity panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSet(99, 0) did not panic with wrong arity")
			}
		}()
		tbl.MustSet(99, 0)
	})
}

func TestTableWithDifferentTypes(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("uint64 table", func(t *testing.T) {
		tbl := newTable[uint64](1, 1000, a)
		tbl.MustSet(100, 0)
		if *tbl.MustAt(0) != 100 {
			t.Error("uint64 table set/get mismatch")
		}
	})

	t.Run("prob table", func(t *testing.T) {
		tbl := newTable[prob](1, 1000, a)
		tbl.MustSet(0.5, 0)
		if *tbl.MustAt(0) != 0.5 {
			t.Error("prob table set/get mismatch")
		}
	})

	t.Run("logprob table", func(t *testing.T) {
		tbl := newTable[logprob](1, 1000, a)
		tbl.MustSet(logprob(math.Log(0.5)), 0)
		got := float64(*tbl.MustAt(0))
		want := math.Log(0.5)
		if math.Abs(got-want) > 1e-10 {
			t.Errorf("logprob table set/get: got %v, want %v", got, want)
		}
	})
}
