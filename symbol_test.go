package ngrams_test

import (
	"testing"

	"github.com/vhespanha/ngrams"
)

func TestSymbolFromByte(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		want  bool
	}{
		{"valid byte a", byte('a'), true},
		{"valid byte z", byte('z'), true},
		{"invalid byte A", byte('A'), false},
		{"invalid byte 0", byte('0'), false},
		{"invalid byte Ü", byte('Ü'), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := ngrams.SymbolFromByte(tt.input)
			if ok != tt.want {
				t.Errorf("got %t, want %t", ok, tt.want)
			}
		})
	}
}

func TestSymbolsFromString(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		inputN      int
		wantErr     bool
	}{
		{"valid string 'foo'", "foo", len("foo"), false},
		{"valid string 'aaaaaaaaaaaaaaaa'", "aaaaaaaaaaaaaaaa", len("aaaaaaaaaaaaaaaa"), false},
		{"invalid string 'AAAAAAAAAAAAAAAA'", "AAAAAAAAAAAAAAAA", len("AAAAAAAAAAAAAAAA"), true},
		{"invalid string 'bar'", "bar", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ngrams.SymbolsFromString(tt.inputString, tt.inputN)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
