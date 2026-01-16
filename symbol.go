package ngrams

import "fmt"

type symbol uint8

func SymbolFromByte(b byte) (symbol, bool) {
	if b >= 'a' && b <= 'z' {
		return symbol(b - 'a'), true
	}
	return 0, false
}

func SymbolsFromString(s string, n int) ([]symbol, error) {
	symbols := make([]symbol, len(s))
	for i := 0; i < n; i++ {
		symbol, ok := SymbolFromByte(s[i])
		if !ok {
			return nil, fmt.Errorf("invalid symbol %q", s[i])
		}
		symbols[i] = symbol
	}
	return symbols, nil
}
