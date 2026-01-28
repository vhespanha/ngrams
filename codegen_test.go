package ngrams

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenGo(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("raw table generates valid Go", func(t *testing.T) {
		rt := NewRawTable(1, 100, a)
		rt.MustSet(50, 0)

		code := rt.GenGo("testpkg")

		// Verify the code is syntactically valid Go
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Errorf("GenGo() generated invalid Go code: %v\nCode:\n%s", err, code)
		}

		// Verify essential semantic elements exist
		if !strings.Contains(code, "package testpkg") {
			t.Error("GenGo() missing package declaration")
		}
		if !strings.Contains(code, "uint64") {
			t.Error("GenGo() missing uint64 type for raw table")
		}
	})

	t.Run("prob table generates valid Go", func(t *testing.T) {
		pt := NewProbTable(1, 100, a)
		pt.MustSetProbFromCount(50, 0)

		code := pt.GenGo("testpkg")

		// Verify the code is syntactically valid Go
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Errorf("GenGo() generated invalid Go code: %v\nCode:\n%s", err, code)
		}

		if !strings.Contains(code, "float64") {
			t.Error("GenGo() missing float64 type for prob table")
		}
	})

	t.Run("logprob table generates valid Go", func(t *testing.T) {
		lpt := NewLogProbTable(1, 100, a)
		lpt.MustSetLogProbFromCount(50, 0)

		code := lpt.GenGo("testpkg")

		// Verify the code is syntactically valid Go
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", code, 0)
		if err != nil {
			t.Errorf("GenGo() generated invalid Go code: %v\nCode:\n%s", err, code)
		}

		if !strings.Contains(code, "float64") {
			t.Error("GenGo() missing float64 type for logprob table")
		}
	})
}

func TestGenC(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("raw table generates C code", func(t *testing.T) {
		rt := NewRawTable(1, 100, a)
		rt.MustSet(50, 0)

		code := rt.GenC()

		// Verify C code contains required includes and declarations
		if !strings.Contains(code, "#include <stdint.h>") {
			t.Error("GenC() missing stdint.h include")
		}
		if !strings.Contains(code, "ALPHABET_SIZE") {
			t.Error("GenC() missing ALPHABET_SIZE define")
		}
		if !strings.Contains(code, "int8_t decode") {
			t.Error("GenC() missing decode table")
		}
		if !strings.Contains(code, "uint8_t encode") {
			t.Error("GenC() missing encode table")
		}
		if !strings.Contains(code, "uint64_t") {
			t.Error("GenC() missing uint64_t type for raw table")
		}
	})

	t.Run("prob table generates C code", func(t *testing.T) {
		pt := NewProbTable(1, 100, a)
		pt.MustSetProbFromCount(50, 0)

		code := pt.GenC()

		if !strings.Contains(code, "double") {
			t.Error("GenC() missing double type for prob table")
		}
	})

	t.Run("logprob table with infinity generates C code", func(t *testing.T) {
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
