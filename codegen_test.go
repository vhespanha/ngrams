package ngrams

import (
	"strings"
	"testing"
)

func TestGenGo(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("raw table", func(t *testing.T) {
		rt := NewRawTable(1, 100, a)
		rt.MustSet(50, 0)

		code := rt.GenGo("testpkg")

		// Check for expected content
		if !strings.Contains(code, "package testpkg") {
			t.Error("GenGo() missing package declaration")
		}
		if !strings.Contains(code, "n = 1") {
			t.Error("GenGo() missing n constant")
		}
		if !strings.Contains(code, "alphabetSize = 26") {
			t.Error("GenGo() missing alphabetSize constant")
		}
		if !strings.Contains(code, "total = 100") {
			t.Error("GenGo() missing total constant")
		}
		if !strings.Contains(code, "var decode") {
			t.Error("GenGo() missing decode table")
		}
		if !strings.Contains(code, "var encode") {
			t.Error("GenGo() missing encode table")
		}
		if !strings.Contains(code, "var freqs") {
			t.Error("GenGo() missing freqs table")
		}
		if !strings.Contains(code, "[26]uint64") {
			t.Error("GenGo() missing uint64 type for raw table")
		}
	})

	t.Run("prob table", func(t *testing.T) {
		pt := NewProbTable(1, 100, a)
		pt.MustSetProbFromCount(50, 0)

		code := pt.GenGo("testpkg")

		if !strings.Contains(code, "[26]float64") {
			t.Error("GenGo() missing float64 type for prob table")
		}
	})

	t.Run("logprob table", func(t *testing.T) {
		lpt := NewLogProbTable(1, 100, a)
		lpt.MustSetLogProbFromCount(50, 0)

		code := lpt.GenGo("testpkg")

		if !strings.Contains(code, "[26]float64") {
			t.Error("GenGo() missing float64 type for logprob table")
		}
	})
}

func TestGenC(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("raw table", func(t *testing.T) {
		rt := NewRawTable(1, 100, a)
		rt.MustSet(50, 0)

		code := rt.GenC()

		// Check for expected content
		if !strings.Contains(code, "#include <stdint.h>") {
			t.Error("GenC() missing stdint.h include")
		}
		if !strings.Contains(code, "#define ALPHABET_SIZE 26") {
			t.Error("GenC() missing ALPHABET_SIZE define")
		}
		if !strings.Contains(code, "static const int8_t decode") {
			t.Error("GenC() missing decode table")
		}
		if !strings.Contains(code, "static const uint8_t encode") {
			t.Error("GenC() missing encode table")
		}
		if !strings.Contains(code, "uint64_t freqs") {
			t.Error("GenC() missing uint64_t freqs for raw table")
		}
	})

	t.Run("prob table", func(t *testing.T) {
		pt := NewProbTable(1, 100, a)
		pt.MustSetProbFromCount(50, 0)

		code := pt.GenC()

		if !strings.Contains(code, "double freqs") {
			t.Error("GenC() missing double freqs for prob table")
		}
	})

	t.Run("logprob table with infinity", func(t *testing.T) {
		lpt := NewLogProbTable(1, 100, a)
		// Don't set any values, so they'll be -Inf

		code := lpt.GenC()

		if !strings.Contains(code, "#include <math.h>") {
			t.Error("GenC() missing math.h include")
		}
		if !strings.Contains(code, "-INFINITY") {
			t.Error("GenC() missing -INFINITY for unset logprob values")
		}
	})
}

func TestTabWidth(t *testing.T) {
	tests := []struct {
		s        string
		expected int
	}{
		{"", 0},
		{"a", 1},
		{"abc", 3},
		{"\t", 8},
		{"a\t", 9},
		{"\t\t", 16},
	}

	for _, tt := range tests {
		got := tabWidth(tt.s)
		if got != tt.expected {
			t.Errorf("tabWidth(%q) = %d, want %d", tt.s, got, tt.expected)
		}
	}
}
