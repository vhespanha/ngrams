package ngrams

import (
	"math"
	"testing"
)

func TestNewRawTable(t *testing.T) {
	a := NewLowercaseAlphabet()
	rt := NewRawTable(2, 1000, a)

	if rt.n != 2 {
		t.Errorf("NewRawTable(2, ...).n = %d, want 2", rt.n)
	}
	if rt.total != 1000 {
		t.Errorf("NewRawTable(2, 1000, ...).total = %d, want 1000", rt.total)
	}
	if len(rt.freqs) != 26*26 {
		t.Errorf("NewRawTable(2, ...).freqs len = %d, want %d", len(rt.freqs), 26*26)
	}
}

func TestRawTableToProb(t *testing.T) {
	a := NewLowercaseAlphabet()
	rt := NewRawTable(1, 100, a)

	// Set up some frequencies
	rt.MustSet(50, 0) // 'a' has frequency 50
	rt.MustSet(30, 1) // 'b' has frequency 30
	rt.MustSet(20, 2) // 'c' has frequency 20

	pt := rt.ToProb()

	if pt.n != 1 {
		t.Errorf("ToProb().n = %d, want 1", pt.n)
	}
	if pt.total != 100 {
		t.Errorf("ToProb().total = %d, want 100", pt.total)
	}

	tests := []struct {
		symbol   symbol
		expected prob
	}{
		{0, 0.5},
		{1, 0.3},
		{2, 0.2},
	}

	for _, tt := range tests {
		got := *pt.MustAt(tt.symbol)
		if math.Abs(float64(got-tt.expected)) > 1e-10 {
			t.Errorf("ToProb() for symbol %d = %v, want %v", tt.symbol, got, tt.expected)
		}
	}
}

func TestRawTableToLogProb(t *testing.T) {
	a := NewLowercaseAlphabet()
	rt := NewRawTable(1, 100, a)

	// Set up some frequencies
	rt.MustSet(50, 0) // 'a' has frequency 50
	rt.MustSet(30, 1) // 'b' has frequency 30
	rt.MustSet(20, 2) // 'c' has frequency 20

	lpt := rt.ToLogProb()

	if lpt.n != 1 {
		t.Errorf("ToLogProb().n = %d, want 1", lpt.n)
	}
	if lpt.total != 100 {
		t.Errorf("ToLogProb().total = %d, want 100", lpt.total)
	}

	tests := []struct {
		symbol   symbol
		expected float64
	}{
		{0, math.Log(0.5)},
		{1, math.Log(0.3)},
		{2, math.Log(0.2)},
	}

	for _, tt := range tests {
		got := float64(*lpt.MustAt(tt.symbol))
		if math.Abs(got-tt.expected) > 1e-10 {
			t.Errorf("ToLogProb() for symbol %d = %v, want %v", tt.symbol, got, tt.expected)
		}
	}

	// Check that unset frequencies result in -Infinity
	got := float64(*lpt.MustAt(3))
	if !math.IsInf(got, -1) {
		t.Errorf("ToLogProb() for unset symbol = %v, want -Inf", got)
	}
}

func TestRawTableBigram(t *testing.T) {
	a := NewLowercaseAlphabet()
	rt := NewRawTable(2, 1000, a)

	// Set up a bigram frequency
	rt.MustSet(100, 0, 1) // 'ab' has frequency 100

	pt := rt.ToProb()
	got := *pt.MustAt(0, 1)
	expected := prob(0.1) // 100/1000

	if math.Abs(float64(got-expected)) > 1e-10 {
		t.Errorf("ToProb() for bigram ab = %v, want %v", got, expected)
	}
}
