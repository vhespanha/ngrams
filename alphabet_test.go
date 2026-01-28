package ngrams

import (
	"testing"
)

func TestExpandRange(t *testing.T) {
	tests := []struct {
		name     string
		spec     string
		expected []rune
	}{
		{
			name:     "single character",
			spec:     "a",
			expected: []rune{'a'},
		},
		{
			name:     "simple range",
			spec:     "a-c",
			expected: []rune{'a', 'b', 'c'},
		},
		{
			name:     "mixed range and singles",
			spec:     "a-c0",
			expected: []rune{'a', 'b', 'c', '0'},
		},
		{
			name:     "multiple ranges",
			spec:     "a-c0-2",
			expected: []rune{'a', 'b', 'c', '0', '1', '2'},
		},
		{
			name:     "lowercase spec",
			spec:     LowercaseSpec,
			expected: []rune("abcdefghijklmnopqrstuvwxyz"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandRange(tt.spec)
			if len(result) != len(tt.expected) {
				t.Errorf("expandRange(%q) = %v (len %d), want %v (len %d)",
					tt.spec, result, len(result), tt.expected, len(tt.expected))
				return
			}
			for i, r := range result {
				if r != tt.expected[i] {
					t.Errorf("expandRange(%q)[%d] = %q, want %q",
						tt.spec, i, r, tt.expected[i])
				}
			}
		})
	}
}

func TestNewAlphabet(t *testing.T) {
	tests := []struct {
		name         string
		spec         string
		expectedSize int
	}{
		{
			name:         "lowercase",
			spec:         LowercaseSpec,
			expectedSize: 26,
		},
		{
			name:         "lowercase alphanum",
			spec:         LowercaseAlphanumSpec,
			expectedSize: 36,
		},
		{
			name:         "alphanumeric",
			spec:         AlphanumericSpec,
			expectedSize: 62,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newAlphabet(tt.spec)
			if a.size != tt.expectedSize {
				t.Errorf("newAlphabet(%q).size = %d, want %d",
					tt.spec, a.size, tt.expectedSize)
			}
			if a.spec != tt.spec {
				t.Errorf("newAlphabet(%q).spec = %q, want %q",
					tt.spec, a.spec, tt.spec)
			}
		})
	}
}

func TestAlphabetConstructors(t *testing.T) {
	t.Run("NewLowercaseAlphabet", func(t *testing.T) {
		a := NewLowercaseAlphabet()
		if a.size != 26 {
			t.Errorf("NewLowercaseAlphabet().size = %d, want 26", a.size)
		}
	})

	t.Run("NewLowercaseAlphanumAlphabet", func(t *testing.T) {
		a := NewLowercaseAlphanumAlphabet()
		if a.size != 36 {
			t.Errorf("NewLowercaseAlphanumAlphabet().size = %d, want 36", a.size)
		}
	})

	t.Run("NewAlphanumericAlphabet", func(t *testing.T) {
		a := NewAlphanumericAlphabet()
		if a.size != 62 {
			t.Errorf("NewAlphanumericAlphabet().size = %d, want 62", a.size)
		}
	})

	t.Run("NewPrintableAlphabet", func(t *testing.T) {
		a := NewPrintableAlphabet()
		if a.size != 82 {
			t.Errorf("NewPrintableAlphabet().size = %d, want 82", a.size)
		}
	})
}

func TestSymbolFromByte(t *testing.T) {
	a := NewLowercaseAlphabet()

	tests := []struct {
		name       string
		b          byte
		wantSymbol symbol
		wantOk     bool
	}{
		{
			name:       "valid a",
			b:          'a',
			wantSymbol: 0,
			wantOk:     true,
		},
		{
			name:       "valid z",
			b:          'z',
			wantSymbol: 25,
			wantOk:     true,
		},
		{
			name:       "invalid digit",
			b:          '0',
			wantSymbol: 0,
			wantOk:     false,
		},
		{
			name:       "invalid uppercase",
			b:          'A',
			wantSymbol: 0,
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ok := a.SymbolFromByte(tt.b)
			if ok != tt.wantOk {
				t.Errorf("SymbolFromByte(%q) ok = %v, want %v", tt.b, ok, tt.wantOk)
			}
			if ok && s != tt.wantSymbol {
				t.Errorf("SymbolFromByte(%q) symbol = %d, want %d", tt.b, s, tt.wantSymbol)
			}
		})
	}
}

func TestSymbolsFromString(t *testing.T) {
	a := NewLowercaseAlphabet()

	tests := []struct {
		name        string
		s           string
		wantSymbols []symbol
		wantErr     bool
	}{
		{
			name:        "valid abc",
			s:           "abc",
			wantSymbols: []symbol{0, 1, 2},
			wantErr:     false,
		},
		{
			name:        "valid xyz",
			s:           "xyz",
			wantSymbols: []symbol{23, 24, 25},
			wantErr:     false,
		},
		{
			name:        "empty string",
			s:           "",
			wantSymbols: []symbol{},
			wantErr:     false,
		},
		{
			name:        "invalid character",
			s:           "ab0",
			wantSymbols: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbols, err := a.SymbolsFromString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("SymbolsFromString(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(symbols) != len(tt.wantSymbols) {
					t.Errorf("SymbolsFromString(%q) = %v, want %v", tt.s, symbols, tt.wantSymbols)
					return
				}
				for i := range symbols {
					if symbols[i] != tt.wantSymbols[i] {
						t.Errorf("SymbolsFromString(%q)[%d] = %d, want %d",
							tt.s, i, symbols[i], tt.wantSymbols[i])
					}
				}
			}
		})
	}
}

func TestAlphabetEncodeDecode(t *testing.T) {
	a := NewLowercaseAlphabet()

	// Test that encoding and decoding are consistent
	for i := 0; i < a.size; i++ {
		encoded := a.encode[i]
		decoded := a.decode[encoded]
		if int(decoded) != i {
			t.Errorf("encode/decode mismatch: encode[%d]=%d, decode[%d]=%d, want %d",
				i, encoded, encoded, decoded, i)
		}
	}
}
